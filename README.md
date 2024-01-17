# Matrix GPT

Matrix GPT is a Matrix chatbot that uses OpenAI for real-time chatting.

![](./.github/img.png)
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

This bot supports the following commands:

- `!image[-natural/-vivid]`: This command will create and return an image based on the text you provide. The default style is "Natural".
- `!reset [text]`: This command will reset the user's history. If you provide text after the `!reset` command, the bot generates a response using GPT, based on this input text.
- `[text]`: If you simply input text without any specific command, the bot will automatically generate a GPT-based response related to the text provided.

### Additional Notes

- You can use short aliases for a command; for example, `!i` for `!image`, or `!iv` for `!image-vivid`.
- If you need to stop any ongoing processing, you can just delete your message from the chat`.
- In case of errors, the bot reacts with a ❌. If you notice this, please check logs.
