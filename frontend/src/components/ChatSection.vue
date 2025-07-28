<template>
    <section class="flex flex-col h-full w-full bg-gray-50 font-roboto">
        <!-- Chat Header -->
        <ChatHeader
            v-if="chatStore.currentChatUser && !isCurrentChatSecret"
            :user="chatStore.currentChatUser"
            :backend-base-url="backendBaseUrl"
            :is-secret-chat="false"
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
            @load-more-messages="handleLoadMoreMessages"
        />

        <!-- Message Input -->
        <MessageInput
            v-if="chatStore.currentChatUser"
            v-model="newMessage"
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
    
    const currentUserId = userStore.user_id;
    const targetUserId = chatStore.currentChatUser.id;
    
    // Check if this is a secret chat by looking in secret chats
    return chatStore.secretChats?.some(chat => 
        (chat.user_1 === currentUserId && chat.user_2 === targetUserId) ||
        (chat.user_2 === currentUserId && chat.user_1 === targetUserId)
    ) || false;
});

// Get the current secret chat object
const currentSecretChat = computed(() => {
    if (!isCurrentChatSecret.value || !chatStore.currentChatUser?.id) return null;
    
    const currentUserId = userStore.user_id;
    const targetUserId = chatStore.currentChatUser.id;
    
    return chatStore.secretChats?.find(chat => 
        (chat.user_1 === currentUserId && chat.user_2 === targetUserId) ||
        (chat.user_2 === currentUserId && chat.user_1 === targetUserId)
    ) || null;
});

// WebSocket management
const { establishConnection, sendMessage: sendWebSocketMessage, closeConnection, getConnectionStatus } =
    useWebSocket();

// Message pagination
const { loadNextPage, loadInitialMessages } = useMessagePagination();

// Message deletion
const { updateMessageId } = useMessageDeletion();

// Watch for chat user changes to manage WebSocket connections
watch(
    () => chatStore.currentChatUser?.id,
    async (newUserId, oldUserId) => {
        if (oldUserId) {
            // Close previous connection explicitly
            closeConnection();
            console.log("Closed previous WebSocket connection for user:", oldUserId);
        }

        if (newUserId) {
            console.log("Setting up WebSocket for new user:", newUserId);
            
            // Add a small delay to ensure chat is properly added to store
            await nextTick();
            
            // Try to get chat data with retries
            let chatData = null;
            let retries = 0;
            const maxRetries = 3;
            
            while (!chatData && retries < maxRetries) {
                chatData = getChatData(newUserId);
                if (!chatData) {
                    console.log(`Chat data not found, retry ${retries + 1}/${maxRetries}`);
                    await new Promise(resolve => setTimeout(resolve, 100)); // Wait 100ms
                    retries++;
                }
            }
            
            if (chatData) {
                console.log("Establishing WebSocket connection for:", chatData);
                establishConnection(chatData, handleIncomingMessage);
                
                // Wait a bit for connection to establish
                await new Promise(resolve => setTimeout(resolve, 200));
                console.log("WebSocket connection status:", getConnectionStatus());
            } else {
                console.log("No chat data found for user after retries:", newUserId);
                console.log("Available chats:", chatStore.chats);
            }
        }
    }
);

// Get chat data for WebSocket connection
const getChatData = (targetUserId) => {
    console.log("Looking for chat data for user:", targetUserId);
    console.log("Current user:", userStore.user_id);
    console.log("All chats in store:", chatStore.chats);
    console.log("Current chat user:", chatStore.currentChatUser);
    
    let chat = null;
    
    // Try to find by chat_id if available in currentChatUser
    if (chatStore.currentChatUser?.chat_id) {
        console.log("Trying to find chat by chat_id:", chatStore.currentChatUser.chat_id);
        chat = chatStore.chats?.find((c) => c.id === chatStore.currentChatUser.chat_id);
        if (chat) {
            console.log("Found chat by chat_id:", chat);
            const senderId = userStore.user_id;
            // Find the other participant
            const receiverId = chat.participants?.find((id) => id !== senderId) || targetUserId;
            return {
                chatId: chat.id,
                senderId,
                receiverId,
                backendBaseUrl,
            };
        }
    }

    // Try to find by participants
    console.log("Trying to find chat by participants");
    chat = chatStore.chats?.find(
        (c) => {
            console.log("Checking chat:", c);
            console.log("Chat participants:", c.participants);
            console.log("Target user:", targetUserId, "Current user:", userStore.user_id);
            
            return c.participants &&
                   c.participants.includes(targetUserId) &&
                   c.participants.includes(userStore.user_id);
        }
    );

    if (!chat) {
        console.log("Chat not found in store");
        return null;
    }

    if (!chat.participants || chat.participants.length < 2) {
        console.log("Invalid chat structure:", chat);
        return null;
    }

    const senderId = userStore.user_id;
    const receiverId = targetUserId;

    console.log("Found chat for WebSocket:", {
        chatId: chat.id,
        senderId,
        receiverId,
        chat
    });

    return {
        chatId: chat.id,
        senderId,
        receiverId,
        backendBaseUrl,
    };
};

// Handle incoming messages
const handleIncomingMessage = (data) => {
    const message = parseIncomingMessage(data);

    // Check if this is a confirmation of a sent message (same content and sender)
    const existingMessage = chatStore.messages.find(
        (msg) =>
            msg.content === message.content &&
            msg.sender_id === message.sender_id &&
            msg.id &&
            msg.id.startsWith("temp-")
    );

    if (existingMessage && message.id && !message.id.startsWith("temp-")) {
        // Update the temp ID with the real ID from backend
        chatStore.updateMessageId(existingMessage.id, message.id);
        console.log(
            `Updated temp ID ${existingMessage.id} to real ID ${message.id}`
        );
    } else {
        // This is a new message from someone else
        chatStore.addMessage(message);
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

    const targetUserId = chatStore.currentChatUser?.id;
    console.log("Sending message to user:", targetUserId);
    
    const chatData = getChatData(targetUserId);

    if (!chatData) {
        console.error("No chat data found for sending message");
        return;
    }

    console.log("Sending message with chat data:", chatData);

    // Check WebSocket connection status
    const connectionStatus = getConnectionStatus();
    console.log("Connection status before sending:", connectionStatus);
    
    if (!connectionStatus.isConnected) {
        console.log("WebSocket not connected, attempting to reconnect...");
        establishConnection(chatData, handleIncomingMessage);
        
        // Wait for connection
        await new Promise(resolve => setTimeout(resolve, 500));
        
        const newStatus = getConnectionStatus();
        if (!newStatus.isConnected) {
            console.error("Failed to establish WebSocket connection");
            return;
        }
    }

    // Create temporary ID for immediate display
    const tempId = `temp-${Date.now()}-${Math.random()
        .toString(36)
        .substr(2, 9)}`;

    const messageData = {
        id: tempId,
        chat_id: chatData.chatId,
        sender_id: chatData.senderId,
        receiver_id: chatData.receiverId,
        content: newMessage.value,
        created_at: new Date().toISOString(),
    };

    // Add message to store immediately with temp ID
    chatStore.addMessage(messageData);

    // Send via WebSocket
    const success = sendWebSocketMessage(messageData.content);
    console.log("WebSocket send result:", success);
    
    newMessage.value = "";
};

// Get current chat ID
const getCurrentChatId = () => {
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
        await loadNextPage(chatId);
    }
};
</script>

<style scoped>
.font-roboto {
    font-family: "Roboto", Arial, sans-serif;
}
</style>
