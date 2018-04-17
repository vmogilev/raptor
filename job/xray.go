package job

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-xray-sdk-go/header"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func (a *app) init() {
	// start INIT segment so we can annotate with values
	topCTX := a.ctx
	ctxINIT, initSEG := xray.BeginSubsegment(topCTX, "INIT")
	a.ctx = ctxINIT

	// annotate xray with our values so it's easy to filter:
	//	annotation.RuleID = "one"
	//	annotation.RequestID = "72861f13-3ec2-11e8-b266-3fbfa0ef4b01"
	a.annotate("RuleID", a.s3Event.RuleID)
	a.annotate("RequestID", a.s3Event.RequestID)

	initSEG.Close(nil)

	// switch back to the original ctx
	a.ctx = topCTX
}

func (a *app) annotate(key string, val interface{}) {
	// this allows searching xray with a filter:
	//	annotation.key = "val"
	err := xray.AddAnnotation(a.ctx, key, val)
	if err != nil {
		a.log.Printf("ERROR: failed to annotate xray with %s/%v: %v\n", key, val, err)
	}
}

// lambdaRequestID - parse the actual Lambda Request ID from context
func lambdaRequestID(ctx context.Context, id string) string {
	lc, ok := lambdacontext.FromContext(ctx)
	if ok {
		return lc.AwsRequestID
	}
	return id
}

// xrayTraceID - parse Xray Trace ID from context
func xrayTraceID(ctx context.Context, log *log.Logger) string {
	if traceHeaderValue := ctx.Value(xray.LambdaTraceHeaderKey); traceHeaderValue != nil {
		traceHeader := traceHeaderValue.(string)
		return header.FromString(traceHeader).TraceID
	}
	return ""
}
