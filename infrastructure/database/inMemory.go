package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
)

// Rows mongo rows
type Rows struct {
	si     *redis.ScanIterator
	client *redis.Client
}

// Next gives redis iterator cursor next
func (r *Rows) Next() bool {
	return r.si.Next()
}

// Err redis iterator Err
func (r *Rows) Err() error {
	return r.si.Err()
}

// Val returns the val
func (r *Rows) Val() string {
	return r.si.Val()
}

// Scan scans a value
func (r *Rows) Scan(v interface{}) error {
	key := r.si.Val()
	data, err := r.client.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), v)
}

// InMemoryClient holds redis client
type InMemoryClient struct {
	client *redis.Client
	db     string
}

// NewInMemoryClient returns a new Client
func NewInMemoryClient(host string, password string, db *int) (*InMemoryClient, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     host,     // use default Addr
		Password: password, // no password set
		DB:       *db,      // use default DB
	})

	_, err := c.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &InMemoryClient{client: c}, nil
}

// Insert insets into databse
func (c *InMemoryClient) Insert(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	err := c.client.Set(key, data, expiration).Err()
	return err
}

// FindOne finds something by primary key ID
func (c *InMemoryClient) FindOne(ctx context.Context, key string) (interface{}, error) {
	data, err := c.client.Get(key).Result()
	return data, err
}

// Find finds something
func (c *InMemoryClient) Find(ctx context.Context, key string) *Rows {
	si := c.client.Scan(0, key, 0).Iterator()
	return &Rows{si: si, client: c.client}
}

// Update updates database
func (c *InMemoryClient) Update(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	err := c.client.SetXX(key, data, 0).Err()
	return err
}
