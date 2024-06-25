# RSS Feed Aggregator

This project is an RSS feed aggregator built with Go. It allows users to add and follow RSS feeds, and periodically fetches and stores new posts from these feeds.

## Features

- User authentication
- Add and manage RSS feeds
- Automatically fetch new posts from followed feeds
- RESTful API for accessing aggregated content

## Technologies Used

- Go
- PostgreSQL
- sqlc for database query generation
- goose for database migrations

## Setup

1. Clone the repository:

   ```
   git clone https://github.com/Delvoid/rss.git
   ```

2. Install dependencies:

   ```
   go mod tidy
   ```

3. Set up the database:

   - Create a PostgreSQL database
   - Update the database connection string in your `.env` file

4. Run database migrations:

   ```
   goose postgres "your-connection-string" up
   ```

5. Build and run the application:
   ```
   go build
   ./rss-feed-aggregator
   ```

## API Endpoints

- `POST /v1/users`: Create a new user
- `GET /v1/users`: Get authenticated user's information
- `POST /v1/feeds`: Create a new feed
- `GET /v1/feeds`: Get all feeds
- `POST /v1/feed_follows`: Follow a feed
- `GET /v1/feed_follows`: Get user's followed feeds
- `DELETE /v1/feed_follows/{feedFollowID}`: Unfollow a feed
- `GET /v1/posts`: Get posts from followed feeds
