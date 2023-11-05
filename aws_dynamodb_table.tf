resource "aws_dynamodb_table" "entries" {
  name           = "${var.project_name}_entries"
  hash_key       = "dataRecordId"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  attribute {
    name = "dataRecordId"
    type = "S"
  }
}

resource "aws_dynamodb_table" "forms" {
  name           = "${var.project_name}_forms"
  hash_key       = "FormId"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  attribute {
    name = "FormId"
    type = "S"
  }
}

resource "aws_iam_policy" "lambda_dynamodb_policy" {
  name        = "${var.project_name}_lambda_dynamodb_policy"
  path        = "/"
  description = "IAM policy for dynamoDB access from a lambda"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "dynamodb:ConditionCheckItem",
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query",
          "dynamodb:Scan",
        ]
        Effect = "Allow",
        Resource : [
            "${aws_dynamodb_table.forms.arn}",
            "${aws_dynamodb_table.entries.arn}"
        ]
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_dynamodb_policy_attachment" {
  role       = aws_iam_role.lambda_iam_role.name
  policy_arn = aws_iam_policy.lambda_dynamodb_policy.arn
}