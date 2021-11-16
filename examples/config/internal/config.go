package internal

/*
Flags:
      --bar-example-string string   interesting value
  -e, --bar-enabled                 enables something
  -f, --foo string                  something something something... (default "gopher")
*/

type Bar struct {
	Example string `flag:"exampleString,,,interesting value"`
	Enabled bool   `flag:"enable,e,false,enables something"`
}

type Config struct {
	Foo string `flag:"foo,f,gopher,something something something..."`
	Bar Bar    `flag:"bar"`
	Baz string `flag:"-"`
}

func (c *Config) Defaults() {
	c.Baz = "another"
}
