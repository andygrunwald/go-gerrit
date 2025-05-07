package gerrit_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

func TestChangesService_ListFiles(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.String(), "/changes/123/revisions/456/files/?base=7"; got != want {
			t.Errorf("request URL:\ngot:  %q\nwant: %q", got, want)
		}
		fmt.Fprint(w, `{
		  "/COMMIT_MSG": {
		    "status": "A",
		    "lines_inserted": 7,
		    "size_delta": 551,
		    "size": 551
		  },
		  "gerrit-server/RefControl.java": {
		    "lines_inserted": 5,
		    "lines_deleted": 3,
		    "size_delta": 98,
		    "size": 23348
		  }
		}`)
	}))
	defer ts.Close()

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	got, _, err := client.Changes.ListFiles(ctx, "123", "456", &gerrit.FilesOptions{
		Base: "7",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]gerrit.FileInfo{
		"/COMMIT_MSG": {
			Status:        "A",
			LinesInserted: 7,
			SizeDelta:     551,
			Size:          551,
		},
		"gerrit-server/RefControl.java": {
			LinesInserted: 5,
			LinesDeleted:  3,
			SizeDelta:     98,
			Size:          23348,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("client.Changes.ListFiles:\ngot:  %+v\nwant: %+v", got, want)
	}
}

func TestChangesService_ListFilesReviewed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.String(), "/changes/123/revisions/456/files/?q=abc&reviewed=true"; got != want {
			t.Errorf("request URL:\ngot:  %q\nwant: %q", got, want)
		}
		fmt.Fprint(w, `["/COMMIT_MSG","gerrit-server/RefControl.java"]`)
	}))
	defer ts.Close()

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	got, _, err := client.Changes.ListFilesReviewed(ctx, "123", "456", &gerrit.FilesOptions{
		Q: "abc",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"/COMMIT_MSG", "gerrit-server/RefControl.java"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("client.Changes.ListFilesReviewed:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestChangesService_GetPatch(t *testing.T) {
	rawPatch := `diff --git a/COMMIT_MSG b/COMMIT_MSG
index 123..456 100644
--- a/COMMIT_MSG
+++ b/COMMIT_MSG
@@ -1,1 +1,1 @@
-Old subject
+New subject for A
diff --git a/fileA.txt b/fileA.txt
new file mode 100644
index 0000000..abc 100644
--- /dev/null
+++ b/fileA.txt
@@ -0,0 +1 @@
+Content for A
`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.String(), "/changes/123/revisions/456/patch"; got != want {
			t.Errorf("request URL:\ngot:  %q\nwant: %q", got, want)
		}
		encodedPatchContent := base64.StdEncoding.EncodeToString(
			[]byte(rawPatch))
		fmt.Fprint(w, encodedPatchContent)
	}))
	defer ts.Close()

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	got, _, err := client.Changes.GetPatch(ctx, "123", "456", nil)
	if err != nil {
		t.Fatal(err)
	}
	want := base64.StdEncoding.EncodeToString([]byte(rawPatch))
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("client.Changes.GetPatch:\ngot:  %q\nwant: %q", *got, want)
	}
}
