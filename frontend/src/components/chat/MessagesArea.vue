<template>
    <div class="flex-1 overflow-y-auto p-4 space-y-6">
        <!-- Messages List -->
        <div v-if="messages && messages.length">
            <div
                v-for="msg in messages"
                :key="msg.id"
                :class="msg.sender_id === currentUserId ? 'justify-end' : 'justify-start'"
                class="flex items-end mb-2 gap-2"
            >
                <!-- Avatar (left for received, right for sent) -->
                <template v-if="msg.sender_id !== currentUserId">
                    <img
                        v-if="otherUserAvatar"
                        :src="getAvatarUrl(otherUserAvatar)"
                        class="w-8 h-8 rounded-full object-cover border"
                        alt="Avatar"
                    />
                    <img
                        v-else
                        src="/src/assets/default-avatar.jpeg"
                        class="w-8 h-8 rounded-full object-cover border"
                        alt="Default Avatar"
                    />
                </template>

                <!-- Message Bubble -->
                <div
                    :class="msg.sender_id === currentUserId ? 'bg-green-500 text-white' : 'bg-white text-gray-800 border'"
                    class="rounded-lg px-4 py-2 max-w-xs"
                >
                    <div class="text-base font-semibold break-words whitespace-pre-line">
                        {{ msg.content }}
                    </div>
                    <div 
                        :class="msg.sender_id === currentUserId ? 'text-white' : 'text-gray-500'" 
                        class="text-xs text-right mt-1"
                    >
                        {{ formatTime(msg.timestamp) }}
                    </div>
                </div>

                <!-- Avatar for sent messages (right side) -->
                <template v-if="msg.sender_id === currentUserId">
                    <img
                        v-if="userAvatar"
                        :src="getAvatarUrl(userAvatar)"
                        class="w-8 h-8 rounded-full object-cover border"
                        alt="Avatar"
                    />
                    <img
                        v-else
                        src="/src/assets/default-avatar.jpeg"
                        class="w-8 h-8 rounded-full object-cover border"
                        alt="Default Avatar"
                    />
                </template>
            </div>
        </div>

        <!-- Empty State -->
        <div v-else class="text-gray-400 text-center mt-10">
            No messages yet. Start the conversation!
        </div>
    </div>
</template>

<script setup>
const props = defineProps({
    messages: {
        type: Array,
        default: () => []
    },
    currentUserId: {
        type: [String, Number],
        required: true
    },
    backendBaseUrl: {
        type: String,
        required: true
    },
    userAvatar: {
        type: String,
        default: null
    },
    otherUserAvatar: {
        type: String,
        default: null
    }
});

// Get avatar URL with proper formatting
const getAvatarUrl = (avatarUrl) => {
    if (!avatarUrl) return null;
    return avatarUrl.startsWith('/') ? avatarUrl : `${props.backendBaseUrl}/static/${avatarUrl}`;
};

// Format timestamp
const formatTime = (ts) => {
    // If ts is a number, treat as timestamp
    if (typeof ts === "number") {
        return new Date(ts).toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
        });
    }
    
    // If ts is a string (Go time.Time), try to parse
    if (typeof ts === "string") {
        const parsed = Date.parse(ts);
        if (!isNaN(parsed)) {
            return new Date(parsed).toLocaleTimeString([], {
                hour: "2-digit",
                minute: "2-digit",
            });
        }
        // Fallback: show raw string
        return ts;
    }
    
    return "";
};
</script> 