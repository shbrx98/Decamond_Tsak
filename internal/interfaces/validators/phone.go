package validators

import (
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/utils"
)

// ValidatePhone فرمت شماره را به E.164 نرمال و اعتبارسنجی می‌کند.
// خروجی: شماره نرمال‌شده یا خطا.
func ValidatePhone(phone string) (string, error) {
    return utils.NormalizePhoneNumber(phone)
}