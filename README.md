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

By default, this will retrieve the 10 most recent emails from each account. You can customize the number of emails to read using the `--count` or `-n` flag:

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