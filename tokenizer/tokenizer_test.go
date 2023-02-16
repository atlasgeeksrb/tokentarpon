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

	Test = true
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
	Test = true
	t.Fatal("argggh")
}

func TestGetToken(t *testing.T) {
	Test = true
	t.Fatal("argggh")
}

/*
// test fixtures
func TestAdd(t *testing.T) {
	// go test sets pwd as package being tested
	// my_folder is a relative path; consider using test_fixtures
	// use it to load config, model data, binary data etc.
	data := filepath.Join("my_folder", "add_data.json")

	// do stuff with data

	// golden files
	// compare a file output to a golden file - expected file output
	// simple use a .golden extension for a file
	// then
	for _, tc := range testcases {
		golden := filepath.Join("my_folder", "add_data.json")

		expected, _ := ioutil.Readfile(golden)
		if !bytes.Equal(actual, expected) {
			// fail
		}

	}

	// test helper functions should never return an error
	// just fail the test if there's an error
	// see t.helper

}

// over-parameterize structs, for us in testing
type ServerOpts struct {
	CachePath string
	Port int

	Test bool
	// consider having a test bool
	// use this for testing;
	// for example production code would
	// test this value and if set to true,
	// skip authentication
}

// testing sub-processes, options:
// 1. test the actual process
// in this case guard against the process not being available i.e. is git here?

// 2. mock the subprocess
// make the *exec.Cmd configurable, pass in a custom one
// let the caller modify it
func helperProcess (s ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs.append(cs, s...)
	env := []string{
		"GO_WANT_HELPER_PROCESS=1",
	}

	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(env, os.Environ()...)
	return cmd
}
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}

		args = args[1:]
	}
}

// include testing.go or testing_*.go files
// exported APIs that provide mocks, test harnesses, helpers
// allows other packages to test using this package,
// without reinventing components needed to use the package in a test
// e.g. TestConfig(t) provides valid testing config
// TestServer(t) (net.Addr, io.Closer) fully started in-memory server and a closer to close it
// ...and tests for things the user of our package might create,
// to confirm that they meet the spec required by our package

// go test is a great workflow tool
// explore this

// if you have a complex pluggable system
// you can write a custom framework
// within go test, rather than a separate test harness



func testTempFile(t *testing.T) (string, func()) {
	t.Helper()

	tf, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tf.Close()

	return tf.Name(), func() { os.Remove(tf.Name())}
}

func TestThing(t *testing.T) {
	tf, tfclose := testTempFile(t)
	defer tfclose()
}

// don't mock net.Conn! go ahead and create a network connection
func TestConn(t *testing.T) (client, server net.Conn) {
	t.Helper()

	// one-connection example
	// also easy to make a n-connection example
	// easy to test any protocol
	// easy to return the listener
	// easy to test IPv6 if needed

	ln, err := net.listen("tcp", "127.0.0.1:0")

	var server net.Conn
	go fnc() {
		defer ln.Close()
		server, err = ln.Accept()
	}()

	client, err := net.Dial("tcp", ln.Addr().String())
	return client, server
}

func TestGetToken(t *testing.T) {

	var stringValue = "Some test string"
	token, err := GetToken("mydomain", stringValue)
	if err != nil {
		t.Log("error should be nil", err)
		t.Fail()
	}
	if len(strings.TrimSpace(token.Uuid)) == 0 {
		t.Log("got empty token, should not be empty", err)
		t.Fail()
	}
}

func TestGetTokens(t *testing.T) {
	tokens, err := GetUsers()
	if err != nil {
		t.Log("error should be nil", err)
		t.Fail()
	}
	if len(tokens) == 0 {
		t.Log("no tokens found, there should be some tokens", err)
		t.Fail()
	}
}

func TestGetTokenValues(t *testing.T) {
	newUser := Token{
		Uuid:     "abc123",
		Username: existingUsername,
		Password: userPassword,
		Name:     userName,
	}
	// expect error when we re-use a username
	_, err := CreateUser(newUser)
	if nil == err {
		t.Log("error should not be nil for create user with existing name", err)
		t.Fail()
	}
	newUser.Username = "superbrandnewuser"
	token, err := CreateUser(newUser)
	if nil != err {
		t.Log("error should be nil for create user with new username", err)
		t.Fail()
	}
	if token.Uuid == newUser.Uuid {
		t.Log("new token uuid should be set by CreateToken function", err)
		t.Fail()
	}

}
func TestCreateToken(t *testing.T) {
}
func TestCreateTokens(t *testing.T) {
}
func TestDeleteToken(t *testing.T) {
}
func TestEncryptValues(t *testing.T) {
}
func TestDecryptValues(t *testing.T) {
}
*/
