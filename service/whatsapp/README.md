# WhatsApp

WhatsApp notification service using [tulir/whatsmeow](https://github.com/tulir/whatsmeow).

## Usage

```go
wa := whatsapp.New()
wa.AddReceivers("628123456789@s.whatsapp.net")

ctx := context.Background()

// Login via QR code (scan with WhatsApp mobile app)
if err := wa.LoginWithQRCode(ctx, "whatsapp.db"); err != nil {
    log.Fatal(err)
}
defer wa.Disconnect()

// Or login via 8-digit pairing code
// if err := wa.LoginWithPairingCode(ctx, "628123456789", "whatsapp.db"); err != nil {
//     log.Fatal(err)
// }

notifier := notify.New()
notifier.UseServices(wa)
notifier.Send(ctx, "Subject", "Message body")
JID Format
Personal chat: 628123456789@s.whatsapp.net
Group chat: 1234567890-1234567890@g.us
Session
Session data is stored in a SQLite database file. Re-authentication is only needed once per device.
