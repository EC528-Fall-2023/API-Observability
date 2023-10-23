import React, { useState, useEffect } from 'react';
import { Button, Form, Container, Header, Table } from 'semantic-ui-react';
import axios from 'axios';
import './index.css';

function App() {
  const [name, setName] = useState('');
  const [age, setAge] = useState('');
  const [commute_method, setcommute_method] = useState('');
  const [college, setCollege] = useState('');
  const [hobbies, setHobbies] = useState('');
  const [APIdata, setAPIdata] = useState([]);
  const [refresh, setRefresh] = useState([]);
  const [apiMode, setApiMode] = useState('REST'); // Default mode is REST

  const handleSubmit = (e) => {
    e.preventDefault();

    const objt = { name, age, commute_method, college, hobbies };
    const apiUrl = apiMode === 'REST' ? 'http://localhost:4000/adduser' : 'https://sheet.best/api/sheets/ab206fd6-ee69-4b1d-a56f-f29c0ba70176';

    axios.post(apiUrl, objt).then((response) => {
      console.log(response);
      // Add logic to handle the response, e.g., update state or show a message
    });
  };

  const handleUpdate = (e) => {
    e.preventDefault();

    const objt = { name, age, commute_method, college, hobbies };

    axios
        .put(
            `http://localhost:4000/updateuser/${name}`, // Update with your Go server route
            objt
        )
        .then((response) => {
          console.log(response);
          // Add logic to handle the response
        });
  };

  const handleSearch = (e) => {
    e.preventDefault();
  
    // Perform the search using the name in the state
    axios.get(`http://localhost:4000/getuser/${name}`)
      .then((response) => {
        setAPIdata([response.data]); // Wrap the response data in an array
        // Add logic to handle the search response
      })
      .catch((error) => {
        console.error('Error searching:', error);
        // Add logic to handle search error
      });
  };
  

  const handleAll = () => {
    // Fetch all data from the server and reset the APIdata state
    axios.get('http://localhost:4000/getuser')
      .then((incomingData) => {
        setAPIdata(incomingData.data);
      })
      .catch((error) => {
        console.error('Error fetching all data:', error);
        // Add logic to handle error
      });
  };


  const handleDelete = (e) => {
    e.preventDefault();

    const apiUrl = `http://localhost:4000/deleteuser/${name}`; // Update with your Go server route

    axios
        .post(apiUrl)
        .then((response) => {
          console.log(response);
          // Add logic to handle the response
        });
  };

  useEffect(() => {
    axios.get('http://localhost:4000/getuser')
        .then((incomingData) => {
          setAPIdata(incomingData.data);
        })
  }, [])

  return (
      <Container fluid className="container">
        <Header as="h2">React Google Sheet</Header>
              <div>
              <label>API Mode:</label>
              <select onChange={(e) => setApiMode(e.target.value)} value={apiMode}>
                <option value="REST">REST</option>
                <option value="gRPC">gRPC</option>
              </select>
            </div>
        <Form className="form">
          <Form.Field>
            <label>Name</label>
            <input
                placeholder="Enter your Name"
                onChange={(e) => setName(e.target.value)}
            />
          </Form.Field>
          <Form.Field>
            <label>Age</label>
            <input
                placeholder="How old are you?"
                onChange={(e) => setAge(e.target.value)}
            />
          </Form.Field>
          <Form.Field>
            <label>Commute Method</label>
            <input
                placeholder="How do you commute to college?"
                onChange={(e) => setcommute_method(e.target.value)}
            />
          </Form.Field>
          <Form.Field>
            <label>College</label>
            <input
                placeholder="What college/university do you go to?"
                onChange={(e) => setCollege(e.target.value)}
            />
          </Form.Field>
          <Form.Field>
            <label>Hobbies</label>
            <input
                placeholder="Enter your Hobby"
                onChange={(e) => setHobbies(e.target.value)}
            />
          </Form.Field>

          <Button color="blue" type="submit" onClick={handleSubmit}>
            Submit
          </Button>
          <Button color="blue" type="submit" onClick={handleUpdate}>
            Update
          </Button>
          <Button color="blue" type="submit" onClick={handleDelete}>
            Delete
          </Button>
          <Button color="blue" type="submit" onClick={handleSearch}>
            Search
          </Button>
          <Button color="blue" onClick={handleAll}>
            Show All
          </Button>
        </Form>
        <Table fixed>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Name</Table.HeaderCell>
              <Table.HeaderCell>Age</Table.HeaderCell>
              <Table.HeaderCell>Commute Method</Table.HeaderCell>
              <Table.HeaderCell>College</Table.HeaderCell>
              <Table.HeaderCell>Hobbies</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {APIdata.map((data, index) => (
              <Table.Row key={index}>
                <Table.Cell>{data.name}</Table.Cell>
                <Table.Cell>{data.age}</Table.Cell>
                <Table.Cell>{data.commute_method}</Table.Cell>
                <Table.Cell>{data.college}</Table.Cell>
                <Table.Cell>{data.hobbies}</Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </Container>
  );
}

export default App;
