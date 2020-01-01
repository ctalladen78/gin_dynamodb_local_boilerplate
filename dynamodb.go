package main

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DbController struct {
	conn *dynamodb.DynamoDB
}

// uses localhost only
func InitDbConnection(h string) *DbController {
	return &DbController{
		conn: dynamodb.New(session.New(&aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String(h),
		})),
	}
}

// use dynamodbattribute for marshal unmarshal
// we use castTo to define the struct which will contain return value
func (ctrl *DbController) GetItem(key string, table string, castTo interface{}) error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
	}
	res, err := ctrl.conn.GetItem(input)
	if err != nil {
		return err
	}
	err = dynamodbattribute.UnmarshalMap(res.Item, &castTo)
	if err != nil {
		return err
	}
	return nil
}

// ensure item follows attribute value schema
func (ctrl *DbController) PutItem(table string, item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	// TODO recheck if item follows table schema
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}
	_, err = ctrl.conn.PutItem(input)
	if err != nil {
		return err
	}
	return err
}

// pass in an empty attribute value struct which will be populated as a result
func (ctrl *DbController) List(table string, castTo interface{}) error {
	if ctrl.conn == nil {
		return errors.New("db connection error")
	}
	// get all items in table
	scanOutput, err := ctrl.conn.Scan(&dynamodb.ScanInput{
		TableName: aws.String(table),
	})
	if err != nil {
		return err
	}
	err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &castTo)
	if err != nil {
		return err
	}
	return err
}

// itemKey {"key":{S:aws.String("val")}}
func (s *DbController) Update(table string, itemKey interface{}, newItem interface{}) error {
	var err error
	var keyMap map[string]*dynamodb.AttributeValue
	var toUpdate map[string]*dynamodb.AttributeValue
	// http://gist.github.com/doncicuto
	// setup key
	keyMap, err = dynamodbattribute.MarshalMap(itemKey)
	if err != nil {
		return errors.New("Itemkey error")
	}
	// setup new object input
	toUpdate, err = dynamodbattribute.MarshalMap(newItem)
	if err != nil {
		return errors.New("newItem error")
	}
	itemInput := &dynamodb.UpdateItemInput{
		Key:                       keyMap,
		TableName:                 aws.String(table),
		UpdateExpression:          aws.String(""),
		ExpressionAttributeValues: toUpdate,
		ReturnValues:              aws.String("Update Successful"),
	}
	result, err := s.conn.UpdateItem(itemInput)
	if err != nil {
		return errors.New("Update error")
	}

	type UpdatedItem struct {
		// value map[string]*dynamodb.AttributeValue
		value Todo
	}

	// TODO print result
	updatedAttributes := &UpdatedItem{}

	// convert db result into inmemory struct
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedAttributes)

	log.Printf("%s", updatedAttributes.value.Action)
	return nil
}
