# Doctor Service API

This is a microservice for managing doctors on the medical platform. It provides a REST API for creating and retrieving doctor records using Clean Architecture. The project is structured in such a way that it can be scaled in the future.

## Technologies

* Language: Go 1.25.5
* Framework: Gin for HTTP routing
* Database: MongoDB
* Architecture: Clean Architecture (Model, Repository, Usecase, Transport)

## How to Run the Project

### Requirements
1. Go installed (version 1.25.5 or higher).
2. A running MongoDB server on the default port (localhost:27017). The service will automatically create the doctor_db database and the doctors collection.

### Running the service
In the root directory of the project, run the following command in your terminal:
```bash
  go run ./cmd/doctor-s/main.go
```
The service will start on port 8081. It also supports Graceful Shutdown to safely close connections when you stop the server.

---

## REST API Endpoints

### 1. Create a Doctor
Adds a new doctor to the database.

* URL: /doctors
* Method: POST
* Request Body (JSON Example):
```json
  {
      "full_name": "John Doe",
      "specialization": "Therapist",
      "email": "johndoe@hospital.com"
  }
```
**Business Rules (Validation):**
* The full_name field is required.
* The email field is required and must contain an @ symbol and a valid domain format.
* The email must be unique across the system.

**Responses:**
* 200 OK: Doctor created successfully. Returns the doctor object with the generated ID.
* 400 Bad Request: Validation error (missing fields, invalid email format, or email already exists).
* 500 Internal Server Error: Server or database error.

---

### 2. Get Doctor by ID
Retrieves information about a specific doctor using their unique identifier.

* URL: /doctors/:id
* Method: GET

**Responses:**
* 200 OK: Successful request. Returns the doctor JSON object.
* 404 Not Found: Returns an error message ("there is no doctor like this") if the doctor with this ID is not found.
* 500 Internal Server Error: Server error.

---

### 3. Get All Doctors
Returns an array of all doctors registered in the system.

* URL: /doctors
* Method: GET

**Responses:**
* 200 OK: Successful request. Returns an array of objects. If there are no doctors in the database, it returns an empty array "[]" instead of null.
* 500 Internal Server Error: Error fetching data.

---

## Project Structure (Clean Architecture)

* `cmd/doctor-s/main.go`: The entry point of the application. Handles DB connection, application layer setup, and server launch.
* `internal/model`: Defines data structures (Doctor) and common system errors.
* `internal/repository`: The layer responsible for interacting with the MongoDB database.
* `internal/usecase`: The business logic layer containing data validation rules, such as email format validation via regular expressions.
* `internal/transport/http`: Handlers for processing HTTP requests via the Gin framework and managing response status codes.