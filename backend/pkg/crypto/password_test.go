package crypto

import (
	"testing"
)

func TestPassword(t *testing.T) {
	hash, err := HashPassword("Admin@123")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Fresh hash: %s", hash)
	if !CheckPassword("Admin@123", hash) {
		t.Fatal("Fresh hash check failed")
	}

	hardcoded := "$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK"
	if !CheckPassword("Admin@123", hardcoded) {
		t.Errorf("Hardcoded hash check failed!")
	} else {
		t.Logf("Hardcoded hash check SUCCESS!")
	}
}
