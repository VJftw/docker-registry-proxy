provider "aws" {
  region  = "eu-west-2"
  version = "~> 2.67"
}

data "aws_route53_zone" "main" {
  name = "vjpatel.me."
}