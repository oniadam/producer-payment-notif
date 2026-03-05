package controllers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DekripEnv(c *gin.Context) {
	passphrase := "mNsRjOIdbyj1X2i6lLFJ5KE/evhYQIEz"

	// Baca konten terenkripsi dari file .env.enc
	encodedContent, err := ioutil.ReadFile(".env.enc")
	if err != nil {
		log.Fatal(err)
	}

	// Decode konten terenkripsi dari base64
	encryptedContent, err := base64.StdEncoding.DecodeString(string(encodedContent))
	if err != nil {
		log.Fatal(err)
	}

	// Dekripsi konten
	decryptedContent, err := decrypt(encryptedContent, passphrase)
	if err != nil {
		log.Fatal(err)
	}

	// Simpan konten dekripsi ke file baru
	// err = ioutil.WriteFile(".env.dec", decryptedContent, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	envMap := parseEnvToJSON(string(decryptedContent))

	// Mengembalikan hasil dekripsi dalam JSON
	c.JSON(http.StatusOK, envMap)

}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Fungsi untuk mengonversi isi env ke dalam JSON
func parseEnvToJSON(envContent string) gin.H {
	lines := strings.Split(envContent, "\n")
	envMap := gin.H{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Lewati baris kosong atau komentar
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			envMap[key] = value
		}
	}

	return envMap
}
