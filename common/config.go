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
		Port                    int64  `toml:"port" default:8080`                       // the port to run on
		StaticDir               string `toml:"static-dir" default:"~/public_html"`      // Where index.html lives (Default: $HOME/public_html)
		LogFile                 string `toml:"log-file" default:"~/http.log"`           // logfile to log (Default: ~/http.log)
		DisablePinephoneBinding bool   `toml:"disable-pinephone-binding" default:false` // disable-pinephone-binding=false
		BitcoinClient           bool   `toml:"bitcoin-client" default:true`             // bitcoin-client=true
		// [bitcoind] section in the `--config` file that defines Bitcoind's setup
		Bitcoind Bitcoind `toml:"bitcoind"`
		// auth-scheme key
		AuthScheme string `toml:"auth-scheme" default:"none"` // either use omitempty or default (https://godoc.org/github.com/pelletier/go-toml)
	}
	// Bitcoind config (enter some default values)
	// NOTE: Keep in mind that this is **not yet encrypted**, so best to keep it _local_
	Bitcoind struct {
		Host string `toml:"host" default:"localhost"`
		Port int64  `toml:"port" default:8332`
		User string `toml:"user" default:"lncm"`
		Pass string `toml:"pass" default:"lncmrocks"`
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
