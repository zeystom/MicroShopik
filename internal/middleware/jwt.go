package middleware

import (
	"MicroShopik/internal/domain"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized,
					map[string]string{"error": "missing token"})
			}

			if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				return c.JSON(http.StatusUnauthorized,
					map[string]string{"error": "invalid authorization header format"})
			}

			tokenString := authHeader[7:]

			token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{},
				func(t *jwt.Token) (interface{}, error) {
					return []byte(secret), nil
				})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized,
					map[string]string{"error": "invalid token"})
			}

			claims := token.Claims.(*domain.JWTClaims)
			c.Set("user_id", claims.UserID)
			c.Set("user", claims)
			return next(c)
		}
	}
}
