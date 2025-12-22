package kafka

import (
	"github.com/semanggilab/webcore-go/app/config"
	"github.com/semanggilab/webcore-go/app/loader"
)

type KafkaConsumerLoader struct{}

func (l *KafkaConsumerLoader) ClassName() string {
	return "KafkaConsumer"
}

func (l *KafkaConsumerLoader) Init(args ...any) (loader.Library, error) {
	config := args[0].(config.KafkaConfig)
	receiver := args[1].(KafkaReceiver)

	kc, err := NewKafkaConsumer(&config, receiver)
	if err != nil {
		return nil, err
	}

	err = kc.Install(args...)
	if err != nil {
		return nil, err
	}

	kc.Connect()
	return kc, nil
}

type KafkaProducerLoader struct{}

func (l *KafkaProducerLoader) ClassName() string {
	return "KafkaProducer"
}

func (l *KafkaProducerLoader) Init(args ...any) (loader.Library, error) {
	config := args[0].(config.KafkaConfig)

	kc, err := NewKafkaProducer(&config)
	if err != nil {
		return nil, err
	}

	err = kc.Install(args...)
	if err != nil {
		return nil, err
	}

	kc.Connect()
	return kc, nil
}
