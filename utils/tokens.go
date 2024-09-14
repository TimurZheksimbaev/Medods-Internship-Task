package utils

import (
	"fmt"
	"log"
	"net/smtp"

	"golang.org/x/crypto/bcrypt"
)

func HashRefreshToken(token string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedToken), nil
}

func CompareRefreshTokens(hashedToken, token string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
	return err == nil
}

func SendEmail(oldIP, newIP, userID string) {
	// данные отправителя
	from := "email@example.com"
	password := "password123"

	// адрес получателя
	to := []string{
		"recipient@example.com",
	}

	// SMTP-сервер и порт
	smtpHost := "smtp.gmail.com" 
	smtpPort := "587"

	// авторизуемся
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// сообщение
	message := []byte(fmt.Sprintf("Ваш IP адрес изменился. Старый адрес: %s Новый адрес: %s", oldIP, newIP))

	// отправляем сообщение
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Printf("Failed to send email to user: %s", userID)
	}

	log.Printf("Successfully sent email to user: %s !", userID)
}



func LogExit(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func Log(err error) {
	if err != nil {
		log.Println(err)
	}
}

func LogMessage(message string) {
	log.Println(message)
}