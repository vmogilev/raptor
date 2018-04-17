// +build integration

package job

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	testingAWSAccountID = "000000000000"
	testBucket          = "raptor-test-bucket"
)

func TestS3Events(t *testing.T) {
	sess := BootAws()

	if err := isThisTestingAccount(sess); err != nil {
		t.Fatal(err)
	}

	// test goes here

}

func isThisTestingAccount(sess *session.Session) error {
	// need to make sure we are running integration test in testing account
	whoAmI := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}
	result, err := whoAmI.GetCallerIdentity(input)
	if err != nil {
		return fmt.Errorf("ERROR: Unable to validate AWS session identify %v", err)
	}
	if *result.Account != testingAWSAccountID {
		return fmt.Errorf("ERROR: Can't run Integration test in a non-testing account# %s (arn: %s)", *result.Account, *result.Arn)
	}
	return nil
}
