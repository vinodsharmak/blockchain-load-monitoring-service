package helpers

// func TestTxPoolContent(t *testing.T) {
// 	t.Run("Received Tx pool Content", func(t *testing.T) {
// 		exepectedResult := model.Response{
// 			Result: model.Result{
// 				Pending: map[common.Address]map[int]model.TxBody{
// 					common.HexToAddress("0x0216d5032f356960cd3749c31ab34eeff21b3395"): {
// 						806: {
// 							From:  "0x0216d5032f356960cd3749c31ab34eeff21b3395",
// 							Gas:   "0x5208",
// 							Hash:  "0xaf953a2d01f55cfe080c0c94150a60105e8ac3d51153058a1f03dd239dd08586",
// 							Nonce: "0x326",
// 							To:    "0x7f69a91a3cf4be60020fb58b893b7cbb65376db8",
// 						},
// 					},
// 				},
// 				Queued: map[common.Address]map[int]model.TxBody{
// 					common.HexToAddress("0x0216d5032f356960cd3749c31ab34eeff21b3395"): {
// 						3: {
// 							From:  "0x976a3fc5d6f7d259ebfb4cc2ae75115475e9867c",
// 							Gas:   "0x15f90",
// 							Hash:  "0x57b30c59fc39a50e1cba90e3099286dfa5aaf60294a629240b5bbec6e2e66576",
// 							Nonce: "0x3",
// 							To:    "0x346fb27de7e7370008f5da379f74dd49f5f2f80f",
// 						},
// 					},
// 				},
// 			},
// 		}
// 		bytes, err := json.Marshal(exepectedResult)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
// 			rw.WriteHeader(http.StatusOK)
// 			_, err := rw.Write(bytes)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 		}))
// 		config.BlockchainURL = server.URL
// 		config.ChainID = "12345"
// 		constants.Response_format = "application/json"

// 		result, err := TxPoolContent()

// 		Ok(t, err)
// 		if !reflect.DeepEqual(exepectedResult, result) {
// 			t.Errorf("GetTxPoolContent = %v, want %v", exepectedResult, result)
// 		}
// 	})

// }

// func TestTxPoolStatus(t *testing.T) {
// 	t.Run("Received Tx pool status", func(t *testing.T) {
// 		exepectedResult := model.TxPoolStatusResponse{
// 			Result: model.TxPoolstatusresult{
// 				Pending: "0x00",
// 				Queued:  "0x00",
// 			},
// 		}
// 		bytes, err := json.Marshal(exepectedResult)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
// 			rw.WriteHeader(http.StatusOK)
// 			_, err := rw.Write(bytes)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 		}))
// 		config.BlockchainURL = server.URL
// 		config.ChainID = "12345"
// 		constants.Response_format = "application/json"

// 		result, err := TxPoolstatus()

// 		Ok(t, err)
// 		if !reflect.DeepEqual(exepectedResult, result) {
// 			t.Errorf("GetTxPoolStatus = %v, want %v", exepectedResult, result)
// 		}
// 	})

// }
