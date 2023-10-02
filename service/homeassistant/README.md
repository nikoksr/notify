# HomeAssistant service

[HomeAssistant](https://www.home-assistant.io) is an application that handles home automation and integration

## Usage

[HomeAssistant webhook trigger](https://www.home-assistant.io/docs/automation/trigger/#webhook-trigger)

```go
// Create a HomeAssistant service
haService := homeassistant.New()

// Add webhook
haService.AddWebhook("https://url-to-home-assistant", "<webhook_id>", "<http_method>")

// Tell our notifier to use the service.
notify.UseServices(haService)

// Send a test message.
_ = notify.Send(
    context.Background(),
    "Subject/Title",
    "The actual message - Hello, you awesome gophers! :)",
)
```
