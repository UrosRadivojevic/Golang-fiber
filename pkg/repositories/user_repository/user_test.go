package user_repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/repositories/user_repository"
	"go.mongodb.org/mongo-driver/bson"
)

func TestRegisterSuccess(t *testing.T) {
	// arrange
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	assert := require.New(t)
	c := container.New("testing")
	col := c.GetMongoCollection("users")
	userRepo := user_repository.New(col)
	user := model.User{
		Username:  "urke",
		Password:  "urke1234",
		Firstname: "Uros",
	}
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})

	// act
	id, err := userRepo.Register(context.Background(), user)
	userCheck, _ := userRepo.GetByUsername(context.Background(), user.Username)
	// assert
	assert.NoError(err)
	assert.Equal(user.Firstname, userCheck.Firstname)
	assert.Equal(user.Username, userCheck.Username)
	assert.Equal(id.Hex(), userCheck.ID.Hex())
}

func TestGetByUsernameSuccess(t *testing.T) {
	// arrange
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	assert := require.New(t)
	c := container.New("testing")
	col := c.GetMongoCollection("users")
	userRepo := user_repository.New(col)
	u := model.User{
		Username:  "urke",
		Password:  "urke1234",
		Firstname: "Uros",
	}
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	id, _ := userRepo.Register(context.Background(), u)

	// act
	user, err := userRepo.GetByUsername(context.Background(), u.Username)

	// assert
	assert.NoError(err)
	assert.Equal(u.Firstname, user.Firstname)
	assert.Equal(u.Username, user.Username)
	assert.Equal(id.Hex(), user.ID.Hex())
}

func TestExistsSuccess(t *testing.T) {
	// arrange
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	assert := require.New(t)
	c := container.New("testing")
	col := c.GetMongoCollection("users")
	userRepo := user_repository.New(col)
	u := model.User{
		Username:  "urke",
		Password:  "urke1234",
		Firstname: "Uros",
	}
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	_, _ = userRepo.Register(context.Background(), u)

	// act
	exist := userRepo.Exists(context.Background(), u.Username)

	// assert
	assert.True(exist)
}

func TestExistsFail(t *testing.T) {
	// arrange
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	assert := require.New(t)
	c := container.New("testing")
	col := c.GetMongoCollection("users")
	userRepo := user_repository.New(col)
	u := model.User{
		Username:  "urke",
		Password:  "urke1234",
		Firstname: "Uros",
	}
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	_, _ = userRepo.Register(context.Background(), u)

	// act
	exist := userRepo.Exists(context.Background(), "test")

	// assert
	assert.False(exist)
}
