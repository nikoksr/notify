<h1 align="center">Welcome to notify (WIP) 👋</h1>
<p>
  <a href="#" target="_blank">
    <img alt="GitHub tag (latest SemVer)" src="https://img.shields.io/github/v/tag/nikoksr/notify">
  </a>
  <a href="#" target="_blank">
    <img alt="Lines of code" src="https://img.shields.io/tokei/lines/github/nikoksr/notify">
  </a>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> A dead simple Go library for sending notifications to various messaging platforms.

## Install

```sh
go get -u github.com/nikoksr/notify
```

## Example usage

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

## Roadmap

* [ ] Add tests
* [ ] Add more notification services

## Libraries in use

* [github.com/bwmarrin/discordgo](github.com/bwmarrin/discordgo)
* [github.com/jordan-wright/email](github.com/jordan-wright/email)
* [github.com/go-telegram-bot-api/telegram-bot-api](github.com/go-telegram-bot-api/telegram-bot-api)

## Author

👤 **Niko Köser**

* Twitter: [@nikoksr](https://twitter.com/nikoksr)
* Github: [@nikoksr](https://github.com/nikoksr)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/nikoksr/notify/issues). You can also take a look at the [contributing guide](https://github.com/nikoksr/notify/blob/main/CONTRIBUTING.md).

## Show your support

Give a ⭐️ if this project helped you!

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_
