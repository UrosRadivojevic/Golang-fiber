package repositories

// napravi interjfejs i strukturka koja ce da implementira interjfejs

//apstraktuj db lajer

//movie repository interfejs i strutkura

//struktura sve implementira

import (
	"context"

	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NetflixInterface interface {
	InsertOneMovie(movie requests.CreateMovieRequest) (primitive.ObjectID, error)
	GetAllMovies() ([]model.Netflix, error)
	UpdateOneMovie(movieId string) error
	DeleteOneMovie(movieId string) error
	GetOneMovie(movieId string) (model.Netflix, error)
}

type Netflix struct {
	col *mongo.Collection
}

func New(col *mongo.Collection) *Netflix {
	return &Netflix{
		col: col,
	}
}
func (n *Netflix) InsertOneMovie(movie requests.CreateMovieRequest) (primitive.ObjectID, error) {
	m, err := n.col.InsertOne(context.Background(), movie)
	if err != nil {
		return primitive.NilObjectID, err
	}
	id := m.InsertedID.(primitive.ObjectID)
	return id, nil
}

func (n *Netflix) GetAllMovies() ([]model.Netflix, error) {
	cursor, err := n.col.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var movies []model.Netflix

	for cursor.Next(context.Background()) {
		var movie model.Netflix
		err := cursor.Decode(&movie)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	defer cursor.Close(context.Background())
	return movies, nil
}
func (n *Netflix) GetOneMovie(movieId string) (model.Netflix, error) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	val := n.col.FindOne(context.Background(), filter)
	var movie model.Netflix
	if err := val.Decode(&movie); err != nil {
		return model.Netflix{}, err
	}
	return movie, nil
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
