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

### Architecture Breakdown

#### 1. API Gateway
- The mock application features both **REST** and **gRPC** API services, accessed via a **ReactJS client**.
- We've integrated **Protobuf** with both REST and gRPC API services for scalability.
- For metric extraction, the API Gateway will be instrumented with **Prometheus** during the developmental phase. However, the final version will rely on natively exposed metrics from GoLang and Protobuf.

#### 2. InfluxDB
- This component serves as a repository for metrics, performance, and usage data collected from the API Gateway instrumentation.
- The front-end observability dashboard will query **InfluxDB** to present the stored data to users.

#### 3. Central Registry
- Designed to provide scalability, especially for users managing multiple API services.
- Functions akin to zero-configuration networking: API services broadcast their status and essential details (name, service type, domain, port, metadata) upon changing their online status.
- The Central Registry, acting as a listener, captures this information, updating a structured list that tracks the current status of all API services. This registry communicates with the front-end dashboard, ensuring that users can access an up-to-date list of services. This list, in turn, helps in querying InfluxDB for metrics.

#### 4. Front-End Observability Dashboard
- This dashboard offers a user-friendly interface, allowing users to monitor all system APIs seamlessly.
- Users can filter APIs based on type or name and delve deeper into specific APIs for detailed metrics and insights.

### System Architecure
![Observability Architecure](https://github.com/EC528-Fall-2023/API-Observability/assets/113144839/35e18823-ff49-4dd3-9faa-75c4b882792a)
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


** **
