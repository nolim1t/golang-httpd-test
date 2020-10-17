package common

type (
    Config struct {
        // Port the service will run on
        Port int64 `toml:"port"`
        LogFile string `toml:"log-file"`
    }
)
