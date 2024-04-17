package check

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func ValidateCarYear(year int) error {
	if year <= 0 || year > time.Now().Year()+1 {
		return errors.New("year is not valid")
	}
	return nil
}

func Validategmail(gmail string) error {
	// Check if the email ends with either "@gmail.com" or "@gmail.ru"
	if !strings.HasSuffix(gmail, "@gmail.com") && !strings.HasSuffix(gmail, "@gmail.ru") {
		return errors.New("email address must end with @gmail.com or @gmail.ru")
	}

	return nil
}

func ValidatePassword(newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("=====password length must be at least 8 characters")
	}

	var hasUppercase, hasLowercase, hasDigit, hasSymbol bool

	for _, char := range newPassword {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsNumber(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	if !hasUppercase {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLowercase {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if !hasSymbol {
		return errors.New("password must contain at least one symbol")
	}

	return nil
}

func ValidatePhone(phone string) error {
	pattern := `^\d{12}$`

	re := regexp.MustCompile(pattern)

	if !re.MatchString(phone) {
		return errors.New("phone number must consist of 12 digits separated by spaces")
	}

	return nil
}
