package bitcoind

/*
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
*/

/*
Reference:
https://developer.bitcoin.org/reference/rpc/index.html
Add to https://gitlab.com/nolim1t/golang-httpd-test/-/issues for stuff to implement.

Interfacing with this package (add to the file)

-----
import (
        "gitlab.com/nolim1t/golang-httpd-test/bitcoind"
)

type (
        BitcoinClient interface {
                BlockCount() (int64, error)
                BlockchainInfo() (bitcoind.BlockchainInfoResponse, error)
                NetworkInfo() (bitcoind.NetworkInfoResponse, err error)
                GetTransactionInfo(string) (bitcoind.VerboseTransactionInfo, error)
                GetMempoolContents() (mempoolcontents []string, err error)
                PushTransaction(hex string) (txid string, err error)
                GetBestBlockHash() (blockhash string, err error)
                GetBlockHashByHeight(height int64) (string, error)
                GetBlock(hash string) (bitcoind.BitcoinBlockResponse, error)
                GetMempoolInfo() (bitcoind.MempoolInfoResponse, error)
                GetMiningInfo() (bitcoind.MiningInfoResponse, error)
                GetPeerInfo() ([]PeerInfo, error)
                GetBlockStats(int64) (bitcoind.BlockStatsResponse, error)
        }
)

var (
        btcClient BitcoinClient
)
*/

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
	MethodGetRawTransaction     = "getrawtransaction"
	MethodGetMempoolContents    = "getrawmempool"
	// Transaction Broadcasting
	// https://developer.bitcoin.org/reference/rpc/sendrawtransaction.html
	MethodBroadcastTx = "sendrawtransaction"
	// Blockchain hash stuff
	MethodGetBlock        = "getblock" // verbosity = 1
	MethodGetBestBlock    = "getbestblockhash"
	MethodGetHashByHeight = "getblockhash"
	MethodGetMempool      = "getmempoolinfo"
	MethodGetMiningInfo   = "getmininginfo"
	MethodGetPeerInfo     = "getpeerinfo"
	MethodGetBlockStats   = "getblockstats"

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

	// Input Transactions (Unspent UTXOs to build TX from)
	TransactionInput struct {
		TransactionID string `json:"txid"`
		VoutID        int64  `json:"vout"`
		Sequence      int64  `json:"sequence"`
	}
	// scriptPubKey struct in Transaction output
	ScriptPubKeyObj struct {
		ASMCode              string   `json:"asm"`
		HexCode              string   `json:"hex"`
		ScriptType           string   `json:"type"`
		RequiredSigs         int64    `json:"reqsigs"`
		TransactionAddresses []string `json:"addresses"`
	}
	// New UTXO to move transaction to
	TransactionOutput struct {
		TransactionValue float64         `json:"value"`
		TransactionIndex int64           `json:"n"`
		ScriptPubKey     ScriptPubKeyObj `json:"scriptPubKey"`
	}

	// UTXO
	// Response for getrawtransaction
	VerboseTransactionInfo struct {
		TransactionID   string              `json:"txid"`
		TransactionHash string              `json:"hash"`
		TransactionSize int64               `json:"size"`
		TransactionHex  string              `json:"hex"`
		Confirmations   int64               `json:"confirmations,omitempty"`
		Time            int64               `json:"time,omitempty"`
		Blocktime       int64               `json:"blocktime,omitempty"`
		Blockhash       string              `json:"blockhash,omitempty"`
		Vin             []TransactionInput  `json:"vin"`
		Vout            []TransactionOutput `json:"vout"`
	}

	// Struct for bitcoin block
	BitcoinBlockResponse struct {
		Hash              string   `json:"hash"`
		Confirmations     int64    `json:"confirmations"`
		Size              int64    `json:"size"`
		StrippedSize      int64    `json:"strippedsize"`
		Weight            int64    `json:"weight"`
		Height            int64    `json:"height"`
		Version           int64    `json:"version"`
		VersionHex        string   `json:"versionhex"`
		MerkleRoot        string   `json:"merkleroot"`
		Transactions      []string `json:"tx"`
		Time              int64    `json:"time"`
		MedianTime        int64    `json:"mediantime"`
		Nonce             int64    `json:"nonce"`
		Bits              string   `json:"bits"`
		Difficulty        float64  `json:"difficulty"`
		Chainwork         string   `json:"chainwork"`
		PreviousBlockHash string   `json:"previousblockhash"`
		NextBlockHash     string   `json:"nextblockhash"`
	}
	// Get mempool info struct
	MempoolInfoResponse struct {
		Size          int64   `json:"size"`
		Bytes         int64   `json:"bytes"`
		Usage         int64   `json:"usage"`
		MaxMempool    int64   `json:"maxmempool"`
		MempoolMinFee float64 `json:"mempoolminfee"`
		MinRelayTxFee float64 `json:"minrelaytxfee"`
	}
	// NetworkList struct
	NetworkList struct {
		Name                      string `json:"string"`
		Limited                   bool   `json:"limited"`
		Reachable                 bool   `json:"reachable"`
		Proxy                     string `json:"proxy"`
		ProxyRandomizeCredentials bool   `json:"proxy_randomize_credentials"`
	}
	// AddressList struct
	AddressList struct {
		Address string `json:"address"`
		Port    int64  `json:"port"`
		Score   int64  `json:"score"`
	}
	// getnetworkinfo struct
	NetworkInfoResponse struct {
		Version         int64         `json:"version"`
		SubVersion      string        `json:"subversion"`
		ProtocolVersion int64         `json:"protocolversion"`
		LocalServices   string        `json:"localservices"`
		LocalRelay      bool          `json:"localrelay"`
		Connections     int64         `json:"connections"`
		NetworkActive   bool          `json:"networkactive"`
		Networks        []NetworkList `json:"networks"`
		RelayFee        float64       `json:"relayfee"`
		IncrementalFee  float64       `json:"incrementalfee"`
		LocalAddresses  []AddressList `json:"localaddresses"`
		Warnings        string        `json:"warnings"`
	}
	// getmininginfo
	MiningInfoResponse struct {
		Blocks                  int64   `json:"blocks"`
		CurrentBlockWeight      int64   `json:"currentblockweight"`
		CurrentBlockTransaction int64   `json:"currentblocktx"`
		Difficulty              float64 `json:"difficulty"`
		NetworkHashPs           float64 `json:"networkhashps"`
		PooledTransaction       int64   `json:"pooledtx"`
		Chain                   string  `json:"chain"`
		Warnings                string  `json:"warnings"`
	}
	// PeerInfo struct
	PeerInfo struct {
		Id             int64    `json:"id"`
		Addr           string   `json:"addr"`
		AddrBind       string   `json:"addrbind"`
		AddrLocal      string   `json:"addrlocal"`
		Services       string   `json:"services"`
		ServicesName   []string `json:"servicesnames"`
		RelayTx        bool     `json:"relaytxes"`
		LastSend       int64    `json:"lastsend"`
		LastRecv       int64    `json:"lastrcv"`
		BytesSent      int64    `json:"bytessent"`
		BytesRecv      int64    `json:"bytesrecv"`
		ConnTime       int64    `json:"conntime"`
		TimeOffset     int64    `json:"timeoffset"`
		PingTime       float64  `json:"pingtime"`
		MinPing        float64  `json:"minping"`
		PingWait       float64  `json:"pingwait"`
		Version        int64    `json:"version"`
		SubVer         string   `json:"subver"`
		Inbound        bool     `json:"inbound"`
		AddNode        bool     `json:"addnode"`
		StartingHeight int64    `json:"startingheight"`
		BanScore       int64    `json:"banscore"`
		SyncedHeaders  int64    `json:"synced_headers"`
		SyncedBlocks   int64    `json:"synced_blocks"`
		InFlight       []int64  `json:"inflight"`
		WhiteListed    bool     `json:"whitelisted"`
		MinFeeFilter   float64  `json:"minfeefilter"`
	}
	BlockStatsResponse struct {
		AvgFee        int64   `json:"avgfee"`
		AvgFeeRate    int64   `json:"avgfeerate"`
		AvgFeeSize    int64   `json:"avgfeesize"`
		Blockhash     string  `json:"blockhash"`
		FeeRates      []int64 `json:"feerate_percentiles"`
		Height        int64   `json:"height"`
		Ins           int64   `json:"ins"`
		MaxFee        int64   `json:"maxfee"`
		MaxFeeRate    int64   `json:"maxfeerate"`
		MaxTxSize     int64   `json:"maxtxsize"`
		MedianFee     int64   `json:"medianfee"`
		MedianTime    int64   `json:"mediantime"`
		MedianTxSize  int64   `json:"mediantxsize"`
		MinFee        int64   `json:"minfee"`
		MinFeeRate    int64   `json:"minfeerate"`
		MinTxSize     int64   `json:"mintxsize"`
		Outs          int64   `json:"outs"`
		Subsidy       int64   `json:"subsidy"`
		SWTotalSize   int64   `json:"swtotal_size"`
		SWTotalWeight int64   `json:"swtotal_weight"`
		SWTxs         int64   `json:"swtxs"`
		Time          int64   `json:"time"`
		TotalOut      int64   `json:"total_out"`
		TotalSize     int64   `json:"total_size"`
		TotalWeight   int64   `json:"total_weight"`
		TotalFee      int64   `json:"total_fee"`
		Txs           int64   `json:"txs"`
		UTXOIncrease  int64   `json:"utxo_increase"`
		UTXOSizeInc   int64   `json:"utxo_size_inc"`
	}
)

// Methods
// GetBlockstats
func (b Bitcoind) GetBlockStats(height int64) (blockstats BlockStatsResponse, err error) {
	res, _ := b.sendRequest(MethodGetBlockStats, height)
	err = json.Unmarshal(res, &blockstats)

	return
}

// BlockCount
func (b Bitcoind) BlockCount() (count int64, err error) {
	res, _ := b.sendRequest(MethodGetBlockCount)
	err = json.Unmarshal(res, &count)

	return
}

func (b Bitcoind) GetPeerInfo() (peerinfo []PeerInfo, err error) {
	res, err := b.sendRequest(MethodGetPeerInfo)
	err = json.Unmarshal(res, &peerinfo)

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

// NetworkInfo
func (b Bitcoind) NetworkInfo() (nwinforesp NetworkInfoResponse, err error) {
	res, err := b.sendRequest(MethodGetNetworkInfo)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &nwinforesp)

	return
}

// Get transaction Info
func (b Bitcoind) GetTransactionInfo(txid string) (txinfo VerboseTransactionInfo, err error) {
	res, err := b.sendRequest(MethodGetRawTransaction, txid, 1)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &txinfo)

	return
}

// Get raw mempool
func (b Bitcoind) GetMempoolContents() (mempoolcontents []string, err error) {
	res, err := b.sendRequest(MethodGetMempoolContents)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &mempoolcontents)

	return
}

// MethodGetMempool
func (b Bitcoind) GetMempoolInfo() (mempoolinfo MempoolInfoResponse, err error) {
	res, err := b.sendRequest(MethodGetMempool)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &mempoolinfo)

	return
}

// Broadcast TX
func (b Bitcoind) PushTransaction(hex string) (txid string, err error) {
	res, err := b.sendRequest(MethodBroadcastTx, hex)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &txid)

	return
}

// Get best block hash
func (b Bitcoind) GetBestBlockHash() (blockhash string, err error) {
	res, err := b.sendRequest(MethodGetBestBlock)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &blockhash)

	return
}

// get block hash by height
func (b Bitcoind) GetBlockHashByHeight(height int64) (blockhash string, err error) {
	res, err := b.sendRequest(MethodGetHashByHeight, height)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &blockhash)

	return
}

// getblock (MethodGetBlock)
func (b Bitcoind) GetBlock(hash string) (blockinfo BitcoinBlockResponse, err error) {
	res, err := b.sendRequest(MethodGetBlock, hash, 1)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &blockinfo)

	return
}

// GetMiningInfo
func (b Bitcoind) GetMiningInfo() (mininginfo MiningInfoResponse, err error) {
	res, err := b.sendRequest(MethodGetMiningInfo)
	if err != nil {
		return
	}
	err = json.Unmarshal(res, &mininginfo)
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
	//fmt.Printf("Raw Response from %s: %s\n", method, resBody.Result)

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
