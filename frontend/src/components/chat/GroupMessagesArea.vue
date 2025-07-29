<template>
    <div
        ref="messagesContainer"
        class="flex-1 overflow-y-auto px-4 py-6 space-y-3 bg-gradient-to-br from-gray-50 via-blue-50 to-indigo-50"
        @scroll="handleScroll"
    >
        <!-- Loading indicator for older messages -->
        <div v-if="isLoadingMessages" class="text-center py-6">
            <div
                class="inline-flex items-center px-6 py-3 bg-white rounded-full shadow-lg border border-gray-200"
            >
                <div
                    class="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-500 mr-3"
                ></div>
                <span class="text-gray-700 text-sm font-medium"
                    >Loading messages...</span
                >
            </div>
        </div>

        <!-- Messages List -->
        <div v-if="messages && messages.length" class="space-y-3">
            <template v-for="msg in messages" :key="msg.id">
                <div
                    v-if="msg.content !== '' || msg.content_address !== ''"
                    :class="
                        msg.sender_id === currentUserId
                            ? 'justify-end'
                            : 'justify-start'
                    "
                    class="flex items-end gap-3 px-2"
                >
                    <!-- Avatar for received messages -->
                    <template v-if="msg.sender_id !== currentUserId">
                        <div class="flex-shrink-0">
                            <img
                                :src="getAvatarUrl(msg.sender_avatar || otherUserAvatar)"
                                class="w-10 h-10 rounded-full object-cover border-2 border-white shadow-md select-none pointer-events-none"
                                alt="Avatar"
                            />
                        </div>
                    </template>

                    <!-- Message Bubble -->
                    <div
                        :class="
                            msg.sender_id === currentUserId
                                ? 'bg-gradient-to-r from-blue-500 to-blue-600 text-white shadow-lg'
                                : 'bg-white text-gray-900 border border-gray-200 shadow-md'
                        "
                        class="rounded-2xl px-5 py-3 max-w-xs sm:max-w-md lg:max-w-lg relative group hover:shadow-lg transition-all duration-200"
                    >
                        <!-- Sender name for group messages -->
                        <div v-if="msg.sender_id !== currentUserId" class="text-xs text-gray-500 mb-1 font-medium">
                            {{ msg.sender_name || 'Unknown User' }}
                        </div>
                        
                        <div
                            class="text-base break-words whitespace-pre-wrap leading-relaxed font-medium"
                        >
                            {{ msg.content }}
                        </div>
                        <div
                            :class="
                                msg.sender_id === currentUserId
                                    ? 'text-blue-100'
                                    : 'text-gray-500'
                            "
                            class="text-xs text-right mt-2 font-medium"
                        >
                            {{
                                formatTime(msg.created_at) ||
                                formatTime(Date.now())
                            }}
                        </div>

                        <!-- Delete button -->
                        <div
                            v-if="canDeleteMessage(msg.id, currentUserId)"
                            class="absolute -top-2 -right-2 opacity-0 group-hover:opacity-100 transition-all duration-200 transform scale-90 group-hover:scale-100"
                        >
                            <button
                                @click="handleDeleteMessage(msg.id)"
                                :disabled="isDeleting"
                                class="bg-red-500 hover:bg-red-600 text-white w-8 h-8 rounded-full flex items-center justify-center text-sm shadow-lg cursor-pointer transition-all duration-200 hover:shadow-xl"
                                title="Delete message"
                            >
                                <span class="material-icons text-sm"
                                    >delete</span
                                >
                            </button>
                        </div>
                    </div>

                    <!-- Avatar for sent messages -->
                    <template v-if="msg.sender_id === currentUserId">
                        <div class="flex-shrink-0">
                            <img
                                :src="getAvatarUrl(userAvatar)"
                                class="w-10 h-10 rounded-full object-cover border-2 border-white shadow-md select-none pointer-events-none"
                                alt="Avatar"
                            />
                        </div>
                    </template>
                </div>
            </template>
        </div>

        <!-- Empty state -->
        <div v-else class="text-center py-8">
            <div class="text-gray-500">
                <span class="material-icons text-6xl mb-4 block">chat_bubble_outline</span>
                <p class="text-lg font-medium">No messages yet</p>
                <p class="text-sm">Start the conversation!</p>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue';
import { useMessageDeletion } from '../../composables/useMessageDeletion';

const props = defineProps({
    messages: {
        type: Array,
        default: () => []
    },
    currentUserId: {
        type: String,
        required: true
    },
    backendBaseUrl: {
        type: String,
        required: true
    },
    userAvatar: {
        type: String,
        default: ''
    },
    otherUserAvatar: {
        type: String,
        default: ''
    },
    chatId: {
        type: String,
        required: true
    },
    isLoadingMessages: {
        type: Boolean,
        default: false
    }
});

const emit = defineEmits(['load-more-messages']);

const messagesContainer = ref(null);
const { handleDeleteMessage, isDeleting } = useMessageDeletion();

const getAvatarUrl = (avatarUrl) => {
    if (!avatarUrl) return '/src/assets/default-avatar.jpg';
    if (avatarUrl.startsWith('http')) return avatarUrl;
    return `${props.backendBaseUrl}/static/${avatarUrl}`;
};

const formatTime = (timestamp) => {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

const canDeleteMessage = (messageId, currentUserId) => {
    // For now, allow deletion of own messages
    return messageId && messageId.startsWith('temp-');
};

const handleScroll = () => {
    if (messagesContainer.value) {
        const { scrollTop } = messagesContainer.value;
        if (scrollTop === 0) {
            emit('load-more-messages');
        }
    }
};

onMounted(() => {
    nextTick(() => {
        if (messagesContainer.value) {
            messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
        }
    });
});
</script> 