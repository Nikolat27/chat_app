import { ref, computed } from 'vue';
import { useChatStore } from '../stores/chat';
import { useUserStore } from '../stores/users';
import { useKeyPair } from './useKeyPair';
import { useE2EE } from './useE2EE';
import axiosInstance from '../axiosInstance';

export function useSecretChatEncryption() {
    const chatStore = useChatStore();
    const userStore = useUserStore();
    const { 
        generateSecretChatKeyPair,
        exportSecretChatPublicKey,
        hasSecretChatKeys
    } = useKeyPair();
    
    const {
        encryptMessage,
        decryptMessage,
        uploadPublicKey,
        getSecretChatData,
        uploadEncryptedSymmetricKeys,
        handleUserBApproval,
        loadSymmetricKeyForUserA,
        loadSymmetricKeyForUserB,
        isChatReadyForMessaging,
        hasSymmetricKey
    } = useE2EE();

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
                await uploadPublicKey(secretChatId, publicKey);
                console.log('✅ Successfully uploaded public key');
                
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

    // Encrypt message for sending using symmetric key
    const encryptMessageForSending = async (message, chatId) => {
        try {
            console.log('🔐 Encrypting message for sending using symmetric key');
            
            // Check if chat is ready for messaging
            const isReady = await isChatReadyForMessaging(chatId);
            if (!isReady) {
                throw new Error('Chat is not ready for messaging. Keys are not finalized yet.');
            }
            
            // Check if symmetric key is available
            let hasKey = await hasSymmetricKey(chatId);
            if (!hasKey) {
                console.log('⚠️ Symmetric key not in memory, attempting to load it...');
                // Try to load the symmetric key
                const loaded = await loadSecretChatSymmetricKey(chatId);
                if (loaded) {
                    hasKey = await hasSymmetricKey(chatId);
                }
            }
            
            if (!hasKey) {
                throw new Error('No symmetric key available for this chat. Please wait for keys to be finalized.');
            }
            
            // Encrypt the message with the symmetric key
            const encryptedMessage = await encryptMessage(message, chatId);
            
            console.log('✅ Message encrypted for sending');
            return encryptedMessage;
        } catch (error) {
            console.error('❌ Error encrypting message for sending:', error);
            throw new Error(`Failed to encrypt message: ${error.message}`);
        }
    };

    // Decrypt incoming messages using symmetric key
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
            console.log('🔐 Found secret chat, attempting to decrypt message with symmetric key');
            console.log('🔍 Message details:', {
                chat_id: message.chat_id,
                content_length: message.content?.length,
                is_encrypted: message.content?.includes('==') || message.content?.length > 100
            });

            // Check if symmetric key is available
            const hasKey = await hasSymmetricKey(message.chat_id);
            console.log('🔍 Symmetric key available:', hasKey);
            
            if (!hasKey) {
                console.log('⚠️ Symmetric key not available, attempting to load it...');
                const loaded = await loadSecretChatSymmetricKey(message.chat_id);
                console.log('🔍 Symmetric key loading result:', loaded);
            }

            // Decrypt the message with symmetric key
            const decryptedContent = await decryptMessage(message.content, message.chat_id);
            console.log('✅ Successfully decrypted message content:', decryptedContent);

            return {
                ...message,
                content: decryptedContent
            };
        } catch (error) {
            console.error('❌ Error during message decryption:', error);
            console.error('❌ Error details:', {
                message: error.message,
                stack: error.stack
            });
            return message;
        }
    };

    // Validate secret chat for encryption
    const validateSecretChatForEncryption = async (chatId) => {
        if (!isCurrentChatSecret()) {
            return { valid: true, message: null };
        }

        const secretChat = getCurrentSecretChat();
        if (!secretChat) {
            return { valid: false, message: 'Secret chat not found' };
        }

        // Check if chat is ready for messaging
        const isReady = await isChatReadyForMessaging(chatId);
        if (!isReady) {
            return { valid: false, message: 'Chat is not ready for messaging. Keys are not finalized yet.' };
        }

        return { valid: true, message: null };
    };

    // Load symmetric key for a secret chat
    const loadSecretChatSymmetricKey = async (secretChatId) => {
        try {
            console.log('🔐 Loading symmetric key for secret chat:', secretChatId);
            
            // Get secret chat data to determine which user we are
            const chatData = await getSecretChatData(secretChatId);
            const currentUserId = userStore.user_id;
            
            console.log('🔍 Chat data:', {
                chatId: secretChatId,
                currentUserId,
                user_1: chatData.user_1,
                user_2: chatData.user_2,
                key_finalized: chatData.key_finalized,
                user_1_encrypted_symmetric_key: !!chatData.user_1_encrypted_symmetric_key,
                user_2_encrypted_symmetric_key: !!chatData.user_2_encrypted_symmetric_key
            });
            
            // Determine if we are User A or User B
            const isUserA = chatData.user_1 === currentUserId;
            console.log('🔍 User role:', isUserA ? 'User A' : 'User B');
            
            if (isUserA) {
                // We are User A, load our symmetric key
                await loadSymmetricKeyForUserA(secretChatId);
                console.log('✅ Successfully loaded symmetric key for User A');
            } else {
                // We are User B, load our symmetric key
                await loadSymmetricKeyForUserB(secretChatId);
                console.log('✅ Successfully loaded symmetric key for User B');
            }
            
            return true;
        } catch (error) {
            console.error('❌ Error loading symmetric key for secret chat:', error);
            console.error('❌ Error details:', {
                message: error.message,
                stack: error.stack
            });
            // Don't throw error, just return false to indicate no symmetric key available
            return false;
        }
    };

    // Handle User B approval (generates symmetric key)
    const handleUserBApprovalForSecretChat = async (secretChatId) => {
        try {
            console.log('🔐 Handling User B approval for secret chat:', secretChatId);
            
            // Handle User B approval and symmetric key generation
            await handleUserBApproval(secretChatId);
            
            console.log('✅ Successfully handled User B approval');
            return true;
        } catch (error) {
            console.error('❌ Error handling User B approval:', error);
            throw error;
        }
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
        loadSecretChatSymmetricKey,
        handleUserBApprovalForSecretChat,
        refreshSecretChatData,
    };
} 