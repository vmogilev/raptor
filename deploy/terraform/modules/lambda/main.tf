provider "aws" {
}

resource "aws_iam_role" "raptor-role" {
  name = "raptor-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "raptor-s3-policy" {
    name        = "raptor-s3-policy"
    description = "raptor-s3-policy"
    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "0",
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket",
                "s3:GetBucketLocation"
            ],
            "Resource": "arn:aws:s3:::${var.s3_bucket_name}"
        },
        {
            "Sid": "1",
            "Effect": "Allow",
            "Action": "s3:*",
            "Resource": "arn:aws:s3:::${var.s3_bucket_name}/*"
        }
    ]
}
EOF
}

resource "aws_iam_policy" "raptor-xray-policy" {
    name        = "raptor-xray-policy"
    description = "raptor-xray-policy"
    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": {
        "Effect": "Allow",
        "Action": [
            "xray:PutTraceSegments",
            "xray:PutTelemetryRecords"
        ],
        "Resource": [
            "*"
        ]
    }
}
EOF
}

resource "aws_iam_policy" "raptor-dynamodb-tables-policy" {
    name        = "raptor-dynamodb-policy"
    description = "grants access to all tables prefixed by raptor_*"
    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "dynamodb:BatchGetItem",
                "dynamodb:BatchWriteItem",
                "dynamodb:DeleteItem",
                "dynamodb:GetItem",
                "dynamodb:PutItem",
                "dynamodb:Query",
                "dynamodb:UpdateItem"
            ],
            "Resource": [
                "arn:aws:dynamodb:*:*:table/raptor_*"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "raptor-role-policy-attach-1" {
  role = "${aws_iam_role.raptor-role.name}"
  policy_arn = "${aws_iam_policy.raptor-s3-policy.arn}"
}

resource "aws_iam_role_policy_attachment" "raptor-role-policy-attach-2" {
  role = "${aws_iam_role.raptor-role.name}"
  policy_arn = "arn:aws:iam::aws:policy/AWSLambdaExecute"
}

resource "aws_iam_role_policy_attachment" "raptor-role-policy-attach-3" {
  role = "${aws_iam_role.raptor-role.name}"
  policy_arn = "${aws_iam_policy.raptor-xray-policy.arn}"
}

resource "aws_iam_role_policy_attachment" "raptor-role-policy-attach-4" {
  role = "${aws_iam_role.raptor-role.name}"
  policy_arn = "${aws_iam_policy.raptor-dynamodb-tables-policy.arn}"
}

resource "aws_lambda_function" "raptor" {
  filename         = "../../../build/raptor.zip"
  function_name    = "raptor"
  role             = "${aws_iam_role.raptor-role.arn}"
  handler          = "raptor"
  source_code_hash = "${base64sha256(file("../../../build/raptor.zip"))}"
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 30
  reserved_concurrent_executions = 50
  publish          = true

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      BUCKET_NAME = "${var.s3_bucket_name}"
    }
  }
}

resource "aws_lambda_permission" "raptor-bucket" {
  statement_id  = "1"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.raptor.arn}"
  principal     = "s3.amazonaws.com"
  source_arn    = "arn:aws:s3:::${var.s3_bucket_name}"
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = "${var.s3_bucket_name}"

  lambda_function {
    id                  = "one"
    lambda_function_arn = "${aws_lambda_function.raptor.arn}"
    events              = ["s3:ObjectCreated:*","s3:ObjectRemoved:*"]
    filter_suffix       = "/ONE"
  }

  lambda_function {
    id                  = "two"
    lambda_function_arn = "${aws_lambda_function.raptor.arn}"
    events              = ["s3:ObjectCreated:*","s3:ObjectRemoved:*"]
    filter_suffix       = "/TWO"
  }

  lambda_function {
    id                  = "three"
    lambda_function_arn = "${aws_lambda_function.raptor.arn}"
    events              = ["s3:ObjectCreated:*","s3:ObjectRemoved:*"]
    filter_suffix       = "/THREE"
  }
}