# Use an official Node 14 as a base image
FROM node:14

# Set the working directory in the container
WORKDIR /usr/src/app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of your application's code
COPY . .

# Make port 3001 available to the world outside this container
EXPOSE 3001

# Run server.js when the container launches
CMD ["node", "server.js"]