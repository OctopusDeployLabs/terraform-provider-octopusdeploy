---
page_title: "2. Setup"
subcategory: "Guides"
---

# 2. Setup

To use the Terraform provider, you may want to test it if you don't have a dev environment of Octopus Deploy.

To do this, let's set up Octopus Deploy instead of a Docker Container.

## Prerequisites
To follow along, you should have the following:
1. Docker installed on the OS that you're using
2. Familiarity with Docker

## The Dockerfile

To create an Octopus Deploy environment, the easiest way is with Docker Compose.

The docker-compose code below will create an Octopus Deploy server and a local SQL server.

You can access it once the container is up and running by going to `http://localhost:8080`

*Please note this is for development purposes and should not be used in production*

```
version: '3'
services:
  octopus:
    ports:
      - "8080:8080"
      - "10943:10943"
    environment:
      ADMIN_USERNAME: admin
      ADMIN_EMAIL: test@gmail.com
      ADMIN_PASSWORD: Password01!
      ACCEPT_EULA: Y
      DB_CONNECTION_STRING: Server=mssql,1433;Database=Octopus;User Id=SA;Password=Password01!;ConnectRetryCount=6
      CONNSTRING: Server=mssql,1433;Database=Octopus;User Id=SA;Password=Password01!;ConnectRetryCount=6
      MASTER_KEY: 6EdU6IWsCtMEwk0kPKflQQ==
    image: octopusdeploy/octopusdeploy:latest
    labels:
      autoheal: true
    depends_on:
      - mssql
  mssql:
    environment:
      ACCEPT_EULA: Y
      SA_PASSWORD: Password01!
      MSSQL_PID: Express
    image: mcr.microsoft.com/mssql/server:2017-latest-ubuntu
  autoheal:
    image: willfarrell/autoheal:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```