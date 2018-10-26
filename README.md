# Clarifai Challenge

Can run the application directly. Will open the browser when the rankings are processed.

Can run using Docker.

The API key can be changed from the config file.

### Windows
`run.bat`

### Linux and Mac
`run.sh`

### Docker
`docker build -t clarify-challenge . && docker run -p 8080:8080 clarify-challenge`

## Configuration
API key in `config/config.json`.
