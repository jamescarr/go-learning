# Example Interacting with Graylog
This is an example of using sling to interact with graylog, implementing GET, POST and DELETE 
operations.

I actually started this to integrate with graylog but came to realize that graylog uses swagger
and you can just generate service stubs using [swagger codegen](https://github.com/swagger-api/swagger-codegen#to-generate-a-sample-client-library). 

## Running It
- `docker-compose up -d`
- `go test`
