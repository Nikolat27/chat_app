<template>
    <div
        ref="messagesContainer"
        class="flex-1 overflow-y-auto p-4 space-y-6"
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
                    ></circle>
                    <path
                        class="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                </svg>
                <span class="text-gray-600 text-sm font-medium"
                    >Loading messages...</span
                >
            </div>
        </div>

        <!-- Messages List -->
        <div v-if="messages && messages.length">
            <template v-for="msg in messages" :key="msg.id">
                <div
                    v-if="msg.content !== '' || msg.content_address !== ''"
                    :class="
                        msg.sender_id === currentUserId
                            ? 'justify-end'
                            : 'justify-start'
                    "
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
                        :class="
                            msg.sender_id === currentUserId
                                ? 'bg-green-500 text-white'
                                : 'bg-white text-gray-800 border'
                        "
                        class="rounded-lg px-4 py-2 max-w-xs relative group"
                    >
                        <div
                            class="text-base font-semibold break-words whitespace-pre-line"
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
                            {{ formatTime(msg.created_at) }}
                        </div>

                        <!-- Delete button (only for own messages with real IDs) -->
                        <div
                            v-if="canDeleteMessage(msg.id, currentUserId)"
                            class="absolute -top-2 -right-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
                        >
                            <button
                                @click="handleDeleteMessage(msg.id)"
                                :disabled="isDeleting"
                                class="cursor-pointer bg-red-500 hover:bg-red-600 text-white rounded-full w-6 h-6 flex items-center justify-center text-xs shadow-lg"
                                title="Delete message"
                            >
                                ×
                            </button>
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
            </template>
        </div>

        <!-- Empty State -->
        <div v-else class="text-gray-400 text-center mt-10">
            No messages yet. Start the conversation!
        </div>

        <!-- Delete Confirmation Modal -->
        <div
            v-if="showDeleteModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
        >
            <div
                class="bg-white rounded-xl shadow-xl w-72 p-6 space-y-4 text-center"
            >
                <h3 class="text-lg font-semibold text-gray-800">
                    Delete Message
                </h3>
                <p class="text-sm text-gray-600">
                    Do you want to delete this message just for yourself or for
                    everyone?
                </p>
                <div class="flex flex-col gap-2">
                    <button
                        @click="deleteForSender"
                        class="cursor-pointer bg-gray-200 hover:bg-gray-300 text-gray-800 py-2 px-4 rounded-md"
                    >
                        Delete for me
                    </button>
                    <button
                        @click="deleteForAll"
                        class="cursor-pointer bg-red-500 hover:bg-red-600 text-white py-2 px-4 rounded-md"
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

// Pagination composable
const { isLoadingMessages, shouldLoadMore, loadNextPage } =
    useMessagePagination();

// Message deletion composable
const { isDeleting, deleteMessage, canDeleteMessage } = useMessageDeletion();

// Template refs
const messagesContainer = ref(null);

// Scroll handler for infinite scroll
const handleScroll = (event) => {
    const { scrollTop, scrollHeight, clientHeight } = event.target;

    console.log(
        `Scroll: top=${scrollTop}, height=${scrollHeight}, client=${clientHeight}, shouldLoad=${shouldLoadMore(
            scrollTop,
            scrollHeight,
            clientHeight
        )}`
    );

    if (shouldLoadMore(scrollTop, scrollHeight, clientHeight)) {
        emit("load-more-messages");
    }
};

// Auto-scroll to bottom when new messages arrive
const scrollToBottom = () => {
    if (messagesContainer.value) {
        messagesContainer.value.scrollTop =
            messagesContainer.value.scrollHeight;
    }
};

// Watch for new messages and scroll to bottom
watch(
    () => props.messages.length,
    () => {
        // Use nextTick to ensure DOM is updated
        nextTick(() => {
            scrollToBottom();
        });
    }
);

// Scroll to bottom when component mounts
onMounted(() => {
    nextTick(() => {
        scrollToBottom();
    });
});

// Handle message deletion
const handleDeleteMessage = (messageId) => {
    messageToDeleteId.value = messageId;
    showDeleteModal.value = true;
};

// Get avatar URL with proper formatting
const getAvatarUrl = (avatarUrl) => {
    if (!avatarUrl) return null;
    return avatarUrl.startsWith("/")
        ? avatarUrl
        : `${props.backendBaseUrl}/static/${avatarUrl}`;
};

// Format timestamp
const formatTime = (ts) => {
    if (!ts) return "";

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

const showDeleteModal = ref(false);
const messageToDeleteId = ref(null);

// Cancel modal
const cancelDelete = () => {
    showDeleteModal.value = false;
    messageToDeleteId.value = null;
};

// Delete for sender only
const deleteForSender = async () => {
    if (messageToDeleteId.value && props.chatId) {
        await deleteMessageByType("sender");
    }
};

// Delete for all
const deleteForAll = async () => {
    if (messageToDeleteId.value && props.chatId) {
        await deleteMessageByType("all");
    }
};

// Generic delete handler
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

    try {
        await axiosInstance.delete(url);

        showMessage("Message deleted successfully. Reload if you want");
        emit("remove-message", messageToDeleteId.value);
    } catch (err) {
        console.error("Delete failed:", err);
    } finally {
        showDeleteModal.value = false;
        messageToDeleteId.value = null;
    }
};
</script>
