package middleware

import "net/http"

// AuthMiddleware checks for a dummy authentication token
func AuthMiddleware(next CustomHandler) CustomHandler {
	return func(ctx *HttpContext) error {
		token := ctx.Request.Header.Get("Authorization")
		if token != "Bearer valid-token" {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return nil
		}
		// Set UserID in HttpContext for further use
		ctx.UserID = "12345"
		next(ctx)
		return nil
	}
}
