<template>
    <div>
        <!-- Header -->
        <div class="mb-6">
            <div class="flex items-center gap-3 mb-2">
                <span class="material-icons text-purple-600 text-xl">lock</span>
                <h3 class="text-lg font-bold text-gray-800">Secret Chats</h3>
                <div class="flex-1"></div>
                <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-full">
                    {{ secretChats ? secretChats.length : 0 }} chats
                </span>
            </div>
            <p class="text-sm text-gray-600">
                End-to-end encrypted conversations for maximum privacy
            </p>
        </div>

        <!-- Confirmation Modal -->
        <ConfirmModal
            :is-visible="showDeleteModal"
            title="Delete Secret Chat"
            subtitle="This action cannot be undone"
            :message="deleteModalMessage"
            confirm-text="Delete Secret Chat"
            :is-loading="isDeleting !== null"
            @close="showDeleteModal = false"
            @confirm="confirmDeleteSecretChat"
        />

        <!-- Empty State -->
        <div
            v-if="!secretChats || secretChats.length === 0"
            class="text-center py-12 px-6"
        >
            <div class="mb-4">
                <span class="material-icons text-6xl text-gray-300">lock</span>
            </div>
            <h4 class="text-lg font-semibold text-gray-600 mb-2">No Secret Chats Yet</h4>
            <p class="text-sm text-gray-500 mb-6 leading-relaxed">
                Start your first encrypted conversation by creating a secret chat with another user.
            </p>
            <div class="bg-purple-50 rounded-xl p-4 border border-purple-200">
                <div class="flex items-center gap-2 mb-2">
                    <span class="material-icons text-purple-600 text-sm">security</span>
                    <span class="text-sm font-semibold text-purple-700">Privacy Features</span>
                </div>
                <ul class="text-xs text-purple-600 space-y-1">
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">check_circle</span>
                        End-to-end encryption
                    </li>
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">check_circle</span>
                        Messages not stored on server
                    </li>
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">check_circle</span>
                        Unique key pairs per chat
                    </li>
                </ul>
            </div>
        </div>

        <!-- Chat List -->
        <div v-else class="space-y-4">
            <div
                v-for="chat in secretChats"
                :key="chat.id"
                class="group bg-white rounded-2xl shadow-sm border border-gray-200 cursor-pointer hover:shadow-lg hover:border-purple-300 transition-all duration-300 overflow-hidden"
                @click="handleSecretChatClick(chat)"
            >
                <div class="p-4">
                    <div class="flex items-center gap-4">
                        <!-- Avatar with Status -->
                        <div class="relative">
                                                    <img
                            src="/src/assets/default-secret-chat-avatar.jpg"
                            alt="Secret Chat Avatar"
                            class="w-12 h-12 rounded-full object-cover border-2 border-purple-300 shadow-sm group-hover:border-purple-400 transition-colors duration-200 select-none pointer-events-none"
                        />
                            <div class="absolute -bottom-1 -right-1 w-5 h-5 bg-purple-500 rounded-full flex items-center justify-center">
                                <span class="material-icons text-white text-xs">lock</span>
                            </div>
                        </div>

                        <!-- Chat Info -->
                        <div class="flex-1 min-w-0">
                            <div class="flex items-center gap-2 mb-1">
                                <span class="font-semibold text-gray-800 truncate">
                                    {{ getOtherUsername(chat) }}
                                </span>
                                <div 
                                    :class="chat.user_2_accepted 
                                        ? 'bg-green-100 text-green-700 border-green-200' 
                                        : 'bg-orange-100 text-orange-700 border-orange-200'"
                                    class="px-2 py-1 rounded-full text-xs font-medium border flex items-center gap-1"
                                >
                                    <span class="material-icons text-xs">
                                        {{ chat.user_2_accepted ? 'check_circle' : 'pending' }}
                                    </span>
                                    {{ chat.user_2_accepted ? 'Active' : 'Pending' }}
                                </div>
                            </div>
                            
                            <div class="flex items-center gap-4 text-xs text-gray-500">
                                <div class="flex items-center gap-1">
                                    <span class="material-icons text-xs">schedule</span>
                                    {{ formatDate(chat.created_at) }}
                                </div>
                                <div class="flex items-center gap-1">
                                    <span class="material-icons text-xs">security</span>
                                    Encrypted
                                </div>
                            </div>
                        </div>

                        <!-- Action Buttons -->
                        <div class="flex items-center gap-2">
                            <button
                                v-if="shouldShowApprove(chat)"
                                @click.stop="approveSecretChat(chat)"
                                :disabled="isApproving === chat.id"
                                class="bg-gradient-to-r from-purple-500 to-pink-500 text-white text-sm px-4 py-2 rounded-lg font-medium hover:from-purple-600 hover:to-pink-600 transition-all duration-200 flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                            >
                                <span v-if="isApproving !== chat.id" class="material-icons text-sm">check</span>
                                <div v-else class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                                {{ isApproving === chat.id ? 'Approving...' : 'Approve' }}
                            </button>
                            <span v-else class="material-icons text-gray-400 group-hover:text-purple-500 transition-colors duration-200">chevron_right</span>
                            
                            <!-- Delete Button -->
                            <div class="flex items-center justify-center">
                                <button
                                    @click.stop="handleDeleteSecretChat(chat)"
                                    :disabled="isDeleting === chat.id"
                                    class="opacity-0 group-hover:opacity-100 transition-all duration-200 p-3 text-red-500 hover:bg-red-50 rounded-full hover:text-red-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md"
                                    title="Delete secret chat"
                                >
                                    <span v-if="isDeleting !== chat.id" class="material-icons text-base">delete</span>
                                    <div v-else class="animate-spin rounded-full h-5 w-5 border-b-2 border-red-500"></div>
                                </button>
                            </div>
                        </div>
                    </div>

                    <!-- Status Message -->
                    <div v-if="!chat.user_2_accepted" class="mt-3 p-3 bg-orange-50 border border-orange-200 rounded-lg">
                        <div class="flex items-center gap-2 text-xs text-orange-700">
                            <span class="material-icons text-xs">info</span>
                            <span class="font-medium">Waiting for approval</span>
                        </div>
                        <p class="text-xs text-orange-600 mt-1 leading-relaxed">
                            The other user needs to approve this secret chat before you can start messaging.
                        </p>
                    </div>

                    <div v-else class="mt-3 p-3 bg-green-50 border border-green-200 rounded-lg">
                        <div class="flex items-center gap-2 text-xs text-green-700">
                            <span class="material-icons text-xs">verified</span>
                            <span class="font-medium">Chat is active</span>
                        </div>
                        <p class="text-xs text-green-600 mt-1 leading-relaxed">
                            Your messages are end-to-end encrypted and secure.
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from "vue";
import axiosInstance from "@/axiosInstance";
import { showError, showMessage } from "@/utils/toast";
import { useKeyPair } from "@/composables/useKeyPair";
import ConfirmModal from "../ui/ConfirmModal.vue";

const props = defineProps({
    secretChats: Array,
    secretUsernames: Object,
    backendBaseUrl: String,
    currentUserId: [String, Number],
});

const emit = defineEmits(["chat-clicked", "secret-chat-deleted"]);

const isApproving = ref(null);
const isDeleting = ref(null);
const showDeleteModal = ref(false);
const deleteModalMessage = ref("");
const chatToDelete = ref(null);
const { generateSecretChatKeyPair, hasSecretChatKeys } = useKeyPair();

const shouldShowApprove = (chat) => {
    return (
        chat.user_2 === props.currentUserId && chat.user_2_accepted === false
    );
};

async function approveSecretChat(chat) {
    isApproving.value = chat.id;
    
    try {
        await axiosInstance.post(`/api/secret-chat/approve/${chat.id}`);
        showMessage("Secret chat approved successfully! You can now start messaging.");
        chat.user_2_accepted = true;
    } catch (err) {
        showError("Failed to approve secret chat. Please try again.");
    } finally {
        isApproving.value = null;
    }
}

// Handle secret chat click with key generation
const handleSecretChatClick = async (chat) => {
    try {
        // Check if we already have keys for this chat
        const hasKeys = await hasSecretChatKeys(chat.id);
        
        if (!hasKeys) {
            // Generate new key pair for this secret chat
            const publicKey = await generateSecretChatKeyPair(chat.id);
            
            // Send public key to backend
            await axiosInstance.post(`/api/secret-chat/add-public-key/${chat.id}`, {
                public_key: publicKey
            });
            
            showMessage("Encryption keys generated and exchanged successfully!");
        }
        
        // Emit the chat click event
        emit("chat-clicked", chat);
        
    } catch (error) {
        console.error("Error handling secret chat click:", error);
        showError("Failed to set up encryption for secret chat. Please try again.");
    }
};

// Get the other user's username
const getOtherUsername = (chat) => {
    return props.secretUsernames[chat.id] || "Unknown User";
};

// Format date for display
const formatDate = (dateString) => {
    if (!dateString) return "Unknown";
    
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 1) {
        return "Today";
    } else if (diffDays === 2) {
        return "Yesterday";
    } else if (diffDays <= 7) {
        return `${diffDays - 1} days ago`;
    } else {
        return date.toLocaleDateString();
    }
};

// Handle secret chat deletion
const handleDeleteSecretChat = (chat) => {
    chatToDelete.value = chat;
    deleteModalMessage.value = `Are you sure you want to delete the secret chat with ${getOtherUsername(chat)}? This action cannot be undone.`;
    showDeleteModal.value = true;
};

// Confirm secret chat deletion
const confirmDeleteSecretChat = async () => {
    if (!chatToDelete.value) return;
    
    isDeleting.value = chatToDelete.value.id;
    
    try {
        await axiosInstance.delete(`/api/secret-chat/delete/${chatToDelete.value.id}`);
        showMessage("Secret chat deleted successfully!");
        emit("secret-chat-deleted", chatToDelete.value.id);
        showDeleteModal.value = false;
    } catch (error) {
        console.error("Failed to delete secret chat:", error);
        showError(error.response?.data?.detail || "Failed to delete secret chat. Please try again.");
    } finally {
        isDeleting.value = null;
        chatToDelete.value = null;
    }
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style>
