import localforage from "localforage";

export function useKeyPair() {
    const generateKeyPair = async () => {
        const keyPair = await window.crypto.subtle.generateKey(
            {
                name: "ECDH",
                namedCurve: "P-256",
            },
            true,
            ["deriveKey", "deriveBits"]
        );

        const publicKey = await window.crypto.subtle.exportKey(
            "jwk",
            keyPair.publicKey
        );
        const privateKey = await window.crypto.subtle.exportKey(
            "jwk",
            keyPair.privateKey
        );

        await localforage.setItem("secretChatPrivateKey", privateKey);
        
        // Base64 encode the public key for backend transmission
        const publicKeyString = JSON.stringify(publicKey);
        const base64PublicKey = btoa(publicKeyString);
        
        return base64PublicKey;
    };

    const generateSecretChatKeyPair = async (secretChatId) => {
        const keyPair = await window.crypto.subtle.generateKey(
            {
                name: "RSA-OAEP",
                modulusLength: 2048,
                publicExponent: new Uint8Array([1, 0, 1]),
                hash: "SHA-256",
            },
            true,
            ["encrypt", "decrypt"]
        );

        const publicKey = await window.crypto.subtle.exportKey(
            "jwk",
            keyPair.publicKey
        );
        const privateKey = await window.crypto.subtle.exportKey(
            "jwk",
            keyPair.privateKey
        );

        // Store both public and private keys with the specific secret chat ID
        await localforage.setItem(`secretChatPrivateKey_${secretChatId}`, privateKey);
        await localforage.setItem(`secretChatPublicKey_${secretChatId}`, publicKey);
        
        // Return the public key as JWK (not base64 encoded)
        return publicKey;
    };

    const getPrivateKey = async () => {
        return await localforage.getItem("secretChatPrivateKey");
    };

    const getSecretChatPrivateKey = async (secretChatId) => {
        return await localforage.getItem(`secretChatPrivateKey_${secretChatId}`);
    };

    const exportPublicKey = async () => {
        // This would be used to share your public key with the other user
        const privateKey = await getPrivateKey();
        if (privateKey) {
            // Convert private key back to CryptoKey to get the public key
            const cryptoKey = await window.crypto.subtle.importKey(
                "jwk",
                privateKey,
                {
                    name: "ECDH",
                    namedCurve: "P-256",
                },
                true,
                ["deriveKey", "deriveBits"]
            );
            
            // Get the public key from the key pair
            const publicKey = await window.crypto.subtle.exportKey(
                "jwk",
                cryptoKey.publicKey
            );
            
            // Base64 encode the public key for backend transmission
            const publicKeyString = JSON.stringify(publicKey);
            const base64PublicKey = btoa(publicKeyString);
            
            return base64PublicKey;
        }
        return null;
    };

    const exportSecretChatPublicKey = async (secretChatId) => {
        const privateKey = await getSecretChatPrivateKey(secretChatId);
        if (privateKey) {
            // For RSA, we need to store both public and private keys
            // Let me check if we have the public key stored
            const publicKey = await localforage.getItem(`secretChatPublicKey_${secretChatId}`);
            if (publicKey) {
                return publicKey;
            }
            
            // If we don't have the public key stored, we can't extract it from the private key
            // This means we need to regenerate the key pair
            console.warn('Public key not found, regenerating key pair for chat:', secretChatId);
            await generateSecretChatKeyPair(secretChatId);
            return await localforage.getItem(`secretChatPublicKey_${secretChatId}`);
        }
        return null;
    };

    // Encrypt a message with a public key (for direct RSA encryption)
    const encryptMessage = async (message, publicKeyJwk, chatId) => {
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
            
            // Encrypt the message
            const encodedMessage = new TextEncoder().encode(message);
            const encryptedData = await window.crypto.subtle.encrypt(
                {
                    name: "RSA-OAEP"
                },
                publicKey,
                encodedMessage
            );
            
            // Return base64 encoded encrypted message
            return btoa(String.fromCharCode(...new Uint8Array(encryptedData)));
        } catch (error) {
            console.error('Error encrypting message:', error);
            throw new Error('Failed to encrypt message');
        }
    };

    // Decrypt a message with private key (for direct RSA decryption)
    const decryptMessage = async (encryptedMessageBase64, privateKeyJwk, chatId) => {
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
            
            // Decode base64 encrypted message
            const encryptedMessage = Uint8Array.from(atob(encryptedMessageBase64), c => c.charCodeAt(0));
            
            // Decrypt the message
            const decryptedData = await window.crypto.subtle.decrypt(
                {
                    name: "RSA-OAEP"
                },
                privateKey,
                encryptedMessage
            );
            
            // Convert back to string
            return new TextDecoder().decode(decryptedData);
        } catch (error) {
            console.error('Error decrypting message:', error);
            throw new Error('Failed to decrypt message');
        }
    };

    const clearKeys = async () => {
        await localforage.removeItem("secretChatPrivateKey");
    };

    const clearSecretChatKeys = async (secretChatId) => {
        await localforage.removeItem(`secretChatPrivateKey_${secretChatId}`);
    };

    const clearAllKeys = async () => {
        // Clear all secret chat keys
        const keys = await localforage.keys();
        for (const key of keys) {
            if (key.startsWith('secretChatPrivateKey')) {
                await localforage.removeItem(key);
            }
        }
        // Also clear the general key
        await localforage.removeItem("secretChatPrivateKey");
        console.log("All secret chat keys cleared from LocalForage");
    };

    const hasSecretChatKeys = async (secretChatId) => {
        const privateKey = await getSecretChatPrivateKey(secretChatId);
        return !!privateKey;
    };

    return {
        generateKeyPair,
        generateSecretChatKeyPair,
        getPrivateKey,
        getSecretChatPrivateKey,
        exportPublicKey,
        exportSecretChatPublicKey,
        encryptMessage,
        decryptMessage,
        clearKeys,
        clearSecretChatKeys,
        clearAllKeys,
        hasSecretChatKeys,
    };
}

// Keep the individual exports for backward compatibility
export async function generateKeyPair() {
    const { generateKeyPair: genKeyPair } = useKeyPair();
    return await genKeyPair();
}

export async function getPrivateKey() {
    const { getPrivateKey: getKey } = useKeyPair();
    return await getKey();
}
