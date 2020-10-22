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

import (
	"fmt"
)

// Define Constants (Accessible from outside this package)
const (
	MagnetometerX  = "/sys/bus/iio/devices/iio\\:device3/in_magn_x_raw"
	MagnetometerY  = "/sys/bus/iio/devices/iio\\:device3/in_magn_y_raw"
	MagnetometerZ  = "/sys/bus/iio/devices/iio\\:device3/in_magn_z_raw"
)

func getMagneticHeading(axis string) string {
	var output = "-9999999"
	switch axis {
	case "X":
		output = stringtofile(MagnetometerX)
	case "Y":
		output = stringtofile(MagnetometerY)
	default:
		output = stringtofile(MagnetometerZ)
	}
	fmt.Print(output)
	return output
}
