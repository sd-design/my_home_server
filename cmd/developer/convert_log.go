package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
)

func encryptGCM(plaintext []byte, key string) ([]byte, error) {
	cipherKey := []byte(key) // Преобразуем строку в []byte

	block, err := aes.NewCipher(cipherKey)
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

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptGCM(ciphertext []byte, key string) ([]byte, error) {
	// 1. Преобразуем строковый ключ в срез байт
	cipherKey := []byte(key)

	// 2. Проверяем длину ключа (должна быть 16, 24 или 32 байта)
	keyLen := len(cipherKey)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, fmt.Errorf(
			"некорректная длина ключа: %d байт. Должно быть 16 (AES-128), 24 (AES-192) или 32 (AES-256) байта",
			keyLen,
		)
	}

	// 3. Создаём блочный шифр по ключу
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания шифра: %w", err)
	}

	// 4. Создаём AEAD (GCM) на основе блочного шифра
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации GCM: %w", err)
	}

	// 5. Извлекаем размер nonce (он фиксирован для данного шифра)
	nonceSize := gcm.NonceSize()

	// 6. Проверяем, что ciphertext достаточно длинный (nonce + зашифрованные данные)
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf(
			"шифротекст слишком короткий: %d байт, требуется минимум %d байт (nonce)",
			len(ciphertext),
			nonceSize,
		)
	}

	// 7. Разделяем nonce и остальную часть шифротекста
	nonce := ciphertext[:nonceSize]
	encryptedData := ciphertext[nonceSize:]

	// 8. Расшифровываем данные
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка расшифрования (возможно, неверный ключ или повреждённые данные): %w", err)
	}

	return plaintext, nil
}

func actionIn(keyFile *string, logFile *string) {
	dataKey, err := os.ReadFile(*keyFile)
	if err != nil {
		fmt.Println("Ошибка чтения файла ключа:", err)
		os.Exit(1)
	}
	contentKey := string(dataKey)

	//fmt.Println(contentKey)

	file, err := os.Open(*logFile)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		os.Exit(1)
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
	plaintext := []byte(content)

	ciphertext, err := encryptGCM(plaintext, contentKey)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	fmt.Printf("Encrypted: %x\n", ciphertext)

	hexString := hex.EncodeToString(ciphertext)
	errWrite := os.WriteFile("convert.log", []byte(hexString), 0644)
	if errWrite != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Файл успешно записан")

}

func actionOut(keyFile *string, logFile *string) {
	dataKey, err := os.ReadFile(*keyFile)
	if err != nil {
		fmt.Println("Ошибка чтения файла ключа:", err)
		os.Exit(1)
	}
	contentKey := string(dataKey)

	//fmt.Println(contentKey)

	file, err := os.Open("convert.log")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		os.Exit(1)
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

	fmt.Println(content)
	ciphertext, err := hex.DecodeString(content)
	plaintext2, err := decryptGCM(ciphertext, contentKey)
	if err != nil {
		fmt.Println("Ошибка расшифрования:", err)
		return
	}

	fmt.Printf("Расшифрованный текст:\n %s\n", plaintext2)

}

func main() {

	direction := flag.String("action", "NONE", "direction")
	keyFile := flag.String("key", "encrypt.key", "key file")
	logFile := flag.String("log", "development.log", "log file")

	flag.Parse()

	fmt.Println(*direction)

	if *direction == "NONE" {
		fmt.Println("Не указан параметр \"-action=\" направление кодирования/декодирования")
		os.Exit(1)
	}

	switch *direction {
	case "in":
		//Шишифруем
		actionIn(keyFile, logFile)
	case "out":
		//Расшифровываем
		actionOut(keyFile, logFile)

	}

}
