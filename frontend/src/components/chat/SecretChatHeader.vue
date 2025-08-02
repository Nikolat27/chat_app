<template>
    <div
        class="flex items-center gap-4 px-6 py-4 border-gray-200 border-b bg-white shadow-sm"
    >
        <!-- Avatar with Secret Chat Indicator -->
        <div class="relative">
            <img
                src="/src/assets/default-secret-chat-avatar.jpg"
                alt="Secret Chat Avatar"
                class="w-12 h-12 rounded-full object-cover border-2 border-purple-300 shadow-sm select-none pointer-events-none"
            />
            <div class="absolute -bottom-1 -right-1 w-6 h-6 bg-purple-500 rounded-full flex items-center justify-center shadow-md">
                <span class="material-icons text-white text-xs">lock</span>
            </div>
        </div>

        <!-- User Info -->
        <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
                <span class="font-semibold text-gray-800 text-lg leading-tight truncate">
                    {{ getOtherUsername() }}
                </span>
                <div class="flex items-center gap-1">
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
                
                <div class="flex items-center gap-1 text-purple-600">
                    <span class="material-icons text-xs">verified</span>
                    <span>End-to-end encrypted</span>
                </div>
            </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center gap-2">
            <button 
                class="w-10 h-10 text-purple-600 hover:bg-purple-50 rounded-full transition-colors duration-200 cursor-pointer flex items-center justify-center"
                title="Secret chat info"
                @click="showSecretChatInfo = true"
            >
                <span class="material-icons text-lg">info</span>
            </button>
            
            <button 
                class="w-10 h-10 text-gray-500 hover:bg-gray-100 rounded-full transition-colors duration-200 cursor-pointer flex items-center justify-center"
                title="More options"
            >
                <span class="material-icons text-lg">more_vert</span>
            </button>
        </div>
    </div>

    <!-- Secret Chat Info Modal -->
    <SecretChatInfoModal
        :is-visible="showSecretChatInfo"
        :user="secretChatUser"
        :backend-base-url="backendBaseUrl"
        @close="showSecretChatInfo = false"
    />
</template>

<script setup>
import { ref, computed } from "vue";
import SecretChatInfoModal from "./SecretChatInfoModal.vue";

const props = defineProps({
    secretChat: {
        type: Object,
        required: true,
    },
    secretUsernames: {
        type: Object,
        default: () => ({}),
    },
    backendBaseUrl: {
        type: String,
        required: true,
    },
    currentUserId: {
        type: [String, Number],
        required: true,
    },
});

const showSecretChatInfo = ref(false);

// Create a user object for the secret chat info modal
const secretChatUser = computed(() => {
    return {
        id: props.secretChat.user_1 === props.currentUserId ? props.secretChat.user_2 : props.secretChat.user_1,
        username: getOtherUsername(),
        avatar_url: null, // Secret chats use default avatar
        secret_chat_id: props.secretChat.id,
    };
});

// Get the other user's username
const getOtherUsername = () => {
    return props.secretUsernames[props.secretChat.id] || "Unknown User";
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style> 