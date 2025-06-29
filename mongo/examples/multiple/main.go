package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/gone-io/gone/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// UserEvent represents a user activity event
type UserEvent struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string                 `bson:"user_id" json:"user_id"`
	EventType string                 `bson:"event_type" json:"event_type"`
	Data      map[string]interface{} `bson:"data" json:"data"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
}

// AnalyticsData represents analytics information
type AnalyticsData struct {
	ID       primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Metric   string                 `bson:"metric" json:"metric"`
	Value    float64                `bson:"value" json:"value"`
	Date     time.Time              `bson:"date" json:"date"`
	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// MultiDBService demonstrates using multiple MongoDB connections
type MultiDBService struct {
	gone.Flag
	defaultDB   *mongoDriver.Client `gone:"*"`
	analyticsDB *mongoDriver.Client `gone:"*,mongo-analytics"`
	logsDB      *mongoDriver.Client `gone:"*,mongo-logs"`
	logger      gone.Logger         `gone:"*"`
}

// CreateUser creates a user in the main database
func (s *MultiDBService) CreateUser(name, email string) error {
	collection := s.defaultDB.Database("myapp").Collection("users")

	user := bson.M{
		"name":       name,
		"email":      email,
		"created_at": time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		s.logger.Errorf("Failed to create user: %v", err)
		return err
	}

	s.logger.Infof("User created in main DB: %s", name)
	return nil
}

// LogUserEvent logs a user event to the logs database
func (s *MultiDBService) LogUserEvent(userID, eventType string, data map[string]interface{}) error {
	collection := s.logsDB.Database("logs").Collection("user_events")

	event := &UserEvent{
		UserID:    userID,
		EventType: eventType,
		Data:      data,
		Timestamp: time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), event)
	if err != nil {
		s.logger.Errorf("Failed to log user event: %v", err)
		return err
	}

	s.logger.Infof("User event logged: %s - %s", userID, eventType)
	return nil
}

// RecordAnalytics records analytics data to the analytics database
func (s *MultiDBService) RecordAnalytics(metric string, value float64, metadata map[string]interface{}) error {
	collection := s.analyticsDB.Database("analytics").Collection("metrics")

	data := &AnalyticsData{
		Metric:   metric,
		Value:    value,
		Date:     time.Now(),
		Metadata: metadata,
	}

	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		s.logger.Errorf("Failed to record analytics: %v", err)
		return err
	}

	s.logger.Infof("Analytics recorded: %s = %f", metric, value)
	return nil
}

// GetUserEvents retrieves user events from logs database
func (s *MultiDBService) GetUserEvents(userID string, limit int64) ([]*UserEvent, error) {
	collection := s.logsDB.Database("logs").Collection("user_events")

	ctx := context.Background()
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID}, &options.FindOptions{
		Limit: &limit,
		Sort:  bson.M{"timestamp": -1}, // Sort by timestamp descending
	})
	if err != nil {
		s.logger.Errorf("Failed to get user events: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []*UserEvent
	for cursor.Next(ctx) {
		var event UserEvent
		if err := cursor.Decode(&event); err != nil {
			s.logger.Errorf("Failed to decode event: %v", err)
			continue
		}
		events = append(events, &event)
	}

	if err := cursor.Err(); err != nil {
		s.logger.Errorf("Cursor error: %v", err)
		return nil, err
	}

	s.logger.Infof("Retrieved %d events for user %s", len(events), userID)
	return events, nil
}

// GetAnalyticsSummary retrieves analytics summary from analytics database
func (s *MultiDBService) GetAnalyticsSummary(metric string, days int) (float64, error) {
	collection := s.analyticsDB.Database("analytics").Collection("metrics")

	// Calculate date range
	startDate := time.Now().AddDate(0, 0, -days)

	// Aggregation pipeline to sum values
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"metric": metric,
				"date": bson.M{
					"$gte": startDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$value"},
			},
		},
	}

	ctx := context.Background()
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		s.logger.Errorf("Failed to aggregate analytics: %v", err)
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Total float64 `bson:"total"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			s.logger.Errorf("Failed to decode analytics result: %v", err)
			return 0, err
		}
	}

	s.logger.Infof("Analytics summary for %s (last %d days): %f", metric, days, result.Total)
	return result.Total, nil
}

// DemoService demonstrates multi-database operations
type DemoService struct {
	gone.Flag
	multiDBService *MultiDBService `gone:"*"`
	logger         gone.Logger     `gone:"*"`
}

// RunDemo runs a demonstration of multiple database operations
func (d *DemoService) RunDemo() error {
	d.logger.Infof("Starting multiple MongoDB connections demo...")

	// Create a user in main database
	err := d.multiDBService.CreateUser("John Doe", "john@example.com")
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	userID := "user123"

	// Log some user events
	events := []struct {
		eventType string
		data      map[string]interface{}
	}{
		{"login", map[string]interface{}{"ip": "192.168.1.1", "device": "mobile"}},
		{"page_view", map[string]interface{}{"page": "/dashboard", "duration": 45}},
		{"purchase", map[string]interface{}{"product_id": "prod123", "amount": 29.99}},
		{"logout", map[string]interface{}{"session_duration": 1800}},
	}

	for _, event := range events {
		err := d.multiDBService.LogUserEvent(userID, event.eventType, event.data)
		if err != nil {
			d.logger.Errorf("Failed to log event %s: %v", event.eventType, err)
		}
		time.Sleep(100 * time.Millisecond) // Small delay between events
	}

	// Record some analytics
	analyticsData := []struct {
		metric   string
		value    float64
		metadata map[string]interface{}
	}{
		{"daily_active_users", 1250, map[string]interface{}{"source": "web"}},
		{"revenue", 29.99, map[string]interface{}{"currency": "USD", "product": "subscription"}},
		{"page_views", 5, map[string]interface{}{"user_id": userID}},
		{"session_duration", 1800, map[string]interface{}{"user_id": userID}},
	}

	for _, analytics := range analyticsData {
		err := d.multiDBService.RecordAnalytics(analytics.metric, analytics.value, analytics.metadata)
		if err != nil {
			d.logger.Errorf("Failed to record analytics %s: %v", analytics.metric, err)
		}
	}

	// Retrieve user events
	userEvents, err := d.multiDBService.GetUserEvents(userID, 10)
	if err != nil {
		d.logger.Errorf("Failed to get user events: %v", err)
	} else {
		d.logger.Infof("Retrieved %d user events", len(userEvents))
		for _, event := range userEvents {
			d.logger.Infof("Event: %s at %s", event.EventType, event.Timestamp.Format(time.RFC3339))
		}
	}

	// Get analytics summary
	revenueSummary, err := d.multiDBService.GetAnalyticsSummary("revenue", 7)
	if err != nil {
		d.logger.Errorf("Failed to get revenue summary: %v", err)
	} else {
		d.logger.Infof("Total revenue (last 7 days): $%.2f", revenueSummary)
	}

	pageViewsSummary, err := d.multiDBService.GetAnalyticsSummary("page_views", 7)
	if err != nil {
		d.logger.Errorf("Failed to get page views summary: %v", err)
	} else {
		d.logger.Infof("Total page views (last 7 days): %.0f", pageViewsSummary)
	}

	d.logger.Infof("Multiple MongoDB connections demo completed!")
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
		s.logger.Infof("Starting multiple MongoDB connections demo after server start...")
		if err := s.demo.RunDemo(); err != nil {
			s.logger.Errorf("Demo failed: %v", err)
		}
	})
}

func main() {
	gone.Serve()
}
