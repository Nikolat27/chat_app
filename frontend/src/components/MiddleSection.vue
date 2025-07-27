<template>
    <section class="w-1/3 bg-white border-r p-4 overflow-y-auto">
        <!-- Chats Tab -->
        <div v-if="activeTab === 'chats'">
            <ChatsTab 
                :chat-store="chatStore"
                :user-store="userStore"
                :backend-base-url="backendBaseUrl"
                @open-chat="handleOpenChat"
            />
        </div>

        <!-- Groups Tab -->
        <div v-else-if="activeTab === 'groups'">
            <GroupsTab />
        </div>

        <!-- Settings Tab -->
        <div v-else-if="activeTab === 'settings'">
            <SettingsTab />
        </div>
    </section>
</template>

<script setup>
import { defineProps, onMounted } from "vue";
import { useChatStore } from "../stores/chat";
import { useUserStore } from "../stores/users";
import axiosInstance from "@/axiosInstance";
import { useMessagePagination } from "../composables/useMessagePagination";
import ChatsTab from "./tabs/ChatsTab.vue";
import GroupsTab from "./tabs/GroupsTab.vue";
import SettingsTab from "./tabs/SettingsTab.vue";

const props = defineProps({ activeTab: String });
const chatStore = useChatStore();
const userStore = useUserStore();
const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;

// Message pagination
const { loadInitialMessages } = useMessagePagination();

// Load initial chat data
onMounted(async () => {
    try {
        const response = await axiosInstance.get("/api/user/get-chats");
        chatStore.setChats(response.data.chats);
        chatStore.setAvatarUrls(response.data.avatar_urls);
        chatStore.setUsernames(response.data.usernames);
    } catch (error) {
        console.error("Failed to fetch chats:", error);
    }
});

// Handle opening a chat (from existing chat or search result)
const handleOpenChat = async (user) => {
    chatStore.setChatUser(user);
    
    // Update store with user data if available
    if (user.id && user.username) {
        chatStore.updateUserData(user.id, user.username, user.avatar_url);
    }

    // Find existing chat or create new one
    const existingChat = findExistingChat(user.id);
    
    if (existingChat) {
        establishWebSocketConnection(existingChat, user.id);
        // Load initial messages for existing chat
        await loadInitialMessages(existingChat.id);
    } else {
        await createNewChat(user);
    }
};

// Find existing chat between current user and target user
const findExistingChat = (targetUserId) => {
    const currentUserId = userStore.user_id;
    return chatStore.chats.find(chat => 
        chat.participants && 
        chat.participants.includes(targetUserId) && 
        chat.participants.includes(currentUserId)
    );
};

// Create new chat with target user
const createNewChat = async (user) => {
    try {
        const response = await axiosInstance.post("/api/chat/create", {
            participant_id: user.id,
        });
        
        if (response.data?.chat) {
            const newChat = response.data.chat;
            chatStore.chats.push(newChat);
            chatStore.updateChatData(newChat.id, user.username, user.avatar_url);
            establishWebSocketConnection(newChat, user.id);
            // Load initial messages for new chat (will be empty initially)
            await loadInitialMessages(newChat.id);
        }
    } catch (error) {
        console.error("Failed to create chat:", error);
    }
};

// Establish WebSocket connection for chat
const establishWebSocketConnection = (chat, targetUserId) => {
    const currentUserId = userStore.user_id;
    const senderId = currentUserId;
    const receiverId = targetUserId;
    
    if (!chat.id || !senderId || !receiverId) {
        console.error("Missing required IDs for WebSocket connection");
        return;
    }

    const wsUrl = `${backendBaseUrl.replace(/^http/, "ws")}/api/websocket/chat/add/${chat.id}?sender_id=${senderId}&receiver_id=${receiverId}`;
    const chatSocket = new WebSocket(wsUrl);

    chatSocket.onopen = () => {
        console.log("WebSocket connected for chat:", chat.id);
    };
    
    chatSocket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        chatStore.addMessage(data);
    };
    
    chatSocket.onclose = () => {
        console.log("WebSocket closed for chat:", chat.id);
    };
    
    chatSocket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
};
</script>
