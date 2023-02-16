package crypto

import (
	"crypto/aes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

func GetHashForByteArray(value []byte) string {
	hasher := sha1.New()
	hasher.Write(value)
	sha1_hash := hex.EncodeToString(hasher.Sum(nil))
	return sha1_hash
}

func GetHashForString(value string) string {
	hasher := sha1.New()
	hasher.Write([]byte(value))
	sha1_hash := hex.EncodeToString(hasher.Sum(nil))
	return sha1_hash
}

func EncryptAES(plaintext string, key string) (string, error) {
	// create cipher
	bkey := []byte(key)
	c, err := aes.NewCipher(bkey)
	if nil != err {
		return "", err
	}

	newplain := base64.StdEncoding.EncodeToString([]byte(plaintext))

	// allocate space for ciphered data
	out := make([]byte, len(newplain))

	// encrypt
	c.Encrypt(out, []byte(newplain))

	// return hex string
	return hex.EncodeToString(out), nil
}

func DecryptAES(encryptedValue string, key string) (string, error) {
	bkey := []byte(key)
	ciphertext, _ := hex.DecodeString(encryptedValue)

	c, err := aes.NewCipher(bkey)
	if nil != err {
		return "", err
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	return s, nil
}
