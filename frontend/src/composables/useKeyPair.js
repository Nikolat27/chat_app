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

        // Store private key with the specific secret chat ID
        await localforage.setItem(`secretChatPrivateKey_${secretChatId}`, privateKey);
        
        // Base64 encode the public key for backend transmission
        const publicKeyString = JSON.stringify(publicKey);
        const base64PublicKey = btoa(publicKeyString);
        
        return base64PublicKey;
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
