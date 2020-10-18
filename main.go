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
    "fmt"
    "flag"
    "os"

    // External libraries
    // mine
    "gitlab.com/nolim1t/golang-httpd-test/common"
    "gitlab.com/nolim1t/golang-httpd-test/pineclient"

    // github
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/gzip"
)

// Globals
var (
    version, gitHash string
    conf    common.Config
    showVersion    = flag.Bool("version", false, "Show version and exit")

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
}

func info(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "pong",
    })
}

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
// Main entrypoint
func main() {
    router := gin.Default()
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://nolim1t.co"},
        AllowHeaders:     []string{"Origin"},
        ExposeHeaders:    []string{"Content-Length"},
    }))
    router.Use(gzip.Gzip(gzip.DefaultCompression))

    r := router.Group("/api")
    r.GET("/info", info)
    // Pinephone stuff
    r.GET("/batteryStatus", batStatus)
    r.GET("/batteryCapacity", batCapacity)

    err := router.Run(fmt.Sprintf(":%d", 3000))
    if err != nil {
        panic(err)
    }
}
