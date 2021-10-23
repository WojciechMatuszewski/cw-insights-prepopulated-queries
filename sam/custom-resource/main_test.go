package main

import (
	"context"
	"sam-app/custom-resource/mock"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()
	t.Run("Returns error when 'QueryName' parameter is not provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		queryService := mock.NewMockQueryService(ctrl)
		h := newHandler(queryService)

		_, _, err := h(ctx, cfn.Event{
			ResourceProperties: map[string]interface{}{
				"QueryString": "query_string",
			},
		})
		assert.Error(t, err)
		assert.Equal(t, "QueryName parameter not found", err.Error())

		queryService.EXPECT().PutQueryDefinition(ctx, nil).Times(0)
		queryService.EXPECT().DeleteQueryDefinition(ctx, nil).Times(0)
	})

	t.Run("Returns error when 'QueryString' parameter is not provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		queryService := mock.NewMockQueryService(ctrl)
		h := newHandler(queryService)

		_, _, err := h(ctx, cfn.Event{
			ResourceProperties: map[string]interface{}{
				"QueryName": "query_name",
			},
		})

		assert.Error(t, err)
		assert.Equal(t, "QueryString parameter not found", err.Error())

		queryService.EXPECT().PutQueryDefinition(ctx, nil).Times(0)
		queryService.EXPECT().DeleteQueryDefinition(ctx, nil).Times(0)
	})

	t.Run("Creates the query on the create event", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		queryService := mock.NewMockQueryService(ctrl)
		h := newHandler(queryService)

		queryService.EXPECT().PutQueryDefinition(
			ctx,
			&cloudwatchlogs.PutQueryDefinitionInput{
				Name:        aws.String("test_query_name"),
				QueryString: aws.String("test_query_string"),
			},
		).Return(
			&cloudwatchlogs.PutQueryDefinitionOutput{QueryDefinitionId: aws.String("test_query_definition_id")},
			nil,
		)

		physicalResourceID, data, err := h(ctx, cfn.Event{
			RequestType: cfn.RequestCreate,
			ResourceProperties: map[string]interface{}{
				"QueryString": "test_query_string",
				"QueryName":   "test_query_name",
			},
		})

		assert.NoError(t, err)
		assert.Empty(t, data)
		assert.Equal(t, "test_query_definition_id", physicalResourceID)
	})

	t.Run("Updates the query on the update event", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		queryService := mock.NewMockQueryService(ctrl)
		h := newHandler(queryService)

		queryService.EXPECT().PutQueryDefinition(
			ctx,
			&cloudwatchlogs.PutQueryDefinitionInput{
				Name:              aws.String("test_query_name"),
				QueryString:       aws.String("test_query_string"),
				QueryDefinitionId: aws.String("test_query_definition_id"),
			},
		).Return(
			&cloudwatchlogs.PutQueryDefinitionOutput{QueryDefinitionId: aws.String("test_query_definition_id")},
			nil,
		)

		physicalResourceID, data, err := h(ctx, cfn.Event{
			RequestType: cfn.RequestUpdate,
			ResourceProperties: map[string]interface{}{
				"QueryString": "test_query_string",
				"QueryName":   "test_query_name",
			},
			PhysicalResourceID: "test_query_definition_id",
		})

		assert.NoError(t, err)
		assert.Empty(t, data)
		assert.Equal(t, "test_query_definition_id", physicalResourceID)
	})

	t.Run("Deletes the query on the delete event", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		queryService := mock.NewMockQueryService(ctrl)
		h := newHandler(queryService)

		queryService.EXPECT().DeleteQueryDefinition(
			ctx,
			&cloudwatchlogs.DeleteQueryDefinitionInput{
				QueryDefinitionId: aws.String("test_query_definition_id"),
			},
		).Return(
			&cloudwatchlogs.DeleteQueryDefinitionOutput{},
			nil,
		)

		physicalResourceID, data, err := h(ctx, cfn.Event{
			RequestType: cfn.RequestDelete,
			ResourceProperties: map[string]interface{}{
				"QueryString": "test_query_string",
				"QueryName":   "test_query_name",
			},
			PhysicalResourceID: "test_query_definition_id",
		})

		assert.NoError(t, err)
		assert.Empty(t, data)
		assert.Empty(t, physicalResourceID)
	})
}
