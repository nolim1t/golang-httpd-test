package pineclient // Define package name
/*
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
*/

// Define imports
import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Define Constants (Accessible from outside this package)
const (
	BatteryPathName  = "/sys/class/power_supply/axp20x-battery"
	BatteryIsPresent = BatteryPathName + "/present"
	BatteryCapacity  = BatteryPathName + "/capacity"
)

// Local methods
func stringtofile(filename string) string {
	byte_output, err := ioutil.ReadFile(filename)
	if err == nil {
		return strings.Trim(string(byte_output), "\n")
	} else {
		return "-1"
	}
}

// Static Methods
// Status of Battery (-1 if file not present, 0 if battery not present, 1 if battery is present)
func GetStatus() string {
	return stringtofile(BatteryIsPresent)
}

// Capacity of Battery (-1 if file not present, or percentage)
func GetCapacity() string {
	return fmt.Sprintf("%s", stringtofile(BatteryCapacity))
}
