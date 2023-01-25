package notifyslack

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("NotifySlack", notifySlack)
}

// helloHTTP is an HTTP Cloud Function with a request parameter.
// curl -X POST -H 'Content-type: application/json' --data '{"text":"Hello, World!"}' https://hooks.slack.com/services/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
// https://api.slack.com/apps/A04LK1W2MCK/incoming-webhooks?
func notifySlack(writer http.ResponseWriter, r *http.Request) {
	url := "https://hooks.slack.com/services/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"
	var jsonStr = []byte(`{"text":"Triggering Illuminating Calculation Wrap Up"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(writer, "err is %v", err)
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	fmt.Fprintf(writer, "status recieved from slack notification POST is %v\n", resp.Status)
}
