<template>
    <section class="flex flex-col h-full w-full bg-gray-50 font-roboto">
        <!-- Chat Header -->
        <ChatHeader
            v-if="chatStore.currentChatUser && !isCurrentChatSecret"
            :user="chatStore.currentChatUser"
            :backend-base-url="backendBaseUrl"
            :is-secret-chat="false"
            @delete-chat="handleDeleteChat"
        />
        
        <!-- Secret Chat Header -->
        <SecretChatHeader
            v-if="currentSecretChat && isCurrentChatSecret"
            :secret-chat="currentSecretChat"
            :secret-usernames="chatStore.secretUsernames"
            :backend-base-url="backendBaseUrl"
            :current-user-id="userStore.user_id"
        />

        <!-- Messages Area -->
        <MessagesArea
            v-if="chatStore.currentChatUser"
            :messages="chatStore.messages"
            :current-user-id="userStore.user_id"
            :backend-base-url="backendBaseUrl"
            :user-avatar="userStore.avatar_url"
            :other-user-avatar="chatStore.currentChatUser.avatar_url"
            :chat-id="getCurrentChatId()"
            :is-secret-chat="isCurrentChatSecret"
            :is-secret-chat-approved="isSecretChatApproved"
            @load-more-messages="handleLoadMoreMessages"
        />

        <!-- Message Input -->
        <MessageInput
            v-if="chatStore.currentChatUser"
            v-model="newMessage"
            :is-secret-chat="isCurrentChatSecret"
            :is-secret-chat-approved="isSecretChatApproved"
            @send="sendMessage"
        />

        <!-- No Chat Selected State -->
        <NoChatSelected v-if="!chatStore.currentChatUser" />
    </section>
</template>

<script setup>
import { ref, watch, computed, nextTick } from "vue";
import { useChatStore } from "../stores/chat";
import { useUserStore } from "../stores/users";
import ChatHeader from "./chat/ChatHeader.vue";
import SecretChatHeader from "./chat/SecretChatHeader.vue";
import NoChatSelected from "./chat/NoChatSelected.vue";
import MessagesArea from "./chat/MessagesArea.vue";
import MessageInput from "./chat/MessageInput.vue";
import { useWebSocket } from "../composables/useWebSocket";
import { useMessagePagination } from "../composables/useMessagePagination";
import { useMessageDeletion } from "../composables/useMessageDeletion";
import { useE2EE } from "../composables/useE2EE";
import { useSecretChatEncryption } from "../composables/useSecretChatEncryption";
import axiosInstance from "../axiosInstance";
import { showMessage, showError } from "../utils/toast";

const chatStore = useChatStore();
const userStore = useUserStore();
const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;
const newMessage = ref("");

// Check if current chat is a secret chat
const isCurrentChatSecret = computed(() => {
    if (!chatStore.currentChatUser?.id) return false;
    
    // If the current chat user has a secret chat ID, it's a secret chat
    if (chatStore.currentChatUser.secret_chat_id) {
        return true;
    }
    
    // If no secret_chat_id is present, it's definitely not a secret chat
    return false;
});

// Get the current secret chat object
const currentSecretChat = computed(() => {
    if (!isCurrentChatSecret.value || !chatStore.currentChatUser?.secret_chat_id) return null;
    
    // Find the secret chat by its ID
    return chatStore.secretChats?.find(chat => chat.id === chatStore.currentChatUser.secret_chat_id) || null;
});

// Check if the current secret chat is approved
const isSecretChatApproved = computed(() => {
    if (!isCurrentChatSecret.value) return true; // Not a secret chat, so "approved"
    
    const secretChat = currentSecretChat.value;
    if (!secretChat) return false; // No secret chat found, not approved
    
    return secretChat.user_2_accepted === true;
});

// WebSocket management
const { establishConnection, sendMessage: sendWebSocketMessage, closeConnection, getConnectionStatus } =
    useWebSocket();

// Message pagination
const { loadNextPage, loadInitialMessages, loadInitialSecretChatMessages } = useMessagePagination();

// Message deletion
const { updateMessageId } = useMessageDeletion();

// E2EE
const { encryptMessage, decryptMessage, loadChatSymmetricKey } = useE2EE();
const { loadSecretChatSymmetricKey, validateSecretChatForEncryption, encryptMessageForSending } = useSecretChatEncryption();

// Watch for chat user changes to manage WebSocket connections
watch(
    () => chatStore.currentChatUser?.id,
    async (newUserId, oldUserId) => {
        console.log('🔄 Chat user changed:', { 
            oldUserId, 
            newUserId, 
            currentChatUser: chatStore.currentChatUser,
            isSecretChat: isCurrentChatSecret.value 
        });
        
        if (oldUserId) {
            // Close previous connection explicitly
            closeConnection();
        }

        if (newUserId) {
            // Add a small delay to ensure chat is properly added to store
            await nextTick();
            // Try to get chat data with retries
            let chatData = null;
            let retries = 0;
            const maxRetries = 3;
            while (!chatData && retries < maxRetries) {
                chatData = getChatData(newUserId);
                if (!chatData) {
                    await new Promise(resolve => setTimeout(resolve, 100)); // Wait 100ms
                    retries++;
                }
            }
            if (chatData) {
                establishConnection(chatData, handleIncomingMessage);
                // Wait a bit for connection to establish
                await new Promise(resolve => setTimeout(resolve, 200));
            }
        }
    }
);

// Watch for when currentChatUser becomes null to close WebSocket
watch(
    () => chatStore.currentChatUser,
    (newChatUser, oldChatUser) => {
        if (oldChatUser && !newChatUser) {
            console.log('🔄 Chat user cleared, closing WebSocket connection');
            closeConnection();
        }
    }
);

// Get chat data for WebSocket connection
const getChatData = (targetUserId) => {
    // Check if this is a secret chat
    if (chatStore.currentChatUser?.secret_chat_id) {
        const senderId = userStore.user_id;
        const receiverId = targetUserId;
        return {
            chatId: chatStore.currentChatUser.secret_chat_id,
            senderId,
            receiverId,
            backendBaseUrl,
            isSecretChat: true,
        };
    }

    let chat = null;
    // Try to find by chat_id if available in currentChatUser
    if (chatStore.currentChatUser?.chat_id) {
        chat = chatStore.chats?.find((c) => c.id === chatStore.currentChatUser.chat_id);
        if (chat) {
            const senderId = userStore.user_id;
            // Find the other participant
            const receiverId = chat.participants?.find((id) => id !== senderId) || targetUserId;
            return {
                chatId: chat.id,
                senderId,
                receiverId,
                backendBaseUrl,
                isSecretChat: false,
            };
        }
    }
    // Try to find by participants
    chat = chatStore.chats?.find(
        (c) =>
            c.participants &&
            c.participants.includes(targetUserId) &&
            c.participants.includes(userStore.user_id)
    );
    if (!chat) {
        return null;
    }
    if (!chat.participants || chat.participants.length < 2) {
        return null;
    }
    const senderId = userStore.user_id;
    const receiverId = targetUserId;
    return {
        chatId: chat.id,
        senderId,
        receiverId,
        backendBaseUrl,
        isSecretChat: false,
    };
};

// Handle incoming messages
const handleIncomingMessage = async (data) => {
    const message = parseIncomingMessage(data);

    // Decrypt message if this is a secret chat
    let decryptedContent = message.content;
    if (isCurrentChatSecret.value && message.content) {
        try {
            // For secret chats, use the secret_chat_id instead of chat_id
            const chatId = chatStore.currentChatUser?.secret_chat_id || message.chat_id;
            console.log('🔐 Decrypting incoming message for secret chat:', {
                message_chat_id: message.chat_id,
                secret_chat_id: chatStore.currentChatUser?.secret_chat_id,
                used_chat_id: chatId,
                content_length: message.content.length,
                is_secret: message.is_secret
            });
            
            // For WebSocket messages in secret chats, always try to decrypt if content looks encrypted
            // WebSocket messages might not have the is_secret flag, so we rely on the chat context
            const shouldDecrypt = message.content.length > 20 && 
                /^[A-Za-z0-9+/=]+$/.test(message.content) && // Base64 pattern
                message.content.length % 4 === 0; // Base64 length check
            
            console.log('🔍 Incoming WebSocket message encryption check:', {
                shouldDecrypt,
                length: message.content.length,
                isSecret: message.is_secret,
                isBase64: /^[A-Za-z0-9+/=]+$/.test(message.content),
                isSecretChat: isCurrentChatSecret.value
            });
            
            if (shouldDecrypt) {
                // Check if we have the symmetric key, if not try to load it
                const { hasSymmetricKey } = useE2EE();
                const { loadSecretChatSymmetricKey } = useSecretChatEncryption();
                const keyAvailable = await hasSymmetricKey(chatId);
                
                if (!keyAvailable) {
                    console.log('🔐 Symmetric key not found, attempting to load...');
                    await loadSecretChatSymmetricKey(chatId);
                }
                
                decryptedContent = await decryptMessage(message.content, chatId);
                console.log('✅ Successfully decrypted incoming WebSocket message');
            } else {
                console.log('🔍 Incoming WebSocket message appears to be plaintext or not encrypted');
            }
        } catch (error) {
            console.error('Error decrypting WebSocket message:', error);
            // If decryption fails, show encrypted content or error message
            decryptedContent = '[Encrypted message - decryption failed]';
        }
    }

    const decryptedMessage = {
        ...message,
        content: decryptedContent
    };

    // Check if this is a confirmation of a sent message (same content and sender)
    const existingMessage = chatStore.messages.find(
        (msg) =>
            msg.content === decryptedMessage.content &&
            msg.sender_id === decryptedMessage.sender_id &&
            msg.id &&
            msg.id.startsWith("temp-")
    );

    if (existingMessage && decryptedMessage.id && !decryptedMessage.id.startsWith("temp-")) {
        // Update the temp ID with the real ID from backend
        chatStore.updateMessageId(existingMessage.id, decryptedMessage.id);
    } else {
        // This is a new message from someone else
        chatStore.addMessage(decryptedMessage);
    }
};

// Parse incoming message data
const parseIncomingMessage = (data) => {
    // If the backend sends the message object directly
    if (data.content && typeof data.content === "object") {
        return data.content;
    }

    // If the backend sends content as a JSON string
    if (typeof data.content === "string" && data.content.startsWith("{")) {
        try {
            return JSON.parse(data.content);
        } catch (e) {
            return { content: data.content };
        }
    }

    // If the backend sends the message object as the root data
    if (data.chat_id && data.sender_id && data.content) {
        return data;
    }

    // Default case
    return data;
};

// Send message
const sendMessage = async () => {
    if (!newMessage.value.trim()) return;

    // Check if this is a secret chat that's not approved
    if (isCurrentChatSecret.value && !isSecretChatApproved.value) {
        showError('Cannot send messages in unapproved secret chat');
        return;
    }

    const targetUserId = chatStore.currentChatUser?.id;
    
    const chatData = getChatData(targetUserId);

    if (!chatData) {
        return;
    }

    // Check WebSocket connection status
    const connectionStatus = getConnectionStatus();
    
    if (!connectionStatus.isConnected) {
        establishConnection(chatData, handleIncomingMessage);
        
        // Wait for connection
        await new Promise(resolve => setTimeout(resolve, 500));
        
    }

    // Create temporary ID for immediate display
    const tempId = `temp-${Date.now()}-${Math.random()
        .toString(36)
        .substr(2, 9)}`;

    let messageContent = newMessage.value;
    
    // Encrypt message if this is a secret chat
    if (isCurrentChatSecret.value) {
        try {
            // Validate that the chat is ready for messaging
            const validation = await validateSecretChatForEncryption(chatData.chatId);
            if (!validation.valid) {
                showError(validation.message);
                return;
            }
            
            // Now encrypt the message
            messageContent = await encryptMessageForSending(newMessage.value, chatData.chatId);
        } catch (error) {
            console.error('Error encrypting message:', error);
            showError('Failed to encrypt message: ' + error.message);
            return;
        }
    }

    const messageData = {
        id: tempId,
        chat_id: chatData.chatId,
        sender_id: chatData.senderId,
        receiver_id: chatData.receiverId,
        content: messageContent,
        created_at: new Date().toISOString(),
    };

    // Add message to store immediately with temp ID (store decrypted content for display)
    chatStore.addMessage({
        ...messageData,
        content: newMessage.value, // Store decrypted content for display
    });

    // Send encrypted content via WebSocket
    const success = sendWebSocketMessage(messageContent);
    
    newMessage.value = "";
};

// Get current chat ID
const getCurrentChatId = () => {
    // Check if this is a secret chat
    if (chatStore.currentChatUser?.secret_chat_id) {
        return chatStore.currentChatUser.secret_chat_id;
    }
    
    const targetUserId = chatStore.currentChatUser?.id;
    const chat = chatStore.chats?.find(
        (c) =>
            c.participants &&
            c.participants.includes(targetUserId) &&
            c.participants.includes(userStore.user_id)
    );
    return chat?.id || null;
};

// Handle loading more messages
const handleLoadMoreMessages = async () => {
    const chatId = getCurrentChatId();
    if (chatId) {
        // Check if this is a secret chat
        if (chatStore.currentChatUser?.secret_chat_id) {
            // For secret chats, we might need a different approach for pagination
            // For now, we'll use the regular loadNextPage but with the secret chat ID
            await loadNextPage(chatId);
        } else {
            await loadNextPage(chatId);
        }
    }
};

// Handle chat deletion
const handleDeleteChat = async (user) => {
    try {
        // Find the chat ID for this user
        const chatId = getCurrentChatId();
        if (!chatId) {
            showError('No chat ID found for this user');
            return;
        }

        // Delete the chat from backend
        const response = await axiosInstance.delete(`/api/chat/delete/${chatId}`);
        
        if (response.status === 200) {
            // Remove chat from store
            const chatIndex = chatStore.chats.findIndex(chat => chat.id === chatId);
            if (chatIndex !== -1) {
                chatStore.chats.splice(chatIndex, 1);
            }
            
            // Clear current chat user
            chatStore.clearChatUser();
            // Close WebSocket connection
            closeConnection();
            
            showMessage('Chat deleted successfully');
        }
    } catch (error) {
        console.error('Error deleting chat:', error);
        showError('Failed to delete chat. Please try again.');
    }
};
</script>

<style scoped>
.font-roboto {
    font-family: "Roboto", Arial, sans-serif;
}
</style>
