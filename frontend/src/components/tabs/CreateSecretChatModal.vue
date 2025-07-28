<template>
    <div
        class="fixed inset-0 bg-gradient-to-br from-purple-100 via-white to-purple-200 bg-opacity-90 backdrop-blur-sm flex items-center justify-center z-50 p-4"
        @click.self="$emit('close')"
    >
        <div
            class="bg-white rounded-3xl shadow-2xl p-8 w-full max-w-md relative font-sans transform transition-all duration-300 scale-100"
        >
            <!-- Close Button -->
            <button
                class="absolute top-4 right-4 text-gray-400 hover:text-purple-700 hover:bg-purple-50 p-2 rounded-full transition-all duration-200 cursor-pointer"
                @click="$emit('close')"
                aria-label="Close"
            >
                <span class="material-icons text-xl">close</span>
            </button>

            <!-- Header -->
            <div class="mb-8 text-center">
                <div class="mb-4">
                    <span class="material-icons text-4xl text-purple-500 mb-3">lock</span>
                </div>
                <h3
                    class="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-purple-600 to-pink-500 mb-3 tracking-tight"
                >
                    Start Secret Chat
                </h3>
                <p class="text-sm text-gray-600 leading-relaxed">
                    Create an end-to-end encrypted conversation with enhanced privacy.
                </p>
            </div>

            <!-- Search Form -->
            <div class="mb-6">
                <label class="block mb-3 text-purple-700 font-semibold text-sm">
                    <span class="material-icons align-middle mr-1 text-base">person_search</span>
                    Enter Username
                </label>
                <div class="relative">
                    <input
                        v-model="enteredUsername"
                        type="text"
                        placeholder="Type username to search..."
                        class="w-full border-2 border-purple-200 rounded-xl px-4 py-3 mb-4 focus:outline-none focus:ring-2 focus:ring-purple-300 focus:border-purple-400 transition-all duration-200 text-gray-700"
                        @keyup.enter="searchUser"
                        :disabled="isSearching"
                    />
                    <div v-if="isSearching" class="absolute right-3 top-1/2 transform -translate-y-1/2">
                        <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-purple-500"></div>
                    </div>
                </div>
                <button
                    class="bg-gradient-to-r from-purple-500 to-pink-500 text-white px-6 py-3 rounded-xl w-full font-semibold shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-200 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                    @click="searchUser"
                    :disabled="!enteredUsername.trim() || isSearching"
                >
                    <span v-if="!isSearching" class="material-icons align-middle mr-2 text-base">search</span>
                    <span v-else class="material-icons align-middle mr-2 text-base animate-spin">hourglass_empty</span>
                    {{ isSearching ? 'Searching...' : 'Search User' }}
                </button>
                
                <!-- Warning Messages -->
                <div
                    v-if="selfSearchWarning"
                    class="mt-4 p-3 bg-red-50 border border-red-200 rounded-xl text-center text-sm font-medium text-red-700 flex items-center justify-center"
                >
                    <span class="material-icons mr-2 text-red-500">warning</span>
                    You can't start a secret chat with yourself!
                </div>
                
                <div
                    v-if="searchError"
                    class="mt-4 p-3 bg-orange-50 border border-orange-200 rounded-xl text-center text-sm font-medium text-orange-700 flex items-center justify-center"
                >
                    <span class="material-icons mr-2 text-orange-500">info</span>
                    {{ searchError }}
                </div>
            </div>

            <!-- Search Results -->
            <div
                v-if="user && !selfSearchWarning && !searchError"
                class="mt-6 p-5 rounded-2xl bg-gradient-to-r from-purple-50 to-pink-50 border-2 border-purple-200 cursor-pointer hover:bg-gradient-to-r hover:from-purple-100 hover:to-pink-100 hover:scale-[1.02] hover:shadow-lg transition-all duration-300 flex flex-col gap-3 items-start shadow-md"
                @click="selectUser(user)"
            >
                <div class="flex items-center gap-4 mb-2 w-full">
                    <div class="relative">
                        <img
                            v-if="user.avatar_url"
                            :src="
                                user.avatar_url.startsWith('http')
                                    ? user.avatar_url
                                    : `${backendBaseUrl}/static/${user.avatar_url}`
                            "
                            alt="Avatar"
                            class="w-14 h-14 rounded-full object-cover border-2 border-purple-300 shadow-lg select-none pointer-events-none"
                        />
                        <img
                            v-else
                            src="/src/assets/default-avatar.jpg"
                            alt="Default Avatar"
                            class="w-14 h-14 rounded-full object-cover border-2 border-purple-300 shadow-lg select-none pointer-events-none"
                        />
                        <div class="absolute -bottom-1 -right-1 w-6 h-6 bg-purple-500 rounded-full flex items-center justify-center">
                            <span class="material-icons text-white text-sm">lock</span>
                        </div>
                    </div>
                    <div class="flex-1">
                        <span class="font-bold text-purple-800 text-lg block">
                            {{ user.username }}
                        </span>
                        <span class="text-xs text-purple-600 font-medium">
                            Available for secret chat
                        </span>
                    </div>
                </div>
                
                <div class="text-xs text-gray-500 mb-2 w-full">
                    <div class="flex items-center gap-2">
                        <span class="material-icons text-xs">schedule</span>
                        Member since: {{ user.created_at ? new Date(user.created_at).toLocaleDateString() : 'Unknown' }}
                    </div>
                </div>
                
                <div class="w-full bg-purple-100 rounded-lg p-3 text-xs text-purple-700">
                    <div class="flex items-center gap-2 mb-1">
                        <span class="material-icons text-xs">security</span>
                        <span class="font-semibold">End-to-End Encrypted</span>
                    </div>
                    <p class="text-purple-600 leading-relaxed">
                        Messages will be encrypted with your unique key pair for maximum privacy.
                    </p>
                </div>
                
                <button class="w-full bg-gradient-to-r from-purple-600 to-pink-600 text-white py-2 px-4 rounded-lg font-semibold hover:from-purple-700 hover:to-pink-700 transition-all duration-200 flex items-center justify-center gap-2">
                    <span class="material-icons text-sm">lock</span>
                    Start Secret Chat
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from "vue";
import axiosInstance from "../../axiosInstance";
import { showError } from "@/utils/toast";

const props = defineProps({
    backendBaseUrl: String,
});

const emit = defineEmits(["close", "user-selected"]);

const enteredUsername = ref("");
const user = ref(null);
const selfSearchWarning = ref(false);
const searchError = ref("");
const isSearching = ref(false);

// Get current user's username from localStorage or a global store
let currentUsername = null;
try {
    const userObj = JSON.parse(localStorage.getItem("user"));
    currentUsername = userObj?.username || null;
} catch (e) {
    currentUsername = null;
}

const searchUser = async () => {
    selfSearchWarning.value = false;
    searchError.value = "";
    user.value = null;
    isSearching.value = true;
    
    if (!enteredUsername.value.trim()) {
        isSearching.value = false;
        return;
    }

    // Check if searching for self
    if (
        currentUsername &&
        enteredUsername.value.trim().toLowerCase() ===
            currentUsername.toLowerCase()
    ) {
        selfSearchWarning.value = true;
        isSearching.value = false;
        return;
    }

    try {
        const response = await axiosInstance.get(
            `/api/user/search?q=${enteredUsername.value}`
        );
        if (!response.data || Object.keys(response.data).length === 0) {
            user.value = null;
            searchError.value = "Username not found. Please check the spelling and try again.";
        } else {
            user.value = response.data;
        }
    } catch (error) {
        user.value = null;
        if (error?.response?.data?.Detail) {
            searchError.value = error.response.data.Detail;
        } else {
            searchError.value = "Failed to search for user. Please check your connection and try again.";
        }
    } finally {
        isSearching.value = false;
    }
};

const selectUser = (u) => {
    emit("user-selected", u);
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style>
