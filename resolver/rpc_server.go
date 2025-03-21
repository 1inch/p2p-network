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
	pb "github.com/1inch/p2p-network/proto/resolver"
	"github.com/1inch/p2p-network/resolver/types"
	ecies "github.com/ecies/go/v2"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

var (
	errEmptyRequest   = errors.New("empty request")
	errEmptyRequestId = errors.New("empty request id")
	errEmptyPayload   = errors.New("empty payload")
	errEmptyPublicKey = errors.New("empty public key")
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

	if err != nil {
		return s.buildResolverResponseWithErr(req, err), nil
	}

	jsonReq, err := s.getJsonRequest(req)
	if err != nil {
		return s.buildResolverResponseWithErr(req, err), nil
	}

	resp, err := s.processRequest(jsonReq)

	if err != nil {
		return s.buildResolverResponseWithErr(req, err), nil
	}

	if req.Encrypted {
		pubKeyDecompressed, err := ethCrypto.DecompressPubkey(req.PublicKey)
		if err != nil {
			return s.buildResolverResponseWithErr(req, err), nil
		}
		pubKeyBytes := ethCrypto.FromECDSAPub(pubKeyDecompressed)

		pubKey, err := ecies.NewPublicKeyFromBytes(pubKeyBytes)
		if err != nil {
			return s.buildResolverResponseWithErr(req, err), nil
		}

		resp, err = encryption.Encrypt(resp, pubKey)
		if err != nil {
			return s.buildResolverResponseWithErr(req, err), nil
		}
	}
	return &pb.ResolverResponse{
		Id:        req.Id,
		Encrypted: req.Encrypted,
		Result: &pb.ResolverResponse_Payload{
			Payload: resp,
		},
	}, nil
}

func (s *Server) validateResolverRequest(req *pb.ResolverRequest) error {
	// return after this check, because maybe nil pointer exception in next checks
	// maybe this check is useless. let it stay for reinsuranceClick to apply
	if req == nil {
		return errEmptyRequest
	}

	if req.Id == "" {
		return errEmptyRequestId
	}

	if len(req.Payload) == 0 {
		return errEmptyPayload
	}

	if req.Encrypted {
		if len(req.PublicKey) == 0 {
			return errEmptyPublicKey
		}
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

func (s *Server) buildResolverResponseWithErr(req *pb.ResolverRequest, err error) *pb.ResolverResponse {
	return &pb.ResolverResponse{
		Id: req.Id,
		Result: &pb.ResolverResponse_Error{
			Error: &pb.Error{
				Code:    s.getErrorCodeByErr(err),
				Message: err.Error(),
			},
		},
	}
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

func (s *Server) getErrorCodeByErr(err error) pb.ErrorCode {
	if errors.Is(err, errEmptyRequest) ||
		errors.Is(err, errEmptyRequestId) ||
		errors.Is(err, errEmptyPayload) ||
		errors.Is(err, errEmptyPublicKey) ||
		errors.Is(err, errWrongParamCount) ||
		errors.Is(err, errInvalidFormatAddress) ||
		errors.Is(err, errUnrecognizedMethod) {

		return pb.ErrorCode_ERR_INVALID_MESSAGE_FORMAT
	}

	return pb.ErrorCode_ERR_GRPC_EXECUTION_FAILED
}
