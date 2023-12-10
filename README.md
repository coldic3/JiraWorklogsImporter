# Jira Worklogs Importer

Use this simple CLI app to import your work logs into the [Worklogs](https://marketplace.atlassian.com/apps/1219004/worklogs-time-tracking-and-reports) Jira app.

For now, the only import file supported is an export file from a [Toggl](https://toggl.com) app. See [toggl_sample_export.csv](toggl_sample_export.csv) file.

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
go run main.go --import=toggl_sample_export.csv
```
