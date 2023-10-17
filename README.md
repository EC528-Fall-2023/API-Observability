# API-Observability

BU CS528: Cloud Computing

Project: API Observability

Team Memebers: Jack Culley, Zachary Tan, Weiru He, Kai Imery, Young-Chan Cho

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

The architecture of our solution starts with a data collection layer, which consists of an API gateway and instrumentation. The API gateway will capture data about each incoming request, such as request headers, payload, and metadata, while the instrumentation will collect data on request handling, response times, errors, and other relevant metrics. This information will then be passed into a data ingestion and processing layer, which includes a message queue that will be responsible for publishing data events into the queue. From here, the information will be either stream-processed or batch-processed. The stream processing will be done to handle the real-time analytics and log/event processing. The batch processing will be done to handle historical data that can be used for analytics, machine learning, or generating reports. The next layer of our solution is data storage. The data storage layer contains a time-series database which is responsible for storing time-series data for metrics such as request/response times, error rates, and traffic volume. Additionally, there will be log/trace storage which will provide detailed information regarding the log and trace data that is collected. The next layer is data visualization and dashboarding. This will be used to visualize and monitor all of the information that is being collected and gives the ability to create alert systems for different metrics. 
 
![Screenshot 2023-09-21 141620](https://github.com/EC528-Fall-2023/API-Observability/assets/114025961/df137ab1-b13f-492d-9b5b-d22bbfe8f5f9)


## 5. Acceptance criteria

#### Inventory Creation:

- The system will host gRPC and RESTful API endpoints on a cloud-based infrastructure.
- Users can create and manage an API inventory, capturing key details such as name, description, version, owner, and endpoint URL.
- These APIs will also facilitate user-related data retrieval via the Google Sheets API.

#### Frontend UI:
- The implementation shall incorporate Grafana as the platform for crafting interactive and adaptable dashboards, charts, and graphs.
- These visualizations will facilitate the analysis and presentation of data associated with the API.

#### Usage Metrics:
- Prometheus shall be employed to implement a suite of 10 to 12 straightforward metrics, serving as a vital component for monitoring and observing system performance and health.

#### Security and Access Control (Stretch Goal):
- The system should enforce access control to ensure that only authorized personnel can view and modify API details.
- Sensitive information related to APIs should be protected.

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
- Develop a mechanism to poll and capture data for new users.
- Implement monitoring and alerting mechanisms for the observability system.

### 23F 
#### Documentation and Final System Testing
- Document the system architecture, API specifications, and database schema.
- Perform final system testing to identify and address any remaining issues.

## 7.   Demos:

### Sprint 1
- https://youtu.be/fMQ8F1_lqdI

### Sprint 2
- https://youtu.be/z9pNmTYq1sk


## General comments

Remember that you can always add features at the end of the semester, but you can't go back in time and gain back time you spent on features that you couldn't complete.

** **
