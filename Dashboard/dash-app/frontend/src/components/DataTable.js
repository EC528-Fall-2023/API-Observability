import React, { useState } from 'react';
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Button, Collapse, Box,Typography,MenuItem, Select, InputLabel, FormControl, OutlinedInput, Checkbox, ListItemText  } from '@mui/material';
import { styled } from '@mui/system';
import { useEffect } from 'react';
import LineChart from './LineChart'
import TopGraph from './TopGraph'

const StyledTableCell = styled(TableCell)(({ theme }) => ({
  fontWeight: 'bold',
  backgroundColor: '#F0F0F0',
  padding: '10px',
}));
const StyledButton = styled(Button)(({ theme, requestType }) => {
  let backgroundColor;
  let hoverColor;

  switch (requestType) {
    case 'GET':
      backgroundColor = '#4CAF50'; // Green shade for GET
      hoverColor = '#45a049';
      break;
    case 'POST':
      backgroundColor = '#2196F3'; // Blue shade for POST
      hoverColor = '#0b7dda';
      break;
    case 'PUT':
      backgroundColor = '#FFC107'; // Yellow shade for PUT
      hoverColor = '#e6a800';
      break;
    case 'DELETE':
      backgroundColor = '#f44336'; // Red shade for DELETE
      hoverColor = '#da190b';
      break;
    default:
      backgroundColor = theme.palette.primary.main;
      hoverColor = theme.palette.primary.dark;
  }

  return {
    backgroundColor: backgroundColor,
    '&:hover': {
      backgroundColor: hoverColor,
    },
    borderRadius: '20px',
    textTransform:'capitalize',
    fontSize: '1.0rem',
    padding: '6px 12px',
  };
});

const DashboardInfo = ({ label, value }) => {
  return (
    <Paper elevation={3} style={styles.dashboardInfoContainer}>
      <Typography variant="h6" style={styles.dashboardLabel}>
        {label}
      </Typography>
      <Typography variant="h4" style={styles.dashboardValue}>
        {value}
      </Typography>
    </Paper>
  );
};

const styles = {
  dashboardRow: {
    display: 'flex',
    justifyContent: 'space-between', // this will space out your items evenly
    alignItems: 'center',
    marginBottom: '20px', // add space between this row and the next content
  },
  dashboardInfoContainer: {
    padding: '20px',
    width: 'calc(33.333% - 20px)', 
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#fff', // match your background color
    borderRadius: '10px', // optional, if you want rounded corners
    boxShadow: '0px 4px 10px rgba(0,0,0,0.1)', // adjust as needed for shadow
  },
  dashboardLabel: {
    color: '#999', // grey color for label
    marginBottom: '5px', // space between label and value
  },
  dashboardValue: {
    color: '#333', // darker color for value
    fontWeight: 'bold', // bold font for value
  }
};



const mockData = [
  {
    type: 'POST',
    protocol: 'REST',
    path: '/adduser',
    details: {
      counter: 400,
      responseRate: 2.73,
      errorRate: 12,
      requestRate: 120,
      payloadSize: 2.5
    },
    variables: [
      { name: 'NAME', type: 'string', description: 'Name of User' },
      { name: 'AGE', type: 'integer', description: 'Age of User' },
      { name: 'COLLEGE', type: 'string', description: 'College Attended' },
      { name: 'COMMUTE', type: 'string', description: 'Method of Commute' },
      { name: 'HOBBY', type: 'string', description: 'Favorite Hobby' }
    ]
  },
  {
    type: 'PUT',
    protocol: 'REST',
    path: '/updateuser/:name',
    details: {
      counter: 365,
      responseRate: 4.53,
      errorRate: 17,
      requestRate: 110,
      payloadSize: 1.9
    },
    variables: [
      { name: 'NAME', type: 'string', description: 'Name of User' }
    ]
  },
  {
    type: 'DELETE',
    protocol: 'REST',
    path: '/deleteuser/:name',
    details: {
      counter: 101,
      responseRate: 0.73,
      errorRate: 2,
      requestRate: 135,
      payloadSize: 3.3
    },
    variables: [
      { name: 'NAME', type: 'string', description: 'Name of User' }
    ]
  },
  {
    type: 'GET',
    protocol: 'REST',
    path: '/getuser/:name',
    details: {
      counter: 234,
      responseRate: 1.64,
      errorRate: 9,
      requestRate: 186,
      payloadSize: 2.4
    },
    variables: [
      { name: 'NAME', type: 'string', description: 'Name of User' }
    ]
  }
];


// Helper function to format time from seconds to HH:MM:SS
function formatTime(seconds) {
  const pad = (num) => num.toString().padStart(2, '0');
  const hours = pad(Math.floor(seconds / 3600));
  const minutes = pad(Math.floor((seconds % 3600) / 60));
  const secondsLeft = pad(seconds % 60);
  return `${hours}:${minutes}:${secondsLeft}`;
}

const DetailedInfo = ({ variables }) => {
  return (
    <Box sx={{ padding: 2, width: '100%', height: '100%' }}>
      <TableRow>
        <TableCell sx={{ fontWeight: 'bold' }}>Name & Type</TableCell>
        <TableCell sx={{ fontWeight: 'bold' }}>Description</TableCell>
      </TableRow>
      {variables.map(variable => (
        <TableRow key={variable.name}>
          <TableCell>{variable.name} ({variable.type})</TableCell>
          <TableCell>{variable.description}</TableCell>
        </TableRow>
      ))}
    </Box>
  );
};


const DataTable = () => {
  const [openRow, setOpenRow] = useState(null);
  const [data, setData] = useState([]); // Replace mockData with this state
  const [onlineAPI, setOnline]=useState(0);
  const [regAPI, setReg]=useState(0);
  const [selectedGraphTypes, setSelectedGraphTypes] = useState([]);
  const [graphdata, setGraphdata] = useState([]);
  const [selectedRow, setSelectedRow] = useState(null);

  const [uptime, setUptime] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setUptime(prevUptime => prevUptime + 1); // Increment uptime every second
    }, 1000);

    return () => clearInterval(interval); // Cleanup on unmount
  }, []);

  const handleGraphTypeChange = (event) => {
    const {
      target: { value },
    } = event;
    setSelectedGraphTypes(typeof value === 'string' ? value.split(',') : value);
  };

  useEffect(() => {
    if (selectedRow != null) {
      const selectedRowData = data.find(row => row.path === selectedRow);
      if (selectedRowData) {
        // New code for setting multiple graph types
        const newGraphData = selectedGraphTypes.map(type => ({
          type,
          data: selectedRowData.history.map(entry => ({
            timestamp: entry.timestamp,
            value: entry[type] // The bracket notation allows for dynamic property access
          })),
        }));
        setGraphdata(newGraphData); // Now you're setting an array of graph data objects
      }
    }
  }, [selectedGraphTypes, selectedRow, data]);
  
  
  console.log(graphdata)

  useEffect(() => {
  
    const fetchDataFromServer = async () => {
      try {
          const response = await fetch('http://localhost:3001/data/explorer');
          if (response.ok) {
              const fetchedData = await response.json();  // Rename to fetchedData to avoid shadowing
              const convertedData = convertBackendDataToMockFormat(fetchedData); 
              setData(convertedData);
          } else {
              console.error('Failed to fetch data:', response.statusText);
          }
      } catch (error) {
          console.error('Error fetching data:', error);
      }
  };
  
  
    fetchDataFromServer();
  
    // Refetch every 10 seconds
    const intervalId = setInterval(fetchDataFromServer, 1000);
  
    return () => {
      clearInterval(intervalId);  // Cleanup interval on component unmount
    };
  }, []);

  function convertBackendDataToMockFormat(data) {
    let onlineCount = 0;
    let registeredCount = 0;
    let stamp;
    const newData = [];
  
    for (const path in data) {
      const latestEntry = data[path][data[path].length - 1];  // Get the last entry for each path
      //console.log(latestEntry)
  
      const historyData = data[path].map(entry => {
        return {
          timestamp: entry.timestamp,
          count: entry.count,
          request_rate: entry.request_rate,
          response_rate: entry.response_rate
        };
      });

      if (latestEntry.status === 'Online') {
        onlineCount++; // Increment online APIs count
        registeredCount++;
      }
      else{
        registeredCount++;
      }
      console.log(onlineCount)
      


      newData.push({
        type: 'GET',  // hardcoded as per your requirement
        protocol: 'REST',  // hardcoded as per your requirement
        path: path,
        details: {
          timestamp: latestEntry.timestamp,
          counter: latestEntry.count,  // using mockData as default if no data from backend
          responseRate: latestEntry.response_rate,
          errorRate: 72,
          requestRate: latestEntry.request_rate,
          payloadSize: 2.7,
          status:latestEntry.status
        },
        variables: [
          { name: 'NAME', type: 'string', description: 'Name of User' }
        ],
        history: historyData
      });
    }
    // setOnlineAPIs(onlineCount);
    // setRegisteredAPIs(registeredCount);
    setOnline(onlineCount);
    setReg(registeredCount);
    return newData;
  }



  return (
    <>
      <div style={styles.dashboardRow}>
      <DashboardInfo label="Uptime:" value={formatTime(uptime)} />
      <DashboardInfo label="Online APIs:" value={onlineAPI} />
      <DashboardInfo label="Registered APIs:" value={regAPI} />
    </div>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 400 }} aria-label="customized table">
          <TableHead>
            <TableRow size="small">
              <StyledTableCell>TYPE</StyledTableCell>
              <StyledTableCell>PROTOCOL</StyledTableCell>
              <StyledTableCell>PATH</StyledTableCell>
              <StyledTableCell>COUNTER (Requests)</StyledTableCell>
              <StyledTableCell>RESPONSE RATE (s)</StyledTableCell>
              <StyledTableCell>ERROR RATE (%)</StyledTableCell>
              <StyledTableCell>REQUEST RATE (req/min)</StyledTableCell>
              <StyledTableCell>PAYLOAD SIZE (KB)</StyledTableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map((row) => (
              <React.Fragment key={row.path}>
                <TableRow size="small">
                  <TableCell>
                    <StyledButton
                      variant="contained"
                      requestType={row.type}
                      onClick={() => {
                        setOpenRow(row.path === openRow ? null : row.path);
                        setSelectedRow(row.path === openRow ? null : row.path); // Update selectedRow state here
                      }}
                    >
                      {row.type}
                    </StyledButton>
                  </TableCell>
                  <TableCell sx={{ color: 'grey', fontSize: '25px' }}>{row.protocol}</TableCell>
                  <TableCell>{row.path}</TableCell>
                  <TableCell>{row.details.counter}</TableCell>
                  <TableCell>{row.details.responseRate}</TableCell>
                  <TableCell>{row.details.errorRate}</TableCell>
                  <TableCell>{row.details.requestRate}</TableCell>
                  <TableCell>{row.details.payloadSize}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell colSpan={8}>
                    <Collapse in={openRow === row.path} timeout="auto" unmountOnExit>
                      <div style={{ display: 'flex', justifyContent: 'space-between', width: '100%', height: '100%' }}>
                        <div style={{ flex: 1, paddingRight: '5%', paddingLeft: '10%', width: '100%' }}>
                          <DetailedInfo variables={row.variables} />
                        </div>
                        <div><FormControl sx={{ m: 1, width: 300 }}>
                                <InputLabel id="multiple-checkbox-label">Graph Types</InputLabel>
                                <Select
                                  labelId="multiple-checkbox-label"
                                  id="multiple-checkbox"
                                  multiple
                                  value={selectedGraphTypes}
                                  onChange={handleGraphTypeChange}
                                  input={<OutlinedInput label="Graph Types" />}
                                  renderValue={(selected) => selected.join(', ')}
                                >
                                  {['response_rate', 'count', /* other types */].map((type) => (
                                    <MenuItem key={type} value={type}>
                                      <Checkbox checked={selectedGraphTypes.indexOf(type) > -1} />
                                      <ListItemText primary={type} />
                                    </MenuItem>
                                  ))}
                                </Select>
                              </FormControl>
                              </div>
                        <div style={{ flex: 2 }}>
                        <LineChart dataSets={graphdata} title="Graph Title" />
                        </div>
                      </div>
                    </Collapse>
                  </TableCell>
                </TableRow>
              </React.Fragment>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );
};

export default DataTable;