package amazonsns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/pkg/errors"
)

// SNSSendMessageAPI Basic interface to send messages through SNS.
//
//go:generate mockery --name=SNSSendMessageAPI --output=. --case=underscore --inpackage
type SNSSendMessageAPI interface {
	SendMessage(ctx context.Context,
		params *sns.PublishInput,
		optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSSendMessageClient Client specific for SNS using aws sdk v2.
type SNSSendMessageClient struct {
	client *sns.Client
}

// SendMessage Client specific for SNS using aws sdk v2.
func (s SNSSendMessageClient) SendMessage(ctx context.Context,
	params *sns.PublishInput,
	optFns ...func(*sns.Options),
) (*sns.PublishOutput, error) {
	return s.client.Publish(ctx, params, optFns...)
}

// AmazonSNS Basic structure with SNS information
type AmazonSNS struct {
	sendMessageClient SNSSendMessageAPI
	queueTopics       []string
}

// New creates a new AmazonSNS
func New(accessKeyID, secretKey, region string) (*AmazonSNS, error) {
	credProvider := credentials.NewStaticCredentialsProvider(accessKeyID, secretKey, "")
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(credProvider),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	client := sns.NewFromConfig(cfg)
	return &AmazonSNS{
		sendMessageClient: SNSSendMessageClient{client: client},
	}, nil
}

// AddReceivers takes queue urls and adds them to the internal topics
// list. The Send method will send a given message to all those
// Topics.
func (s *AmazonSNS) AddReceivers(queues ...string) {
	s.queueTopics = append(s.queueTopics, queues...)
}

// Send message to everyone on all topics
func (s AmazonSNS) Send(ctx context.Context, subject, message string) error {
	// For each topic
	for _, topic := range s.queueTopics {
		// Create new input with subject, message and the specific topic
		input := &sns.PublishInput{
			Subject:  aws.String(subject),
			Message:  aws.String(message),
			TopicArn: aws.String(topic),
		}
		// Send the message
		_, err := s.sendMessageClient.SendMessage(ctx, input)
		if err != nil {
			return errors.Wrapf(err, "failed to send message using Amazon SNS to ARN TOPIC '%s'", topic)
		}
	}
	return nil
}
