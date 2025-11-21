package models

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
type Actor struct {
	ID         int     `bson:"tmdb_id,omitempty" json:"id"`
	UID        string  `bson:"uid,omitempty" json:"uid,omitempty"`
	Title      string  `bson:"title,omitempty" json:"title,omitempty"`
	Name       string  `bson:"name,omitempty" json:"name,omitempty"`
	Popularity float64 `bson:"popularity" json:"popularity"`
	Known_for  []Movie `bson:"known_for" json:"known_for"`
}
type User struct {
	ID           string  `bson:"user_id" json:"id"`
	Username     string  `bson:"username" json:"username"`
	Email        string  `bson:"email" json:"email"`
	PasswordHash string  `bson:"password_hash" json:"-"` // never send password to client
	Favorites    []Movie `bson:"favorites" json:"favorites"`
}
