package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
	assert.Len(t, hash, 60) // bcrypt hashes are 60 characters
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("")
	require.NoError(t, err)
	require.NotEmpty(t, hash)
}

func TestHashPassword_SpecialCharacters(t *testing.T) {
	passwords := []string{
		"!@#$%^&*()",
		"password with spaces",
		"Unicode: 🚀🎉",
		"1234567890",
		"a",
	}

	for _, password := range passwords {
		t.Run("password: "+password, func(t *testing.T) {
			hash, err := HashPassword(password)
			require.NoError(t, err)
			require.NotEmpty(t, hash)
			assert.NotEqual(t, password, hash)
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testPassword123"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	// Test correct password
	err = ComparePassword(hash, password)
	assert.NoError(t, err)

	// Test incorrect password
	err = ComparePassword(hash, "wrongPassword")
	assert.Error(t, err)
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	password := ""
	hash, err := HashPassword(password)
	require.NoError(t, err)

	// Test correct empty password
	err = ComparePassword(hash, password)
	assert.NoError(t, err)

	// Test non-empty password against empty password hash
	err = ComparePassword(hash, "somePassword")
	assert.Error(t, err)
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	password := "testPassword123"

	// Test with invalid hash
	err := ComparePassword("invalid_hash", password)
	assert.Error(t, err)

	// Test with empty hash
	err = ComparePassword("", password)
	assert.Error(t, err)

	// Test with malformed hash
	err = ComparePassword("not_a_bcrypt_hash", password)
	assert.Error(t, err)
}

func TestPasswordHashing_Consistency(t *testing.T) {
	password := "testPassword123"

	// Hash the same password multiple times
	hash1, err := HashPassword(password)
	require.NoError(t, err)

	hash2, err := HashPassword(password)
	require.NoError(t, err)

	// Hashes should be different (due to salt)
	assert.NotEqual(t, hash1, hash2)

	// But both should verify correctly
	err = ComparePassword(hash1, password)
	assert.NoError(t, err)

	err = ComparePassword(hash2, password)
	assert.NoError(t, err)
}

func TestPasswordHashing_Performance(t *testing.T) {
	password := "testPassword123"

	// Test that hashing doesn't take too long
	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	// Verify the hash works
	err = ComparePassword(hash, password)
	assert.NoError(t, err)
}

func TestPasswordVerification_EdgeCases(t *testing.T) {
	// Test long password (but within bcrypt limits)
	longPassword := "this_is_a_reasonable_length_password_for_testing"

	hash, err := HashPassword(longPassword)
	require.NoError(t, err)

	err = ComparePassword(hash, longPassword)
	assert.NoError(t, err)

	// Test password with null bytes
	nullPassword := "test\x00password"
	hash, err = HashPassword(nullPassword)
	require.NoError(t, err)

	err = ComparePassword(hash, nullPassword)
	assert.NoError(t, err)
}
