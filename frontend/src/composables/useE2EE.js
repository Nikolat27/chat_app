import { ref } from 'vue';
import localforage from 'localforage';
import { useKeyPair } from './useKeyPair';

export function useE2EE() {
    const { generateKeyPair, exportSecretChatPublicKey, getSecretChatPrivateKey } = useKeyPair();
    
    // In-memory storage for decrypted symmetric keys (more secure)
    const symmetricKeys = ref(new Map());
    
    // Generate a symmetric key for a chat
    const generateSymmetricKey = async () => {
        try {
            // Generate a random 256-bit key
            const key = await window.crypto.subtle.generateKey(
                {
                    name: "AES-GCM",
                    length: 256
                },
                true,
                ["encrypt", "decrypt"]
            );
            
            // Export the key as raw bytes
            const rawKey = await window.crypto.subtle.exportKey("raw", key);
            return Array.from(new Uint8Array(rawKey));
        } catch (error) {
            console.error('Error generating symmetric key:', error);
            throw new Error('Failed to generate symmetric key');
        }
    };
    
    // Encrypt symmetric key with a user's public key
    const encryptSymmetricKey = async (symmetricKey, publicKeyJwk) => {
        try {
            // Import the public key
            const publicKey = await window.crypto.subtle.importKey(
                "jwk",
                publicKeyJwk,
                {
                    name: "RSA-OAEP",
                    hash: "SHA-256"
                },
                false,
                ["encrypt"]
            );
            
            // Encrypt the symmetric key
            const encryptedKey = await window.crypto.subtle.encrypt(
                {
                    name: "RSA-OAEP"
                },
                publicKey,
                new Uint8Array(symmetricKey)
            );
            
            // Return base64 encoded encrypted key
            return btoa(String.fromCharCode(...new Uint8Array(encryptedKey)));
        } catch (error) {
            console.error('Error encrypting symmetric key:', error);
            throw new Error('Failed to encrypt symmetric key');
        }
    };
    
    // Decrypt symmetric key with user's private key
    const decryptSymmetricKey = async (encryptedKeyBase64, privateKeyJwk) => {
        try {
            // Import the private key
            const privateKey = await window.crypto.subtle.importKey(
                "jwk",
                privateKeyJwk,
                {
                    name: "RSA-OAEP",
                    hash: "SHA-256"
                },
                false,
                ["decrypt"]
            );
            
            // Decode base64 encrypted key
            const encryptedKey = Uint8Array.from(atob(encryptedKeyBase64), c => c.charCodeAt(0));
            
            // Decrypt the symmetric key
            const decryptedKey = await window.crypto.subtle.decrypt(
                {
                    name: "RSA-OAEP"
                },
                privateKey,
                encryptedKey
            );
            
            return Array.from(new Uint8Array(decryptedKey));
        } catch (error) {
            console.error('Error decrypting symmetric key:', error);
            throw new Error('Failed to decrypt symmetric key');
        }
    };
    
    // Get or create symmetric key for a chat
    const getSymmetricKey = async (chatId) => {
        // Check if we already have the key in memory
        if (symmetricKeys.value.has(chatId)) {
            return symmetricKeys.value.get(chatId);
        }
        
        // Check if we have it cached in localforage
        const cachedKey = await localforage.getItem(`chatKey_${chatId}`);
        if (cachedKey) {
            // Import the cached key
            const key = await window.crypto.subtle.importKey(
                "raw",
                new Uint8Array(cachedKey),
                {
                    name: "AES-GCM",
                    length: 256
                },
                false,
                ["encrypt", "decrypt"]
            );
            symmetricKeys.value.set(chatId, key);
            return key;
        }
        
        return null;
    };
    
    // Store symmetric key in memory and optionally cache it
    const storeSymmetricKey = async (chatId, symmetricKeyBytes, cacheInStorage = true) => {
        try {
            // Import the key
            const key = await window.crypto.subtle.importKey(
                "raw",
                new Uint8Array(symmetricKeyBytes),
                {
                    name: "AES-GCM",
                    length: 256
                },
                false,
                ["encrypt", "decrypt"]
            );
            
            // Store in memory
            symmetricKeys.value.set(chatId, key);
            
            // Optionally cache in localforage
            if (cacheInStorage) {
                await localforage.setItem(`chatKey_${chatId}`, symmetricKeyBytes);
            }
            
            return key;
        } catch (error) {
            console.error('Error storing symmetric key:', error);
            throw new Error('Failed to store symmetric key');
        }
    };
    
    // Encrypt a message with the symmetric key
    const encryptMessage = async (message, chatId) => {
        try {
            const key = await getSymmetricKey(chatId);
            if (!key) {
                throw new Error('No symmetric key found for chat');
            }
            
            // Generate a random IV
            const iv = window.crypto.getRandomValues(new Uint8Array(12));
            
            // Encrypt the message
            const encodedMessage = new TextEncoder().encode(message);
            const encryptedData = await window.crypto.subtle.encrypt(
                {
                    name: "AES-GCM",
                    iv: iv
                },
                key,
                encodedMessage
            );
            
            // Combine IV and encrypted data
            const combined = new Uint8Array(iv.length + encryptedData.byteLength);
            combined.set(iv);
            combined.set(new Uint8Array(encryptedData), iv.length);
            
            // Return base64 encoded
            return btoa(String.fromCharCode(...combined));
        } catch (error) {
            console.error('Error encrypting message:', error);
            throw new Error('Failed to encrypt message');
        }
    };
    
    // Decrypt a message with the symmetric key
    const decryptMessage = async (encryptedMessageBase64, chatId) => {
        try {
            const key = await getSymmetricKey(chatId);
            if (!key) {
                throw new Error('No symmetric key found for chat');
            }
            
            // Decode base64
            const combined = Uint8Array.from(atob(encryptedMessageBase64), c => c.charCodeAt(0));
            
            // Extract IV (first 12 bytes) and encrypted data
            const iv = combined.slice(0, 12);
            const encryptedData = combined.slice(12);
            
            // Decrypt the message
            const decryptedData = await window.crypto.subtle.decrypt(
                {
                    name: "AES-GCM",
                    iv: iv
                },
                key,
                encryptedData
            );
            
            // Convert back to string
            return new TextDecoder().decode(decryptedData);
        } catch (error) {
            console.error('Error decrypting message:', error);
            throw new Error('Failed to decrypt message');
        }
    };
    
    // Initialize E2EE for a new chat
    const initializeChatE2EE = async (chatId, participants) => {
        try {
            // Generate symmetric key
            const symmetricKey = await generateSymmetricKey();
            
            // Get our public key
            const ourPublicKey = await exportSecretChatPublicKey(chatId);
            if (!ourPublicKey) {
                throw new Error('No public key available');
            }
            
            // Encrypt symmetric key with our public key
            const encryptedKeyForUs = await encryptSymmetricKey(symmetricKey, ourPublicKey);
            
            // Store the key locally
            await storeSymmetricKey(chatId, symmetricKey);
            
            // Return encrypted key for backend storage
            return {
                chatId,
                encryptedKey: encryptedKeyForUs,
                participants
            };
        } catch (error) {
            console.error('Error initializing chat E2EE:', error);
            throw error;
        }
    };
    
    // Load symmetric key for an existing chat
    const loadChatSymmetricKey = async (chatId, encryptedKeyBase64) => {
        try {
            // Get our private key
            const privateKey = await getSecretChatPrivateKey(chatId);
            if (!privateKey) {
                throw new Error('No private key available');
            }
            
            // Decrypt the symmetric key
            const symmetricKey = await decryptSymmetricKey(encryptedKeyBase64, privateKey);
            
            // Store the key
            await storeSymmetricKey(chatId, symmetricKey);
            
            return true;
        } catch (error) {
            console.error('Error loading chat symmetric key:', error);
            throw error;
        }
    };
    
    // Clear symmetric key from memory
    const clearSymmetricKey = (chatId) => {
        symmetricKeys.value.delete(chatId);
    };
    
    // Clear all symmetric keys
    const clearAllSymmetricKeys = () => {
        symmetricKeys.value.clear();
    };
    
    return {
        generateSymmetricKey,
        encryptSymmetricKey,
        decryptSymmetricKey,
        getSymmetricKey,
        storeSymmetricKey,
        encryptMessage,
        decryptMessage,
        initializeChatE2EE,
        loadChatSymmetricKey,
        clearSymmetricKey,
        clearAllSymmetricKeys
    };
} 