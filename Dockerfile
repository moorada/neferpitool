## TO DO

## We specify the base image we need for our
FROM ubuntu:latest 
## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app
## We copy everything in the root directory
## into our /app directory
ADD . /app
## We specify that we now wish to execute 
## any further commands inside our /app

WORKDIR /app/cmd/
