<template>
    <section class="flex flex-col h-full w-full bg-gray-50 font-roboto">
        <!-- Header: Contact Info -->
        <div
            v-if="chatStore.currentChatUser"
            class="flex items-center gap-3 p-4 border-b bg-white"
        >
            <img
                v-if="chatStore.currentChatUser.avatar_url"
                :src="`${backendBaseUrl}/static/${chatStore.currentChatUser.avatar_url}`"
                class="w-10 h-10 rounded-full object-cover border"
                alt="Avatar"
            />
            <img
                v-else
                src="/src/assets/default-avatar.jpeg"
                class="w-10 h-10 rounded-full object-cover border"
                alt="Default Avatar"
            />
            <div>
                <div class="font-bold text-green-700">
                    {{ chatStore.currentChatUser.username }}
                </div>
                <div class="text-xs text-gray-500">Online</div>
            </div>
        </div>
        <div v-else class="flex-1 flex items-center justify-center bg-gray-200">
            <div class="text-gray-500 text-lg font-semibold">
                No chat selected. Please select a chat to start messaging.
            </div>
        </div>
        <!-- Messages -->
        <div
            v-if="chatStore.currentChatUser"
            class="flex-1 overflow-y-auto p-4 space-y-6"
        >
            <div
                v-if="
                    Array.isArray(chatStore.messages) &&
                    chatStore.messages.length
                "
            >
                <div
                    v-for="msg in chatStore.messages"
                    :key="msg.id"
                    :class="
                        msg.sender_id === userId
                            ? 'justify-end'
                            : 'justify-start'
                    "
                    class="flex mb-2"
                >
                    <div
                        :class="
                            msg.sender_id === userId
                                ? 'bg-green-500 text-white'
                                : 'bg-white text-gray-800 border'
                        "
                        class="rounded-lg px-4 py-2 max-w-xs"
                    >
                        <div
                            class="text-base font-semibold break-words whitespace-pre-line"
                        >
                            {{ msg.content }}
                        </div>
                        <div class="text-xs text-white text-right mt-1">
                            {{ formatTime(msg.timestamp) }}
                        </div>
                    </div>
                </div>
            </div>
            <div v-else class="text-gray-400 text-center mt-10">
                No messages yet. Start the conversation!
            </div>
        </div>
        <!-- Input -->
        <div
            v-if="chatStore.currentChatUser"
            class="p-4 border-t bg-white flex gap-2"
        >
            <input
                v-model="newMessage"
                type="text"
                placeholder="Type a message..."
                class="flex-1 border rounded px-3 py-2 focus:outline-none"
                @keyup.enter="sendMessage"
            />
            <button
                @click="sendMessage"
                class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
            >
                Send
            </button>
        </div>
    </section>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import { useChatStore } from "../stores/chat";
import axios from "../axiosInstance";
import { useUserStore } from "../stores/users";
const chatStore = useChatStore();
const userStore = useUserStore();
const userId = userStore.user_id;
const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;
const newMessage = ref("");
let chatSocket = null;

// Helper to get current userId and receiverId
function getChatIds() {
    const chat = chatStore.chats?.find(
        (c) =>
            c.participants &&
            c.participants.includes(chatStore.currentChatUser?.id)
    );
    if (!chat || !chat.participants || chat.participants.length < 2) return {};
    const senderId = chatStore.currentChatUser?.id;
    const receiverId = chat.participants.find((id) => id !== senderId);
    return { chatId: chat?.id, senderId, receiverId };
}

// Create a new chat with a target user
async function createChat(targetUserId) {
    try {
        const response = await axios.post("/api/chat/create", {
            target_user: targetUserId,
        });
        // No data returned, just proceed to fetch messages
        fetchMessages();
    } catch (error) {
        // Handle error (show toast, etc.)
        console.error("Failed to create chat:", error);
    }
}

function fetchMessages() {
    chatStore.setMessages([
        {
            id: 1,
            sender: userId,
            content: "Hello! How are you?",
            timestamp: Date.now() - 60000,
        },
        {
            id: 2,
            sender: chatStore.currentChatUser?.id,
            content: "I'm good, thanks! How about you?",
            timestamp: Date.now() - 30000,
        },
        {
            id: 3,
            sender: userId,
            content: "Doing well! Ready to chat.",
            timestamp: Date.now() - 10000,
        },
        {
            id: 4,
            sender: "123123",
            content: "Doing well! Ready to chat.",
            timestamp: Date.now() - 10000,
        },
    ]);
}

// Manage WebSocket connection when chat changes
watch(
    () => chatStore.currentChatUser?.id,
    (newId, oldId) => {
        if (chatSocket) {
            chatSocket.close();
            chatSocket = null;
        }
        if (newId) {
            const { chatId, senderId, receiverId } = getChatIds();
            if (!chatId || !senderId || !receiverId) return;
            const wsUrl = `${backendBaseUrl.replace(
                /^http/,
                "ws"
            )}/api/websocket/chat/add/${chatId}?sender_id=${senderId}&receiver_id=${receiverId}`;
            chatSocket = new WebSocket(wsUrl);
            chatSocket.onopen = () => {
                console.log("WebSocket connected");
            };
            chatSocket.onmessage = (event) => {
                const data = JSON.parse(event.data);
                const message = parseIncomingMessage(data);
                chatStore.addMessage(message);
            };
            chatSocket.onclose = () => {
                console.log("WebSocket closed");
            };
            chatSocket.onerror = (error) => {
                console.error("WebSocket error:", error);
            };
        }
    }
);

function parseIncomingMessage(data) {
    // Handles backend sending content as stringified JSON or plain text
    if (typeof data.content === "string" && data.content.startsWith("{")) {
        try {
            return JSON.parse(data.content);
        } catch (e) {
            return { content: data.content };
        }
    }
    return data;
}

function sendMessage() {
    if (!newMessage.value.trim() || !chatSocket || chatSocket.readyState !== 1)
        return;
    const { receiverId, chatId } = getChatIds();
    // Generate a temporary id and timestamp
    const tempId = `temp-${Date.now()}-${Math.floor(Math.random() * 10000)}`;
    const timestamp = new Date().toISOString();
    const msgObj = {
        id: tempId,
        chat_id: chatId,
        sender_id: userId,
        receiver_id: receiverId,
        content: newMessage.value,
        timestamp,
    };
    chatSocket.send(JSON.stringify(msgObj));
    // Only add optimistic message if backend does NOT echo it back immediately
    // chatStore.addMessage(msgObj);
    newMessage.value = "";
}

function formatTime(ts) {
    // If ts is a number, treat as timestamp
    if (typeof ts === "number") {
        return new Date(ts).toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
        });
    }
    // If ts is a string (Go time.Time), try to parse
    if (typeof ts === "string") {
        // Try to parse ISO8601/RFC3339
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
}
</script>

<style scoped>
.font-roboto {
    font-family: "Roboto", Arial, sans-serif;
}
</style>
