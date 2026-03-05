package controllers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

func EnkripEnv(c *gin.Context) {
	passphrase := "mNsRjOIdbyj1X2i6lLFJ5KE/evhYQIEz"
	// Baca konten file .env
	content, err := ioutil.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Enkripsi konten
	encryptedContent, err := encrypt(content, passphrase)
	if err != nil {
		log.Fatal(err)
	}

	// Simpan konten terenkripsi ke file baru
	encodedContent := base64.StdEncoding.EncodeToString(encryptedContent)
	err = ioutil.WriteFile(".env.enc", []byte(encodedContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}
