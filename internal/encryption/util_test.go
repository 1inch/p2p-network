package encryption

import (
	"testing"

	ecies "github.com/ecies/go/v2"
	"github.com/stretchr/testify/suite"
)

type EncryptionTestSuite struct {
	suite.Suite

	relayer  *ecies.PrivateKey
	resolver *ecies.PrivateKey
}

func (s *EncryptionTestSuite) SetupTest() {
	//  keys
	relayer, err := GenerateKeyPair()
	s.Require().NoError(err)
	s.relayer = relayer

	resolver, err := GenerateKeyPair()
	s.Require().NoError(err)
	s.resolver = resolver
}

func (s *EncryptionTestSuite) TearDownTest() {
}

func TestEncryptionTestSuite(t *testing.T) {
	suite.Run(t, new(EncryptionTestSuite))
}

func (s *EncryptionTestSuite) TestEncryption() {
	payload := []byte("test payload")

	encrypted, err := Encrypt(payload, s.resolver.PublicKey)

	s.Require().NoError(err)

	decrypted, err := Decrypt(encrypted, s.resolver)
	s.Require().NoError(err)

	encrypted, err = Encrypt(decrypted, s.relayer.PublicKey)
	s.Require().NoError(err)

	decrypted, err = Decrypt(encrypted, s.relayer)
	s.Require().NoError(err)
	s.Require().Equal(payload, decrypted)
}
