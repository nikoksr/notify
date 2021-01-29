<div align="center">
<img
    width=40%
    src="assets/gopher-letter.svg"
    alt="notify logo"
/>

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/nikoksr/notify?color=success&label=version&sort=semver)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikoksr/notify)](https://goreportcard.com/report/github.com/nikoksr/notify)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/37fdff3c275c4a72a3a061f2d0ec5553)](https://www.codacy.com/gh/nikoksr/notify/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nikoksr/notify&amp;utm_campaign=Badge_Grade)
[![Maintainability](https://api.codeclimate.com/v1/badges/b3afd7bf115341995077/maintainability)](https://codeclimate.com/github/nikoksr/notify/maintainability)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify)

</div>

> <p align="center">A dead simple Go library for sending notifications to various messaging services.</p>

<h1></h1>

## About <a id="about"></a>

*Notify* arose from my own need for one of my api server running in production to be able to notify me when, for example, an error occurred. The library is kept as simple as possible to allow a quick integration and easy usage.

## Disclaimer <a id="disclaimer"></a>

Any misuse of this library is your own liability and responsibility and cannot be attributed to the authors of this library.  See [license](LICENSE) for more.

Spamming through the use of this library **may get you permanently banned** on most supported platforms.

## Install <a id="install"></a>

```sh
go get -u github.com/nikoksr/notify
```

## Example usage <a id="usage"></a>

```go
// The notifier we're gonna send our messages to
notifier := notify.New()

// Create a telegram service. Ignoring error for demo simplicity
telegramService, _ := telegram.New("your_telegram_api_token")

// Passing a telegram chat id as receiver for our messages.
// Basically where should our message be sent to?
telegramService.AddReceivers(-1234567890)

// Tell our notifier to use the telegram service. You can repeat the above process
// for as many services as you like and just tell the notifier to use them.
// Its kinda like using middlewares for api servers.
notifier.UseService(telegramService)

// Send a test message
_ = notifier.Send(
	"Message Subject/Title",
	"The actual message. Hello, you awesome gophers! :)",
)
```

## Supported services <a id="supported_services"></a>

- *Discord*
- *Email*
- *Microsoft Teams*
- *Slack*
- *Telegram*

## Roadmap <a id="roadmap"></a>

- [ ] Add tests (see [#1](https://github.com/nikoksr/notify/issues/1))
- [ ] Add more notification services

## Credits <a id="credits"></a>

- Discord support: [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)
- Email support: [jordan-wright/email](https://github.com/jordan-wright/email)
- Microsoft Teams support: [atc0005/go-teams-notify](https://github.com/atc0005/go-teams-notify)
- Slack support: [slack-go/slack](https://github.com/slack-go/slack)
- Telegram support: [go-telegram-bot-api/telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
- Logo: [MariaLetta/free-gophers-pack](https://github.com/MariaLetta/free-gophers-pack)

## Author <a id="author"></a>

**Niko Köser**

* Twitter: [@nikoksr](https://twitter.com/nikoksr)
* Github: [@nikoksr](https://github.com/nikoksr)

## Contributing <a id="contributing"></a>

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/nikoksr/notify/issues). You can also take a look at the [contributing guide](https://github.com/nikoksr/notify/blob/main/CONTRIBUTING.md).

## Show your support <a id="support"></a>

Give a ⭐️ if this project helped you!

