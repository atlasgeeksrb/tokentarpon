package tokencrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
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

// Encryption functions credit
// https://gist.github.com/17twenty/b7a050d6a3ed991db0433d4a1fc50de7
func EncryptAES(plaintextValue string, keyValue string) (string, error) {

	// Byte array of the string
	plaintext := []byte(plaintextValue)

	// Key
	key := []byte(keyValue)

	keyError := checkKey(key)
	if keyError != nil {
		return "", keyError
	}

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return string(ciphertext), nil
}

func DecryptAES(encryptedValue string, keyValue string) (string, error) {

	// Byte array of the string
	ciphertext := []byte(encryptedValue)

	// Key
	key := []byte(keyValue)

	keyError := checkKey(key)
	if keyError != nil {
		return "", keyError
	}

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		tooshortErr := errors.New("text is too short")
		return "", tooshortErr
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func checkKey(key []byte) error {
	// check key size
	keylen := len(key)
	if keylen != 16 && keylen != 32 && keylen != 64 {
		badKeySize := errors.New("key size must be 16, 32 or 64")
		return badKeySize
	}
	// any other checks?

	return nil
}
