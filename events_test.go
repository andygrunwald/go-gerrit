package gerrit_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/andygrunwald/go-gerrit"
)

var (
	fakeEvents = []byte(`
	{"submitter":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"newRev":"0000000000000000000000000000000000000000","patchSet":{"number":"1","revision":"0000000000000000000000000000000000000000","parents":["0000000000000000000000000000000000000000"],"ref":"refs/changes/1/1/1","uploader":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"createdOn":1470000000,"author":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"isDraft":false,"kind":"TRIVIAL_REBASE","sizeInsertions":10,"sizeDeletions":0},"change":{"project":"test","branch":"master","id":"Iffffffffffffffffffffffffffffffffffffffff","number":"1","subject":"subject","owner":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"url":"https://localhost/1","commitMessage":"commitMessage\n\nline2\n\nChange-Id: Iffffffffffffffffffffffffffffffffffffffff\n","status":"MERGED"},"type":"change-merged","eventCreatedOn":1470000000}
	{"author":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"comment":"Patch Set 1:\n\n(2 comments)\n\nSome comment","patchSet":{"number":"1","revision":"0000000000000000000000000000000000000000","parents":["0000000000000000000000000000000000000000"],"ref":"refs/changes/1/1/1","uploader":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"createdOn":1470000000,"author":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"isDraft":false,"kind":"REWORK","sizeInsertions":4,"sizeDeletions":-2},"change":{"project":"test","branch":"master","id":"Iffffffffffffffffffffffffffffffffffffffff","number":"1","subject":"subject","owner":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"url":"https://localhost/1","commitMessage":"commitMessage\n\nChange-Id: Iffffffffffffffffffffffffffffffffffffffff\n","status":"NEW"},"type":"comment-added","eventCreatedOn":1470000000}`)

	fakeEventsWithError = []byte(`
	{"submitter":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"newRev":"0000000000000000000000000000000000000000","patchSet":{"number":"1","revision":"0000000000000000000000000000000000000000","parents":["0000000000000000000000000000000000000000"],"ref":"refs/changes/1/1/1","uploader":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"createdOn":1470000000,"author":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"isDraft":false,"kind":"TRIVIAL_REBASE","sizeInsertions":10,"sizeDeletions":0},"change":{"project":"test","branch":"master","id":"Iffffffffffffffffffffffffffffffffffffffff","number":"1","subject":"subject","owner":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"url":"https://localhost/1","commitMessage":"commitMessage\n\nline2\n\nChange-Id: Iffffffffffffffffffffffffffffffffffffffff\n","status":"MERGED"},"type":"change-merged","eventCreatedOn":1470000000}
	{"author":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"comment":"Patch Set 1:\n\n(2 comments)\n\nSome comment","patchSet":{"number":"1","revision":"0000000000000000000000000000000000000000","parents":["0000000000000000000000000000000000000000"],"ref":"refs/changes/1/1/1","uploader":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"createdOn":1470000000,"author":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"isDraft":false,"kind":"REWORK","sizeInsertions":4,"sizeDeletions":-2},"change":{"project":"test","branch":"master","id":"Iffffffffffffffffffffffffffffffffffffffff","number":"1","subject":"subject","owner":{"name":"Foo Bar","email":"fbar@example.com","username":"fbar"},"url":"https://localhost/1","commitMessage":"commitMessage\n\nChange-Id: Iffffffffffffffffffffffffffffffffffffffff\n","status":"NEW"},"type":"comment-added","eventCreatedOn":1470000000}
	{"author":1}`)
)

func TestEventsLogService_GetEvents_NoDateRange(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/plugins/events-log/events/", func(writer http.ResponseWriter, request *http.Request) {
		if _, err := writer.Write(fakeEvents); err != nil {
			t.Error(err)
		}
	})

	options := &gerrit.EventsLogOptions{}
	events, _, _, err := testClient.EventsLog.GetEvents(options)
	if err != nil {
		t.Error(err)
	}

	if len(events) != 2 {
		t.Error("Expected 2 events")
	}

	// Basic test
	for i, event := range events {
		switch i {
		case 0:
			if event.Type != "change-merged" {
				t.Error("Expected event type to be `change-merged`")
			}
		case 1:
			if event.Type != "comment-added" {
				t.Error("Expected event type to be `comment-added`")
			}
		}
	}
}

func TestEventsLogService_GetEvents_DateRangeFromAndTo(t *testing.T) {
	setup()
	defer teardown()

	to := time.Now()
	from := to.AddDate(0, 0, -7)

	testMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()

		fromFormat := from.Format("2006-01-02 15:04:05")
		if query.Get("t1") != fromFormat {
			t.Errorf("%s != %s", query.Get("t1"), fromFormat)
		}

		toFormat := to.Format("2006-01-02 15:04:05")
		if query.Get("t2") != toFormat {
			t.Errorf("%s != %s", query.Get("t2"), toFormat)
		}

		if _, err := writer.Write(fakeEvents); err != nil {
			t.Error(err)
		}
	})

	options := &gerrit.EventsLogOptions{From: from, To: to}
	_, _, _, err := testClient.EventsLog.GetEvents(options)
	if err != nil {
		t.Error(err)
	}
}

func TestEventsLogService_GetEvents_DateRangeFromOnly(t *testing.T) {
	setup()
	defer teardown()

	to := time.Now()
	from := to.AddDate(0, 0, -7)

	testMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()

		fromFormat := from.Format("2006-01-02 15:04:05")
		if query.Get("t1") != fromFormat {
			t.Errorf("%s != %s", query.Get("t1"), fromFormat)
		}

		if query.Get("t2") != "" {
			t.Error("Did not expect t2 to be set")
		}

		if _, err := writer.Write(fakeEvents); err != nil {
			t.Error(err)
		}
	})

	options := &gerrit.EventsLogOptions{From: from}
	_, _, _, err := testClient.EventsLog.GetEvents(options)
	if err != nil {
		t.Error(err)
	}
}

func TestEventsLogService_GetEvents_DateRangeToOnly(t *testing.T) {
	setup()
	defer teardown()

	to := time.Now()

	testMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()

		toFormat := to.Format("2006-01-02 15:04:05")
		if query.Get("t2") != toFormat {
			t.Errorf("%s != %s", query.Get("t2"), toFormat)
		}

		if query.Get("t1") != "" {
			t.Error("Did not expect t1 to be set")
		}

		if _, err := writer.Write(fakeEvents); err != nil {
			t.Error(err)
		}
	})

	options := &gerrit.EventsLogOptions{To: to}
	_, _, _, err := testClient.EventsLog.GetEvents(options)
	if err != nil {
		t.Error(err)
	}
}

func TestEventsLogService_GetEvents_UnmarshalError(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if _, err := writer.Write(fakeEventsWithError); err != nil {
			t.Error(err)
		}
	})

	options := &gerrit.EventsLogOptions{IgnoreUnmarshalErrors: true}
	events, _, failures, err := testClient.EventsLog.GetEvents(options)
	if err != nil {
		t.Error(err)
	}
	if len(failures) != 1 {
		t.Error("Expected 1 failures")
	}
	if len(events) != 2 {
		t.Error(len(events))
	}
}
