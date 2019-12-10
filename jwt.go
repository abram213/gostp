package system

import (
	"fmt"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	*jwt.Token
}

var MySigningKey = []byte("u495CqrE*ZY!zR%8Wv7oIvvjUg")

var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return MySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler:  CustomJWTError,
})

func JWTParse(t string) (JWT, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return MySigningKey, nil
	})

	if err != nil {
		fmt.Printf("token parse problem - %v", err)
		return JWT{}, fmt.Errorf("token parse problem: %v", err)
	}
	return JWT{token}, nil
}

func CustomJWTError(w http.ResponseWriter, r *http.Request, err string) {
	CommonHeader(w)
	http.Error(w, err, http.StatusUnauthorized)
}

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

func (token *JWT) GetUserID() (float64, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("claims problem")
		return 0, fmt.Errorf("claims problem")
	}
	return claims["user_id"].(float64), nil
}

func GenerateToken(userId uint, expiresIn int64) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	claims["user_id"] = userId
	claims["expires_in"] = expiresIn
	token.Claims = claims

	tokenString, _ := token.SignedString(MySigningKey)

	return tokenString
}

func RefreshUserTokens(user User, userTokens *UserTokens) {
	expiresIn := time.Now().Add(time.Minute * 30).Unix()
	accessToken := GenerateToken(user.ID, expiresIn)
	refreshToken := GenerateToken(user.ID, time.Now().Add(time.Hour*24*30).Unix())

	user.Token.RefreshToken = refreshToken
	Db.Save(&user)

	userTokens.AccessToken = accessToken
	userTokens.RefreshToken = refreshToken
	userTokens.AccessExpiresIn = expiresIn
}

// GetUserIDClaim returns UserId
func GetUserIDClaim(tokenString string) (float64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return MySigningKey, nil
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
