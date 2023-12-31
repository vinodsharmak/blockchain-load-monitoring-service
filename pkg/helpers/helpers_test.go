package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/config"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/constants"
	"bitbucket.org/gath3rio/blockchain-load-monitoring-service.git/pkg/model"
	"github.com/ethereum/go-ethereum/common"
)

func TestTxPoolContent(t *testing.T) {
	t.Run("Received Tx pool Content", func(t *testing.T) {
		exepectedResult := model.Response{
			Result: model.Result{
				Pending: map[common.Address]map[uint64]model.TxBody{
					common.HexToAddress("0x0216d5032f356960cd3749c31ab34eeff21b3395"): {
						806: {
							From:  common.HexToAddress("0x0216d5032f356960cd3749c31ab34eeff21b3395"),
							Gas:   "0x5208",
							Hash:  common.HexToHash("0xaf953a2d01f55cfe080c0c94150a60105e8ac3d51153058a1f03dd239dd08586"),
							Nonce: "0x326",
							To:    common.HexToAddress("0x7f69a91a3cf4be60020fb58b893b7cbb65376db8"),
						},
					},
				},
				Queued: map[common.Address]map[uint64]model.TxBody{
					common.HexToAddress("0x0216d5032f356960cd3749c31ab34eeff21b3395"): {
						3: {
							From:  common.HexToAddress("0x0216d5032f356960cd3749c31ab34eeff21b3395"),
							Gas:   "0x5208",
							Hash:  common.HexToHash("0xaf953a2d01f55cfe080c0c94150a60105e8ac3d51153058a1f03dd239dd08586"),
							Nonce: "0x326",
							To:    common.HexToAddress("0x7f69a91a3cf4be60020fb58b893b7cbb65376db8"),
						},
					},
				},
			},
		}
		bytes, err := json.Marshal(exepectedResult)
		if err != nil {
			t.Fatal(err)
		}
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			_, err := rw.Write(bytes)
			if err != nil {
				t.Error(err)
			}
		}))
		config.BlockchainURL = server.URL
		config.ChainID = "12345"
		constants.Response_format = "application/json"

		result, err := TxPoolContent()

		Ok(t, err)
		if !reflect.DeepEqual(exepectedResult, result) {
			t.Errorf("GetTxPoolContent = %v, want %v", exepectedResult, result)
		}
	})

}

func TestTxPoolStatus(t *testing.T) {
	t.Run("Received Tx pool status", func(t *testing.T) {
		exepectedResult := model.TxPoolStatusResponse{
			Result: model.TxPoolstatusresult{
				Pending: "0x00",
				Queued:  "0x00",
			},
		}
		bytes, err := json.Marshal(exepectedResult)
		if err != nil {
			t.Fatal(err)
		}
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			_, err := rw.Write(bytes)
			if err != nil {
				t.Error(err)
			}
		}))
		config.BlockchainURL = server.URL
		config.ChainID = "12345"
		constants.Response_format = "application/json"

		result, err := TxPoolstatus()

		Ok(t, err)
		if !reflect.DeepEqual(exepectedResult, result) {
			t.Errorf("GetTxPoolStatus = %v, want %v", exepectedResult, result)
		}
	})

}
