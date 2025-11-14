resource "aws_cloudwatch_log_group" "backend" {
  name              = "/ecs/${var.project_name}-backend"
  retention_in_days = 7

  kms_key_id = null

  tags = {
    Name = "${var.project_name}-backend-logs"
  }
}

resource "aws_cloudwatch_log_group" "frontend" {
  name              = "/ecs/${var.project_name}-frontend"
  retention_in_days = 7

  kms_key_id = null

  tags = {
    Name = "${var.project_name}-frontend-logs"
  }
}

resource "aws_iam_policy" "cloudwatch_logs" {
  name        = "${var.project_name}-cloudwatch-logs"
  description = "Allows ECS tasks to write logs to CloudWatch Logs"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = [
          "${aws_cloudwatch_log_group.backend.arn}:*",
          "${aws_cloudwatch_log_group.frontend.arn}:*"
        ]
      }
    ]
  })

  tags = {
    Name = "${var.project_name}-cloudwatch-logs"
  }
}