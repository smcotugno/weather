# weather-app

This Microservices Weather application is composed of three containers as follows:

  weather-proxy: nginx based web server that provides the access to the user via port 6060 to access the weather-ui interface
  weather-ui: nginx based web server providing the User Interface to the weather application
  weather-api:a golang based web server that serves the REST API to the weather-ui microservice.
              The REST API implementation gets historical data from the Weather Channel via the REST API provided at www.wunderground.com

There are two different deployment models provided.
  - Docker deployment using Docker Compose (see weather/container/docker-compose.yaml)
  - Kubernetes deployment - see the weather/k8s for the yaml files for kubernettes deployment
