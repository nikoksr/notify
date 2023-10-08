# APNS2

[APNS2](https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/sending_notification_requests_to_apns) sending notifications through the apns server directly to iphone devices

## Usage

```go
// Create a apns2 service. `service.p12` or `service.pem` is generated when you install the application.
apns2Service := apns2.New(P12File("/cert/service.p12",""),"<com.myapp.topic>")
apns2Service = apns2.New(P12Bytes([]byte{},""),"<com.myapp.topic>")
apns2Service = apns2.New(PemFile("/cert/service.pem",""),"<com.myapp.topic>")
apns2Service = apns2.New(PemBytes([]byte{},""),"<com.myapp.topic>")

// Add devices
apns2Service.AddReceivers("<token1>","<token2>")

// Tell our notifier to use the apns2 service.
notify.UseServices(apns2Service)

// Send a test message.
_ = notify.Send(
    context.Background(),
    "Subject/Title",
    "The actual message - Hello, you awesome gophers! :)",
)
```

