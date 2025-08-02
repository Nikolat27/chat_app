<template>
    <div
        class="flex items-center gap-4 px-6 py-4 border-gray-200 border-b bg-white shadow-sm"
    >
        <!-- Avatar with Secret Chat Indicator -->
        <div class="relative">
            <img
                v-if="user.avatar_url"
                :src="`${backendBaseUrl}/static/${user.avatar_url}`"
                class="w-12 h-12 rounded-full object-cover border-2 border-gray-200 shadow-sm select-none pointer-events-none"
                alt="Avatar"
            />
            <img
                v-else
                src="/src/assets/default-avatar.jpg"
                class="w-12 h-12 rounded-full object-cover border-2 border-gray-200 shadow-sm select-none pointer-events-none"
                alt="Default Avatar"
            />
            
            <!-- Secret Chat Indicator -->
            <div v-if="isSecretChat" class="absolute -bottom-1 -right-1 w-6 h-6 bg-purple-500 rounded-full flex items-center justify-center shadow-md">
                <span class="material-icons text-white text-xs">lock</span>
            </div>
        </div>

        <!-- User Info -->
        <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
                <span class="font-semibold text-gray-800 text-lg leading-tight truncate">
                    {{ user.username }}
                </span>
                <div v-if="isSecretChat" class="flex items-center gap-1">
                    <div class="px-2 py-1 bg-purple-100 text-purple-700 rounded-full text-xs font-medium border border-purple-200 flex items-center gap-1">
                        <span class="material-icons text-xs">security</span>
                        Secret
                    </div>
                </div>
            </div>
            
            <div class="flex items-center gap-3 text-xs">
                <div class="flex items-center gap-1 text-green-600">
                    <span class="material-icons text-xs">circle</span>
                    <span class="font-medium">Online</span>
                </div>
                
                <div v-if="isSecretChat" class="flex items-center gap-1 text-purple-600">
                    <span class="material-icons text-xs">verified</span>
                    <span>End-to-end encrypted</span>
                </div>
                
                <div v-else class="flex items-center gap-1 text-gray-500">
                    <span class="material-icons text-xs">chat</span>
                    <span>Regular chat</span>
                </div>
            </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center gap-2">
            <button 
                v-if="isSecretChat"
                class="w-10 h-10 text-purple-600 hover:bg-purple-50 rounded-full transition-colors duration-200 cursor-pointer flex items-center justify-center"
                title="Secret chat info"
                @click="showSecretChatInfo = true"
            >
                <span class="material-icons text-lg">info</span>
            </button>
            
            <!-- Delete Chat Button -->
            <button 
                @click="showDeleteModal = true"
                class="w-10 h-10 text-red-500 hover:bg-red-50 rounded-full transition-colors duration-200 cursor-pointer flex items-center justify-center"
                title="Delete chat"
            >
                <span class="material-icons text-lg">delete</span>
            </button>
        </div>
    </div>

    <!-- Secret Chat Info Modal -->
    <SecretChatInfoModal
        v-if="isSecretChat"
        :is-visible="showSecretChatInfo"
        :user="user"
        :backend-base-url="backendBaseUrl"
        @close="showSecretChatInfo = false"
    />

    <!-- Delete Chat Confirmation Modal -->
    <ConfirmModal
        :is-visible="showDeleteModal"
        title="Delete Chat"
        subtitle="This action cannot be undone"
        :message="`Are you sure you want to delete the chat with ${user.username}? This action cannot be undone.`"
        confirm-text="Delete Chat"
        :is-loading="isDeleting"
        @close="showDeleteModal = false"
        @confirm="confirmDeleteChat"
    />
</template>

<script setup>
import { computed, ref } from "vue";
import SecretChatInfoModal from "./SecretChatInfoModal.vue";
import ConfirmModal from "../ui/ConfirmModal.vue";

const props = defineProps({
    user: {
        type: Object,
        required: true,
    },
    backendBaseUrl: {
        type: String,
        required: true,
    },
    isSecretChat: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits(['delete-chat']);

const showSecretChatInfo = ref(false);
const showDeleteModal = ref(false);
const isDeleting = ref(false);

const handleDeleteChat = () => {
    emit('delete-chat', props.user);
};

const confirmDeleteChat = async () => {
    isDeleting.value = true;
    try {
        await handleDeleteChat();
        showDeleteModal.value = false;
    } finally {
        isDeleting.value = false;
    }
};

// You can add more computed properties here if needed
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style>
