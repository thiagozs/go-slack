package options

type Options func(s *OptionsParams) error

type OptionsParams struct {
	Token string
	Debug bool
}

func CfgDebug(debug bool) Options {
	return func(c *OptionsParams) error {
		c.Debug = debug
		return nil
	}
}

func CfgToken(token string) Options {
	return func(c *OptionsParams) error {
		c.Token = token
		return nil
	}
}

// ------ getters

func (c *OptionsParams) GetDebug() bool {
	return c.Debug
}

func (c *OptionsParams) GetToken() string {
	return c.Token
}

// ------ setters

func (c *OptionsParams) SetDebug(debug bool) {
	c.Debug = debug
}

func (c *OptionsParams) SetToken(token string) {
	c.Token = token
}
