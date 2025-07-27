<template>
    <div>
        <h2 class="text-lg font-bold text-green-600 mb-4">Chats</h2>

        <!-- Action Buttons -->
        <div class="flex gap-2 mb-4">
            <button
                class="bg-green-500 text-white px-3 py-1 rounded hover:bg-green-600"
                @click="showCreateChat = true"
            >
                Create Chat
            </button>
            <button
                class="bg-purple-500 text-white px-3 py-1 rounded hover:bg-purple-600"
                @click="showCreateSecretChat = true"
            >
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
        />
    </div>
</template>

<script setup>
import { ref } from "vue";
import CreateChatModal from "./CreateChatModal.vue";
import CreateSecretChatModal from "./CreateSecretChatModal.vue";
import ChatList from "./ChatList.vue";
import axiosInstance from "../../axiosInstance";
import { showError, showMessage } from "@/utils/toast";

const props = defineProps({
    chatStore: Object,
    userStore: Object,
    backendBaseUrl: String,
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
        showMessage("Secret chat created successfully!");
        if (response.data?.chat) {
            // Optionally update chatStore or emit event
            emit("open-chat", user);
        }
    } catch (error) {
        showError(error.response.data.detail);
    }
    showCreateSecretChat.value = false;
};

// Handle clicking on existing chat
const handleChatClick = (chat) => {
    const currentUserId = props.userStore.user_id;
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
</script>
