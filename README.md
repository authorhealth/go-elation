# go-elation

## Getting Started

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/authorhealth/go-elation"
)

func main() {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	tokenURL := os.Getenv("TOKEN_URL")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	baseURL := os.Getenv("BASE_URL")

	client := elation.NewHttpClient(httpClient, tokenURL, clientID, clientSecret, baseURL)

	res := &elation.Response[[]*elation.Patient]{}
	var err error

	res, _, err = client.Patients().Find(context.Background(), &elation.FindPatientsOptions{
		Pagination: res.PaginationNextWithLimit(1),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Results[0].FirstName)
}
```
