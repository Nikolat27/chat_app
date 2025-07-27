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
                class="p-3 bg-white rounded-lg shadow-sm border border-gray-200 cursor-pointer hover:bg-gray-50 transition duration-150"
                @click="$emit('chat-clicked', chat)"
            >
                <div class="flex items-center gap-3">
                    <img
                        v-if="getOtherUserAvatar(chat)"
                        :src="`${backendBaseUrl}/static/${getOtherUserAvatar(
                            chat
                        )}`"
                        alt="Avatar"
                        class="w-10 h-10 rounded-full object-cover border border-gray-300 shadow-sm"
                    />
                    <img
                        v-else
                        src="/src/assets/default-avatar.jpeg"
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
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
const props = defineProps({
    chats: Array,
    avatarUrls: Object,
    usernames: Object,
    backendBaseUrl: String,
    currentUserId: [String, Number],
});

const emit = defineEmits(["chat-clicked"]);

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
</script>
