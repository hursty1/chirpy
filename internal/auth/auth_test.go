package auth

import (
	"encoding/hex"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRefreshToken(t *testing.T){
	key := make([]byte, 32)
	hex_key := hex.EncodeToString(key)
	token, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("Error making token")
	}
	if hex_key == token {
		t.Fatalf("Key is the default value not cryptographic")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword1"
	hashed_pw, err := HashPassword(password)
	if err != nil {
		t.Errorf("Failed to hash password: %s and produced error: %s", password, err)
	}
	if hashed_pw == "" {
		t.Errorf("Password is a zero length string something has gone wrong while hashing %s", password)
	}
	// fmt.Println(hashed_pw)
	//Success
}

func TestCreateJWT(t *testing.T) {
	cases := []struct {
		userID uuid.UUID
		tokenSecret string
		expiresIn time.Duration
		
	}{
		{
		userID: uuid.New(),
		tokenSecret: "allyourbases",
		expiresIn: time.Duration(time.Second*10),
		},
	}

	for _, c := range cases {
		jwt, err := MakeJWT(c.userID, c.tokenSecret, c.expiresIn)
		if err != nil {
			t.Fatalf("Unable to create JWT %s", err)

		}
		userUUID, err := ValidateJWT(jwt,c.tokenSecret)
		if userUUID != c.userID {
			t.Fatalf("Unable to validate JWT %s", err)
		}
		
	}
}

func TestExpiredJWT(t *testing.T) {
	cases := []struct {
		userID uuid.UUID
		tokenSecret string
		expiresIn time.Duration
		
	}{
		{
		userID: uuid.New(),
		tokenSecret: "allyourbases",
		expiresIn: time.Duration(time.Minute*-10),
		},
	}

	for _, c := range cases {
		jwt, err := MakeJWT(c.userID, c.tokenSecret, c.expiresIn)
		if err != nil {
			t.Fatalf("Unable to create JWT %s", err)

		}
		userUUID, err := ValidateJWT(jwt,c.tokenSecret)
		if userUUID == c.userID {
			t.Fatalf("JWT Should be invalid %s", err)
		}
		
	}
}


func TestAuthHeader(t *testing.T) {
	c1:= http.Header{}
	c1.Set("authorization", "Bearer sdkdlsl3l12")
	cases := []struct {
		header http.Header
		expecting string
		
	}{
		{
		header: c1,
		expecting: "sdkdlsl3l12",
		},
	}

	for _, c := range cases {
		tok, err := GetBearerToken(c.header)
		if err != nil {
			t.Fatalf("Unable to get header: %s", err)
		}
		if tok != c.expecting {
			t.Fatalf("Token %s does not match expecting %s", tok, c.expecting)
		}
		
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword1"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Unable to hash password2: %s", err)
	}
	password2:= "lots12312312"
	hashed2, err := HashPassword(password2)
	if err != nil {
		t.Fatalf("Unable to hash password2: %s", err)
	}
	password3:= "lots12312312"
	hashed3, err := HashPassword("")
	if err != nil {
		t.Fatalf("Unable to hash password2: %s", err)
	}
	cases := []struct {
		name string
		input1 string
		input2 string
		expected bool
	}{
		{
		name: "Test Normal Password",
		input1: hashed,
		input2: password,
		expected: false,
		},
		{
		name: "Test Lots of numbers",
		input1: hashed2,
		input2: password2,
		expected: false,
		},
		{
		name: "Test Blank Password",
		input1: hashed3,
		input2: password3,
		expected: true,
		},
	}

	for _, c := range cases {
		err := CheckPasswordHash(c.input1, c.input2)
		if err == nil && c.expected{ 
			t.Errorf("%s: expected error but recieved: %s", c.name, err)
		}
		if err != nil && !c.expected{
			t.Errorf("%s: Error not expected but recieved error from checkpasswordhash %s", c.name, err)
		}
		//err and expected match
	}


}