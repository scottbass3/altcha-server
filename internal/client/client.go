package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/altcha-org/altcha-lib-go"
)

type Client struct {
	hmacKey     string
	maxNumber   int64
	algorithm   altcha.Algorithm
	salt        string
	saltLength  int
	expire      time.Duration
	checkExpire bool
}

func New(hmacKey string, maxNumber int64, algorithm string, salt string, saltLength int, expire time.Duration, checkExpire bool) (*Client, error) {
	if len(hmacKey) < 16 {
		return nil, errors.New("ALTCHA_HMAC_KEY must be at least 16 characters")
	}
	alg := altcha.Algorithm(algorithm)
	switch alg {
	case altcha.SHA1, altcha.SHA256, altcha.SHA512:
	default:
		return nil, fmt.Errorf("unsupported algorithm %q: must be SHA-1, SHA-256, or SHA-512", algorithm)
	}
	return &Client{
		hmacKey:     hmacKey,
		maxNumber:   maxNumber,
		algorithm:   alg,
		salt:        salt,
		saltLength:  saltLength,
		expire:      expire,
		checkExpire: checkExpire,
	}, nil
}

func (c *Client) Generate() (altcha.Challenge, error) {
	expiration := time.Now().Add(c.expire)
	
	options := altcha.ChallengeOptions{
		HMACKey:    c.hmacKey,
		MaxNumber:  c.maxNumber,
		Algorithm:  c.algorithm,
		SaltLength: c.saltLength,
		Expires:    &expiration,
	}

	if len(c.salt) > 0 {
		options.Salt = c.salt
	}

	return altcha.CreateChallenge(options)
}

func (c *Client) Solve(ctx context.Context, challenge string) (*altcha.Solution, error) {
	return altcha.SolveChallenge(challenge, c.salt, c.algorithm, int(c.maxNumber), 0, ctx.Done())
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

