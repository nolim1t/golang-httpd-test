package proto

// Some sample pathnames (this pertains to my pinephone)
const (
    BatteryPathName = "/sys/class/power_supply/axp20x-battery"
    BatteryIsPresent = BatteryPathName + "/present"
    BatteryCapacity = BatteryPathName + "/capacity"
)

func getStatus() (string) {
    var status = "1"

    return status
}
