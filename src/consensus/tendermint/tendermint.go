package tendermint

import (
    "fmt"
    "os"

    _tmConfig "github.com/tendermint/tendermint/config"
    _tmProxy "github.com/tendermint/tendermint/proxy"
    _p2p "github.com/tendermint/tendermint/p2p"
    _node "github.com/tendermint/tendermint/node"
    _types "github.com/tendermint/tendermint/abci/types"
    _privval "github.com/tendermint/tendermint/privval"
    rpcclient "github.com/tendermint/tendermint/rpc/client"
    "github.com/bear987978897/evm-lite/src/config"
    "github.com/bear987978897/evm-lite/src/service"
    "github.com/bear987978897/evm-lite/src/state"
    "github.com/sirupsen/logrus"
    "github.com/tendermint/tendermint/libs/log"
)

type Tendermint struct {
    config      *config.TmConfig
    node        *_node.Node
    ethService  *service.Service
    ethState    *state.State
    logger      *logrus.Logger
    rpcclient   *rpcclient.HTTP
}

func NewTendermint(config *config.TmConfig, logger *logrus.Logger) *Tendermint {
    return &Tendermint{
	config: config,
	logger: logger,
    }
}

func DefaultNewNodeWithApp(config *_tmConfig.Config, app _types.Application, logger log.Logger) (*_node.Node, error) {
    // Generate node PrivKey
    nodeKey, err := _p2p.LoadOrGenNodeKey(config.NodeKeyFile())
    if err != nil {
	return nil, err
    }

    // Convert old PrivValidator if it exists.
    oldPrivVal := config.OldPrivValidatorFile()
    newPrivValKey := config.PrivValidatorKeyFile()
    newPrivValState := config.PrivValidatorStateFile()
    if _, err := os.Stat(oldPrivVal); !os.IsNotExist(err) {
	oldPV, err := _privval.LoadOldFilePV(oldPrivVal)
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


    return _node.NewNode(
        config,
        _privval.LoadOrGenFilePV(newPrivValKey, newPrivValState),
        nodeKey,
        _tmProxy.NewLocalClientCreator(app),
        _node.DefaultGenesisDocProviderFunc(config),
        _node.DefaultDBProvider,
        _node.DefaultMetricsProvider(config.Instrumentation),
        logger,
    )
}

/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

func (t *Tendermint) Init(state *state.State, service *service.Service) error {

    t.logger.Debug("INIT")

    t.ethState = state
    t.ethService = service

    abciApp := NewABCIProxy(state, service, service.GetSubmitCh(), t.logger)
    logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
    realConfig := t.config.ToRealTmConfig()
    node, err := DefaultNewNodeWithApp(realConfig, abciApp, logger)
    t.node = node
    t.rpcclient = rpcclient.NewHTTP(realConfig.RPC.ListenAddress, "/websocket");

    if err != nil {
        panic(err.Error())
    }

    return nil
}

func (t *Tendermint) Run() error {
    submitCh := t.ethService.GetSubmitCh()
    if err := t.node.Start(); err != nil {
        return fmt.Errorf("Failed to start node: %v", err)
    }
    t.logger.Info("Started node", "nodeInfo", t.node.Switch().NodeInfo())

    for {
        select {
        case tx := <-submitCh:
            _, err := t.rpcclient.BroadcastTxSync(tx)

            if err != nil {
                fmt.Errorf("Failed to broacast transaction: %v", err)
            }
        }
    }

    return nil
}

func (t *Tendermint) Info() (map[string]string, error) {
    tmInfo := t.node.NodeInfo().(_p2p.DefaultNodeInfo)
    info := map[string]string{
        "type"                  : "tendermint",
        "id"                    : string(tmInfo.ID()),
        "laddr"                 : tmInfo.ListenAddr,
        "network"               : tmInfo.Network,
        "version"               : tmInfo.Version,
        "moniker"               : tmInfo.Moniker,
        "tx_index"              : tmInfo.Other.TxIndex,
        "rpc_laddr"             : tmInfo.Other.RPCAddress,
    }
    return info, nil
}
