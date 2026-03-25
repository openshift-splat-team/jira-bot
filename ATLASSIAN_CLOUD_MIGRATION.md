# Atlassian Cloud Migration Guide

This guide will help you migrate jira-bot from self-hosted Jira (issues.redhat.com) to Atlassian Cloud.

## Changes Made

The jira-bot has been updated to automatically detect and support both:
- **Self-hosted Jira** (e.g., issues.redhat.com) - Uses Bearer token authentication
- **Atlassian Cloud** (e.g., yourcompany.atlassian.net) - Uses Basic Auth with email + API token

## Key Differences

### Authentication

**On-Premise (self-hosted):**
- Uses Bearer token authentication
- Environment variable: `JIRA_PERSONAL_ACCESS_TOKEN`
- Example: `export JIRA_PERSONAL_ACCESS_TOKEN=your_bearer_token`

**Atlassian Cloud:**
- Uses Basic Auth with email + API token
- Environment variables: `JIRA_EMAIL` and `JIRA_API_TOKEN`
- Example:
  ```bash
  export JIRA_EMAIL=your.email@company.com
  export JIRA_API_TOKEN=your_cloud_api_token
  ```

### Custom Field IDs

Custom field IDs are different between on-premise and Cloud instances:

| Field | On-Premise ID | Cloud ID (may vary) |
|-------|---------------|---------------------|
| Story Points | customfield_12310243 | customfield_10028 |
| Status Summary | customfield_12320841 | customfield_10841 |
| Epic Link | customfield_12311140 | customfield_10014 |
| Feature Link | customfield_12318341 | customfield_10341 |

**Note:** Cloud custom field IDs may vary by organization. If the default Cloud field IDs don't work, you can find the correct IDs by:
1. Going to an issue in your browser
2. Appending `/rest/api/3/issue/KEY-123` to the URL
3. Searching for the field names in the JSON response

Update the constants in `pkg/util/jira.go` if needed.

## Migration Steps

### 1. Generate Atlassian Cloud API Token

1. Log in to your Atlassian Cloud account
2. Visit: https://id.atlassian.com/manage-profile/security/api-tokens
3. Click "Create API token"
4. Give it a name (e.g., "jira-bot")
5. Copy the generated token - you won't be able to see it again!

### 2. Update Environment Variables

**For Atlassian Cloud:**
```bash
export JIRA_BASE_URL=https://yourcompany.atlassian.net
export JIRA_EMAIL=your.email@company.com
export JIRA_API_TOKEN=your_atlassian_cloud_api_token_here
export JIRA_PROJECT=SPLAT
export JIRA_BOARD="SPLAT - Scrum Board"
```

**For on-premise (existing setup):**
```bash
export JIRA_BASE_URL=https://issues.redhat.com  # Optional, this is the default
export JIRA_PERSONAL_ACCESS_TOKEN=your_bearer_token
export JIRA_PROJECT=SPLAT
export JIRA_BOARD="SPLAT - Scrum Board"
```

### 3. Test the Connection

Run a simple command to verify connectivity:

```bash
./jira-bot issue --help
```

The bot will log which authentication method and field IDs it's using:
- `Connecting to Atlassian Cloud: ...` indicates Cloud mode
- `Connecting to on-premise Jira: ...` indicates on-premise mode

### 4. Verify Custom Field IDs (Cloud only)

If you're using Atlassian Cloud and encounter errors related to custom fields, you may need to verify the field IDs:

```bash
# Use curl or browser to fetch an issue and examine the fields
curl -u your.email@company.com:your_api_token \
  https://yourcompany.atlassian.net/rest/api/3/issue/SPLAT-123
```

Look for fields like "Story Points", "Epic Link", etc. and note their custom field IDs. If they differ from the defaults in `pkg/util/jira.go`, update the constants:

```go
// Custom field IDs for Atlassian Cloud
const (
	CloudFieldStoryPoints   = "customfield_XXXXX"  // Update with your field ID
	CloudFieldStatusSummary = "customfield_XXXXX"  // Update with your field ID
	CloudFieldEpicLink      = "customfield_XXXXX"  // Update with your field ID
	CloudFieldFeatureLink   = "customfield_XXXXX"  // Update with your field ID
)
```

Then rebuild the bot:
```bash
./hack/build.sh
```

## Automatic Detection

The jira-bot automatically detects Cloud vs on-premise by checking if `.atlassian.net` is in the base URL. This means you can switch between environments just by changing the `JIRA_BASE_URL` environment variable.

## Troubleshooting

### "Atlassian Cloud requires JIRA_EMAIL and JIRA_API_TOKEN environment variables"

- Verify you've set both `JIRA_EMAIL` and `JIRA_API_TOKEN`
- Ensure `JIRA_BASE_URL` contains `.atlassian.net`

### "On-premise Jira requires JIRA_PERSONAL_ACCESS_TOKEN environment variable"

- Set the `JIRA_PERSONAL_ACCESS_TOKEN` environment variable
- This happens when connecting to non-Cloud Jira instances

### "unable to execute query" or 401 errors

**Cloud:**
- Verify your API token is correct (regenerate if needed)
- Check that `JIRA_EMAIL` matches your Atlassian account

**On-premise:**
- Verify your bearer token is still valid
- Check that you have access to the Jira instance

### "unable to get project" or 403 errors

- You may not have permissions to view the project
- Contact your Jira admin to grant access

### Custom field errors (nil pointer, field not found)

- The custom field IDs may be different in your Cloud instance
- Follow step 4 above to verify and update the field IDs

## Rolling Back

If you need to temporarily switch back to on-premise:

1. Update environment variables:
   ```bash
   export JIRA_BASE_URL=https://issues.redhat.com
   export JIRA_PERSONAL_ACCESS_TOKEN=your_bearer_token
   # Unset Cloud variables
   unset JIRA_EMAIL
   unset JIRA_API_TOKEN
   ```

2. Re-run the bot - it will automatically use on-premise mode

## Additional Resources

- [Atlassian Cloud API Documentation](https://developer.atlassian.com/cloud/jira/platform/rest/v3/)
- [Manage API tokens](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/)
- [Jira REST API v3 Reference](https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/)
