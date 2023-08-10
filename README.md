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

Follow these simple steps to interact with the bot:

1. **Add the bot to your contact list.**

   After the bot has been added, begin your interaction by sending your questions, commands or comments.

2. **Send commands.**

   You can use the following commands to communicate with the bot:

   - **Generate an Image:** `!image [text]` - This command creates an image based on the provided text.
   - **Reset User History:** `!reset [text]` - This command resets the user's command history. If text is provided following the reset command, the bot will generate a GPT-based response based on this text.
   - **Send a Text Message:** `[text]` - Send any text to the bot and it will generate a GPT-based response relevant to your text.

3. **Identify error responses.**

   If there are any errors in processing your requests or commands, the bot will respond with a ❌ reaction.