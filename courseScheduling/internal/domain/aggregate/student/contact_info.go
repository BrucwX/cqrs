package student

import (
	"net/mail"
	"regexp"
)

// ContactInfo 联系方式值对象（手机号必填且符合规则，邮箱选填）
type ContactInfo struct {
	phone string
	email string
}

var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

func NewContactInfo(phone, email string) (ContactInfo, error) {
	if !phoneRegex.MatchString(phone) {
		return ContactInfo{}, ErrInvalidPhone
	}
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			return ContactInfo{}, ErrInvalidEmail
		}
	}
	return ContactInfo{phone: phone, email: email}, nil
}

func (c ContactInfo) Phone() string { return c.phone }
func (c ContactInfo) Email() string { return c.email }
