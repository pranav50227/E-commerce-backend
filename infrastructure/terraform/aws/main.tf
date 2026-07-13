provider "aws" {
  region = "us-east-1"
}

# 1. VPC Configuration
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  tags = {
    Name = "ecommerce-vpc"
  }
}

# 2. Security Groups
resource "aws_security_group" "web_sg" {
  name        = "ecommerce-web-sg"
  description = "Allow web traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# 3. RDS Postgres Instance
resource "aws_db_instance" "postgres" {
  allocated_storage    = 20
  engine               = "postgres"
  engine_version       = "16"
  instance_class       = "db.t4g.micro"
  db_name              = "ecommerce"
  username             = "postgres"
  password             = "supersecurepassword123"
  parameter_group_name = "default.postgres16"
  skip_final_snapshot  = true
}

# 4. EC2 Application Server
resource "aws_instance" "app_server" {
  ami           = "ami-0c7217cdde317cfec" # Ubuntu 22.04 LTS
  instance_type = "t3.medium"
  security_groups = [
    aws_security_group.web_sg.name
  ]

  tags = {
    Name = "ecommerce-app-server"
  }
}
