# Gmail CLI Tool

Gmail CLI is a command-line interface tool for managing and reading emails from multiple Gmail accounts using OAuth2 authentication.

## Features

- Manage multiple Gmail accounts
- Read recent emails from configured accounts
- Concurrent email retrieval for efficient operation

## Installation

To install the Gmail CLI tool, make sure you have Go installed on your system (version 1.16 or later), then run:

```bash
go install github.com/leon123858/gmail-cli@latest
```

## Usage

### Configuration

1. Go to the Google Cloud Console
2. Create a new project or select an existing one
3. Enable the Gmail API for your project
4. Create OAuth2 credentials (OAuth client ID For Web)
   - Authorized JavaScript origins: `http://localhost:8080`
   - Authorized redirect URIs: `http://localhost:8080/callback`
5. Instead of downloading the credentials, note down the Client ID and Client Secret
6. Set up your credentials using one of the methods below:

```bash
gmail-cli config set <your-client-id> <your-client-secret>
```

### Adding an account

To add a new Gmail account:

```bash
gmail-cli config add your.email@gmail.com
```

### Deleting an account

To remove a Gmail account:

```bash
gmail-cli config delete your.email@gmail.com
```

### Reading emails

To read emails from all configured accounts:

```bash
gmail-cli run read
```

By default, this will retrieve the 10 most recent emails **today** from each account. You can customize the number of emails to read using the `--count` or `-n` flag:

```bash
gmail-cli run read --count 20
```

## Command Structure

- `gmail-cli`: Root command
    - `config`: Manage configuration
        - `add <email>`: Add a new email account
        - `delete <email>`: Delete an email account
    - `run`: Run Gmail operations
        - `read`: Read emails from configured accounts

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.