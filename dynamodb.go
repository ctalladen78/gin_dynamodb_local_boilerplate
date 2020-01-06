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
func (ctrl *DbController) GetItem(id string, todo string, table string) (interface{}, error) {
	// https://github.com/ace-teknologi/memzy
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go
	var pkey = map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
		"todo": {
			S: aws.String(todo),
		},
	}
	// building pkey for search query
	// keys = &TodoObject{"id":id,"todo":todo}
	// var pkey = map[string]*dynamodb.AttributeValue{}
	// option 1
	// for k, v := range keys {
	// av, err := dynamodbattribute.Marshal(v)
	// 	pkey[k] = av
	// }
	// option 2
	// pkey, err := dynamodbattribute.MarshalMap(keys)
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
	// https://stackoverflow.com/questions/38151687/dynamodb-adding-non-key-attributes/56177142
	avMap, err := dynamodbattribute.MarshalMap(todo) // conver todo item to av map
	log.Printf("AV Map %v", avMap)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      avMap,
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

// using dynamodb.Scan as opposed to dynamodb.Query
func (ctrl *DbController) ScanFilter(table string) {

}

// itemKey {"partitionKey":{S:aws.String("val")},{"partitionKey":{S:aws.String("val")}}
func (ctrl *DbController) Update(table string, oldItem *TodoObject, newItemValue string) (interface{}, error) {
	var err error
	// var keyMapAV2 map[string]*dynamodb.AttributeValue
	// var toUpdate map[string]*dynamodb.AttributeValue
	// http://gist.github.com/doncicuto
	// setup key
	// keyMapAV, err := dynamodbattribute.MarshalMap(&TodoObject{
	// 	CreatedAt: "2020-01-06T18:38:01+07:00",
	// 	Todo:      "texasasasas", // oldItem.Todo
	// 	Id:        "BucketName",  // oldItem.Id
	// })
	keyMapAV := map[string]*dynamodb.AttributeValue{
		"id":   {S: aws.String("texasasasas")},
		"todo": {S: aws.String("BucketName")},
	}
	log.Printf("AV Map %v", keyMapAV)
	if err != nil {
		return nil, errors.New("itemkey error")
	}
	if err != nil {
		return nil, errors.New("newItem error")
	}
	// https://aws.amazon.com/blogs/developer/introducing-amazon-dynamodb-expression-builder-in-the-aws-sdk-for-go/
	// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.UpdateExpressions.html
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go#L33
	itemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(table),
		Key:       keyMapAV, // match key attributes per table definition
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {S: aws.String(newItemValue)},
		}, // set new value
		ExpressionAttributeNames: map[string]*string{"#T": aws.String("createdat")}, // attribute must not be part of the key
		// ConditionExpression:	"attribute_exists("#T"),
		// https://gist.github.com/doncicuto/d623ec0e74bf6ea0db7c364d88507393#file-main-go-L63
		ReturnValues:     aws.String("ALL_NEW"),     // enum of ReturnValue class
		UpdateExpression: aws.String("set #T = :t"), // SET,REMOVE the attribute to update

	}
	result, err := ctrl.conn.UpdateItem(itemInput)
	if err != nil {
		return nil, err
	}

	// type UpdatedItem struct {
	// 	newItem TodoObject
	// }

	// TODO print resulting updated attributes
	// updatedAttributes := &UpdatedItem{}

	// convert db result into inmemory struct
	// err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedAttributes)

	log.Printf("%s", result.Attributes)
	// return updatedAttributes.newItem, nil
	return nil, nil
}
