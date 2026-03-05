package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"
)

func LoadEncryptedEnv(filePath, passphrase string) error {
	// Baca konten terenkripsi dari file .env.enc
	encodedContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Decode konten terenkripsi dari base64
	encryptedContent, err := base64.StdEncoding.DecodeString(string(encodedContent))
	if err != nil {
		return err
	}

	// Dekripsi konten
	decryptedContent, err := Decrypt(encryptedContent, passphrase)
	if err != nil {
		return err
	}

	// Tulis konten yang didekripsi ke file sementara .env
	tempFilePath := ".env.temp"
	err = ioutil.WriteFile(tempFilePath, decryptedContent, 0644)
	if err != nil {
		return err
	}
	defer os.Remove(tempFilePath) // Hapus file sementara setelah selesai

	// Muat variabel lingkungan dari file sementara .env
	return godotenv.Overload(tempFilePath)
}

func Decrypt(data []byte, passphrase string) ([]byte, error) {
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
