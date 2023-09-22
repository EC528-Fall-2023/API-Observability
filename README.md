# API-Observability

## 1.   Vision and Goals Of The Project:

Our vision is to empower organizations with the tools and capabilities needed to ensure the performance, security, and reliability of their APIs in real-time. We aim to create a seamless and user-friendly observability solution that provides immediate insights into API behavior, enabling faster decision-making, proactive issue resolution, and enhanced user experiences.

#### Real-Time API Monitoring

- **Goal:** Enable real-time monitoring of API traffic, performance, and behavior.
- **Objective:** Ensure that you can react to issues and opportunities as they happen, minimizing downtime and maximizing user satisfaction.

#### Comprehensive Data Collection

- **Goal:** Collect and store comprehensive data about API interactions and system health.
- **Objective:** Provide rich insights into API usage, errors, latency, and other critical metrics for better decision-making.

#### Scalability and Performance

- **Goal:** Create a scalable and high-performance observability system.
- **Objective:** Ensure that the system can handle increasing data volumes and maintain responsiveness under heavy loads.

#### Security and Compliance

- **Goal:** Implement robust security measures and compliance with data privacy regulations.
- **Objective:** Protect sensitive data, user privacy, and adhere to industry-specific compliance standards.

#### User-Friendly Visualization

- **Goal:** Provide intuitive and customizable dashboards for data visualization.
- **Objective:** Enable users to easily access and interpret API observability data, with the ability to create personalized views.

#### Anomaly Detection and Insights

- **Goal:** Develop machine learning-based anomaly detection and insights generation.
- **Objective:** Automatically identify unusual patterns and provide actionable insights to address issues and opportunities.

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

The Scope places a boundary around the solution by detailing the range of features and functions of the project. This section helps to clarify the solution scope and can explicitly state what will not be delivered as well.

It should be specific enough that you can determine that e.g. feature A is in-scope, while feature B is out-of-scope.

### In-Scope Features

#### API Health Checks

- Scheduled and on-demand checks to ensure the availability and responsiveness of APIs.
- Detailed reporting on API uptime, downtimes, and performance metrics.

#### Alerting and Notification System

- Proactive alerts on potential issues or breaches based on user-defined thresholds.
- Integration with communication platforms like Slack, email, and SMS for immediate notifications.

#### Data Retention and Backup

- Customizable data retention policies allowing users to determine how long their data is stored.
- Automated backup systems ensuring no loss of crucial observability data.

#### Interactive Query Interface

- Users can directly query their stored data, allowing for custom insights and reports.
- Support for common query languages to ensure ease of use.

#### Collaboration Tools

- Features that enable team members to collaborate on detected issues, share insights, and annotate data.
- Integration with project management tools for streamlined workflows.

#### Mobile Access

- Mobile applications or responsive design ensuring stakeholders can monitor and react to API issues on-the-go.

#### Custom Integrations

- SDKs or plugins to extend the platform's capabilities, integrate with third-party tools, or support newer technologies.

#### Advanced Reporting

- Generation of weekly, monthly, or custom duration reports.
- Export options in various formats such as PDF, CSV, etc.

### Out-of-Scope Features

#### Multi-language support for the dashboard interface

- Refers to the inclusion of multiple language options in the dashboard interface. It allows users to view the dashboard in their preferred language and is especially beneficial for products with an international user base. 

#### Training modules or workshops for teams

- Involves providing educational modules or workshops specifically tailored for teams. It helps users or team members understand how to effectively utilize the product or system.

#### Focusing on cloud-based solutions rather than on-premises deployment

- Suggests that the focus will not be on on-premises deployment but exclusively on cloud-based solutions.

** **

## 4. Solution Concept

The architecture of our solution starts with a data collection layer, which consists of an API gateway and instrumentation. The API gateway will capture data about each incoming request, such as request headers, payload, and metadata, while the instrumentation will collect data on request handling, response times, errors, and other relevant metrics. This information will then be passed into a data ingestion and processing layer, which includes a message queue that will be responsible for publishing data events into the queue. From here, the information will be either stream-processed or batch-processed. The stream processing will be done to handle the real-time analytics and log/event processing. The batch processing will be done to handle historical data that can be used for analytics, machine learning, or generating reports. The next layer of our solution is data storage. The data storage layer contains a time-series database which is responsible for storing time-series data for metrics such as request/response times, error rates, and traffic volume. Additionally, there will be log/trace storage which will provide detailed information regarding the log and trace data that is collected. The next layer is data visualization and dashboarding. This will be used to visualize and monitor all of the information that is being collected and gives the ability to create alert systems for different metrics. 
 
![Screenshot 2023-09-21 141620](https://github.com/EC528-Fall-2023/API-Observability/assets/114025961/df137ab1-b13f-492d-9b5b-d22bbfe8f5f9)


## 5. Acceptance criteria

#### Inventory Creation:
- The system should allow users to create and maintain an inventory of API services.
- The inventory should capture essential details for each API, such as name, description, version, owner, and endpoint URL.

#### API Type Support:
- The system should support the tracking of APIs developed in both REST and gRPC.

#### REST API Details:
- For REST APIs, the system should capture information like HTTP methods (GET, POST, PUT, DELETE), request/response formats (JSON, XML), and authentication mechanisms.

#### gRPC API Details:
- For gRPC APIs, the system should record details about the service methods, request/response message types, and authentication mechanisms.

#### Usage Metrics:
- The system should collect usage metrics for each API, including the number of requests, response times, and error rates.
- Users should be able to view historical usage data and trends for informed decision-making.

#### Security and Access Control:
- The system should enforce access control to ensure that only authorized personnel can view and modify API details.
- Sensitive information related to APIs should be protected.

## 6.  Release Planning:

### Sprint 23A (Sept 22 to Sept 29)
- Meet with mentor to learn demands and requirements for the project
- Familiarize ourselves with gRPC and Observability
- Investigate and plan cloud development considerations, including cloud service selection
- Investigate metric measurement applications, and evaluate intergration options with metrics and visualization tools.

## General comments

Remember that you can always add features at the end of the semester, but you can't go back in time and gain back time you spent on features that you couldn't complete.

** **
