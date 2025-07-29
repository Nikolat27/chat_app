import { ref } from 'vue';
import localforage from 'localforage';
import { useKeyPair } from './useKeyPair';
import { useE2EEStore } from '../stores/e2ee';
import axiosInstance from '../axiosInstance';

export function useE2EE() {
    const { generateKeyPair, exportSecretChatPublicKey, getSecretChatPrivateKey } = useKeyPair();
    const e2eeStore = useE2EEStore();
    
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
        return e2eeStore.getSymmetricKey(chatId);
    };
    
    // Check if symmetric key is available for a chat
    const hasSymmetricKey = async (chatId) => {
        return e2eeStore.hasSymmetricKey(chatId);
    };
    
    // Store symmetric key in memory only (more secure)
    const storeSymmetricKey = async (chatId, symmetricKeyBytes) => {
        try {
            console.log('🔐 Storing symmetric key for chat:', chatId);
            
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
            
            // Store in global store
            e2eeStore.storeSymmetricKey(chatId, key);
            
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
            console.log('🔐 Starting message decryption for chat:', chatId);
            console.log('🔍 Encrypted message length:', encryptedMessageBase64?.length);
            
            const key = await getSymmetricKey(chatId);
            if (!key) {
                throw new Error('No symmetric key found for chat');
            }
            console.log('✅ Symmetric key found for decryption');
            
            // Decode base64
            console.log('🔐 Decoding base64 encrypted message...');
            const combined = Uint8Array.from(atob(encryptedMessageBase64), c => c.charCodeAt(0));
            console.log('🔍 Combined data length:', combined.length);
            
            // Extract IV (first 12 bytes) and encrypted data
            const iv = combined.slice(0, 12);
            const encryptedData = combined.slice(12);
            console.log('🔍 IV length:', iv.length, 'Encrypted data length:', encryptedData.length);
            
            // Decrypt the message
            console.log('🔐 Decrypting with AES-GCM...');
            const decryptedData = await window.crypto.subtle.decrypt(
                {
                    name: "AES-GCM",
                    iv: iv
                },
                key,
                encryptedData
            );
            
            // Convert back to string
            const result = new TextDecoder().decode(decryptedData);
            console.log('✅ Successfully decrypted message:', result);
            return result;
        } catch (error) {
            console.error('❌ Error decrypting message:', error);
            console.error('❌ Error details:', {
                message: error.message,
                stack: error.stack
            });
            throw new Error('Failed to decrypt message');
        }
    };
    
    // Upload public key to backend
    const uploadPublicKey = async (chatId, publicKeyJwk) => {
        try {
            console.log('🔐 Uploading public key for chat:', chatId);
            
            // Base64 encode the public key using a more robust method
            const publicKeyString = JSON.stringify(publicKeyJwk);
            const publicKeyBase64 = btoa(unescape(encodeURIComponent(publicKeyString)));
            
            const response = await axiosInstance.post(`/api/secret-chat/add-public-key/${chatId}`, {
                public_key: publicKeyBase64
            });
            
            console.log('✅ Successfully uploaded public key');
            return response.data;
        } catch (error) {
            console.error('Error uploading public key:', error);
            throw error;
        }
    };
    
    // Get secret chat data from backend
    const getSecretChatData = async (chatId) => {
        try {
            console.log('🔐 Getting secret chat data for chat:', chatId);
            
            const response = await axiosInstance.get(`/api/secret-chat/get/${chatId}`);
            
            console.log('✅ Successfully retrieved secret chat data');
            return response.data;
        } catch (error) {
            console.error('Error getting secret chat data:', error);
            throw error;
        }
    };
    
    // Upload encrypted symmetric keys to backend
    const uploadEncryptedSymmetricKeys = async (chatId, user1EncryptedKey, user2EncryptedKey) => {
        try {
            console.log('🔐 Uploading encrypted symmetric keys for chat:', chatId);
            
            const response = await axiosInstance.post(`/api/secret-chat/add-symmetric-key/${chatId}`, {
                user_1_encrypted_symmetric_key: user1EncryptedKey,
                user_2_encrypted_symmetric_key: user2EncryptedKey
            });
            
            console.log('✅ Successfully uploaded encrypted symmetric keys');
            return response.data;
        } catch (error) {
            console.error('Error uploading encrypted symmetric keys:', error);
            throw error;
        }
    };
    
    // Handle User B's approval and symmetric key generation
    const handleUserBApproval = async (chatId) => {
        try {
            console.log('🔐 Handling User B approval for chat:', chatId);
            
            // Get secret chat data to check if User A's public key is available
            const chatData = await getSecretChatData(chatId);
            
            if (!chatData.user_1_public_key) {
                throw new Error('User A\'s public key not available yet');
            }
            
            // Parse User A's public key
            const userAPublicKeyString = decodeURIComponent(escape(atob(chatData.user_1_public_key)));
            const userAPublicKeyJwk = JSON.parse(userAPublicKeyString);
            
            // Generate symmetric key
            const symmetricKey = await generateSymmetricKey();
            console.log('✅ Generated symmetric key');
            
            // Encrypt symmetric key for User A
            const userAEncryptedKey = await encryptSymmetricKey(symmetricKey, userAPublicKeyJwk);
            console.log('✅ Encrypted symmetric key for User A');
            
            // Get our public key and encrypt symmetric key for User B (ourselves)
            const ourPublicKey = await exportSecretChatPublicKey(chatId);
            if (!ourPublicKey) {
                throw new Error('No public key available');
            }
            
            const userBEncryptedKey = await encryptSymmetricKey(symmetricKey, ourPublicKey);
            console.log('✅ Encrypted symmetric key for User B');
            
            // Upload both encrypted keys to backend
            await uploadEncryptedSymmetricKeys(chatId, userAEncryptedKey, userBEncryptedKey);
            console.log('✅ Uploaded both encrypted symmetric keys to backend');
            
            // Store symmetric key locally for User B
            await storeSymmetricKey(chatId, symmetricKey);
            console.log('✅ Stored symmetric key locally for User B');
            
            return true;
        } catch (error) {
            console.error('Error handling User B approval:', error);
            throw error;
        }
    };
    
    // Load symmetric key for User A (first user)
    const loadSymmetricKeyForUserA = async (chatId) => {
        try {
            console.log('🔐 Loading symmetric key for User A in chat:', chatId);
            
            // Get secret chat data
            const chatData = await getSecretChatData(chatId);
            
            if (!chatData.user_1_encrypted_symmetric_key) {
                throw new Error('No encrypted symmetric key available for User A');
            }
            
            // Get our private key
            const privateKey = await getSecretChatPrivateKey(chatId);
            if (!privateKey) {
                throw new Error('No private key available');
            }
            
            // Decrypt the symmetric key
            const symmetricKey = await decryptSymmetricKey(chatData.user_1_encrypted_symmetric_key, privateKey);
            console.log('✅ Decrypted symmetric key for User A');
            
            // Store the key locally
            await storeSymmetricKey(chatId, symmetricKey);
            console.log('✅ Stored symmetric key locally for User A');
            
            return true;
        } catch (error) {
            console.error('Error loading symmetric key for User A:', error);
            throw error;
        }
    };
    
    // Load symmetric key for User B (second user)
    const loadSymmetricKeyForUserB = async (chatId) => {
        try {
            console.log('🔐 Loading symmetric key for User B in chat:', chatId);
            
            // Get secret chat data
            const chatData = await getSecretChatData(chatId);
            console.log('🔍 Chat data for User B:', {
                has_user_2_encrypted_symmetric_key: !!chatData.user_2_encrypted_symmetric_key,
                key_finalized: chatData.key_finalized
            });
            
            if (!chatData.user_2_encrypted_symmetric_key) {
                throw new Error('No encrypted symmetric key available for User B');
            }
            
            // Get our private key
            const privateKey = await getSecretChatPrivateKey(chatId);
            console.log('🔍 Private key available:', !!privateKey);
            if (!privateKey) {
                throw new Error('No private key available');
            }
            
            // Decrypt the symmetric key
            console.log('🔐 Decrypting symmetric key...');
            const symmetricKey = await decryptSymmetricKey(chatData.user_2_encrypted_symmetric_key, privateKey);
            console.log('✅ Decrypted symmetric key for User B');
            console.log('🔍 Symmetric key length:', symmetricKey.length);
            
            // Store the key locally
            console.log('🔐 Storing symmetric key...');
            await storeSymmetricKey(chatId, symmetricKey);
            console.log('✅ Stored symmetric key locally for User B');
            
            return true;
        } catch (error) {
            console.error('❌ Error loading symmetric key for User B:', error);
            console.error('❌ Error details:', {
                message: error.message,
                stack: error.stack
            });
            throw error;
        }
    };
    
    // Check if chat is ready for messaging (key finalized)
    const isChatReadyForMessaging = async (chatId) => {
        try {
            const chatData = await getSecretChatData(chatId);
            return chatData.key_finalized === true;
        } catch (error) {
            console.error('Error checking if chat is ready for messaging:', error);
            return false;
        }
    };
    
    // Clear symmetric key from memory
    const clearSymmetricKey = (chatId) => {
        e2eeStore.clearSymmetricKey(chatId);
    };
    
    // Clear all symmetric keys
    const clearAllSymmetricKeys = () => {
        e2eeStore.clearAllSymmetricKeys();
    };
    
    return {
        generateSymmetricKey,
        encryptSymmetricKey,
        decryptSymmetricKey,
        getSymmetricKey,
        hasSymmetricKey,
        storeSymmetricKey,
        encryptMessage,
        decryptMessage,
        uploadPublicKey,
        getSecretChatData,
        uploadEncryptedSymmetricKeys,
        handleUserBApproval,
        loadSymmetricKeyForUserA,
        loadSymmetricKeyForUserB,
        isChatReadyForMessaging,
        clearSymmetricKey,
        clearAllSymmetricKeys
    };
} 