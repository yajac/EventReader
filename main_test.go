package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler(t *testing.T) {

	if testing.Short() {
		request := events.APIGatewayProxyRequest{}
		var paths = map[string]string{
			"location": "Richmond",
		}

		request.PathParameters = paths
		expectedResponse := events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: "\"eventId\":",
		}

		response, err := Handler(request)

		assert.Equal(t, response.Headers, expectedResponse.Headers)
		assert.Contains(t, response.Body, expectedResponse.Body)
		assert.Contains(t, response.Body, "\"description\":")
		assert.Contains(t, response.Body, "\"title\":")
		assert.Equal(t, err, nil)
	}

}
