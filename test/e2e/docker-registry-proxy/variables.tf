variable "base_domain" {
    type = string
    description = "the base domain to expose ingresses over"
}

variable "docker_registry_proxy_image" {
    type = string
    description = "the docker registry proxy image to use"
    default = "vjftw/docker-registry-proxy:latest"
}

variable "docker_registry_proxy_authentication_verifier" {
    type = string
    description = "The authentication verification configuration to use. (empty grants anonymous access)"
    default = ""
}

variable "docker_registry_proxy_gcpverifier_project_ids" {
    type = string
    description = "The GCP Project IDs to accept"
    default = ""
}

variable "docker_registry_proxy_awsverifier_account_ids" {
    type = string
    description = "The AWS Account IDs to accep"
    default = ""
}
