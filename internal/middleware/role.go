package middleware

import (
	"MicroShopik/internal/domain"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			user := c.Get("user")
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found in context"})
			}

			claims, ok := user.(*domain.JWTClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user claims"})
			}

			for _, userRole := range claims.Roles {
				if userRole == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error":         "insufficient permissions",
				"required_role": role,
			})
		}
	}
}

func RequireAnyRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user")
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found in context"})
			}

			claims, ok := user.(*domain.JWTClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user claims"})
			}
			for _, requiredRole := range roles {
				for _, userRole := range claims.Roles {
					if userRole == requiredRole {
						return next(c)
					}
				}
			}
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"error":          "insufficient permissions",
				"required_roles": roles,
			})
		}
	}
}
