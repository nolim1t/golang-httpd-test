package pineclient

// Some sample pathnames (this pertains to my pinephone)
const (
    BatteryPathName = "/sys/class/power_supply/axp20x-battery"
    BatteryIsPresent = BatteryPathName + "/present"
    BatteryCapacity = BatteryPathName + "/capacity"
)

// Methods
func GetStatus() (string) {
    var status = "1"

    return status
}

