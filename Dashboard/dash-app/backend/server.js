#!/usr/bin/env node
const express = require('express')
const cors = require('cors')
const csvParser = require('papaparse');

const { InfluxDB, flux } = require('@influxdata/influxdb-client')

// vars to connect to bucket in influxdb
const baseURL = "http://34.86.236.100/"; // url of your local InfluxDB instance
const influxToken = "I_UycfPULIG3VFr6eT-b0EzSIESMVb6rxZlS3n49zwHAcmpjPXQPS4u0eaZNY69hsWIVErE--T3lodcHQyx5rA=="; // token
const orgID = "8d3c99041893ac29"; // org id
const bucket = "testing"; // name of your bucket

// connect to influxdb
const influxDB = new InfluxDB({ url: baseURL, token: influxToken })
const queryApi = influxDB.getQueryApi(orgID)

// start the server
const app = express();
app.use(cors())
const port = 3001;



const dataExplorerQuery = `
from(bucket: "testing")
  |> range(start: -5m)
  |> filter(fn: (r) => r["_measurement"] == "REST")
  |> filter(fn: (r) => r["_field"] == "latency" or r["_field"] == "request_count")
  |> aggregateWindow(every: 1m, fn: min, createEmpty: false)
  |> yield(name: "mean")
  |> fill(usePrevious: true)`


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

        if (data['string_1'] === 'request_count' || data['string_1'] === 'counter') {
          tempData[endpoint][timestamp].count = Number(data['double']);
        } else if (data['string_1'] === 'latency') {
          tempData[endpoint][timestamp].durationSum = Number(data['double']);
        }
      });

      // Convert temporary data to the desired format
      latestData = {};
      for (const [endpoint, timestamps] of Object.entries(tempData)) {
        const validDataEntries = Object.entries(timestamps).reduce((entries, [timestamp, data]) => {
          if (data.count !== 0 || data.durationSum !== 0) {
            entries.push({
              timestamp: timestamp,
              count: data.count,
              request_rate: data.count / 60,
              response_rate: data.count === 0 ? null : data.durationSum / data.count,
              status: 'Online'
            });
          }
          return entries;
        }, []);

        if (validDataEntries.length > 0) {
          latestData[endpoint] = validDataEntries;
        }
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
