package rpc

import (
	"context"
	"log/slog"
	"testing"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ResolverTestSuite struct {
	suite.Suite

	server *grpc.Server
}

func (s *ResolverTestSuite) SetupTest() {
	server, err := Start(&Config{Port: 8001, LogLevel: slog.LevelInfo, Testing: false})
	if err != nil {
		slog.Info("###Failed to start server", "error", err)
		return
	}
	slog.Info("### Server started")
	s.server = server
}

func (s *ResolverTestSuite) TearDownTest() {
	s.server.GracefulStop()
}

func TestResolverTestSuite(t *testing.T) {
	suite.Run(t, new(ResolverTestSuite))
}

func (s *ResolverTestSuite) TestExecute() {
	conn, err := grpc.NewClient("localhost:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Info("### error", "err", err)
	}

	client := pb.NewExecuteClient(conn)
	s.Require().NotNil(client)
	req := &pb.ResolverRequest{Id: "1", Payload: []byte("test"), Encrypted: false}

	slog.Info("###about to execute")
	var callOpts []grpc.CallOption = make([]grpc.CallOption, 0)
	resp, err := client.Execute(context.Background(), req, callOpts...)
	slog.Info("test output", "resp", resp)
	s.Require().Equal(string(resp.Payload), "test")
	defer conn.Close()
}
