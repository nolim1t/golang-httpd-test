package btcprice

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
Pre-requisites

Imports (including standard library)
---
import (
    "gitlab.com/nolim1t/golang-httpd-test/common"
    "net/http"
    io/ioutil"
)
*/
import (
	"gitlab.com/nolim1t/golang-httpd-test/common"
	"io/ioutil"
	"net/http"
)

func GetPriceFeed(conf common.Config) ([]byte, error) {
	resp, err := http.Get(conf.BtcPriceApi)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	} else {
		return body, nil
	}

}
