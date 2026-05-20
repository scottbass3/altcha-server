package config

type Config struct {
	BaseUrl				string	`env:"ALTCHA_BASE_URL" envDefault:""`
	Port 				string	`env:"ALTCHA_PORT" envDefault:"3333"`
	HmacKey				string	`env:"ALTCHA_HMAC_KEY"`
	MaxNumber			int64	`env:"ALTCHA_MAX_NUMBER" envDefault:"1000000"`
	Algorithm			string	`env:"ALTCHA_ALGORITHM" envDefault:"SHA-256"`
	Salt				string	`env:"ALTCHA_SALT"`
	Expire				string	`env:"ALTCHA_EXPIRE" envDefault:"600s"`
	CheckExpire			bool	`env:"ALTCHA_CHECK_EXPIRE" envDefault:"true"`
	Debug				bool	`env:"ALTCHA_DEBUG" envDefault:"false"`
	DisableValidation	bool	`env:"ALTCHA_DISABLE_VALIDATION" envDefault:"false"`
}
