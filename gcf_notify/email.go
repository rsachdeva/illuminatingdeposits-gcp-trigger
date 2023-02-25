//nolint:revive,stylecheck // using underscore in package name for clarity
package gcf_notify

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var errDataCorrruption = errors.New("Data corruption detected")

func sendEmailThroughSendGrid(sendgridKey, dataRecvd string) error {
	senderEmail := os.Getenv("SENDER_EMAIL")
	from := mail.NewEmail("Rohit Sachdeva", senderEmail)
	subject := "Your Illuminating Deposits Delta Calculation Analytics is Ready!"
	receiverEmail := os.Getenv("RECEIVER_EMAIL")
	to := mail.NewEmail("Testing Email", receiverEmail)
	plainTextContent := fmt.Sprintf("Here are latest overall delta by date: %s", dataRecvd)
	htmlContent := fmt.Sprintf("Here are latest overall delta by date: <br> <strong>%s</strong>", dataRecvd)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendgridKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("sendgrid.Send: %w", err)
	}

	log.Println(response.StatusCode)
	log.Println(response.Body)
	log.Println(response.Headers)

	return nil
}

// accessSecretVersion gets information about the given secret version. It does not
// include the payload data.
func accessSecretVersion(ctx context.Context, name string) (string, error) {
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"
	// Create the client.
	// ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %writer", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)

	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		return "", errDataCorrruption
	}

	// WARNING: Do not print the secret in a production environment - this snippet
	// is showing how to access the secret material.
	// fmt.Fprintf(writer, "Plaintext: %s\n", string(result.Payload.Data))

	return string(result.Payload.Data), nil
}
