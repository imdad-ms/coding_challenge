provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "cf_s3_bucket" {
  bucket = "my-bucket"
}
// as acl = "private" is deprected from terraform that's why have tos et these ones

resource "aws_s3_bucket_ownership_controls" "example" {
  bucket = aws_s3_bucket.cf_s3_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "example" {
  bucket = aws_s3_bucket.cf_s3_bucket.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_acl" "example" {
  depends_on = [
    aws_s3_bucket_ownership_controls.example,
    aws_s3_bucket_public_access_block.example,
  ]

  bucket = aws_s3_bucket.cf_s3_bucket.id
  acl    = "public-read"
}

// create security group
resource "aws_security_group" "allow_https" {
  name        = "allow_https"

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


// crate web server ec2 instance
resource "aws_instance" "web_server" {
    ami           = "ami-0c55b159cbfafe1f0" # need a valid AMI here
    instance_type = "t2.micro"
    security_groups = [aws_security_group.allow_https.name]

    user_data = <<-EOF
                  #!/bin/bash
                  yum update -y
                  yum install -y httpd
                  systemctl start httpd
                  systemctl enable httpd
                  echo '<html><head><title>Hello World</title></head><body><h1>Hello World!</h1></body></html>' > /var/www/html/index.html
                  EOF

    tags = {
      Name = "HTTP Web Server"
    }
}

// need to add code to create cloud front distribution and also make sure that s3 bucket is only accessible form cloud front
