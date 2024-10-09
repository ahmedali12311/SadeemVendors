# Use an official Node runtime as a parent image
FROM node:14

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the React app
RUN npm run build

# Install a web server to serve the built app
RUN npm install -g serve

# Define environment variable
ENV PORT=8080

# Expose the port the app runs on
EXPOSE 8080

# Command to run the app
CMD ["serve", "-s", "build", "-l", "8080"]
