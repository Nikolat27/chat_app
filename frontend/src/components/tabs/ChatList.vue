<template>
    <div>
        <!-- Empty State -->
        <div
            v-if="!chats || chats.length === 0"
            class="text-gray-500"
        >
            No chats yet.
        </div>

        <!-- Chat List -->
        <div v-else>
            <div
                v-for="chat in chats"
                :key="chat.id"
                class="mb-2 p-2 border rounded cursor-pointer hover:bg-gray-100"
                @click="$emit('chat-clicked', chat)"
            >
                <div class="flex items-center gap-2">
                    <img
                        v-if="getOtherUserAvatar(chat)"
                        :src="`${backendBaseUrl}/static/${getOtherUserAvatar(chat)}`"
                        alt="Avatar"
                        class="w-8 h-8 rounded-full object-cover border"
                    />
                    <img
                        v-else
                        src="/src/assets/default-avatar.jpeg"
                        alt="Default Avatar"
                        class="w-8 h-8 rounded-full object-cover border"
                    />
                    <span class="font-bold text-green-900">
                        {{ getOtherUsername(chat) }}
                    </span>
                </div>
                <div class="text-xs text-gray-600">
                    Created: {{ new Date(chat.created_at).toLocaleString() }}
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

const emit = defineEmits(['chat-clicked']);

// Get the other user's ID from chat participants
const getOtherUserId = (chat) => {
    if (!chat.participants || chat.participants.length < 2 || !props.currentUserId) {
        return null;
    }
    return chat.participants.find(id => id !== props.currentUserId);
};

// Get the other user's username
const getOtherUsername = (chat) => {
    return props.usernames[chat.id] || 'Unknown User';
};

// Get the other user's avatar
const getOtherUserAvatar = (chat) => {
    return props.avatarUrls[chat.id] || null;
};
</script> 