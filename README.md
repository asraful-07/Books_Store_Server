# Bookio Store — Go Backend

Backend API for the Bookio Store using Go and MongoDB Atlas.

## Features

- Book listing & details
- Category retrieval
- Wishlist (favorites) per user
- Modular structure for easy expansion (auth, orders, users, etc.)

## Tech Stack

- Go (1.21+)
- MongoDB Atlas
- Gorilla Mux for routing

## Setup

1. Copy `.env.example` to `.env` and fill in:
   ```env
   MONGO_URI=...
   DATABASE_NAME=bookio
   PORT=8080
   JWT_SECRET=...
   ```
