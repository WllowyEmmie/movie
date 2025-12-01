package mongo_handler

import (
	"context"
	"errors"
	auth "movieapp/Auth"
	"movieapp/Hash"
	"movieapp/models"
	"movieapp/py"
	"movieapp/tmdb"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	DB *mongo.Database
}

func (h *Handler) Signup(c echo.Context) error {
	users := h.DB.Collection("users")

	var req models.SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSONPretty(400, echo.Map{"error": "Invalid Request"}, " ")
	}
	count, err := users.CountDocuments(
		c.Request().Context(),
		bson.M{"email": req.Email},
	)
	if err != nil {
		return c.JSONPretty(500, echo.Map{"error": "Database error"}, " ")
	}
	if count > 0 {
		return c.JSON(http.StatusConflict, echo.Map{"error": "user already exists"})
	}
	hash, err := Hash.HashPassword(req.Password)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "Failed to hash password"})
	}
	user := models.User{
		ID:              primitive.NewObjectID(),
		Email:           req.Email,
		PasswordHash:    hash,
		PreferredGenres: []int{},
		Favorites:       []models.Movie{},
		CreatedAt:       time.Now(),
	}
	if _, err := users.InsertOne(c.Request().Context(), user); err != nil {
		return c.JSONPretty(500, echo.Map{"error": "Failed to create user"}, " ")
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User created successfully",
		"user_id": user.ID.Hex(),
	})
}
func (h *Handler) Login(c echo.Context) error {
	users := h.DB.Collection("users")

	var req models.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSONPretty(400, echo.Map{"error": "Invalid request"}, " ")
	}
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email and password required"})
	}
	ctx := c.Request().Context()
	user, err := h.FindUserByEmail(ctx, users, req.Email)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User doesn't exists"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}
	if !Hash.CheckPassword(req.Password, user.PasswordHash) {
		return c.JSON(400, echo.Map{"error": "Invalid password"})
	}
	token, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Unable to generate token"})
	}
	return c.JSON(200, echo.Map{
		"message": "Login successful",
		"user_id": user.ID.Hex(),
		"token":   token,
	})
}
func (h *Handler) FindUserByEmail(ctx context.Context, users *mongo.Collection, email string) (models.User, error) {
	var user models.User
	err := users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (h *Handler) GetUser(c echo.Context) (models.User, error) {
	ctx := c.Request().Context()
	uid, ok := c.Get("userID").(primitive.ObjectID)
	if !ok {
		return models.User{}, c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	users := h.DB.Collection("users")
	var user models.User
	err := users.FindOne(ctx, bson.M{"_id": uid}).Decode(&user)
	if err != nil {
		return models.User{}, c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}
	return user, nil
}

func (h *Handler) AddToFavorites(c echo.Context) error {
	var req models.UpdateFavoriteRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(400, echo.Map{"error": "Invalid request"})
	}

	movie := req.Movie
	if movie.Name == "" {
		return c.JSON(http.StatusNotAcceptable, echo.Map{"error": "Movie name is cannot be empty"})
	}

	user, err := h.GetUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	update := bson.M{
		"$push": bson.M{
			"favorites": movie,
		},
	}
	_, err = h.DB.Collection("users").UpdateByID(
		c.Request().Context(),
		user.ID,
		update,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update favorites"})
	}
	updatedUser, _ := h.GetUser(c)
	return c.JSON(http.StatusCreated,
		echo.Map{
			"message":  "Favorites Updated",
			"movie":    movie,
			"favorite": updatedUser.Favorites,
		})
}
func (h *Handler) RecommendFromUserFavorites(c echo.Context) error {
	user, err := h.GetUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	if len(user.Favorites) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "No favorites found for user"})
	}
	recs, err := py.RecommendFromFavorites(user.Favorites)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK,
		echo.Map{
			"message":        "Reccomedation Genrated Successfully",
			"favorite":       user.Favorites,
			"reccomendation": recs,
		})
}
func (h *Handler) UpdatePreferredGenres(c echo.Context) error {
	mapMovieGenre := tmdb.MovieGenreMap

	var req models.UpdatePreferredGenresRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(400, echo.Map{"error": "Invalid request"})
	}
	genre := req.Genre

	user, err := h.GetUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	for key, value := range mapMovieGenre {
		if strings.EqualFold(value, genre) {
			update := bson.M{
				"$addToSet": bson.M{
					"preferred_genres": key,
				},
			}
			_, err = h.DB.Collection("users").UpdateByID(
				c.Request().Context(),
				user.ID,
				update,
			)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update favorites"})
			}
			user.PreferredGenres = append(user.PreferredGenres, key)
			return c.JSON(http.StatusCreated,
				echo.Map{
					"message":          "Preferred genres updated",
					"genre":            value,
					"preferred_genres": user.PreferredGenres,
				})
		}
	}
	return c.JSON(http.StatusBadRequest, echo.Map{
		"error":            "Invalid Genre",
		"available_genres": tmdb.MovieGenreMap,
	})
}
