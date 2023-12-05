# Jira Worklogs Importer

Use this simple CLI app to import your work logs into the [Worklogs](https://marketplace.atlassian.com/apps/1219004/worklogs-time-tracking-and-reports) Jira app.

## Setup

1. Create the `.env` file:
```bash
cp .env.dist .env
```

1. Open the `.env` file in a text editor and replace the placeholders with your actual data:
* `ATLASSIAN_DOMAIN`: Your Atlassian domain.
* `EMAIL`: Your email address used for Atlassian access.
* `API_TOKEN`: Your API token for Atlassian access.

## Usage

Run:
```bash
go run main.go
```
