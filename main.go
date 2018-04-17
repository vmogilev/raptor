package main

import (
	"github.com/InVisionApp/hound-agent/job"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(job.Do)
}
