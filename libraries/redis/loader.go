package redis

import (
	"github.com/semanggilab/webcore-go/app/config"
	"github.com/semanggilab/webcore-go/app/loader"
)

type RedisLoader struct {
	Redis *Redis
}

func (l *RedisLoader) ClassName() string {
	return "Redis"
}

func (l *RedisLoader) Init(args ...any) (loader.Library, error) {
	config := args[0].(config.RedisConfig)
	redis := NewRedis(config)
	err := redis.Install(args...)
	if err != nil {
		return nil, err
	}

	redis.Connect()

	l.Redis = redis
	return redis, nil
}
