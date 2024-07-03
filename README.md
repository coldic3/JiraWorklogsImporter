# Jira Worklogs Importer

Use this simple CLI app to import your work logs into the [Worklogs](https://marketplace.atlassian.com/apps/1219004/worklogs-time-tracking-and-reports) Jira app.

For now, you can import your work logs only from Toggl or from CSV file.

The import file supported is an export file from a [Toggl](https://toggl.com) app. See [toggl_sample_export.csv](toggl_sample_export.csv) file.

## Setup

### Basic setup for a single project

1. Create the `.env` file:
```bash
cp .env.dist .env
```

1. Open the `.env` file in a text editor and replace the placeholders with your actual data.

Env variables prefixed with `ATLASSIAN_` are required:
* `ATLASSIAN_DOMAIN`: Your Atlassian domain.
* `ATLASSIAN_EMAIL`: Your email address used for Atlassian access.
* `ATLASSIAN_API_TOKEN`: Your API token for Atlassian access.

Env variables prefixed with `TOGGL_` are optional if you use an `--import` option only:
* `TOGGL_API_TOKEN`: Your API token for Toggl access.
* `TOGGL_USER_ID`: Your user ID in Toggl.
* `TOGGL_CLIENT_ID`: The client ID in Toggl from whom you want to export your work logs.
* `TOGGL_WORKSPACE_ID`: Your workspace ID in Toggl from whom you want to export your work logs.

### Advanced setup for many project

A very common case is that you want to import your worklogs into different Jira accounts from different Toggle accounts
depending on the project. You can handle it with the `--project` option! We will consider two cases.

#### CASE 1: One Toggl account and many Jira accounts

In this case, you most likely work on many Jira projects and want to import all from your one Toggl account where you
log everything. Let's say you work on two projects `sylius` and `symfony`. In your `.env` define common env variables so
all prefixed with `TOGGL_`. In `.env.sylius` and `.env.symfony` specify env variables specific for your Jira account so
prefixed with `ATLASSIAN_`. Now, if you run the command with `--project="symfony"`, files `.env` and `.env.symfony` will
be used.

#### CASE 2: A different Toggl account and Jira account for each project

Similar to the case above but you leave `.env` empty. Instead of that file, you define all env vars right in
`.env.<your-project-name>`.


## Usage

Import from Toggl:
```
./bin/jira-worklogs-importer --since="2024-05-17" --until="2024-05-17"
```

Import from file:
```
./bin/jira-worklogs-importer --import=toggl_sample_export.csv
```

Use `--dry-run` option to see the export before being imported.

By default, it takes description from Toggl / CSV file and combines it into Jira issue key and work log description.
Therefore, your descriptions should follow the pattern [^(.*?)\s*(?:\((.*?)\))?$](https://regex101.com/r/YUvRCq/1). For
example these are valid descriptions: `XYZ-8 (resolving conflicts in the PR)`, `XYZ-8`.
If you are not going to follow this pattern, worry not! You can replace the regex with your own with `DESCRIPTION_REGEX`
env variable.

## Troubleshooting

### Where can I find these user, client and workspace IDs in Toggl?

Well, these IDs are not visible in the Toggl UI however you can easily see it in the address bar in your browser.

1. Open detailed reports in Toggl: https://track.toggl.com/reports/detailed
2. In filters select as follows:
   1. Member => you
   2. Client => the client from whom you want to export you work logs
3. Your URL in the address bar in your browser should be like this: https://track.toggl.com/reports/detailed/{TOGGL_WORKSPACE_ID}/clients/{TOGGL_CLIENT_ID}/period/thisWeek/users/{TOGGL_USER_ID}
