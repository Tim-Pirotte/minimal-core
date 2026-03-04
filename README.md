# Minimal Core
The base Minimal compiler on which the Minimal language extensions are added.

## Project Structure
- **cmd:** Applications used for testing
- **src:** The root of all layers and extensions
  - **startup:** Cli system to select an app
    - **commands:** Apps for startup
      - **templates:** App to create projects from templates
        - **stores:** Places to save templates
          - **application-data:** Uses the app data directories to save templates
          - **directory:** Uses a directory at the location of the executable for templates
          - **embedded:** Templates baked into the executable
  - **config:** A config loader
  - **internal-logging:** Logging for apps that shouldn't be shown to users by default
  - **user-messaging:** Sends messages that should be shown to the user via output channels
    - **outputs:** The different channels to show the messages
      - **log-renderer:** Nicely rendered messages to display in the terminal
  - **ui:** UI layer via which can be communicated with the user
    - **user-interfaces:** UI implementations
      - **cli:** UI via CLI
      - **tui:** UI via TUI *(Not meant to be used!)*
  - **tokenizer:** Layer to convert text into a series of tokens
    - **matchers:** Extensions that match on textual input
      - **identifiers:** Match text that doesn't start with a digit
      - **indentation:** Keep track of indents
      - **numbers:** Match decimal, hex, octal and binary numbers
      - **string-literals:** Match strings with interpolation and multiline support
      - **symbols:** A trie to match exact strings
      - **white-space:** Ignore white space
  
*The docs of the language itself can be found at: TODO*
