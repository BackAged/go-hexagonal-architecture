package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Insert insets into databse
func (c *Client) Insert(ctx context.Context, col string, data interface{}) (interface{}, error) {
	r, err := c.client.Database(c.db).Collection(col).InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return r.InsertedID, nil
}

// FindByID finds something by primary key ID
func (c *Client) FindByID(ctx context.Context, col, id string) (*Row, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filtr := bson.M{"_id": objID}
	res := c.client.Database(c.db).Collection(col).FindOne(ctx, filtr)
	return &Row{s: res}, nil
}

// Find multiple things
func (c *Client) Find(ctx context.Context, col string, q bson.M) (*Rows, error) {
	opt := options.Find()
	cur, err := c.client.Database(c.db).Collection(col).Find(ctx, q, opt)
	if err != nil {
		return nil, err
	}
	return &Rows{c: cur}, nil
}
