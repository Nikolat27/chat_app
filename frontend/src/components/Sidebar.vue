<template>
    <aside class="w-20 bg-white border-r flex flex-col items-center py-4">
        <button
            class="mb-6 p-2 rounded hover:bg-green-100 cursor-pointer"
            @click="$emit('changeTab', 'chats')"
            :class="{ 'bg-green-200': activeTab === 'chats' }"
        >
            <span class="material-icons">chat</span>
        </button>
        <button
            class="mb-6 p-2 rounded hover:bg-green-100 cursor-pointer"
            @click="$emit('changeTab', 'groups')"
            :class="{ 'bg-green-200': activeTab === 'groups' }"
        >
            <span class="material-icons">group</span>
        </button>
        <button
            class="mb-6 p-2 rounded hover:bg-green-100 cursor-pointer"
            @click="$emit('changeTab', 'settings')"
            :class="{ 'bg-green-200': activeTab === 'settings' }"
        >
            <span class="material-icons">settings</span>
        </button>
        <button
            v-if="userStore.token"
            class="mt-auto p-2 mb-4 rounded hover:bg-green-100 cursor-pointer"
            @click="$emit('logout')"
        >
            <span class="material-icons">logout</span>
        </button>
        <img
            :src="userStore.avatar_url ? `${backendBaseUrl}/static/${userStore.avatar_url}` : defaultAvatar"
            alt="User Avatar"
            class="w-12 h-12 rounded-full object-cover border mb-2"
        />
        <div class="text-xs text-gray-700 font-semibold mb-4">{{ userStore.username }}</div>
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
