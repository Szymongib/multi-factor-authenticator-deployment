package test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
)

var (
	testSuite *TestSuite
	config    Config
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	err := envconfig.InitWithPrefix(&config, "APP")
	exitOnError(err)

	testSuite, err = NewTestSuite(config)
	exitOnError(err)

	code := m.Run()
	os.Exit(code)
}

func exitOnError(err error) {
	if err != nil {
		logrus.Errorf("Error: %s", err.Error())
		os.Exit(1)
	}
}
