package database

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	QueryFunc func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

func (m *mockDynamoDB) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(input)
	}
	return nil, errors.New("QueryFunc not implemented")
}
