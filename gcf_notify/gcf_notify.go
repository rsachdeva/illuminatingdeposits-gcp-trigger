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

	return nil
}
