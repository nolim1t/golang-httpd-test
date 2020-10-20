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
	"strconv"
)

// Define Constants (Accessible from outside this package)
const (
	ThermalsBasePath = "/sys/class/thermal"
	CPUBasePath      = ThermalsBasePath + "/thermal_zone0"
	GPUBasePath      = ThermalsBasePath + "/thermal_zone1"
	CPUTemp          = CPUBasePath + "/temp"
	GPUTemp          = GPUBasePath + "/temp"
)

// Static Methods
// CPU Temp (-1 if file not present, -2 if conversion error, or CPU temperature)
func GetCPUTemp() string {
	return gettemp("CPU")
}
func GetGPUTemp() string {
	return gettemp("GPU")
}

func gettemp(thermalclass string) string {
	var temp = "-3"
	if thermalclass == "GPU" {
		temp = stringtofile(GPUTemp)
	} else {
		temp = stringtofile(CPUTemp)
	}
	if temp == "-1" {
		return temp
	} else {
		tempint, err := strconv.Atoi(temp)
		if err != nil {
			return "-2"
		} else {
			return fmt.Sprintf("%d", (tempint / 1000))
		}
	}
}
