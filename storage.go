package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

type Planet struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Climate     string             `bson:"climate"`
	Terrain     string             `bson:"terrain"`
	FilmCounter int                `bson:"film_counter"`
}

type Filter struct {
	Key   string
	Value string
}

type Storage struct {
	Planets *mongo.Collection
}

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "Planet not found"
}

func NewStorage(url string) (*Storage, error) {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}

	planets := client.Database("swapi").Collection("planets")
	return &Storage{planets}, nil
}

func (s Storage) Find(filters []Filter) ([]Planet, error) {
	query := make(bson.D, len(filters))

	for _, f := range filters {
		if f.Key == "id" {
			val, err := primitive.ObjectIDFromHex(f.Value)

			if err != nil {
				return nil, err
			}

			query = append(query, bson.E{Key: "_id", Value: val})
		} else {
			query = append(query, bson.E{Key: f.Key, Value: f.Value})
		}
	}

	planets := make([]Planet, 0)
	cur, err := s.Planets.Find(ctx, query)

	if err != nil {
		return nil, err
	}

	cur.All(ctx, &planets)

	return planets, nil
}

func (s Storage) Add(p Planet) error {
	p.ID = primitive.NewObjectID()
	_, err := s.Planets.InsertOne(ctx, p)

	if err != nil {
		return err
	}

	return nil
}

func (s Storage) Remove(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return &ErrNotFound{}
	}

	filter := bson.M{"_id": oid}

	res, err := s.Planets.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return &ErrNotFound{}
	}

	return nil
}
