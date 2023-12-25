package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const ENV_PREFIX = "FLEMQ_"

type Config struct {
	Addr string `default:":22123"`
	TLS  struct {
		Enabled  bool `default:"false"`
		CertFile string
		KeyFile  string
	}
	Connection struct {
		RWTimeout     time.Duration `default:"60s"`
		RecvChunkSize int           `default:"1024"`
	}
}

// NewConfig returns a new Config struct with default values
// It loads environment variables with the prefix "FLEMQ_"
func NewConfig() Config {
	var config Config

	config.Addr = loadEnv(ENV_PREFIX, "ADDR", ":22123").(string)

	config.TLS.Enabled = loadEnv(ENV_PREFIX, "TLS_ENABLED", false).(bool)
	config.TLS.CertFile = loadEnv(ENV_PREFIX, "TLS_CERT_FILE", "not_set").(string)
	config.TLS.KeyFile = loadEnv(ENV_PREFIX, "TLS_KEY_FILE", "not_set").(string)

	s := loadEnv(ENV_PREFIX, "RW_TIMEOUT_SEC", 60).(int)
	config.Connection.RWTimeout = time.Duration(s) * time.Second
	fmt.Println("TIMEOUT:", os.Getenv("FLEMQ_RW_TIMEOUT_SEC"), config.Connection.RWTimeout)
	config.Connection.RecvChunkSize = loadEnv(ENV_PREFIX, "RECV_CHUNK_SIZE", 1024).(int)

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
