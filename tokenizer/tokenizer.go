package tokenizer

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"tokentarpon/crypto"
	"tokentarpon/tokenizer/datastore"

	"github.com/google/uuid"
)

var DefaultMaxLength int = 255
var MaxRecords int64 = 100
var UnitTest = false
var tokenRecordType = "token"

type Token struct {
	Uuid           string `bson:"uuid" json:"uuid"`
	DomainUuid     string `bson:"domainUuid" json:"domainUuid"`
	Value          string `bson:"value" json:"value"`
	EncryptedValue string `bson:"encryptedValue" json:"encryptedValue"`
	IsDeleted      bool   `json:"isDeleted" bson:"isDeleted"`
	DocumentType   string `bson:"documentType" json:"documentType"`
	Version        string `bson:"version" json:"version"`
	Created        int64  `json:"created" bson:"created"`
	Updated        int64  `json:"updated" bson:"updated"`
	Check          string `bson:"check" json:"check"`
}

type Token_v001 struct {
	Uuid           string `bson:"uuid" json:"uuid"`
	DomainUuid     string `bson:"domainUuid" json:"domainUuid"`
	Value          string `bson:"value" json:"value"`
	EncryptedValue string `bson:"encryptedValue" json:"encryptedValue"`
	IsDeleted      bool   `json:"isdeleted" bson:"isdeleted"`
	DocumentType   string `bson:"documentType" json:"documentType"`
	Version        string `bson:"version" json:"version"`
	Created        int64  `json:"created" bson:"created"`
	Updated        int64  `json:"updated" bson:"updated"`
	Check          string `bson:"check" json:"check"`
}

type TokenQuery struct {
	DomainUuid string   `bson:"domainUuid" json:"domainUuid"`
	Uuids      []string `bson:"uuids" json:"uuids"`
}

type TokenError struct {
	Token Token  `bson:"token" json:"token"`
	Error string `bson:"error" json:"error"`
}

var (
	ErrEmptyValue      = errors.New("cannot store empty value")
	ErrValueTooBig     = errors.New("value too large for storage")
	ErrNoMatchingToken = errors.New("no token found for provided domain and token id")
)

var CollectionName string = "community"

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

	if UnitTest {
		return tok, nil
	}

	datastore.CollectionName = CollectionName
	errInsert := datastore.InsertRecord(tokenRecordType, tok)
	return tok, errInsert
}

func CreateTokens(domainUuid string, tokens []Token) ([]Token, []TokenError) {
	var createdTokens []Token
	var errorTokens []TokenError
	datastore.CollectionName = CollectionName

	for _, tokenObj := range tokens {
		if len(strings.TrimSpace(tokenObj.DomainUuid)) == 0 {
			e := TokenError{Token: tokenObj, Error: "Missing Domain ID"}
			errorTokens = append(errorTokens, e)
		} else if tokenObj.DomainUuid != domainUuid {
			e := TokenError{Token: tokenObj, Error: "Invalid Domain ID"}
			errorTokens = append(errorTokens, e)
		} else if len(strings.TrimSpace(tokenObj.Value)) == 0 {
			e := TokenError{Token: tokenObj, Error: "Missing Token Value"}
			errorTokens = append(errorTokens, e)
		} else {
			aUuid := uuid.New()
			tokenObj.Uuid = aUuid.String()
			tokenObj.Created = time.Now().Unix()

			if UnitTest {
				createdTokens = append(createdTokens, tokenObj)
			} else {
				errInsert := datastore.InsertRecord(tokenRecordType, &tokenObj)
				if errInsert != nil {
					errmsg := fmt.Sprint(errInsert)
					e := TokenError{Token: tokenObj, Error: errmsg}
					errorTokens = append(errorTokens, e)
				} else {
					createdTokens = append(createdTokens, tokenObj)
				}
			}
		}
	}
	return createdTokens, errorTokens
}

func GetToken(domainUuid string, tokenUuid string) (Token, error) {
	var tok Token
	err := errors.New("data: need domain id, token id")
	if len(strings.TrimSpace(domainUuid)) == 0 {
		return tok, err
	}
	if len(strings.TrimSpace(tokenUuid)) == 0 {
		return tok, err
	}

	if UnitTest {
		tok.DomainUuid = domainUuid
		tok.Uuid = tokenUuid
		return tok, nil
	}

	filter := datastore.MakeDomainQuery(domainUuid, "uuid", tokenUuid, true)
	datastore.CollectionName = CollectionName
	geterr := datastore.GetRecord(filter, &tok)
	return tok, geterr
}

func DeleteToken(domainUuid string, tokenUuid string) (Token, error) {
	var empty Token
	err := errors.New("data: token incomplete, need domain id, token id, value")
	if len(strings.TrimSpace(domainUuid)) == 0 {
		return empty, err
	}
	if len(strings.TrimSpace(tokenUuid)) == 0 {
		return empty, err
	}

	var tokenObj Token = Token{
		DomainUuid: domainUuid,
		Uuid:       tokenUuid,
		IsDeleted:  true,
	}
	if UnitTest {
		return tokenObj, nil
	}

	filter := datastore.MakeDomainQuery(domainUuid, "uuid", tokenUuid, false)
	datastore.CollectionName = CollectionName
	geterr := datastore.GetRecord(filter, &tokenObj)
	if geterr != nil {
		return empty, geterr
	}

	tokenObj.IsDeleted = true
	updateResult, err := datastore.UpdateRecord(tokenRecordType, filter, "and", tokenObj)
	updatedToken := updateResult.(Token)
	return updatedToken, err
}

func CreateMultiTokenQuery(tokenQuery TokenQuery) []datastore.DataQueryGroup {

	var filters = make([]datastore.DataQueryGroup, 2)
	var nvq datastore.DataQueryGroup

	nvq.Operator = "and"
	nvq.DataQueries = make([]datastore.DataQuery, 1)
	nvq.DataQueries[0].FieldName = "domainUuid"
	nvq.DataQueries[0].FieldValue = tokenQuery.DomainUuid
	nvq.DataQueries[0].Wildcard = false
	nvq.DataQueries[0].CaseSensitive = true
	filters[0] = nvq

	nvq.Operator = "or"
	nvq.DataQueries = make([]datastore.DataQuery, len(tokenQuery.Uuids))
	for idx, uuid := range tokenQuery.Uuids {
		nvq.DataQueries[idx].FieldName = "uuid"
		nvq.DataQueries[idx].FieldValue = uuid
		nvq.DataQueries[idx].Wildcard = false
		nvq.DataQueries[idx].CaseSensitive = true
	}
	filters[1] = nvq
	return filters
}

func GetTokens(domainUuid string, start int64, limit int64) ([]Token, error) {
	var empty, tokens []Token
	err := errors.New("data: need domain id")
	if len(strings.TrimSpace(domainUuid)) == 0 {
		return empty, err
	}

	if UnitTest {
		return empty, err
	}

	datastore.CollectionName = CollectionName
	//record := &Token{}
	//my := &My{}
	var token Token
	filter := datastore.MakeSimpleQuery("domainUuid", domainUuid, false)
	records, geterr := datastore.GetRecords(filter, "and", start, limit, token)
	if geterr != nil {
		return nil, datastore.ErrQueryError
	}

	for _, x := range records {
		tokens = append(tokens, x.(Token))
	}

	return tokens, nil
}

func GetTokenValues(tokenQuery TokenQuery) ([]string, error) {
	var empty []string
	//var token Token
	err := errors.New("data: need domain id")
	if len(strings.TrimSpace(tokenQuery.DomainUuid)) == 0 {
		return empty, err
	}
	err = errors.New("data: need uuids")
	if len(tokenQuery.Uuids) == 0 {
		return empty, err
	}

	if UnitTest {
		return empty, err
	}

	filter := CreateMultiTokenQuery(tokenQuery)
	datastore.CollectionName = CollectionName
	var token Token
	records, geterr := datastore.GetRecords(filter, "and", 0, MaxRecords, token)
	if geterr != nil {
		return nil, datastore.ErrQueryError
	}

	tokenValues := make([]string, 0)
	// return the token values at the same indices as their uuids were presented
	for _, uuid := range tokenQuery.Uuids {
		for _, tv := range records {
			t := tv.(Token)
			if uuid == t.Uuid {
				tokenValues = append(tokenValues, t.Value)
			}
		}
	}
	return tokenValues, nil
}

func EncryptValue(plaintext string, key string) (string, error) {
	return crypto.EncryptAES(plaintext, key)
}

func DecryptValue(encryptedValue string, key string) (string, error) {
	return crypto.DecryptAES(encryptedValue, key)
}

func EncryptValues(plaintext []string, key string) []string {
	encryptedValues := make([]string, 0)
	for _, plain := range plaintext {
		encrypted, err := EncryptValue(plain, key)
		if nil == err {
			encryptedValues = append(encryptedValues, encrypted)
		}
	}
	return encryptedValues
}

func DecryptValues(encryptedValues []string, key string) []string {
	plaintext := make([]string, 0)
	for _, enc := range encryptedValues {
		decrypted, err := DecryptValue(enc, key)
		if nil == err {
			plaintext = append(plaintext, decrypted)
		}
	}
	return plaintext
}
