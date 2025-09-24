package database

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Database struct {
	*dynamodb.Client
}

func NewDatabase(cfg aws.Config, options ...func(*dynamodb.Options)) (database *Database) {
	return &Database{dynamodb.NewFromConfig(cfg, options...)}
}

// unmarshalDynamodbMapSlice unmarshals the items in a DynamoDB QueryOutput into a slice of structs.
// It takes a queryOutput representing the result of a DynamoDB query, and a pointer to a slice
// of structs (slicePtr) where the unmarshalled items will be stored.
//
// Example usage:
//
//	var skillsTools []models.SkillsTools
//	err := unmarshalDynamodbMapSlice(queryOutput, &skillsTools)
//	if err != nil {
//	    log.Printf("error in unmarshalling: %v", err)
//	}
func unmarshalDynamodbMapSlice(queryOutput *dynamodb.QueryOutput, slicePtr interface{}) error {
	items := queryOutput.Items
	sliceValue := reflect.ValueOf(slicePtr).Elem()
	elementType := sliceValue.Type().Elem()

	for _, item := range items {
		newItem := reflect.New(elementType).Interface()
		if err := attributevalue.UnmarshalMap(item, newItem); err != nil {
			return fmt.Errorf("error in deserializing dynamodb item: %v", elementType.Name())
		}
		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(newItem).Elem()))
	}
	return nil
}

// GetItem retrieves an item from DynamoDB based on the provided personalWebsiteType
// and sortValue. It takes an initialized DynamoDB client (svc), the personalWebsiteType
// and sortValue to uniquely identify the item, and a pointer to the struct (itemPtr)
// where the retrieved item will be unmarshalled.
//
// Example usage:
//
//	var item models.Item
//	err := GetItem(ctx, svc, "tableName", "type", "sortValue", &item)
//	if err != nil {
//	    log.Printf("error retrieving item: %v", err)
//	}
func GetItem(ctx context.Context, svc *dynamodb.Client, tableName string, personalWebsiteType string, sortValue string, itemPtr interface{}) (err error) {
	inputGet := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"personalWebsiteType": &types.AttributeValueMemberS{Value: personalWebsiteType},
			"sortValue":           &types.AttributeValueMemberS{Value: sortValue},
		},
		TableName: aws.String(tableName),
	}
	result, err := svc.GetItem(ctx, inputGet)
	if err != nil {
		log.Printf("error in DynamoDB GetItem func: %v", err)
		return err
	}
	err = attributevalue.UnmarshalMap(result.Item, itemPtr)
	if err != nil {
		log.Printf("error in DynamoDB UnmarshalMap func: %v", err)
		return err
	}
	return nil
}

func DeleteItem(ctx context.Context, svc *dynamodb.Client, tableName string, personalWebsiteType string, sortValue string) (err error) {
	key := map[string]types.AttributeValue{
		"personalWebsiteType": &types.AttributeValueMemberS{Value: personalWebsiteType},
		"sortValue":           &types.AttributeValueMemberS{Value: sortValue},
	}
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}
	_, err = svc.DeleteItem(ctx, input)
	if err != nil {
		log.Printf("error in DynamoDB DeleteItem func: %v", err)
		return err
	}
	return nil
}
