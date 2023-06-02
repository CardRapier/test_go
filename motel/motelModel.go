package motel

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

type Room struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type CreateRoom struct {
	Name string `bson:"name" json:"name"`
}

type Motel struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Telephone   string             `json:"telephone"`
	Location    Location           `json:"location"`
	// Rooms       []Room   `json:"rooms"`
}

type CreateMotel struct {
	// ID          primitive.ObjectID `bson:"_id, omitempty" json:"id, omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Telephone   string   `json:"telephone"`
	Location    Location `json:"location"`
	// Rooms       []CreateRoom `json:"rooms" bson:"rooms, omit"`
}
