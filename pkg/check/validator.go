package check

import (
	"errors"
	"regexp"
	"time"
)

func ValidateCarYear(year int) error {
	if year <= 0 || year > time.Now().Year()+1 {
		return errors.New("year is not valid")
	}
	return nil
}

func ValidateGmailCustomer(e string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,3}$`)
    return emailRegex.MatchString(e)
}


func ValidatePhoneNumberOfCustomer(phone string) bool {
	if 12 < len(phone) && len(phone) <= 13{
		phoneregex:= regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
		return phoneregex.MatchString(phone)
	}
	return false
}