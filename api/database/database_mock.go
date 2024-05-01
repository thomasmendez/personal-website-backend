package database

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	QueryFunc  func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	PutFunc    func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetFunc    func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	UpdateFunc func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
}

func (m *mockDynamoDB) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(input)
	}
	return nil, errors.New("QueryFunc not implemented")
}
func (m *mockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	if m.PutFunc != nil {
		return m.PutFunc(input)
	}
	return nil, errors.New("PutItem not implemented")
}
func (m *mockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	if m.GetFunc != nil {
		return m.GetFunc(input)
	}
	return nil, errors.New("GetItem not implemented")
}
func (m *mockDynamoDB) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(input)
	}
	return nil, errors.New("UpdateItem not implemented")
}
