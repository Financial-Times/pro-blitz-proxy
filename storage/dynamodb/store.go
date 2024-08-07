package dynamodb

import (
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	ID      string              `json:"ID"`
	Headers map[string][]string `json:"Headers"`
	Data    []byte              `json:"Data"`
}

type CacheStore struct {
	conf Config
	db   *dynamodb.DynamoDB
}

func NewCacheStore(sess *session.Session, conf Config) *CacheStore {
	return &CacheStore{
		conf: conf,
		db:   dynamodb.New(sess),
	}
}

func NewCacheStoreWithConfig(conf Config) (*CacheStore, error) {
	sess, err := NewSession(conf)
	if err != nil {
		return nil, err
	}
	return NewCacheStore(sess, conf), nil
}

func NewCacheStoreWithEnvConfig() (*CacheStore, error) {
	conf, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewCacheStoreWithConfig(conf)
}

func (s *CacheStore) Init() error {
	exists, err := s.TableExists(s.conf.TableName)
	if err != nil {
		return fmt.Errorf("could not init dynamodb store: %w", err)
	}
	if !exists {
		err := s.CreateTable()
		if err != nil {
			return fmt.Errorf("could not init dynamodb store: %w", err)
		}
	}
	return nil
}

func (s *CacheStore) TableExists(tableName string) (bool, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	_, err := s.db.DescribeTable(input)
	if err != nil {
		// Check if the error is because the table does not exist
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			return false, nil
		}
		return false, fmt.Errorf("could not check if table exists: %w", err)
	}
	return true, nil
}

func (s *CacheStore) CreateTable() error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(s.conf.TableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
	}
	if s.conf.ProvisionedWriteCapacity >= 0 && s.conf.ProvisionedReadCapactiy >= 0 {
		input.BillingMode = aws.String(dynamodb.BillingModeProvisioned)
		input.ProvisionedThroughput = &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(s.conf.ProvisionedReadCapactiy),
			WriteCapacityUnits: aws.Int64(s.conf.ProvisionedWriteCapacity),
		}
	} else {
		input.BillingMode = aws.String(dynamodb.BillingModePayPerRequest)
	}

	_, err := s.db.CreateTable(input)
	if err != nil {
		return fmt.Errorf("unable to create table: %w", err)
	}
	return nil
}

func (s *CacheStore) Exists(key string) bool {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(s.conf.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(key),
			},
		},
	}

	result, err := s.db.GetItem(input)
	if err != nil {
		slog.Error("error check existence of an item", slog.String("err", err.Error()))
		// @todo - this should return error
		return false
	}

	if result.Item != nil {
		return true
	}
	return false
}

func (s *CacheStore) Save(key string, data []byte, headers map[string][]string) error {
	item, err := dynamodbattribute.MarshalMap(Item{
		ID:      key,
		Headers: headers,
		Data:    data,
	})
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(s.conf.TableName),
		Item:      item,
	})
	if err != nil {
		slog.Error("error save item", slog.String("err", err.Error()))
		return fmt.Errorf("unable to write data: %w", err)
	}

	return nil
}

func (s *CacheStore) Get(key string) ([]byte, map[string][]string, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(s.conf.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(key),
			},
		},
	}

	result, err := s.db.GetItem(input)
	if err != nil {
		slog.Error("failed get item", slog.String("err", err.Error()))
		return nil, nil, fmt.Errorf("failed to get item: %w", err)
	}

	var item Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unamarshal item: %w", err)
	}
	return item.Data, item.Headers, nil
}

func (s *CacheStore) DeleteById(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(s.conf.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	}

	_, err := s.db.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}
