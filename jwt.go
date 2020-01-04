package gostp

import (
	"fmt"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

// JWT token struct
type JWT struct {
	*jwt.Token
}

type JWTMiddleware struct {
	Options jwtmiddleware.Options
}

// JwtMiddleware - middleware which validates token
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(Settings.SigningKey), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler:  CustomJWTError,
})

// JWTHandler gets http request, checks jwt token (if it's correct and not expired)
func JWTHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Let secure process the request. If it returns an error,
		// that indicates the request should not continue.
		err := JwtMiddleware.CheckJWT(w, r)

		// If there was an error, do not continue.
		if err != nil {
			fmt.Println("wow, error")
			return
		}

		rawToken, err := JwtMiddleware.Options.Extractor(r)
		token, err := JWTParse(rawToken)
		if token.IsExpired() {
			Header(w)
			http.Error(w, `{"error":"token expired"}`, 401)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// JWTParse - parses jwt token
func JWTParse(t string) (JWT, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(Settings.SigningKey), nil
	})

	if err != nil {
		fmt.Printf("token parse problem - %v", err)
		return JWT{}, fmt.Errorf("token parse problem: %v", err)
	}
	return JWT{token}, nil
}

// CustomJWTError - returns error if validation fails
func CustomJWTError(w http.ResponseWriter, r *http.Request, err string) {
	http.Error(w, err, http.StatusUnauthorized)
}

// IsExpired - checks if token is expired
func (token *JWT) IsExpired() bool {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("claim problem")
		return true
	}
	expiresIn := claims["expires_in"].(float64)
	nowUnixTime := time.Now().Unix()
	if int64(expiresIn)-nowUnixTime < 0 {
		fmt.Println("token expired")
		return true
	}
	return false
}

// GetUserID - gets user id from jwt token
func (token *JWT) GetUserID() (float64, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("claims problem")
		return 0, fmt.Errorf("claims problem")
	}
	return claims["user_id"].(float64), nil
}

// GenerateToken - generates new token
func GenerateToken(userID uint, expiresIn int64) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	claims["user_id"] = userID
	claims["expires_in"] = expiresIn
	fmt.Println("generated expiresIn:", expiresIn)
	token.Claims = claims

	tokenString, _ := token.SignedString([]byte(Settings.SigningKey))

	return tokenString
}

// RefreshUserTokens - refreshes user tokens
func RefreshUserTokens(user User, userTokens *UserTokens) {
	accessExpiresIn := time.Now().Add(time.Minute * time.Duration(Settings.JWTaccessExpiration)).Unix()
	refreshExpiresIn := time.Now().Add(time.Minute * time.Duration(Settings.JWTrefreshExpiration)).Unix()
	accessToken := GenerateToken(user.ID, accessExpiresIn)
	refreshToken := GenerateToken(user.ID, refreshExpiresIn)

	user.Token.RefreshToken = refreshToken
	Db.Save(&user)

	userTokens.AccessToken = accessToken
	userTokens.RefreshToken = refreshToken
	userTokens.AccessExpiresIn = accessExpiresIn
}

// GetUserIDClaim returns UserId
func GetUserIDClaim(tokenString string) (float64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(Settings.SigningKey), nil
	})

	if err != nil {
		fmt.Printf("token parce problem - %v", err)
		return 0, fmt.Errorf("token parce problem: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("claim error: %v", ok)
	}
	return claims["user_id"].(float64), nil
}
