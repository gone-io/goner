package mongo

import (
	"context"
	"fmt"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/g"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"
	"time"
)

var clientMap sync.Map

type Config struct {
	URI                    string        `json:"uri"`
	Database               string        `json:"database"`
	Username               string        `json:"username"`
	Password               string        `json:"password"`
	AuthSource             string        `json:"authSource"`
	MaxPoolSize            uint64        `json:"maxPoolSize"`
	MinPoolSize            uint64        `json:"minPoolSize"`
	MaxConnIdleTime        time.Duration `json:"maxConnIdleTime"`
	ConnectTimeout         time.Duration `json:"connectTimeout"`
	SocketTimeout          time.Duration `json:"socketTimeout"`
	ServerSelectionTimeout time.Duration `json:"serverSelectionTimeout"`
}

func (config Config) ToMongoOptions() *options.ClientOptions {
	opts := options.Client()

	if config.URI != "" {
		opts.ApplyURI(config.URI)
	}

	if config.Username != "" && config.Password != "" {
		credential := options.Credential{
			Username: config.Username,
			Password: config.Password,
		}
		if config.AuthSource != "" {
			credential.AuthSource = config.AuthSource
		}
		opts.SetAuth(credential)
	}

	if config.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(config.MaxPoolSize)
	}

	if config.MinPoolSize > 0 {
		opts.SetMinPoolSize(config.MinPoolSize)
	}

	if config.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(config.MaxConnIdleTime)
	}

	if config.ConnectTimeout > 0 {
		opts.SetConnectTimeout(config.ConnectTimeout)
	}

	if config.SocketTimeout > 0 {
		opts.SetSocketTimeout(config.SocketTimeout)
	}

	if config.ServerSelectionTimeout > 0 {
		opts.SetServerSelectionTimeout(config.ServerSelectionTimeout)
	}

	return opts
}

func provide(tagConf string, param struct {
	configure gone.Configure `gone:"configure"`
}) (*mongo.Client, error) {
	if value, ok := clientMap.Load(tagConf); ok {
		return value.(*mongo.Client), nil
	}

	var prefix = "mongo"
	_, keys := gone.TagStringParse(tagConf)
	if len(keys) > 0 && keys[0] != "" {
		prefix = strings.TrimSpace(keys[0])
	}

	var config Config
	err := param.configure.Get(prefix, &config, "")
	g.PanicIfErr(gone.ToErrorWithMsg(err, fmt.Sprintf("get %s config err", prefix)))
	client, err := mongo.Connect(context.Background(), config.ToMongoOptions())
	g.PanicIfErr(gone.ToErrorWithMsg(err, "failed to connect to MongoDB"))

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	g.PanicIfErr(gone.ToErrorWithMsg(err, "failed to ping MongoDB"))
	clientMap.Store(tagConf, client)
	return client, nil
}

// Load *mongo.Client provider
func Load(loader gone.Loader) error {
	return loader.Load(gone.WrapFunctionProvider(provide))
}
