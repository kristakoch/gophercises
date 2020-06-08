package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// NewFileVault is a factory function for FileVaults.
func NewFileVault(encodingKey string, filePath string) (*FileVault, error) {
	if len(encodingKey) != 32 {
		return nil, errors.New("encoding key must have a length of 32 (hex-encoded 16 character string)")
	}
	if filePath == "" {
		return nil, errors.New("filepath cannot be empty")
	}

	// Make sure there is a secrets file at filePath.
	if err := assertSecretsFile(filePath); err != nil {
		return nil, fmt.Errorf("error creating secrets file at %s, err: %s", filePath, err.Error())
	}

	var fv FileVault
	fv.fp = filePath
	fv.key = encodingKey

	return &fv, nil
}

// FileVault holds data for getting & storing
// encrypted data locally.
type FileVault struct {
	fp  string
	key string
}

// AllSecrets holds all stored secrets.
type AllSecrets struct {
	Secrets []Secret
}

// Secret defines the structure of stored secrets.
type Secret struct {
	Key   string
	Value string
}

// Set adds a key-value pair to the FileVault.
func (fv *FileVault) Set(key, val string) error {

	// Grab the current secrets.
	fb, err := fv.readEncrypted()
	if err != nil {
		return err
	}

	var curSecrets AllSecrets
	if len(fb) > 0 {
		if err := json.Unmarshal(fb, &curSecrets); err != nil {
			return err
		}
	}

	if key == "" || val == "" || key == "=" || val == "=" {
		return fmt.Errorf("failed to create new secret with key '%s' and value '%s' because either key or value is invalid", key, val)
	}

	// Add the new secret, updating the key value
	// if it already exists.
	var found bool
	for i, s := range curSecrets.Secrets {
		if s.Key == key {
			log.Printf("===> found duplicate - updating value for key '%s' to '%s'", key, val)
			curSecrets.Secrets[i].Value = val
			found = true
		}
	}

	if !found {
		log.Printf("===> adding new secret '%s: %s'", key, val)
		curSecrets.Secrets = append(curSecrets.Secrets,
			Secret{
				Key:   key,
				Value: val,
			})
	}

	secretBytes, err := json.Marshal(curSecrets)
	if err != nil {
		return err
	}

	// Encode and write the secrets into to the file.
	if err = fv.writeEncrypted(secretBytes); err != nil {
		return err
	}

	return nil
}

// Get retrieves the value associated with the key.
func (fv *FileVault) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("key not given for lookup in secrets store")
	}

	// Grab the current secrets.
	fb, err := fv.readEncrypted()
	if err != nil {
		return "", err
	}

	if len(fb) == 0 {
		return "", errors.New("secrets store is empty")
	}

	var curSecrets AllSecrets
	if err := json.Unmarshal(fb, &curSecrets); err != nil {
		return "", err
	}

	for _, s := range curSecrets.Secrets {
		if s.Key == key {
			return s.Value, nil
		}
	}

	return "", fmt.Errorf("result not found for key '%s'", key)
}

// Delete removes the entry with a matching key
// from the secrets store.
func (fv *FileVault) Delete(key string) error {

	// Grab the current secrets.
	fb, err := fv.readEncrypted()
	if err != nil {
		return err
	}

	var curSecrets AllSecrets
	if len(fb) > 0 {
		if err := json.Unmarshal(fb, &curSecrets); err != nil {
			return err
		}
	}

	if key == "" {
		return errors.New("failed to delete secret because key value is empty")
	}

	// Add the new secret, updating the key value
	// if it already exists.
	var found bool
	for i, s := range curSecrets.Secrets {
		if s.Key == key {
			log.Printf("===> deleting secret with key '%s'", key)
			curSecrets.Secrets = append(curSecrets.Secrets[:i], curSecrets.Secrets[i+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("no secret with key '%s' in secrets store")
	}

	secretBytes, err := json.Marshal(curSecrets)
	if err != nil {
		return err
	}

	// Encode and write the secrets into to the file.
	if err = fv.writeEncrypted(secretBytes); err != nil {
		return err
	}

	return nil
}

// ListAll logs all currently stored secrets.
func (fv *FileVault) ListAll() error {
	// Grab the current secrets.
	fb, err := fv.readEncrypted()
	if err != nil {
		return err
	}

	if len(fb) == 0 {
		return errors.New("secrets store is empty")
	}

	var curSecrets AllSecrets
	if err := json.Unmarshal(fb, &curSecrets); err != nil {
		return err
	}

	for _, s := range curSecrets.Secrets {
		fmt.Println(s.Key + ": " + s.Value)
	}

	return nil
}

// readEncrypted gets the file bytes and decrypts them.
func (fv *FileVault) readEncrypted() ([]byte, error) {

	// Retreive bytes with the secrets.
	fb, err := ioutil.ReadFile(fv.fp)
	if err != nil {
		return fb, err
	}

	// If file is empty, no need to decrypt.
	if len(fb) == 0 {
		return nil, nil
	}

	// Must convert to base 16 string
	// for decode to work.
	b16Bytes := fmt.Sprintf("%x", fb)

	// Create the necessary ingredients for decryption.
	key, _ := hex.DecodeString(fv.key)
	ciphertext, _ := hex.DecodeString(b16Bytes)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short, err: %s", err.Error())
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// writeEncrypted creates a file that contains the bytes of
// an encrypted string.
func (fv *FileVault) writeEncrypted(plaintext []byte) error {
	key, _ := hex.DecodeString(fv.key)

	// Create the necessary ingredients for encryption.
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt.
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Once encrypted, the secrets can be written
	// to the file.
	if err := ioutil.WriteFile(fv.fp, ciphertext, 0777); err != nil {
		return err
	}

	return nil
}

// assertSecretsFile makes sure a file exists
// to add secrets into.
func assertSecretsFile(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {

		// Create the path.
		fp, _ := filepath.Split(path)
		os.MkdirAll(fp, os.ModePerm)

		// Create the file.
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		log.Println("===> created secrets file at", path)
	}

	return nil
}
