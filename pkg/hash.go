package pkg

import "golang.org/x/crypto/bcrypt"

func HashPassword(plainPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
