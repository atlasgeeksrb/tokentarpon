package tokencrypto

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	Uuid         string `bson:"uuid" json:"uuid"`
	Value        string `bson:"value" json:"value"`
	DocumentType string `bson:"documentType" json:"documentType"`
	Version      string `bson:"version" json:"version"`
	Created      int64  `json:"created" bson:"created"`
	Updated      int64  `json:"updated" bson:"updated"`
	Check        string `bson:"check" json:"check"`
}

func TestGetHashForByteArray(t *testing.T) {

	nownow := time.Now().Unix()
	var x Foo = Foo{
		Uuid:         "abc123",
		Value:        "some rando value",
		DocumentType: "Foo",
		Version:      "001",
		Created:      nownow,
		Updated:      nownow,
	}
	j, _ := json.Marshal(x)
	hash := GetHashForByteArray(j)
	x.Check = hash

	// hash can't be empty
	assert.NotEqual(t, strings.TrimSpace(hash), "")

	// re-hashing the same data must yield an identical hash
	assert.Equal(t, GetHashForByteArray(j), hash)

	// hashing data with slight change must yield a different hash
	x.Updated = time.Now().Unix()
	j, _ = json.Marshal(x)
	assert.NotEqual(t, strings.TrimSpace(GetHashForByteArray(j)), hash)

}

func TestGetHashForString(t *testing.T) {

	teststring := "This is a delishously hashable string!"
	hash := GetHashForString(teststring)

	// hash can't be empty
	assert.NotEqual(t, hash, "")

	// re-hashing the same data must yield an identical hash
	assert.Equal(t, GetHashForString(teststring), hash)

	// hashing data with slight change must yield a different hash
	teststring += " "
	assert.NotEqual(t, GetHashForString(teststring), hash)

}

func TestEncryptAES(t *testing.T) {

	plaintextValue := "this is some value!"

	// The key argument should be the AES key,
	// either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	key := ";kldfpo87-28374isu;dfjhZXJCVG786"

	encryptedValue, encryptionErr := EncryptAES(plaintextValue, key)
	if nil != encryptionErr {
		t.Fatal(encryptionErr)
	}

	// encryptedValue can't be empty
	assert.NotEqual(t, strings.TrimSpace(encryptedValue), "")

	// encryptedValue can't equal plaintextValue
	assert.NotEqual(t, strings.TrimSpace(encryptedValue), plaintextValue)

	// it's hard to test encryption without confirming that decryption is also working
	decryptedValue, decryptionErr := DecryptAES(encryptedValue, key)
	if nil != decryptionErr {
		t.Fatal(decryptionErr)
	}

	// decryptedValue must equal plaintextValue
	assert.Equal(t, strings.TrimSpace(decryptedValue), plaintextValue)

	// decryptedValue must not equal encryptedValue
	assert.NotEqual(t, decryptedValue, encryptedValue)

	// less happy paths
	// test with wrong size key
	key = "anewkey"
	_, encryptionErr = EncryptAES(plaintextValue, key)
	assert.Equal(t, encryptionErr, errors.New("key size must be 16, 32 or 64"))

	//@todo
	// what's the largest string we can encrypt? should know this, and test against it

}

func TestDecryptAES(t *testing.T) {
	t.Skip()
}
