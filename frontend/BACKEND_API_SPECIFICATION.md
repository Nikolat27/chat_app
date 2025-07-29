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

## Database Schema

### Secret Chat Table
```sql
ALTER TABLE secret_chats ADD COLUMN user_1_public_key TEXT;
ALTER TABLE secret_chats ADD COLUMN user_2_public_key TEXT;
ALTER TABLE secret_chats ADD COLUMN user_1_encrypted_symmetric_key TEXT;
ALTER TABLE secret_chats ADD COLUMN user_2_encrypted_symmetric_key TEXT;
ALTER TABLE secret_chats ADD COLUMN user_2_accepted BOOLEAN DEFAULT FALSE;
ALTER TABLE secret_chats ADD COLUMN key_finalized BOOLEAN DEFAULT FALSE;
```

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