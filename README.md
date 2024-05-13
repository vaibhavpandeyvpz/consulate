<h1 align="center">consulate</h1>

[![Screenshot](https://raw.githubusercontent.com/vaibhavpandeyvpz/consulate/main/screenshot.png)](https://raw.githubusercontent.com/vaibhavpandeyvpz/consulate/main/screenshot.png)

<p align="center">
Headless enquiry management system in <a href="https://go.dev/">Go</a> built on <a href="https://slack.com/intl/en-in/">Slack</a> and <a href="https://exotel.com/">Exotel</a> for inbound marketing teams.
Record collected enquiries, place and record outbound calls, store follow-ups etc.
</p>

## Development

Make sure you have [Docker](https://www.docker.com/) installed on your workstation.
For the IDE, I highly recommend using [GoLand](https://www.jetbrains.com/go/) i.e., my go to choice for [Go](https://go.dev) development.

Firstly, download or clone the project using [Git](https://git-scm.com/) and then run following commands in project folder:

```shell
# create .env file in project
cp .env.dist .env

# update NGROK_AUTHTOKEN in .env

# create app config file
cp config.dist.yml config.yml

# update values in config.yml

# create ngrok config file
cp ngrok.dist.yml ngrok.yml

# update ngrok domain in ngrok.yml

# create Slack app manifest
cp slack.dist.yml slack.yml

# update ngrok domain in slack.yml

# start all services
docker compose up -d
```

Go to [api.slack.com](https://api.slack.com/), create a new app using provided manifest and install it on a [Slack](https://slack.com/intl/en-in/) workspace.
Once done, get the **signing secret** as well as **bot access token** to update in `.env` file and restart the services using below command:

```shell
# stop running services
docker compose down

# restart services
docker compose up -d
```

## Deployment

For deployment, using [Docker](https://www.docker.com/) is the easiest way to go.
There's a bundled `Dockerfile` that builds and exposes the server on port `8080` (can be configured using `PORT` environment variable).

To build the container for deployment, use below command:

```shell
docker build -t consulate .
```
