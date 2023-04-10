# Minecraft Players Watcher

Minecraft Players Watcher is a Go application that monitors a Minecraft server log file and extracts player usernames. The usernames are stored in a Redis database and exposed through a simple web API.

## Features

- Monitors the Minecraft `latest.log` file for changes.
- Extracts player usernames from log entries.
- Stores unique player usernames in a Redis database.
- Provides a web API endpoint to retrieve the list of player usernames.

## Prerequisites

- Docker and Docker Compose installed on your system.
- Minecraft server log file (`latest.log`) accessible from the host machine.

## Usage

### Configuration

Configuration for the application is provided through environment variables. The following environment variables are required:

- `SERVER_PORT`: The port on which the web API will listen (e.g., 8080).
- `LOG_FILE`: The path to the Minecraft `latest.log` file inside the Docker container (e.g., /app/logs/latest.log).

### Building and Running

1. Clone this repository to your local machine.

git clone https://github.com/your-github-username/minecraft-players-watcher.git
cd minecraft-players-watcher


2. Update the `docker-compose.yml` file with the correct absolute path to the `latest.log` file on your host machine.

3. Build and run the Docker Compose setup.

```bash
docker-compose up --build
```


4. Access the web API endpoint to retrieve the list of player usernames.

http://localhost:SERVER_PORT/api/users



Replace `SERVER_PORT` with the port number you set in the configuration.

## API

### Get Player Usernames

Returns the list of player usernames.

- URL: `/api/users`
- Method: `GET`
- Response: JSON array containing player usernames.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)

