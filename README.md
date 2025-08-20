# Chirpy
A server for an X/Twitter type of application named Chirpy. Below are the available REST API endpoints:
- Create an account
- Update their account
- Login
- Refresh tokens
- Revoke refresh tokens
- Webhook to upgrade their account to 'Chirpy Red' via a fictious payment company
- Create/Post a 'Chirp'
- Delete a 'Chirp'
- Retrieve a 'Chirp'
- Retrieve several 'Chirps' with filtering options
- Admins can view health of the app
- Admins can view metrics for number of file server visits
- Server running check

## Installation
Inside a Go module:
- go get github.com/JA50N14/chirpy

## Quick Start
1. Create a .env file inside the root directory of your Go module. Add the following variables
    - DB_URL: path to your PostgreSQL database
    - PLATFORM: "dev"
    - JWT_SECRET: A JWT Secret
    - POLKA_KEY: API Key. Used for upgrading a users account to 'Chirpy Red'
2. Start the server. From root directory of Go module run:
    - Go run .

## REST API Endpoints
- POST /api/users
    - Purpose: Create a new user
    - Payload: email, password
- PUT /api/users
    - Purpose: Update a logged in users email or password
    - Payload: email, password
    - Header: Access Token
- POST /api/login
    -Purpose: Login a user
    -Payload: email, password
- POST /api/refresh
    - Header: Refresh Token
- POST /api/revoke
    - Header: Refresh Token
- POST /api/polka/webhooks
    - Payload: event, data {user_id}
    - Header: POLKA_KEY
    - Note: This is a fictious third party payment app
- POST /api/chirps
    - Payload: body
    - Header: JWT Token
- GET /api/chirps
    - URL filter options: author_id, sort
        - author_id=authorIDHere
        - sort=asc or sort-desc 
- GET /api/chirps/{chirpID}
    - URL Parameter: chirpID you want to retrieve
- Delete /api/chirps/{chirpID}
    - Header: Access Token
    - URL Parameter: chirpID belonging to logged in user to delete
- GET /admin/metrics
    - No headers or body needed
    - Returns number of hits to /app/ api endpoint
- POST /admin/reset
    - Just need PLATFORM="dev" set in your .env file
- /app/
    - No headers or body needed
- GET /api/healthz
    - No headers or body needed