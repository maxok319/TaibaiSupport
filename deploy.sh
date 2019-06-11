#!/bin/bash
# define env var default value.
export rabbitmq_addr=59.110.158.254
export rabbitmq_user=taibai-support
export rabbitmq_passwd=taibai-support
export classroom_region=3
export consul_server_addr=
export consul_client_no=
export host_addr=

#docker-compose down
docker-compose build --build-arg branch_name=master
docker-compose up
