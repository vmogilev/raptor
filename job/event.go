package job

import "time"

// Event - S3 event that we pass around lambda function
type Event struct {
	EventName       string    `json:"eventName"`
	EventTime       time.Time `json:"eventTime"`
	S3Bucket        string    `json:"s3Bucket"`
	S3Key           string    `json:"s3Key"`
	ETag            string    `json:"eTag"`
	RuleID          string    `json:"ruleID"`
	RequestID       string    `json:"requestID"`
	PrincipalID     string    `json:"principalID"`
	SourceIPAddress string    `json:"sourceIPAddress"`
	AwsRegion       string    `json:"awsRegion"`
}
