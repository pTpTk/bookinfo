package main

import (
	"time"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

type Rating struct {
	Id	   string `dynamodbav:"id"`
	Rating int    `dynamodbav:"rating"`
}

func DeleteTables(ctx context.Context, db *dynamodb.Client) {
	db.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String("Ratings"),
	})

    NEwaiter := dynamodb.NewTableNotExistsWaiter(db)
	err := NEwaiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String("Ratings")}, 5*time.Minute)
		
	if err != nil {
		panic(err)
	}
}

func CreateRatings(ctx context.Context, db *dynamodb.Client) {
	_, err := db.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName: aws.String("Ratings"),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		fmt.Printf("Couldn't create table Ratings. Here's why: %v\n", err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(db)
		err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String("Ratings")}, 5*time.Minute)
		if err != nil {
			fmt.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
	}
}

func PopulateRatings(ctx context.Context, db *dynamodb.Client) {

    newRatings := []Rating{
		Rating{"1", 5,},
		Rating{"2", 4,},
	}

    for _,r := range newRatings {
		av, err := attributevalue.MarshalMap(r)
		if err != nil {
			fmt.Println(err)
		}

        _, err = db.PutItem(ctx, &dynamodb.PutItemInput{
            TableName: aws.String("Ratings"),
            Item: av,
        })

        if err != nil {
            panic(err)
        }
    }
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	db := dynamodb.NewFromConfig(cfg)

	fmt.Println("Deleting old tables")
	
	DeleteTables(context.TODO(), db)

	fmt.Println("Done")
	fmt.Println("Creating tables")

	CreateRatings(context.TODO(), db)

    fmt.Println("Done")
    fmt.Println("Populating tables")

	PopulateRatings(context.TODO(), db)

    fmt.Println("Done")
}