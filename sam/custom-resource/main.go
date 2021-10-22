package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	h := newHandler()
	lambda.Start(cfn.LambdaWrap(h))
}

// How does named return parameter work?
func newHandler() cfn.CustomResourceFunction {
	return func(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
		queryName, ok := event.ResourceProperties["QueryName"].(string)
		if !ok {
			return physicalResourceID, data, errors.New("QueryName parameter not found")
		}

		queryString, ok := event.ResourceProperties["QueryString"].(string)
		if !ok {
			return physicalResourceID, data, errors.New("QueryString parameter not found")
		}

		return
	}
}
