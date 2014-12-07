listen
======

Data NYC concert listings with sample audio. 

Hosted at: http://listen.brnstz.com/

## Related code

Go API client for ohmyrockness.com: https://github.com/brnstz/ohmy

## Quickstart
```bash

# Clone, get dependencies and build the binary
git clone git@github.com:brnstz/listen.git
cd listen
go get ./...
go install github.com/brnstz/listen

# Set environment and run a Rabbit MQ server
docker run -d -p 5672:5672 -p 15672:15672 dockerfile/rabbitmq
export AMQP_URL="amqp://guest:guest@localhost:5672"
export AWS_ACCESS_KEY_ID='your access key'
export AWS_SECRET_ACCESS_KEY='your secret key'
export AWS_DEFAULT_REGION='your aws region'
export AWS_S3_BUCKET='your bucket'
export LISTEN_STATIC_DIR=`pwd`/html

# Run it
listen

# Look at it: http://localhost:8084/
```

