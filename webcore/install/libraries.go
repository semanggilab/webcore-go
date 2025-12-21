package install

import (
	"github.com/semanggilab/webcore-go/app/core"
	"github.com/semanggilab/webcore-go/lib/postgres"
	"github.com/semanggilab/webcore-go/lib/pubsub"
	"github.com/semanggilab/webcore-go/lib/redis"
)

var APP_LIBRARIES = map[string]core.LibraryLoader{
	"db:postgres": &postgres.PostgresLoader{},
	"redis":       &redis.RedisLoader{},
	"pubsub":      &pubsub.PubSubLoader{},

	// Add your library here
}
