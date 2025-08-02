# Secure Group Chat

A group chat and messaging app built with Go, MongoDB.

---

## Features

- **User Registration & Login**
    - Register new users, log in, log out
    - Authentication and session check

- **User Profile**
    - Search for users
    - View user info
    - Upload avatar
    - Delete account
    - View personal chats, secret chats, and groups

- **Direct Chats**
    - Create 1-on-1 chats
    - Send text and images
    - View chat messages
    - Delete chats
    - Real-time chat via WebSockets


- **Secret Chats**
    - Create end-to-end encrypted chats (E2EE)
    - Exchange public and symmetric keys
    - Approve or reject secret chat requests
    - View and delete secret chat messages
    - Real-time E2EE messaging via WebSockets

- **Group Chats**
    - Create and update groups (name, description, avatar)
    - Join groups by invite link
    - Add, remove, ban, or unban users
    - View group messages and members
    - Leave or delete group
    - Real-time messaging via WebSockets

- **Secret Groups**
    - Create end-to-end encrypted group chats
    - Group owner generates and shares secret key
    - Manage members, join with invite link
    - Messages are encrypted/decrypted on the client
    - Real-time E2EE group messaging via WebSockets

- **Messages**
    - Edit messages
    - Delete messages (for self or all)
    - Upload images in chat and group messages

- **Saved Messages**
    - Create, edit, and delete personal saved messages

- **Approvals**
    - Submit, edit, or delete approvals for group/secret chat access
    - View sent and received approvals

- **File Uploads**
    - Upload and serve static files (avatars, msg images)

---

## Security Note

> **All regular group and chats msgs are encrypted on server side for extra protection.**

---

## Tech Stack

- Go (Golang)
- MongoDB
- [Chi](https://github.com/go-chi/chi) router
- WebSockets (real-time messaging)
- Static file server for uploads

---