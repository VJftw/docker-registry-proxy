locals {
  namespace           = "test"
  stage               = "e2e"
  name                = "dockerregistryproxy"
  tags                = {}
  region              = "eu-west-2"
  kubernetes_version  = ""
  desired_nodes       = 1
  node_disk_size      = 20
  node_instance_types = ["t3.medium"]
}

module "vpc" {
  source     = "git::https://github.com/cloudposse/terraform-aws-vpc.git?ref=tags/0.8.1"
  namespace  = local.namespace
  stage      = local.stage
  name       = local.name
  cidr_block = "172.16.0.0/16"
  tags       = local.tags
}

data "aws_availability_zones" "available" {
  state = "available"
}

module "subnets" {
  source               = "git::https://github.com/cloudposse/terraform-aws-dynamic-subnets.git?ref=tags/0.19.0"
  availability_zones   = data.aws_availability_zones.available.names
  namespace            = local.namespace
  stage                = local.stage
  name                 = local.name
  vpc_id               = module.vpc.vpc_id
  igw_id               = module.vpc.igw_id
  cidr_block           = module.vpc.vpc_cidr_block
  nat_gateway_enabled  = false
  nat_instance_enabled = false
  tags                 = {
    "kubernetes.io/cluster/${local.namespace}-${local.stage}-${local.name}-cluster" = "shared"
  }
  public_subnets_additional_tags = {
    "kubernetes.io/role/elb" = "1"
  }
}

module "eks_cluster" {
  source             = "git::https://github.com/cloudposse/terraform-aws-eks-cluster.git?ref=0.22.0"
  namespace          = local.namespace
  stage              = local.stage
  name               = local.name
  tags               = local.tags
  region             = local.region
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.subnets.public_subnet_ids
  kubernetes_version = local.kubernetes_version
}

# Ensure ordering of resource creation to eliminate the race conditions when applying the Kubernetes Auth ConfigMap.
# Do not create Node Group before the EKS cluster is created and the `aws-auth` Kubernetes ConfigMap is applied.
# Otherwise, EKS will create the ConfigMap first and add the managed node role ARNs to it,
# and the kubernetes provider will throw an error that the ConfigMap already exists (because it can't update the map, only create it).
# If we create the ConfigMap first (to add additional roles/users/accounts), EKS will just update it by adding the managed node role ARNs.
data "null_data_source" "wait_for_cluster_and_kubernetes_configmap" {
  inputs = {
    cluster_name             = module.eks_cluster.eks_cluster_id
    kubernetes_config_map_id = module.eks_cluster.kubernetes_config_map_id
  }
}

module "eks_node_group" {
  source         = "git::https://github.com/cloudposse/terraform-aws-eks-node-group.git?ref=tags/0.4.0"
  namespace      = local.namespace
  stage          = local.stage
  name           = local.name
  subnet_ids     = module.subnets.public_subnet_ids
  cluster_name   = data.null_data_source.wait_for_cluster_and_kubernetes_configmap.outputs["cluster_name"]
  instance_types = local.node_instance_types
  desired_size   = local.desired_nodes
  min_size       = local.desired_nodes
  max_size       = local.desired_nodes
  disk_size      = local.node_disk_size
}