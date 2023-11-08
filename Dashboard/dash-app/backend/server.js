#!/usr/bin/env node
const express = require('express')
const cors = require('cors')
const csvParser = require('papaparse');

const { InfluxDB, flux } = require('@influxdata/influxdb-client')

// vars to connect to bucket in influxdb
const baseURL = "http://localhost:8086"; // url of your local InfluxDB instance
const influxToken = "hzIhrWi_NhfpE-32wfVTkebF98z_X_jcuOcDmGlzhjxOgYT9uMnNKpa1gRBaDmE7OrcMC7PQij6-DqEiTFeKIw=="; // token
const orgID = "1aaa9955ce11b8e7"; // org id
const bucket = "metrics"; // name of your bucket

// connect to influxdb
const influxDB = new InfluxDB({ url: baseURL, token: influxToken })
const queryApi = influxDB.getQueryApi(orgID)

// start the server
const app = express();
app.use(cors())
const port = 3001;



const dataExplorerQuery = `
  from(bucket: "metrics")
  |> range(start: -5m)
  |> filter(fn: (r) => r["_measurement"] == "http_request_duration_seconds" or r["_measurement"] == "http_requests_total")
  |> filter(fn: (r) => r["_field"] == "sum" or r["_field"] == "count" or r["_field"] == "counter")
  |>  aggregateWindow(every: 1m, fn: min, createEmpty: false)
  |> yield(name: "mean")`;

  let latestData = [];

  // Extracted function to query InfluxDB and update `latestData`
  function fetchData() {
    let csvData = [];
  
    queryApi.queryLines(dataExplorerQuery, {
      next(line) {
        csvData.push(line);
      },
      error(error) {
        console.error(error);
        console.log('\nFinished querying ERROR');
      },
      complete() {
        console.log('\nFinished querying SUCCESS');
      
        const parsedData = csvParser.parse(csvData.join('\n'), { header: true, skipEmptyLines: true }).data;
  
        // Unwanted keys to ignore
        const unwantedKeys = ["true", "", "path"];
  
        // Temporary storage to collect data points for each timestamp
        const tempData = {};
  
        parsedData.forEach(data => {
          const endpoint = data['string_3'];
          const timestamp = data['dateTime:RFC3339_2'];
  
          if (unwantedKeys.includes(endpoint)) {
            return; // Skip unwanted keys
          }
  
          if (!tempData[endpoint]) {
            tempData[endpoint] = {};
          }
  
          if (!tempData[endpoint][timestamp]) {
            tempData[endpoint][timestamp] = {
              count: 0,
              durationSum: 0
            };
          }
  
          if (data['string_1'] === 'count' || data['string_1'] === 'counter') {
            tempData[endpoint][timestamp].count = Number(data['double']);
          } else if (data['string_1'] === 'sum') {
            tempData[endpoint][timestamp].durationSum = Number(data['double']);
          }
        });
  
        // Convert temporary data to the desired format
        latestData = {};
        for (const [endpoint, timestamps] of Object.entries(tempData)) {
          latestData[endpoint] = Object.entries(timestamps).map(([timestamp, data]) => {
            return {
              timestamp: timestamp,
              count: data.count,
              request_rate: data.count / 60,
              response_rate: data.durationSum / data.count,
              status:'Online'
            };
          });
        }
      },
    });
  }
  
  // Fetch data immediately on start and then every 10 seconds
  fetchData();
  setInterval(fetchData, 10 * 1000);
  
  app.get('/data/explorer', (req, res) => {
    res.json(latestData);
  });
  
  app.listen(port, () => {
    console.log(`listening on port :${port}`);
  });