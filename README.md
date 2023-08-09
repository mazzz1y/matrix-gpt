# Matrix GPT

Matrix GPT is a Matrix chatbot that uses OpenAI for real-time chatting.

![](./.github/img.png)
## Installation

### Docker

Run the Docker container:

```bash
docker run -d --name matrix-gpt \
  -e MATRIX_PASSWORD="matrix password" \
  -e MATRIX_ID="matrix id" \
  -e MATRIX_URL="matrix server url" \
  -e OPENAI_TOKEN="openai token" \
  -e SQLITE_PATH="persistent path for sqlite database"
  -e USER_IDS="allowed user ids"
  ghcr.io/mazzz1y/matrix-gpt:latest

```
## Configuration

You can configure GPT Matrix using the following environment variables:

- `SERVER_URL`: The URL to the Matrix homeserver.
- `USER_ID`: Your Matrix user ID for the bot.
- `PASSWORD`: The password for your Matrix bot's account.
- `SQLITE_PATH`: Path to SQLite database for end-to-end encryption.
- `HISTORY_EXPIRE`: Duration after which chat history expires.
- `GPT_MODEL`: The OpenAI GPT model being used.
- `GPT_HISTORY_LIMIT`: Limit for number of chat messages retained in history.
- `GPT_TIMEOUT`: Duration for OpenAI API timeout.
- `GPT_MAX_ATTEMPTS`: Maximum number of attempts for GPT API retries.
- `GPT_USER_IDS`: List of authorized user IDs for the bot.

Alternatively, you can set these options using command-line flags. Run `./matrix-gpt --help` for more
information.

## Usage

1. Begin by adding the bot to your contact list.
   Once added, you can start interacting with it. Simply send your questions or commands, and the bot will respond.
2. If at any point you wish to reset the context of the conversation, use the `!reset` command.
   Send `!reset` as a message to the bot, and it will clear the existing context, allowing you to start a fresh
   conversation.

## TODO

* Add image generation using DALL·E model