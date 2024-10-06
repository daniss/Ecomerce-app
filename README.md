# Microservices-based E-commerce Application with Monitoring

## Project Overview

This is a microservices-based e-commerce platform built using **Go**, **Kubernetes**, **PostgreSQL**, and **Docker**. The project includes key features such as user management, product management, and order processing, with each microservice communicating via APIs. 

Additionally, **Prometheus** and **Grafana** are integrated for monitoring system performance and generating alerts when thresholds are exceeded.

## Features

- **Microservices**:
  - **User Service**: Handles user registration and authentication.
  - **Product Service**: Manages product inventory and details.
  - **Order Service**: Handles order placement and interacts with the product service to update inventory.
  
- **Monitoring and Alerts**:
  - **Prometheus** collects metrics on resource usage and request performance.
  - **Grafana** visualizes metrics with custom dashboards and sends alerts based on pre-defined conditions.

## Architecture

The application is composed of the following microservices:

1. **User Service** - Manages users and authentication.
2. **Product Service** - Manages product inventory.
3. **Order Service** - Handles user orders and updates product quantities.
4. **PostgreSQL Database** - A single shared database for all services (future potential for microservice-specific databases and scaling).

Each service communicates via RESTful APIs, and all services are containerized using Docker and deployed on Kubernetes.

## Deployment

### Requirements

- **Docker**
- **Kubernetes (Minikube or other cluster manager)**
- **kubectl**
- **Postman (for API testing)**

### Steps to Run

1. **Set up Kubernetes cluster**:
   ```bash
   minikube start
   ```

2. **Deploy PostgreSQL**:
   Deploy the `postgres-service` by applying the corresponding Kubernetes YAML configuration file:
   ```bash
   kubectl apply -f postgres-deployment.yaml
   ```

3. **Deploy Microservices**:
   Deploy each microservice by applying the respective deployment and service YAML files:
   ```bash
   kubectl apply -f user-service.yaml
   kubectl apply -f product-service.yaml
   kubectl apply -f order-service.yaml
   ```

4. **Configure Ingress**:
   Set up the ingress controller to route traffic to the respective services:
   ```bash
   kubectl apply -f ingress.yaml
   ```

5. **Monitoring with Prometheus and Grafana**:
   - Deploy **Prometheus** and **Grafana** using their respective configuration files:
     ```bash
     kubectl apply -f grafana.yaml
     ```

   - Access **Grafana** via the browser to view dashboards and set up alerts.

6. **Accessing the Application**:
   Access the application through the configured ingress:
   - **User Service**: http://ecommerce.local/user
   - **Product Service**: http://ecommerce.local/product
   - **Order Service**: http://ecommerce.local/order

## Monitoring and Alerts

### Dashboards

Grafana provides visual dashboards for metrics such as:

- CPU and memory usage
- Request/response time for each microservice
- Database query performance
- Error rates

### Alerts

Prometheus and Grafana are configured to send alerts for:

- Microservice failures
- High error rates
- Resource exhaustion (CPU/memory)

## Future Improvements

- **Database Separation**: Option to move to a separate database for each microservice as the system scales.
- **Scaling**: Implement autoscaling for services based on resource usage.
- **Security Enhancements**: Add JWT-based authentication between services.
