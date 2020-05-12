package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")

	HTMLPageBody =`
  <html>
  <head>
    <meta charset="utf-8"/>
  </head>
  <body>
    <form method="POST">
      Please enter your name:
      <input type="text" name="name"/>
      <br/>
      <input type="submit" />
    </form>
  </body>
  </html>
`
	HTMLSubmissionBody = `
  <html>
  <head>
    <meta charset="utf-8"/>
  </head>
  <body>
    <h1>Thanks</h1>
    <p>We received your submission</p>
  </body>
  </html>
`
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Log the request
	requestJSON, _ := json.Marshal(request)
	log.Println(string(requestJSON))

	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	// Get the appropriate HTML body
	var body string
	if request.HTTPMethod == "GET" {
		body = HTMLPageBody
	} else {
		body = HTMLSubmissionBody
	}

	return events.APIGatewayProxyResponse{
		Body: body,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
