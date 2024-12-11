package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"testing"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/resolver/types"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const defaultBalance = 555

type ResolverTestSuite struct {
	suite.Suite

	server *grpc.Server
	client pb.ExecuteClient
	conn   *grpc.ClientConn
}

func (s *ResolverTestSuite) SetupTest() {
	listener := bufconn.Listen(1024 * 1024)

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	server, err := newServer(&TestApiHandler{})
	if err != nil {
		slog.Error("newServer failed", "error", err)
		return
	}

	pb.RegisterExecuteServer(grpcServer, server)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			s.Require().Fail("Failed to start grpc server")
			return
		}
	}()
	slog.Info("### Server started")
	s.server = grpcServer

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(
		func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Info("### error", "err", err)
	}

	s.conn = conn
	s.client = pb.NewExecuteClient(conn)
}

func (s *ResolverTestSuite) TearDownTest() {
	s.server.GracefulStop()
	s.conn.Close()
}

func TestResolverTestSuite(t *testing.T) {
	suite.Run(t, new(ResolverTestSuite))
}

type TestApiHandler struct{}

func (h *TestApiHandler) Process(req *types.JsonRequest) *types.JsonResponse {
	switch req.Method {
	case "GetWalletBalance":
		return &types.JsonResponse{Id: req.Id, Result: defaultBalance, Error: nil}
	default:
		return &types.JsonResponse{Id: req.Id, Result: 0, Error: errors.New("Unrecognized method")}
	}
}
func (s *ResolverTestSuite) getWalletBalancePayload() []byte {
	jsonReq := &types.JsonRequest{Id: "1", Method: "GetWalletBalance", Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"}}
	byteArr, _ := json.Marshal(jsonReq)

	return byteArr
}

func (s *ResolverTestSuite) TestExecute() {

	req := &pb.ResolverRequest{Id: "1", Payload: s.getWalletBalancePayload(), Encrypted: false}

	slog.Info("###about to execute")
	resp, err := s.client.Execute(context.Background(), req)
	if err != nil {
		slog.Info("Execute error", "error", err)
		s.Require().Fail("execute error")
	}
	slog.Info("test output", "resp", resp)
	var jsonResp types.JsonResponse
	err = json.Unmarshal(resp.Payload, &jsonResp)
	s.Require().NoError(err)
	s.Require().Equal(jsonResp.Id, req.Id)
	s.Require().Equal(int(jsonResp.Result.(float64)), defaultBalance)
}
