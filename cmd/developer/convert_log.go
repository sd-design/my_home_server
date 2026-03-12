package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
)

func encryptGCM(plaintext, key byte) (byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make(byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptGCM(ciphertext, key byte) (byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}

func main() {

	keyFile := flag.String("key", "./encrypt.key", "key file")
	logFile := ".development.log"

	dataKey, err := os.ReadFile(keyFile)
	if err != nil {
		fmt.Println("Ошибка чтения файла ключа:", err)
		return
	}
	contentKey := string(dataKey)

	key := byte(contentKey) // 32 байта для AES-256

	file, err := os.Open("example.txt")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении log файла:", err)
	}

	plaintext := byte(content)

	ciphertext, err := encryptGCM(plaintext, key)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	fmt.Printf("Encrypted: %x\n", ciphertext)

	decrypted, err := decryptGCM(ciphertext, key)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}
	fmt.Printf("Decrypted: %s\n", decrypted)

}
