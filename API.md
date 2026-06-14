# Chatroom API Documentation

## Authentication
Most endpoints require authentication via cookie (set after login).

## Endpoints

### 1. Create User
**POST** `api/users`

Creates a new user.

**Request Body:**
```
json
{
    "username": "string",
    "passworrd": "string"
}
```

**Response Body Example:**
```
    {
		"id":        9cc4f409-bdb7-4226-bc77-3bfe941b1e7c,
		CreatedAt:   2026-05-23 13:18:55.131498,
		UpdatedAt:   2026-05-23 13:18:55.131498,
		Username:    test123,
	}
```
