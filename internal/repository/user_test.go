package repository

import (
	"context"
	"testing"

	"github.com/nsaltun/userapi/internal/model"
	"github.com/nsaltun/userapi/pkg/lib/db/mongohandler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMyRepositoryFunction(t *testing.T) {
	ctx := context.Background()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("MyRepositoryFunction", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(),
			mtest.CreateSuccessResponse(),
			// mtest.CreateSuccessResponse(bson.E{
			// 	"value", bson.M{"_id": "custom123", "key": 24},
			// }),
		)

		mongoWrapper := &mongohandler.MongoDBWrapper{
			Database: mt.DB,
		}

		repo, err := NewUserRepository(mongoWrapper)
		require.NoError(t, err)

		user := &model.User{FirstName: "test_firstname"}
		err = repo.Create(ctx, user)

		assert.NoError(t, err, "Should have successfully run")
	})
}
