## Contributing to notify
We want to make contributing to this project as easy and transparent as
possible.

## Project structure

- `service` - Contains definitions for the underlying notification services.
  - `service/discord` - Discord notification service.
  - `service/mail` - Email notification service.
  - `service/msteams` - Microsoft Teams notification service.
  - `service/slack` - Slack notification service.
  - `service/telegram` - Telegram notification service

## Commits

Commit messages should be well formatted, and to make that "standardized", we
are using Conventional Commits.

You can follow the documentation on
[their website](https://www.conventionalcommits.org).

## Pull Requests
We actively welcome your pull requests.

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes (`make test`).
5. Make sure your code lints (`make lint`).
6. Make sure your code is well formatted (`make fmt`).

## Issues
We use GitHub issues to track public bugs. Please ensure your description is
clear and has sufficient instructions to be able to reproduce the issue.

## License
By contributing to notify, you agree that your contributions will be licensed
under the LICENSE file in the root directory of this source tree.
