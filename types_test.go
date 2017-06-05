package gerrit_test

import (
	"encoding/json"
	"testing"
	"github.com/andygrunwald/go-gerrit"
)

func TestTypesNumber_String(t *testing.T) {
	number := gerrit.Number("7")
	if number.String() != "7" {
		t.Fatalf("%s != 7", number.String())
	}
}

func TestTypesNumber_Int(t *testing.T) {
	number := gerrit.Number("7")
	integer, err := number.Int()
	if err != nil {
		t.Fatal(err)
	}
	if integer != 7 {
		t.Fatalf("%d != 7", integer)
	}
}

func TestTypesNumber_UnmarshalJSON_String(t *testing.T) {
	var number gerrit.Number
	if err := json.Unmarshal([]byte(`"7"`), &number); err != nil {
		t.Fatal(err)
	}
	if number.String() != "7" {
		t.Fatalf("%s != 7", number.String())
	}
}

func TestTypesNumber_UnmarshalJSON_Int(t *testing.T) {
	var number gerrit.Number
	if err := json.Unmarshal([]byte("7"), &number); err != nil {
		t.Fatal(err)
	}
	integer, err := number.Int()
	if err != nil {
		t.Fatal(err)
	}
	if integer != 7 {
		t.Fatalf("%d != 7", integer)
	}
}
