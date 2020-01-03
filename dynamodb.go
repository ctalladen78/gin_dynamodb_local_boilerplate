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

// get item by id, todo as per table schema
func (ctrl *DbController) GetItem(v string, t string, table string) (interface{}, error) {
	// https://github.com/ace-teknologi/memzy
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go
	var pkey = map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(v),
		},
		"todo": {
			S: aws.String(t),
		},
	}
	// av, err := dynamodbattribute.MarshalMap(pkey)
	// if err != nil {
	// 	return nil, errors.New("INVALID key")
	// }
	input := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       pkey,
	}
	log.Println("GET ITEM input %s", input)
	res, err := ctrl.conn.GetItem(input)
	log.Println("GET ITEM output %s", res)
	if err != nil {
		return nil, err
	}
	var castTo *TodoObject
	err = dynamodbattribute.UnmarshalMap(res.Item, &castTo)
	if err != nil {
		return nil, err
	}
	return castTo, nil
}

// ensure item follows attribute value schema
func (ctrl *DbController) PutItem(table string, todo interface{}) error {
	av, err := dynamodbattribute.MarshalMap(todo)
	// TODO recheck if item follows table schema
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}
	log.Printf("Put Input %v", input)
	o, err := ctrl.conn.PutItem(input)
	if err != nil {
		return err
	}
	log.Printf("Put output %v", o)
	return err
}

// pass in an empty attribute value struct which will be populated as a result
func (ctrl *DbController) List(table string) (interface{}, error) {
	if ctrl.conn == nil {
		return nil, errors.New("db connection error")
	}
	// get all items in table
	scanOutput, err := ctrl.conn.Scan(&dynamodb.ScanInput{
		TableName: aws.String(table),
	})
	if err != nil {
		return nil, err
	}
	var castTo []*TodoObject
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go
	err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &castTo)
	if err != nil {
		return nil, err
	}
	return castTo, nil
}

// itemKey {"partitionKey":{S:aws.String("val")},{"partitionKey":{S:aws.String("val")}}
func (ctrl *DbController) Update(table string, itemKey *TodoObject, newValue string) (interface{}, error) {
	var err error
	var keyMap map[string]*dynamodb.AttributeValue
	var toUpdate *dynamodb.AttributeValue
	// http://gist.github.com/doncicuto
	// setup key
	keyMap, err = dynamodbattribute.MarshalMap(itemKey)
	if err != nil {
		return nil, errors.New("Itemkey error")
	}
	// setup new object input
	toUpdate, err = dynamodbattribute.Marshal(newValue)
	if err != nil {
		return nil, errors.New("newItem error")
	}
	itemInput := &dynamodb.UpdateItemInput{
		Key:              keyMap,
		TableName:        aws.String(table),
		UpdateExpression: aws.String(""),
		// ExpressionAttributeValues: toUpdate,
		ReturnValues: aws.String("Update Successful"),
	}
	result, err := ctrl.conn.UpdateItem(itemInput)
	if err != nil {
		return nil, errors.New("Update error")
	}

	type UpdatedItem struct {
		// value map[string]*dynamodb.AttributeValue
		value TodoObject
	}

	// TODO print result
	updatedAttributes := &UpdatedItem{}

	// convert db result into inmemory struct
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedAttributes)

	log.Printf("%s", updatedAttributes.value.Todo)
	return updatedAttributes.value, nil
}
