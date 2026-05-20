package common

import (
	"fmt"
	"time"

	"github.com/scottbass3/altcha-server/internal/client"
	"github.com/scottbass3/altcha-server/internal/config"
)

func NewClientFromConfig(cfg config.Config) (*client.Client, error) {
	expire, err := time.ParseDuration(cfg.Expire)
	if err != nil {
		return nil, fmt.Errorf("invalid ALTCHA_EXPIRE %q: %w", cfg.Expire, err)
	}
	return client.New(cfg.HmacKey, cfg.MaxNumber, cfg.Algorithm, cfg.Salt, expire, cfg.CheckExpire)
}
