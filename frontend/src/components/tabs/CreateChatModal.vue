<template>
    <div class="fixed inset-0 bg-gray-900 bg-opacity-20 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg shadow-lg p-6 w-80 relative">
            <!-- Close Button -->
            <button
                class="absolute top-2 right-2 text-gray-400 hover:text-gray-700"
                @click="$emit('close')"
            >
                <span class="material-icons">close</span>
            </button>

            <!-- Header -->
            <h3 class="text-lg font-bold text-green-600 mb-4">
                Start New Chat
            </h3>

            <!-- Search Form -->
            <div class="mb-4">
                <label class="block mb-2 text-green-700 font-semibold">
                    Enter Username
                </label>
                <input
                    v-model="enteredUsername"
                    type="text"
                    placeholder="Type username..."
                    class="w-full border rounded px-2 py-1 mb-4"
                    @keyup.enter="searchUser"
                />
                <button
                    class="bg-green-500 text-white px-3 py-1 rounded hover:bg-green-600 w-full"
                    @click="searchUser"
                    :disabled="!enteredUsername"
                >
                    Search
                </button>
            </div>

            <!-- Search Results -->
            <div
                v-if="users.length"
                class="mt-4 p-3 rounded bg-gray-50 border cursor-pointer hover:bg-gray-200 transition"
                @click="selectUser(users[0])"
            >
                <div class="flex items-center gap-2 mb-2">
                    <img
                        v-if="users[0].avatar_url"
                        :src="`${backendBaseUrl}/static/${users[0].avatar_url}`"
                        alt="Avatar"
                        class="w-10 h-10 rounded-full object-cover border"
                    />
                    <img
                        v-else
                        src="/src/assets/default-avatar.jpeg"
                        alt="Default Avatar"
                        class="w-10 h-10 rounded-full object-cover border"
                    />
                    <span class="font-bold text-green-700">
                        {{ users[0].username }}
                    </span>
                </div>
                <div class="text-xs text-gray-600 mb-1">
                    Created: {{ new Date(users[0].created_at).toLocaleString() }}
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from "vue";
import axiosInstance from "@/axiosInstance";
import { showError } from "@/utils/toast";

const props = defineProps({
    backendBaseUrl: String,
});

const emit = defineEmits(['close', 'user-selected']);

const enteredUsername = ref("");
const users = ref([]);

// Search for user by username
const searchUser = async () => {
    if (!enteredUsername.value.trim()) return;

    try {
        const response = await axiosInstance.get(`/api/user/search?q=${enteredUsername.value}`);
        
        if (!response.data || Object.keys(response.data).length === 0) {
            users.value = [];
            showError("Username is invalid or not found.");
        } else {
            users.value = [response.data];
        }
    } catch (error) {
        console.error("Error fetching user:", error);
        showError("Failed to fetch user. Please try again.");
        users.value = [];
    }
};

// Select user from search results
const selectUser = (user) => {
    emit('user-selected', user);
};
</script> 