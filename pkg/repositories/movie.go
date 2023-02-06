package repositories

// napravi interjfejs i strukturka koja ce da implementira interjfejs

//apstraktuj db lajer

//movie repository interfejs i strutkura

//struktura sve implementira

import (
	"context"

	"github.com/urosradivojevic/health/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NetflixInterface interface {
	InsertOneMovie(movie model.Netflix) error
	GetAllMovies() ([]primitive.M, error)
	UpdateOneMovie(movieId string) error
	DeleteOneMovie(movieId string) error
}

type Netflix struct {
	col *mongo.Collection
}

func New(col *mongo.Collection) *Netflix {
	return &Netflix{
		col: col,
	}
}
func (n *Netflix) InsertOneMovie(movie model.Netflix) error {
	_, err := n.col.InsertOne(context.Background(), movie)
	if err != nil {
		return err
	}
	return nil
}

func (n *Netflix) GetAllMovies() ([]primitive.M, error) {
	cursor, err := n.col.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var movies []primitive.M

	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	defer cursor.Close(context.Background())
	return movies, nil
}

func (n *Netflix) UpdateOneMovie(movieId string) error {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	_, err := n.col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	return err
}

func (n *Netflix) DeleteOneMovie(movieId string) error {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}

	_, err := n.col.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return err
}
