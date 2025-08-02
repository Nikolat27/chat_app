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