package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vmogilev/raptor/job"
)

func main() {
	lambda.Start(job.Do)
}
