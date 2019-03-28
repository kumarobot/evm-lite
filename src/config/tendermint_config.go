package config

import (
	"fmt"

	_tmConfig "github.com/tendermint/tendermint/config"
)

var (
	defaultTmDir = fmt.Sprintf("%s/tendermint", DefaultDataDir)
)

// BabbleConfig contains the configuration of a Babble node
type TmConfig struct {
	// Directory containing priv_key.pem and peers.json files
	DataDir string `mapstructure:"datadir"`
}

// DefaultBabbleConfig returns the default configuration for a Babble node
func DefaultTmConfig() *TmConfig {
	return &TmConfig{
		DataDir: defaultTmDir,
	}
}

// SetDataDir updates the babble configuration directories if they were set to
// to default values.
func (c *TmConfig) SetDataDir(datadir string) {
	if c.DataDir == defaultTmDir {
		c.DataDir = datadir
	}
}

// ToRealBabbleConfig converts an evm-lite/src/config.BabbleConfig to a
// babble/src/babble.BabbleConfig as used by Babble
func (c *TmConfig) ToRealTmConfig() *_tmConfig.Config {
	tmConfig := _tmConfig.DefaultConfig()
	tmConfig.SetRoot(c.DataDir)
	return tmConfig
}
