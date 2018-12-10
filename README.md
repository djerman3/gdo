# gdo - garage door opener

I'm learning AWS Lambda + API Gateway + authentication with this 
homebrew garage door gateway.

The API is hosted in my AWS account, with a React frontend in a static S3 bucket.

My raspberry pi with its piFace cape will post the destination IP to the API as my home IP changes.

The lambda functions will authenticate/authorize users and call server functions to get door state and usage statistics.
