package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func main() {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		// One might use structural loggig for better output in CloudWatch logs console.
	)
	if err != nil {
		panic(err)
	}

	cwLogsService := cloudwatchlogs.NewFromConfig(cfg)
	h := newHandler(cwLogsService)
	lambda.Start(cfn.LambdaWrap(h))
}

type QueryService interface {
	PutQueryDefinition(
		ctx context.Context,
		params *cloudwatchlogs.PutQueryDefinitionInput,
		optsFns ...func(*cloudwatchlogs.Options),
	) (*cloudwatchlogs.PutQueryDefinitionOutput, error)

	DeleteQueryDefinition(
		ctx context.Context,
		params *cloudwatchlogs.DeleteQueryDefinitionInput,
		optsFns ...func(*cloudwatchlogs.Options),
	) (*cloudwatchlogs.DeleteQueryDefinitionOutput, error)
}

func newHandler(queryService QueryService) cfn.CustomResourceFunction {
	return func(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
		queryName, ok := event.ResourceProperties["QueryName"].(string)
		if !ok {
			return physicalResourceID, data, errors.New("QueryName parameter not found")
		}

		queryString, ok := event.ResourceProperties["QueryString"].(string)
		if !ok {
			return physicalResourceID, data, errors.New("QueryString parameter not found")
		}

		if event.RequestType == cfn.RequestCreate {
			out, err := queryService.PutQueryDefinition(
				ctx,
				&cloudwatchlogs.PutQueryDefinitionInput{
					Name:        aws.String(queryName),
					QueryString: aws.String(queryString),
				},
			)
			if err != nil {
				return physicalResourceID, data, err
			}

			return *out.QueryDefinitionId, data, nil
		}

		if event.RequestType == cfn.RequestUpdate {
			out, err := queryService.PutQueryDefinition(
				ctx,
				&cloudwatchlogs.PutQueryDefinitionInput{
					Name:              aws.String(queryName),
					QueryString:       aws.String(queryString),
					QueryDefinitionId: aws.String(event.PhysicalResourceID),
				},
			)
			if err != nil {
				return physicalResourceID, data, err
			}

			return *out.QueryDefinitionId, data, nil
		}

		if event.RequestType == cfn.RequestDelete {
			_, err := queryService.DeleteQueryDefinition(
				ctx,
				&cloudwatchlogs.DeleteQueryDefinitionInput{
					QueryDefinitionId: aws.String(event.PhysicalResourceID),
				},
			)
			if err != nil {
				return physicalResourceID, data, err
			}
		}

		return
	}
}
