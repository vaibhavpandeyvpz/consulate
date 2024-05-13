<h1 align="center">consulate</h1>

[![Screenshot](https://raw.githubusercontent.com/vaibhavpandeyvpz/consulate/main/screenshot.png)](https://raw.githubusercontent.com/vaibhavpandeyvpz/consulate/main/screenshot.png)

<p align="center">
Headless enquiry management system in <a href="https://go.dev/">Go</a> built on <a href="https://slack.com/intl/en-in/">Slack</a> and <a href="https://exotel.com/">Exotel</a> for inbound marketing teams.
Record collected enquiries, place and record outbound calls, store follow-ups etc.
</p>

## Usage

Firstly, go to [api.slack.com](https://api.slack.com/), create a new app using provided manifest (see [slack.dist.yml](slack.dist.yml)) and install it on a [Slack](https://slack.com/intl/en-in/) workspace.
Once done, get the **signing secret** as well as **bot access token** to update in `config.yml` file.

Grab a binary from the latest release for your platform from [this page](https://github.com/vaibhavpandeyvpz/consulate/releases/latest).
In the same folder as binary, create a `config.yml` file from the sample in the repository using below command:

```shell
wget -O config.yml https://raw.githubusercontent.com/vaibhavpandeyvpz/consulate/main/config.dist.yml
```

Update your [Slack](https://slack.com/intl/en-in/) and [Exotel](https://exotel.com/) credentials in `config.yml` file and start the app server using below command:

```shell
./consulate -config=config.yml
```

Since [Slack](https://slack.com/intl/en-in/) needs to communicate to your app for certain functionality, its recommended to run this on a server and install an [SSL](https://letsencrypt.org/) certificate.

## Development

Make sure you have [Docker](https://www.docker.com/) installed on your workstation.
For the IDE, I highly recommend using [GoLand](https://www.jetbrains.com/go/) i.e., my go to choice for [Go](https://go.dev) development.

Download or clone the project using [Git](https://git-scm.com/) and then run following commands in project folder:

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

## Deployment

For deployment, using pre-built binary from [Releases](https://github.com/vaibhavpandeyvpz/consulate/releases) section is the easiest way to go.

You could also use [Docker](https://www.docker.com/) for deployment. There's a bundled `Dockerfile` that builds and exposes the server on port `8080` (can be configured using `PORT` environment variable).

To build the [Docker](https://www.docker.com/) container locally, use below command:

```shell
docker build -t consulate .
```
