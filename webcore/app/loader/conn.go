package loader

import (
	"context"
)

type Library interface {
	Install(args ...any) error
	Uninstall() error
}

type Connector interface {
	Library
	Connect() error
	Close() error
}

type DbMap map[string]any

type IDatabase interface {
	Connector

	// Common methods that both databases should implement
	Ping(ctx context.Context) error

	// Get the underlying connection for specific operations
	GetConnection() any
	GetDriver() string
	GetName() string

	Count(ctx context.Context, table string, filter DbMap) (int64, error)
	Find(ctx context.Context, table string, column []string, filter DbMap, sort map[string]int, limit int64, skip int64) ([]DbMap, error)
	FindOne(ctx context.Context, result any, table string, column []string, filter DbMap, sort map[string]int) error
	InsertOne(ctx context.Context, table string, data any) (any, error)
	Update(ctx context.Context, table string, filter DbMap, data any) (int64, error)
	UpdateOne(ctx context.Context, table string, filter DbMap, data any) (int64, error)
	Delete(ctx context.Context, table string, filter DbMap) (any, error)
	DeleteOne(ctx context.Context, table string, filter DbMap) (any, error)
}

type IRedis interface {
	Connector

	GetClient() any
}

type IPubSub interface {
	Connector

	Publish(ctx context.Context, channel string, message any, attributes map[string]string) error
	RegisterReceiver(receiver PubSubReceiver) (<-chan any, error)
}

// PubSubMessage represents a PubSub message
type PubSubMessage struct {
	ID         string
	Data       []byte
	Attributes map[string]string
}

type PubSubReceiver interface {
	Consume(ctx context.Context, messages []*PubSubMessage) (map[string]bool, error)
}

type IKafka interface {
	Connector

	Publish(ctx context.Context, topic string, message any) error
	Consume(ctx context.Context, topic string) (<-chan any, error)
}

type KafkaConsumer interface {
	Consume(ctx context.Context, message []byte) (bool, error)
}
