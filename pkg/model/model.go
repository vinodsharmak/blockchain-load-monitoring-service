package model

import (
	"github.com/ethereum/go-ethereum/common"
)

type TxBody struct {
	From             common.Address `json:"from"`
	To               common.Address `json:"to"`
	Gas              string         `json:"gas"`
	Hash             common.Hash    `json:"hash"`
	Nonce            string         `json:"nonce"`
	FoundAtEpochTime int64
}

type Result struct {
	Pending map[common.Address]map[uint64]TxBody `json:"pending"`
	Queued  map[common.Address]map[uint64]TxBody `json:"queued"`
}

type Response struct {
	Result Result `json:"result"`
}

type TxPoolContentStuckMail struct {
	PendingCount   int
	QueuedCount    int
	PendingContent []TxBody
	QueuedContent  []TxBody
}
type TxPoolstatusresult struct {
	Pending string
	Queued  string
}

type TxPoolStatusResponse struct {
	Result TxPoolstatusresult `json:"result"`
}
