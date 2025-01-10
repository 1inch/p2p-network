// Package encryption provides basic functionality for payload encryption
package encryption

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log/slog"

	ecies "github.com/ecies/go/v2"
)

// Encryption type constants
const (
	Secp256k1 string = "ECDSA PUBLIC KEY"
	RSA4096   string = "RSA PUBLIC KEY"
)

// ErrUnknownType is an error type for unknown encryption type string
var ErrUnknownType = errors.New("Unknown encryption type")

// GenerateKeyPair generates asymmetric key pair given the encryption type
func GenerateKeyPair(keyType string) (crypto.PrivateKey, error) {
	switch keyType {
	case Secp256k1:
		return ecies.GenerateKey()
	case RSA4096:
		return rsa.GenerateKey(rand.Reader, 4096)
	default:
		return nil, ErrUnknownType
	}
}

// Encrypt encrypts payload using public key PEM
func Encrypt(payload []byte, pemBytes []byte) ([]byte, error) {
	pub, err := FromPEM(pemBytes)

	if err != nil {
		return nil, err
	}
	switch v := pub.(type) {
	case *ecies.PublicKey:
		slog.Info("### encrypt ecies")
		return ecies.Encrypt(v, payload)
	case *rsa.PublicKey:
		return rsa.EncryptOAEP(sha256.New(), rand.Reader, v, payload, nil)
	default:
		return nil, nil
	}
}

// FromPEM extracts public key from PEM
func FromPEM(pemBytes []byte) (crypto.PublicKey, error) {
	decodedKey, _ := pem.Decode(pemBytes)
	slog.Info("### FromPEM", decodedKey.Type)
	switch decodedKey.Type {
	case Secp256k1:
		slog.Info("Recover ecies key")
		return ecies.NewPublicKeyFromBytes(decodedKey.Bytes)
	case RSA4096:
		return x509.ParsePKIXPublicKey(decodedKey.Bytes)
	default:
		return nil, ErrUnknownType
	}
}

// ToPEM produces PEM with public given corresponding private key
func ToPEM(priv crypto.PrivateKey) ([]byte, error) {
	var keyType string
	var keyBytes []byte
	switch v := priv.(type) {
	case *ecies.PrivateKey:
		keyType = Secp256k1
		keyBytes = v.PublicKey.Bytes(false)
	case *rsa.PrivateKey:
		keyType = RSA4096
		keyBytes, _ = x509.MarshalPKIXPublicKey(&v.PublicKey)
	}
	var buf bytes.Buffer
	b := &pem.Block{Type: keyType, Bytes: keyBytes}
	err := pem.Encode(&buf, b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decrypt decrypts a payload using provided private key
func Decrypt(payload []byte, pk crypto.PrivateKey) ([]byte, error) {
	switch v := pk.(type) {
	case *ecies.PrivateKey:
		return ecies.Decrypt(v, payload)
	case *rsa.PrivateKey:
		return rsa.DecryptOAEP(sha256.New(), nil, v, payload, nil)
	default:
		return nil, nil
	}
}
