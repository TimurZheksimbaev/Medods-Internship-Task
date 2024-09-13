package utils

import (
	"log"

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


func SendEmailWarning(userID string) {
	// Mock email sending
	log.Printf("Sending email warning to user %s\n", userID)
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