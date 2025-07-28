<template>
    <div
        class="fixed inset-0 bg-gradient-to-br from-green-100 via-white to-green-200 bg-opacity-80 flex items-center justify-center z-50"
    >
        <div
            class="bg-white rounded-2xl shadow-2xl p-8 w-96 relative font-sans"
        >
            <!-- Close Button -->
            <button
                class="absolute top-3 right-3 text-gray-400 hover:text-green-700 text-2xl transition-colors cursor-pointer"
                @click="$emit('close')"
                aria-label="Close"
            >
                <span class="material-icons">close</span>
            </button>

            <!-- Header -->
            <div class="mb-6">
                <h3
                    class="text-2xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-green-600 to-emerald-500 mb-2 tracking-tight"
                >
                    Start New Chat
                </h3>
                <p class="text-sm text-gray-500">
                    Find a user to start a conversation.
                </p>
            </div>

            <!-- Search Form -->
            <div class="mb-6">
                <label class="block mb-2 text-green-700 font-semibold">
                    Enter Username
                </label>
                <input
                    v-model="enteredUsername"
                    type="text"
                    placeholder="Type username..."
                    class="w-full border border-green-200 rounded-xl px-3 py-2 mb-4 focus:outline-none focus:ring-2 focus:ring-green-300 transition"
                    @keyup.enter="searchUser"
                />
                <button
                    class="bg-gradient-to-r from-green-500 to-emerald-500 text-white px-4 py-2 rounded-xl w-full font-semibold shadow hover:scale-105 hover:shadow-lg transition-all duration-150 cursor-pointer"
                    @click="searchUser"
                    :disabled="!enteredUsername"
                >
                    <span class="material-icons align-middle mr-1 text-base"
                        >search</span
                    >
                    Search
                </button>
                <div
                    v-if="selfSearchWarning"
                    class="mt-3 p-2 bg-red-100 text-red-700 rounded text-center text-sm font-semibold border border-red-200"
                >
                    You can't start a chat with yourself!
                </div>
            </div>

            <!-- Search Results -->
            <div
                v-if="users.length && !selfSearchWarning"
                class="mt-4 p-4 rounded-xl bg-green-50 border border-green-200 cursor-pointer hover:bg-green-100 hover:scale-[1.03] transition-all flex flex-col gap-2 items-start shadow-sm"
                @click="selectUser(users[0])"
            >
                <div class="flex items-center gap-3 mb-1">
                    <img
                        v-if="users[0].avatar_url"
                        :src="
                            users[0].avatar_url.startsWith('http')
                                ? users[0].avatar_url
                                : `${backendBaseUrl}/static/${users[0].avatar_url}`
                        "
                        alt="Avatar"
                        class="w-12 h-12 rounded-full object-cover border border-green-200 shadow select-none pointer-events-none"
                    />
                    <img
                        v-else
                        src="/src/assets/default-avatar.jpg"
                        alt="Default Avatar"
                        class="w-12 h-12 rounded-full object-cover border border-green-200 shadow select-none pointer-events-none"
                    />
                    <span class="font-bold text-green-700 text-lg">
                        {{ users[0].username }}
                    </span>
                </div>
                <div class="text-xs text-gray-500 mb-1">
                    Created:
                    {{
                        users[0].created_at
                            ? new Date(users[0].created_at).toLocaleString()
                            : ""
                    }}
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from "vue";
import axiosInstance from "../../axiosInstance";
import { showError, showMessage } from "../../utils/toast";

const props = defineProps({
    backendBaseUrl: String,
});

const emit = defineEmits(["close", "user-selected"]);

const enteredUsername = ref("");
const users = ref([]);
const selfSearchWarning = ref(false);

// Get current user's username from localStorage or a global store
let currentUsername = null;
try {
    const user = JSON.parse(localStorage.getItem("user"));
    currentUsername = user?.username || null;
} catch (e) {
    currentUsername = null;
}

// Search for user by username
const searchUser = async () => {
    selfSearchWarning.value = false;
    users.value = [];
    if (!enteredUsername.value.trim()) return;

    // Check if searching for self
    if (
        currentUsername &&
        enteredUsername.value.trim().toLowerCase() ===
            currentUsername.toLowerCase()
    ) {
        selfSearchWarning.value = true;
        return;
    }

    try {
        const response = await axiosInstance.get(
            `/api/user/search?q=${enteredUsername.value}`
        );
        if (!response.data || Object.keys(response.data).length === 0) {
            users.value = [];
            showError("Username is invalid or not found.");
        } else {
            users.value = [response.data];
        }
    } catch (error) {
        users.value = [];
        if (error?.response?.data?.detail) {
            showError(error.response.data.detail);
        } else {
            showError("Failed to fetch user. Please try again.");
        }
    }
};

const selectUser = (user) => {
    emit("user-selected", user);
};
</script>
