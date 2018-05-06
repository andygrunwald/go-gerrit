package gerrit_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/andygrunwald/go-gerrit"
)

func TestTimestamp(t *testing.T) {
	const jsonData = `{
	"subject": "net/http: write status code in Redirect when Content-Type header set",
	"created": "2018-05-04 17:24:39.000000000",
	"updated": "0001-01-01 00:00:00.000000000",
	"submitted": "2018-05-04 18:01:10.000000000",
	"_number": 111517
}
`
	type ChangeInfo struct {
		Subject   string            `json:"subject"`
		Created   gerrit.Timestamp  `json:"created"`
		Updated   gerrit.Timestamp  `json:"updated"`
		Submitted *gerrit.Timestamp `json:"submitted,omitempty"`
		Omitted   *gerrit.Timestamp `json:"omitted,omitempty"`
		Number    int               `json:"_number"`
	}
	ci := ChangeInfo{
		Subject:   "net/http: write status code in Redirect when Content-Type header set",
		Created:   gerrit.Timestamp{Time: time.Date(2018, 5, 4, 17, 24, 39, 0, time.UTC)},
		Updated:   gerrit.Timestamp{},
		Submitted: &gerrit.Timestamp{Time: time.Date(2018, 5, 4, 18, 1, 10, 0, time.UTC)},
		Omitted:   nil,
		Number:    111517,
	}

	// Try decoding JSON data into a ChangeInfo struct.
	var v ChangeInfo
	err := json.Unmarshal([]byte(jsonData), &v)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := v, ci; !reflect.DeepEqual(got, want) {
		t.Errorf("decoding JSON data into a ChangeInfo struct:\ngot:\n%v\nwant:\n%v", got, want)
	}

	// Try encoding a ChangeInfo struct into JSON data.
	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	e.SetIndent("", "\t")
	err = e.Encode(ci)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := buf.String(), jsonData; got != want {
		t.Errorf("encoding a ChangeInfo struct into JSON data:\ngot:\n%v\nwant:\n%v", got, want)
	}
}

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
