<template>
    <div class="flex items-center gap-3 p-4 border-b bg-white">
        <img
            v-if="user.avatar_url"
            :src="`${backendBaseUrl}/static/${user.avatar_url}`"
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
                {{ user.username }}
            </div>
            <div class="text-xs text-gray-500">Online</div>
        </div>
    </div>
</template>

<script setup>
defineProps({
    user: {
        type: Object,
        required: true
    },
    backendBaseUrl: {
        type: String,
        required: true
    }
});
</script> 