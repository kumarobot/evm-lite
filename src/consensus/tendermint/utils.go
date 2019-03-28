package tendermint

import (
    "fmt"
    "os"

    node "github.com/tendermint/tendermint/node"
    tmConfig "github.com/tendermint/tendermint/config"
    tmProxy "github.com/tendermint/tendermint/proxy"
    types "github.com/tendermint/tendermint/abci/types"
    privval "github.com/tendermint/tendermint/privval"
    p2p "github.com/tendermint/tendermint/p2p"
    "github.com/tendermint/tendermint/libs/log"
)


// Generate new tendermint node with abci application interface
func DefaultNewNodeWithApp(config *tmConfig.Config, app types.Application, logger log.Logger) (*node.Node, error) {
    // Generate node PrivKey
    nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
    if err != nil {
        return nil, err
    }

    // Convert old PrivValidator if it exists.
    oldPrivVal := config.OldPrivValidatorFile()
    newPrivValKey := config.PrivValidatorKeyFile()
    newPrivValState := config.PrivValidatorStateFile()
    if _, err := os.Stat(oldPrivVal); !os.IsNotExist(err) {
        oldPV, err := privval.LoadOldFilePV(oldPrivVal)
        if err != nil {
            return nil, fmt.Errorf("Error reading OldPrivValidator from %v: %v\n", oldPrivVal, err)
        }
        logger.Info(
            "Upgrading PrivValidator file",
            "old", oldPrivVal,
            "newKey", newPrivValKey,
            "newState", newPrivValState,
        )
        oldPV.Upgrade(newPrivValKey, newPrivValState)
    }


    return node.NewNode(
        config,
        privval.LoadOrGenFilePV(newPrivValKey, newPrivValState),
        nodeKey,
        tmProxy.NewLocalClientCreator(app),
        node.DefaultGenesisDocProviderFunc(config),
        node.DefaultDBProvider,
        node.DefaultMetricsProvider(config.Instrumentation),
        logger,
    )
}
