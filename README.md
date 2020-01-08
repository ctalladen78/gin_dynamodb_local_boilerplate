
// create table schema 
```
aws dynamodb create-table \
--endpoint-url http://localhost:8000 \
--key-schema AttributeName=BucketName,KeyType=HASH AttributeName=CreatedAt,KeyType=RANGE \
--attribute-definitions AttributeName=BucketName,AttributeType=S \
AttributeName=CreatedAt,AttributeType=S,AttributeName=Todo,AttributeType=S \
--global-secondary-indexes file://lsi.json
--billing-mode PAY_PER_REQUEST \
--table-name Test
```

// add item with cli
```
aws dynamodb put-item --endpoint-url http://localhost:8000 \
--table-name Test --item file://item.json
```

// http requests

```
GET http://localhost:5000/user?userid=third&todo=action
GET http://localhost:5000/userlist
POST http://localhost:5000/user, --form-data {"todo":"value"}
POST http://localhost:5000/user/edit --form-data {"todo":"newvalue"}
GET http://localhost:5000/query?userid=1BUCKET-username123
```


```
{
    "TableDescription": {
        "TableArn": "arn:aws:dynamodb:ddblocal:000000000000:table/Test",
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
        "ItemCount": 0,
        "CreationDateTime": 1577906088.431
    }
}
```