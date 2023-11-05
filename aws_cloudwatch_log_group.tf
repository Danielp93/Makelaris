resource "aws_cloudwatch_log_group" "cloudwatch" {
  name              = "/aws/lambda/${var.project_name}_lambda_fn"
  retention_in_days = 3
}

resource "aws_iam_policy" "lambda_cloudwatch_policy" {
  name        = "${var.project_name}_lambda_cloudwatch_policy"
  path        = "/"
  description = "IAM policy for cloudwatch logging from a lambda"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      { 
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Effect = "Allow",
        Resource: "${aws_cloudwatch_log_group.cloudwatch.arn}"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_cloudwatch_policy_attachment" {
  role       = aws_iam_role.lambda_iam_role.name
  policy_arn = aws_iam_policy.lambda_cloudwatch_policy.arn
}