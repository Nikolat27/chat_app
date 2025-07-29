# Backend API Specification for E2EE Implementation

## Current Endpoints (Already Implemented)

### 1. Create Secret Chat
```
POST /api/secret-chat/create
Body: { "target_user": user_id }
Response: { "chat": { "id": chat_id, ... } }
```

## E2EE Endpoints

### 2. Upload Public Key
```
POST /api/secret-chat/add-public-key/{secret_chat_id}
Body: { "public_key": base64_encoded_public_key }
Response: { "message": "success" }
```

### 3. Upload Encrypted Symmetric Keys (User B only)
```
POST /api/secret-chat/add-symmetric-key/{secret_chat_id}
Body: { 
    "user_1_encrypted_symmetric_key": base64_encoded_encrypted_key_for_user_1,
    "user_2_encrypted_symmetric_key": base64_encoded_encrypted_key_for_user_2
}
Response: { "message": "success" }
```

### 4. Get Secret Chat Data
```
GET /api/secret-chat/get/{secret_chat_id}
Response: { 
    "id": chat_id,
    "user_1": user_1_id,
    "user_2": user_2_id,
    "user_1_public_key": base64_encoded_public_key,
    "user_2_public_key": base64_encoded_public_key,
    "user_1_encrypted_symmetric_key": base64_encoded_encrypted_key,
    "user_2_encrypted_symmetric_key": base64_encoded_encrypted_key,
    "user_2_accepted": boolean,
    "key_finalized": boolean
}
```

## Group API Endpoints

### 1. Create Group
```
POST /api/group/create
Body: {
    "name": "string",
    "description": "string (optional)",
    "type": "public" | "private",
    "avatar_url": "string (optional)"
}
Response: {
    "group": {
        "id": group_id,
        "name": "string",
        "description": "string",
        "type": "public" | "private",
        "avatar_url": "string",
        "created_by": user_id,
        "created_at": "timestamp",
        "member_count": 1,
        "invite_code": "string (for private groups)"
    }
}
```

### 2. Create Secret Group
```
POST /api/group/create-secret
Body: {
    "name": "string",
    "description": "string (optional)",
    "avatar_url": "string (optional)"
}
Response: {
    "group": {
        "id": group_id,
        "name": "string",
        "description": "string",
        "type": "secret",
        "avatar_url": "string",
        "created_by": user_id,
        "created_at": "timestamp",
        "member_count": 1,
        "invite_code": "string"
    }
}
```

### 3. Get User Groups
```
GET /api/group/user-groups
Response: {
    "groups": [
        {
            "id": group_id,
            "name": "string",
            "description": "string",
            "type": "public" | "private" | "secret",
            "avatar_url": "string",
            "created_by": user_id,
            "created_at": "timestamp",
            "member_count": number,
            "role": "admin" | "moderator" | "member",
            "is_secret": boolean
        }
    ]
}
```

### 4. Join Group
```
POST /api/group/join
Body: {
    "group_id": group_id (for public groups)
    OR
    "invite_code": "string" (for private/secret groups)
}
Response: {
    "message": "Successfully joined group",
    "group": {
        "id": group_id,
        "name": "string",
        "description": "string",
        "type": "public" | "private" | "secret",
        "avatar_url": "string",
        "member_count": number,
        "role": "member"
    }
}
```

### 5. Leave Group
```
POST /api/group/leave/{group_id}
Response: {
    "message": "Successfully left group"
}
```

### 6. Get Group Details
```
GET /api/group/{group_id}
Response: {
    "group": {
        "id": group_id,
        "name": "string",
        "description": "string",
        "type": "public" | "private" | "secret",
        "avatar_url": "string",
        "created_by": user_id,
        "created_at": "timestamp",
        "member_count": number,
        "members": [
            {
                "user_id": user_id,
                "username": "string",
                "avatar_url": "string",
                "role": "admin" | "moderator" | "member",
                "joined_at": "timestamp"
            }
        ],
        "invite_code": "string (for private/secret groups)"
    }
}
```

### 7. Search Groups
```
GET /api/group/search?q={search_term}
Response: {
    "groups": [
        {
            "id": group_id,
            "name": "string",
            "description": "string",
            "type": "public" | "private" | "secret",
            "avatar_url": "string",
            "member_count": number,
            "is_member": boolean
        }
    ]
}
```

### 8. Update Group
```
PUT /api/group/{group_id}
Body: {
    "name": "string (optional)",
    "description": "string (optional)",
    "avatar_url": "string (optional)"
}
Response: {
    "message": "Group updated successfully",
    "group": {
        "id": group_id,
        "name": "string",
        "description": "string",
        "type": "public" | "private" | "secret",
        "avatar_url": "string"
    }
}
```

### 9. Delete Group
```
DELETE /api/group/{group_id}
Response: {
    "message": "Group deleted successfully"
}
```

### 10. Add Member to Group
```
POST /api/group/{group_id}/add-member
Body: {
    "user_id": user_id
}
Response: {
    "message": "Member added successfully"
}
```

### 11. Remove Member from Group
```
POST /api/group/{group_id}/remove-member
Body: {
    "user_id": user_id
}
Response: {
    "message": "Member removed successfully"
}
```

### 12. Change Member Role
```
POST /api/group/{group_id}/change-role
Body: {
    "user_id": user_id,
    "role": "admin" | "moderator" | "member"
}
Response: {
    "message": "Role changed successfully"
}
```

## Database Schema

### Groups Table
```sql
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL DEFAULT 'public', -- 'public', 'private', 'secret'
    avatar_url VARCHAR(500),
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    invite_code VARCHAR(50) UNIQUE,
    is_active BOOLEAN DEFAULT TRUE
);
```

### Group Members Table
```sql
CREATE TABLE group_members (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- 'admin', 'moderator', 'member'
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, user_id)
);
```

### Group Messages Table
```sql
CREATE TABLE group_messages (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    sender_id INTEGER NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text', -- 'text', 'image', 'file'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_encrypted BOOLEAN DEFAULT FALSE,
    encrypted_content TEXT -- for secret groups
);
```

### Secret Group Keys Table
```sql
CREATE TABLE secret_group_keys (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    public_key TEXT NOT NULL,
    encrypted_symmetric_key TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, user_id)
);
```

## Group Types

### Public Groups
- Anyone can join without invitation
- Visible in group search
- No invite codes needed

### Private Groups
- Require invite code to join
- Not visible in public search
- Admin can generate new invite codes

### Secret Groups
- End-to-end encrypted messaging
- Require invite code to join
- Not visible in public search
- Each member has their own encrypted copy of the group symmetric key

## Security Considerations for Groups

- ✅ Public groups are discoverable but require explicit join
- ✅ Private groups require invite codes
- ✅ Secret groups use E2EE similar to secret chats
- ✅ Group admins can manage members and roles
- ✅ Members can leave groups at any time
- ✅ Group creators are automatically admins
- ✅ Invite codes are unique and secure
- ✅ Secret groups require key exchange for each member

## Workflow

### User A (Chat Creator):
1. Creates secret chat
2. Generates key pair locally
3. Uploads public key via `/api/secret-chat/add-public-key/{chat_id}`
4. Waits for User B to approve and generate symmetric key

### User B (Chat Approver):
1. Approves secret chat
2. Generates key pair locally
3. Uploads public key via `/api/secret-chat/add-public-key/{chat_id}`
4. Generates symmetric key K
5. Encrypts K with User A's public key → K_encA
6. Encrypts K with User B's public key → K_encB
7. Uploads both encrypted keys via `/api/secret-chat/add-symmetric-key/{chat_id}`
8. Backend sets `key_finalized = true`

### Both Users (Messaging):
1. Fetch chat data via `/api/secret-chat/get/{chat_id}`
2. Decrypt their respective encrypted symmetric key with their private key
3. Cache symmetric key in memory
4. Use symmetric key for all message encryption/decryption

## Example Implementation (Python/FastAPI)

```python
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import base64

class PublicKeyRequest(BaseModel):
    public_key: str

class SymmetricKeyRequest(BaseModel):
    user_1_encrypted_symmetric_key: str
    user_2_encrypted_symmetric_key: str

@app.post("/api/secret-chat/add-public-key/{chat_id}")
async def add_public_key(chat_id: int, request: PublicKeyRequest):
    # Store public key for current user
    user_id = get_current_user_id()
    secret_chat = get_secret_chat(chat_id)

    if not secret_chat:
        raise HTTPException(status_code=404, detail="Secret chat not found")

    # Update database with public key for this user
    if secret_chat.user_1 == user_id:
        secret_chat.user_1_public_key = request.public_key
    else:
        secret_chat.user_2_public_key = request.public_key

    db.commit()
    return {"message": "success"}

@app.post("/api/secret-chat/add-symmetric-key/{chat_id}")
async def add_symmetric_key(chat_id: int, request: SymmetricKeyRequest):
    user_id = get_current_user_id()
    secret_chat = get_secret_chat(chat_id)

    if not secret_chat:
        raise HTTPException(status_code=404, detail="Secret chat not found")

    # Only User B can upload symmetric keys
    if secret_chat.user_1 == user_id:
        raise HTTPException(status_code=403, detail="Only User B can upload symmetric keys")

    # Store both encrypted symmetric keys
    secret_chat.user_1_encrypted_symmetric_key = request.user_1_encrypted_symmetric_key
    secret_chat.user_2_encrypted_symmetric_key = request.user_2_encrypted_symmetric_key
    secret_chat.key_finalized = True

    db.commit()
    return {"message": "success"}

@app.get("/api/secret-chat/get/{chat_id}")
async def get_secret_chat(chat_id: int):
    user_id = get_current_user_id()
    secret_chat = get_secret_chat(chat_id)

    if not secret_chat:
        raise HTTPException(status_code=404, detail="Secret chat not found")

    # Check if user is participant
    if secret_chat.user_1 != user_id and secret_chat.user_2 != user_id:
        raise HTTPException(status_code=403, detail="Not a participant")

    return {
        "id": secret_chat.id,
        "user_1": secret_chat.user_1,
        "user_2": secret_chat.user_2,
        "user_1_public_key": secret_chat.user_1_public_key,
        "user_2_public_key": secret_chat.user_2_public_key,
        "user_1_encrypted_symmetric_key": secret_chat.user_1_encrypted_symmetric_key,
        "user_2_encrypted_symmetric_key": secret_chat.user_2_encrypted_symmetric_key,
        "user_2_accepted": secret_chat.user_2_accepted,
        "key_finalized": secret_chat.key_finalized
    }
```

## Security Considerations

- ✅ Encrypted keys only stored on backend
- ✅ Backend cannot decrypt messages
- ✅ Each user has their own encrypted copy of symmetric key
- ✅ Private keys never leave the client
- ✅ Symmetric keys generated on client side
- ✅ Base64 encoding for MongoDB compatibility
- ✅ Handles offline scenarios perfectly
- ✅ User B generates and distributes symmetric key
- ✅ Key finalized flag prevents tampering
- ✅ No messaging allowed until keys are finalized 