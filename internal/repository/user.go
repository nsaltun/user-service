package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nsaltun/userapi/internal/model"
	"github.com/nsaltun/userapi/pkg/lib/db/mongohandler"
	"github.com/nsaltun/userapi/pkg/lib/errwrap"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// userRepository implementor
type userRepository struct {
	collection *mongo.Collection
}

// NewUserRepository returns new instance to be able to use UserRepository interface methods.
//
// Creates index in this method
func NewUserRepository(db *mongohandler.MongoDBWrapper) (UserRepository, error) {
	repo := &userRepository{db.Collection("users")}
	err := repo.createIndexes()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// createIndexes creates indexes specific to the User collection
//
// Creating index for `email`(unique) and `nickName`(unique) and `country`.
func (r *userRepository) createIndexes() error {
	// Define index models
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}}, // Ascending index on email
			Options: options.Index().SetUnique(true),  // Unique constraint
		},
		{
			Keys:    bson.D{{Key: "country", Value: 1}}, // Ascending index on country
			Options: options.Index(),                    // Background creation
		},
		{
			Keys:    bson.D{{Key: "nickName", Value: 1}}, // Ascending index on nickName
			Options: options.Index().SetUnique(true),     // Unique constraint
		},
	}

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		slog.ErrorContext(ctx, "Error creating indexes for users collection", slog.Any("error", err))
		return err
	}

	slog.InfoContext(ctx, "Indexes created successfully for users collection.")
	return nil
}

// Create a new user
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	user.Id = uuid.NewString() // Generate a new UUID
	user.Status = model.UserStatus_Active
	user.Meta = model.NewMeta()
	_, err := r.collection.InsertOne(ctx, user)

	// empty password to not return in the api response
	// NOTE: empty password before logging error to not leak password in logs
	user.Password = ""

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			slog.InfoContext(ctx, "already exists with the same nickname or email.", slog.Any("error", err))
			return errwrap.ErrConflict.SetMessage("already exists with the same nickname or email")
		}
		slog.ErrorContext(ctx, "mongo create user error", slog.Any("error", err), slog.Any("user", user))
		return errwrap.ErrInternal.SetMessage("internal error").SetOriginError(err)
	}

	return nil
}

// Update user by id
//
// Excluding password from the response.
//
// Sanitizing fields for update.
//
// Checking uniqueness of user data.
//
// - Returns NotFound when record is not found
//
// - Returns Conflict when duplicated key error occured or uniqueness violated
//
// - Returns internal error for other error cases
//
// Returns updated user when it is successful with updated `UpdatedAt` and `Version` field
func (r *userRepository) Update(ctx context.Context, id string, user *model.User) (*model.User, error) {
	err := r.checkUniqueness(ctx, id, user)
	if err != nil {
		return nil, err
	}

	// findOneAndUpdate options
	opt := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetProjection(bson.M{
			"password": 0, //exclude password from the response
		})

	user.Meta.Update()

	// keeping some fields from updating.
	userM := sanitizeUserForUpdate(user)

	// filter by ID
	filter := bson.M{"_id": id} // Using string ID (UUID)

	// Use MongoDB's $set operator to update fields
	updatedUserM := r.collection.FindOneAndUpdate(ctx,
		filter,
		bson.M{"$set": userM, "$inc": bson.M{"version": 1}},
		opt)

	if err := updatedUserM.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errwrap.ErrNotFound.SetMessage("record not found")
		} else if mongo.IsDuplicateKeyError(err) {
			slog.InfoContext(ctx, "user update failed with duplicate key error", slog.Any("error", err), slog.Any("userBson", userM))
			return nil, errwrap.ErrConflict.SetMessage("unique constraint violated").SetOriginError(err)
		}

		slog.InfoContext(ctx, "user update failed.", slog.Any("error", err), slog.Any("userBson", userM))
		return nil, errwrap.ErrInternal.SetMessage("internal error").SetOriginError(err)
	}

	var updatedUser *model.User
	if err := updatedUserM.Decode(&updatedUser); err != nil {
		slog.InfoContext(ctx, "error while decoding bson user to user model", slog.Any("error", err), slog.Any("updatedUserBson", updatedUser))
		return nil, errwrap.ErrInternal.SetMessage("user decode error").SetOriginError(err)
	}

	return updatedUser, nil
}

// ListByFilter fetches users based on a dynamic filter with pagination and cursor
func (r *userRepository) ListByFilter(ctx context.Context, filter bson.M, limit int, offset int) ([]model.User, int64, error) {
	var users []model.User

	// Find the total count of documents that match the filter
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, "error from mongo while counting documents", slog.Any("error", err.Error()))
		return nil, 0, errwrap.ErrInternal.SetMessage("internal error").SetOriginError(err)
	}

	// Define MongoDB options for pagination
	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetProjection(bson.M{
			"password": 0,
		})

	// Query the database using the provided filter and options
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		slog.InfoContext(ctx, "error from mongo while finding docs by filter.", slog.Any("error", err))
		if err == mongo.ErrNoDocuments {
			return nil, 0, errwrap.ErrNotFound.SetMessage("user record not found")
		}
		return nil, 0, errwrap.ErrInternal.SetMessage("internal error").SetOriginError(err)
	}
	defer cursor.Close(ctx)

	// Decode users from the cursor
	if err := cursor.All(ctx, &users); err != nil {
		slog.InfoContext(ctx, "error from mongo cursor while getting docs by filter.", slog.Any("error", err))
		return nil, 0, errwrap.ErrInternal.SetMessage("internal error").SetOriginError(err)
	}

	return users, totalCount, nil
}

// Delete user by id
func (r *userRepository) Delete(ctx context.Context, id string) error {
	// findOneAndUpdate options
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	// keeping some fields from updating.
	userM := sanitizeUserForDelete()

	// filter by ID
	filter := bson.M{"_id": id} // Using string ID (UUID)

	// Use MongoDB's $set operator to update fields
	updatedUserM := r.collection.FindOneAndUpdate(ctx,
		filter,
		bson.M{"$set": userM},
		opt)

	if err := updatedUserM.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return errwrap.ErrNotFound.SetMessage("record not found")
		}
		slog.ErrorContext(ctx, "mongo error while deleting user", slog.Any("error", err), slog.Any("id", id))
		return errwrap.ErrInternal.SetMessage("internal error").SetOriginError(err)
	}

	return nil
}

// Get user by id
func (r *userRepository) Get(ctx context.Context, id string) (*model.User, error) {
	filter := bson.M{"_id": id}
	found := r.collection.FindOne(ctx, filter)
	if found.Err() != nil {
		return nil, found.Err()
	}

	var user *model.User
	if err := found.Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

// checkUniqueness checks uniqueness by filtering with unique constraint fields.
//
// Returns Conflict error if duplicated record exists.
func (r *userRepository) checkUniqueness(ctx context.Context, id string, user *model.User) error {
	// Pre-check for uniqueness
	filter := bson.M{
		"$or": []bson.M{
			{"nickName": user.NickName},
			{"email": user.Email},
		},
		"_id": bson.M{"$ne": id}, // Exclude the current document
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, "failed to check user uniqueness from mongo", slog.Any("error", err.Error()), slog.Any("userFilter", filter))
		return errwrap.ErrInternal.SetMessage("internal error").SetOriginError(err)
	}
	if count > 0 {
		return errwrap.ErrConflict.SetMessage("nickname or email should be unique").SetOriginError(err)
	}

	return nil
}

// sanitizeUserForUpdate excludes fields that should not be updated
func sanitizeUserForUpdate(user *model.User) bson.M {
	// Manually create the update map, allowing only specific fields
	return bson.M{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"nickName":  user.NickName,
		"email":     user.Email,
		"country":   user.Country,
		"updatedAt": user.UpdatedAt,
		"status":    user.Status,
	}
}

// sanitizeUserForDelete excludes fields that should not be updated while deleting a user
func sanitizeUserForDelete() bson.M {
	return bson.M{
		"status":    model.UserStatus_Inactive,
		"updatedAt": time.Now().UTC(),
	}
}
