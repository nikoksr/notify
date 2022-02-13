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

*Notify* was born out of my own need to have my API servers running in production be able to notify me when critical errors occur. Of course, _Notify_ can be used for any other purpose as well. The library is kept as simple as possible for quick integration and ease of use.

## Disclaimer <a id="disclaimer"></a>

Any misuse of this library is your own liability and responsibility and cannot be attributed to the authors of this library.  See [license](LICENSE) for more.

Spamming through the use of this library **may get you permanently banned** on most supported platforms.

Since Notify is highly dependent on the consistency of the supported external services and the corresponding latest client libraries, we cannot guarantee its reliability nor its consistency, and therefore you should probably not use or rely on Notify in critical scenarios.

## Install <a id="install"></a>

```sh
go get -u github.com/nikoksr/notify
```

## Example usage <a id="usage"></a>

```go
// Create a telegram service. Ignoring error for demo simplicity.
telegramService, _ := telegram.New("your_telegram_api_token")

// Passing a telegram chat id as receiver for our messages.
// Basically where should our message be sent?
telegramService.AddReceivers(-1234567890)

// Tell our notifier to use the telegram service. You can repeat the above process
// for as many services as you like and just tell the notifier to use them.
// Inspired by http middlewares used in higher level libraries.
notify.UseServices(telegramService)

// Send a test message.
_ = notify.Send(
	context.Background(),
	"Subject/Title",
	"The actual message - Hello, you awesome gophers! :)",
)
```

## Supported services <a id="supported_services"></a>

> Please create feature requests for missing services (see [#3](https://github.com/nikoksr/notify/issues/3) for example)

| Service                                                      | Path                                                                                 | Credits                                                                                         |
|--------------------------------------------------------------|--------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------|
| [Amazon SES](https://aws.amazon.com/ses)                     | [service/amazonses](https://github.com/nikoksr/notify/tree/main/service/amazonses)   | [aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)                                       |
| [Amazon SNS](https://aws.amazon.com/sns)                     | [service/amazonsns](https://github.com/nikoksr/notify/tree/main/service/amazonsns)   | [aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)                                       |
| [DingTalk](https://www.dingtalk.com)                         | [service/dinding](https://github.com/nikoksr/notify/tree/main/service/dingding)      | [blinkbean/dingtalk](https://github.com/blinkbean/dingtalk)                                     |
| [Discord](https://discord.com)                               | [service/discord](https://github.com/nikoksr/notify/tree/main/service/discord)       | [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)                                     |
| [Email](https://wikipedia.org/wiki/Email)                    | [service/mail](https://github.com/nikoksr/notify/tree/main/service/mail)             | [jordan-wright/email](https://github.com/jordan-wright/email)                                   |
| [Line](https://line.me)                                      | [service/line](https://github.com/nikoksr/notify/tree/main/service/line)             | [line/line-bot-sdk-go](https://github.com/line/line-bot-sdk-go)                                 |
| [Line Notify](https://notify-bot.line.me)                    | [service/line](https://github.com/nikoksr/notify/tree/main/service/line)             | [utahta/go-linenotify](https://github.com/utahta/go-linenotify)                                 |
| [Mailgun](https://www.mailgun.com)                           | [service/mailgun](https://github.com/nikoksr/notify/tree/main/service/mailgun)       | [mailgun/mailgun-go](https://github.com/mailgun/mailgun-go)                                     |
| [Microsoft Teams](https://www.microsoft.com/microsoft-teams) | [service/msteams](https://github.com/nikoksr/notify/tree/main/service/msteams)       | [atc0005/go-teams-notify](https://github.com/atc0005/go-teams-notify)                           |
| [Plivo](https://www.plivo.com)                               | [service/plivo](https://github.com/nikoksr/notify/tree/main/service/plivo)           | [plivo/plivo-go](https://github.com/plivo/plivo-go)                                             |
| [Pushbullet](https://www.pushbullet.com)                     | [service/pushbullet](https://github.com/nikoksr/notify/tree/main/service/pushbullet) | [cschomburg/go-pushbullet](https://github.com/cschomburg/go-pushbullet)                         |
| [RocketChat](https://rocket.chat)                            | [service/rocketchat](https://github.com/nikoksr/notify/tree/main/service/rocketchat) | [RocketChat/Rocket.Chat.Go.SDK](https://github.com/RocketChat/Rocket.Chat.Go.SDK)               |
| [SendGrid](https://sendgrid.com)                             | [service/sendgrid](https://github.com/nikoksr/notify/tree/main/service/sendgrid)     | [sendgrid/sendgrid-go](https://github.com/sendgrid/sendgrid-go)                                 |
| [Slack](https://slack.com)                                   | [service/slack](https://github.com/nikoksr/notify/tree/main/service/slack)           | [slack-go/slack](https://github.com/slack-go/slack)                                             |
| [Syslog](https://wikipedia.org/wiki/Syslog)                  | [service/syslog](https://github.com/nikoksr/notify/tree/main/service/syslog)         | [log/syslog](https://pkg.go.dev/log/syslog)                                                     |
| [Telegram](https://telegram.org)                             | [service/telegram](https://github.com/nikoksr/notify/tree/main/service/telegram)     | [go-telegram-bot-api/telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) |
| [TextMagic](https://www.textmagic.com)                       | [service/textmagic](https://github.com/nikoksr/notify/tree/main/service/textmagic)   | [textmagic/textmagic-rest-go-v2](https://github.com/textmagic/textmagic-rest-go-v2)             |
| [Twitter](https://twitter.com)                               | [service/twitter](https://github.com/nikoksr/notify/tree/main/service/twitter)       | [dghubble/go-twitter](https://github.com/dghubble/go-twitter)                                   |
| [WeChat](https://www.wechat.com)                             | [service/wechat](https://github.com/nikoksr/notify/tree/main/service/wechat)         | [silenceper/wechat](https://github.com/silenceper/wechat)                                       |
| [WhatsApp](https://www.whatsapp.com)                         | [service/whatsapp](https://github.com/nikoksr/notify/tree/main/service/whatsapp)     | [Rhymen/go-whatsapp](https://github.com/Rhymen/go-whatsapp)                                     |

## Contributing <a id="contributing"></a>

Contributions, issues and feature requests are very welcome! Feel free to check [issues page](https://github.com/nikoksr/notify/issues). Please also take a look at the [contribution guidelines](https://github.com/nikoksr/notify/blob/main/CONTRIBUTING.md).

## Show your support <a id="support"></a>

Give a ⭐️ if you like this project!
