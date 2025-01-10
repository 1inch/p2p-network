// Package resolver implements the gRPC server
package resolver

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	"github.com/1inch/p2p-network/internal/encryption"
	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/resolver/types"
	ecies "github.com/ecies/go/v2"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errEmptyParam = errors.New("empty parameter")
)

// Server represents gRPC server.
type Server struct {
	pb.UnimplementedExecuteServer

	privateKey *ecies.PrivateKey

	logger *slog.Logger

	handler ApiHandler
}

// newServer creates new RpcServer.
func newServer(cfg *Config) (*Server, error) {
	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	})
	logger := slog.New(loggerHandler)
	var handler ApiHandler

	switch {
	case cfg.Apis.Default.Enabled:
		{
			handler = NewDefaultApiHandler(cfg.Apis.Default, logger)
		}
	case cfg.Apis.Infura.Enabled:
		{
			handler = NewInfuraApiHandler(cfg.Apis.Infura, logger)
		}
	default:
		logger.Error("expect someone handler api in config")
		return nil, errNoHandlerApiInConfig
	}

	var privKey *ecies.PrivateKey
	if len(cfg.PrivateKey) > 0 {
		privKeyBytes, err := hex.DecodeString(cfg.PrivateKey)
		if err != nil {
			logger.Error("error decoding private key hex", slog.Any("err", err))
			return nil, err
		}
		privKey = ecies.NewPrivateKeyFromBytes(privKeyBytes)
	} else {
		privKeyGenerated, err := encryption.GenerateKeyPair()

		if err != nil {
			return nil, err
		}
		privKey = privKeyGenerated
	}
	return &Server{privateKey: privKey, logger: logger.With("module", "rpc-server"), handler: handler}, nil
}

// Execute executes ResolverRequest.
func (s *Server) Execute(ctx context.Context, req *pb.ResolverRequest) (*pb.ResolverResponse, error) {
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

	if req.Encrypted {
		pubKeyDecompressed, err := ethCrypto.DecompressPubkey(req.PublicKey)
		if err != nil {
			return nil, s.formatGrpcError(response, err)
		}
		pubKeyBytes := ethCrypto.FromECDSAPub(pubKeyDecompressed)

		pubKey, err := ecies.NewPublicKeyFromBytes(pubKeyBytes)
		if err != nil {
			return nil, s.formatGrpcError(response, err)
		}

		resp, err = encryption.Encrypt(resp, pubKey)
		if err != nil {
			return nil, s.formatGrpcError(response, err)
		}
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
	var payload []byte
	if req.Encrypted {
		decrypted, err := encryption.Decrypt(req.Payload, s.privateKey)
		if err != nil {
			return nil, err
		}
		payload = decrypted
	} else {
		payload = req.Payload
	}
	err := json.Unmarshal(payload, &jsonReq)
	if err != nil {
		s.logger.Error("failed unmarshal request payload")
		return nil, err
	}
	return &jsonReq, nil
}

func (s *Server) processRequest(jsonReq *types.JsonRequest) ([]byte, error) {
	jsonResponses, err := s.handler.Process(jsonReq)
	if err != nil {
		s.logger.Error("failed process request in handler")
		return nil, err
	}

	byteArr, err := json.Marshal(jsonResponses)
	if err != nil {
		s.logger.Error("failed marshal json responses ")
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
		return status.New(codes.Internal, "failed format err status").Err()
	}
	return errStatus.Err()
}

func (s *Server) itsKnownError(err error) bool {
	return errors.Is(err, errEmptyBlock) ||
		errors.Is(err, errEmptyAddress) ||
		errors.Is(err, errUnrecognizedMethod) ||
		errors.Is(err, errWrongParamCount)
}
