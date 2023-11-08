# Dashboard App

## Overview

This Dashboard App is a full-stack application designed to display metrics and analytics from various APIs, sourced from InfluxDB. The back end is built with Node.js, using Express for server-side logic and InfluxDB as the data store. The front end is developed using React and Material-UI for a responsive user interface.

## Features

- Real-time data fetching from InfluxDB.
- Periodic data refresh every 10 seconds.
- Interactive charts and graphs for data visualization.
- REST API endpoints to serve the latest data.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Before you begin, ensure you have the following installed:

- Node.js
- npm (comes with Node.js)
- InfluxDB
- npm install to get dependencies

### Setup
Adjust the value of the following variables

- INFLUX_URL=http://localhost:8086
- INFLUX_TOKEN=your_influxdb_token
- INFLUX_ORG=your_org_id
- INFLUX_BUCKET=your_bucket_name

### Run

- node server.js
- npm start

