package model

import (
	"go.mongodb.org/mongo-driver/bson"
)

type UserStatus int

const (
	UserStatus_Active   UserStatus = 1 //Default
	UserStatus_Inactive UserStatus = 2
)

// User represents the user model
type User struct {
	Id        string           `bson:"_id,omitempty" json:"id"` // UUID as string
	FirstName string           `bson:"firstName" json:"firstName"`
	LastName  string           `bson:"lastName" json:"lastName"`
	NickName  string           `bson:"nickName" json:"nickName"`
	Password  string           `bson:"password" json:"password,omitempty"`
	Email     string           `bson:"email" json:"email"`
	Country   string           `bson:"country" json:"country"`
	Status    UserStatus       `bson:"status" json:"status"`
	Meta      `bson:",inline"` // Embed Meta fields directly
}

// UserFilter defines the criteria to filter users in MongoDB
type UserFilter struct {
	Id        string     `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	NickName  string     `json:"nickName"`
	Email     string     `json:"email"`
	Country   string     `json:"country"`
	Status    UserStatus `json:"status"`
}

// ParseUserFilter converts a UserFilter into a MongoDB filter
func (f *UserFilter) ToBson() bson.M {
	//TODO: Sorting might be added as well
	mongoFilter := bson.M{}

	if f.Id != "" {
		mongoFilter["_id"] = f.Id // Exact match for id
	}
	// Use exact matches for fields to utilize indexes
	if f.FirstName != "" {
		// Use prefix match instead of full regex if possible
		mongoFilter["firstName"] = bson.M{"$regex": "^" + f.FirstName, "$options": "i"} // Prefix matching
	}
	if f.LastName != "" {
		mongoFilter["lastName"] = bson.M{"$regex": "^" + f.LastName, "$options": "i"}
	}
	if f.NickName != "" {
		mongoFilter["nickName"] = f.NickName // Exact match
	}
	if f.Email != "" {
		mongoFilter["email"] = f.Email // Exact match for email
	}
	if f.Country != "" {
		mongoFilter["country"] = f.Country // Exact match for country code
	}
	if f.Status == 0 {
		mongoFilter["status"] = UserStatus_Active // Set active as default
	} else if f.Status > 0 {
		mongoFilter["status"] = f.Status //Set value coming from userFilter
	}

	return mongoFilter
}
