# Gator CLI - RSS Feed Aggregator

## Prerequisites

Before running the Gator CLI, ensure you have the following installed on your system:

1. PostgreSQL
2. Go (version 1.23.1 or later)

## Installation

To install the Gator CLI, run the following command:

`go install github.com/Jschles1/gator@latest`

This will install the `gator` command-line tool on your system.

## Configuration

1. Create a `.gatorconfig.json` file in your home directory with the following structure:

`{
  "port": "8080",
  "db_url": "postgresql://username:password@localhost:5432/database_name?sslmode=disable"
}`

Replace `username`, `password`, and `database_name` with your PostgreSQL credentials and desired database name.

2. Set up the database by running the SQL scripts provided in the `sql/schema` directory.

## Running the Program

To run the Gator CLI, use the following syntax:

`gator [command] [arguments]`

Here are some available commands:

- `register`: Create a new user account
  `gator register [username] [email]`

- `login`: Log in to an existing account
  `gator login [username]`

- `addfeed`: Add a new RSS feed (requires login)
  `gator addfeed [feed_name] [feed_url]`

- `feeds`: List all available feeds
  `gator feeds`

- `follow`: Follow a feed (requires login)
  `gator follow [feed_url]`

- `browse`: View recent posts from followed feeds (requires login)
  `gator browse`

For a complete list of commands and their usage, refer to the source code or run `gator help`.

Enjoy using Gator CLI to aggregate and manage your favorite RSS feeds!


