package common

import (
    "os"
    "os/user"
    "path/filepath"
    "strings"
)
type (
    Config struct {
        // Port the service will run on
        Port int64 `toml:"port"`
        LogFile string `toml:"log-file"`
        OffChainOnly bool `toml:"disable-pinephone-binding"`
    }
)

func CleanAndExpandPath(path string) string {
    if path == "" {
        return ""
    }
    if strings.HasPrefix(path, "~") {
        var homeDir string
        u, err := user.Current()
        if err == nil {
            homeDir = u.HomeDir
        } else {
            homeDir = os.Getenv("HOME")
        }
        path = strings.Replace(path, "~", homeDir, 1)
    }

    return filepath.Clean(os.ExpandEnv(path))
}
