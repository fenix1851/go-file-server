# Go File Server

This is a simple file server implemented in Go that serves files and directories from a specified path. It also includes user authentication using JWT and supports basic user registration and login functionalities.

## Features

- Serve files and directories using a web interface.
- User authentication using JWT (JSON Web Tokens).
- Basic user registration and login system.
- Template-based UI using HTML templates.

## Getting Started

1. Clone the repository:

git clone https://github.com/fenix1851/go-file-server

2. If you want just launch the server, click on the executable file. Otherwise, if you want to build the server yourself, go to the next step.

3. Get the dependencies:    

go get 

4. Build the server:

go build -o go-file-server main.go

5. Run the server:

./go-file-server

6. Open your browser and go to http://localhost:4001

## Usage

- Access the file server by navigating to `http://localhost:4001` in your web browser.
- Use the login and registration forms to authenticate and register users.
- Upon successful login, you'll be able to access the file server interface and browse files and directories.

## Contributing

Contributions are welcome! If you find any issues or want to add new features, please feel free to submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).