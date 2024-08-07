package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession(c Config) (*session.Session, error) {
	sess_config := &aws.Config{
		Region: &c.Region,
	}
	if c.Endpoint != "" {
		sess_config.Endpoint = &c.Endpoint
	}
	if c.AccessKeyID != "" || c.SecretAccessKey != "" {
		sess_config.Credentials = credentials.NewStaticCredentials(
			c.AccessKeyID,
			c.SecretAccessKey,
			"",
		)
	}
	sess, err := session.NewSession(sess_config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sessio: %w", err)
	}
	return sess, nil
}
