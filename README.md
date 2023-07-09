# URL Shortener Application

This is a URL Shortener application written in Golang. The application provides functionality to shorten long URLs into short, easy-to-use aliases. It also supports generating QR codes for the shortened URLs and keeps track of the click statistics.

## Setup and Configuration

To run the URL Shortener application, please follow the instructions below:

1. Make sure you have Go installed on your system. You can download and install Go from the [official website](https://golang.org/).

2. Clone the repository to your local machine:

   ```
   git clone https://github.com/seyifunmiSF/URL-SHORTENER.git
   ```

3. Navigate to the project directory:

   ```
   cd url-shortener
   ```

4. Create a `.env` file in the project root directory and configure the following environment variables:

   ```
   DB_CONNECTION_STRING=<your-database-connection-string>
   ```

   The `DB_CONNECTION_STRING` environment variable should contain the connection string for your PostgreSQL database.

5. Install the required dependencies using the following command:

   ```
   go mod download
   ```

6. Start the application by running the following command:

   ```
   go run main.go
   ```

7. The application will start running on `http://localhost:8080`.

## Dependencies

The URL Shortener application uses the following dependencies:

- [gin-gonic/gin](https://github.com/gin-gonic/gin): A web framework for building APIs in Go.
- [go-redis/redis/v8](https://github.com/go-redis/redis/v8): A Golang Redis client library.
- [jinzhu/gorm](https://github.com/jinzhu/gorm): An ORM library for Golang.
- [joho/godotenv](https://github.com/joho/godotenv): A Go library for loading environment variables from a `.env` file.
- [lib/pq](https://github.com/lib/pq): A pure Go PostgreSQL driver for the database/sql package.

## Project Structure

The project is structured as follows:

- `domain`: Contains domain models and interfaces.
- `handlers`: Handles the incoming HTTP requests and calls the corresponding services.
- `repositories`: Implements the data access layer using the PostgreSQL database and Redis.
- `services`: Implements the business logic and provides services to handle the requests.

## API Endpoints

The URL Shortener application exposes the following API endpoints:

- `POST /shorten`: Shortens a long URL and returns the shortened URL.
- `GET /history`: Retrieves the history of shortened URLs.
- `GET /s/:alias`: Redirects the user to the original URL associated with the given alias.
- `GET /shorten/:alias/stats`: Retrieves click statistics for a shortened URL.
- `GET /static/:filepath`: Serves static files such as images.

## Configuration

The URL Shortener application uses a configuration file to set up certain parameters. The configuration can be found in the `main.go` file, and the available options are:

- `BaseURL`: The base URL of the application.
- `QRCodeDirectory`: The directory where QR code images are stored.
- `UrlImagePath`: The path for serving static URL images.

Please modify these values according to your specific setup and requirements.

## Contributing

Contributions to the URL Shortener application are welcome! If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

