<template>
    <section class="w-1/3 bg-white border-r p-4 overflow-y-auto">
        <div v-if="activeTab === 'chats'">
            <h2 class="text-lg font-bold text-green-600 mb-4">Chats</h2>
            <div class="flex gap-2 mb-4">
                <button
                    class="bg-green-500 text-white px-3 py-1 rounded hover:bg-green-600"
                    @click="showCreateChat = true"
                >
                    Create Chat
                </button>
                <button
                    class="bg-gray-300 text-gray-700 px-3 py-1 rounded cursor-not-allowed"
                    disabled
                >
                    Create Secret Chat
                </button>
            </div>
            <!-- Chat creation modal -->
            <div
                v-if="showCreateChat"
                class="fixed inset-0 bg-gray-900 bg-opacity-20 flex items-center justify-center z-50"
            >
                <div class="bg-white rounded-lg shadow-lg p-6 w-80 relative">
                    <button
                        class="absolute top-2 right-2 text-gray-400 hover:text-gray-700"
                        @click="showCreateChat = false"
                    >
                        <span class="material-icons">close</span>
                    </button>
                    <h3 class="text-lg font-bold text-green-600 mb-4">
                        Start New Chat
                    </h3>
                    <label class="block mb-2 text-green-700 font-semibold"
                        >Enter Username</label
                    >
                    <input
                        v-model="enteredUsername"
                        type="text"
                        placeholder="Type username..."
                        class="w-full border rounded px-2 py-1 mb-4"
                    />
                    <button
                        class="bg-green-500 text-white px-3 py-1 rounded hover:bg-green-600 w-full"
                        @click="searchUser"
                        :disabled="!enteredUsername"
                    >
                        Search
                    </button>
                    <div
                        v-if="users.length"
                        class="mt-4 p-3 rounded bg-gray-50 border cursor-pointer hover:bg-gray-200 transition"
                        @click="openChatWithUser(users[0])"
                    >
                        <div class="flex items-center gap-2 mb-2">
                            <img
                                v-if="users[0].avatar_url"
                                :src="`${backendBaseUrl}/static/${users[0].avatar_url}`"
                                alt="Avatar"
                                class="w-10 h-10 rounded-full object-cover border"
                            />
                            <span class="font-bold text-green-700">{{
                                users[0].username
                            }}</span>
                        </div>
                        <div class="text-xs text-gray-600 mb-1">
                            Created:
                            {{ new Date(users[0].created_at).toLocaleString() }}
                        </div>
                    </div>
                </div>
            </div>
            <!-- End of chat creation modal -->
            <div
                v-if="chatStore.chats && chatStore.chats.length === 0"
                class="text-gray-500"
            >
                No chats yet.
            </div>
            <div v-else-if="chatStore.chats && chatStore.chats.length">
                <div
                    v-for="chat in chatStore.chats"
                    :key="chat.id"
                    class="mb-2 p-2 border rounded cursor-pointer hover:bg-gray-100"
                    @click="handleChatClick(chat)"
                >
                    <div class="flex items-center gap-2">
                        <img
                            v-if="
                                chatStore.avatarUrls &&
                                chatStore.avatarUrls[chat.id]
                            "
                            :src="`${backendBaseUrl}/static/${
                                chatStore.avatarUrls[chat.id]
                            }`"
                            alt="Avatar"
                            class="w-8 h-8 rounded-full object-cover border"
                        />
                        <span class="font-bold text-green-900">{{
                            chatStore.usernames[chat.id]
                        }}</span>
                    </div>
                    <div class="text-xs text-gray-600">
                        Created:
                        {{ new Date(chat.created_at).toLocaleString() }}
                    </div>
                </div>
            </div>
        </div>
        <div v-else-if="activeTab === 'groups'">
            <h2 class="text-lg font-bold text-green-600 mb-4">Groups</h2>
            <!-- List of groups will go here -->
            <div class="text-gray-500">No groups yet.</div>
        </div>
        <div v-else-if="activeTab === 'settings'">
            <h2 class="text-lg font-bold text-green-600 mb-4">User Settings</h2>
            <UserAvatarUpload />
        </div>
    </section>
</template>

<script setup>
import { defineProps, ref, onMounted } from "vue";
import UserAvatarUpload from "./UserAvatarUpload.vue";
import axiosInstance from "@/axiosInstance";
import { showError } from "@/utils/toast";
import { useChatStore } from "../stores/chat";
import { useUserStore } from "../stores/users";

const props = defineProps({ activeTab: String });
const showCreateChat = ref(false);
const enteredUsername = ref("");
const users = ref([]);
const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;
const chatStore = useChatStore();

const userStore = useUserStore();

function handleChatClick(chat) {
    const currentUserId = userStore.user_id;
    if (!chat.participants || chat.participants.length < 2 || !currentUserId) {
        return;
    }

    const otherUserId = chat.participants.find((id) => id !== currentUserId);
    if (!otherUserId) return;
    openChatWithUser({ id: otherUserId });
}

onMounted(async () => {
    try {
        const response = await axiosInstance.get("/api/user/get-chats");
        chatStore.setChats(response.data.chats);
        chatStore.setAvatarUrls(response.data.avatar_urls);
        chatStore.setUsernames(response.data.usernames);
    } catch (error) {
        console.error("Failed to fetch chats:", error);
    }
});

function searchUser() {
    axiosInstance
        .get(`/api/user/search?q=${enteredUsername.value}`)
        .then((resp) => {
            if (!resp.data || Object.keys(resp.data).length === 0) {
                users.value = [];
                showError("Username is invalid or not found.");
            } else {
                users.value = [resp.data];
            }
        })
        .catch((error) => {
            console.error("Error fetching user:", error);
            showError("Failed to fetch user. Please try again.");
        });
}

function openChatWithUser(user) {
    const chatStore = useChatStore();
    chatStore.setChatUser(user);
    showCreateChat.value = false;

    // Find the selected chat object by participants
    const selectedChat = chatStore.chats.find(
        (chat) => chat.participants && chat.participants.includes(user.id)
    );
    if (
        !selectedChat ||
        !selectedChat.participants ||
        selectedChat.participants.length < 2
    )
        return;
    const chatId = selectedChat.id;
    // senderId is the current user (user.id), receiverId is the other participant
    const senderId = user.id;
    const receiverId = selectedChat.participants.find((id) => id !== user.id);
    if (!chatId || !senderId || !receiverId) return;

    // Establish WebSocket connection
    const wsUrl = `${backendBaseUrl.replace(
        /^http/,
        "ws"
    )}/api/websocket/chat/add/${chatId}?sender_id=${senderId}&receiver_id=${receiverId}`;
    const chatSocket = new WebSocket(wsUrl);

    chatSocket.onopen = () => {
        console.log("WebSocket connected");
    };
    chatSocket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        // Handle incoming message, e.g.:
        chatStore.addMessage(data);
    };
    chatSocket.onclose = () => {
        console.log("WebSocket closed");
    };
    chatSocket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
}
</script>
