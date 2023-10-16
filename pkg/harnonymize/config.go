package harnonymize

import (
	"fmt"
)

var (
	ErrNotHARFile = fmt.Errorf("not a har file")
)

type Config struct {
	BlockContentKeywords []string
}

func New() *Config {
	return &Config{}
}
