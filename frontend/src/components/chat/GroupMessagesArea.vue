<template>
    <div
        ref="messagesContainer"
        class="flex-1 overflow-y-auto bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50"
        @scroll="handleScroll"
    >
        <!-- Secret Key Error Banner -->
        <div v-if="showSecretKeyError" class="sticky top-0 z-20 bg-red-50 border-b border-red-200">
            <div class="flex items-center justify-between p-4">
                <div class="flex items-center gap-3">
                    <span class="material-icons text-red-600">warning</span>
                    <div>
                        <div class="text-red-800 font-medium text-sm">Secret Key Required</div>
                        <div class="text-red-600 text-xs">You need to enter the secret key to read encrypted messages</div>
                    </div>
                </div>
                <button
                    @click="openSecretKeyModal"
                    class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-sm font-medium"
                >
                    Enter Key
                </button>
            </div>
        </div>

        <!-- Loading indicator for older messages -->
        <div v-if="isLoadingMessages" class="sticky top-0 z-10 bg-white/80 backdrop-blur-sm border-b border-gray-200">
            <div class="flex items-center justify-center py-4">
                <div class="inline-flex items-center px-4 py-2 bg-white rounded-full shadow-lg border border-gray-200">
                    <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-500 mr-2"></div>
                    <span class="text-gray-600 text-sm font-medium">Loading messages...</span>
                </div>
            </div>
        </div>

        <!-- Messages Container -->
        <div class="px-6 py-8 space-y-4">
            <!-- Messages List -->
            <template v-if="messages && messages.length">
                <div
                    v-for="msg in filteredMessages"
                    :key="msg.id"
                    :class="
                        msg.sender_id === currentUserId
                            ? 'flex justify-end'
                            : 'flex justify-start'
                    "
                    class="group"
                >
                    <div class="flex items-end gap-3 max-w-[70%]">
                        <!-- Avatar for received messages -->
                        <div v-if="msg.sender_id !== currentUserId" class="flex-shrink-0">
                            <img
                                :src="getAvatarUrl(msg.sender_avatar)"
                                class="w-8 h-8 rounded-full object-cover border-2 border-white shadow-sm"
                                alt="Avatar"
                                @error="handleAvatarError"
                            />
                        </div>

                        <!-- Message Content -->
                        <div class="flex flex-col">
                            <!-- Sender name for group messages -->
                            <div v-if="msg.sender_id !== currentUserId" class="flex items-center gap-2 mb-1">
                                <span class="text-xs font-semibold text-gray-700">
                                    {{ msg.sender_name || 'Unknown User' }}
                                </span>
                                <span class="text-xs text-gray-400">
                                    {{ formatTime(msg.created_at) }}
                                </span>
                            </div>

                            <!-- Message Bubble -->
                            <div
                                :class="
                                    msg.sender_id === currentUserId
                                        ? 'bg-gradient-to-r from-blue-500 to-blue-600 text-white'
                                        : 'bg-white text-gray-800 border border-gray-200'
                                "
                                class="rounded-2xl px-4 py-3 shadow-sm hover:shadow-md transition-all duration-200 relative"
                            >
                                <div class="text-base leading-relaxed break-words whitespace-pre-wrap">
                                    {{ msg.content }}
                                </div>

                                <!-- Time for sent messages -->
                                <div
                                    v-if="msg.sender_id === currentUserId"
                                    class="text-xs text-blue-100 mt-2 text-right"
                                >
                                    {{ formatTime(msg.created_at) }}
                                </div>

                                <!-- Delete button for sent messages -->
                                <div
                                    v-if="canDeleteMessage(msg.id, currentUserId) && msg.sender_id === currentUserId"
                                    class="absolute -top-1 -right-1 opacity-0 group-hover:opacity-100 transition-all duration-200"
                                >
                                    <button
                                        @click="handleDeleteMessage(msg.id)"
                                        :disabled="isDeleting"
                                        class="bg-red-500 hover:bg-red-600 text-white w-6 h-6 rounded-full flex items-center justify-center text-xs shadow-lg cursor-pointer transition-all duration-200"
                                        title="Delete message"
                                    >
                                        <span class="material-icons text-xs">delete</span>
                                    </button>
                                </div>
                            </div>
                        </div>

                        <!-- Avatar for sent messages -->
                        <div v-if="msg.sender_id === currentUserId" class="flex-shrink-0">
                            <img
                                :src="getAvatarUrl(userAvatar)"
                                class="w-8 h-8 rounded-full object-cover border-2 border-white shadow-sm"
                                alt="Avatar"
                                @error="handleAvatarError"
                            />
                        </div>
                    </div>
                </div>
            </template>

            <!-- Empty state -->
            <div v-else class="flex items-center justify-center h-full">
                <div class="text-center">
                    <div class="inline-flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mb-4">
                        <span class="material-icons text-2xl text-blue-600">chat_bubble_outline</span>
                    </div>
                    <h3 class="text-lg font-semibold text-gray-700 mb-2">No messages yet</h3>
                    <p class="text-sm text-gray-500">Start the conversation!</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch, computed } from 'vue';
import { useMessageDeletion } from '../../composables/useMessageDeletion';
import { useSecretGroupE2EE } from '../../composables/useSecretGroupE2EE';

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
    },
    isSecretGroup: {
        type: Boolean,
        default: false
    }
});

const emit = defineEmits(['load-more-messages', 'open-secret-key-modal']);

const messagesContainer = ref(null);
const { handleDeleteMessage, isDeleting } = useMessageDeletion();
const { hasGroupSecretKey } = useSecretGroupE2EE();
const previousMessageCount = ref(0);
const isPaginationLoading = ref(false);
const showSecretKeyError = ref(false);

// Check if user has entered the secret key
const checkSecretKey = async () => {
    if (!props.isSecretGroup) {
        showSecretKeyError.value = false;
        return;
    }
    
    try {
        const hasKey = await hasGroupSecretKey(props.chatId);
        showSecretKeyError.value = !hasKey;
    } catch (error) {
        console.error('Error checking secret key:', error);
        showSecretKeyError.value = true; // Show error if we can't check
    }
};

const openSecretKeyModal = () => {
    emit('open-secret-key-modal');
};

// Filter messages to only show valid ones
const filteredMessages = computed(() => {
    if (!props.messages) return [];
    return props.messages.filter(msg => 
        msg && (msg.content !== '' || msg.content_address !== '')
    );
});

const getAvatarUrl = (avatarUrl) => {
    if (!avatarUrl) {
        return '/src/assets/default-avatar.jpg';
    }
    if (avatarUrl.startsWith('http')) {
        return avatarUrl;
    }
    return avatarUrl;
};

const handleAvatarError = (event) => {
    event.target.src = '/src/assets/default-avatar.jpg';
};

const formatTime = (timestamp) => {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

const canDeleteMessage = (messageId, currentUserId) => {
    return messageId && messageId.startsWith('temp-');
};

const handleScroll = () => {
    if (messagesContainer.value) {
        const { scrollTop } = messagesContainer.value;
        if (scrollTop === 0) {
            isPaginationLoading.value = true;
            emit('load-more-messages');
        }
    }
};

const scrollToBottom = () => {
    if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
    }
};

const preserveScrollPosition = () => {
    if (messagesContainer.value) {
        const currentScrollTop = messagesContainer.value.scrollTop;
        const currentScrollHeight = messagesContainer.value.scrollHeight;
        nextTick(() => {
            if (messagesContainer.value) {
                const newScrollHeight = messagesContainer.value.scrollHeight;
                const heightDifference = newScrollHeight - currentScrollHeight;
                messagesContainer.value.scrollTop = currentScrollTop + heightDifference;
            }
        });
    }
};

// Watch for new messages and scroll to bottom
watch(
    () => props.messages.length,
    (newLength, oldLength) => {
        nextTick(() => {
            if (newLength > oldLength) {
                if (isPaginationLoading.value) {
                    preserveScrollPosition();
                    isPaginationLoading.value = false;
                } else {
                    scrollToBottom();
                }
            }
            previousMessageCount.value = newLength;
        });
    }
);

// Watch for new messages and scroll to bottom
watch(
    () => props.messages[props.messages.length - 1]?.id,
    () => {
        nextTick(() => {
            if (!isPaginationLoading.value) {
                scrollToBottom();
            }
        });
    }
);

onMounted(() => {
    scrollToBottom();
    checkSecretKey();
});

// Watch for changes in secret group status
watch(() => props.isSecretGroup, () => {
    checkSecretKey();
});

// Watch for changes in chatId
watch(() => props.chatId, () => {
    checkSecretKey();
});
</script>

<style scoped>
/* Custom scrollbar */
.overflow-y-auto::-webkit-scrollbar {
    width: 6px;
}

.overflow-y-auto::-webkit-scrollbar-track {
    background: transparent;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
    background: rgba(156, 163, 175, 0.3);
    border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
    background: rgba(156, 163, 175, 0.5);
}
</style> 