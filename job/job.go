package job

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type app struct {
	s3Event *Event
	sess    *session.Session
	s3sess  *s3.S3
	ctx     context.Context
	log     *log.Logger
	db      *dynamodb.DynamoDB
}

// Do - entry point for lambda function
func Do(ctx context.Context, i input) (err error) {
	xray.Configure(xray.Config{
		LogLevel: "info", // default
	})

	if len(i.Records) < 1 {
		return fmt.Errorf("ERROR: S3 event has no records - cnt: %d", len(i.Records))
	}

	r := i.Records[0]

	// XAmzRequestID that we get here is the ID of S3 EVENT,
	// but we want the log to be prefixed with the actual Lambda Request ID
	requesID := lambdaRequestID(ctx, r.ResponseElements.XAmzRequestID)
	log := log.New(os.Stdout, requesID+" ", log.LstdFlags)

	// we still want to capture the original S3 EVENT ID here
	// and then switch it to Lambda ID so it's used throughout ...
	log.Printf("START Lambda request RuleID=%s Context=%#v XrayTraceID=%s %s\n", r.S3.ConfigurationID, ctx, xrayTraceID(ctx, log), r.ResponseElements.XAmzRequestID)

	sess := BootAws()

	// need to validate identity of our AWS Role Credentials
	whoAmI := sts.New(sess)
	xray.AWS(whoAmI.Client)

	input := &sts.GetCallerIdentityInput{}
	result, err := whoAmI.GetCallerIdentityWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("ERROR: Unable to validate AWS session identify %v", err)
	}
	log.Printf("Aws Identify ID=%s | Arn=%s\n", *result.Account, *result.Arn)

	s3event := Event{
		EventName:       r.EventName,
		EventTime:       r.EventTime,
		S3Bucket:        r.S3.Bucket.Name,
		S3Key:           r.S3.Object.Key,
		ETag:            r.S3.Object.ETag,
		RuleID:          r.S3.ConfigurationID,
		RequestID:       requesID,
		PrincipalID:     r.UserIdentity.PrincipalID,
		SourceIPAddress: r.RequestParameters.SourceIPAddress,
		AwsRegion:       r.AwsRegion,
	}

	// Pass s3 sess down the stack so we reuse it
	s3sess := s3.New(sess)
	xray.AWS(s3sess.Client)

	// Pass dynamoDB session down the stack as well
	db := dynamodb.New(sess)
	xray.AWS(db.Client)

	a := &app{
		s3Event: &s3event,
		sess:    sess,
		s3sess:  s3sess,
		ctx:     ctx,
		log:     log,
		db:      db,
	}

	// annotate xray
	a.init()
	return a.execute()
}
