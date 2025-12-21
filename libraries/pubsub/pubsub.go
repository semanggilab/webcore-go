package pubsub

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/semanggilab/webcore-go/app/config"
	"github.com/semanggilab/webcore-go/app/loader"
	"github.com/semanggilab/webcore-go/app/logger"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// PubSub represents shared Google PubSub connection
type PubSub struct {
	Client    *pubsub.Client
	Config    config.PubSubConfig
	Receivers []loader.PubSubReceiver
}

// NewPubSub creates a new PubSub connection
func NewPubSub(ctx context.Context, config config.PubSubConfig) (*PubSub, error) {
	var client *pubsub.Client
	var err error

	if config.ProjectID == "" || config.Topic == "" || config.Subscription == "" {
		return nil, fmt.Errorf("PubSub config project_id, topic, and subscription cannot be empty")
	}

	// Configure PubSub client options
	opts := []option.ClientOption{}

	if config.EmulatorHost != "" {
		opts = append(opts, option.WithEndpoint(config.EmulatorHost), option.WithoutAuthentication())
	}

	if config.Credentials != "" {
		// In a real implementation, you would load credentials from the provided string
		// For now, we'll use the default credentials
		opts = append(opts, option.WithCredentialsFile(config.Credentials))
	}

	client, err = pubsub.NewClient(ctx, config.ProjectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create PubSub client: %v", err)
	}

	return &PubSub{
		Client:    client,
		Config:    config,
		Receivers: []loader.PubSubReceiver{},
	}, nil
}

func (ps *PubSub) Install(args ...any) error {
	// Tidak melakukan apa-apa
	return nil
}

func (ps *PubSub) Connect() error {
	// Tidak melakukan apa-apa proses konek hanya dilakukan saat di mode consumer pull message atau publish message di mode producer
	return nil
}

// Close closes the PubSub connection
func (ps *PubSub) Close() error {
	if ps.Client != nil {
		return ps.Client.Close()
	}
	return nil
}

func (ps *PubSub) Uninstall() error {
	// Tidak melakukan apa-apa
	return nil
}

func (ps *PubSub) RegisterReceiver(receiver loader.PubSubReceiver) {
	ps.Receivers = append(ps.Receivers, receiver)
}

func (ps *PubSub) StartReceiving(ctx context.Context) {
	// Ensure topic and subscription exist
	if !ps.EnsureTopicExists(ctx) {
		logger.Error("Topic tidak ditemukan")
		return
	}

	if !ps.EnsureSubscriptionExists(ctx) {
		logger.Error("Subscription tidak ditemukan")
		return
	}

	if len(ps.Receivers) == 0 {
		return
	}

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		// Main consumer loop
		for {
			select {
			case <-ctx.Done():
				logger.Info("Consumer context cancelled, shutting down...")
				return
			case <-signalChan:
				logger.Info("Received termination signal, shutting down...")
				return
			default:
				// Pull messages from PubSub
				messages, err := ps.PullMessages(ctx, ps.Config.MaxMessagesPerPull)
				if err != nil {
					logger.Error("Error pulling messages", "error", err)
					// Wait before retrying
					time.Sleep(ps.Config.SleepTimeBetweenPulls)
					continue
				}

				if len(messages) == 0 {
					// No messages received, wait before pulling again
					time.Sleep(ps.Config.SleepTimeBetweenPulls)
					continue
				}

				logger.Debug("Received messages", "count", len(messages))

				// pubsubMsgs := map[string]pubsub.Message{}
				msgs := []*loader.PubSubMessage{}
				for _, msg := range messages {
					// pubsubMsgs[msg.ID] = *msg
					msgs = append(msgs, &loader.PubSubMessage{
						ID:         msg.ID,
						Data:       msg.Data,
						Attributes: msg.Attributes,
					})
				}

				for _, c := range ps.Receivers {
					go c.Consume(ctx, msgs)
					// go func() {
					// 	status, err := c.Consume(ctx, msgs)
					// 	if err != nil {
					// 		logger.Error("Error consuming messages", "error", err)
					// 	}

					// 	for id, sts := range status {
					// 		if sts {

					// 		} else {
					// 			// logger.Debug("Error consuming messages", "error", err)
					// 		}

					// 		pmsg := pubsubMsgs[id]
					// 		pmsg.Ack()
					// 	}

					// }()
				}
			}
		}
	}()
}

// EnsureTopicExists checks if the topic exists
func (ps *PubSub) EnsureTopicExists(ctx context.Context) bool {
	return ps.GetTopicInfo(ctx) != nil
}

// EnsureSubscriptionExists checks if the subscription exists
func (ps *PubSub) EnsureSubscriptionExists(ctx context.Context) bool {
	return ps.GetSubscriptionInfo(ctx) != nil
}

// GetTopicInfo returns information about the topic
func (ps *PubSub) GetTopicInfo(ctx context.Context) *pubsubpb.Topic {
	req := &pubsubpb.ListTopicsRequest{
		Project: fmt.Sprintf("projects/%s", ps.Config.ProjectID),
	}
	it := ps.Client.TopicAdminClient.ListTopics(ctx, req)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			continue
		}

		slog.Debug("Found topic", "name", topic.Name)
		if topic.Name == fmt.Sprintf("projects/%s/topics/%s", ps.Config.ProjectID, ps.Config.Topic) {
			return topic
		}
	}

	return nil
}

// GetSubscriptionInfo returns information about the subscription
func (ps *PubSub) GetSubscriptionInfo(ctx context.Context) *pubsubpb.Subscription {
	exists := ps.ListSubscriptions(ctx)

	for _, sub := range exists {
		if sub != nil && sub.Name == fmt.Sprintf("projects/%s/subscriptions/%s", ps.Config.ProjectID, ps.Config.Subscription) {
			return sub
		}
	}

	return nil
}

// ListSubscriptions lists all subscriptions for this topic
func (ps *PubSub) ListSubscriptions(ctx context.Context) []*pubsubpb.Subscription {
	var subs []*pubsubpb.Subscription
	req := &pubsubpb.ListSubscriptionsRequest{
		Project: fmt.Sprintf("projects/%s", ps.Config.ProjectID),
	}
	it := ps.Client.SubscriptionAdminClient.ListSubscriptions(ctx, req)
	for {
		s, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			continue
		}
		subs = append(subs, s)
	}
	return subs
}

// PublishMessage publishes a message to the topic
func (ps *PubSub) PublishMessage(ctx context.Context, data []byte, attributes map[string]string) (string, error) {
	publisher := ps.Client.Publisher(ps.Config.Topic)
	result := publisher.Publish(ctx, &pubsub.Message{
		Data:       data,
		Attributes: attributes,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	msgID, err := result.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to publish message: %v", err)
	}

	logger.Debug("Published message with msgID %s", msgID)
	return msgID, nil
}

// PublishMessages publishes multiple messages to the topic
func (ps *PubSub) PublishMessages(ctx context.Context, messages [][]byte, attributes map[string]string) ([]string, error) {
	results := []string{}

	var err error
	var msgID string
	i := 0
	for _, msg := range messages {
		msgID, err = ps.PublishMessage(ctx, msg, attributes)
		if msgID != "" {
			results[i] = msgID
			i++
		}
	}

	return results, err
}

// PullMessages pulls messages from the subscription
func (ps *PubSub) PullMessages(ctx context.Context, maxMessages int) ([]*pubsub.Message, error) {
	subscriber := ps.Client.Subscriber(ps.Config.Subscription)
	subscriber.ReceiveSettings.MaxOutstandingMessages = maxMessages

	messages := make([]*pubsub.Message, 0)

	err := subscriber.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		messages = append(messages, msg)
		msg.Ack() // Acknowledge the message
	})

	if err != nil {
		return nil, fmt.Errorf("failed to pull messages: %v", err)
	}

	return messages, nil
}
