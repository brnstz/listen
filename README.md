listen
======

NYC concert listings with sample audio: http://listen.brnstz.com/

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

