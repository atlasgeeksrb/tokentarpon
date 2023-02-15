package tokenizer

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var DefaultMaxLength int = 255

var (
	ErrEmptyValue      = errors.New("cannot store empty value")
	ErrValueTooBig     = errors.New("value too large for storage")
	ErrNoMatchingToken = errors.New("no token found for provided domain and token id")
)

type Token struct {
	Uuid           string `bson:"uuid" json:"uuid"`
	DomainUuid     string `bson:"domainUuid" json:"domainUuid"`
	Value          string `bson:"value" json:"value"`
	EncryptedValue string `bson:"encryptedValue" json:"encryptedValue"`
	Created        int64  `json:"created" bson:"created"`
	IsDeleted      bool   `json:"isdeleted" bson:"isdeleted"`
}

func CreateToken(domainUuid string, value string) (Token, error) {
	var tok Token

	err := errors.New("data: token incomplete, need domain id, value")
	if len(strings.TrimSpace(domainUuid)) == 0 {
		return tok, err
	}
	if len(strings.TrimSpace(value)) == 0 {
		return tok, err
	}
	aUuid := uuid.New()
	tok.DomainUuid = domainUuid
	tok.Value = value
	tok.Uuid = aUuid.String()
	tok.Created = time.Now().Unix()

	//@todo store token
	return tok, nil
}

func GetToken(domainUuid string, tokenUuid string) (Token, error) {
	var tok Token
	return tok, nil
}

func DeleteToken(domainuuid string, tokenUuid string) error {
	return nil
}
