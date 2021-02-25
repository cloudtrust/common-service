package validation

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// ObfuscatePhoneNumber hides the middle part of a phone number
// It keeps the country code and the 2 last digits
func ObfuscatePhoneNumber(phoneNumber string) string {
	// Use country ZZ to make parsing fail if phone number is not parsable
	var metadata, err = phonenumbers.Parse(phoneNumber, "ZZ")
	var runes = []rune(phoneNumber)
	var length = len(runes)
	var prefix string
	if err != nil {
		// Should not execute this if number is a valid E164
		if length < 6 {
			return phoneNumber
		}
		prefix = string(runes[0:3])
	} else {
		prefix = fmt.Sprintf("+%d", metadata.GetCountryCode())
	}
	var suffix = string(runes[length-2:])
	var middle = length - len(prefix) - len(suffix)
	return fmt.Sprintf("%s%s%s", prefix, strings.Repeat("*", middle), suffix)
}
