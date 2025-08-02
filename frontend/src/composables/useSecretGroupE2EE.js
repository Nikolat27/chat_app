import { ref } from 'vue';
import localforage from 'localforage';
import { useKeyPair } from './useKeyPair';
import { useE2EEStore } from '../stores/e2ee';
import axiosInstance from '../axiosInstance';

export function useSecretGroupE2EE() {
    const e2eeStore = useE2EEStore();
    
    // Generate a simple secret key for a group
    const generateGroupSecretKey = async () => {
        try {
            console.log('🔐 Generating secret key for group');
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
            const keyBase64 = btoa(String.fromCharCode(...new Uint8Array(rawKey)));
            
            console.log('✅ Generated secret key for group');
            return keyBase64;
        } catch (error) {
            console.error('❌ Error generating group secret key:', error);
            throw new Error('Failed to generate group secret key');
        }
    };
    
    // Store secret key in localForage
    const storeGroupSecretKey = async (groupId, secretKey) => {
        try {
            console.log('🔐 Storing secret key for group:', groupId);
            console.log('🔐 Secret key length:', secretKey?.length);
            await localforage.setItem(`secret_group_key_${groupId}`, secretKey);
            console.log('✅ Secret key stored in localForage');
            return true;
        } catch (error) {
            console.error('❌ Error storing group secret key:', error);
            throw error;
        }
    };
    
    // Get secret key from localForage
    const getGroupSecretKey = async (groupId) => {
        try {
            console.log('🔐 Getting secret key for group:', groupId);
            const secretKey = await localforage.getItem(`secret_group_key_${groupId}`);
            console.log('🔐 Retrieved secret key from localForage:', secretKey ? 'exists' : 'not found');
            console.log('✅ Secret key retrieved from localForage');
            return secretKey;
        } catch (error) {
            console.error('❌ Error getting group secret key:', error);
            return null;
        }
    };
    
    // Check if secret key exists for a group
    const hasGroupSecretKey = async (groupId) => {
        const secretKey = await getGroupSecretKey(groupId);
        return !!secretKey;
    };
    
    // Encrypt a message with the secret key
    const encryptGroupMessage = async (message, secretKeyBase64) => {
        try {
            if (!secretKeyBase64) {
                throw new Error('No secret key provided for encryption');
            }
            
            // Convert base64 key to CryptoKey
            const keyBytes = Uint8Array.from(atob(secretKeyBase64), c => c.charCodeAt(0));
            const cryptoKey = await window.crypto.subtle.importKey(
                'raw',
                keyBytes,
                {
                    name: 'AES-GCM',
                    length: 256
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
    
    // Decrypt a message with the secret key
    const decryptGroupMessage = async (encryptedMessageBase64, secretKeyBase64) => {
        try {
            if (!secretKeyBase64) {
                throw new Error('No secret key provided for decryption');
            }
            
            // Convert base64 key to CryptoKey
            const keyBytes = Uint8Array.from(atob(secretKeyBase64), c => c.charCodeAt(0));
            const cryptoKey = await window.crypto.subtle.importKey(
                'raw',
                keyBytes,
                {
                    name: 'AES-GCM',
                    length: 256
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
    
    // Initialize secret group with a new secret key
    const initializeSecretGroup = async (groupId) => {
        try {
            console.log('🔐 Initializing secret group:', groupId);
            
            // Generate a new secret key
            const secretKey = await generateGroupSecretKey();
            console.log('🔐 Generated secret key for group:', groupId);
            
            // Store the secret key
            await storeGroupSecretKey(groupId, secretKey);
            console.log('🔐 Stored secret key for group:', groupId);
            
            console.log('✅ Secret group initialized with key');
            return secretKey;
        } catch (error) {
            console.error('❌ Error initializing secret group:', error);
            throw error;
        }
    };
    
    // Send encrypted message to secret group
    const sendSecretGroupMessage = async (message, groupId) => {
        try {
            console.log('🔐 Sending encrypted message to secret group:', groupId);
            
            // Get the secret key for this group
            const secretKey = await getGroupSecretKey(groupId);
            if (!secretKey) {
                throw new Error('No secret key available for this group');
        }
    
            // Encrypt the message
            const encryptedMessage = await encryptGroupMessage(message, secretKey);
            console.log('✅ Message encrypted successfully');
            
            // Return the encrypted message
            return {
                encrypted_content: encryptedMessage
            };
            
        } catch (error) {
            console.error('❌ Error sending secret group message:', error);
            throw error;
        }
    };
    
    // Copy secret key to clipboard
    const copySecretKeyToClipboard = async (groupId) => {
        try {
            const secretKey = await getGroupSecretKey(groupId);
            if (!secretKey) {
                throw new Error('No secret key available for this group');
            }
            
            await navigator.clipboard.writeText(secretKey);
            console.log('✅ Secret key copied to clipboard');
            return true;
        } catch (error) {
            console.error('❌ Error copying secret key:', error);
            throw error;
        }
    };
    
    // Enter secret key for a group
    const enterSecretKey = async (groupId, secretKey) => {
        try {
            console.log('🔐 Entering secret key for group:', groupId);
            
            // Validate the key format (should be base64)
            if (!/^[A-Za-z0-9+/=]+$/.test(secretKey)) {
                throw new Error('Invalid secret key format');
            }
            
            // Test the key by trying to decrypt a test message
            const testMessage = "test";
            const encryptedTest = await encryptGroupMessage(testMessage, secretKey);
            const decryptedTest = await decryptGroupMessage(encryptedTest, secretKey);
            
            if (decryptedTest !== testMessage) {
                throw new Error('Invalid secret key');
            }
            
            // Store the valid key
            await storeGroupSecretKey(groupId, secretKey);
            console.log('✅ Secret key entered and validated');
            return true;
        } catch (error) {
            console.error('❌ Error entering secret key:', error);
            throw error;
        }
    };
    
    // Clear secret key for a group
    const clearGroupSecretKey = async (groupId) => {
        try {
            console.log('🔐 Clearing secret key for group:', groupId);
            await localforage.removeItem(`secret_group_key_${groupId}`);
            console.log('✅ Secret key cleared from localForage');
            return true;
        } catch (error) {
            console.error('❌ Error clearing group secret key:', error);
            throw error;
        }
    };
    
    return {
        generateGroupSecretKey,
        storeGroupSecretKey,
        getGroupSecretKey,
        hasGroupSecretKey,
        encryptGroupMessage,
        decryptGroupMessage,
        initializeSecretGroup,
        sendSecretGroupMessage,
        copySecretKeyToClipboard,
        enterSecretKey,
        clearGroupSecretKey
    };
} 