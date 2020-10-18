package pineclient

import (
    "fmt"
    "io/ioutil"
    "strings"
)

// Some sample pathnames (this pertains to my pinephone)
const (
    BatteryPathName = "/sys/class/power_supply/axp20x-battery"
    BatteryIsPresent = BatteryPathName + "/present"
    BatteryCapacity = BatteryPathName + "/capacity"
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
func GetStatus() (string) {
    return stringtofile(BatteryIsPresent)
}

func GetCapacity() (string) {
    return fmt.Sprintf("%s", stringtofile(BatteryCapacity))
}

