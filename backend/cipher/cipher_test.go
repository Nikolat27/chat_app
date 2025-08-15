package cipher

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// Test with valid environment variable
	os.Setenv("ENCRYPTION_SECRET_KEY", "test-secret-key-for-testing-only")
	defer os.Unsetenv("ENCRYPTION_SECRET_KEY")

	cipher := New()
	if cipher == nil {
		t.Fatal("Expected cipher instance, got nil")
	}
	if cipher.Aead == nil {
		t.Fatal("Expected AEAD cipher, got nil")
	}
}

func TestNewPanic(t *testing.T) {
	// Test panic when environment variable is missing
	os.Unsetenv("ENCRYPTION_SECRET_KEY")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when ENCRYPTION_SECRET_KEY is not set")
		}
	}()

	New()
}

func TestEncryptDecrypt(t *testing.T) {
	// Set up test environment
	os.Setenv("ENCRYPTION_SECRET_KEY", "test-secret-key-for-testing-only")
	defer os.Unsetenv("ENCRYPTION_SECRET_KEY")

	cipher := New()

	tests := []struct {
		name      string
		plainText string
	}{
		{"Empty string", ""},
		{"Simple text", "Hello, World!"},
		{"Long text", "This is a very long text that should be encrypted and decrypted properly. It contains multiple sentences and various characters like !@#$%^&*()_+-=[]{}|;':\",./<>?"},
		{"Unicode text", "Hello 世界! Привет мир! 🌍"},
		{"Special characters", "!@#$%^&*()_+-=[]{}|;':\",./<>?"},
		{"Numbers", "1234567890"},
		{"Mixed content", "Hello123!@#世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plainTextBytes := []byte(tt.plainText)

			// Encrypt
			encrypted, err := cipher.Encrypt(plainTextBytes)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			// Verify encrypted data is different from plain text
			if string(encrypted) == tt.plainText {
				t.Error("Encrypted data should not be the same as plain text")
			}

			// Verify encrypted data is longer than plain text (due to nonce)
			if len(encrypted) <= len(plainTextBytes) {
				t.Error("Encrypted data should be longer than plain text due to nonce")
			}

			// Decrypt
			decrypted, err := cipher.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			// Verify decrypted data matches original
			if string(decrypted) != tt.plainText {
				t.Errorf("Decrypted data doesn't match original. Expected: %s, Got: %s", tt.plainText, string(decrypted))
			}
		})
	}
}

func TestEncryptDecryptConsistency(t *testing.T) {
	// Set up test environment
	os.Setenv("ENCRYPTION_SECRET_KEY", "test-secret-key-for-testing-only")
	defer os.Unsetenv("ENCRYPTION_SECRET_KEY")

	cipher := New()
	plainText := "Test message for consistency check"
	plainTextBytes := []byte(plainText)

	// Encrypt the same data multiple times
	encrypted1, err := cipher.Encrypt(plainTextBytes)
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	encrypted2, err := cipher.Encrypt(plainTextBytes)
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	// Encrypted data should be different due to random nonce
	if string(encrypted1) == string(encrypted2) {
		t.Error("Encrypted data should be different due to random nonce")
	}

	// But both should decrypt to the same plain text
	decrypted1, err := cipher.Decrypt(encrypted1)
	if err != nil {
		t.Fatalf("First decryption failed: %v", err)
	}

	decrypted2, err := cipher.Decrypt(encrypted2)
	if err != nil {
		t.Fatalf("Second decryption failed: %v", err)
	}

	if string(decrypted1) != string(decrypted2) {
		t.Error("Both decryptions should produce the same result")
	}

	if string(decrypted1) != plainText {
		t.Errorf("Decrypted data doesn't match original. Expected: %s, Got: %s", plainText, string(decrypted1))
	}
}

func TestDecryptInvalidData(t *testing.T) {
	// Set up test environment
	os.Setenv("ENCRYPTION_SECRET_KEY", "test-secret-key-for-testing-only")
	defer os.Unsetenv("ENCRYPTION_SECRET_KEY")

	cipher := New()

	tests := []struct {
		name        string
		invalidData []byte
		expectError bool
	}{
		{"Empty data", []byte{}, true},
		{"Too short data", []byte{1, 2, 3, 4, 5}, true},
		{"Exactly 12 bytes (only nonce)", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, true},
		{"Corrupted data", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cipher.Decrypt(tt.invalidData)
			if tt.expectError && err == nil {
				t.Error("Expected error for invalid data, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestDecryptWrongKey(t *testing.T) {
	// Create two ciphers with different keys
	os.Setenv("ENCRYPTION_SECRET_KEY", "first-secret-key")
	cipher1 := New()

	os.Setenv("ENCRYPTION_SECRET_KEY", "second-secret-key")
	cipher2 := New()

	// Reset to first key for cleanup
	os.Setenv("ENCRYPTION_SECRET_KEY", "first-secret-key")
	defer os.Unsetenv("ENCRYPTION_SECRET_KEY")

	plainText := "Test message"
	plainTextBytes := []byte(plainText)

	// Encrypt with first cipher
	encrypted, err := cipher1.Encrypt(plainTextBytes)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Try to decrypt with second cipher (different key)
	_, err = cipher2.Decrypt(encrypted)
	if err == nil {
		t.Error("Expected error when decrypting with wrong key, got nil")
	}
}

func TestLargeData(t *testing.T) {
	// Set up test environment
	os.Setenv("ENCRYPTION_SECRET_KEY", "test-secret-key-for-testing-only")
	defer os.Unsetenv("ENCRYPTION_SECRET_KEY")

	cipher := New()

	// Create large data
	largeData := make([]byte, 10000)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// Encrypt
	encrypted, err := cipher.Encrypt(largeData)
	if err != nil {
		t.Fatalf("Encryption of large data failed: %v", err)
	}

	// Decrypt
	decrypted, err := cipher.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decryption of large data failed: %v", err)
	}

	// Verify
	if len(decrypted) != len(largeData) {
		t.Errorf("Decrypted data length doesn't match. Expected: %d, Got: %d", len(largeData), len(decrypted))
	}

	for i := range largeData {
		if decrypted[i] != largeData[i] {
			t.Errorf("Data mismatch at position %d. Expected: %d, Got: %d", i, largeData[i], decrypted[i])
			break
		}
	}
}

// Benchmark tests
func BenchmarkEncrypt(b *testing.B) {
	os.Setenv("ENCRYPTION_SECRET_KEY", "test-secret-key-for-benchmarking")
	defer os.Unsetenv("ENCRYPTION_SECRET_KEY")

	cipher := New()
	data := []byte("Benchmark test data for encryption")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cipher.Encrypt(data)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}
	}
}

func BenchmarkDecrypt(b *testing.B) {
	os.Setenv("ENCRYPTION_SECRET_KEY", "test-secret-key-for-benchmarking")
	defer os.Unsetenv("ENCRYPTION_SECRET_KEY")

	cipher := New()
	data := []byte("Benchmark test data for decryption")
	encrypted, err := cipher.Encrypt(data)
	if err != nil {
		b.Fatalf("Setup encryption failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cipher.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}

func BenchmarkEncryptDecrypt(b *testing.B) {
	os.Setenv("ENCRYPTION_SECRET_KEY", "test-secret-key-for-benchmarking")
	defer os.Unsetenv("ENCRYPTION_SECRET_KEY")

	cipher := New()
	data := []byte("Benchmark test data for encrypt-decrypt cycle")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypted, err := cipher.Encrypt(data)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}

		_, err = cipher.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
} 