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

#### Recommendation <a id="recommendation"></a>

In this example, we use the global `Send()` function. Similar to most logging libraries such as
[zap](https://github.com/uber-go/zap), we provide global functions for convenience. However, as with most logging
libraries, we also recommend avoiding the use of global functions as much as possible. Instead, use one of our versatile
constructor functions to create a new local `Notify` instance and pass it down the stream.

Read the [library docs](https://pkg.go.dev/github.com/nikoksr/notify#section-documentation) for more information.

## Contributing <a id="contributing"></a>

Yes, please! Contributions of all kinds are very welcome! Feel free to check our [open issues](https://github.com/nikoksr/notify/issues). Please also take a look at the [contribution guidelines](https://github.com/nikoksr/notify/blob/main/CONTRIBUTING.md).

> Psst, don't forget to check the list of [missing services](https://github.com/nikoksr/notify/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3Aaffects%2Fservices+label%3A%22help+wanted%22+no%3Aassignee) waiting to be added by you or create [a new issue](https://github.com/nikoksr/notify/issues/new?assignees=&labels=affects%2Fservices%2C+good+first+issue%2C+hacktoberfest%2C+help+wanted%2C+type%2Fenhancement%2C+up+for+grabs&template=service-request.md&title=feat%28service%29%3A+Add+%5BSERVICE+NAME%5D+service) if you want a new service to be added.

## Supported services <a id="supported_services"></a>

> Click [here](https://github.com/nikoksr/notify/issues/new?assignees=&labels=affects%2Fservices%2C+good+first+issue%2C+hacktoberfest%2C+help+wanted%2C+type%2Fenhancement%2C+up+for+grabs&template=service-request.md&title=feat%28service%29%3A+Add+%5BSERVICE+NAME%5D+service) to request a missing service.

| Service                                                                        | Path                                     | Credits                                                                                         |
|--------------------------------------------------------------------------------|------------------------------------------|-------------------------------------------------------------------------------------------------|
| [Amazon SES](https://aws.amazon.com/ses)                                       | [service/amazonses](service/amazonses)   | [aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)                                       |
| [Amazon SNS](https://aws.amazon.com/sns)                                       | [service/amazonsns](service/amazonsns)   | [aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)                                       |
| [Bark](https://apps.apple.com/us/app/bark-customed-notifications/id1403753865) | [service/bark](service/bark)             | -                                                                                               |
| [DingTalk](https://www.dingtalk.com)                                           | [service/dinding](service/dingding)      | [blinkbean/dingtalk](https://github.com/blinkbean/dingtalk)                                     |
| [Discord](https://discord.com)                                                 | [service/discord](service/discord)       | [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)                                     |
| [Email](https://wikipedia.org/wiki/Email)                                      | [service/mail](service/mail)             | [jordan-wright/email](https://github.com/jordan-wright/email)                                   |
| [Firebase Cloud Messaging](https://firebase.google.com/docs/cloud-messaging)   | [service/fcm](service/fcm)               | [appleboy/go-fcm](https://github.com/appleboy/go-fcm)                                           |
| [Lark](https://www.larksuite.com/)                                             | [service/lark](service/lark)             | [go-lark/lark](https://github.com/go-lark/lark)                                                 |
| [Line](https://line.me)                                                        | [service/line](service/line)             | [line/line-bot-sdk-go](https://github.com/line/line-bot-sdk-go)                                 |
| [Line Notify](https://notify-bot.line.me)                                      | [service/line](service/line)             | [utahta/go-linenotify](https://github.com/utahta/go-linenotify)                                 |
| [Mailgun](https://www.mailgun.com)                                             | [service/mailgun](service/mailgun)       | [mailgun/mailgun-go](https://github.com/mailgun/mailgun-go)                                     |
| [Microsoft Teams](https://www.microsoft.com/microsoft-teams)                   | [service/msteams](service/msteams)       | [atc0005/go-teams-notify](https://github.com/atc0005/go-teams-notify)                           |
| [Plivo](https://www.plivo.com)                                                 | [service/plivo](service/plivo)           | [plivo/plivo-go](https://github.com/plivo/plivo-go)                                             |
| [Pushbullet](https://www.pushbullet.com)                                       | [service/pushbullet](service/pushbullet) | [cschomburg/go-pushbullet](https://github.com/cschomburg/go-pushbullet)                         |
| [RocketChat](https://rocket.chat)                                              | [service/rocketchat](service/rocketchat) | [RocketChat/Rocket.Chat.Go.SDK](https://github.com/RocketChat/Rocket.Chat.Go.SDK)               |
| [SendGrid](https://sendgrid.com)                                               | [service/sendgrid](service/sendgrid)     | [sendgrid/sendgrid-go](https://github.com/sendgrid/sendgrid-go)                                 |
| [Slack](https://slack.com)                                                     | [service/slack](service/slack)           | [slack-go/slack](https://github.com/slack-go/slack)                                             |
| [Syslog](https://wikipedia.org/wiki/Syslog)                                    | [service/syslog](service/syslog)         | [log/syslog](https://pkg.go.dev/log/syslog)                                                     |
| [Telegram](https://telegram.org)                                               | [service/telegram](service/telegram)     | [go-telegram-bot-api/telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) |
| [TextMagic](https://www.textmagic.com)                                         | [service/textmagic](service/textmagic)   | [textmagic/textmagic-rest-go-v2](https://github.com/textmagic/textmagic-rest-go-v2)             |
| [Twilio](https://www.twilio.com/)                                              | [service/twilio](service/twilio)         | [kevinburke/twilio-go](https://github.com/kevinburke/twilio-go)                                 |
| [Twitter](https://twitter.com)                                                 | [service/twitter](service/twitter)       | [dghubble/go-twitter](https://github.com/dghubble/go-twitter)                                   |
| [WeChat](https://www.wechat.com)                                               | [service/wechat](service/wechat)         | [silenceper/wechat](https://github.com/silenceper/wechat)                                       |
| [WhatsApp](https://www.whatsapp.com)                                           | [service/whatsapp](service/whatsapp)     | [Rhymen/go-whatsapp](https://github.com/Rhymen/go-whatsapp)                                     |
| [Matrix](https://www.matrix.org)                                               | [service/matrix](service/matrix)         | [mautrix/go](https://github.com/mautrix/go)                                                     |


## Special Thanks <a id="special_thanks"></a>

### Maintainers <a id="maintainers"></a>

- [@svaloumas](https://github.com/svaloumas)

### Logo <a id="logo"></a>

The [logo](https://github.com/MariaLetta/free-gophers-pack) was made by the amazing [MariaLetta](https://github.com/MariaLetta).

## Similar projects <a id="similar_projects"></a>

> Just to clarify, Notify was not inspired by any other project. I created it as a tiny subpackage of a larger project and only later decided to make it a standalone project. In this section I just want to mention other great projects.

  - [containrrr/shoutrrr](https://github.com/containrrr/shoutrrr)
  - [caronc/apprise](https://github.com/caronc/apprise)

## Show your support <a id="support"></a>

Give a ⭐️ if you like this project!
