package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNon200Response non 200 status code in response
	ErrSessionCreateFailed = errors.New("Failed to create AWS Session")

	HTMLPageBody =`
    <html>
    <head>
      <meta charset="utf-8"/>
    </head>
    <body>
      <h1>Thanks</h1>
      <p>We received your submission</p>
      <p>Reference: %s</p>
      </p>
    </body>
    </html>
  `
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Log the request
	requestJSON, _ := json.Marshal(request)
	log.Println(string(requestJSON))

	// Get variables
	var requestMap map[string]interface{}
	_ = json.Unmarshal(requestJSON, &requestMap)
	requestContext := requestMap["requestContext"].(map[string]interface{})
	requestID := requestContext["requestId"].(string)

	// Format HTML
	formattedHTML := fmt.Sprintf(HTMLPageBody, requestID)

	// Open AWS Session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Upload the File
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("UPLOAD_S3_BUCKET")),
		Key: aws.String(requestID),
		Body: bytes.NewReader(requestJSON),
	})
	if err != nil {
		// Print the error and exit.
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body: formattedHTML,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
