# API-Observability

BU CS528: Cloud Computing

Project: API Observability

Team Members: Jack Culley, Zachary Tan, Weiru He, Kai Imery, Young-Chan Cho

Mentor: Surya Jayanthi

## 1.   Vision and Goals Of The Project:

Our vision is to empower organizations with the tools and capabilities needed to ensure the performance, security, and reliability of their APIs in real-time. We aim to create a seamless and user-friendly observability solution that provides immediate insights into API behavior, enabling faster decision-making, proactive issue resolution, and enhanced user experiences.

### Must Have:

#### Real-Time API Monitoring

- **Goal:** Enable real-time monitoring of API traffic, performance, and behavior.
- **Objective:** Ensure that you can react to issues and opportunities as they happen, minimizing downtime and maximizing user satisfaction.

#### Comprehensive Data Collection

- **Goal:** Collect and store comprehensive data about API interactions and system health.
- **Objective:** Provide rich insights into API usage, errors, latency, and other critical metrics for better decision-making.

#### User-Friendly Visualization

- **Goal:** Provide intuitive and customizable dashboards for data visualization.
- **Objective:** Enable users to easily access and interpret API observability data, with the ability to create personalized views.

#### Anomaly Detection and Insights

- **Goal:** Develop machine learning-based anomaly detection and insights generation.
- **Objective:** Automatically identify unusual patterns and provide actionable insights to address issues and opportunities.

### Optional:

#### Scalability and Performance

- **Goal:** Create a scalable and high-performance observability system.
- **Objective:** Ensure that the system can handle increasing data volumes and maintain responsiveness under heavy loads.

#### Security and Compliance

- **Goal:** Implement robust security measures and compliance with data privacy regulations.
- **Objective:** Protect sensitive data, user privacy, and adhere to industry-specific compliance standards.

#### Integration and Extensibility

- **Goal:** Facilitate easy integration with external tools and services.
- **Objective:** Allow you to leverage the observability data in your existing workflows, such as incident management and analytics.

#### Documentation and Support

- **Goal:** Provide comprehensive documentation and support resources.
- **Objective:** Empower users to effectively implement and maximize the benefits of the observability system.


## 2. Users/Personas Of The Project:

This product would be designed for various stakeholders involved in development and upkeep of API development and management. The product how ever will target: 

 - **API Developers**: Individuals responsible for designing, building, and maintaining REST and gRPC APIs so that they can efficiently manage and document their system.
 - **API Consumers**: Developers or applications that rely on their organization's APIs who want easy-access to comprehensive API documentation and authentication methods.
 - **DevOps Teams**: Operation teams that aim to deploy and maintain APIs in the cloud for high availability and performance tracking
 - **Product Managers**: Individuals who monitor the metrics, and user-behaviors of their system to make informed data-driven decisions. 


## 3.   Scope and Features Of The Project:

- Centralized Register Development: Development of a centralized register for managing REST and gRPC APIs, including documentation, versioning, and access controls.
- API Usage Monitoring: Implementation of API usage monitoring features, capturing metrics like request/response times, error rates, and request counts.
- Security Integration: Integration of security features, including access controls and authentication mechanisms for secure API management.
- Cloud Development Practices: Incorporation of cloud development practices, such as containerization and cloud deployment, for scalability and flexibility.
- Compliance and Governance: Support for compliance tracking and reporting to ensure adherence to organizational and industry standards in API governance.
** **

## 4. Solution Concept

### Overview
We aim to develop an observability solution for API services, built around an open-source GoLang module. This module will offer extensive metrics, data storage, and visualization by leveraging several components to ensure scalability. Our development strategy is modular, focusing initially on a mock application. This ensures thorough design and testing phases before broadening the solution into a comprehensive GoLang module.

### System Architecture Overview

This document provides an overview of the system architecture for a comprehensive observability platform designed to monitor REST and gRPC API services.

### Components

The system consists of the following main components:

#### API Service (User-Provided)

- Must be a REST or gRPC API service.
- **Observability Module**: Contains functions for registration and instrumentation.
  - **Registration Function**: Runs within a go routine to connect the API service to the Central Register.
  - **Instrumentation Function**: Collects metrics from the API service and sends them to the Central Register.

#### Central Register

- **Server**: Stores details of registered services and collects metrics.
- **/services WebSocket**: Maintains a list of registered services with online/offline status through heartbeat checks.
- **/register WebSocket**: Used for service registration and maintaining a persistent connection for status monitoring.
- **/metrics WebSocket**: Receives metrics from API services and writes them to InfluxDB.
- **System Metrics**: Collects and provides its own system metrics to ensure the system's health.

#### InfluxDB

- **Time-Series Database**: Stores metrics from both the API services and the Central Register.

#### Metrics Dashboard

- **Visualization Interface**: Displays metrics from InfluxDB.
- **Service Metrics**: Pulls metrics for each API service based on the list retrieved from the Central Register.
- **System Metrics**: Displays the health and status of the Central Register itself.

### Data Flow

1. The API service, equipped with the observability module, registers with the Central Register using the `/register` WebSocket.
2. The API service sends collected metrics to the Central Register via the `/metrics` WebSocket.
3. The Central Register writes these metrics to InfluxDB.
4. The Central Register also collects its own system metrics and stores them in InfluxDB.
5. The Metrics Dashboard retrieves a list of services from the Central Register.
6. The Dashboard queries InfluxDB to retrieve and display metrics for each listed API service, as well as the system metrics for the Central Register.

### Observability Module Functions

- **Registration Function**: This asynchronous function ensures the API service is registered and consistently connected to the Central Register.
- **Instrumentation Function**: This function is responsible for gathering relevant metrics from the API service operations and transmitting them for centralized storage and monitoring.



### System Architecure
![image](https://github.com/EC528-Fall-2023/API-Observability/assets/113144839/3bc9271a-b907-43bb-bab1-23a97a0e08fe)

### Mock Application Architecture
![Mock Application Architecture](https://github.com/EC528-Fall-2023/API-Observability/assets/113144839/1f7d0c4d-6627-457c-a6fa-af90d978a06f)



## 5. Acceptance Criteria

### Inventory Creation:
- The system will host gRPC and RESTful API endpoints on a cloud-based infrastructure.
- Users can create and manage an API inventory, capturing key details such as name, description, version, owner, and endpoint URL.
- These APIs will also facilitate user-related data retrieval via the Google Sheets API.

### Frontend UI:
- Instead of Grafana, the implementation will leverage a custom dashboard tailored to the specific needs and aesthetics of the project.
- This custom dashboard will provide interactive visualizations, charts, and graphs that are suited to the analysis and presentation of data associated with the API.

### Usage Metrics:
- Rather than relying on Prometheus, the system will instrument the API services by collecting natively exposed metrics from GoLang and Protobuf. 
- A central registry will be created to manage and oversee the metrics collected.
- This will serve as a core component for monitoring and observing system performance and health.

### Security and Access Control (Stretch Goal):
- The system should enforce access control to ensure that only authorized personnel can view and modify API details.
- Sensitive information related to APIs should be protected.

### Distribution:
- The end product will be packaged as an installable library.
- Developers and teams can easily incorporate this library into their codebase, enabling them to monitor the metrics of their APIs seamlessly.


## 6.  Release Planning:

### 23A 
#### Project Kickoff and Planning
- Kick off the project with a team meeting.
- Define roles and responsibilities.
- Review and refine the project scope.
- Create a detailed sprint plan.

#### API Development and Cloud Framework Exploration
- Start developing example APIs using gRPC and REST.
- Explore and evaluate cloud platforms for hosting.

### 23B 
#### Frontend UI Wireframing and User Testing Planning
- Begin wireframing the frontend UI components.
- Plan for user testing sessions.

#### Frontend and Backend Connectivity
- Connect frontend UI to gRPC and REST APIs
- Allow frontend to access Google Sheets API through gRPC and Rest

#### Establish Desired Metrics and API Integration with Prometheus
- Define the key performance metrics.
- Begin implementing API endpoints to integrate with Prometheus.

### 23C (MVP)
#### Time Series Database Schema, Backend Persistent Storage Integration, and API Catalog
- Design and implement the time series database schema.
- Integrate backend persistent storage with the observability UI.
- List all APIs in a tabular form.

#### Metric Visualizations and User Registration Class/Header
- Develop visualizations for the desired metrics.
- Implement a class or header for user registration and authentication.

#### MVP Testing and Anomaly Detection Implementation
- Perform testing to ensure the MVP meets basic requirements.
- Integrate anomaly detection mechanisms for key metrics.

### 23D 
#### Polling Mechanism for New Users and Monitoring/Alerting
- Set up Prometheus during development phase for REST and gRPC and Integrate Protobuf for both services.
- Set up and connect InfluxDB for metric storage from API Gateway, ensure dashboard can query data from InfluxDB.
- Create Central Registry for tracking API services' status, and connect Central Registry to the dashboard for real-time updates.

### 23F 
#### Documentation and Final System Testing
- Document the system architecture, API specifications, and database schema.
- Perform final system testing to identify and address any remaining issues.

## 7.   Demos:

### Sprint 1
- https://youtu.be/fMQ8F1_lqdI

### Sprint 2
- https://youtu.be/z9pNmTYq1sk
- https://docs.google.com/presentation/d/1L15A24P99b2qUJilOyeG5eDxHLJTsWDs3QDSy5C90Kw/edit?usp=sharing

### Sprint 3
- https://youtu.be/_5oXQZSA_pc
- https://docs.google.com/presentation/d/1Dt1SpIZKvS8tZmLFcuJYILQoJt9-rrGBiSVRH7_OtLc/edit?usp=sharing

### Sprint 4
- https://youtu.be/Zq9Yxi2Zk7U
- https://docs.google.com/presentation/d/1Q5FWIcdhTDRx7qI9gEtD5dZ6F1lDaFT6khgk_J9iC_o/edit?usp=sharing

### Sprint 5
- https://youtu.be/bbgZpS2PvHc
- https://docs.google.com/presentation/d/1UHGiNDGGkUJeq2H33PovCq71lBuRnhyitBJui3Gx9a4/edit?usp=sharing

### Final Demo
- https://docs.google.com/presentation/d/1II75PosUCc73MHUXuw6IjlxI9E9cSricR8MzWkkaQrA/edit?usp=sharing

## 8.   Further Instructions on Installation and Usage
Attached [here](https://drive.google.com/file/d/17NFmHx3TcPzD-7W9Xvvs5cRgVGjYTfKv/view?usp=sharing) is a document with detailed descriptions of each component of this project, along with instructions on installation and each function



** **
