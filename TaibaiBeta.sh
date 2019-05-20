#!/bin/bash
# define env var default value.
export rabbitmq_addr=host.docker.internal
export rabbitmq_user=taibai-support
export rabbitmq_passwd=taibai-support
export classroom_region=1

docker-compose build --build-arg branch_name=beta
docker-compose up