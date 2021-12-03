package internal

import (
	"github.com/go-playground/validator/v10"
)

/*
Flags:
      --bar-example-string string   interesting value
  -e, --bar-enabled                 enables something
  -f, --foo string                  something something something... (default "gopher")
*/

type Bar struct {
	Example string `flag:"exampleString" flag-desc:"interesting value"`
	Enabled bool   `flag:"enable" flag-short:"e" flag-val:"false" flag-desc:"enables something"`
}

type Config struct {
	Foo string `flag:"foo" flag-short:"f" flag-val:"gopher" flag-desc:"something something something..."`
	Bar Bar    `flag:"bar"`
	Baz string `flag:"-"`
}

func (c *Config) PreValidate(_ *validator.Validate) {
	c.Baz = "another"
}
