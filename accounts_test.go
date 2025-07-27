package gerrit_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestAddSSHKey tests the addition of an SSH key to an account.
func TestAddSSHKey(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/accounts/self/sshkeys", func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Ensure Content-Type is text/plain
		if r.Header.Get("Content-Type") != "text/plain" {
			t.Errorf("Expected Content-Type 'text/plain', got %s", r.Header.Get("Content-Type"))
		}

		// Read body and validate SSH key
		expectedSSHKey := "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA0T...YImydZAw== john.doe@example.com"
		body, _ := io.ReadAll(r.Body)
		receivedSSHKey := strings.TrimSpace(string(body))

		if receivedSSHKey != expectedSSHKey {
			t.Errorf("Expected SSH key '%s', but received '%s'", expectedSSHKey, receivedSSHKey)
		}

		// Mock successful JSON response
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"seq": 2,
			"ssh_public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA0T...YImydZAw== john.doe@example.com",
			"encoded_key": "AAAAB3NzaC1yc2EAAAABIwAAAQEA0T...YImydZAw==",
			"algorithm": "ssh-rsa",
			"comment": "john.doe@example.com",
			"valid": true
		}`)
	})

	ctx := context.Background()
	sshKey := "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA0T...YImydZAw== john.doe@example.com"

	// Use testClient.Accounts instead of undefined Accounts variable
	keyInfo, _, err := testClient.Accounts.AddSSHKey(ctx, "self", sshKey)
	if err != nil {
		t.Fatalf("AddSSHKey returned error: %v", err)
	}

	// Verify SSH key information in the response
	if keyInfo.SSHPublicKey != sshKey {
		t.Errorf("Expected SSH key '%s', got '%s'", sshKey, keyInfo.SSHPublicKey)
	}

	if keyInfo.Valid != true {
		t.Errorf("Expected key validity to be true, got false")
	}

	if keyInfo.Comment != "john.doe@example.com" {
		t.Errorf("Expected comment 'john.doe@example.com', got '%s'", keyInfo.Comment)
	}
}
