# PUFA Computing Backend

This is the backend repository for the PUFA Computing Backend project. It provides the server-side logic and API endpoints necessary to support the [front-end application](https://github.com/PUFA-Computing/Frontend).

## Table of Contents

- [Description](#description)
- [Technologies Used](#technologies-used)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Description

This backend repository is built using Golang and PostgreSQL. It serves as the server-side component of the [Project Name] project, providing API endpoints for data retrieval, manipulation, and authentication. The backend includes RBAC (Role-Based Access Control) authorization and utilizes JWT (JSON Web Tokens) for secure user authentication.

## Technologies Used

- Golang: The backend server is implemented in Golang, a statically typed programming language known for its simplicity and performance.
- PostgreSQL: The database management system used for storing and managing application data.
- Nginx: A web server used for reverse proxying and serving static files.
- RBAC Authorization: Role-Based Access Control is implemented to manage user access rights and permissions.
- JWT (JSON Web Tokens): JWT is used for secure authentication and authorization.

## Installation

To run the backend server locally, follow these steps:

1. Clone this repository to your local machine:

   ```bash
   git clone https://github.com/your-organization/backend.git
2. Navigate to the project directory:
   ```bash
   cd backend
4. Install dependencies:
   ```bash
   go mod tidy
6. Set up the PostgreSQL database and configure the database connection in the .env file.
8. Build the application:
   ```bash
   go build -o main cmd/app/main.go

## Usage

To start the backend server, run the following command:
```bash
go run cmd/app/main.go
```
The server will start listening on the specified port (default: 8080) and will be ready to handle incoming requests.

## Contributing

Contributions to this project are welcome! Feel free to open an issue or submit a pull request.

## License
This project is licensed under the [Creative Commons Attribution-NonCommercial (CC BY-NC) license](LICENSE).
