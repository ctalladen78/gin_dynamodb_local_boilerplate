
// create table schema 
```
aws dynamodb create-table \
--endpoint-url http://localhost:8000 \
--key-schema AttributeName=id,KeyType=HASH AttributeName=todo,KeyType=RANGE \
--attribute-definitions AttributeName=id,AttributeType=S \
AttributeName=todo,AttributeType=S \
--billing-mode PAY_PER_REQUEST \
--table-name Test
```

// add item
```
aws dynamodb put-item --endpoint-url http://localhost:8000 \
--table-name Test --item file://item.json
```



```
{
    "TableDescription": {
        "TableArn": "arn:aws:dynamodb:ddblocal:000000000000:table/Test",
        "AttributeDefinitions": [
            {
                "AttributeName": "id",
                "AttributeType": "S"
            },
            {
                "AttributeName": "todo",
                "AttributeType": "S"
            }
        ],
        "ProvisionedThroughput": {
            "NumberOfDecreasesToday": 0,
            "WriteCapacityUnits": 0,
            "LastIncreaseDateTime": 0.0,
            "ReadCapacityUnits": 0,
            "LastDecreaseDateTime": 0.0
        },
        "TableSizeBytes": 0,
        "TableName": "Test",
        "BillingModeSummary": {
            "LastUpdateToPayPerRequestDateTime": 1653090.799,
            "BillingMode": "PAY_PER_REQUEST"
        },
        "TableStatus": "ACTIVE",
        "KeySchema": [
            {
                "KeyType": "HASH",
                "AttributeName": "id"
            },
            {
                "KeyType": "RANGE",
                "AttributeName": "todo"
            }
        ],
        "ItemCount": 0,
        "CreationDateTime": 1577906088.431
    }
}
```