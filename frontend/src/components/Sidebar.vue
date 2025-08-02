<template>
    <aside
        class="w-20 bg-white border-r border-gray-200 flex flex-col items-center py-4"
    >
        <!-- Navigation Buttons -->
        <button
            class="mb-5 w-12 h-12 rounded-full hover:bg-green-200 transition cursor-pointer flex items-center justify-center"
            @click="$emit('changeTab', 'chats')"
            :class="{ 'bg-white shadow-md': activeTab === 'chats' }"
        >
            <span class="material-icons" :class="{ 'text-green-600': activeTab === 'chats', 'text-gray-700': activeTab !== 'chats' }">chat</span>
        </button>
        <button
            class="mb-5 w-12 h-12 rounded-full hover:bg-green-200 transition cursor-pointer flex items-center justify-center"
            @click="$emit('changeTab', 'groups')"
            :class="{ 'bg-white shadow-md': activeTab === 'groups' }"
        >
            <span class="material-icons" :class="{ 'text-green-600': activeTab === 'groups', 'text-gray-700': activeTab !== 'groups' }">group</span>
        </button>

        <button
            class="mb-5 w-12 h-12 rounded-full hover:bg-green-200 transition cursor-pointer flex items-center justify-center"
            @click="$emit('changeTab', 'saved-messages')"
            :class="{ 'bg-white shadow-md': activeTab === 'saved-messages' }"
        >
            <span class="material-icons" :class="{ 'text-green-600': activeTab === 'saved-messages', 'text-gray-700': activeTab !== 'saved-messages' }">bookmark</span>
        </button>

        <button
            class="mb-5 w-12 h-12 rounded-full hover:bg-green-200 transition cursor-pointer flex items-center justify-center"
            @click="$emit('changeTab', 'settings')"
            :class="{ 'bg-white shadow-md': activeTab === 'settings' }"
        >
            <span class="material-icons" :class="{ 'text-green-600': activeTab === 'settings', 'text-gray-700': activeTab !== 'settings' }">settings</span>
        </button>

        <!-- Logout Button -->
        <button
            v-if="userStore.username"
            class="mt-auto mb-4 w-12 h-12 rounded-full hover:bg-green-200 transition cursor-pointer flex items-center justify-center"
            @click="$emit('logout')"
        >
            <span class="material-icons text-red-500">logout</span>
        </button>

        <!-- User Avatar -->
        <img
            :src="
                userStore.avatar_url
                    ? `${backendBaseUrl}/static/${userStore.avatar_url}`
                    : defaultAvatar
            "
            alt="User Avatar"
            class="w-12 h-12 rounded-full object-cover border border-gray-200 shadow-sm mb-2 select-none pointer-events-none"
        />

        <!-- Username -->
        <div
            class="text-xs text-gray-700 font-semibold text-center px-1 break-words"
        >
            {{ userStore.username }}
        </div>
    </aside>
</template>

<script setup>
import { useUserStore } from "../stores/users";
import { defineProps } from "vue";
import defaultAvatar from "../assets/default-avatar.jpg";
const props = defineProps({ activeTab: String });
const userStore = useUserStore();
const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style>
