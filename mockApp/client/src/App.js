import React, { useState, useEffect } from 'react';
import { Button, Form, Container, Header, Table } from 'semantic-ui-react';

//import userService from './UserServiceClient';

import axios from 'axios';
import './index.css';

function App() {
  const [name, setName] = useState('');
  const [age, setAge] = useState('');
  const [commuteMethod, setcommuteMethod] = useState('');
  const [college, setCollege] = useState('');
  const [hobbies, setHobbies] = useState('');
  const [APIdata, setAPIdata] = useState([]);
  //const [refresh, setRefresh] = useState([]);
  const [apiMode, setApiMode] = useState('REST'); // Default mode is REST

  const { UserServiceClient } = require('./user_grpc_web_pb');
  const { CreateUserRequest, DeleteUserRequest, GetUserRequest, UpdateUserRequest, User } = require('./user_pb.js');
  //var client = new UserServiceClient('http://localhost:9090', null, null);
  var client = new UserServiceClient('https://envoy-proxy-necuf5ddgq-ue.a.run.app', null, null);
  const [user, setUser] = useState(null);
  
  //const handleSubmit = (e) => {
  const handleSubmit = (e = { preventDefault: () => {} }) => {
    e.preventDefault();

    const objt = { name, age, commuteMethod, college, hobbies };
    //const apiUrl = apiMode === 'REST' ? 'http://localhost:4000/adduser' : 'https://sheet.best/api/sheets/ab206fd6-ee69-4b1d-a56f-f29c0ba70176';
    const apiUrl = apiMode === 'REST' ? 'https://rest-envoy-necuf5ddgq-ue.a.run.app/adduser' : 'https://sheet.best/api/sheets/ab206fd6-ee69-4b1d-a56f-f29c0ba70176';
    //https://rest-apigo-main-6j7fqbeloq-ue.a.run.app/

    axios.post(apiUrl, objt).then((response) => {
      console.log(response);
      // Add logic to handle the response, e.g., update state or show a message
    });
  };

  //const handleUpdate = (e) => {
  const handleUpdate = (e = { preventDefault: () => {} }) => {
    e.preventDefault();

    const objt = { name, age, commuteMethod, college, hobbies };
    //`http://localhost:4000/updateuser/${name}`, // Update with your Go server route
    
    axios
        .put(
            `https://rest-envoy-necuf5ddgq-ue.a.run.app/updateuser/${name}`,
            objt
        )
        .then((response) => {
          console.log(response);
          // Add logic to handle the response
        });
  };

  //const handleSearch = (e) => {
  const handleSearch = (e = { preventDefault: () => {} }) => {
    e.preventDefault();
  
    // Perform the search using the name in the state
    //axios.get(`http://localhost:4000/getuser/${name}`)
    axios.get(`https://rest-envoy-necuf5ddgq-ue.a.run.app/getuser/${name}`)
      .then((response) => {
        console.log(response.data);
        setAPIdata([response.data]); // Wrap the response data in an array
        // Add logic to handle the search response
      })
      .catch((error) => {
        console.error('Error searching:', error);
        // Add logic to handle search error
      });
  };
  

  const handleAll  = (e = { preventDefault: () => {} }) => {
    // Fetch all data from the server and reset the APIdata state
    //axios.get('http://localhost:4000/getuser')
    axios.get('https://rest-envoy-necuf5ddgq-ue.a.run.app/getuser')
      .then((incomingData) => {
        if (Array.isArray(incomingData.data) && incomingData.data.length > 0) {
          // Exclude the first row (header) and update the state
          setAPIdata(incomingData.data.slice(1));
        } else {
          // Handle the case where data is not as expected
          console.error('Received data is not in expected array format');
        }
      })
      .catch((error) => {
        console.error('Error fetching all data:', error);
        // Add logic to handle error
      });
  };


  const handleDelete = (e = { preventDefault: () => {} }) => {
    e.preventDefault();

    const apiUrl = `https://rest-apigo-main-6j7fqbeloq-ue.a.run.app/deleteuser/${name}`; // Update with your Go server route
    axios
        .post(apiUrl)
        .then((response) => {
          console.log(response);
          // Add logic to handle the response
        });
  };

  const handleGetUser = () => {
    let request = new GetUserRequest();
    request.setName(name);

    client.getUser(request, {}, (err, response) => {
      if (err) {
        console.error(err);
        return;
      }
      //setUser(response.toObject());
      //setAPIdata([response.data]);
      setAPIdata([response.toObject()]);
      console.log(request);
    });
  };

  const handleCreateUser = () => {
    // Create a User message and set its fields
    let userMessage = new User();
    userMessage.setName(name);
    userMessage.setAge(age);  // Assuming User has a setAge method
    userMessage.setCommutemethod(commuteMethod); // and so on for other fields...
    userMessage.setCollege(college);
    userMessage.setHobbies(hobbies);
  
    // Create a CreateUserRequest and set its user message
    let request = new CreateUserRequest();
    request.setUser(userMessage);
  
    client.createUser(request, {}, (err, response) => {
      if (err) {
        console.error(err);
        alert('Failed to create user: ' + err.message);
        return;
      }
      alert('User successfully created!');
      // Reset the form fields after successful creation
      //setName('');
      //setAge('');
      //setcommuteMethod('');
      //setCollege('');
      //setHobbies('');
    });
  };

  const handleUpdateUser = () => {
    let request = new UpdateUserRequest();
    
    // Create a User message and set its fields
    let userMessage = new User();
    userMessage.setName(name)
    userMessage.setAge(age); // Assuming User has a setAge method
    userMessage.setCommutemethod(commuteMethod); // and so on for other fields...
    userMessage.setCollege(college);
    userMessage.setHobbies(hobbies);
    
    // Set the name and user message in UpdateUserRequest
    request.setName(name);
    request.setUser(userMessage);
    
    client.updateUser(request, {}, (err, response) => {
      if (err) {
        console.error(err);
        alert('Failed to update user: ' + err.message);
        return;
      }
      alert('User successfully updated!');
    });
  };

  const handleDeleteUser = () => {
    let request = new DeleteUserRequest();
    request.setName(name);

    client.deleteUser(request, {}, (err, response) => {
        if (err) {
            console.error(err);
            alert('Failed to delete user: ' + err.message);
            return;
        }
        alert('User successfully deleted!');

        // Clear the displayed user details
        setUser(null);

        // Clear the input field, if needed
        setName('');

        
    });
  };

  const handleSubmitChoice = (e) => {
    if (apiMode === 'REST') {
      handleSubmit(e);
    } else {
      handleCreateUser(e);
    }
  };

  const handleUpdateChoice = (e) => {
    if (apiMode === 'REST') {
      handleUpdate(e);
    } else {
      handleUpdateUser(e);
    }
  };

  const handleDeleteChoice = (e) => {
    if (apiMode === 'REST') {
      handleDelete(e);
    } else {
      handleDeleteUser(e);
    }
  };

  const handleSearchChoice = (e) => {
    if (apiMode === 'REST') {
      handleSearch(e);
    } else {
      handleGetUser(e);
    }
  };


  useEffect(() => {
    //axios.get('http://localhost:4000/getuser')
    axios.get('https://rest-apigo-main-6j7fqbeloq-ue.a.run.app/getuser')
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
                onChange={(e) => setcommuteMethod(e.target.value)}
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

          <Button color="blue" type="submit" onClick={handleSubmitChoice}>
            Submit
          </Button>
          <Button color="blue" type="submit" onClick={handleUpdateChoice}>
            Update
          </Button>
          <Button color="blue" type="submit" onClick={handleDeleteChoice}>
            Delete
          </Button>
          <Button color="blue" type="submit" onClick={handleSearchChoice}>
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
                <Table.Cell>{data.commuteMethod || data.commutemethod}</Table.Cell>
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

