package crypto

import "golang.org/x/crypto/bcrypt"

// Hash hashes a string using the bcrypt algorithm
func Hash(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}

// CompareHashAndData compares a hash to clear data and returns an error if the comparison fails.
func CompareHashAndData(hash string, data string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
}
