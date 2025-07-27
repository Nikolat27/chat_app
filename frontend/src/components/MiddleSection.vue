<template>
    <section class="w-1/3 bg-white border-r p-4 overflow-y-auto">
        <!-- Chats Tab -->
        <div v-if="activeTab === 'chats'">
            <ChatsTab
                :chat-store="chatStore"
                :user-store="userStore"
                :backend-base-url="backendBaseUrl"
                :secret-chats="chatStore.secretChats"
                :secret-usernames="chatStore.secretUsernames"
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
import { defineProps, onMounted, ref } from "vue";
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

const { loadInitialMessages } = useMessagePagination();
const chatsLoaded = ref(false);

onMounted(async () => {
    try {
        await fetchUserChats();
        await fetchUserSecretChats();
    } catch (error) {
        console.error("Failed to fetch chats:", error);
    } finally {
        chatsLoaded.value = true;
    }
});

async function fetchUserChats() {
    try {
        const response = await axiosInstance.get("/api/user/get-chats");
        chatStore.setChats(response.data.chats);
        chatStore.setAvatarUrls(response.data.avatar_urls);
        chatStore.setUsernames(response.data.usernames);
    } catch (error) {
        console.error("Failed to fetch user chats:", error);
    }
}

async function fetchUserSecretChats() {
    try {
        const response = await axiosInstance.get("/api/user/get-secret-chats");
        chatStore.setSecretChats(response.data.secret_chats);
        chatStore.setSecretUsernames(response.data.secret_usernames);
    } catch (error) {
        console.error("Failed to fetch user secret chats:", error);
    }
}

const handleOpenChat = async (user) => {
    chatStore.setChatUser(user);

    if (user.id && user.username) {
        chatStore.updateUserData(user.id, user.username, user.avatar_url);
    }

    const existingChat = findExistingChat(user.id);

    if (existingChat) {
        establishWebSocketConnection(existingChat, user.id);
        await loadInitialMessages(existingChat.id);
    } else {
        await createNewChat(user);
    }
};

const findExistingChat = (targetUserId) => {
    const currentUserId = userStore.user_id;
    return chatStore.chats?.find(
        (chat) =>
            chat.participants &&
            chat.participants.includes(targetUserId) &&
            chat.participants.includes(currentUserId)
    );
};

const createNewChat = async (user) => {
    try {
        const response = await axiosInstance.post("/api/chat/create", {
            target_user: user.id,
        });

        if (response.data?.chat) {
            const newChat = response.data.chat;
            chatStore.chats.push(newChat);
            chatStore.updateChatData(
                newChat.id,
                user.username,
                user.avatar_url
            );
            establishWebSocketConnection(newChat, user.id);
            await loadInitialMessages(newChat.id);
        }
    } catch (error) {
        console.error("Failed to create chat:", error);
    }
};

const establishWebSocketConnection = (chat, targetUserId) => {
    const currentUserId = userStore.user_id;
    const senderId = currentUserId;
    const receiverId = targetUserId;

    if (!chat.id || !senderId || !receiverId) {
        console.error("Missing required IDs for WebSocket connection");
        return;
    }

    const wsUrl = `${backendBaseUrl.replace(
        /^http/,
        "ws"
    )}/api/websocket/chat/add/${
        chat.id
    }?sender_id=${senderId}&receiver_id=${receiverId}`;
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
