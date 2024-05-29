package repository

import (
	"fmt"
)

type VCard struct {
	Name     string
	Email    string
	Homepage string
}

var _ fmt.Stringer = (*VCard)(nil)

func (vcard VCard) String() string {
	vCard := "BEGIN:VCARD\nVERSION:3.0\n"

	if vcard.Name != "" {
		vCard += fmt.Sprintf("FN:%s\n", vcard.Name)
	}
	if vcard.Email != "" {
		vCard += fmt.Sprintf("EMAIL:%s\n", vcard.Email)
	}
	if vcard.Homepage != "" {
		vCard += fmt.Sprintf("URL:%s\n", vcard.Homepage)
	}

	vCard += "END:VCARD"
	return vCard
}

