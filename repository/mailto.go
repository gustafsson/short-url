package repository

import (
	"fmt"
	"net/url"
)

type MailTo struct {
	Address string
	Subject string
	Body    string
}

var _ fmt.Stringer = (*VCard)(nil)

func (e MailTo) String() string {
	// Initialize the base mailto link with the email address
	mailto := "mailto:" + url.QueryEscape(e.Address)

	// Prepare the query parameters
	params := url.Values{}
	if e.Subject != "" {
		params.Add("subject", e.Subject)
	}
	if e.Body != "" {
		params.Add("body", e.Body)
	}

	// If there are any parameters, add them to the mailto link
	if len(params) > 0 {
		mailto += "?" + params.Encode()
	}

	return mailto
}
