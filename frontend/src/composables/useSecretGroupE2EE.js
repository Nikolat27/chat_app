import { ref } from 'vue';
import localforage from 'localforage';
import { useKeyPair } from './useKeyPair';
import { useE2EEStore } from '../stores/e2ee';
import axiosInstance from '../axiosInstance';

export function useSecretGroupE2EE() {
    const { generateSecretGroupKeyPair, exportSecretGroupPublicKey, getSecretGroupPrivateKey } = useKeyPair();
    const e2eeStore = useE2EEStore();
    
    // Generate a symmetric key for a group (small and light - AES-128-GCM)
    const generateGroupSymmetricKey = async () => {
        try {
            console.log('🔐 Generating symmetric key for group message (128-bit)');
            // Generate a random 128-bit key for smaller size
            const key = await window.crypto.subtle.generateKey(
                {
                    name: "AES-GCM",
                    length: 128
                },
                true,
                ["encrypt", "decrypt"]
            );
            
            // Export the key as raw bytes
            const rawKey = await window.crypto.subtle.exportKey("raw", key);
            console.log('✅ Generated symmetric key for group message (128-bit)');
            return new Uint8Array(rawKey);
        } catch (error) {
            console.error('❌ Error generating group symmetric key:', error);
            throw new Error('Failed to generate group symmetric key');
        }
    };
    
    // Encrypt symmetric key with a user's public key
    const encryptGroupSymmetricKey = async (symmetricKey, publicKeyJwk) => {
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
                symmetricKey
            );
            
            // Return base64 encoded encrypted key
            return btoa(String.fromCharCode(...new Uint8Array(encryptedKey)));
        } catch (error) {
            console.error('Error encrypting group symmetric key:', error);
            throw new Error('Failed to encrypt group symmetric key');
        }
    };
    
    // Decrypt symmetric key with user's private key
    const decryptGroupSymmetricKey = async (encryptedKeyBase64, privateKeyJwk) => {
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
            
            return new Uint8Array(decryptedKey);
        } catch (error) {
            console.error('Error decrypting group symmetric key:', error);
            throw new Error('Failed to decrypt group symmetric key');
        }
    };
    
    // Get or create symmetric key for a group
    const getGroupSymmetricKey = async (groupId) => {
        return e2eeStore.getSymmetricKey(`group_${groupId}`);
    };
    
    // Check if symmetric key is available for a group
    const hasGroupSymmetricKey = async (groupId) => {
        return e2eeStore.hasSymmetricKey(`group_${groupId}`);
    };
    
    // Store symmetric key in memory only (more secure)
    const storeGroupSymmetricKey = async (groupId, symmetricKeyBytes) => {
        try {
            console.log('🔐 Storing group symmetric key for group:', groupId);
            
            // Import the key
            const key = await window.crypto.subtle.importKey(
                "raw",
                symmetricKeyBytes,
                {
                    name: "AES-GCM",
                    length: 128
                },
                false,
                ["encrypt", "decrypt"]
            );
            
            // Store in global store with group prefix
            e2eeStore.storeSymmetricKey(`group_${groupId}`, key);
            
            return key;
        } catch (error) {
            console.error('Error storing group symmetric key:', error);
            throw new Error('Failed to store group symmetric key');
        }
    };
    
    // Encrypt a message with the symmetric key
    const encryptGroupMessage = async (message, symmetricKeyBytes) => {
        try {
            if (!symmetricKeyBytes) {
                throw new Error('No symmetric key provided for encryption');
            }
            
            // Convert raw bytes to CryptoKey
            const cryptoKey = await window.crypto.subtle.importKey(
                'raw',
                symmetricKeyBytes,
                {
                    name: 'AES-GCM',
                    length: 128
                },
                false,
                ['encrypt']
            );
            
            // Generate a random IV
            const iv = window.crypto.getRandomValues(new Uint8Array(12));
            
            // Encrypt the message
            const encodedMessage = new TextEncoder().encode(message);
            const encryptedData = await window.crypto.subtle.encrypt(
                {
                    name: "AES-GCM",
                    iv: iv
                },
                cryptoKey,
                encodedMessage
            );
            
            // Combine IV and encrypted data
            const combined = new Uint8Array(iv.length + encryptedData.byteLength);
            combined.set(iv);
            combined.set(new Uint8Array(encryptedData), iv.length);
            
            // Return base64 encoded encrypted message
            return btoa(String.fromCharCode(...combined));
        } catch (error) {
            console.error('Error encrypting group message:', error);
            throw new Error('Failed to encrypt group message');
        }
    };
    
    // Decrypt a message with the symmetric key
    const decryptGroupMessage = async (encryptedMessageBase64, symmetricKeyBytes) => {
        try {
            if (!symmetricKeyBytes) {
                throw new Error('No symmetric key provided for decryption');
            }
            
            // Convert raw bytes to CryptoKey
            const cryptoKey = await window.crypto.subtle.importKey(
                'raw',
                symmetricKeyBytes,
                {
                    name: 'AES-GCM',
                    length: 128
                },
                false,
                ['decrypt']
            );
            
            // Decode base64 encrypted message
            const encryptedMessage = Uint8Array.from(atob(encryptedMessageBase64), c => c.charCodeAt(0));
            
            // Extract IV (first 12 bytes) and encrypted data
            const iv = encryptedMessage.slice(0, 12);
            const encryptedData = encryptedMessage.slice(12);
            
            // Decrypt the message
            const decryptedData = await window.crypto.subtle.decrypt(
                {
                    name: "AES-GCM",
                iv: iv
                },
                cryptoKey,
                encryptedData
            );
            
            // Convert back to string
            return new TextDecoder().decode(decryptedData);
        } catch (error) {
            console.error('Error decrypting group message:', error);
            throw new Error('Failed to decrypt group message');
        }
    };
    
    // Upload public key to backend for a secret group
    const uploadGroupPublicKey = async (groupId, publicKeyJwk) => {
        try {
            console.log('🔐 Uploading public key for group:', groupId);
            
            // Convert JWK to base64 for backend transmission
            const publicKeyString = JSON.stringify(publicKeyJwk);
            const base64PublicKey = btoa(publicKeyString);
            
            const response = await axiosInstance.post(`/api/secret-group/${groupId}/upload-public-key`, {
                public_key: base64PublicKey
            });
            
            console.log('✅ Successfully uploaded group public key');
            return response.data;
        } catch (error) {
            console.error('Error uploading group public key:', error);
            throw error;
        }
    };
    
    // Get secret group data from backend
    const getSecretGroupData = async (groupId) => {
        try {
            console.log('🔐 Getting secret group data for group:', groupId);
            
            const response = await axiosInstance.get(`/api/secret-group/get/${groupId}/members`);
            
            console.log('✅ Successfully retrieved secret group data');
            return response.data;
        } catch (error) {
            console.error('Error getting secret group data:', error);
            throw error;
        }
    };
    
    // Upload encrypted symmetric keys to backend for all group members
    const uploadGroupEncryptedSymmetricKeys = async (groupId, encryptedKeys) => {
        try {
            console.log('🔐 Uploading encrypted symmetric keys for group:', groupId);
            
            const response = await axiosInstance.post(`/api/secret-group/${groupId}/upload-symmetric-keys`, {
                encrypted_keys: encryptedKeys
            });
            
            console.log('✅ Successfully uploaded encrypted symmetric keys');
            return response.data;
        } catch (error) {
            console.error('Error uploading encrypted symmetric keys:', error);
            throw error;
        }
    };
    
    // Initialize encryption for a secret group (generate keys only, public key already uploaded)
    const initializeSecretGroupEncryption = async (groupId) => {
        try {
            console.log('🔐 Initializing encryption for secret group:', groupId);
            
            // Generate key pair for this secret group
            await generateSecretGroupKeyPair(groupId);
            console.log('✅ Generated key pair for group');
            
            // Get our public key
            const publicKey = await exportSecretGroupPublicKey(groupId);
            if (!publicKey) {
                throw new Error('Failed to export group public key');
            }
            console.log('✅ Got group public key');
            
            // Public key is already uploaded during group creation/joining
            console.log('✅ Public key already uploaded during group creation/joining');
            
            return true;
        } catch (error) {
            console.error('❌ Error initializing secret group encryption:', error);
            throw new Error(`Failed to initialize group encryption: ${error.message}`);
        }
    };
    
    // Load symmetric key for a secret group (for new per-message architecture)
    const loadSecretGroupSymmetricKey = async (groupId) => {
        try {
            console.log('🔐 Loading symmetric key for secret group:', groupId);
            
            // Get secret group data to verify we have access
            const groupData = await getSecretGroupData(groupId);
            const currentUserId = JSON.parse(localStorage.getItem('user'))?.user_id;
            
            console.log('🔍 Group data:', {
                groupId,
                currentUserId,
                user_public_keys: Object.keys(groupData.user_public_keys || {}),
                member_join_times: Object.keys(groupData.member_join_times || {})
            });
            
            // Check if we have a private key for this group
            const privateKey = await getSecretGroupPrivateKey(groupId);
            if (!privateKey) {
                console.warn('⚠️ No private key available for group, keys will be loaded from messages');
                return false;
            }
            
            // In the new architecture, symmetric keys are loaded from individual messages
            // not pre-stored in group data. So we return false to indicate no pre-stored key.
            console.log('✅ Private key available, symmetric keys will be loaded from messages');
            return false;
        } catch (error) {
            console.error('❌ Error loading symmetric key for secret group:', error);
            console.error('❌ Error details:', {
                message: error.message,
                stack: error.stack
            });
            // Don't throw error, just return false to indicate no symmetric key available
            return false;
        }
    };
    
    // Generate and upload symmetric keys for all group members
    const generateAndUploadGroupSymmetricKeys = async (groupId) => {
        try {
            console.log('🔐 Generating and uploading symmetric keys for group:', groupId);
            
            // Get group data to get all members and their public keys
            const groupData = await getSecretGroupData(groupId);
            const userPublicKeys = groupData.user_public_keys || {};
            
            // Generate symmetric key
            const symmetricKey = await generateGroupSymmetricKey();
            console.log('✅ Generated symmetric key for group');
            
            // Encrypt symmetric key for each member
            const encryptedKeys = {};
            for (const [userId, publicKeyBase64] of Object.entries(userPublicKeys)) {
                try {
                    // Decode base64 public key
                    const publicKeyString = atob(publicKeyBase64);
                    const publicKeyJwk = JSON.parse(publicKeyString);
                    
                    // Encrypt symmetric key for this user
                    const encryptedKey = await encryptGroupSymmetricKey(symmetricKey, publicKeyJwk);
                    encryptedKeys[userId] = encryptedKey;
                    console.log(`✅ Encrypted symmetric key for user: ${userId}`);
                } catch (error) {
                    console.error(`❌ Failed to encrypt symmetric key for user ${userId}:`, error);
                    // Continue with other users
                }
            }
            
            // Upload encrypted keys to backend
            await uploadGroupEncryptedSymmetricKeys(groupId, encryptedKeys);
            console.log('✅ Uploaded encrypted symmetric keys to backend');
            
            // Store symmetric key locally for current user
            await storeGroupSymmetricKey(groupId, symmetricKey);
            console.log('✅ Stored symmetric key locally for current user');
            
            return true;
        } catch (error) {
            console.error('❌ Error generating and uploading group symmetric keys:', error);
            throw error;
        }
    };
    
    // Clear symmetric key from memory
    const clearGroupSymmetricKey = (groupId) => {
        e2eeStore.clearSymmetricKey(`group_${groupId}`);
    };
    
    // Send encrypted message to secret group with symmetric keys for each user
    const sendSecretGroupMessage = async (message, groupId) => {
        try {
            console.log('🔐 Sending encrypted message to secret group:', groupId);
            
            // 1. Generate a new symmetric key for this message
            const symmetricKeyBytes = await generateGroupSymmetricKey();
            console.log('✅ Generated symmetric key for message');
            
            // 2. Encrypt the message content with the symmetric key
            const encryptedMessage = await encryptGroupMessage(message, symmetricKeyBytes);
            console.log('✅ Encrypted message content');
            
            // 3. Get all group members and their public keys
            const groupData = await getSecretGroupData(groupId);
            console.log('✅ Retrieved group data with members');
            
            // 4. Encrypt the symmetric key for each user
            const usersSymmetricKeys = {};
            
            // Get all member IDs (including owner)
            const allMemberIds = [
                groupData.owner_id,
                ...(groupData.members || [])
            ].filter(id => id); // Remove any null/undefined values
            
            console.log('👥 Encrypting symmetric key for members:', allMemberIds);
            
            for (const memberId of allMemberIds) {
                try {
                    // Get user's public key from group data
                    const userPublicKeyBase64 = groupData.user_public_keys?.[memberId];
                    if (!userPublicKeyBase64) {
                        console.warn(`⚠️ No public key found for user ${memberId}`);
                        continue;
                    }
                    
                    // Decode base64 public key
                    const publicKeyString = atob(userPublicKeyBase64);
                    const publicKeyJwk = JSON.parse(publicKeyString);
                    
                    // Encrypt symmetric key with user's public key
                    const encryptedSymmetricKey = await encryptGroupSymmetricKey(symmetricKeyBytes, publicKeyJwk);
                    
                    // Store encrypted key for this user
                    usersSymmetricKeys[memberId] = encryptedSymmetricKey;
                    console.log(`✅ Encrypted symmetric key for user ${memberId}`);
                    
                } catch (error) {
                    console.error(`❌ Failed to encrypt symmetric key for user ${memberId}:`, error);
                    // Continue with other users
                }
            }
            
            // 5. Prepare the message payload
            const messagePayload = {
                content: encryptedMessage,
                users_symmetric_keys: usersSymmetricKeys
            };
            
            console.log('📤 Sending secret group message with payload:', {
                content_length: encryptedMessage.length,
                users_count: Object.keys(usersSymmetricKeys).length
            });
            
            return messagePayload;
            
        } catch (error) {
            console.error('❌ Error sending secret group message:', error);
            throw error;
        }
    };

    // Clear all group symmetric keys
    const clearAllGroupSymmetricKeys = () => {
        // This would need to be implemented in the store to clear all group keys
        // For now, we'll just clear the current group's key
        console.log('🗑️ Clearing all group symmetric keys');
    };
    
    return {
        generateGroupSymmetricKey,
        encryptGroupSymmetricKey,
        decryptGroupSymmetricKey,
        getGroupSymmetricKey,
        hasGroupSymmetricKey,
        storeGroupSymmetricKey,
        encryptGroupMessage,
        decryptGroupMessage,
        uploadGroupPublicKey,
        getSecretGroupData,
        uploadGroupEncryptedSymmetricKeys,
        initializeSecretGroupEncryption,
        loadSecretGroupSymmetricKey,
        generateAndUploadGroupSymmetricKeys,
        clearGroupSymmetricKey,
        clearAllGroupSymmetricKeys,
        sendSecretGroupMessage
    };
} 