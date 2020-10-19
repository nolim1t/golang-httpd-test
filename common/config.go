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
        Port int64 `toml:"port"` // the port to run on
        StaticDir string `toml:"static-dir"` // Where index.html lives
        LogFile string `toml:"log-file"` // logfile to log
        DisablePinephoneBinding bool `toml:"disable-pinephone-binding"`
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
