<template>
    <div>
        <!-- Empty State -->
        <div
            v-if="!secretChats || secretChats.length === 0"
            class="text-gray-500 text-sm text-center mt-6 italic"
        >
            No secret chats yet.
        </div>

        <!-- Chat List -->
        <div v-else class="space-y-3 mt-2">
            <div
                v-for="chat in secretChats"
                :key="chat.id"
                class="p-3 bg-white rounded-lg shadow-sm border border-gray-200 cursor-pointer hover:bg-gray-50 transition duration-150"
                @click="$emit('chat-clicked', chat)"
            >
                <div class="flex items-center gap-3">
                    <img
                        src="/src/assets/default-secret-chat-avatar.jpg"
                        alt="Default Avatar"
                        class="w-10 h-10 rounded-full object-cover border border-gray-300 shadow-sm"
                    />
                    <div class="flex flex-col">
                        <span class="font-semibold text-gray-800">
                            {{ getOtherUsername(chat) }}
                        </span>
                        <span class="text-xs text-gray-500">
                            Created:
                            {{ new Date(chat.created_at).toLocaleString() }}
                        </span>
                        <span
                            class="text-xs mt-1"
                            :class="
                                chat.user_2_accepted
                                    ? 'text-green-600'
                                    : 'text-red-600'
                            "
                        >
                            {{
                                chat.user_2_accepted
                                    ? "Approved"
                                    : "Pending Approval"
                            }}
                        </span>
                    </div>
                </div>

                <!-- Approve Button -->
                <div v-if="shouldShowApprove(chat)" class="mt-3 text-right">
                    <button
                        @click.stop="approveSecretChat(chat)"
                        class="cursor-pointer bg-blue-500 hover:bg-blue-600 text-white text-sm px-3 py-1 rounded-md transition"
                    >
                        Approve Secret Chat
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import axiosInstance from "@/axiosInstance";
import { showError, showMessage } from "@/utils/toast";

const props = defineProps({
    secretChats: Array,
    secretUsernames: Object,
    backendBaseUrl: String,
    currentUserId: [String, Number],
});

const emit = defineEmits(["chat-clicked"]);

const shouldShowApprove = (chat) => {
    return (
        chat.user_2 === props.currentUserId && chat.user_2_accepted === false
    );
};

async function approveSecretChat(chat) {
    try {
        await axiosInstance.post(`/api/secret-chat/approve/${chat.id}`);
        showMessage("Secret chat approved!");
        chat.user_2_accepted = true;
    } catch (err) {
        showError("Failed to approve chat.");
    }
}

// Get the other user's username
const getOtherUsername = (chat) => {
    return props.secretUsernames[chat.id] || "Unknown User";
};
</script>
