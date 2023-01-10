package model

import "github.com/ethereum/go-ethereum/common"

type txBody struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Gas   string `json:"gas"`
	Hash  string `json:"hash"`
	Nonce string `json:"nonce"`
}

type result struct {
	Pending map[common.Address]map[int]txBody `json:"pending"`
	Queued  map[common.Address]map[int]txBody `json:"queued"`
}

type Response struct {
	Result result `json:"result"`
}
