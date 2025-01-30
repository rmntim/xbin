# XBin - A Pastebin Clone in Golang

XBin is a simple Pastebin clone written in Go (Golang). It allows users to store and share text snippets easily. The application provides a RESTful API for creating, reading, and deleting pastes, as well as a web interface for easy interaction.

## Features

- **Create Pastes**: Store text snippets with optional expiration times.
- **Read Pastes**: Retrieve pastes using a unique URL.
- **Delete Pastes**: Remove pastes manually or automatically after expiration.
- **RESTful API**: Integrate with other applications using the provided API.
- **Web Interface**: A simple web interface for creating and viewing pastes.

## Installation

### Prerequisites

- Go 1.23 or higher
- Git

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/xbin.git
   cd xbin
   ```

2. Build the application:

   ```bash
   go build -o xbin
   ```

3. Run the application:

   ```bash
   ./xbin
   ```

   By default, the application will start on `http://localhost:8080`.

## Usage

### Web Interface

1. **Create a Bin**:
   - Navigate to `http://localhost:8080`.
   - Enter your text in the provided textarea.
   - Optionally, set an expiration time.
   - Click "Bin".

2. **View a Bin**:
   - After creating a paste, you will be redirected to a unique URL (e.g., `http://localhost:8080/bin/abc123`).
   - You can share this URL with others to allow them to view the paste.

### API Endpoints

- **Create a Paste**:
  - **Endpoint**: `POST /bin`
  - **Body**:
    ```jsonc
    {
      "content": "Your text here",
      "expiration": "10m" // Optional: e.g., "10m", "1h", "24h"
    }
    ```
  - **Response**:
    ```jsonc
    {
      "url": "http://localhost:8080/bin/abc123"
    }
    ```

- **Read a Paste**:
  - **Endpoint**: `GET /bin/{slug}`
  - **Response**:
    A HTML page with your code.

## Configuration

XBin can be configured using command-line flags:

- `-port`: The port on which the application will run (default: `8080`). Can also be specified in `PORT` environment variable.
- `-storagePath`: The directory where pastes will be stored (default: `./data/bins.db`).
- `-env`: The environment of the application, either `dev` or `prod` (default: `prod`).

Also [Turso](https://turso.tech) may be used as a database replica provider.
Specify both `TUSRO_TOKEN` and `TUSRO_URL` environment variables to enable this feature.
An example `.env` file is provided.

Example:

```bash
./xbin -port 3000 -env dev
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [Pastebin](https://pastebin.com).
- Uses [go-libsql](https://github.com/tursodatabase/go-libsql) as database driver.
