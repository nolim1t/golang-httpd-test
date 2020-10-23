package bitcoind

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// common utilities
	// if commented out then we must redefine the following structs as outlined below
	// But must redefine common.Bitcoind as something else so it doesnt conflict
	"gitlab.com/nolim1t/golang-httpd-test/common"
)

/*
   Config Notes in case common isn't available:

   Inside common/config.go expecting to see the following struct
   Config struct {
       BitcoinClient           bool    `toml:"bitcoin-client"`
       // [bitcoind] section in the `--config` file that defines Bitcoind's setup
       Bitcoind Bitcoind `toml:"bitcoind"`
   }

   // Bitcoind config (common.Bitcoind)
   Bitcoind struct {
       Host string `toml:"host"`
       Port int64  `toml:"port"`
       User string `toml:"user"`
       Pass string `toml:"pass"`
   }

*/
const (
	DefaultHostname = "localhost"
	DefaultPort     = 8332
	DefaultUsername = "lncm"

	// Methods
	MethodGetBlockCount         = "getblockcount"
	MethodGetBlockchainInfo     = "getblockchaininfo"
	MethodGetNetworkInfo        = "getnetworkinfo"
	MethodGetNewAddress         = "getnetaddress"
	MethodImportAddress         = "importaddress"
	MethodListReceivedByAddress = "listreceivedbyaddress"

	Bech32 = "bech32"
)

type (
	Bitcoind struct {
		url, user, pass string
	}

	requestBody struct {
		JSONRPC string        `json: "jsonrpc"`
		ID      string        `json:"id"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}

	responseBody struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Code    int    `json:"code,omitempty"`
			Message string `json:"message,omitempty"`
		} `json:"error,omitempty"`
	}

	// Bitcoin structs
	// Response for 'getblockchainfo'
	// omit: softforks section
	BlockchainInfoResponse struct {
		Chain                string  `json:"chain"`
		Blocks               int64   `json:"blocks"`
		Headers              int64   `json:"headers"`
		BlockHash            string  `json:"bestblockhash"`
		Difficulty           float64 `json:"difficulty"`
		MedianTime           int64   `json:"mediantime"`
		VerificationProgress float64 `json:"verificationprogress"`
		InitialBlockDownload bool    `json:"initialblockdownload"`
		ChainWork            string  `json:"chainwork"`
		SizeOnDisk           int64   `json:"size_on_disk"`
		Pruned               bool    `json:"pruned,omitempty"`
		PruneHeight          int64   `json:"pruneheight,omitempty"`
		AutomaticPruning     bool    `json:"automatic_pruning,omitempty"`
		PruneTargetSize      int64   `json:"prune_target_size,omitempty"`
		ChainWarnings        string  `json:"warnings"`
	}
)

// Methods
// BlockCount
func (b Bitcoind) BlockCount() (count int64, err error) {
	res, _ := b.sendRequest(MethodGetBlockCount)
	err = json.Unmarshal(res, &count)

	return
}
// BlockchainInfo
func (b Bitcoind) BlockchainInfo() (blockresp BlockchainInfoResponse, err error) {
    res, err := b.sendRequest(MethodGetBlockchainInfo)
    if err != nil {
        return
    }
    err = json.Unmarshal(res, &blockresp)
    
    return
}

// sendRequest
func (b Bitcoind) sendRequest(method string, params ...interface{}) (response []byte, err error) {
	reqBody, err := json.Marshal(requestBody{
		JSONRPC: "1.0",
		Method:  method,
		Params:  params,
	})

	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", b.url, bytes.NewReader(reqBody))
	if err != nil {
		fmt.Printf("Error making request to %s", b.url)
		return
	}
	req.SetBasicAuth(b.user, b.pass)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer func() { _ = res.Body.Close() }()
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var resBody responseBody

	err = json.Unmarshal(resBytes, &resBody)
	if err != nil {
		return
	}
	fmt.Printf("Raw Response from %s: %s\n", method, resBody.Result)

	if resBody.Error != nil {
		return nil, fmt.Errorf("bitcoind error (%d): %s", resBody.Error.Code, resBody.Error.Message)
	}

	return resBody.Result, nil
}

// Create new object of Bitcoind client
func New(conf common.Bitcoind) (Bitcoind, error) {
	// Check if theres a bitcoin conf defined
	if conf.Host == "" {
		conf.Host = DefaultHostname
	}
	if conf.Port == 0 {
		conf.Port = DefaultPort
	}
	if conf.User == "" {
		conf.User = DefaultUsername
	}
	client := Bitcoind{
		url:  fmt.Sprintf("http://%s:%d", conf.Host, conf.Port),
		user: conf.User,
		pass: conf.Pass,
	}
	fmt.Printf("Creating bitcoin client... %s\n", client.url)
	_, err := client.BlockCount()
	if err != nil {
		return Bitcoind{}, fmt.Errorf("can't connect to Bitcoind: %w", err)
	}

	return client, nil
}
