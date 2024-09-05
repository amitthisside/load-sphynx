
# Load Sphynx

Building a load balancer named "Load Sphynx" from scratch using Go. **Note:** This is for educational purposes and should not be used in production.

## Prerequisites

Ensure that the following dependencies are installed on your system:

-   **Go**: Install from [golang.org](https://golang.org/doc/install).
-   **Python**: Install from [python.org](https://www.python.org/downloads/).
-   **Flask**: Installation steps at [flask.palletsprojects.com](https://flask.palletsprojects.com/en/3.0.x/installation/).

## Setting Up

### Clone the Repository

First, clone the repository and navigate into the `load-sphynx` directory:
```console
git clone https://github.com/amitthisside/load-sphynx.git

cd load-sphynx
```
### Install Go Dependencies

1.  **Initialize Go modules:**
    ```console
    go mod init github.com/amitthisside/load-sphynx
    ```
2.  **Tidy up the module dependencies:**
    ```console
    go mod tidy
    ```

### Setting Up Python Backend Servers

1.  **Create a virtual environment (optional but recommended):**
    ```console
    virtualenv env

    source env/bin/activate 
    ```
2.  **Install Flask:**
    ```console
    pip install flask 
    ```
3.  **Start multiple instances of `server.py`:**
    ```console
    for i in {1..5}; do python server.py "server-$i" "500$i" &; done 
    ```
    This will start 5 instances of the Flask server, each on a different port (from 5001 to 5005).
    

### Configurations

Ensure that your `server_conf.json` file is correctly set up to define the backend server configurations. Example:

```json
{
    "servers": [
      "http://localhost:5000",
      "http://localhost:5001",
      "http://localhost:5002",
      "http://localhost:5003",
      "http://localhost:5004"
    ]
  }
```

## Running the Load Sphynx

1.  **Start the Load Sphynx load balancer:**
    ```console
    go run main.go 
    ``` 
    This will run the load balancer, which fetches configurations from `server_conf.json`.
    
2.  **Bombard the Load Sphynx with requests to test load balancing:**
    ```console
    for i in {1..20}; do curl http://localhost:8000; done 
    ```

  This will send 20 requests to the load balancer, which will distribute them among the running servers.
    

## Stopping the Backend Servers

To terminate all running instances of the servers:
```console
pkill -9 python
```
