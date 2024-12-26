package encryption

import (
	"crypto"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EncryptionTestSuite struct {
	suite.Suite

	relayerSecp256k1     crypto.PrivateKey
	resolverSecp256k1    crypto.PrivateKey
	relayerSecp256k1PEM  []byte
	resolverSecp256k1PEM []byte

	relayerRsa     crypto.PrivateKey
	resolverRsa    crypto.PrivateKey
	relayerRsaPEM  []byte
	resolverRsaPEM []byte
}

func (s *EncryptionTestSuite) SetupTest() {
	// Secp256k1 keys
	relayerSecp256k1, err := GenerateKeyPair(Secp256k1)
	s.Require().NoError(err)
	s.relayerSecp256k1 = relayerSecp256k1

	relayerSecp256k1PEM, err := ToPEM(relayerSecp256k1)
	s.Require().NoError(err)
	s.relayerSecp256k1PEM = relayerSecp256k1PEM

	resolverSecp256k1, err := GenerateKeyPair(Secp256k1)
	s.Require().NoError(err)
	s.resolverSecp256k1 = resolverSecp256k1

	resolverSecp256k1PEM, err := ToPEM(resolverSecp256k1)
	s.Require().NoError(err)
	s.resolverSecp256k1PEM = resolverSecp256k1PEM

	// RSA keys
	relayerRsa, err := GenerateKeyPair(RSA4096)
	s.Require().NoError(err)
	s.relayerRsa = relayerRsa

	relayerRsaPEM, err := ToPEM(relayerRsa)
	s.Require().NoError(err)
	s.relayerRsaPEM = relayerRsaPEM

	resolverRsa, err := GenerateKeyPair(RSA4096)
	s.Require().NoError(err)
	s.resolverRsa = resolverRsa

	resolverRsaPEM, err := ToPEM(resolverRsa)
	s.Require().NoError(err)
	s.resolverRsaPEM = resolverRsaPEM
}

func (s *EncryptionTestSuite) TearDownTest() {
}

func TestEncryptionTestSuite(t *testing.T) {
	suite.Run(t, new(EncryptionTestSuite))
}

func (s *EncryptionTestSuite) TestEncryption() {
	payload := []byte("test payload")

	// Secp256k1
	encrypted, err := Encrypt(payload, s.resolverSecp256k1PEM)

	s.Require().NoError(err)

	decrypted, err := Decrypt(encrypted, s.resolverSecp256k1)
	s.Require().NoError(err)

	encrypted, err = Encrypt(decrypted, s.relayerSecp256k1PEM)
	s.Require().NoError(err)

	decrypted, err = Decrypt(encrypted, s.relayerSecp256k1)
	s.Require().NoError(err)
	s.Require().Equal(payload, decrypted)

	// RSA
	encrypted, err = Encrypt(payload, s.resolverRsaPEM)
	s.Require().NoError(err)

	decrypted, err = Decrypt(encrypted, s.resolverRsa)
	s.Require().NoError(err)

	encrypted, err = Encrypt(decrypted, s.relayerRsaPEM)
	s.Require().NoError(err)

	decrypted, err = Decrypt(encrypted, s.relayerRsa)
	s.Require().NoError(err)
	s.Require().Equal(payload, decrypted)
}
