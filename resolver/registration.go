package resolver

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"
	"math/big"
	"net"

	"github.com/1inch/p2p-network/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	errInvalidFormatAddress = errors.New("invalid format for address")
	errInvalidFormatIp      = errors.New("invalid format for ip")
)

type RegistrationResolver struct {
	logger       *slog.Logger
	ethClient    *ethclient.Client
	nodeRegistry *contracts.NodeRegistry
}

func NewRegistrationResolver(logger *slog.Logger, ethUrl string, contractAddr string) (*RegistrationResolver, error) {
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		logger.Error("error when try create eth client with url", slog.Any("err", err.Error()))
		return nil, err
	}
	err = validateAddress(contractAddr)
	if err != nil {
		return nil, err
	}
	nodeRegistry, err := contracts.NewNodeRegistry(common.HexToAddress(contractAddr), client)
	if err != nil {
		logger.Error("error when create contract node registry")
		return nil, err
	}
	return &RegistrationResolver{
		logger:       logger,
		ethClient:    client,
		nodeRegistry: nodeRegistry,
	}, nil
}

// TODO Register maybe need added configuration limit gas price and gas limit,
func (r *RegistrationResolver) Register(ctx context.Context, account string, hexPrivKey string, ip string, encodedPublicKey string) (*common.Hash, error) {
	err := validateAddress(account)
	if err != nil {
		r.logger.Error("error when validate account address", slog.Any("err", err.Error()))
		return nil, err
	}

	err = validateIp(ip)
	if err != nil {
		r.logger.Error("error when validate ip", slog.Any("err", err.Error()))
		return nil, err
	}

	accountAddr := common.HexToAddress(account)
	nonce, err := r.ethClient.NonceAt(ctx, accountAddr, nil)
	if err != nil {
		r.logger.Error("error when get nonce by account", slog.Any("err", err.Error()))
		return nil, err
	}
	transactnOps := &bind.TransactOpts{
		From:  accountAddr,
		Nonce: big.NewInt(int64(nonce)),
		Signer: func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return r.funcForSign(hexPrivKey, tx)
		},
	}

	publicKey, err := decodePublicKey(encodedPublicKey)
	if err != nil {
		r.logger.Error("error when decode public key")
		return nil, err
	}

	tx, err := r.nodeRegistry.RegisterResolver(transactnOps, ip, publicKey)
	if err != nil {
		r.logger.Error("error when call contract method 'RegisterResolver'", slog.Any("err", err.Error()))
		return nil, err
	}

	txHash := tx.Hash()

	return &txHash, nil
}

func validateIp(ip string) error {
	netIp := net.ParseIP(ip)
	if netIp.To4() == nil {
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

func decodePublicKey(encodedPublicKey string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedPublicKey)
}

func (r *RegistrationResolver) funcForSign(hexPrivKey string, tx *types.Transaction) (*types.Transaction, error) {
	chainId, err := r.ethClient.ChainID(context.Background())
	if err != nil {
		r.logger.Error("error when try get chain id", slog.Any("err", err.Error()))
		return nil, err
	}
	privKey, _ := crypto.HexToECDSA(hexPrivKey)
	signer := types.NewLondonSigner(chainId)
	trans, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		return nil, err
	}

	return trans, nil
}
