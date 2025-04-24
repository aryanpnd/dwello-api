# Dwello API Documentation

## Base URL
```
http://localhost:8080/api
```

---

## User Routes

### Register or Login User
**POST** `/users/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "name": "John Doe",
  "profile_pic": "https://example.com/profile.jpg",
  "location": "New York",
  "preferred_location": "Manhattan"
}
```

**Response:**
- **201 Created**: User successfully registered.
- **200 OK**: User already exists, returns the existing user.

---

### Get User by Email
**GET** `/users/:email`

**Response:**
- **200 OK**: Returns user details.
- **404 Not Found**: User not found.

---

### Update User Location
**PUT** `/users/:email/location`

**Request Body:**
```json
{
  "location": "New York"
}
```

**Response:**
- **200 OK**: Location updated.
- **400 Bad Request**: Invalid request body.
- **500 Internal Server Error**: Update failed.

---

### Update User Preferred Location
**PUT** `/users/:email/preferred-location`

**Request Body:**
```json
{
  "preferred_location": "Manhattan"
}
```

**Response:**
- **200 OK**: Preferred location updated.
- **400 Bad Request**: Invalid request body.
- **500 Internal Server Error**: Update failed.

---

### Get Liked Properties
**GET** `/users/:email/liked-properties`

**Response:**
- **200 OK**: Returns a list of liked properties.
- **404 Not Found**: User not found.
- **500 Internal Server Error**: Failed to fetch properties.

---

### Get Posted Properties
**GET** `/users/:email/posted-properties`

**Response:**
- **200 OK**: Returns a list of posted properties.
- **404 Not Found**: User not found.
- **500 Internal Server Error**: Failed to fetch properties.

---

## Property Routes

### Create Property
**POST** `/properties`

**Request Body:**
```json
{
  "title": "Beautiful Apartment",
  "description": "A spacious apartment in Manhattan.",
  "price": 2500,
  "location": "Manhattan",
  "owner_email": "user@example.com",
  "owner_name": "John Doe",
  "owner_pic": "https://example.com/profile.jpg",
  "thumbnail": "https://example.com/thumbnail.jpg",
  "pictures": ["https://example.com/pic1.jpg", "https://example.com/pic2.jpg"]
}
```

**Response:**
- **201 Created**: Property successfully created.
- **400 Bad Request**: Invalid request body.
- **403 Forbidden**: User is not the owner.
- **500 Internal Server Error**: Failed to create property.

---

### Update Property
**PUT** `/properties/:id`

**Request Body:**
```json
{
  "title": "Updated Apartment Title",
  "description": "Updated description.",
  "price": 3000
}
```

**Response:**
- **200 OK**: Property updated.
- **400 Bad Request**: Invalid property ID or request body.
- **403 Forbidden**: User is not the owner.
- **500 Internal Server Error**: Failed to update property.

---

### Delete Property
**DELETE** `/properties/:id`

**Response:**
- **200 OK**: Property deleted.
- **400 Bad Request**: Invalid property ID.
- **403 Forbidden**: User is not the owner.
- **500 Internal Server Error**: Failed to delete property.

---

### Like Property
**POST** `/properties/:id/like`

**Response:**
- **200 OK**: Property liked.
- **400 Bad Request**: Invalid property ID.
- **500 Internal Server Error**: Failed to like property.

---

### Unlike Property
**POST** `/properties/:id/unlike`

**Response:**
- **200 OK**: Property unliked.
- **400 Bad Request**: Invalid property ID.
- **500 Internal Server Error**: Failed to unlike property.

---

### Search Properties
**GET** `/properties/search`

**Query Parameters:**
- `location` (optional): Filter by location.
- `min_price` (optional): Minimum price.
- `max_price` (optional): Maximum price.

**Response:**
- **200 OK**: Returns a list of properties matching the criteria.
- **500 Internal Server Error**: Failed to fetch properties.

---

### Get Homescreen Properties
**GET** `/properties/homescreen`

**Response:**
- **200 OK**: Returns properties based on the user's preferred location.
- **400 Bad Request**: Preferred location not set.
- **500 Internal Server Error**: Failed to fetch properties.

---
