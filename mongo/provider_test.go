package mongo

import (
	"context"
	"github.com/gone-io/gone/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConfig_ToMongoOptions(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		check  func(*testing.T, *options.ClientOptions)
	}{
		{
			name: "basic config with URI",
			config: Config{
				URI:      "mongodb://localhost:27017",
				Database: "testdb",
			},
			check: func(t *testing.T, opts *options.ClientOptions) {
				assert.NotNil(t, opts)
			},
		},
		{
			name: "config with authentication",
			config: Config{
				URI:        "mongodb://localhost:27017",
				Username:   "testuser",
				Password:   "testpass",
				AuthSource: "admin",
			},
			check: func(t *testing.T, opts *options.ClientOptions) {
				assert.NotNil(t, opts)
				auth := opts.Auth
				assert.NotNil(t, auth)
				assert.Equal(t, "testuser", auth.Username)
				assert.Equal(t, "testpass", auth.Password)
				assert.Equal(t, "admin", auth.AuthSource)
			},
		},
		{
			name: "config with pool settings",
			config: Config{
				URI:             "mongodb://localhost:27017",
				MaxPoolSize:     100,
				MinPoolSize:     10,
				MaxConnIdleTime: 30 * time.Minute,
			},
			check: func(t *testing.T, opts *options.ClientOptions) {
				assert.NotNil(t, opts)
				assert.Equal(t, uint64(100), *opts.MaxPoolSize)
				assert.Equal(t, uint64(10), *opts.MinPoolSize)
				assert.Equal(t, 30*time.Minute, *opts.MaxConnIdleTime)
			},
		},
		{
			name: "config with timeouts",
			config: Config{
				URI:                    "mongodb://localhost:27017",
				ConnectTimeout:         10 * time.Second,
				SocketTimeout:          30 * time.Second,
				ServerSelectionTimeout: 5 * time.Second,
			},
			check: func(t *testing.T, opts *options.ClientOptions) {
				assert.NotNil(t, opts)
				assert.Equal(t, 10*time.Second, *opts.ConnectTimeout)
				assert.Equal(t, 30*time.Second, *opts.SocketTimeout)
				assert.Equal(t, 5*time.Second, *opts.ServerSelectionTimeout)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.config.ToMongoOptions()
			tt.check(t, opts)
		})
	}
}

func TestLoad(t *testing.T) {
	_ = os.Setenv("GONE_MONGO", `{"uri":"mongodb://root:example@127.0.0.1:27017/"}`)
	defer func() {
		_ = os.Unsetenv("GONE_MONGO_URI")
	}()

	gone.
		NewApp(Load).
		Test(func(client *mongo.Client, i struct {
			client  *mongo.Client `gone:"*"`
			client2 *mongo.Client `gone:"*,test"`
		}) {
			assert.NotNil(t, client)
			assert.Equal(t, client, i.client)

			collection := client.Database("myapp").Collection("users")

			type User struct {
				ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
				Name      string             `bson:"name" json:"name"`
				Email     string             `bson:"email" json:"email"`
				Age       int                `bson:"age" json:"age"`
				CreatedAt time.Time          `bson:"created_at" json:"created_at"`
				UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
			}

			user := &User{
				Name:      "gone",
				Email:     "gone@goner.fun",
				Age:       4,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			result, err := collection.InsertOne(context.Background(), user)
			if err != nil {
				t.Fatalf("Failed to create user: %v", err)
			}
			user.ID = result.InsertedID.(primitive.ObjectID)

			var updatedUser User

			err = collection.FindOne(context.Background(), bson.M{"email": "gone@goner.fun"}).Decode(&updatedUser)
			if err != nil {
				t.Fatalf("Failed to get user: %v", err)
			}
			assert.Equal(t, user.Name, updatedUser.Name)

			one, err := collection.DeleteOne(context.Background(), bson.M{"_id": user.ID})
			if err != nil {
				t.Fatalf("Failed to delete user: %v", err)
			}
			print(one)
		})
}
