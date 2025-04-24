# 🏡 Dwello API

Dwello API is the backend service for the **Dwello real estate app**, built using **Go** and **MongoDB**. It powers the Dwello Android app (built with Kotlin + Jetpack Compose) and provides RESTful endpoints to manage users, properties, and rental requests.

---

## ✨ Features

### 👤 User Management
- 🔐 Register or login with email.
- 📍 Update user details like location and preferred areas.
- ❤️ View liked and posted properties.

### 🏠 Property Management
- 🛠️ Create, update, or delete properties.
- 🔎 Search by location, price, and more.
- 👍 Like/unlike properties.
- 🏘️ Homescreen recommendations based on preferences.

### 📬 Rental Requests
- 📤 Send rental requests for properties.
- ✅ Accept or ❌ reject rental requests.
- 📦 View rental requests for owned properties.

---

## 🧰 Tech Stack

- **Backend Framework**: [Fiber](https://gofiber.io/) (Go) ⚙️  
- **Database**: MongoDB 🍃  
- **API Docs**: Swagger (via Swag CLI) 📖

---

## 📁 Project Structure

```
dwello-api/
├── config/          # 🔧 Database config
├── db/              # 📂 MongoDB collections
├── docs/            # 🧾 Swagger docs
├── handlers/        # 🪝 Route handlers
├── models/          # 🧬 Data models
├── routes/          # 🚦 Route definitions
├── utils/           # 🧰 Utility functions
├── main.go          # 🚀 App entry point
├── go.mod           # 📦 Go module config
└── go.sum           # 🧮 Dependency checksums
```

---

## 🚀 Getting Started

1. **Clone the repo**:
   ```sh
   git clone https://github.com/your-repo/dwello-api.git
   cd dwello-api
   ```

2. **Install dependencies**:
   ```sh
   go mod tidy
   ```

3. **Configure MongoDB**:  
   Start MongoDB locally or update the connection string in `config/db.go`.

4. **Run the app**:
   ```sh
   go run main.go
   ```

5. **Access the API**:  
   Open your browser at `http://localhost:8080`.

---

## 📚 API Documentation

Swagger docs are available at:  
👉 `http://localhost:8080/swagger/index.html`  
Explore all endpoints, request parameters, and response formats interactively.

---

## 🔗 Example Endpoints

### 👤 User Routes
- `POST /api/users/register` – Register or login
- `GET /api/users/:email` – Get user by email
- `PUT /api/users/:email/location` – Update location

### 🏘️ Property Routes
- `POST /api/properties` – Create a new property
- `GET /api/properties/search` – Search properties
- `POST /api/properties/:id/like` – Like/unlike a property

### 📩 Rental Requests
- `POST /api/properties/:id/request` – Send rental request
- `POST /api/users/rental-requests/:id/handle` – Accept/reject request

---

## 📄 License

This project is licensed under the **MIT License**.
