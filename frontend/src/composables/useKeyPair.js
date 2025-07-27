import localforage from "localforage";

export async function generateKeyPair() {
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
    return publicKey;
}

export async function getPrivateKey() {
    return await localforage.getItem("secretChatPrivateKey");
}
