package test

import (
	"math/rand"

	"github.com/szymongib/multi-factor-authenticator-core/pkg/api/contract"
	"github.com/szymongib/multi-factor-authenticator-core/pkg/client"
)

type Config struct {
	MultiFactorAuthenticatorCoreURL string `envconfig:"default=https://localhost:8000"`

	PasswordMethod struct {
		Name string `envconfig:"default=Password"`
	}

	SkipTLSVerify bool `envconfig:"default=true"`
}

type TestSuite struct {
	CoreClient *client.APIClient
}

func NewTestSuite(config Config) (*TestSuite, error) {

	return &TestSuite{
		CoreClient: client.NewAPIClient(config.MultiFactorAuthenticatorCoreURL, config.SkipTLSVerify),
	}, nil
}

func (ts *TestSuite) GenerateCredentials() contract.UserCredentials {
	randPart := randStringBytes(5)

	return contract.UserCredentials{
		Email:    randPart + "@test.com",
		Password: randPart + "-password",
	}
}

type PasswordAuthCredentials struct {
	Password string `json:"password"`
}

func (ts *TestSuite) GeneratePasswordAuthMethodCredentials() PasswordAuthCredentials {
	randPart := randStringBytes(5)

	return PasswordAuthCredentials{
		Password: "password-" + randPart,
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
