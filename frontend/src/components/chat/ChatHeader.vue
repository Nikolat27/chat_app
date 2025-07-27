<template>
    <div
        class="flex items-center gap-4 px-5 py-3 border-gray-300 border-b bg-gray-50 shadow-sm"
    >
        <img
            v-if="user.avatar_url"
            :src="`${backendBaseUrl}/static/${user.avatar_url}`"
            class="w-11 h-11 rounded-full object-cover border border-gray-300 shadow-sm"
            alt="Avatar"
        />
        <img
            v-else
            src="/src/assets/default-avatar.jpeg"
            class="w-11 h-11 rounded-full object-cover border border-gray-300 shadow-sm"
            alt="Default Avatar"
        />
        <div class="flex flex-col">
            <span class="font-semibold text-gray-800 text-base leading-tight">
                {{ user.username }}
            </span>
            <span class="text-xs text-green-600">Online</span>
        </div>
    </div>
</template>

<script setup>
defineProps({
    user: {
        type: Object,
        required: true,
    },
    backendBaseUrl: {
        type: String,
        required: true,
    },
});
</script>
