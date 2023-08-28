package bot

const (
	helpMsg = `**Commands**
- *!image [text]*: Creates an image based on the provided text.
- *!reset [text]*: Resets the user history. If a text is provided after the reset command, it will generate a GPT response based on this text.
- *[text]*: If only text is provided, the bot will generate a GPT-based response related to that text.

**Notes**
- You can use the first letter of a command as an alias. For example, "!i" for "!image".
- If you wish to terminate the current processing, simply delete your message from the chat.
`
	timeoutMsg        = "Timeout error. Please try again. If issue persists, contact the administrator."
	unknownCommandMsg = "Unknown command. Please use the `!help` command to access the available commands"
)
