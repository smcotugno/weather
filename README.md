Weather Application 

This Microservices Weather application is composed of three microservices implemented with Docker containers.  
The microservices that compose this weather application are as follows:

  - weather-proxy: nginx based web microservice that provides user access to the weather application.

  - weather-ui: nginx based web microservice that provides the User Interface for the weather application

  - weather-api:a golang based web microservice  that serves the REST API to the weather-ui microservice. The REST API implementation gets historical data from the Weather Channel via the REST API provided at www.wunderground.com

There are two different deployment models provided.
  - Docker deployment using Docker Compose (see weather/container/docker-compose.yaml).
  - Kubernetes deployment - see the weather/k8s for the yaml files for kubernettes deployment
