#!/usr/bin/env node
const express = require('express');
const cors = require('cors');
const csvParser = require('papaparse');
const { InfluxDB } = require('@influxdata/influxdb-client');
const WebSocket = require('ws');

// WebSocket server URL
const wsUrl = 'wss://centralreg-necuf5ddgq-ue.a.run.app/services';

// Existing code for InfluxDB connection and variables
const baseURL = "http://35.236.200.122:8086/";
const influxToken = "I_UycfPULIG3VFr6eT-b0EzSIESMVb6rxZlS3n49zwHAcmpjPXQPS4u0eaZNY69hsWIVErE--T3lodcHQyx5rA==";
const orgID = "8d3c99041893ac29";


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
  console.log('Received message:', serviceDetails);
  fetchData();
});

wsClient.on('error', function error(error) {
  console.error('WebSocket error:', error);
});

wsClient.on('close', function close() {
  console.log('WebSocket connection closed');
});



  let latestData = {};
  let protocol="";
  let status="";
  let centralRegMetrics= {
    cpu_usage: 0,
    disk_free: 0,
    disk_total: 0,
    disk_used: 0,
    memory_free: 0,
    memory_total: 0,
    memory_used: 0,
    net_bytes_recv: 0,
    net_bytes_sent: 0,
    net_packets_recv: 0,
    net_packets_sent: 0
  };


  function populateMetrics(dataEntry) {
    const field = dataEntry['string_1'];
    const value = Number(dataEntry['double']);
    // Use bracket notation to use the value of 'field' as the key
    if (field !== 'field' && !isNaN(value) && field !== '') {
      centralRegMetrics[field] = value;
    }
  }
  

function queryCentralRegBucket() {
  let queryCentralReg = `
  from(bucket: "CentralReg")
    |> range(start: -1h)
    |> filter(fn: (r) => r["_measurement"] == "system_metrics")
    |> filter(fn: (r) => r["_field"] == "cpu_usage" or r["_field"] == "disk_free" or r["_field"] == "disk_total" or r["_field"] == "disk_used" or r["_field"] == "memory_free" or r["_field"] == "memory_total" or r["_field"] == "memory_used" or r["_field"] == "net_bytes_recv" or r["_field"] == "net_bytes_sent" or r["_field"] == "net_packets_recv" or r["_field"] == "net_packets_sent")
    |> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
    |> yield(name: "mean")
  `;
  let csvDataCentralReg = [];

  // Query the CentralReg bucket
  queryApi.queryLines(queryCentralReg, {
    next(line) {
      csvDataCentralReg.push(line);
    },
    error(error) {
      console.error('Error querying CentralReg bucket:', error);
    },
    complete() {
      console.log('Finished querying CentralReg bucket');
      let processedData = processCsvData(csvDataCentralReg, 'centralRegMetrics');
    }
  });
}

  function fetchData() {
    queryCentralRegBucket();
    Object.keys(serviceDetails).forEach(serviceName => {
      const serviceInfo = serviceDetails[serviceName];
    const serviceType = serviceInfo.type;
    const serviceStatus = serviceInfo.status;
    //console.log(serviceName)
      fetchDataForService(serviceName, serviceType, serviceStatus); // Pass both name and type to the function
    });
  }


  function fetchDataForService(serviceName,serviceType, serviceStatus) {
    let csvDataTesting = [];
    let queryTesting = constructQuery(serviceName, "combined_metrics");

    
    queryApi.queryLines(queryTesting, {
      next(line) {
        //console.log(serviceStatus)
        protocol=serviceType
        if(serviceStatus==undefined){
          serviceStatus="Down"
        }else{
            status=serviceStatus
        }
  
        csvDataTesting.push(line);
      },
      error(error) {
        console.error('Error querying testing bucket:', error);
      },
      complete() {
        //console.log('Finished querying testing bucket');
        processCsvData(csvDataTesting, 'combined_metrics');
        // Fetch data from "gRPC-Metrics" bucket after completing the first query
      }
    });

  }

  
  function constructQuery(serviceName, bucketName) {
    // Dynamically construct the query based on serviceName and bucketName
    console.log(serviceName)
    let query = `
    from(bucket: "${bucketName}")
      |> range(start: -30d)
      |> filter(fn: (r) => r["_measurement"] == "${serviceName}")
      |> filter(fn: (r) => r["_field"] == "status_code" or r["_field"] == "response_size" or r["_field"] == "request_size" or r["_field"] == "request_count" or r["_field"] == "latency" or r["_field"] == "error_count")
      |> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
      |> yield(name: "mean")
      |> fill(usePrevious: true)`;
    return query;
  }
  
      
  function processCsvData(csvData, bucketName) {
    const parsedData = csvParser.parse(csvData.join('\n'), { header: true, skipEmptyLines: true }).data;

  
    const unwantedKeys = ["true", "", "path","endpoint","string"];
  
    // Temporary storage to collect data points for each timestamp
    const tempData = {};
   
  parsedData.forEach(data => {
    if(bucketName=='centralRegMetrics'){
     // console.log(parsedData);
      populateMetrics(data)
      //console.log(centralRegMetrics)
      return;
      }

    const endpoint = data['string_3'];
    const timestamp = data['dateTime:RFC3339_2'];
    const service= data['string_2']
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
        bucket:service,
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
    
    const validDataEntries = Object.entries(timestamps).reduce((entries, [timestamp, data]) => {
      if (data.count !== 0 || data.durationSum !== 0 || data.request_size !== 0 || data.response_size !== 0 || data.status_code !== 0) {
        entries.push({
          timestamp: timestamp,
          bucket:data.bucket,
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
    }
  }
  //console.log(latestData)
  //console.log(centralRegMetrics)
    }
    
// Fetch data immediately on start and then every 10 seconds
fetchData();
setInterval(fetchData, 10 * 1000);

app.get('/data/explorer', (req, res) => {
  res.json(latestData);
});

app.get('/websocket-entries-count', (req, res) => {
  const count = Object.keys(serviceDetails).length; // Count of unique service names
  res.json({ count }); // Send the count as JSON response
});
app.get('/central-reg-metrics', (req, res) => {
  res.json(centralRegMetrics);
});
app.listen(port, () => {
  console.log(`listening on port :${port}`);
});
