# Jira Worklogs Importer

Use this CLI app to import your work logs into the [Worklogs](https://marketplace.atlassian.com/apps/1219004/worklogs-time-tracking-and-reports) Jira app.

Toggl, Clockify and CSV files are supported!

The CSV file supported is an export file from a [Toggl](https://toggl.com) app. See [toggl_sample_export.csv](toggl_sample_export.csv) file.

## Setup

### Basic setup for a single project

1. Create the `.env` file:
```bash
cp .env.dist .env
```

1. Open the `.env` file in a text editor and replace the placeholders with your actual data.

* `IMPORT_STRATEGY`: Specifies the import strategy to use. It can be `toggl_to_jira`, `clockify_to_jira`, or `csv_to_jira`. Default is `csv_to_jira`.

Environment variables prefixed with `ATLASSIAN_` are always required:
* `ATLASSIAN_DOMAIN`: Your Atlassian domain.
* `ATLASSIAN_EMAIL`: The email address associated with your Atlassian account.
* `ATLASSIAN_API_TOKEN`: Your Atlassian API token.

Environment variables prefixed with `TOGGL_` are required if you use the `toggl_to_jira` import strategy:
* `TOGGL_API_TOKEN`: Your Toggl API token.
* `TOGGL_USER_ID`: Your Toggl user ID.
* `TOGGL_CLIENT_ID`: The client ID from which to export work logs.
* `TOGGL_WORKSPACE_ID`: Your Toggl workspace ID.

Environment variables prefixed with `CLOCKIFY_` are required if you use the `clockify_to_jira` import strategy:
* `CLOCKIFY_API_TOKEN`: Your Clockify API token.
* `CLOCKIFY_USER_ID`: Your Clockify user ID.
* `CLOCKIFY_PROJECT_ID`: The project ID from which to export work logs.
* `CLOCKIFY_WORKSPACE_ID`: Your Clockify workspace ID.

### Advanced setup for many projects

A very common case is that you want to import your worklogs into different Jira accounts from different Toggle accounts
depending on the project. You can handle it with the `--project` option! We will consider two example scenarios.

#### EXAMPLE 1: One Toggl account and many Jira accounts

In this case, you most likely work on many Jira projects and want to import all from your one Toggl account where you
log everything. Let's say you work on two projects `sylius` and `symfony`. In your `.env` define common env variables so
all prefixed with `TOGGL_`. In `.env.sylius` and `.env.symfony` specify env variables specific for your Jira account so
all prefixed with `ATLASSIAN_`. Now, if you run the command with `--project="symfony"`, files `.env` and `.env.symfony`
will be used.

#### EXAMPLE 2: A different Toggl account and Jira account for each project

Similar to the case above but you leave `.env` empty. Instead of that file, you define all env vars right in
`.env.<your-project-name>`.


## Usage

Import from Toggl/Clockify:
```
./bin/jira-worklogs-importer --since="2024-05-17" --until="2024-05-17"
```

Import from file:
```
./bin/jira-worklogs-importer --import=toggl_sample_export.csv
```

By default, it takes a description from Toggl / Clockify / CSV file and combines it into Jira issue key and work log
description.\
Therefore, your descriptions should follow the pattern [^(.*?)\s*(?:\((.*?)\))?$](https://regex101.com/r/YUvRCq/1). For example these are valid
descriptions: `XYZ-8 (resolving conflicts in the PR)`, `XYZ-8`.\
If you are not going to follow this pattern, worry not! You can replace the regex with your own with `DESCRIPTION_REGEX`
env variable.

## Troubleshooting

### Where can I find these user, client and workspace IDs in Toggl?

Well, these IDs are not visible in the Toggl UI however you can easily see it in the address bar in your browser.

1. Open detailed reports in Toggl: https://track.toggl.com/reports/detailed
2. In filters select as follows:
   1. Member => you
   2. Client => the client from whom you want to export you worklogs
3. Your URL in the address bar in your browser should be like this: https://track.toggl.com/reports/detailed/{TOGGL_WORKSPACE_ID}/clients/{TOGGL_CLIENT_ID}/period/thisWeek/users/{TOGGL_USER_ID}
