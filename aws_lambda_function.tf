resource "aws_lambda_function" "lambda_fn" {
  # If the file is not in the current working directory you will need to include a
  # path.module in the filename.
  filename      = "./go/queryWebpage/queryWebpage.zip"
  function_name = "${var.project_name}_queryWebpage"
  role          = aws_iam_role.lambda_iam_role.arn
  handler       = "queryWebpage"

  source_code_hash = filebase64sha256("./go/queryWebpage/queryWebpage.zip")

  runtime = "go1.x"

  depends_on = [
    aws_iam_role_policy_attachment.lambda_dynamodb_policy_attachment,
    aws_iam_role_policy_attachment.lambda_cloudwatch_policy_attachment,
    aws_cloudwatch_log_group.cloudwatch,
    aws_dynamodb_table.entries,
    aws_dynamodb_table.forms
  ]
}

resource "aws_iam_role" "lambda_iam_role" {
  name = "${var.project_name}_lambda_iam_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}