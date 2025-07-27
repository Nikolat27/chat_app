<template>
    <section class="flex flex-col h-full w-full bg-gray-50 font-roboto">
        <!-- Chat Header -->
        <ChatHeader
            v-if="chatStore.currentChatUser"
            :user="chatStore.currentChatUser"
            :backend-base-url="backendBaseUrl"
        />

        <!-- No Chat Selected State -->
        <NoChatSelected v-else />

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
    </section>
</template>

<script setup>
import { ref, watch } from "vue";
import { useChatStore } from "../stores/chat";
import { useUserStore } from "../stores/users";
import ChatHeader from "./chat/ChatHeader.vue";
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

// WebSocket management
const { establishConnection, sendMessage: sendWebSocketMessage } =
    useWebSocket();

// Message pagination
const { loadNextPage, loadInitialMessages } = useMessagePagination();

// Message deletion
const { updateMessageId } = useMessageDeletion();

// Watch for chat user changes to manage WebSocket connections
watch(
    () => chatStore.currentChatUser?.id,
    (newUserId, oldUserId) => {
        if (oldUserId) {
            // Close previous connection
            // This is handled by the composable
        }

        if (newUserId) {
            const chatData = getChatData(newUserId);
            if (chatData) {
                establishConnection(chatData, handleIncomingMessage);
            }
        }
    }
);

// Get chat data for WebSocket connection
const getChatData = (targetUserId) => {
    const chat = chatStore.chats?.find(
        (c) =>
            c.participants &&
            c.participants.includes(targetUserId) &&
            c.participants.includes(userStore.user_id)
    );

    if (!chat || !chat.participants || chat.participants.length < 2) {
        return null;
    }

    const senderId = userStore.user_id;
    const receiverId = targetUserId;

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
const sendMessage = () => {
    if (!newMessage.value.trim()) return;

    const targetUserId = chatStore.currentChatUser?.id;
    const chatData = getChatData(targetUserId);

    if (!chatData) return;

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
    sendWebSocketMessage(messageData.content);
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
