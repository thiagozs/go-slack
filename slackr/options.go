package slackr

func CfgDebug(debug bool) SlackrOptions {
	return func(c *SlackrParams) error {
		c.Debug = debug
		return nil
	}
}

func CfgToken(token string) SlackrOptions {
	return func(c *SlackrParams) error {
		c.Token = token
		return nil
	}
}
