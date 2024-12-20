package resolver

import (
	"testing"

	"github.com/1inch/p2p-network/internal/configs"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestYamlConfigParsing() {
	cfg, err := configs.LoadConfig[Config]("resolver.config.example.yaml")
	s.Require().NoError(err)

	s.Require().Equal(cfg.Port, 8001)
	s.Require().Equal(cfg.Apis.Default.Enabled, true)
	s.Require().Equal(cfg.Apis.Infura.Enabled, true)
	s.Require().Equal(cfg.Apis.Infura.Key, "test-key")
}
