package dynamodb

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Endpoint                 string
	AccessKeyID              string
	SecretAccessKey          string
	Region                   string
	TableName                string
	ProvisionedReadCapactiy  int64
	ProvisionedWriteCapacity int64
}

func NewConfigFromEnv() (Config, error) {
	c := Config{
		Region:    "eu-west-1",
		TableName: "blitz-proxy",
	}
	if v := os.Getenv("AWS_REGION"); v != "" {
		c.Region = v
	}

	if v := os.Getenv("AWS_ENDPOINT"); v != "" {
		c.Endpoint = v
	}

	if v := os.Getenv("AWS_ACCESS_KEY_ID"); v != "" {
		c.AccessKeyID = v
	}

	if v := os.Getenv("AWS_SECRET_ACCESS_KEY"); v != "" {
		c.SecretAccessKey = v
	}

	v := os.Getenv("AWS_DYNAMODB_TABLENAME")
	if v == "" {
		return c, fmt.Errorf("AWS_DYNAMODB_TABLENAME is not set")
	}
	c.TableName = v

	if v := os.Getenv("AWS_DYNAMODB_PROVISIONED_READ_CAPACITY"); v != "" {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return c, fmt.Errorf("invalid config value for provisioned read capacity: '%v'", v)
		}
		c.ProvisionedReadCapactiy = val
	}

	if v := os.Getenv("AWS_DYNAMODB_PROVISIONED_WRITE_CAPACITY"); v != "" {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return c, fmt.Errorf("invalid config value for provisioned write capacity: '%v'", v)
		}
		c.ProvisionedWriteCapacity = val
	}
	return c, nil
}
