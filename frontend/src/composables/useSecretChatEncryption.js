import { ref, computed } from 'vue';
import { useChatStore } from '../stores/chat';
import { useUserStore } from '../stores/users';
import { useKeyPair } from './useKeyPair';
import axiosInstance from '../axiosInstance';

export function useSecretChatEncryption() {
    const chatStore = useChatStore();
    const userStore = useUserStore();
    const { 
        generateSecretChatKeyPair,
        exportSecretChatPublicKey,
        hasSecretChatKeys,
        encryptMessage,
        decryptMessage,
        getSecretChatPrivateKey
    } = useKeyPair();

    // Check if current chat is a secret chat
    const isCurrentChatSecret = () => {
        return chatStore.currentChatUser?.secret_chat_id;
    };

    // Get current secret chat data
    const getCurrentSecretChat = () => {
        if (!isCurrentChatSecret()) return null;
        return chatStore.secretChats?.find(chat => chat.id === chatStore.currentChatUser.secret_chat_id);
    };

    // Initialize encryption for a secret chat (generate keys and upload public key)
    const initializeSecretChatEncryption = async (secretChatId) => {
        try {
            console.log('🔐 Initializing encryption for secret chat:', secretChatId);
            
            // Generate key pair for this secret chat
            await generateSecretChatKeyPair(secretChatId);
            console.log('✅ Generated key pair');
            
            // Get our public key
            const publicKey = await exportSecretChatPublicKey(secretChatId);
            if (!publicKey) {
                throw new Error('Failed to export public key');
            }
            console.log('✅ Got public key');
            
            // Send public key to backend
            try {
                const response = await axiosInstance.post(`/api/secret-chat/add-symmetric-key/${secretChatId}`, {
                    public_key: publicKey
                });
                console.log('✅ Successfully uploaded public key, response:', response.data);
                
                // Refresh secret chat data from backend
                await refreshSecretChatData();
            } catch (uploadError) {
                if (uploadError.response?.status === 409) {
                    console.log('⚠️ Public key already exists, skipping upload');
                } else {
                    throw uploadError;
                }
            }
            
            return true;
        } catch (error) {
            console.error('❌ Error initializing secret chat encryption:', error);
            throw new Error(`Failed to initialize encryption: ${error.message}`);
        }
    };

    // Get recipient's public key for encryption
    const getRecipientPublicKey = (secretChat, currentUserId) => {
        if (!secretChat) return null;
        
        const isUser1 = secretChat.user_1 === currentUserId;
        const recipientKey = isUser1 ? secretChat.user_2_public_key : secretChat.user_1_public_key;
        
        console.log('🔍 Getting recipient public key:', {
            isUser1,
            currentUserId,
            user_1: secretChat.user_1,
            user_2: secretChat.user_2,
            user_1_public_key: secretChat.user_1_public_key,
            user_2_public_key: secretChat.user_2_public_key,
            recipientKey
        });
        
        return recipientKey;
    };

    // Get sender's public key for verification
    const getSenderPublicKey = (secretChat, senderId) => {
        if (!secretChat) return null;
        
        const isUser1 = secretChat.user_1 === senderId;
        return isUser1 ? secretChat.user_1_public_key : secretChat.user_2_public_key;
    };

    // Encrypt message for sending
    const encryptMessageForSending = async (message, recipientUserId) => {
        try {
            console.log('🔐 Encrypting message for sending');
            
            const secretChat = getCurrentSecretChat();
            if (!secretChat) {
                throw new Error('Secret chat not found');
            }
            
            const recipientPublicKey = getRecipientPublicKey(secretChat, userStore.user_id);
            if (!recipientPublicKey) {
                throw new Error('Recipient public key not available. Cannot send encrypted message.');
            }
            
            // Encrypt the message with recipient's public key
            const encryptedMessage = await encryptMessage(message, recipientPublicKey, secretChat.id);
            
            console.log('✅ Message encrypted for sending');
            return encryptedMessage;
        } catch (error) {
            console.error('❌ Error encrypting message for sending:', error);
            throw new Error(`Failed to encrypt message: ${error.message}`);
        }
    };

    // Decrypt incoming messages
    const decryptIncomingMessage = async (message) => {
        console.log('🔍 Starting decryption for message:', message);
        
        if (!message.chat_id) {
            console.log('❌ No chat_id found in message');
            return message;
        }

        // Check if this chat_id belongs to a secret chat
        const secretChat = chatStore.secretChats?.find(chat => chat.id === message.chat_id);
        if (!secretChat) {
            console.log('💬 This is a regular chat message, no decryption needed');
            return message;
        }

        try {
            console.log('🔐 Found secret chat, attempting to decrypt message');

            // Get sender's public key for decryption
            const senderPublicKey = getSenderPublicKey(secretChat, message.sender_id);
            if (!senderPublicKey) {
                console.log('❌ No sender public key found for this message');
                return message;
            }
            
            // Decrypt the message with sender's public key
            const decryptedContent = await decryptMessage(message.content, senderPublicKey, message.chat_id);
            console.log('✅ Successfully decrypted message content:', decryptedContent);

            return {
                ...message,
                content: decryptedContent
            };
        } catch (error) {
            console.error('❌ Error during message decryption:', error);
            return message;
        }
    };

    // Validate secret chat for encryption
    const validateSecretChatForEncryption = () => {
        if (!isCurrentChatSecret()) {
            return { valid: true, message: null };
        }

        const secretChat = getCurrentSecretChat();
        if (!secretChat) {
            return { valid: false, message: 'Secret chat not found' };
        }

        const currentUserId = userStore.user_id;
        const recipientPublicKey = getRecipientPublicKey(secretChat, currentUserId);
        
        if (!recipientPublicKey) {
            return { valid: false, message: 'Recipient public key not available. Cannot send encrypted message.' };
        }

        return { valid: true, message: null };
    };

    // Refresh secret chat data from backend
    const refreshSecretChatData = async () => {
        try {
            console.log('🔄 Refreshing secret chat data from backend...');
            const response = await axiosInstance.get('/api/user/get-secret-chats');
            chatStore.setSecretChats(response.data.secret_chats);
            console.log('✅ Secret chat data refreshed');
        } catch (error) {
            console.error('❌ Error refreshing secret chat data:', error);
        }
    };

    return {
        isCurrentChatSecret,
        getCurrentSecretChat,
        encryptMessageForSending,
        decryptIncomingMessage,
        validateSecretChatForEncryption,
        getRecipientPublicKey,
        getSenderPublicKey,
        initializeSecretChatEncryption,
        getSymmetricKey: null, // Remove this function
        refreshSecretChatData,
    };
} 