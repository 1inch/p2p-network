// Package resolver implements the gRPC server
package resolver

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"log/slog"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/resolver/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	errMessageForInvalidArgumens     = "error when process request"
	fmtDescriptionErrorUnknownMethod = "unknown handle method: %s"
)

var (
	errWrongRequest = errors.New("wrong request body")
	errEmptyParam   = errors.New("empty parameter")
)

// Server represents gRPC server.
type Server struct {
	pb.UnimplementedExecuteServer

	privateKey *rsa.PrivateKey

	logger *slog.Logger

	handler ApiHandler
}

func generateKey() (*rsa.PrivateKey, error) {
	p, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// newServer creates new RpcServer.
func newServer(logger *slog.Logger, handler ApiHandler) (*Server, error) {
	privKey, err := generateKey()
	if err != nil {
		return nil, err
	}
	return &Server{privateKey: privKey, logger: logger.With("module", "server"), handler: handler}, nil
}

// Execute executes ResolverRequest.
func (s *Server) Execute(ctx context.Context, req *pb.ResolverRequest) (*pb.ResolverResponse, error) {
	s.logger.Info("###Incoming request", "id", req.Id)

	err := s.validateResolverRequest(req)
	response := &pb.ResolverResponse{
		Id: req.Id,
	}

	if err != nil {
		return nil, s.formatGrpcError(response, err)
	}

	jsonReq, err := s.getJsonRequest(req)
	if err != nil {
		return nil, s.formatGrpcError(response, err)
	}

	resp, err := s.processRequest(jsonReq)

	if err != nil {
		return nil, s.formatGrpcError(response, err)
	}

	response.Payload = resp

	return response, nil
}

func (s *Server) validateResolverRequest(req *pb.ResolverRequest) error {
	// return after this check, because maybe nil pointer exception in next checks
	if req == nil {
		return errEmptyParam
	}

	if req.Id == "" {
		return errEmptyParam
	}

	if len(req.Payload) == 0 {
		return errEmptyParam
	}

	return nil
}

func (s *Server) getJsonRequest(req *pb.ResolverRequest) (*types.JsonRequest, error) {
	var jsonReq types.JsonRequest
	err := json.Unmarshal(req.Payload, &jsonReq)
	if err != nil {
		s.logger.Error("error when try unmarshal request payload")
		return nil, err
	}
	return &jsonReq, nil
}

func (s *Server) processRequest(jsonReq *types.JsonRequest) ([]byte, error) {
	jsonResponses, err := s.handler.Process(jsonReq)
	if err != nil {
		s.logger.Error("error when process request in handler")
		return nil, err
	}

	byteArr, err := json.Marshal(jsonResponses)
	if err != nil {
		s.logger.Error("error when try marshal json responses ")
		return nil, err
	}

	return byteArr, nil
}

// formatGrpcError function describe create grpc error and map to go error
// pb.ResolverResponse needed for add it in details. In the response stored request id.
func (s *Server) formatGrpcError(resp *pb.ResolverResponse, err error) error {
	var code codes.Code
	switch {
	case s.itsKnownError(err):
		{
			code = codes.InvalidArgument
		}
	default:
		{
			code = codes.Internal
		}
	}

	errStatus, err := status.New(code, err.Error()).WithDetails(resp)
	if err != nil {
		return status.New(codes.Internal, "error when try format err status").Err()
	}
	return errStatus.Err()
}

func (s *Server) itsKnownError(err error) bool {
	return errors.Is(err, errEmptyBlock) ||
		errors.Is(err, errEmptyAddress) ||
		errors.Is(err, errUnrecognizedMethod) ||
		errors.Is(err, errWrongParamCount)
}
