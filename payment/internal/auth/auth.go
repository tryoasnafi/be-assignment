package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
)

// // Supertokens middleware gin impl
// func Verify() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 			c.Next()
// 		})).ServeHTTP(c.Writer, c.Request)
// 		// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
// 		c.Abort()
// 	}
// }

// // This is a function that wraps the supertokens verification function
// // to work the gin
// func VerifySession(options *sessmodels.VerifySessionOptions) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		session.VerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
// 			c.Request = c.Request.WithContext(r.Context())
// 			c.Next()
// 		})(c.Writer, c.Request)
// 		// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
// 		c.Abort()
// 	}
// }

// func SessionInfo(c *gin.Context) {
// 	sessionContainer := session.GetSessionFromRequestContext(c.Request.Context())
// 	if sessionContainer == nil {
// 		c.JSON(500, "no session found")
// 		return
// 	}
// 	sessionData, err := sessionContainer.GetSessionDataInDatabase()
// 	if err != nil {
// 		err = supertokens.ErrorHandler(err, c.Request, c.Writer)
// 		if err != nil {
// 			c.JSON(500, err.Error())
// 			return
// 		}
// 		return
// 	}
// 	c.JSON(200, map[string]interface{}{
// 		"sessionHandle":      sessionContainer.GetHandle(),
// 		"userId":             sessionContainer.GetUserID(),
// 		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
// 		"sessionData":        sessionData,
// 	})
// }

// func GetUserIDFromRequest(r *http.Request) (uuid.UUID, error) {
// 	sessionContainer := session.GetSessionFromRequestContext(r.Context())
// 	if sessionContainer == nil {
// 		return uuid.UUID{}, fmt.Errorf("no session found")
// 	}
// 	userID := sessionContainer.GetUserID()
// 	return uuid.Parse(userID)
// }

var JWKs *jwk.Set

func FetchJWKs() {
	jwkEndpoint := os.Getenv("APP_JWK_ENDPOINT")
	log.Println("JWK Endpoint", jwkEndpoint)
	ctx := context.Background()
	// First, set up the `jwk.Cache` object. You need to pass it a
	// `context.Context` object to control the lifecycle of the background fetching goroutine.
	//
	// Note that by default refreshes only happen very 15 minutes at the
	// earliest. If you need to control this, use `jwk.WithRefreshWindow()`
	cache := jwk.NewCache(ctx)

	// Tell *jwk.Cache that we only want to refresh this JWKS
	// when it needs to (based on Cache-Control or Expires header from
	// the HTTP response). If the calculated minimum refresh interval is less
	// than 15 minutes, don't go refreshing any earlier than 15 minutes.
	cache.Register(jwkEndpoint, jwk.WithMinRefreshInterval(15*time.Minute))

	// Refresh the JWKS once before getting into the main loop.
	// This allows you to check if the JWKS is available before we start
	// a long-running program
	_, err := cache.Refresh(ctx, jwkEndpoint)
	if err != nil {
		log.Printf("failed to refresh Auth JWKS: %s\n", err)
		return
	}
	for {
		keyset, err := cache.Get(ctx, jwkEndpoint)
		if err != nil {
			fmt.Printf("failed to fetch google JWKS: %s\n", err)
			return
		}
		JWKs = &keyset
		time.Sleep(15 * time.Second)
	}
}


// Verify the JWT and return the payload
func VerifyJWT(token string) (map[string]any, error) {
	if JWKs == nil {
		return nil, fmt.Errorf("key set not fetched yet")
	}
	verified, err := jws.Verify([]byte(token), jws.WithKeySet(*JWKs))
	if err != nil {
		return nil, err
	}
	var payload map[string]any
	json.Unmarshal(verified, &payload)
	log.Println("JWT Payload", payload)
	return payload, nil
}

func ExtractJWT(c *gin.Context) (string, error) {
	header := c.Request.Header
	parts := strings.Split(header.Get("Authorization"), "Bearer ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid authorization header")
	}
	return parts[1], nil
}

// It's assume that the request already Verify
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	return uuid.Parse(c.GetString("userID"))
}

func Verify(c *gin.Context) {
	jwt, err := ExtractJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "message": err.Error() })
		c.Abort()
		return
	}
	payload, err := VerifyJWT(jwt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "token is not valid",
		})
		c.Abort()
		return
	}
	if sub, ok := payload["sub"]; ok {
		c.Set("userID", sub.(string))
	}
	c.Next()
}
