package common

import (
	"github.com/scottbass3/altcha-server/internal/client"
	"github.com/scottbass3/altcha-server/internal/config"
)

func NewClientFromConfig(cfg config.Config) (*client.Client, error) {
	expire, err := cfg.ExpireDuration()
	if err != nil {
		return nil, err
	}
	return client.New(cfg.HmacKey, cfg.MaxNumber, cfg.Algorithm, cfg.Salt, cfg.SaltLength, expire, cfg.CheckExpire)
}
