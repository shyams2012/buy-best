package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// SecretKey secret key being used to sign tokens
var (
	SecretKey = []byte("some_deep_secret#129086543")
)

// type User struct {
// 	ID             string        `json:"id"`
// 	CreatedAt      time.Time `json:"createdAt"`
// 	ModifiedAt     time.Time `json:"modifiedAt"`
// 	Username       string        `json:"username"`
// 	Fullname       string        `json:"fullname"`
// 	Role           UserRole      `json:"role"`
// 	PasswordHash   string        `json:"-"`
// 	RefreshCounter int           `json:"-"`
// 	RefreshedTill  int           `json:"-"`
// 	IsActive       bool          `json:"isActive"`
// }

type User struct {
	ID             string      `json:"id"`
	Username       string      `json:"username"`
	Fullname       string      `json:"fullname"`
	CreatedAt      time.Time   `json:"createdAt"`
	ModifiedAt     time.Time   `json:"modifiedAt"`
	StreetNo       *int        `json:"streetNo"`
	ZipCode        *int        `json:"zipCode"`
	City           *string     `json:"city"`
	Mobile         *string     `json:"mobile"`
	Email          string      `json:"email" gorm:"unique"`
	Role           UserRole    `json:"role"`
	PasswordHash   string      `json:"-"`
	RefreshCounter int         `json:"-"`
	RefreshedTill  int         `json:"-"`
	IsActive       bool        `json:"isActive"`
	PaymentMethod  PaymentMode `json:"paymentMethod"`
}

func (u *User) SetPassword(password string) {
	u.PasswordHash = getHash(password)
}

func (u *User) CheckPassword(password string) bool {
	return getHash(password) == u.PasswordHash
}

func (u *User) Token() (*AuthToken, error) {
	authToken, err := u._AuthToken()
	if err != nil {
		return nil, err
	}
	refreshToken, err := u._RefreshToken()
	if err != nil {
		return nil, err
	}

	return &AuthToken{
		Token:        authToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *User) _AuthToken() (string, error) {
	/* Set token claims */
	userBytes, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = string(userBytes)
	claims["exp"] = time.Now().Add(time.Second * 3600 * 24 * 30).Unix()

	authToken, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return authToken, nil
}

func (u *User) _RefreshToken() (string, error) {
	/* Set token claims */
	userBytes, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	counter := u.RefreshCounter + 1
	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = string(userBytes)
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims["counter"] = counter

	refreshToken, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	u.RefreshCounter = counter

	return refreshToken, nil
}

func (u *User) Refresh(refreshCounter int) (*AuthToken, error) {
	if refreshCounter <= u.RefreshedTill {
		return nil, errors.New("old refresh token")
	}

	u.RefreshedTill = refreshCounter

	token, err := u.Token()
	if err != nil {
		return nil, err
	}

	return token, nil
}

func getHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	if tx := db.Where("username = ?", username).First(&user); tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func ParseAuthToken(tokenStr string) (*User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	var user User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiresAt := claims["exp"].(float64)
		if float64(time.Now().Unix()) > expiresAt {
			return nil, errors.New("token expired")
		}

		userStr := claims["user"].(string)
		err = json.Unmarshal([]byte(userStr), &user)
		if err != nil {
			return nil, err
		}

		return &user, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

func ParseRefreshToken(tokenStr string) (*string, *int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	var user User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiresAt := claims["exp"].(float64)
		if float64(time.Now().Unix()) > expiresAt {
			return nil, nil, errors.New("refresh token expired")
		}

		userStr := claims["user"].(string)
		err = json.Unmarshal([]byte(userStr), &user)
		if err != nil {
			return nil, nil, err
		}

		counterF, ok := claims["counter"].(float64)
		if !ok {
			return nil, nil, errors.New("invalid token")
		}
		counter := int(counterF)

		return &user.ID, &counter, nil
	} else {
		return nil, nil, errors.New("invalid token")
	}
}
