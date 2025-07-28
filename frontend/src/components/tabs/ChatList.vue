<template>
    <div>
        <!-- Empty State -->
        <div
            v-if="!chats || chats.length === 0"
            class="text-gray-500 text-sm text-center mt-6 italic"
        >
            No chats yet.
        </div>

        <!-- Chat List -->
        <div v-else class="space-y-3 mt-2">
            <div
                v-for="chat in chats"
                :key="chat.id"
                class="p-4 bg-white rounded-xl shadow-sm border border-gray-200 hover:bg-gray-50 hover:shadow-md transition-all duration-200 group"
            >
                <div class="flex items-center justify-between">
                    <div
                        class="flex items-center gap-3 flex-1 cursor-pointer"
                        @click="$emit('chat-clicked', chat)"
                    >
                        <img
                            v-if="getOtherUserAvatar(chat)"
                            :src="`${backendBaseUrl}/static/${getOtherUserAvatar(
                                chat
                            )}`"
                            alt="Avatar"
                            class="w-10 h-10 rounded-full object-cover border border-gray-300 shadow-sm select-none pointer-events-none"
                        />
                        <img
                            v-else
                            src="/src/assets/default-avatar.jpg"
                            alt="Default Avatar"
                            class="w-10 h-10 rounded-full object-cover border border-gray-300 shadow-sm select-none pointer-events-none"
                        />
                        <div class="flex flex-col">
                            <span class="font-semibold text-gray-800">
                                {{ getOtherUsername(chat) }}
                            </span>
                            <span class="text-xs text-gray-500">
                                Created:
                                {{ new Date(chat.created_at).toLocaleString() }}
                            </span>
                        </div>
                    </div>

                    <!-- Delete Button -->
                    <div class="flex items-center justify-center">
                        <button
                            @click.stop="handleDeleteChat(chat)"
                            :disabled="isDeleting === chat.id"
                            class="w-[48px] h-[48px] opacity-0 group-hover:opacity-100 transition-all duration-200 p-3 text-red-500 hover:bg-red-50 rounded-full hover:text-red-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md"
                            title="Delete chat"
                        >
                            <span
                                v-if="isDeleting !== chat.id"
                                class="material-icons text-base"
                                >delete</span
                            >
                            <div
                                v-else
                                class="animate-spin rounded-full h-5 w-5 border-b-2 border-red-500"
                            ></div>
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Confirmation Modal -->
        <ConfirmModal
            :is-visible="showDeleteModal"
            title="Delete Chat"
            subtitle="This action cannot be undone"
            :message="deleteModalMessage"
            confirm-text="Delete Chat"
            :is-loading="isDeleting !== null"
            @close="showDeleteModal = false"
            @confirm="confirmDeleteChat"
        />
    </div>
</template>

<script setup>
import { ref } from "vue";
import axiosInstance from "../../axiosInstance";
import { showError, showMessage } from "../../utils/toast";
import ConfirmModal from "../ui/ConfirmModal.vue";

const props = defineProps({
    chats: Array,
    avatarUrls: Object,
    usernames: Object,
    backendBaseUrl: String,
    currentUserId: [String, Number],
});

const emit = defineEmits(["chat-clicked", "chat-deleted"]);

const isDeleting = ref(null);
const showDeleteModal = ref(false);
const deleteModalMessage = ref("");
const chatToDelete = ref(null);

// Get the other user's ID from chat participants
const getOtherUserId = (chat) => {
    if (
        !chat.participants ||
        chat.participants.length < 2 ||
        !props.currentUserId
    ) {
        return null;
    }
    return chat.participants.find((id) => id !== props.currentUserId);
};

// Get the other user's username
const getOtherUsername = (chat) => {
    return props.usernames[chat.id] || "Unknown User";
};

// Get the other user's avatar
const getOtherUserAvatar = (chat) => {
    return props.avatarUrls[chat.id] || null;
};

// Handle chat deletion
const handleDeleteChat = (chat) => {
    chatToDelete.value = chat;
    deleteModalMessage.value = `Are you sure you want to delete the chat with ${getOtherUsername(
        chat
    )}? This action cannot be undone.`;
    showDeleteModal.value = true;
};

// Confirm chat deletion
const confirmDeleteChat = async () => {
    if (!chatToDelete.value) return;

    isDeleting.value = chatToDelete.value.id;

    try {
        await axiosInstance.delete(`/api/chat/delete/${chatToDelete.value.id}`);
        showMessage("Chat deleted successfully!");
        emit("chat-deleted", chatToDelete.value.id);
        showDeleteModal.value = false;
    } catch (error) {
        console.error("Failed to delete chat:", error);
        showError(
            error.response?.data?.detail ||
                "Failed to delete chat. Please try again."
        );
    } finally {
        isDeleting.value = null;
        chatToDelete.value = null;
    }
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style>
