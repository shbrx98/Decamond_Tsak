package utils

import (
    "errors"
    "regexp"
    "strings"
)

// NormalizePhoneNumber enforces a simple E.164-like format.
// For production-quality validation use github.com/nyaruka/phonenumbers.
func NormalizePhoneNumber(p string) (string, error) {
    p = strings.TrimSpace(p)
    p = strings.ReplaceAll(p, " ", "")
    p = strings.ReplaceAll(p, "-", "")
    p = strings.ReplaceAll(p, "(", "")
    p = strings.ReplaceAll(p, ")", "")

    if strings.HasPrefix(p, "00") {
        p = "+" + p[2:]
    }
    if !strings.HasPrefix(p, "+") {
        // assume it's a national number; reject or prepend default country code
        // To be safe, require leading '+'
        return "", errors.New("phone must be in E.164 format starting with +")
    }
    re := regexp.MustCompile(`^\+\d{8,15}$`)
    if !re.MatchString(p) {
        return "", errors.New("invalid E.164 phone number")
    }
    return p, nil
}