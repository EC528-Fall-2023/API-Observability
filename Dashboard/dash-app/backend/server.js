#!/usr/bin/env node
const express = require('express');
const cors = require('cors');
const csvParser = require('papaparse');
const { InfluxDB } = require('@influxdata/influxdb-client');
const WebSocket = require('ws');

// WebSocket server URL
const wsUrl = 'wss://centralreg-necuf5ddgq-ue.a.run.app/services';

// Existing code for InfluxDB connection and variables
const baseURL = "http://34.86.236.100/";
const influxToken = "I_UycfPULIG3VFr6eT-b0EzSIESMVb6rxZlS3n49zwHAcmpjPXQPS4u0eaZNY69hsWIVErE--T3lodcHQyx5rA==";
const orgID = "8d3c99041893ac29";
const bucket = "testing";

const influxDB = new InfluxDB({ url: baseURL, token: influxToken });
const queryApi = influxDB.getQueryApi(orgID);

// Existing server setup
const app = express();
app.use(cors());
//const port = 3001;
const port = process.env.PORT || 3001;

// WebSocket connection
const wsClient = new WebSocket(wsUrl);
let serviceDetails ={}; // Store unique service names

wsClient.on('message', function incoming(message) {
  const data = JSON.parse(message.toString());
  data.forEach(item => {
    serviceDetails[item.Name] = {type: item.Type, status:item.Status}; // Store name and its associated type
  });
  console.log('Received message:', data);
  fetchData();
});


wsClient.on('message', function incoming(message) {
  const messageString = message.toString();
  console.log('Received message:', messageString);
});

wsClient.on('error', function error(error) {
  console.error('WebSocket error:', error);
});

wsClient.on('close', function close() {
  console.log('WebSocket connection closed');
});



// const dataExplorerQueryTesting = `
// from(bucket: "testing")
//   |> range(start: -1h)
//   |> filter(fn: (r) => r["_measurement"] == "REST")
//   |> filter(fn: (r) => r["_field"] == "status_code" or r["_field"] == "response_size" or r["_field"] == "request_size" or r["_field"] == "request_count" or r["_field"] == "latency" or r["_field"] == "error_count")
//   |> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
//   |> yield(name: "mean")
//   |> fill(usePrevious: true)`


  const dataExplorerQueryGrpcMetrics = `
  from(bucket: "gRPC-Metrics")
    |> range(start: -24h)
    |> filter(fn: (r) => r["_measurement"] == "gRPCMetrics")
    |> filter(fn: (r) => r["_field"] == "request_count" or r["_field"] == "request_size" or r["_field"] == "response_size" or r["_field"] == "error_rate" or r["_field"] == "duration")
    |> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
    |> yield(name: "mean")
    |> fill(usePrevious: true)`;
  
  let latestData = {};
  let protocol="";
  let status="";

  function fetchData() {
    Object.keys(serviceDetails).forEach(serviceName => {
      const serviceInfo = serviceDetails[serviceName];
    const serviceType = serviceInfo.type;
    const serviceStatus = serviceInfo.status;
      fetchDataForService(serviceName, serviceType, serviceStatus); // Pass both name and type to the function
    });
  }


  function fetchDataForService(serviceName,serviceType, serviceStatus) {
    let csvDataTesting = [];
    let queryTesting = constructQuery(serviceName, "combined_metrics");

    
    queryApi.queryLines(queryTesting, {
      next(line) {
        protocol=serviceType
        status=serviceStatus
        csvDataTesting.push(line);
      },
      error(error) {
        console.error('Error querying testing bucket:', error);
      },
      complete() {
        console.log('Finished querying testing bucket');
        processCsvData(csvDataTesting, 'combined_metrics');
        // Fetch data from "gRPC-Metrics" bucket after completing the first query
      }
    });

  }

  
  function constructQuery(serviceName, bucketName) {
    // Dynamically construct the query based on serviceName and bucketName
    let query = `
    from(bucket: "${bucketName}")
      |> range(start: -1h)
      |> filter(fn: (r) => r["_measurement"] == "${serviceName}")
      |> filter(fn: (r) => r["_field"] == "status_code" or r["_field"] == "response_size" or r["_field"] == "request_size" or r["_field"] == "request_count" or r["_field"] == "latency" or r["_field"] == "error_count")
      |> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
      |> yield(name: "mean")
      |> fill(usePrevious: true)`;
    return query;
  }
  
  
  // function fetchGrpcMetricsData() {
  //   let csvDataGrpcMetrics = [];
  
  //   // Fetch data from "gRPC-Metrics" bucket
  //   queryApi.queryLines(dataExplorerQueryGrpcMetrics, {
  //     next(line) {
  //       protocol="gRPC"
  //       csvDataGrpcMetrics.push(line);
  //     },
  //     error(error) {
  //       console.error('Error querying gRPC-Metrics bucket:', error);
  //     },
  //     complete() {
  //       console.log('Finished querying gRPC-Metrics bucket');
  //       processCsvData(csvDataGrpcMetrics, 'grpcMetrics');
  //     }
  //   });
  // }
      
  function processCsvData(csvData, bucketName) {
    const parsedData = csvParser.parse(csvData.join('\n'), { header: true, skipEmptyLines: true }).data;
    // Unwanted keys to ignore
    const unwantedKeys = ["true", "", "path","endpoint","string"];
  
    // Temporary storage to collect data points for each timestamp
    const tempData = {};
     //console.log(parsedData);
  parsedData.forEach(data => {
    const endpoint = data['string_3'];
    const timestamp = data['dateTime:RFC3339_2'];

    if (unwantedKeys.includes(endpoint)) {
      return; // Skip unwanted keys
    }
    if (data['#datatype'] === '#datatype' || data['#datatype'] === '#group' || data['#datatype'] === '#default') {
      return;
    }

    if (!tempData[endpoint]) {
      tempData[endpoint] = {};
    }

    if (!tempData[endpoint][timestamp]) {
      tempData[endpoint][timestamp] = {
        count: 0,
        durationSum: 0,
        error_count:0,
        request_size:0,
        response_size:0,
        status_code:0
      };
    }

    switch (data['string_1']) {
      case 'request_count':
        tempData[endpoint][timestamp].count = Math.ceil(Number(data['double']));
        break;
      case 'latency':
      case 'duration':
        tempData[endpoint][timestamp].durationSum = Number(data['double']);
        break;
      case 'error_count':
      case 'error_rate':
        tempData[endpoint][timestamp].error_count = Number(data['double']);
        break;
      case 'request_size':
        tempData[endpoint][timestamp].request_size = Number(data['double']);
        break;
      case 'response_size':
        tempData[endpoint][timestamp].response_size = Number(data['double']);
        break;
      case 'status_code':
        tempData[endpoint][timestamp].status_code = Number(data['double']);
        break;
    }
    
    //console.log(tempData)

  });

  // Convert temporary data to the desired format
  
  for (const [endpoint, timestamps] of Object.entries(tempData)) {
    latestData[endpoint] = {};
    console.log(endpoint)
    const validDataEntries = Object.entries(timestamps).reduce((entries, [timestamp, data]) => {
      if (data.count !== 0 || data.durationSum !== 0 || data.request_size !== 0 || data.response_size !== 0 || data.status_code !== 0) {
        entries.push({
          timestamp: timestamp,
          count: data.count,
          request_rate: data.count / 60,
          response_rate: data.durationSum/1000,
          error_count: data.error_count,
          request_size: data.request_size,
          response_size: data.response_size,
          status_code: data.status_code,
          status: status,
          protocol:protocol
        });
      }
      return entries;
    }, []);
  
    if (validDataEntries.length > 0) {
      latestData[endpoint] = validDataEntries;
      //console.log("ENTERED")
    }
  }
        
    }

// Fetch data immediately on start and then every 10 seconds
fetchData();
setInterval(fetchData, 10 * 1000);

app.get('/data/explorer', (req, res) => {
  //console.log("SENT")
  res.json(latestData);
});

app.get('/websocket-entries-count', (req, res) => {
  const count = Object.keys(serviceDetails).length; // Count of unique service names
  res.json({ count }); // Send the count as JSON response
});

app.listen(port, () => {
  console.log(`listening on port :${port}`);
});
