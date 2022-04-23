# Bark

[Bark](https://apps.apple.com/us/app/bark-customed-notifications/id1403753865) is an application allows you to push customed notifications to your iPhone

## Usage

```go
// Create a bark service. `device key` is generated when you install the application
barkService := bark.New("your bark device key", bark.DefaultServer)

// Tell our notifier to use the bark service. You can repeat the above process
// for as many services as you like and just tell the notifier to use them.
notify.UseServices(barkService)

// Send a test message.
_ = notify.Send(
	context.Background(),
	"Subject/Title",
	"The actual message - Hello, you awesome gophers! :)",
)
```

