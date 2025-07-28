<template>
    <div
        ref="messagesContainer"
        class="flex-1 overflow-y-auto px-6 py-4 space-y-4 bg-gray-50"
        @scroll="handleScroll"
    >
        <!-- Loading indicator for older messages -->
        <div v-if="isLoadingMessages" class="text-center py-4">
            <div
                class="inline-flex items-center px-4 py-2 bg-gray-100 rounded-lg"
            >
                <svg
                    class="animate-spin -ml-1 mr-3 h-5 w-5 text-gray-500"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                >
                    <circle
                        class="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        stroke-width="4"
                    />
                    <path
                        class="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    />
                </svg>
                <span class="text-gray-600 text-sm font-medium"
                    >Loading messages...</span
                >
            </div>
        </div>

        <!-- Messages List -->
        <div v-if="messages && messages.length" class="space-y-1">
            <template v-for="msg in messages" :key="msg.id">
                <div
                    v-if="msg.content !== '' || msg.content_address !== ''"
                    :class="
                        msg.sender_id === currentUserId
                            ? 'justify-end'
                            : 'justify-start'
                    "
                    class="flex items-end gap-2"
                >
                    <!-- Avatar for received messages -->
                    <template v-if="msg.sender_id !== currentUserId">
                        <img
                            :src="getAvatarUrl(otherUserAvatar)"
                            class="w-9 h-9 rounded-full object-cover border shadow-sm select-none pointer-events-none"
                            alt="Avatar"
                        />
                    </template>

                    <!-- Message Bubble -->
                    <div
                        :class="
                            msg.sender_id === currentUserId
                                ? 'bg-green-500 text-white'
                                : 'bg-white text-gray-900 border border-gray-200'
                        "
                        class="rounded-2xl px-4 py-2 max-w-xs sm:max-w-md shadow-md relative group"
                    >
                        <div
                            class="text-base break-words whitespace-pre-wrap leading-snug"
                        >
                            {{ msg.content }}
                        </div>
                        <div
                            :class="
                                msg.sender_id === currentUserId
                                    ? 'text-white'
                                    : 'text-gray-500'
                            "
                            class="text-xs text-right mt-1"
                        >
                            {{
                                formatTime(msg.created_at) ||
                                formatTime(Date.now())
                            }}
                        </div>

                        <!-- Delete button -->
                        <div
                            v-if="canDeleteMessage(msg.id, currentUserId)"
                            class="absolute -top-2 -right-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
                        >
                            <button
                                @click="handleDeleteMessage(msg.id)"
                                :disabled="isDeleting"
                                class="bg-red-500 hover:bg-red-600 text-white w-6 h-6 rounded-full flex items-center justify-center text-xs shadow-lg cursor-pointer"
                                title="Delete message"
                            >
                                ×
                            </button>
                        </div>
                    </div>

                    <!-- Avatar for sent messages -->
                    <template v-if="msg.sender_id === currentUserId">
                        <img
                            :src="getAvatarUrl(userAvatar)"
                            class="w-9 h-9 rounded-full object-cover border shadow-sm select-none pointer-events-none"
                            alt="Avatar"
                        />
                    </template>
                </div>
            </template>
        </div>

        <!-- Empty State -->
        <div v-else class="text-gray-400 text-center mt-10 text-sm italic">
            No messages yet. Start the conversation!
        </div>

        <!-- Delete Confirmation Modal -->
        <div
            v-if="showDeleteModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-gray-200 bg-opacity-60"
        >
            <div
                class="bg-white w-80 p-6 rounded-2xl shadow-2xl border border-gray-100 text-center space-y-5"
            >
                <h3 class="text-xl font-bold text-gray-800">Delete Message</h3>
                <p class="text-base text-gray-600 leading-relaxed">
                    Do you want to delete this message just for yourself or for
                    everyone?
                </p>
                <div class="flex flex-col gap-3 text-base">
                    <button
                        @click="deleteForSender"
                        class="cursor-pointer bg-gray-100 hover:bg-gray-200 text-gray-800 py-2 px-4 rounded-md font-semibold"
                    >
                        Delete for me
                    </button>
                    <button
                        @click="deleteForAll"
                        class="cursor-pointer bg-green-500 hover:bg-green-600 text-white py-2 px-4 rounded-md font-semibold"
                    >
                        Delete for everyone
                    </button>
                    <button
                        @click="cancelDelete"
                        class="cursor-pointer text-gray-500 hover:underline text-sm"
                    >
                        Cancel
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from "vue";
import { useMessagePagination } from "../../composables/useMessagePagination";
import { useMessageDeletion } from "../../composables/useMessageDeletion";
import axiosInstance from "../../axiosInstance";
import { showError, showMessage } from "../../utils/toast";
import { useChatStore } from "../../stores/chat";

const props = defineProps({
    messages: {
        type: Array,
        default: () => [],
    },
    currentUserId: {
        type: [String, Number],
        required: true,
    },
    backendBaseUrl: {
        type: String,
        required: true,
    },
    userAvatar: {
        type: String,
        default: null,
    },
    otherUserAvatar: {
        type: String,
        default: null,
    },
    chatId: {
        type: String,
        default: null,
    },
});

const emit = defineEmits(["load-more-messages", "remove-message"]);

const { isLoadingMessages, shouldLoadMore } = useMessagePagination();
const { isDeleting, deleteMessage, canDeleteMessage } = useMessageDeletion();

const messagesContainer = ref(null);

const handleScroll = (event) => {
    const { scrollTop, scrollHeight, clientHeight } = event.target;
    if (shouldLoadMore(scrollTop, scrollHeight, clientHeight)) {
        emit("load-more-messages");
    }
};

const scrollToBottom = () => {
    if (messagesContainer.value) {
        messagesContainer.value.scrollTop =
            messagesContainer.value.scrollHeight;
    }
};

watch(
    () => props.messages.length,
    () => {
        nextTick(() => {
            scrollToBottom();
        });
    }
);

onMounted(() => {
    nextTick(() => {
        scrollToBottom();
    });
});

const handleDeleteMessage = (messageId) => {
    messageToDeleteId.value = messageId;
    showDeleteModal.value = true;
};

const getAvatarUrl = (avatarUrl) => {
    if (!avatarUrl) return "/src/assets/default-avatar.jpg";
    return avatarUrl.startsWith("/")
        ? avatarUrl
        : `${props.backendBaseUrl}/static/${avatarUrl}`;
};

const formatTime = (ts) => {
    if (!ts) return "";
    if (typeof ts === "number") {
        return new Date(ts).toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
        });
    }
    if (typeof ts === "string") {
        const parsed = Date.parse(ts);
        if (!isNaN(parsed)) {
            return new Date(parsed).toLocaleTimeString([], {
                hour: "2-digit",
                minute: "2-digit",
            });
        }
        return ts;
    }
    return "";
};

const showDeleteModal = ref(false);
const messageToDeleteId = ref(null);

const cancelDelete = () => {
    showDeleteModal.value = false;
    messageToDeleteId.value = null;
};

const deleteForSender = async () => {
    if (messageToDeleteId.value && props.chatId) {
        await deleteMessageByType("sender");
    }
};

const deleteForAll = async () => {
    if (messageToDeleteId.value && props.chatId) {
        await deleteMessageByType("all");
    }
};

const deleteMessageByType = async (type) => {
    let url;
    switch (type) {
        case "sender":
            url = `/api/message/delete/sender/${messageToDeleteId.value}`;
            break;
        case "all":
            url = `/api/message/delete/all/${messageToDeleteId.value}`;
            break;
        default:
            return;
    }

    // Store the message for potential restoration
    const chatStore = useChatStore();
    const messageToRestore = chatStore.messages.find(msg => msg.id === messageToDeleteId.value);
    
    // Immediately remove message from Vue array for better UX
    const success = chatStore.deleteMessage(messageToDeleteId.value);
    if (!success) {
        showError("Failed to delete message from local store");
        return;
    }

    try {
        await axiosInstance.delete(url);
        showMessage("Message deleted successfully");
        emit("remove-message", messageToDeleteId.value);
    } catch (err) {
        console.error("Delete failed:", err);
        showError("Failed to delete message from server");
        
        // Restore the message if backend deletion failed
        if (messageToRestore) {
            chatStore.addMessage(messageToRestore);
            showError("Message restored - deletion failed on server");
        }
    } finally {
        showDeleteModal.value = false;
        messageToDeleteId.value = null;
    }
};
</script>
