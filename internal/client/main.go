package client

import (
	"time"

	"github.com/altcha-org/altcha-lib-go"
)

type Client struct {
	hmacKey		string
	maxNumber	int64
	algorithm	altcha.Algorithm
	salt		string
	expire		string
	checkExpire	bool
}

func NewClient(hmacKey string, maxNumber int64, algorithm string, salt string, expire string, checkExpire bool) *Client {
	if len(hmacKey) == 0 {
		panic("HMAC key not found in env")
	}
	return &Client {
		hmacKey:		hmacKey,
		maxNumber:		maxNumber,
		algorithm:		altcha.Algorithm(algorithm),
		salt:			salt,
		expire:			expire,
		checkExpire:	checkExpire,
	}
}

func (c *Client) Generate() (altcha.Challenge, error) {
	expirationDuration, _ := time.ParseDuration(c.expire+"s")
	expiration := time.Now().Add(expirationDuration)
	
	options := altcha.ChallengeOptions{
		HMACKey:	c.hmacKey,
		MaxNumber:	c.maxNumber,
		Algorithm:	c.algorithm,
		Expires: 	&expiration,
	}

	if len(c.salt) > 0 {
		options.Salt = c.salt
	}

	return altcha.CreateChallenge(options)
}

func (c *Client) Solve(challenge string) (*altcha.Solution, error) {
	return altcha.SolveChallenge(challenge, c.salt, c.algorithm, int(c.maxNumber), 0, make(<-chan struct{}))
}

func (c *Client) VerifySolution(payload interface{}) (bool, error) {
	return altcha.VerifySolution(payload, c.hmacKey, c.checkExpire)
}

func (c *Client) VerifyServerSignature(payload interface{}) (bool, altcha.ServerSignatureVerificationData, error) {
	return altcha.VerifyServerSignature(payload, c.hmacKey)
}

func (c *Client) VerifyFieldsHash(formData map[string][]string, fields []string, fieldsHash string) (bool, error) {
	return altcha.VerifyFieldsHash(formData, fields, fieldsHash, c.algorithm)
}