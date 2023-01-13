package model

import "github.com/ethereum/go-ethereum/common"

type TxBody struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Gas   string `json:"gas"`
	Hash  string `json:"hash"`
	Nonce string `json:"nonce"`
}

type Result struct {
	Pending map[common.Address]map[int]TxBody `json:"pending"`
	Queued  map[common.Address]map[int]TxBody `json:"queued"`
}

type Response struct {
	Result Result `json:"result"`
}

type TxPoolstatusresult struct {
	Pending string
	Queued  string
}

type TxPoolStatusResponse struct {
	Result TxPoolstatusresult `json:"result"`
}
