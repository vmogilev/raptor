module "lambda" {
  source = "../modules/lambda"
  s3_bucket_name = "raptor-prod-bucket"
}