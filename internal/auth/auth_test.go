package auth

import "testing"

func TestPasswordHashing(t *testing.T) {
	password := "mypassword123"

	// Test hashing
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	if hash == password {
		t.Error("Hash should not be equal to password")
	}

	// Test correct password
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Error("Password check should succeed with correct password")
	}

	// Test incorrect password
	err = CheckPasswordHash("wrongpassword", hash)
	if err == nil {
		t.Error("Password check should fail with wrong password")
	}
}
