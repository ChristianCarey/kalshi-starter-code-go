# kalshi-starter-code-go

Example Go code for accessing API-authenticated endpoints on Kalshi. This is not an SDK.

## Overview

This repository a basic example for interacting with Kalshi's API endpoints using Go. It demonstrates authentication, making requests, and handling responses in a simple, straightforward manner.

## Prerequisites

- Go 1.19 or later
- A Kalshi account with API credentials

## Installation

Clone the repository:

```bash
git clone https://github.com/ChristianCarey/kalshi-starter-code-go.git
cd kalshi-starter-code-go
```

Install dependencies:

```bash
go mod download
```

## Configuration

Create a `.env` file in the project root with your Kalshi API credentials (see `.env.example` for template):

## Usage

Run the example code:

```bash
go run main.go
```

## Security Notes

- Never commit your private keys or `.env` file containing API credentials
- The `.gitignore` file is configured to exclude sensitive files
- Store your private keys securely and use appropriate file permissions

## Contributing

Feel free to submit issues and pull requests. Please ensure your code follows Go best practices and includes appropriate tests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Disclaimer

This is example code for educational purposes and is not an official Kalshi SDK. Use at your own risk in production environments.
