package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Product Model
type Product struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Price float64            `json:"price,omitempty" bson:"price,omitempty"`
}

func (p *Product) getProduct(db *mongo.Database) error {
	err := db.Collection("products").FindOne(context.TODO(), bson.M{"_id": p.ID}).Decode(&p)
	return err
}

func (p *Product) createProduct(db *mongo.Database) (*mongo.InsertOneResult, error) {
	result, err := db.Collection("products").InsertOne(context.TODO(), p)
	return result, err
}

func (p *Product) updateProduct(db *mongo.Database) error {
	update := bson.M{
		"$set": bson.M{
			"name":  p.Name,
			"price": p.Price,
		},
	}

	_, err := db.Collection("products").UpdateOne(
		context.TODO(),
		bson.M{"_id": p.ID},
		update,
	)
	return err
}

func (p *Product) deleteProduct(db *mongo.Database) error {
	_, err := db.Collection("products").DeleteOne(context.TODO(), bson.M{"_id": p.ID})
	return err
}

func getProducts(db *mongo.Database) ([]Product, error) {
	cursor, err := db.Collection("products").Find(context.TODO(), bson.M{})
	var products []Product

	if err != nil {
		return []Product{}, err
	}

	err = cursor.All(context.TODO(), &products)
	return products, err
}
