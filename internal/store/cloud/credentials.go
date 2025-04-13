package cloud

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CredentialsManager struct {
	configDir string
	key       []byte // 32 bytes for AES-256
}

// generateKey creates a secure encryption key and stores it separately
func generateKey(keyPath string) ([]byte, error) {
	key := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("generating key: %w", err)
	}

	// Save the key with restricted permissions
	if err := os.WriteFile(keyPath, key, 0600); err != nil {
		return nil, fmt.Errorf("saving key: %w", err)
	}

	return key, nil
}

func loadOrGenerateKey(keyPath string) ([]byte, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return generateKey(keyPath)
		}
		return nil, fmt.Errorf("reading key: %w", err)
	}
	return key, nil
}

func NewCredentialsManager() (*CredentialsManager, error) {
	configBase, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("getting config directory: %w", err)
	}

	configDir := filepath.Join(configBase, "merlion")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("creating config directory: %w", err)
	}

	// Load or generate encryption key
	keyPath := filepath.Join(configDir, ".key")
	key, err := loadOrGenerateKey(keyPath)
	if err != nil {
		return nil, fmt.Errorf("handling encryption key: %w", err)
	}

	return &CredentialsManager{
		configDir: configDir,
		key:       key,
	}, nil
}

func (cm *CredentialsManager) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(cm.key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generating nonce: %w", err)
	}

	// Encrypt and prepend nonce
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (cm *CredentialsManager) decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(cm.key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (cm *CredentialsManager) SaveCredentials(creds *Credentials) error {
	if creds == nil {
		return fmt.Errorf("No credentials supply (nil)")
	}

	// Marshal credentials
	data, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("marshaling credentials: %w", err)
	}

	// Encrypt
	encrypted, err := cm.encrypt(data)
	if err != nil {
		return fmt.Errorf("encrypting credentials: %w", err)
	}

	// Save to file with restricted permissions
	credFile := filepath.Join(cm.configDir, "credentials.json")
	if err := os.WriteFile(credFile, encrypted, 0600); err != nil {
		return fmt.Errorf("writing credentials file: %w", err)
	}

	return nil
}

func (cm *CredentialsManager) LoadCredentials() (*Credentials, error) {
	credFile := filepath.Join(cm.configDir, "credentials.json")
	encrypted, err := os.ReadFile(credFile)
	if err != nil {
		return nil, fmt.Errorf("reading credentials file: %w", err)
	}

	// Decrypt
	data, err := cm.decrypt(encrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypting credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("unmarshaling credentials: %w", err)
	}

	return &creds, nil
}
