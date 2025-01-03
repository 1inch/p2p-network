package resolver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"os"
	"testing"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/resolver/types"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const defaultBalance = 555.

type ResolverTestSuite struct {
	suite.Suite

	logger *slog.Logger
	server *grpc.Server
	client pb.ExecuteClient
	conn   *grpc.ClientConn
}

func (s *ResolverTestSuite) SetupTest() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s.logger = logger
	listener := bufconn.Listen(1024 * 1024)

	server, err := newServer(logger, &TestApiHandler{})
	if err != nil {
		logger.Error("newServer failed", "error", err)
		return
	}

	grpcServer := setupRpcServer(listener, server)
	logger.Info("### Server started")
	s.server = grpcServer

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(
		func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Info("### error", "err", err)
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

func (h *TestApiHandler) Process(req *types.JsonRequest) (*types.JsonResponse, error) {
	switch req.Method {
	case "GetWalletBalance":
		{
			if len(req.Params) < 2 {
				return &types.JsonResponse{Id: req.Id, Result: 0}, errWrongParamCount
			}
			return &types.JsonResponse{Id: req.Id, Result: defaultBalance}, nil
		}
	default:
		return &types.JsonResponse{Id: req.Id, Result: 0}, errUnrecognizedMethod
	}
}
func (s *ResolverTestSuite) getWalletBalancePayloadOk() []byte {
	jsonReq := &types.JsonRequest{Id: "1", Method: "GetWalletBalance", Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"}}
	byteArr, _ := json.Marshal(jsonReq)

	return byteArr
}

func (s *ResolverTestSuite) getWalletBalancePayloadUnrecognizedMethod() []byte {
	jsonReq := &types.JsonRequest{Id: "1", Method: "UnrecognizedMethod", Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"}}
	byteArr, _ := json.Marshal(jsonReq)

	return byteArr
}

func (s *ResolverTestSuite) getWalletBalancePayloadNoParams() []byte {
	jsonReq := &types.JsonRequest{Id: "1", Method: "GetWalletBalance", Params: []string{}}
	byteArr, _ := json.Marshal(jsonReq)

	return byteArr
}

func (s *ResolverTestSuite) TestExecutePositive() {
	req := &pb.ResolverRequest{Id: "1", Payload: s.getWalletBalancePayloadOk(), Encrypted: false}

	s.logger.Info("###about to execute")
	resp, err := s.client.Execute(context.Background(), req)
	if err != nil {
		s.logger.Info("Execute error", "error", err)
		s.Require().Fail("execute error")
	}
	s.logger.Info("test output", "resp", resp)
	var jsonResp types.JsonResponse
	err = json.Unmarshal(resp.Payload, &jsonResp)
	s.Require().NoError(err)
	s.Require().Equal(jsonResp.Id, req.Id)
	s.Require().Equal(jsonResp.Result, defaultBalance)
}

type negativeTestCase struct {
	Name            string
	ResolverRequest *pb.ResolverRequest
	ExpectedCode    codes.Code
	ExpectedError   error
}

// i use this approach because negative tests looks like copy-paste with change in the expected data
func (s *ResolverTestSuite) TestExecuteNegativeCases() {
	testCases := []negativeTestCase{
		{
			Name:            "UnrecognizedMethodParameter",
			ResolverRequest: &pb.ResolverRequest{Id: "1", Payload: s.getWalletBalancePayloadUnrecognizedMethod(), Encrypted: false},
			ExpectedCode:    codes.InvalidArgument,
			ExpectedError:   errUnrecognizedMethod,
		},
		{
			Name:            "NoParameterInPayload",
			ResolverRequest: &pb.ResolverRequest{Id: "1", Payload: s.getWalletBalancePayloadNoParams(), Encrypted: false},
			ExpectedCode:    codes.InvalidArgument,
			ExpectedError:   errWrongParamCount,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			s.logger.Info("call execute client method")
			resp, err := s.client.Execute(context.Background(), testCase.ResolverRequest)
			s.Require().NotNil(err, "expect error from client")
			s.Require().Nil(resp, "expect response from client is nil")

			statusErr := status.Convert(err)

			resolverResponse := s.getResolverResponse(statusErr)
			s.Require().NotNil(resolverResponse, "expect resolver response be a first position in status error detais")
			s.Require().Equal(testCase.ResolverRequest.Id, resolverResponse.Id, "expect that id in request and response is equal")
			s.Require().Equal(testCase.ExpectedCode, statusErr.Code())
			s.Require().Equal(testCase.ExpectedError.Error(), statusErr.Message())
		})
	}
}

func (s *ResolverTestSuite) TestInfuraEndpoint() {
	client, err := gethrpc.DialHTTP("https://mainnet.infura.io/v3/a8401733346d412389d762b5a63b0bcf")
	s.Require().NoError(err)
	s.Require().NotNil(client)

	var result string
	err = client.Call(&result, "eth_getBalance", "0x38308C349fd2F9dad31Aa3bFe28015dA3EB67193", "latest")
	s.Require().NoError(err)

	s.Require().NotEmpty(result)
	s.logger.Info("Infura result", "res", result)
}

// parse status error from client
func (s *ResolverTestSuite) getResolverResponse(statusError *status.Status) *pb.ResolverResponse {
	resolverResponse, ok := statusError.Details()[0].(*pb.ResolverResponse)

	if !ok {
		return nil
	}

	return resolverResponse
}
