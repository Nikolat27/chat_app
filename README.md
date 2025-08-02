# Real-Time E2EE Chat App

A modern, privacy-focused real-time chat application.  
Built with **Go** and **MongoDB** on the backend, and **Vue.js** (JavaScript) + **Tailwind CSS** on the frontend.  
All private conversations use **end-to-end encryption (E2EE)** for maximum security.

---

## ⚡️ Features

-   Real-time messaging (WebSocket)
-   **End-to-end encrypted** (E2EE) group & private chats
-   Modern UI with Tailwind CSS
-   User authentication & profiles
-   Fast, lightweight, and scalable

---

## 🛠️ Tech Stack

-   **Backend:** Go (Golang)
-   **Database:** MongoDB
-   **Frontend:** Vue.js (JavaScript) + Tailwind CSS
-   **Real-Time:** WebSockets
-   **Containerization:** Docker & Docker Compose

---

## 🚀 Getting Started

### Prerequisites

-   [Docker](https://docs.docker.com/get-docker/)
-   [Docker Compose](https://docs.docker.com/compose/install/)

---

### 1. Clone the repository

```bash
git clone https://github.com/nikolat27/chat_app.git
cd chat_app
```

### 2. Fill the Environmental variables in compose.yaml file

### 3. Run the docker compose

sudo docker-compose up --build

### Or Run the project locally (without docker)

Install nginx

cd backend, go mod download, make runserver (port 8000 default)

cd frontend, npm install --force, make runserver (port 80 default)

## 🔐 How Encryption Works

All secret chats and secret groups use **end-to-end encryption (E2EE)**, so only the participants can read messages—no one else, not even the server.

### Secret 1-on-1 Chats

-   Each user generates a public/private key pair (stored in the localforage).
-   When a secret chat starts (second user approves the secret chat), users exchange public keys.
-   A **unique symmetric chat key** is generated (by the second user who approves the secret chat).
-   This chat key is **encrypted separately with each participant’s public key** (also stored in the db in encrypted version ) and sent to them.
-   All messages in the secret chat are encrypted and decrypted on your device using the symmetric chat key.
-   The server never sees your private keys or the unencrypted chat key.

### Secret Groups

-   The group owner generates a **symmetric group key** (stored in the localforage).
-   The owner must send that secret key to all users he want
-   Unlock group messaging by entering your secret key
-   This key is securely shared with each group member by encrypting it with their public key.
-   All group messages are encrypted/decrypted using the group key.
-   Only members who have received the group key can read the messages.

### Key Points

-   **Encryption and decryption always happen on your device, never on the server.**
-   **Private keys never leave your device.**
-   If you lose your private key, you lose access to your secret chats/groups.

For a deeper technical dive, check out our docs or open an issue!
