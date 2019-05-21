#!/bin/bash
# define env var default value.
export rabbitmq_addr=host.docker.internal
export rabbitmq_user=taibai-support
export rabbitmq_passwd=taibai-support
export classroom_region=3

docker-compose down
docker-compose build --build-arg branch_name=master
docker-compose up
