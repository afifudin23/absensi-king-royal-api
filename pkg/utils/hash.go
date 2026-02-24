package utils

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	return hashedPassword, err
}

func CheckPassword(password string, hashedPassword string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	if err != nil {
		return false
	}
	return match
}
