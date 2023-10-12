package helper

import "net/mail"

func ValidateEmailAddr(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
