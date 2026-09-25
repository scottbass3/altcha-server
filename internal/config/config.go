package config

import (
	"fmt"
	"strconv"
	"time"
)

type Config struct {
	BaseUrl           string `env:"ALTCHA_BASE_URL" envDefault:""`
	Port              string `env:"ALTCHA_PORT" envDefault:"3333"`
	HmacKey           string `env:"ALTCHA_HMAC_KEY"`
	MaxNumber         int64  `env:"ALTCHA_MAX_NUMBER" envDefault:"1000000"`
	Algorithm         string `env:"ALTCHA_ALGORITHM" envDefault:"SHA-256"`
	Salt              string `env:"ALTCHA_SALT"`
	SaltLength        int    `env:"ALTCHA_SALT_LENGTH" envDefault:"12"`
	Expire            string `env:"ALTCHA_EXPIRE" envDefault:"600s"`
	CheckExpire       bool   `env:"ALTCHA_CHECK_EXPIRE" envDefault:"true"`
	CorsOrigins       string `env:"ALTCHA_CORS_ORIGINS" envDefault:"*"`
	StoreBackend      string `env:"ALTCHA_STORE" envDefault:"memory"`
	RedisURL          string `env:"ALTCHA_REDIS_URL"`
	MemcachedServers  string `env:"ALTCHA_MEMCACHED_SERVERS"`
	Debug             bool   `env:"ALTCHA_DEBUG" envDefault:"false"`
	DisableValidation bool   `env:"ALTCHA_DISABLE_VALIDATION" envDefault:"false"`
}

// ExpireDuration parses Expire as a Go duration, or as seconds when it is a bare integer.
func (c Config) ExpireDuration() (time.Duration, error) {
	if secs, err := strconv.ParseInt(c.Expire, 10, 64); err == nil {
		return time.Duration(secs) * time.Second, nil
	}
	d, err := time.ParseDuration(c.Expire)
	if err != nil {
		return 0, fmt.Errorf("invalid ALTCHA_EXPIRE %q: %w", c.Expire, err)
	}
	return d, nil
}
