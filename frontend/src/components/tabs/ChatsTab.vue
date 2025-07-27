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
                class="bg-gray-300 text-gray-700 px-3 py-1 rounded cursor-not-allowed"
                disabled
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
import ChatList from "./ChatList.vue";

const props = defineProps({
    chatStore: Object,
    userStore: Object,
    backendBaseUrl: String,
});

const emit = defineEmits(['open-chat']);

const showCreateChat = ref(false);

// Handle user selection from search results
const handleUserSelected = (user) => {
    emit('open-chat', user);
    showCreateChat.value = false;
};

// Handle clicking on existing chat
const handleChatClick = (chat) => {
    const currentUserId = props.userStore.user_id;
    if (!chat.participants || chat.participants.length < 2 || !currentUserId) {
        return;
    }

    const otherUserId = chat.participants.find(id => id !== currentUserId);
    if (!otherUserId) return;

    // Create user object from store data using chat.id (backend structure)
    const user = {
        id: otherUserId,
        username: props.chatStore.usernames[chat.id],
        avatar_url: props.chatStore.avatarUrls[chat.id]
    };

    emit('open-chat', user);
};
</script> 