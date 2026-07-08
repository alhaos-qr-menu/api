package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const userIDKey = "userID"

// RequireAuth returns a gin middleware that rejects requests without a
// valid "Authorization: Bearer <token>" header, and stores the
// authenticated user's ID in the gin context for downstream handlers.
func RequireAuth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		claims, err := ParseAccessToken(secret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(userIDKey, claims.UserID)
		c.Next()
	}
}

// UserIDFromContext retrieves the authenticated user's ID previously
// stored by RequireAuth. The second return value is false if no user
// ID is present in the context.
func UserIDFromContext(c *gin.Context) (int64, bool) {
	v, exists := c.Get(userIDKey)
	if !exists {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}
