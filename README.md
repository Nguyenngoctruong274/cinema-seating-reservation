# 🎬 Cinema Seating Reservation System

## 📖 Overview
This project implements a **Cinema Seating Reservation System** that enforces social distancing rules.
It provides RESTful APIs :
- Configure cinema layouts
- Query available seats
- Reserve seats 
- Cancel reservations.

## ⚙️ Setup Instructions

### Prerequisites
- **Go** (version 1.18 or later)
- **Gin Framework** (used for RESTful APIs)
- **Git** (for cloning the repository)

### Steps
1. **Clone the Repository**:
   ```
   git clone <repository-url>
   cd <repository-folder>

2. **Install Dependencies**:
   ```
   go mod tidy
   ```
3. **Run the Application**: Start the server on port 8080:
   ```
    go run main.go
    ```
4. **Access the API**:
    Open your browser and navigate to `http://localhost:8080` to access the API documentation.

5. **Test the Endpoints**:
    You can use the provided endpoints to interact with the system. The main endpoints include:
    - `POST /cinema/configure`: Configure cinema layout
    - `GET /cinema/available-seats`: Get available seats
    - `POST /cinema/reserve`: Reserve seats
    - `POST /cinema/cancel`: Cancel reservations
 Use tools like Postman or cURL to test the API endpoints or go test ./... -v

6. **Document API Postman**:
   https://documenter.getpostman.com/view/36212477/2sB2x2JuCH
