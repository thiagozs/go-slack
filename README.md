# go-slack, simple slack message notifier

> ***Work in progress***
---
> Need improvements
> * [ ] - Logger
> * [ ] - Logger debug
> * [ ] - Options

## Start with code

```golang
lopts := []slackr.SlackrOptions{
    slackr.CfgDebug(false),
    slackr.CfgToken(token),
}

slk, err := slackr.NewSlackClient(lopts)
if err != nil {
    log.Fatal(err)
}
```

Find user by e-mail

```golang
user, err := slk.SearchByEmail(email)
if err != nil {
    log.Fatal(err)
}
```

Find user by approche (fuzzy algorithm)

```golang
users, err := slk.SearchFuzzyMatch(slackr.REALNAME, term)
if err != nil {
    log.Fatal(err)
}
```
## Versioning and license

Our version numbers follow the [semantic versioning specification](http://semver.org/). You can see the available versions by checking the [tags on this repository](https://github.com/thiagozs/go-slack/tags). For more details about our license model, please take a look at the [LICENSE](LICENSE) file.

**2022**, thiagozs.