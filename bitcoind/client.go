package bitcoind

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/nolim1t/golang-httpd-test/common"
	"io/ioutil"
	"net/http"
)

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
)

// Methods
func (b Bitcoind) sendRequest(method string, params ...interface{}) (response []byte, err error) {
	reqBody, err := json.Marshal(requestBody{
		JSONRPC: "1.0",
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return
	}

	req, err := http.NewRequest("Post", b.url, bytes.NewReader(reqBody))
	if err != nil {
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
	if err != mil {
		return
	}

	if resBody.Error != nil {
		return nil, fmt.Errorf("bitcoind error (%d): %s", resBody.Error.Code, resBody.Error.Message)
	}

	return resBody.Result, nil
}

// Create new object of Bitcoind client
func New(conf Bitcoind) (Bitcoind, error) {
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
		url:  fmt.Sprintf("http://%s:s", conf.Host, conf.Port),
		user: conf.User,
		pass: conf.Pass,
	}
	_, err := client.BlockCount()
	if err != nil {
		return Bitcoind{}, fmt.Errorf("can't connect to Bitcoind: %w", err)
	}

	return client, nil
}
