# code-challenge-json-transformer

DevX RUN SEAL JSON Transformer Coding Challenge

## Requirements

- need to install terraform
- then need to init terraform.


### Architecture and Scaleability

for  a secure architecture and scaleable static web application

1. create AWS S3 login to host static website contents
2.  then implement cloudfront to server static contents via CDN to reduce latency and increase speed by catching
3. For AWS certificate management (ACM) can be used for automated renewal
4. Route 53 for DNS management and directing to cloudFront

5. AWS IAM roles need to make sure all services are properly secured
6. Use security group to control incoming traffic (e.g port 443)