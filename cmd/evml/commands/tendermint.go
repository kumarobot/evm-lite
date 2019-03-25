package commands

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/bear987978897/evm-lite/src/engine"
    "github.com/bear987978897/evm-lite/src/consensus/tendermint"
    "github.com/sirupsen/logrus"
)

// AddTendermintFlags adds flags to the Tendermint command
func AddTendermintFlags(cmd *cobra.Command) {
    // TODO
    cmd.Flags().String("tendermint.datadir", config.Tendermint.DataDir, "Tendermint data directory")
    viper.BindPFlags(cmd.Flags())
}

// NewTendermintCmd returns the command that starts EVM-Lite with Tendermint consensus
func NewTendermintCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:    "tendermint",
        Short:  "Run the evm-lite node with Tendermint consensus",
        PreRunE:    func(cmd *cobra.Command, args []string) (err error) {
            // TODO
            config.SetDataDir(config.BaseConfig.DataDir)
            logger.WithFields(logrus.Fields{
                "Babble": config.Babble,
            }).Debug("Config")
            return nil
        },
        RunE: runTendermint,
    }

    AddTendermintFlags(cmd)

    return cmd
}

func runTendermint(cmd *cobra.Command, args []string) error {
    tendermint := tendermint.NewTendermint(config.Tendermint, logger)
    engine, err := engine.NewEngine(*config, tendermint, logger)
    if err != nil {
        return fmt.Errorf("Error building Engine: %s", err)
    }

    engine.Run()

    return nil
}
