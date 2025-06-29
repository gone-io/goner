package main

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/gone-io/gone/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// User represents a user document
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Age       int                `bson:"age" json:"age"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// UserService provides user-related operations
type UserService struct {
	gone.Flag
	mongoClient *mongoDriver.Client `gone:"*"`
	logger      gone.Logger         `gone:"*"`
}

// CreateUser creates a new user
func (s *UserService) CreateUser(name, email string, age int) (*User, error) {
	collection := s.mongoClient.Database("myapp").Collection("users")

	user := &User{
		Name:      name,
		Email:     email,
		Age:       age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		s.logger.Errorf("Failed to create user: %v", err)
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	s.logger.Infof("User created successfully: %s (ID: %s)", name, user.ID.Hex())
	return user, nil
}

// GetUser retrieves a user by email
func (s *UserService) GetUser(email string) (*User, error) {
	collection := s.mongoClient.Database("myapp").Collection("users")

	var user User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocuments) {
			s.logger.Warnf("User not found: %s", email)
			return nil, fmt.Errorf("user not found: %s", email)
		}
		s.logger.Errorf("Failed to get user: %v", err)
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(email string, updates bson.M) error {
	collection := s.mongoClient.Database("myapp").Collection("users")

	updates["updated_at"] = time.Now()
	update := bson.M{"$set": updates}

	result, err := collection.UpdateOne(context.Background(), bson.M{"email": email}, update)
	if err != nil {
		s.logger.Errorf("Failed to update user: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		s.logger.Warnf("User not found for update: %s", email)
		return fmt.Errorf("user not found: %s", email)
	}

	s.logger.Infof("User updated successfully: %s", email)
	return nil
}

// DeleteUser deletes a user by email
func (s *UserService) DeleteUser(email string) error {
	collection := s.mongoClient.Database("myapp").Collection("users")

	result, err := collection.DeleteOne(context.Background(), bson.M{"email": email})
	if err != nil {
		s.logger.Errorf("Failed to delete user: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		s.logger.Warnf("User not found for deletion: %s", email)
		return fmt.Errorf("user not found: %s", email)
	}

	s.logger.Infof("User deleted successfully: %s", email)
	return nil
}

// ListUsers retrieves all users with pagination
func (s *UserService) ListUsers(limit, skip int64) ([]*User, error) {
	collection := s.mongoClient.Database("myapp").Collection("users")

	ctx := context.Background()
	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		s.logger.Errorf("Failed to list users: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			s.logger.Errorf("Failed to decode user: %v", err)
			continue
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		s.logger.Errorf("Cursor error: %v", err)
		return nil, err
	}

	s.logger.Infof("Retrieved %d users", len(users))
	return users, nil
}

// DemoService demonstrates MongoDB operations
type DemoService struct {
	gone.Flag
	userService *UserService `gone:"*"`
	logger      gone.Logger  `gone:"*"`
}

// RunDemo runs a demonstration of MongoDB operations
func (d *DemoService) RunDemo() error {
	d.logger.Infof("Starting MongoDB demo...")

	// Create users
	users := []struct {
		name  string
		email string
		age   int
	}{
		{"Alice Johnson", "alice@example.com", 28},
		{"Bob Smith", "bob@example.com", 32},
		{"Charlie Brown", "charlie@example.com", 25},
	}

	for _, userData := range users {
		user, err := d.userService.CreateUser(userData.name, userData.email, userData.age)
		if err != nil {
			d.logger.Errorf("Failed to create user %s: %v", userData.name, err)
			continue
		}
		d.logger.Infof("Created user: %+v", user)
	}

	// Get a user
	user, err := d.userService.GetUser("alice@example.com")
	if err != nil {
		d.logger.Errorf("Failed to get user: %v", err)
	} else {
		d.logger.Infof("Retrieved user: %+v", user)
	}

	// Update a user
	err = d.userService.UpdateUser("alice@example.com", bson.M{"age": 29})
	if err != nil {
		d.logger.Errorf("Failed to update user: %v", err)
	}

	// List users
	allUsers, err := d.userService.ListUsers(10, 0)
	if err != nil {
		d.logger.Errorf("Failed to list users: %v", err)
	} else {
		d.logger.Infof("All users: %+v", allUsers)
	}

	// Delete a user
	err = d.userService.DeleteUser("charlie@example.com")
	if err != nil {
		d.logger.Errorf("Failed to delete user: %v", err)
	}

	d.logger.Infof("MongoDB demo completed!")
	return nil
}

// AfterServerStart demonstrates using hook to run demo after server start
type AfterServerStart struct {
	gone.Flag
	afterStart gone.AfterStart `gone:"*"`
	demo       *DemoService    `gone:"*"`
	logger     gone.Logger     `gone:"*"`
}

func (s *AfterServerStart) Init() {
	s.afterStart(func() {
		s.logger.Infof("Starting MongoDB demo after server start...")
		if err := s.demo.RunDemo(); err != nil {
			s.logger.Errorf("Demo failed: %v", err)
		}
	})
}

func main() {
	gone.Serve()
}
