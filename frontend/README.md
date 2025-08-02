# Chat Application Frontend

A modern, real-time chat application built with Vue.js, featuring end-to-end encryption and group chat functionality.

## Features

### ✅ Implemented
- **User Authentication**: Login/Register with JWT tokens
- **Real-time Messaging**: WebSocket-based chat with instant message delivery
- **End-to-End Encryption**: Secret chats with RSA encryption
- **Group Creation**: Create public, private, and secret groups
- **Modern UI**: Beautiful, responsive design with smooth animations
- **File Upload**: Avatar upload functionality
- **Message Encryption**: Automatic encryption/decryption for secret chats

### 🔧 Backend Integration Status

#### Group Functionality
- ✅ **Create Group**: `/api/group/create` - Fully implemented
- ✅ **Get User Groups**: `/api/user/get-groups` - Fully implemented
- ⏳ **Join Group**: `/api/group/join` - Frontend ready, backend needed
- ⏳ **Leave Group**: `/api/group/leave/{id}` - Frontend ready, backend needed
- ⏳ **Search Groups**: `/api/group/search` - Frontend ready, backend needed

#### Chat Functionality
- ✅ **Create Chat**: `/api/chat/create` - Fully implemented
- ✅ **Secret Chat**: `/api/secret-chat/create` - Fully implemented
- ✅ **User Search**: `/api/user/search` - Fully implemented
- ✅ **Get Chats**: `/api/user/get-chats` - Fully implemented

## Getting Started

### Prerequisites
- Node.js (v16 or higher)
- npm or yarn

### Installation
```bash
npm install
```

### Development
```bash
npm run dev
```

### Build
```bash
npm run build
```

## Group Creation Guide

### Creating Groups
1. Navigate to the Groups tab in the sidebar
2. Click "Create Group" button
3. Fill in the group details:
   - **Name**: Required
   - **Description**: Optional
   - **Type**: Public (anyone can join) or Private (invite only)
4. Click "Create Group"

### Creating Secret Groups
1. Click "Create Secret Group" button
2. Fill in the group details:
   - **Name**: Required
   - **Description**: Optional
3. Click "Create Secret Group"

### Group Types
- **Public Groups**: Anyone can join without invitation
- **Private Groups**: Require invite code to join
- **Secret Groups**: End-to-end encrypted messaging

## API Endpoints

### Group Endpoints
```bash
POST /api/group/create          # Create a new group
GET  /api/user/get-groups       # Get user's groups
POST /api/group/join           # Join a group (by ID or invite code)
POST /api/group/leave/{id}     # Leave a group
GET  /api/group/search         # Search for groups
GET  /api/group/{id}           # Get group details
```

### Chat Endpoints
```bash
POST /api/chat/create          # Create a new chat
POST /api/secret-chat/create   # Create a secret chat
GET  /api/user/get-chats       # Get user's chats
GET  /api/user/search          # Search for users
```

## Technology Stack

- **Frontend**: Vue.js 3 with Composition API
- **State Management**: Pinia
- **Styling**: Tailwind CSS
- **HTTP Client**: Axios
- **Real-time**: WebSocket
- **Encryption**: Web Crypto API
- **Build Tool**: Vite

## Project Structure

```
src/
├── components/
│   ├── chat/           # Chat-related components
│   ├── tabs/           # Tab components (Chats, Groups, Settings)
│   └── ui/             # Reusable UI components
├── stores/             # Pinia stores
├── composables/        # Vue composables
├── utils/              # Utility functions
└── views/              # Page components
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License.
