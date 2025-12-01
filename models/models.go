package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Movie struct {
	ID             int      `bson:"tmdb_id,omitempty" json:"id"`
	UID            string   `bson:"uid,omitempty" json:"uid,omitempty"`
	Title          string   `bson:"title,omitempty" json:"title,omitempty"`
	Name           string   `bson:"name,omitempty" json:"name,omitempty"`
	Type           string   `bson:"type,omitempty" json:"type,omitempty"`
	Overview       string   `bson:"overview,omitempty" json:"overview,omitempty"`
	Release_Date   string   `bson:"release_date,omitempty" json:"release_date,omitempty"`
	First_Air_Date string   `bson:"air_date,omitempty" json:"first_air_date,omitempty"`
	VoteAverage    float64  `bson:"vote_average,omitempty" json:"vote_average,omitempty"`
	Genre_Ids      []int    `bson:"genre_ids,omitempty" json:"genre_ids,omitempty"`
	Genre          []string `bson:"genres,omitempty" json:"genres,omitempty"`
	Poster_Path    string   `bson:"poster_path,omitempty" json:"poster_path,omitempty"`
}

type PyMovie struct {
	ShowID      *string `json:"show_id"`
	Title       string  `json:"title"`
	Type        *string `json:"type"`
	Director    *string `json:"director"`
	Cast        *string `json:"cast"`
	Country     *string `json:"country"`
	DateAdded   *string `json:"date_added"`
	ReleaseYear *int    `json:"release_year"`
	Rating      *string `json:"rating"`
	Duration    *string `json:"duration"`
	ListedIn    *string `json:"listed_in"`
	Description *string `json:"description"`
}

type Actor struct {
	ID         int     `bson:"tmdb_id,omitempty" json:"id"`
	UID        string  `bson:"uid,omitempty" json:"uid,omitempty"`
	Title      string  `bson:"title,omitempty" json:"title,omitempty"`
	Name       string  `bson:"name,omitempty" json:"name,omitempty"`
	Popularity float64 `bson:"popularity" json:"popularity"`
	Known_for  []Movie `bson:"known_for" json:"known_for"`
}
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username        string             `bson:"username" json:"username"`
	Email           string             `bson:"email" json:"email"`
	PasswordHash    string             `bson:"password_hash" json:"-"`
	PreferredGenres []int              `bson:"preferred_genres" json:"preferred_genres"`
	Favorites       []Movie            `bson:"favorites" json:"favorites"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
}
type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type UpdateFavoriteRequest struct {
	Movie Movie `json:"movie"`
}
type UpdatePreferredGenresRequest struct {
	Genre string `json:"genre"`
}
type Handler struct {
	DB *mongo.Database
}
