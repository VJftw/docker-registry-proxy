data "aws_ecr_authorization_token" "test" {
  registry_id = aws_ecr_repository.test.registry_id
}

resource "aws_ecr_repository" "test" {
  name                 = "test-web"
  image_tag_mutability = "MUTABLE"
}

locals {
  test_image = "nginxdemos/hello:plain-text"
}

resource "null_resource" "test_image_push" {
  triggers = {
    registry_id = aws_ecr_repository.test.registry_id
  }

  provisioner "local-exec" {
    environment = {
      ECR_PASSWORD = data.aws_ecr_authorization_token.test.password
    }

    command = <<EOC
docker pull ${local.test_image}
docker tag ${local.test_image} ${aws_ecr_repository.test.repository_url}:latest
echo "$ECR_PASSWORD" | docker login \
  --username ${data.aws_ecr_authorization_token.test.user_name} \
  --password-stdin \
  ${data.aws_ecr_authorization_token.test.proxy_endpoint}
docker push ${aws_ecr_repository.test.repository_url}:latest
docker logout ${data.aws_ecr_authorization_token.test.proxy_endpoint}
docker rmi ${local.test_image} ${aws_ecr_repository.test.repository_url}:latest
EOC
  
  }
}