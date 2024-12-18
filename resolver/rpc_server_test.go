package resolver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"testing"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/resolver/types"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const defaultBalance = 555
const testApiName = "testApi"

type ResolverTestSuite struct {
	suite.Suite

	server *grpc.Server
	client pb.ExecuteClient
	conn   *grpc.ClientConn
}

func (s *ResolverTestSuite) SetupTest() {
	listener := bufconn.Listen(1024 * 1024)

	server, err := newServer([]ApiHandler{&TestApiHandler{}})
	if err != nil {
		slog.Error("newServer failed", "error", err)
		return
	}

	grpcServer := setupRpcServer(listener, server)
	slog.Info("### Server started")
	s.server = grpcServer

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(
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
	s.Require().NoError(s.conn.Close())
}

func TestResolverTestSuite(t *testing.T) {
	suite.Run(t, new(ResolverTestSuite))
}

type TestApiHandler struct{}

func (h *TestApiHandler) Name() string {
	return "testApi"
}

func (h *TestApiHandler) Process(req *types.JsonRequest) *types.JsonResponse {
	switch req.Method {
	case "GetWalletBalance":
		if len(req.Params[0]) > 0 {
			return &types.JsonResponse{Id: req.Id, Result: defaultBalance}
		} else {
			return &types.JsonResponse{Id: req.Id, Result: 0, Error: "Missing address"}
		}
	default:
		return &types.JsonResponse{Id: req.Id, Result: 0, Error: "Unrecognized method"}
	}
}
func (s *ResolverTestSuite) getWalletBalancePayloadOk() []byte {
	jsonReq := &types.JsonRequest{Id: "1", Method: "GetWalletBalance", Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"}}
	byteArr, _ := json.Marshal(jsonReq)

	return byteArr
}

func (s *ResolverTestSuite) getWalletBalancePayloadMissingMethod() []byte {
	jsonReq := &types.JsonRequest{Id: "1", Method: "MissingMethod", Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"}}
	byteArr, _ := json.Marshal(jsonReq)

	return byteArr
}

func (s *ResolverTestSuite) getWalletBalancePayloadWrongParams() []byte {
	jsonReq := &types.JsonRequest{Id: "1", Method: "GetWalletBalance", Params: []string{"", "latest"}}
	byteArr, _ := json.Marshal(jsonReq)

	return byteArr
}

func (s *ResolverTestSuite) TestExecute() {
	req := &pb.ResolverRequest{Id: "1", Payload: s.getWalletBalancePayloadOk(), Encrypted: false}

	slog.Info("###about to execute")
	resp, err := s.client.Execute(context.Background(), req)
	if err != nil {
		slog.Info("Execute error", "error", err)
		s.Require().Fail("execute error")
	}
	slog.Info("test output", "resp", resp)
	var jsonResp types.JsonResponses
	err = json.Unmarshal(resp.Payload, &jsonResp)
	s.Require().NoError(err)
	s.Require().Equal(jsonResp[testApiName].Id, req.Id)
	s.Require().Equal(int(jsonResp[testApiName].Result.(float64)), defaultBalance)
}

func (s *ResolverTestSuite) TestExecuteMissingMethod() {
	req := &pb.ResolverRequest{Id: "1", Payload: s.getWalletBalancePayloadMissingMethod(), Encrypted: false}

	slog.Info("###about to execute")
	resp, err := s.client.Execute(context.Background(), req)
	if err != nil {
		slog.Info("Execute error", "error", err)
		s.Require().Fail("execute error")
	}
	slog.Info("test output", "resp", resp)
	var jsonResp types.JsonResponses
	err = json.Unmarshal(resp.Payload, &jsonResp)
	s.Require().NoError(err)
	s.Require().Equal(jsonResp[testApiName].Id, req.Id)
	s.Require().Equal("Unrecognized method", jsonResp[testApiName].Error)
}

func (s *ResolverTestSuite) TestExecuteMissingAddress() {
	req := &pb.ResolverRequest{Id: "1", Payload: s.getWalletBalancePayloadWrongParams(), Encrypted: false}

	slog.Info("###about to execute")
	resp, err := s.client.Execute(context.Background(), req)
	if err != nil {
		slog.Info("Execute error", "error", err)
		s.Require().Fail("execute error")
	}
	slog.Info("test output", "resp", resp)
	var jsonResp types.JsonResponses
	err = json.Unmarshal(resp.Payload, &jsonResp)
	s.Require().NoError(err)
	s.Require().Equal(jsonResp[testApiName].Id, req.Id)
	s.Require().Equal("Missing address", jsonResp[testApiName].Error)
}

func (s *ResolverTestSuite) TestInfuraEndpoint() {
	client, err := gethrpc.DialHTTP("https://mainnet.infura.io/v3/a8401733346d412389d762b5a63b0bcf")
	s.Require().NoError(err)
	s.Require().NotNil(client)

	var result string
	err = client.Call(&result, "eth_getBalance", "0x38308C349fd2F9dad31Aa3bFe28015dA3EB67193", "latest")
	s.Require().NoError(err)

	s.Require().NotEmpty(result)
	slog.Info("Infura result", "res", result)
}
