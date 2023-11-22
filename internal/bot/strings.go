package bot

const (
	helpMsg = `**Commands**
- ` + "`!image[-natural/-vivid] [prompt]`" + `: Creates an image based on the provided prompt. The default style is "Natural".
- ` + "`!reset [prompt]`" + `: Resets the user's history. If a prompt is provided after the reset command, the bot will generate a GPT response based on this prompt.
- ` + "`[prompt]`" + `: If only a prompt is provided, the bot will generate a GPT-based response related to that prompt.

**Notes**
- You can use short aliases for a command; for example, ` + "`!i`" + ` for ` + "`!image`" + `, or ` + "`!iv`" + ` for ` + "`!image-vivid`" + `.
- To terminate the current processing, simply delete your message from the chat.
- If there are any errors, the bot will respond with a ❌ reaction. Contact the administrator if this occurs.
`
	timeoutMsg        = "Timeout error. Please try again. If the issue persists, contact the administrator."
	unknownCommandMsg = "Unknown command. Please use the `!help` command to access the available commands."
)
