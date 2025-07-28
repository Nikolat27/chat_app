<template>
    <div>
        <h2 class="text-lg font-bold text-green-600 mb-4">Chats</h2>

        <!-- Action Buttons -->
        <div class="flex gap-3 mb-6">
            <button
                class="bg-gradient-to-r from-green-500 to-emerald-500 hover:from-green-600 hover:to-emerald-600 text-white font-semibold px-4 py-2 rounded-lg shadow-md hover:shadow-lg transition-all duration-200 cursor-pointer flex items-center gap-2"
                @click="showCreateChat = true"
            >
                <span class="material-icons text-sm">chat</span>
                Create Chat
            </button>
            <button
                class="bg-gradient-to-r from-purple-500 to-pink-500 hover:from-purple-600 hover:to-pink-600 text-white font-semibold px-4 py-2 rounded-lg shadow-md hover:shadow-lg transition-all duration-200 cursor-pointer flex items-center gap-2"
                @click="showCreateSecretChat = true"
            >
                <span class="material-icons text-sm">lock</span>
                Create Secret Chat
            </button>
        </div>

        <!-- Chat Creation Modal -->
        <CreateChatModal
            v-if="showCreateChat"
            :backend-base-url="backendBaseUrl"
            @close="showCreateChat = false"
            @user-selected="handleUserSelected"
        />

        <!-- Secret Chat Creation Modal -->
        <CreateSecretChatModal
            v-if="showCreateSecretChat"
            :backend-base-url="backendBaseUrl"
            @close="showCreateSecretChat = false"
            @user-selected="handleSecretUserSelected"
        />

        <!-- Chat List -->
        <ChatList
            :chats="chatStore.chats"
            :avatar-urls="chatStore.avatarUrls"
            :usernames="chatStore.usernames"
            :backend-base-url="backendBaseUrl"
            :current-user-id="userStore.user_id"
            @chat-clicked="handleChatClick"
            @chat-deleted="handleChatDeleted"
        />

        <!-- Secret Chat List -->
        <div v-if="props.secretChats && props.secretChats.length > 0" class="mt-8">
            <SecretChatList
                :secretChats="props.secretChats"
                :secretUsernames="props.secretUsernames"
                :backend-base-url="props.backendBaseUrl"
                :current-user-id="props.userStore.user_id"
                @chat-clicked="handleChatClick"
                @secret-chat-deleted="handleSecretChatDeleted"
            />
        </div>
    </div>
</template>

<script setup>
import { ref } from "vue";
import CreateChatModal from "./CreateChatModal.vue";
import CreateSecretChatModal from "./CreateSecretChatModal.vue";
import ChatList from "./ChatList.vue";
import SecretChatList from "./SecretChatList.vue";
import axiosInstance from "../../axiosInstance";
import { showError, showMessage } from "@/utils/toast";

const props = defineProps({
    chatStore: Object,
    userStore: Object,
    backendBaseUrl: String,
    secretChats: Array,
    secretUsernames: Object,
});

const emit = defineEmits(["open-chat"]);

const showCreateChat = ref(false);
const showCreateSecretChat = ref(false);

// Handle user selection from search results
const handleUserSelected = (user) => {
    emit("open-chat", user);
    showCreateChat.value = false;
};

// Handle user selection for secret chat
const handleSecretUserSelected = async (user) => {
    // Call backend to create secret chat
    try {
        const response = await axiosInstance.post("/api/secret-chat/create", {
            target_user: user.id,
        });
        showMessage("Secret chat created successfully! Waiting for approval.");
        if (response.data?.chat) {
            // Optionally update chatStore or emit event
            emit("open-chat", user);
        }
    } catch (error) {
        showError(error.response?.data?.detail || "Failed to create secret chat. Please try again.");
    }
    showCreateSecretChat.value = false;
};

// Handle clicking on existing chat
const handleChatClick = (chat) => {
    const currentUserId = props.userStore.user_id;
    
    // Check if this is a secret chat
    const isSecretChat = props.secretChats?.some(secretChat => 
        secretChat.id === chat.id
    );
    
    if (isSecretChat) {
        // For secret chats, create a user object with the other user's ID
        const otherUserId = chat.user_1 === currentUserId ? chat.user_2 : chat.user_1;
        const user = {
            id: otherUserId,
            username: props.secretUsernames[chat.id] || "Unknown User",
            avatar_url: null, // Secret chats use default avatar
            secret_chat_id: chat.id,
        };
        emit("open-chat", user);
        return;
    }
    
    // Handle regular chats
    if (!chat.participants || chat.participants.length < 2 || !currentUserId) {
        return;
    }

    const otherUserId = chat.participants.find((id) => id !== currentUserId);
    if (!otherUserId) return;

    // Create user object from store data using chat.id (backend structure)
    const user = {
        id: otherUserId,
        username: props.chatStore.usernames[chat.id],
        avatar_url: props.chatStore.avatarUrls[chat.id],
    };

    emit("open-chat", user);
};

// Handle chat deletion
const handleChatDeleted = (chatId) => {
    // Remove the chat from the store
    const chatIndex = props.chatStore.chats.findIndex(chat => chat.id === chatId);
    if (chatIndex !== -1) {
        props.chatStore.chats.splice(chatIndex, 1);
    }
    
    // Remove associated data
    delete props.chatStore.avatarUrls[chatId];
    delete props.chatStore.usernames[chatId];
};

// Handle secret chat deletion
const handleSecretChatDeleted = (chatId) => {
    // Remove the secret chat from the store
    const chatIndex = props.secretChats.findIndex(chat => chat.id === chatId);
    if (chatIndex !== -1) {
        props.secretChats.splice(chatIndex, 1);
    }
    
    // Remove associated data
    delete props.secretUsernames[chatId];
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style>
