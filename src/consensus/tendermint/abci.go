package tendermint

import(
    "github.com/bear987978897/evm-lite/src/state"
    "github.com/bear987978897/evm-lite/src/service"
    "github.com/tendermint/tendermint/abci/types"
    "github.com/sirupsen/logrus"
)

type abciProxy struct {
    types.BaseApplication

    service     *service.Service
    state       *state.State
    submitCh    chan []byte
    logger      *logrus.Entry
}

func NewABCIProxy (
    state       *state.State,
    service     *service.Service,
    submitCh    chan []byte,
    logger      *logrus.Logger,
) *abciProxy {
    return &abciProxy{
        service:  service,
        state:    state,
        submitCh: submitCh,
        logger:   logger.WithField("module", "tendermint/abci"),
    }
}

/********************************************************
Implement Tendermint ABCI application
TODO: Implement the functions in the interface if needed
*********************************************************/
