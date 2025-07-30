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

    // Generate key pair for a secret group
    const generateSecretGroupKeyPair = async (groupId) => {
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

        // Store both public and private keys with the specific group ID
        await localforage.setItem(`secretGroupPrivateKey_${groupId}`, privateKey);
        await localforage.setItem(`secretGroupPublicKey_${groupId}`, publicKey);
        
        // Return the public key as JWK (not base64 encoded)
        return publicKey;
    };

    const getPrivateKey = async () => {
        return await localforage.getItem("secretChatPrivateKey");
    };

    const getSecretChatPrivateKey = async (secretChatId) => {
        return await localforage.getItem(`secretChatPrivateKey_${secretChatId}`);
    };

    // Get private key for a secret group
    const getSecretGroupPrivateKey = async (groupId) => {
        return await localforage.getItem(`secretGroupPrivateKey_${groupId}`);
    };

    const exportSecretChatPublicKey = async (secretChatId) => {
        return await localforage.getItem(`secretChatPublicKey_${secretChatId}`);
    };

    // Export public key for a secret group
    const exportSecretGroupPublicKey = async (groupId) => {
        return await localforage.getItem(`secretGroupPublicKey_${groupId}`);
    };

    const hasSecretChatKeys = async (secretChatId) => {
        const privateKey = await localforage.getItem(`secretChatPrivateKey_${secretChatId}`);
        const publicKey = await localforage.getItem(`secretChatPublicKey_${secretChatId}`);
        return !!(privateKey && publicKey);
    };

    // Check if secret group keys exist
    const hasSecretGroupKeys = async (groupId) => {
        const privateKey = await localforage.getItem(`secretGroupPrivateKey_${groupId}`);
        const publicKey = await localforage.getItem(`secretGroupPublicKey_${groupId}`);
        return !!(privateKey && publicKey);
    };

    const clearSecretChatKeys = async (secretChatId) => {
        await localforage.removeItem(`secretChatPrivateKey_${secretChatId}`);
        await localforage.removeItem(`secretChatPublicKey_${secretChatId}`);
    };

    // Clear secret group keys
    const clearSecretGroupKeys = async (groupId) => {
        await localforage.removeItem(`secretGroupPrivateKey_${groupId}`);
        await localforage.removeItem(`secretGroupPublicKey_${groupId}`);
    };

    return {
        generateKeyPair,
        generateSecretChatKeyPair,
        generateSecretGroupKeyPair,
        getPrivateKey,
        getSecretChatPrivateKey,
        getSecretGroupPrivateKey,
        exportSecretChatPublicKey,
        exportSecretGroupPublicKey,
        hasSecretChatKeys,
        hasSecretGroupKeys,
        clearSecretChatKeys,
        clearSecretGroupKeys
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
