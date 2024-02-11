package password

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// Hash hashes the given password using bcrypt
func Hash(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// using recommended cost parameters from - https://godoc.org/golang.org/x/crypto/scrypt
	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// return hex-encoded string with salt appended to password
	hashedPassword := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPassword, nil
}

// Compare compares the stored password with the supplied one
func Compare(storedPassword string, suppliedPassword string) (bool, error) {
	pwsalt := strings.Split(storedPassword, ".")

	if len(pwsalt) < 2 {
		return false, fmt.Errorf("did not provide a valid hash")
	}

	// check supplied password salted with hash
	salt, err := hex.DecodeString(pwsalt[1])

	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)

	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	return hex.EncodeToString(shash) == pwsalt[0], nil
}
