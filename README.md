# Raptor

This is an example of [AWS GoLang Lambda](https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/) function that receives an [AWS S3 Event](https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html) and integrates with [AWS XRAY](https://aws.amazon.com/xray/).  

For a full explanation you can visit my blog post: [http://www.hashjoin.com/t/aws-golang-lambda-s3-xray-terraform.html](http://www.hashjoin.com/t/aws-golang-lambda-s3-xray-terraform.html)

The Lambda function and it's IAM Policies and Permissions are fully provisioned by [Terraform](https://www.terraform.io/) module (see `deploy/terraform/modules/lambda/` dir).  And terraform `init` and `apply` phases are in turn managed by the `Makefile`:

```
$ make
build/linux                    Build a linux binary ready to be zip'ed for AWS Lambda Deployment
build/release                  Zip linux binary as AWS Deployment archive
clean                          Remove all build artifacts
clean/linux                    Remove linux build artifacts
deploy/production              Deploy zip'ed archive to AWS production account
deploy/testing                 Deploy zip'ed archive to AWS testing account
help                           Display this help message
test/integration               Integration Testing
```

## Prerequisites

Setup `production` and `testing` profiles in `~/.aws/credentials`.  For example:

```
$ cat ~/.aws/credentials
[testing]
aws_access_key_id=********************
aws_secret_access_key=****************************************

[production]
aws_access_key_id=********************
aws_secret_access_key=****************************************
```