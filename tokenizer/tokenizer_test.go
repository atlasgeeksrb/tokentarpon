package tokenizer

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCreateToken(t *testing.T) {
	testScenarios := []struct {
		givenValue, givenDomain string
		expectedErrorMsg        string
	}{
		{
			givenValue:  "this is a perfectly acceptable string to be tokenized",
			givenDomain: "mydomain",
		},
		{
			givenValue:       "",
			givenDomain:      "mydomain",
			expectedErrorMsg: ErrEmptyValue.Error(),
		},
		{
			givenValue:       "oh my, this message is bigger than the max 01234567890 - oh my, this message is bigger than the max 01234567890 - oh my, this message is bigger than the max 01234567890 - oh my, this message is bigger than the max 01234567890 - oh my, this message is bigger than the max 01234567890 - oh my, this message is bigger than the max 01234567890",
			givenDomain:      "mydomain",
			expectedErrorMsg: ErrValueTooBig.Error(),
		},
	}

	for i, scenario := range testScenarios {

		t.Run(strconv.Itoa(i), func(t *testing.T) {

			// 1./2. given/do
			tok, createerr := CreateToken(scenario.givenDomain, scenario.givenValue)

			// 3.1 expect error message
			if nil != createerr && len(scenario.expectedErrorMsg) > 0 && scenario.expectedErrorMsg != createerr.Error() {
				t.Error(createerr.Error())
				return
			}

			// 3.2 try get normally
			createdtoken, geterr := GetToken(scenario.givenDomain, tok.Uuid)
			if nil != geterr {
				t.Error(geterr.Error())
				return
			} else if want, got := scenario.givenValue, createdtoken.Value; want != got {
				t.Error(fmt.Printf("expect token value %#v but got %#v", want, got))
				return
			}

			// 3.3 try get with wrong domain
			_, geterr2 := GetToken("thisisnotmybeautifuldomain", tok.Uuid)
			if want, got := ErrNoMatchingToken, geterr2; want != got {
				t.Error(fmt.Printf("expect error %#v but got %#v", want, got))
			}

			// 3.3 try delete
			deleteerr := DeleteToken(tok.DomainUuid, tok.Uuid)
			if nil != deleteerr {
				t.Error(deleteerr.Error())
			}

			// 3.4 try delete token that doesn't exist
			deleteerr2 := DeleteToken("thisisnotmybeautifuldomain", tok.Uuid)
			if want, got := ErrNoMatchingToken, deleteerr2; want != got {
				t.Error(fmt.Printf("expect error %#v but got %#v", want, got))
			}

		})
	}
}

func TestDeleteToken(t *testing.T) {
	t.Fatal("argggh")
}

func TestGetToken(t *testing.T) {
	t.Fatal("argggh")
}
