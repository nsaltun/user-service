package repository

import (
	"context"
	"testing"

	"github.com/nsaltun/userapi/internal/model"
	mocks "github.com/nsaltun/userapi/pkg/mocks/lib/db/mongohandler"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	mongoMock := mocks.NewMongoDBWrapper(t)
	collectionMock := mocks.NewCollection(t)
	mongoMock.On("Collection", "users").Return(collectionMock).Once()
	collectionMock.On("CreateManyIndexes", mock.Anything, mock.Anything).Return(nil, nil)

	userRepo, err := NewUserRepository(mongoMock)
	require.NoError(t, err)

	collectionMock.On("InsertOne", mock.Anything, mock.Anything).Return(nil, nil).Once()

	ctx := context.TODO()
	user := &model.User{FirstName: "Enes", NickName: "testEnes", Email: "test@email.com", Password: "pass"}
	err = userRepo.Create(ctx, user)

	require.NoError(t, err)
}
