terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
    profile = "default"
    region = "eu-central-1"
}

variable "project_name" {
    description = "Name of the current project"
    default = "Housing"
}