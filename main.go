package main

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/scrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Password struct {
	gorm.Model
	Username string
	Salt     string
	Hash     string
}

func hashPassword(hashFunc func(password, salt []byte, N, r, p, keyLen int) ([]byte, error), dbFunc func(username string, salt string, hash string) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request body
		type passwordRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var req passwordRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "invalid request")
		}

		// Generate a random salt
		salt := make([]byte, 16)
		if _, err := rand.Read(salt); err != nil {
			return c.String(http.StatusInternalServerError, "failed to generate salt")
		}
		saltHex := hexutil.Encode(salt)

		// Hash the password using scrypt
		hash, err := hashFunc([]byte(req.Password), salt, 32768, 8, 1, 32)
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to hash password")
		}
		hashHex := hex.EncodeToString(hash)

		// Store the salt and hash in the database
		if err := dbFunc(req.Username, saltHex, hashHex); err != nil {
			return c.String(http.StatusInternalServerError, "failed to save password")
		}

		// Return HTTP status OK
		return c.String(http.StatusOK, "")
	}
}

func saveToDB(username string, salt string, hash string) error {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres dbname=postgres password=postgres sslmode=disable"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer func() {
		dbSql, err := db.DB()
		if err != nil {
			fmt.Println(err)
			return
		}
		dbSql.Close()
	}()

	p := Password{
		Username: username,
		Salt:     salt,
		Hash:     hash,
	}
	if err := db.Create(&p).Error; err != nil {
		return err
	}

	return nil
}

func verifyPassword(hashFunc func(password, salt []byte, N, r, p, keyLen int) ([]byte, error), password string, saltHex string, hashHex string) bool {
	salt, err := hexutil.Decode(saltHex)
	if err != nil {
		return false
	}
	hash, err := hex.DecodeString(hashHex)
	if err != nil {
		return false
	}

	computedHash, err := hashFunc([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return false
	}

	return hmac.Equal(hash, computedHash)
}

func getPassword(db *gorm.DB, hashFunc func(password, salt []byte, N, r, p, keyLen int) ([]byte, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request body
		type passwordRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var req passwordRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "invalid request")
		}
		// Get the password hash from the database
		var p Password
		if err := db.Where("username = ?", req.Username).First(&p).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.String(http.StatusNotFound, "username not found")
			}
			return c.String(http.StatusInternalServerError, "failed to get password")
		}

		// Verify the password hash
		if !verifyPassword(hashFunc, req.Password, p.Salt, p.Hash) {
			return c.String(http.StatusUnauthorized, "invalid password")
		}

		// Return HTTP status OK
		return c.String(http.StatusOK, "")
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/password", hashPassword(scrypt.Key, saveToDB))
	e.GET("/password", getPasswordHandler(getPasswordFromDB))

	e.Logger.Fatal(e.Start(":8080"))
}

func getPasswordHandler(getPassword func(string) (*Password, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.QueryParam("username")
		password := c.QueryParam("password")

		p, err := getPassword(username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.String(http.StatusUnauthorized, "invalid username or password")
			}
			return c.String(http.StatusInternalServerError, "failed to query password")
		}

		salt, err := base64.StdEncoding.DecodeString(p.Salt)
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to decode salt")
		}

		hash, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to hash password")
		}

		hashStr := base64.StdEncoding.EncodeToString(hash)

		if hashStr != p.Hash {
			return c.String(http.StatusUnauthorized, "invalid username or password")
		}

		return c.String(http.StatusOK, "login successful")
	}
}

func getPasswordFromDB(username string) (*Password, error) {
	dsn := "host=localhost port=5432 user=postgres dbname=passwords password=password sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %v", err)
	}
	defer sqlDB.Close()

	var p Password
	if err := db.Where("username = ?", username).First(&p).Error; err != nil {
		return nil, err
	}

	return &p, nil
}
