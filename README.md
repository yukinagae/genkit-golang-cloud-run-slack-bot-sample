# genkit-golang-cloud-run-slack-bot-sample

`genkit-golang-cloud-run-slack-bot-sample` is a sample repository for building a ChatGPT Slack bot using Firebase Genkit with Golang, deployed on Google Cloud Run.

- [Requirements](#requirements)
- [Usage](#usage)
- [License](#license)

## Requirements

- **Go**: Follow the [Go - Download and install](https://go.dev/doc/install) to install Go.
- **Genkit**: Follow the [Firebase Genkit - Get started](https://firebase.google.com/docs/genkit/get-started) to install Genkit.
- **Google Cloud CLI (gcloud)**: Follow the [Google Cloud - Install the gcloud CLI](https://cloud.google.com/sdk/docs/install) to install gcloud.
- **ngrok**: Follow the [ngrok - Quickstart](https://ngrok.com/docs/getting-started/) to install ngrok.
- **golangci-lint**: Follow the [golangci-lint - Install](https://golangci-lint.run/welcome/install/) to install golangci-lint.

Verify your installations:

```bash
$ go version
v22.4.1
$ genkit --version
0.5.4
$ gcloud --version
Google Cloud SDK 489.0.0
alpha 2024.08.16
bq 2.1.8
core 2024.08.16
gcloud-crc32c 1.0.0
gsutil 5.30
$ ngrok --version
ngrok version 3.3.0
$ golangci-lint --version
golangci-lint has version 1.60.3 built with go1.23.0 from c2e095c on 2024-08-22T21:45:24Z
```

## Usage

### Run Genkit

Set your API key and start Genkit:

```bash
$ export GOOGLE_GENAI_API_KEY=your_api_key
$ make genkit # Starts Genkit
```

Open your browser and navigate to [http://localhost:4000](http://localhost:4000) to access the Genkit UI.

### Setup Your Slack App

1. Navigate to [Slack - Your Apps](https://api.slack.com/apps) and click `Create New App`.
2. Choose `From an app manifest` option, select a workspace under `Pick a workspace to develop your app`, and then click `Next`.
3. In the app manifest JSON below, replace `[your_app_name]` with your app's name, paste the updated JSON, then proceed by clicking `Next` and `Create`.

```json
{
  "display_information": {
    "name": "[your_app_name]"
  },
  "features": {
    "bot_user": {
      "display_name": "[your_app_name]",
      "always_online": true
    }
  },
  "oauth_config": {
    "scopes": {
      "bot": [
        "app_mentions:read",
        "channels:history",
        "chat:write",
        "files:read"
      ]
    }
  },
  "settings": {
    "event_subscriptions": {
      "request_url": "http://dummy/events",
      "bot_events": ["app_mention"]
    },
    "org_deploy_enabled": false,
    "socket_mode_enabled": false,
    "token_rotation_enabled": false
  }
}
```

4. Navigate to `Settings` and select `Install App`, then click `Install to Workspace` and `Allow` button.
5. Find your `Bot User OAuth Token` under `OAuth & Permissions`. This is your `SLACK_BOT_TOKEN` for later use.
6. Find your `Signing Secret` under `Basic Information`. This is your `SLACK_SIGNING_SECRET` for later use.
7. To add your bot to a Slack channel, use the command:

```bash
/invite @[your_app_name]
```

### Run HTTP Server Locally

To start the local http server, run the following command:

```bash
$ export GOOGLE_GENAI_API_KEY=your_api_key
$ export SLACK_BOT_TOKEN=your_bot_token
$ export SLACK_SIGNING_SECRET=your_signing_secret
$ make dev
```

To make your local http server accessible online, use ngrok to forward port `3000`:

```bash
$ ngrok http 3000
Forwarding https://[your_ngrok_id].ngrok-free.app -> http://localhost:3000
```

This command provides a public URL. Replace [your_ngrok_id] in the URL `https://[your_ngrok_id].ngrok-free.app` with the ID provided by ngrok.

To configure Slack event subscriptions:

1. Go to the `Event Subscriptions` page on your Slack app's dashboard.
2. In the `Request URL` field, enter `https://[your_ngrok_id].ngrok-free.app/slack/events`.
3. Wait for the `Request URL Verified` confirmation, then click the `Save changes` button.

To test in a Slack channel, mention your bot using `@[your_app_name]` followed by a URL, like so:

```bash
@[your_app_name] hello
```

If everything is set up correctly, you should see a response like this:

```text
Hi there! How can I help you today?
```

### Deploy

Set your secret values in `./.env.yaml`:

```bash
$ cp -p ./.env.example.yaml ./.env.yaml
$ vim ./.env.yaml # replace the secrets with your own values
GOOGLE_GENAI_API_KEY: your_api_key
SLACK_BOT_TOKEN: your_bot_token
SLACK_SIGNING_SECRET: your_signing_secret
```

Follow these steps to deploy the application:

```bash
$ gcloud auth application-default login
$ gcloud config set core/project [your-project-id]
$ make deploy
```

**CAUTION**: This deployment uses `.env.yaml` for environment variables, including the API key. This is not recommended for production. Instead, use Google Cloud Secret Manager for better security.

The final step involves linking your deployed application to the Slack app for integration.

To configure Slack event subscriptions:

1. Go to the `Event Subscriptions` page on your Slack app's dashboard.
2. In the `Request URL` field, enter `https://slack-bot-application-[your-cloud-run-id]-uc.a.run.app/slack/events`.
3. Wait for the `Request URL Verified` confirmation, then click the `Save changes` button.

NOTE: Replace `[your-cloud-run-id]` with your Cloud Run service URL value, found in the Cloud Run Console.

To test in a Slack channel, mention your bot using `@[your_app_name]` followed by a URL, like so:

```bash
@[your_app_name] hello
```

### Code Formatting

To ensure your code is properly formatted, run the following command:

```bash
$ make tidy
```

## License

MIT
