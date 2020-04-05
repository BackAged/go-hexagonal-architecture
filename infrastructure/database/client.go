package database

import (
	"context"
	"io"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client mongo db client
type Client struct {
	client *mongo.Client
	db     string
}

// NewClient returns a new mongo db client
func NewClient(uri, db string) (*Client, error) {
	opts := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second)
	c, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	return &Client{client: c, db: db}, nil
}

// Row is mongo row
type Row struct {
	s   *mongo.SingleResult
	mu  sync.RWMutex
	raw bson.Raw
	err error
}

// Close closes mongo row cursor
func (r *Row) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.s = nil
	r.err = nil
	r.raw = nil
	return nil
}

// Next gives mongo row cursor next value
func (r *Row) Next() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.s == nil {
		return false
	}
	r.raw, r.err = r.s.DecodeBytes()
	r.s = nil

	return r.err == nil && r.raw != nil
}

// Scan scans bson to v
func (r *Row) Scan(v interface{}) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.raw == nil {
		return io.EOF
	}
	return bson.Unmarshal(r.raw, v)
}

// Err mongo row error
func (r *Row) Err() error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.err
}

// Rows mongo rows
type Rows struct {
	c *mongo.Cursor
}

// Next gives mongo row cursor next value
func (r *Rows) Next() bool {
	return r.c.Next(context.TODO())
}

// Close closes mongo row cursor
func (r *Rows) Close() error {
	return r.c.Close(context.TODO())
}

// Err mongo rows cursor Err
func (r *Rows) Err() error {
	return r.c.Err()
}

// Scan bson to v
func (r *Rows) Scan(v interface{}) error {
	return r.c.Decode(v)
}
