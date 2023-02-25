// Package gcf_notify  Google cloud function is triggered by message published to pubsub topic `deltaanalyticstopic` by
// gcf_analytics cloud function. Uses google cloud secret manager grabbing sendgrid api key to send email message
// decoded from base64 encoding per JSON message schema received from pubsub.
// Only sends from verified sender.
// Sends to env email address for now, in future version will be able to
// send to email address specified in request
// Depends on gcf_analytics cloud function deployed and having run successfully. See System Diagram for more details.
//
//nolint:revive,stylecheck // using underscore in package name for clarity
package gcf_notify

import (
	"context"
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

// required as per google cloud functions version 2
// https://cloud.google.com/functions/docs/writing/write-http-functions
//
//nolint:gochecknoinits
func init() {
	functions.CloudEvent("NotifyInvestorOfDelta", notifyInvestorOfDelta)
}

const versionName = "projects/illuminatingdeposits-gcp/secrets/sendgrid-api-key/versions/latest"

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// notifyInvestorOfDelta consumes a CloudEvent message and extracts the Pub/Sub message.
func notifyInvestorOfDelta(ctx context.Context, e event.Event) error {
	log.Println("e.Data is", string(e.Data()))

	var msg MessagePublishedData

	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	log.Printf("full msg is %#v", msg)

	dataRecvd := string(msg.Message.Data) // Automatically decoded from base64.
	log.Printf("******Presenting to you the delta notification\n: %s", dataRecvd)

	sendgridKey, err := accessSecretVersion(ctx, versionName)
	if err != nil {
		return fmt.Errorf("accessSecretVersion: %w", err)
	}

	log.Println("sendGridKey is", sendgridKey)

	err = sendEmailThroughSendGrid(sendgridKey, dataRecvd)
	if err != nil {
		return fmt.Errorf("sendEmailThroughSendGrid: %w", err)
	}

	log.Println("email sent")

	return nil
}
