//go:build !test
// +build !test

package repository

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	projectID       = os.Getenv("GCP_PROJECT_ID")
	firestoreClient *firestore.Client
)

func init() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	firestoreClient = client
}

func CheckIfExists(id string) (bool, error) {
	doc, err := firestoreClient.Collection("urls").Doc(id).Get(context.Background())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, err
	}
	return doc.Exists(), nil
}

func SaveURL(id, longURL string) error {
	_, err := firestoreClient.
		Collection("urls").
		Doc(id).
		Set(context.Background(), map[string]interface{}{
			"longUrl": longURL,
		})
	return err
}

func SaveRequest(data map[string]interface{}) error {
	isoTime := time.Now().Format("2006-01-02T15:04:05.000000Z07:00")
	_, err := firestoreClient.
		Collection("requests").
		Doc(isoTime).
		Set(context.Background(), data)
	return err
}

func GetRedirect(id string) (string, []byte, error) {
	doc, err := firestoreClient.
		Collection("urls").
		Doc(id).
		Get(context.Background())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", nil, nil
		} else {
			return "", nil, err
		}
	}
	if !doc.Exists() {
		return "", nil, nil
	}

	data := doc.Data()
	longUrl, _ := data["url"].(string)
	email, emailOk := data["email"].(string)
	name, nameOk := data["name"].(string)

	if !nameOk {
		log.Printf("name(%d) %s", len(name), name)
	}

	if nameOk && name != "" {
		vcard := VCard{Name: name, Email: email, Homepage: longUrl}
		return "", []byte(vcard.String()), nil
	}

	if emailOk && email != "" {
		subject, _ := data["subject"].(string)
		body, _ := data["body"].(string)
		return MailTo{Address: email, Subject: subject, Body: body}.String(), nil, nil
	}

	return longUrl, nil, nil
}
