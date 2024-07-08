# Track Docs

Track Docs is an application designed to manage and track project documents efficiently. Built with Golang, it leverages a microservices architecture to ensure scalability and reliability. The application includes services such as Gin API, kafka for background processing, Redis for caching and Firebase for secure storage.

## Features

- **Document Management**: Upload, store, and organize project documents.
- **Tracking**: Track changes and updates to documents.
- **User Authentication**: Secure user authentication and authorization.
- **Notifications**: Receive notifications for document updates.
- **Search**: Search for documents using various filters.
- **Scalability**: Microservices architecture ensures scalability and flexibility.

## Architecture

Track Docs is composed of several microservices:

- **API**: A Gin-based REST API for handling client requests.
- **Kafka**: Handles queuing for processing tasks and communication between services.
- **Redis**: Caches frequently accessed data to improve performance.
- **Firebase Storage**: Stores and manages project documents.

## Installation

### Prerequisites

- Go 1.22+
- Docker
- Docker Compose
- Make

### Steps

1. **Clone the repository**

    ```sh
    git clone https://github.com/yourusername/track-docs.git
    cd track-docs
    ```

2. **Setup dev environment**

    ```sh
    make setup-dev
    ```

3. **Build and run services with Docker Compose**

    ```sh
    make run
    ```

## Usage

### Accessing the API

Once the services are up, you can access the API at `http://localhost:8080`.
