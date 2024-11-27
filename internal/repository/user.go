package repository

import (
	"context"
	"fmt"
	"log"
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

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db mongohandler.MongoDBWrapper) (UserRepository, error) {
	repo := &userRepository{db.Collection("users")}
	err := repo.createIndexes()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// createIndexes creates indexes specific to the User collection
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
		log.Printf("Error creating indexes for users collection: %v", err)
		return err
	}

	log.Println("Indexes created successfully for users collection.")
	return nil
}

// Create a new user
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	user.Id = uuid.NewString() // Generate a new UUID
	user.Status = model.UserStatus_Active
	user.Meta = model.NewMeta()
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil && mongo.IsDuplicateKeyError(err) {
		slog.Info("already exists with the same nickname or email.", slog.Any("error", err))
		return errwrap.ErrConflict.SetMessage("already exists with the same nickname or email")
	}

	// empty password to not return in the api response
	user.Password = ""
	return err
}

// Update user by id
func (r *userRepository) Update(ctx context.Context, id string, user *model.User) (*model.User, error) {
	err := r.checkUniqueness(ctx, id, user)
	if err != nil {
		return nil, err
	}

	// findOneAndUpdate options
	opt := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetProjection(bson.M{
			"password": 0,
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

	if updatedUserM.Err() != nil {
		return nil, updatedUserM.Err()
	}

	var updatedUser *model.User
	if err := updatedUserM.Decode(&updatedUser); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// ListByFilter fetches users based on a dynamic filter
func (r *userRepository) ListByFilter(ctx context.Context, filter bson.M, limit int, offset int) ([]model.User, int64, error) {
	var users []model.User

	// Find the total count of documents that match the filter
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
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
		if err == mongo.ErrNoDocuments {
			return nil, 0, errwrap.ErrNotFound.SetMessage("user record not found")
		}
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode users from the cursor
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
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

	if updatedUserM.Err() != nil {
		return updatedUserM.Err()
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
		return fmt.Errorf("failed to check uniqueness: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("unique constraint violated")
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
		"password":  user.Password,
		"email":     user.Email,
		"country":   user.Country,
		"updatedAt": user.UpdatedAt,
	}
}

// sanitizeUserForDelete excludes fields that should not be updated while deleting a user
func sanitizeUserForDelete() bson.M {
	return bson.M{
		"status":    model.UserStatus_Inactive,
		"updatedAt": time.Now().UTC(),
	}
}
