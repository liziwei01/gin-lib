#version 1.0
FROM golang:latest

LABEL maintainer="alssylk@gmail.com"

RUN ["mkdir", "-p", "/home/work/github.com/liziwei01/gin-lib"]

WORKDIR /home/work/github.com/liziwei01/gin-lib

COPY . /home/work/github.com/liziwei01/gin-lib

CMD ["/home/work/github.com/liziwei01/gin-lib/docker_run"] 
