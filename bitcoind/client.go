package bitcoind

const (
    DefaultHostname = "localhost"
    DefaultPort = 8332
    DefaultUsername = "lncm"
    
    // Methods
    MethodGetBlockCount         = "getblockcount"
    MethodGetBlockchainInfo     = "getblockchaininfo"
    MethodGetNetworkInfo        = "getnetworkinfo"
    MethodGetNewAddress         = "getnetaddress"
    MethodImportAddress         = "importaddress"
    MethodListReceivedByAddress  = "listreceivedbyaddress"

    Bech32 = "bech32"
)

type (
    Bitcoind struct {
        url, user, pass string
    }

    requestBody struct {
        JSONRPC string          `json: "jsonrpc"`
        ID      string          `json:"id"`
        Method  string          `json:"method"`
        Params  []interface{}   `json:"params"`
    }

    responseBody struct {
        Result  json.RawMessage     `json:"result"`
        Error   *struct {
                    Code    int    `json:"code,omitempty"`
                    Message string `json:"message,omitempty"`
                }   `json:"error,omitempty"`
    }
)


