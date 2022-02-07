# Changelog
All notable changes to this project will be documented in this file.

## [0.20.1] - 2022-02-07

### Dependencies

- Update and tidy dependencies (#206)

## [0.20.0] - 2022-02-07

### Documentation

- Add credits for AWS SNS service

### Features

- Add AWS SNS (#110)

### Miscellaneous Tasks

- Fix minor typo

## [0.19.1] - 2022-02-07

### Dependencies

- Update and tidy dependencies (#202)

## [0.19.0] - 2022-02-07

### CI

- Add new pull-request template
- Fix goimports local prefix

### Documentation

- Add build commit parser
- Apply git-cliff's config update
- Simplify contribution guide
- Add DingTalk credits and remove FOSSA badge

### Features

- Add DingTalk (#183)

## [0.18.0] - 2022-02-07

### Build

- Fix linter command and add implicit install commands
- Add commands to generate changelogs

### CI

- Reintroduce config for readme badge
- Add git-cliff config
- Change version signature of workflow 'ci'
- Add workflow for git-cliff to auto-gen changelogs
- Fix smaller issues and clean-up
- Switch from config-flag to config option in action file
- Drop changelog generating action

### Dependencies

- Update module github.com/aws/aws-sdk-go-v2 to v1.11.1
- Update module github.com/aws/aws-sdk-go-v2/config to v1.10.2
- Update module github.com/aws/aws-sdk-go-v2/service/ses to v1.9.1
- Update module github.com/aws/aws-sdk-go-v2/config to v1.10.3
- Update module github.com/aws/aws-sdk-go-v2 to v1.11.2
- Update module github.com/aws/aws-sdk-go-v2/credentials to v1.6.4
- Update module github.com/aws/aws-sdk-go-v2/service/ses to v1.9.2
- Update module github.com/aws/aws-sdk-go-v2/config to v1.11.0
- Update module github.com/silenceper/wechat/v2 to v2.1.0
- Update module github.com/aws/aws-sdk-go-v2/config to v1.11.1
- Update module github.com/aws/aws-sdk-go-v2/service/ses to v1.10.0
- Update module github.com/slack-go/slack to v0.10.1
- Update module github.com/aws/aws-sdk-go-v2 to v1.12.0
- Update module github.com/aws/aws-sdk-go-v2/service/ses to v1.11.0
- Update module github.com/dghubble/oauth1 to v0.7.1
- Update module github.com/aws/aws-sdk-go-v2/config to v1.12.0
- Update module github.com/aws/aws-sdk-go-v2 to v1.13.0
- Update module github.com/aws/aws-sdk-go-v2/config to v1.13.0
- Update module github.com/aws/aws-sdk-go-v2/service/ses to v1.12.0
- Update module github.com/aws/aws-sdk-go-v2/config to v1.13.1

## [0.17.4] - 2021-11-17

### CI

- Drop stale config & bot
- Add codeql-analysis config
- Remove dependabot and extend renovate config
- Update github action
- Remove config

### Dependencies

- Update module github.com/go-telegram-bot-api/telegram-bot-api to v5
- Update module github.com/line/line-bot-sdk-go to v7.10.0
- Revert version updates for line and telegram
- Update module github.com/aws/aws-sdk-go-v2 to v1.8.1
- Update module github.com/aws/aws-sdk-go-v2/service/ses to v1.5.2
- Update module github.com/aws/aws-sdk-go-v2 to v1.9.0
- Update module github.com/aws/aws-sdk-go-v2/service/ses to v1.6.0
- Update module github.com/mailgun/mailgun-go/v4 to v4.5.3
- Update module github.com/silenceper/wechat/v2 to v2.0.9
- Update module github.com/slack-go/slack to v0.9.5
- Update module github.com/aws/aws-sdk-go-v2 to v1.9.1
- Update module github.com/aws/aws-sdk-go-v2/service/ses to v1.6.1
- Update module github.com/sendgrid/sendgrid-go to v3.10.1
- Update module github.com/aws/aws-sdk-go-v2/service/ses to v1.7.0
- Update github.com/dghubble/go-twitter commit hash to ad02880
- Update module github.com/sendgrid/sendgrid-go to v3.10.3
- Update module github.com/aws/aws-sdk-go-v2/config to v1.9.0
- Update github.com/dghubble/go-twitter commit hash to 93a8679
- Update module github.com/aws/aws-sdk-go-v2 to v1.11.0
- Update module github.com/aws/aws-sdk-go-v2/credentials to v1.6.1
- Update module github.com/aws/aws-sdk-go-v2/config to v1.10.1
- Update module github.com/aws/aws-sdk-go-v2/service/ses to v1.9.0
- Update module github.com/mailgun/mailgun-go/v4 to v4.6.0
- Update module github.com/slack-go/slack to v0.10.0
- Update go mod version and run go mod tidy

## [0.17.3] - 2021-08-13

### Dependencies

- Bump github.com/aws/aws-sdk-go-v2/service/ses
- Bump github.com/slack-go/slack from 0.9.2 to 0.9.4
- Bump github.com/aws/aws-sdk-go-v2/config
- Bump github.com/aws/aws-sdk-go-v2/credentials
- Bump github.com/textmagic/textmagic-rest-go-v2/v2
- Go get -u & go mod tidy

## [0.17.2] - 2021-07-14

### Dependencies

- Go mod tidy

## [0.17.1] - 2021-07-14

### Build

- Simplify Makefile

### CI

- Refactor github action for linting, testing & building

### Dependencies

- Bump github.com/textmagic/textmagic-rest-go-v2/v2
- Bump github.com/aws/aws-sdk-go-v2/service/ses
- Bump github.com/aws/aws-sdk-go-v2/config
- Bump github.com/atc0005/go-teams-notify/v2
- Bump github.com/slack-go/slack from 0.9.1 to 0.9.2
- Bump github.com/mailgun/mailgun-go/v4 from 4.5.1 to 4.5.2
- Bump github.com/plivo/plivo-go

### Documentation

- Change calling order in notify example
- Reorder supported services list alphabetically

### Miscellaneous Tasks

- Remove go generate directives for installing deps

## [0.17.0] - 2021-05-26

### Bug Fixes

- Fix : resolved issues


### Build

- Update golangci-lint config
- Replaced deprecated linter

### Dependencies

- Bump github.com/aws/aws-sdk-go-v2 from 1.3.4 to 1.6.0
- Bump github.com/aws/aws-sdk-go-v2/config
- Bump github.com/aws/aws-sdk-go-v2/service/ses
- Bump github.com/slack-go/slack from 0.9.0 to 0.9.1
- Bump github.com/sendgrid/sendgrid-go
- Bump github.com/silenceper/wechat/v2 from 2.0.5 to 2.0.6
- Bump github.com/aws/aws-sdk-go-v2/service/ses
- Update and tidy modules
- Go mod tidy

### Features

- Add textmagic service
- Formatted code
- Formatted code

### Miscellaneous Tasks

- Run 'make fmt'

### Add

- Textmagic service

## [0.16.1] - 2021-04-27

### Dependencies

- Bump github.com/plivo/plivo-go
- Bump github.com/atc0005/go-teams-notify/v2
- Bump github.com/mailgun/mailgun-go/v4 from 4.4.1 to 4.5.1
- Bump github.com/sendgrid/sendgrid-go
- Go get -u ./... & go mod tidy

## [0.16.0] - 2021-04-27

### Build

- Build-passing commit

- Build-passing commit


### Dependencies

- Bump github.com/aws/aws-sdk-go-v2/service/ses
- Bump actions/cache from v2.1.4 to v2.1.5
- Bump github.com/aws/aws-sdk-go-v2/config
- Bump github.com/aws/aws-sdk-go-v2/credentials
- Bump github.com/slack-go/slack from 0.8.1 to 0.9.0
- Bump github.com/aws/aws-sdk-go-v2 from 1.2.1 to 1.3.4

### Features

- Added support for WeChat

## [0.15.0] - 2021-03-16

### Bug Fixes

- Lint code
- Fix comments in code change variable name
- Update variable names
- Lint code

### Dependencies

- Bump github.com/aws/aws-sdk-go-v2/service/ses
- Bump github.com/aws/aws-sdk-go-v2/credentials
- Bump github.com/aws/aws-sdk-go-v2 from 1.2.0 to 1.2.1
- Bump github.com/aws/aws-sdk-go-v2/config
- Bump github.com/aws/aws-sdk-go-v2/config from 1.1.1 to 1.1.2 (#60)
- Update and tidy go dependencies
- Go mod tidy

### Documentation

- Add rocketchat service to readme
- Add rocketchat service usage example

### Features

- Add rocketchat service

## [0.14.0] - 2021-03-09

### Bug Fixes

- Adjust code style

### Dependencies

- Bump github.com/mailgun/mailgun-go/v4 from 4.3.4 to 4.4.1
- Bump github.com/mailgun/mailgun-go/v4 from 4.3.4 to 4.4.1 (#54)

### Documentation

- Add line service to readme
- Add line service example usage

### Features

- Add line service

## [0.13.0] - 2021-02-23

### Dependencies

- Go get -u & go mod tidy

### Documentation

- Fix wrong type name in comment

### Refactor

- Add context.Context to parameter list of Send method
- Return context.Err when context.IsDone
- Fix faulty error return
- Add context.Context to parameter list of Send method (#51)

## [0.12.0] - 2021-02-17

## [0.11.0] - 2021-02-17

### Dependencies

- Go mod tidy
- Bump github.com/slack-go/slack from 0.8.0 to 0.8.1
- Bump github.com/slack-go/slack from 0.8.0 to 0.8.1 (#46)
- Fix go.mod
- Fix go.sum

### Documentation

- Reword some paragraphs and add missing libs to credits list
- Update README with new supported Amazon SES service

### Features

- Add WhatsApp service
- Add Amazon SES service

### Miscellaneous Tasks

- Make fmt

### Refactor

- Use newest version v2 of AWS Golang SDK
- Separate out login into LoginWithQRCode and LoginWithSessionCredentials
- Construct email input on a single initialization
- Make changes as per review comments

## [0.10.0] - 2021-02-12

## [0.9.0] - 2021-02-12

### Refactor

- Proper variable naming
- Split error checking to separate line

## [0.8.0] - 2021-02-12

### Bug Fixes

- Fix comment typo on SendGrid struct

### Documentation

- Update README with new supported SendGrid service
- Update README with new supported Mailgun service

### Features

- Add SendGrid service
- Add Mailgun service

### Miscellaneous Tasks

- Add punctuation to function comment
- Put actual type names into comments

### Refactor

- Simplify useService function
- Remove AddReceivers function from notify.Notifier interface
- Simplify telegram.AddReceivers function
- Remove unnecessary variable assignments
- Make code more compact

## [0.7.0] - 2021-02-09

### CI

- Change dependabot label

### Dependencies

- Bump actions/cache from v2 to v2.1.4
- Bump github.com/bwmarrin/discordgo from 0.23.1 to 0.23.2
- Bump actions/cache from v2 to v2.1.4 (#36)
- Bump github.com/bwmarrin/discordgo from 0.23.1 to 0.23.2 (#37)

### Features

- Add Plivo service

### Refactor

- Remove unnecessary new-lines

## [0.6.0] - 2021-02-08

## [0.5.1] - 2021-02-07

### Bug Fixes

- Fix miss-spell & run gofmt -s


### Documentation

- Fix UseService usage

### Rename

- Notify.UseService to Notify.UseServices

## [0.5.0] - 2021-02-02

### Bug Fixes

- Fix a small typo in function comment

### Build

- Add make file and go-generate directives
- Remove deprecated config option
- Add dependabot config

### CI

- Add action to test and lint any PRs
- Add go build command to gh build-action
- Add stale-bot config

### Documentation

- Update contribution guidelines
- Add fossa badge to readme

### Features

- Added Pushbullet service (tested)
- Feat(service) Pushbullet SMS - added new service and usage documentation to pushbullet package


### Patch

- Renamed PushBullet to Pushbullet
- Removed validation from pushbullet.go AddReceivers

## [0.4.0] - 2021-01-31

### Bug Fixes

- Remove unused parameter from function New

### Build

- Make codacy ignore markdown files

## [0.3.1] - 2021-01-31

### Bug Fixes

- Remove error return on method new

### Documentation

- Add notes about Microsoft Teams support

### Features

- Add support for Microsoft Teams

### Miscellaneous Tasks

- Fix grammar in method comment
- Apply changes suggested by golangci-lint

### Refactor

- Remove obsolete listener from Telegram struct
- Remove underscore from package name (golint: var-naming)

### UNTESTED

- Twitter service

## [0.3.0] - 2021-01-27

### CI

- Add templates for PR, bug- and feature-request

### Documentation

- Refactor readme
- Fix codacy badge
- Add list of supported services
- Add comments pointing to useful resources
- Update description
- Add notes about slack support

## [0.2.1] - 2021-01-26

### Bug Fixes

- Add gateway intent to discord client

### Dependencies

- Go mod tidy

### Documentation

- Add missing comment to pseudo struct
- Correct function name in comment

## [0.2.0] - 2021-01-25

### Documentation

- Correct faulty telegram chat id
- Add wip note

### Features

- Add support for various authentication methods

## [0.1.0] - 2021-01-25

### Documentation

- Add readme, coc, and contrib guidelines
- Replace gh release with tag shield

### Features

- Init project

### Miscellaneous Tasks

- Go mod tidy

### Refactor

- Comment and clean up code

<!-- generated by git-cliff -->
