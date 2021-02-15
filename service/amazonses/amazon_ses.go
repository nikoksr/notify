package amazonses

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
)

// AmazonSES struct holds necessary data to communicate with the Amazon Simple Email Service API.
type AmazonSES struct {
	client            *ses.SES
	senderAddress     string
	receiverAddresses []*string
}

// New returns a new instance of a AmazonSES notification service.
// You will need an Amazon Simple Email Service API access key and secret.
// See https://docs.aws.amazon.com/ses/latest/DeveloperGuide/get-aws-keys.html
func New(accessKeyID, secretKey, region, senderAddress string) (*AmazonSES, error) {
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretKey, ""),
		Region:      aws.String(region),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &AmazonSES{
		client:            ses.New(sess),
		senderAddress:     senderAddress,
		receiverAddresses: []*string{},
	}, nil
}

// AddReceivers takes email addresses and adds them to the internal address list. The Send method will send
// a given message to all those addresses.
func (a *AmazonSES) AddReceivers(addresses ...string) {
	for _, address := range addresses {
		a.receiverAddresses = append(a.receiverAddresses, aws.String(address))
	}
}

// Send takes a message subject and a message body and sends them to all previously set chats. Message body supports
// html as markup language.
func (a AmazonSES) Send(subject, message string) error {
	input := &ses.SendEmailInput{}
	input.SetSource(a.senderAddress)
	input.SetDestination(&ses.Destination{
		ToAddresses: a.receiverAddresses,
	})
	input.SetMessage(&ses.Message{
		Body: &ses.Body{
			Html: &ses.Content{
				Data: aws.String(message),
			},
			// Text: &ses.Content{
			//     Data:    aws.String(message),
			// },
		},
		Subject: &ses.Content{
			Data: aws.String(subject),
		},
	})

	_, err := a.client.SendEmail(input)
	if err != nil {
		return errors.Wrap(err, "failed to send mail using Amazon SES service")
	}

	return nil
}
