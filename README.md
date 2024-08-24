RSS Aggregator

This project is an RSS aggregator, built using Go, that collects and processes RSS feeds from various sources. It parses the feeds, stores the information in a PostgreSQL database, and serves it via a RESTful API.


Features:

Fetch and Parse RSS Feeds: Automatically fetches RSS feeds from configured sources and parses them.
Store Data in PostgreSQL: Stores parsed feed data into a PostgreSQL database.
RESTful API: Provides a RESTful API to access the stored feed data.
Docker Support: Dockerized setup for easy deployment and scaling.


Requirements:

Go 1.20+: The application is built with Go.
PostgreSQL: For storing the RSS feed data.
Docker: Optional, for running the application and PostgreSQL database in containers.


Installation:

Clone the repository

git clone https://github.com/Anas-Sayed0/rss-agg.git
cd rss-agg

Set up PostgreSQL:
Ensure you have PostgreSQL installed and running. You can also use Docker to run PostgreSQL:

bash
docker-compose up -d

Create a PostgreSQL database:

sql

CREATE DATABASE gobank;
Configure your environment variables in the .env file to match your PostgreSQL setup.

Build and Run
Build the Go application:

bash
go build -o rss-agg


Run the application:

bash
./rss-agg
Running with Docker
Ensure Docker is installed and running on your machine.

Build and run the Docker containers:

bash
docker-compose up --build
API Endpoints
GET /feeds: Fetch all the RSS feeds.
POST /feeds: Add a new RSS feed source.
GET /feeds/{id}: Fetch a specific RSS feed by ID.

Configuration:
Configuration is managed via environment variables:

POSTGRES_USER: PostgreSQL username
POSTGRES_PASSWORD: PostgreSQL password
POSTGRES_DB: PostgreSQL database name
POSTGRES_HOST: PostgreSQL host (default: localhost)
POSTGRES_PORT: PostgreSQL port (default: 5432)


Contributing:
Contributions are welcome! Please submit a pull request or open an issue to discuss changes.