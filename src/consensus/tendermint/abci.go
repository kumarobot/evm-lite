package tendermint

import (
	"encoding/hex"
	"github.com/bear987978897/evm-lite/src/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/abci/types"
)

type ABCIProxy struct {
	types.BaseApplication

	state     *state.State
	logger    *logrus.Entry
	blockHash common.Hash
	txIndex   int
}


func NewABCIProxy(
	state *state.State,
	logger *logrus.Logger,
) *ABCIProxy {
	return &ABCIProxy{
		state:     state,
		logger:    logger.WithField("module", "tendermint/abci"),
		blockHash: common.Hash{},
		txIndex:   0,
	}
}

/********************************************************
Implement Tendermint ABCI application
*********************************************************/
// TODO: Implement CheckTx
func (p *ABCIProxy) CheckTx(tx []byte) types.ResponseCheckTx {
	return types.ResponseCheckTx{Code: types.CodeTypeOK}
}

func (p *ABCIProxy) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	p.blockHash = common.BytesToHash(req.Hash)

	p.logger.Debug("Begin block: ", p.blockHash.String())

	return types.ResponseBeginBlock{}
}

func (p *ABCIProxy) DeliverTx(tx []byte) types.ResponseDeliverTx {
	err := p.state.ApplyTransaction(tx, p.txIndex, p.blockHash)
	if err != nil {
		p.logger.Error("DeliverTx Error: ", err)
		return types.ResponseDeliverTx{Code: 1}
	}
	p.txIndex++

	dst := make([]byte, hex.DecodedLen(len(tx)))
	txHex, _ := hex.Decode(dst, tx)
	p.logger.Debug("DeliverTx: ", txHex)

	return types.ResponseDeliverTx{Code: types.CodeTypeOK}
}

// TODO: Return application state root
func (p *ABCIProxy) Commit() types.ResponseCommit {
	p.txIndex = 0

	hash, err := p.state.Commit()
	if err != nil {
		p.logger.Error("Commit Error: ", err)
		return types.ResponseCommit{}
	}
	//p.logger.Debug("Block commited: ", hash)
	s := hash.Hex()
	//p.logger.Debug("hash.Hex: ", s)
	data, err := hex.DecodeString(s[2:])
	if err != nil{
		p.logger.Debug("data error: ", err)
	}
	//p.logger.Debug("Byte data: ", data)		

	return types.ResponseCommit{Data: data}
	//return types.ResponseCommit{}
}
