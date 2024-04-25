package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type StoreType string

const ENV_PREFIX = "FLEMQ_"
const (
	StoreTypeFqueue StoreType = "fqueue"
	StoreTypeMqueue StoreType = "mqueue"
)

type Config struct {
	Addr       string `default:":22123"`
	Store      StoreConfig
	TLS        TLSConfig
	Connection ConnectionConfig
}

type StoreConfig struct {
	// Type of store to use, available: [mqueue, fqueue]
	Type StoreType `default:"fqueue"`
	// Only used for file store
	Folder string `default:"/tmp/flemq"`
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
	Enabled  bool `default:"false"`
}

type ConnectionConfig struct {
	RWTimeout time.Duration `default:"60s"`
}

// NewConfig returns a new Config struct with default values
// It loads environment variables with the prefix "FLEMQ_"
func NewConfig() Config {
	var config Config

	config.Addr = loadEnv(ENV_PREFIX, "ADDR", ":22123").(string)

	st := loadEnv(ENV_PREFIX, "STORE_TYPE", "fqueue").(string)
	config.Store.Type = StoreType(st)
	config.Store.Folder = loadEnv(ENV_PREFIX, "STORE_FOLDER", "/tmp/flemq").(string)

	config.TLS.Enabled = loadEnv(ENV_PREFIX, "TLS_ENABLED", false).(bool)
	config.TLS.CertFile = loadEnv(ENV_PREFIX, "TLS_CERT_FILE", "not_set").(string)
	config.TLS.KeyFile = loadEnv(ENV_PREFIX, "TLS_KEY_FILE", "not_set").(string)

	s := loadEnv(ENV_PREFIX, "RW_TIMEOUT_SEC", 60).(int)
	config.Connection.RWTimeout = time.Duration(s) * time.Second

	return config
}

// No, I'm not gonna use [one of the thousands Go config libraries](https://github.com/avelino/awesome-go?tab=readme-ov-file#configuration)
func loadEnv(prefix string, name string, def any) any {
	name = fmt.Sprintf("%s%s", prefix, name)

	switch d := def.(type) {
	case string:
		if os.Getenv(name) == "" {
			return d
		}
		return os.Getenv(name)
	case bool:
		if os.Getenv(name) == "" {
			return d
		}
		return os.Getenv(name) == "true"
	case int:
		if os.Getenv(name) == "" {
			return d
		}
		i, err := strconv.Atoi(os.Getenv(name))
		if err != nil {
			panic(err)
		}
		return i
	default:
		panic(fmt.Sprintf("Unknown type for %v", def))
	}
}
