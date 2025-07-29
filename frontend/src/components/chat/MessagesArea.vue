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
                                :src="isSecretChat ? '/src/assets/default-secret-chat-avatar.jpg' : getAvatarUrl(otherUserAvatar)"
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

        <!-- Secret Chat Not Approved Warning -->
        <div
            v-if="isSecretChat && !isSecretChatApproved"
            class="flex flex-col items-center justify-center h-full text-center py-12"
        >
            <div
                class="w-24 h-24 bg-orange-100 rounded-full flex items-center justify-center mb-6"
            >
                <span class="material-icons text-orange-500 text-4xl"
                    >lock</span
                >
            </div>
            <h3 class="text-xl font-semibold text-gray-700 mb-2">
                Secret Chat Pending Approval
            </h3>
            <p class="text-gray-500 text-sm max-w-md mb-4">
                This secret chat is waiting for the other user to approve it.
                You'll be able to send messages once it's approved.
            </p>
            <div
                class="bg-orange-50 border border-orange-200 rounded-lg p-4 max-w-md"
            >
                <div class="flex items-center gap-2 text-orange-700 mb-2">
                    <span class="material-icons text-sm">info</span>
                    <span class="text-sm font-medium">What happens next?</span>
                </div>
                <ul class="text-xs text-orange-600 space-y-1 text-left">
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">schedule</span>
                        <span>The other user will receive a notification</span>
                    </li>
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">check_circle</span>
                        <span>They need to approve the secret chat</span>
                    </li>
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">security</span>
                        <span>Once approved, messages will be encrypted</span>
                    </li>
                </ul>
            </div>
        </div>

        <!-- Empty State -->
        <div
            v-else-if="!messages || !messages.length"
            class="flex flex-col items-center justify-center h-full text-center py-12"
        >
            <div
                class="w-24 h-24 bg-gray-200 rounded-full flex items-center justify-center mb-6"
            >
                <span class="material-icons text-gray-400 text-4xl">chat</span>
            </div>
            <h3 class="text-xl font-semibold text-gray-700 mb-2">
                No messages yet
            </h3>
            <p class="text-gray-500 text-sm max-w-md">
                Start the conversation by sending your first message!
            </p>
        </div>

        <!-- Delete Confirmation Modal -->
        <div
            v-if="showDeleteModal"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-80 backdrop-blur-md"
        >
            <div
                class="bg-white w-96 p-8 rounded-2xl shadow-2xl border border-gray-100 text-center space-y-6"
            >
                <div
                    class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4"
                >
                    <span class="material-icons text-red-600 text-2xl"
                        >delete</span
                    >
                </div>
                <h3 class="text-xl font-bold text-gray-800">Delete Message</h3>
                <p class="text-base text-gray-600 leading-relaxed">
                    Do you want to delete this message just for yourself or for
                    everyone?
                </p>
                <div class="flex flex-col gap-3 text-base">
                    <button
                        @click="deleteForSender"
                        class="cursor-pointer bg-gray-100 hover:bg-gray-200 text-gray-800 py-3 px-6 rounded-lg font-semibold transition-colors duration-200"
                    >
                        Delete for me
                    </button>
                    <button
                        @click="deleteForAll"
                        class="cursor-pointer bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white py-3 px-6 rounded-lg font-semibold transition-all duration-200 shadow-md"
                    >
                        Delete for everyone
                    </button>
                    <button
                        @click="cancelDelete"
                        class="cursor-pointer text-gray-500 hover:text-gray-700 hover:underline text-sm py-2"
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
    isSecretChat: {
        type: Boolean,
        default: false,
    },
    isSecretChatApproved: {
        type: Boolean,
        default: true,
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
    const messageToRestore = chatStore.messages.find(
        (msg) => msg.id === messageToDeleteId.value
    );

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
