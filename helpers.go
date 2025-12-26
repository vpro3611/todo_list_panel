package main

import (
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func Encrypter(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashAndPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}
	return false
}

func EncodeJSONhelper(w http.ResponseWriter, data any) error {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func ConvertToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
