**Hot-Coffee: Coffee Shop Management System**

**Overview**

The **Hot-Coffee** project is a Coffee Shop Management System designed to simulate the backend of a coffee shop's ordering system. It handles managing orders, menu items, inventory, and generating reports using a RESTful API. The data is stored in JSON files (instead of a database) for simplicity and easy persistence.

This project is a great exercise to learn backend development, API design, data storage, and software architecture. It implements layered architecture with a clear separation of concerns for maintainability and scalability.

__Features__

- **RESTful API** to handle orders, menu items, and inventory.
- **Data stored in JSON files** for orders, menu items, and inventory.
- **Error Handling & Validation** for input and resource management.
- **Logging:** Logs events and errors for debugging and system monitoring.
- **Aggregations:** Retrieve total sales and popular menu items.
- **Inventory management:** Updates inventory levels when orders are processed.

__Learning Objectives__

- Develop a __REST API.__
- Understand how to work with __JSON__ data.
- Implement logging using Go's `log/slog` package. 
- Apply __layered software architecture__ to separate concerns.
- Learn how to handle data persistence without a database.

__Architecture__

The application uses a __three-layered architecture__ to promote maintainability and scalability:

1. __Presentation Layer (Handlers):__
- Handles HTTP requests and responses.
- Routes requests to appropriate services.
2. __Business Logic Layer (Services):__
- Contains core business rules and functionality.
- Handles aggregations, computations, and logic for managing orders, menu items, and inventory.
3. __Data Access Layer (Repositories):__
- Handles reading and writing to JSON files for storing and retrieving data.

__Project Structure__
```
hot-coffee/
├── cmd/
│   └── main.go             # Entry point of the application
├── internal/
│   ├── handler/            # HTTP Handlers for API routes
│   ├── service/            # Business logic for orders, menu, and inventory
│   └── dal/                # Data access (repositories for each entity)
├── models/                 # Data models for orders, menu items, inventory
├── data/                   # Directory storing JSON data files
│   ├── orders.json         # Orders data
│   ├── menu_items.json     # Menu items data
│   └── inventory.json      # Inventory data
├── go.mod                  # Go module file
├── go.sum                  # Go sum file
└── README.md               # Project documentation (this file)
```

__Endpoints__

__Orders__
- __POST /orders:__ Create a new order.
    - Request: JSON object with customer name and items.
    - Response: Order ID and status.
- **GET /orders:** Retrieve all orders.
    - Response: List of orders in JSON format.
- **GET /orders/{id}:** Retrieve a specific order by ID.
    - Response: A single order in JSON format.
- **PUT /orders/{id}:** Update an existing order.
    - Request: Updated order data.
    - Response: Updated order details.
- **DELETE /orders/{id}:** Delete an order by ID.
    - Response: Confirmation of deletion.
- **POST /orders/{id}/close:** Close an order.
    - Response: Confirmation of order closure.


__Installation & Usage__

1. **Clone the repository:**
```
git clone https://github.com/your-username/hot-coffee.git
cd hot-coffee
```

2. **Install Go:** Make sure you have Go installed. You can download it from [here.](https://go.dev/doc/install)
3. Build the project:
In the root of the project, run:
```
go build -o hot-coffee .
```
4. Run the application:
You can start the application using the following command:
```
./hot-coffee --port 8080 --dir ./data
```
This starts the server on port 8080 with data stored in the `data/` directory.
5. **Access the API:**
The API is available at `http://localhost:8080.` You can interact with it using a tool like **Postman** or **curl.**
6. **Display Usage Information:**
To get help and see available options, use the `--help` flag:
```
./hot-coffee --help
```
This will display:
```
Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>]
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number (default: 8080)
  --dir S      Path to the data directory (default: ./data)
  ```
