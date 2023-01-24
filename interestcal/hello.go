package interestcal

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloHTTP", helloHTTP)
}

// helloHTTP is an HTTP Cloud Function with a request parameter.
func helloHTTP(writer http.ResponseWriter, r *http.Request) {
	var person struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		fmt.Fprint(writer, "Salut, World!!!!RS")
		return
	}
	if person.Name == "" {
		fmt.Fprint(writer, "Salut, World!!!!RS")
		return
	}
	fmt.Fprintf(writer, "Salut, %s!!\n", html.EscapeString(person.Name))
}
