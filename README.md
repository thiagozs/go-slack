# go-slack, simple slack message notifier

Simple wrapper for slack user notification. The purpose of this project are find the user by some anchor and send a notification. The project has ability to search by approach (fuzzy algorithm).

You need a valid **token** bot of slack.
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

Find user by approach (fuzzy algorithm)

```golang
users, err := slk.SearchFuzzyMatch(slackr.REALNAME, term)
if err != nil {
    log.Fatal(err)
}
```
## Versioning and license

Our version numbers follow the [semantic versioning specification](http://semver.org/). You can see the available versions by checking the [tags on this repository](https://github.com/thiagozs/go-slack/tags). For more details about our license model, please take a look at the [LICENSE](LICENSE) file.

**2022**, thiagozs.