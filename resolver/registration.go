package resolver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/1inch/p2p-network/internal/registry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	errInvalidFormatAddress = errors.New("invalid format for address")
	errInvalidFormatIp      = errors.New("invalid format for ip")
)

// RegistrationResolver describe registration new resolver on blockchain registry
type RegistrationResolver struct {
	logger         *slog.Logger
	cfg            Config
	registryClient *registry.Client
}

// NewRegistrationResolver create new instans RegistrationResolver
func NewRegistrationResolver(logger *slog.Logger, cfg *Config) (*RegistrationResolver, error) {
	err := validateEndpoint(cfg.GrpcEndpoint)
	if err != nil {
		logger.Error("error when validate ip", slog.Any("err", err.Error()))
		return nil, err
	}

	err = validateAddress(cfg.ContractAddress)
	if err != nil {
		return nil, err
	}

	rawUrl := fmt.Sprintf("http://%s", cfg.RpcUrl)
	registryCfg := &registry.Config{
		DialURI:         rawUrl,
		PrivateKey:      cfg.PrivateKey,
		ContractAddress: cfg.ContractAddress,
	}
	registryCli, err := registry.Dial(context.Background(), registryCfg)
	if err != nil {
		logger.Error("error when try create registry client", slog.Any("err", err.Error()))
		return nil, err
	}

	return &RegistrationResolver{
		logger:         logger,
		cfg:            *cfg,
		registryClient: registryCli,
	}, nil
}

// Register workflow for registration resolver to blockchain registry
// TODO Register maybe need added configuration limit gas price and gas limit,
func (r *RegistrationResolver) Register(ctx context.Context) (*common.Hash, error) {
	privateKey, err := crypto.HexToECDSA(r.cfg.PrivateKey)
	if err != nil {
		r.logger.Error("failed map hex to private key")
		return nil, err
	}
	publicKey := crypto.CompressPubkey(&privateKey.PublicKey)

	tx, err := r.registryClient.Registry.RegisterResolver(r.registryClient.Auth, r.cfg.GrpcEndpoint, publicKey)
	if err != nil {
		r.logger.Error("failed call contract method 'RegisterResolver'", slog.Any("err", err.Error()))
		return nil, err
	}

	txHash := tx.Hash()

	err = r.registryClient.WaitForTx(ctx, txHash)
	if err != nil {
		r.logger.Error("failed process transaction for register resolver", slog.Any("err", err.Error()))
		return nil, err
	}

	return &txHash, nil
}

func validateEndpoint(endpoint string) error {
	_, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return errInvalidFormatIp
	}
	return nil
}

func validateAddress(address string) error {
	if !common.IsHexAddress(address) {
		return errInvalidFormatAddress
	}
	return nil
}
