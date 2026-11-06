package common

import (
	"flag"
	"os"
	"strings"
)

type Env string

const (
	Development Env = "development"
	Test        Env = "test"
	Production  Env = "production"
)

func GetEnv() Env {
	e := os.Getenv("GOENV")
	if e != "" {
		return Env(e)
	}
	if strings.HasSuffix(os.Args[0], ".test") {
		return Test
	}
	if flag.Lookup("test.v") != nil {
		return Test
	}

	return Development
}
