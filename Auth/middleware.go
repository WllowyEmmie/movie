package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSONPretty(http.StatusUnauthorized, echo.Map{"error": "Invalid token"}, " ")
			}
			tokenstring := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}
				return jwt_secret_key, nil
			})
			if err != nil || !token.Valid {
				return c.JSONPretty(http.StatusUnauthorized, echo.Map{"error": err.Error()}, " ")
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return c.JSONPretty(http.StatusUnauthorized, echo.Map{"error": "Invalid token claims"}, " ")
			}
			userIDStr, ok := claims["sub"].(string)
			if !ok {
				return c.JSONPretty(http.StatusUnauthorized, echo.Map{"error": "Invalid subject claims"}, " ")
			}
			userID, err := primitive.ObjectIDFromHex(userIDStr)
			if err != nil {
				return c.JSONPretty(http.StatusUnauthorized, echo.Map{"error": "Invalid userID"}, " ")
			}
			c.Set("userID", userID)
			return next(c)
		}
	}

}
