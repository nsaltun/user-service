package model

import "time"

type Meta struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	Version   int32     `bson:"version" json:"version"`
}

func NewMeta() Meta {
	now := time.Now().UTC()
	return Meta{
		CreatedAt: now,
		UpdatedAt: now,
		Version:   0,
	}
}

func (m *Meta) Update() {
	m.UpdatedAt = time.Now().UTC()
}
