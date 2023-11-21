package hashing

import (
	"testing"
)

func TestGenerateCreds(t *testing.T) {
	pass, salt := GenerateCredentials()
	if len(pass) != PASS_LEN {
		t.Fatalf(`GenerateCredentials password %s, want length %d`, pass, PASS_LEN)
	}
	if len(salt) != SALT_LEN {
		t.Fatalf(`GenerateCredentials salt %s, want length %d`, salt, SALT_LEN)
	}
}

func TestHashCorrect(t *testing.T) {
	pass, salt := GenerateCredentials()
	hash1 := GenerateHash(pass, salt)
	hash2 := GenerateHash(pass, salt)
	if hash1 != hash2 {
		t.Fatalf(`GenerateHash should match but got: %s, %s`, hash1, hash2)
	}
}

func TestHashDifferentSalt(t *testing.T) {
	pass, salt := GenerateCredentials()
	hash1 := GenerateHash(pass, salt)
	hash2 := GenerateHash(pass, "abcdefg")
	if hash1 == hash2 {
		t.Fatalf(`GenerateHash should not match with different salt but did: %s, %s`, hash1, hash2)
	}
}

func TestHashDifferentPassword(t *testing.T) {
	pass, salt := GenerateCredentials()
	hash1 := GenerateHash(pass, salt)
	hash2 := GenerateHash("password", salt)
	if hash1 == hash2 {
		t.Fatalf(`GenerateHash should not match with different password but did: %s, %s`, hash1, hash2)
	}
}
