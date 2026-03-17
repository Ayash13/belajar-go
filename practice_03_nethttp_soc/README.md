# Practice 3: Net/HTTP with Separation of Concerns (SoC)

This practice demonstrates how to build a RESTful API using the standard `net/http` package with a clean architecture approach, also known as Separation of Concerns (SoC).

## Architecture

The project is structured into multiple layers to ensure separation of concerns:

- `main.go`: The entry point of the application. It initializes the database connection, sets up the dependencies, configures the router, and starts the HTTP server.
- `database`: Contains the configuration and initialization logic for the database connection (e.g., PostgreSQL).
- `repository`: The data access layer. It handles all interactions with the database, such as executing queries (CRUD operations).
- `service`: The business logic layer. It contains the core rules of the application and coordinates between the repository and the handler.
- `handler`: The presentation layer. It parses incoming HTTP requests, calls the appropriate service methods, and formats the HTTP responses (e.g., JSON encoding/decoding).
- `server`: Contains routing definitions, middleware, and common HTTP helper functions.
- `entity`: Defines the core domain models or structures used across the application.
- `dto` (Data Transfer Object): Defines the structures used for request payloads and response formats to avoid exposing internal entity structures directly.

## How It Works

1. **Request Lifecycle**: An incoming HTTP request hits the server on a specific route (e.g., `/persons`).
2. **Handler**: The configured handler function intercepts the request, reads any path parameters or JSON payload, and converts them to DTOs.
3. **Service**: The handler passes the data to the service layer. The service performs any necessary validation or business logic.
4. **Repository**: If the service needs to read or write data, it calls the repository layer. The repository executes the SQL queries using the injected database connection.
5. **Response**: The repository returns the entity to the service, which may convert it to a response DTO. The handler then sends this data back to the client as an HTTP response (usually in JSON format).

## Running the Application

1. Make sure you have a database running and set the expected credentials in the `.env` file.
2. Run the application:
   ```bash
   go run main.go
   ```
3. The server will start and listen on port `8080`.
