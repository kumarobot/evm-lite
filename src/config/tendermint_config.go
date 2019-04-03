package config

import (
	"fmt"

	tmConfig "github.com/tendermint/tendermint/config"
)

var (
	defaultTmDir = fmt.Sprintf("%s/tendermint", DefaultDataDir)
)

// BabbleConfig contains the configuration of a Babble node
type TmConfig struct {
	// Directory containing priv_key.pem and peers.json files
	DataDir    string `mapstructure:"datadir"`
	RealConfig *tmConfig.Config
}

// DefaultBabbleConfig returns the default configuration for a Babble node
func DefaultTmConfig() *TmConfig {
	var conf = &TmConfig{
		DataDir:    defaultTmDir,
		RealConfig: tmConfig.DefaultConfig().SetRoot(defaultTmDir),
	}
	return conf
}

// SetDataDir updates the tendermint configuration directories if they were set to
// to default values.
func (c *TmConfig) SetDataDir(datadir string) {
	if c.DataDir == defaultTmDir {
		c.DataDir = datadir
		c.RealConfig.SetRoot(datadir)
	}
}

// ToRealBabbleConfig converts an evm-lite/src/config.BabbleConfig to a
// babble/src/babble.BabbleConfig as used by Babble
func (c *TmConfig) ToRealTmConfig() *tmConfig.Config {
	return c.RealConfig
}
