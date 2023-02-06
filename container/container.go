package container

import (
	"context"
	"log"
	"os"

	"github.com/urosradivojevic/health/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Container struct {
	Client           *mongo.Client
	mongoDbDatabases map[string]*mongo.Database
	env              string
}

func New(env string) *Container {
	return &Container{
		mongoDbDatabases: make(map[string]*mongo.Database),
		env:              env,
	}
}
func (c *Container) GetEnviorment() string {
	return c.env
}
func (c *Container) GetMongoClient() *mongo.Client { //ovo pozivam preko kontejnera

	if c.Client != nil {
		return c.Client
	}
	uri := os.Getenv("MONGODB_URI")

	//client option
	clientOption := options.Client().ApplyURI(uri)

	//connect to mongoDB
	client, err := mongo.Connect(context.Background(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	c.Client = client
	return c.Client

}

func (c *Container) GetMongoDatabase() *mongo.Database {
	dbName := os.Getenv("MONGODB_DB")
	mongodbDatabes, exist := c.mongoDbDatabases[dbName]
	if exist {
		return mongodbDatabes
	}

	database := c.GetMongoClient().Database(dbName)
	c.mongoDbDatabases[dbName] = database
	return database
}

func (c *Container) GetMongoCollection(col string) *mongo.Collection {
	return c.GetMongoDatabase().Collection(col)
}

func (c *Container) GetNetflixRepository() repositories.NetflixInterface {
	return repositories.New(c.GetMongoCollection("watchlist"))
}
