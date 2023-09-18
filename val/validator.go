package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var(
	//os caracteres de a-z e 0-9 podem aparecer mais de uma vez na string
	isValidUsername = regexp.MustCompile(`ˆ[a-z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`ˆ[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(value string, minLength, maxLength int) error  {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters",minLength,maxLength)
	}

	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 70); err != nil {
		return err
	}
	if isValidUsername(value) {
		return fmt.Errorf("must contains only letters, digits or undercore")
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 70); err != nil {
		return err
	}
	if isValidFullname(value) {
		return fmt.Errorf("must contains only letters, or spaces")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value,6,100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value,3,200); err != nil {
		return err
	}

	//checar se é um email valido
	if _,err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("should be a valid email address")
	}
	return nil
}