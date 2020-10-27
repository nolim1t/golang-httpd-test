package main

/*
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

*/
import (
	// System Libraries
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"

	// External libraries
	// mine
	"gitlab.com/nolim1t/golang-httpd-test/bitcoind"
	"gitlab.com/nolim1t/golang-httpd-test/common"
	"gitlab.com/nolim1t/golang-httpd-test/pineclient"

	// github
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	// non-github
	"gopkg.in/natefinch/lumberjack.v2"
)

// types
type (
	// How to read from Bitcoin client
	BitcoinClient interface {
		BlockCount() (int64, error)
		BlockchainInfo() (bitcoind.BlockchainInfoResponse, error)
		NetworkInfo() (nwinforesp string, err error)
		GetTransactionInfo(string) (bitcoind.VerboseTransactionInfo, error)
		GetMempoolContents() (mempoolcontents []string, err error)
		PushTransaction(hex string) (txid string, err error)
		GetBestBlockHash() (blockhash string, err error)
		GetBlockHashByHeight(height int64) (blockhash string, err error)
	}
)

// Globals
var (
	version, gitHash string
	// Accessing bitcoinclient
	btcClient BitcoinClient

	conf           common.Config
	showVersion    = flag.Bool("version", false, "Show version and exit")
	configFilePath = flag.String("config", common.DefaultConfigFile, "Path to a config file in TOML format")
)

// Functions

// Init function
func init() {
	flag.Parse()
	versionString := "debug"

	if version != "" && gitHash != "" {
		versionString = fmt.Sprintf("%s (git: %s)", version, gitHash)
	}

	if *showVersion {
		fmt.Println(versionString)
		os.Exit(0)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	// Load configfile
	configFile, err := toml.LoadFile(common.CleanAndExpandPath(*configFilePath))
	// Config file
	if err != nil {
		panic(fmt.Errorf("unable to load %s:\n\t%w", *configFilePath, err))
	}
	err = configFile.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("unable to process %s:\n\t%w", *configFilePath, err))
	}
	// set up logfile
	if conf.LogFile == "" {
		conf.LogFile = common.DefaultLogFile
	}
	fields := log.Fields{
		"version":   versionString,
		"log-file":  conf.LogFile,
		"conf-file": *configFilePath,
	}
	if conf.LogFile != "none" {
		log.SetOutput(&lumberjack.Logger{
			Filename:  common.CleanAndExpandPath(conf.LogFile),
			LocalTime: true,
			Compress:  true,
		})
		log.SetFormatter(&log.JSONFormatter{
			PrettyPrint: false, // so 'jq' always works in 'tail -f'
		})
		log.WithFields(fields).Println("server started")
	}
	// if bitcoin client enabled
	if conf.BitcoinClient {
		btcClient, err = bitcoind.New(conf.Bitcoind)
		if err != nil {
			panic(err)
		}
	}
}

// Test endpoint
func info(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// Bitcoin endpoints
// begin: bitcoin functions
func blockchainInfo(c *gin.Context) {
	blockchainInforesp, err := btcClient.BlockchainInfo()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Can't get blockchain info",
		})
		return
	}

	c.JSON(200, gin.H{
		"message":        "OK",
		"blockchaininfo": blockchainInforesp,
	})
}

func blockchainTxInfo(c *gin.Context) {
	txInforesp, err := btcClient.GetTransactionInfo(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Can't access transaction index: %s", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "OK",
		"txinfo":  txInforesp,
	})
}
func mempoolContents(c *gin.Context) {
	// GetMempoolContents() (mempoolcontents []string, err error)
	mempoolInfo, err := btcClient.GetMempoolContents()
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Can't access mempool: %s", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "OK",
		"mempool": mempoolInfo,
	})
}
func pushTransaction(c *gin.Context) {
	// PushTransaction(hex string) (txid string, err error)
	pushTxRes, err := btcClient.PushTransaction(c.PostForm("hex"))
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Can't broadcast transaction: %s", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "OK",
		"txid":    pushTxRes,
	})
}
func getBestBlockHash(c *gin.Context) {
	// GetBestBlockHash() (blockhash string, err error)
	bestblock, err := btcClient.GetBestBlockHash()
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Error getting the block hash: %s", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"message":   "OK",
		"blockhash": bestblock,
	})
}

// get blockhash by height
func getBlockHashByHeight(c *gin.Context) {
	// GetBlockHashByHeight(height int64) (blockhash string, err error)
	heightInt, errtoInt := strconv.ParseInt(c.Param("id"), 10, 64)
	if errtoInt != nil {
		c.JSON(500, gin.H{
			"message": "Error converting input to integer",
		})
		return
	}
	blockhash, err := btcClient.GetBlockHashByHeight(heightInt)
	if err != nil {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("Error getting the block hash: %s", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"message":   "OK",
		"blockhash": blockhash,
	})
}

// index endpoint
// PinePhone Endpoints
func batStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": pineclient.GetStatus(),
	})
}

func batCapacity(c *gin.Context) {
	c.JSON(200, gin.H{
		"percent": pineclient.GetCapacity(),
	})
}
func cpuTemp(c *gin.Context) {
	c.JSON(200, gin.H{
		"cputemp": pineclient.GetCPUTemp(),
	})
}
func gpuTemp(c *gin.Context) {
	c.JSON(200, gin.H{
		"gputemp": pineclient.GetGPUTemp(),
	})
}

// querystring test
func testQueryString(c *gin.Context) {
	param1 := c.DefaultQuery("param1", "none")
	if param1 == "none" {
		c.JSON(200, gin.H{
			"message": "show this message if there is no query string",
		})
	} else {
		c.JSON(200, gin.H{
			"message": param1,
		})
	}
}

// Main entrypoint
func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	r := router.Group("/api")
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, common.FormatRoutes(router.Routes()))
	})
	r.GET("/info", info)
	if conf.BitcoinClient {
		fmt.Println("Bitcoin client enabled")
		r.GET("/test", testQueryString)
		// Bitcoin
		r.GET("/blockchainInfo", blockchainInfo)        // blockchainInfo
		r.GET("/txid/:id", blockchainTxInfo)            // txid
		r.GET("/mempool", mempoolContents)              // mempool contents
		r.POST("/pushtx", pushTransaction)              // Push transaction
		r.GET("/getblockhash", getBestBlockHash)        // Get best blockhash
		r.GET("/blockheight/:id", getBlockHashByHeight) // get blockhash by height
	} else {
		fmt.Println("Bitcoin client not enabled")
	}
	// Pinephone stuff
	r.GET("/batteryStatus", batStatus)
	r.GET("/batteryCapacity", batCapacity)
	r.GET("/cpuTemp", cpuTemp)
	r.GET("/gpuTemp", gpuTemp)
	if conf.Port == 0 {
		conf.Port = 8080
	}
	var staticFilePath string
	if conf.StaticDir != "" {
		staticFilePath = path.Join(conf.StaticDir, "index.html")
		fmt.Println(conf.StaticDir)
		fmt.Println(staticFilePath)
		router.StaticFile("/", common.CleanAndExpandPath(staticFilePath))
	}
	log.WithFields(log.Fields{
		"routes":      common.FormatRoutes(router.Routes()),
		"port":        conf.Port,
		"static-file": staticFilePath,
	}).Println("gin router defined")
	err := router.Run(fmt.Sprintf(":%d", conf.Port))

	if err != nil {
		panic(err)
	}
}
