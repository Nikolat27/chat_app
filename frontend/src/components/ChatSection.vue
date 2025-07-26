<template>
    <main class="flex-1 p-4 flex flex-col">
        <div class="flex-1 bg-white rounded-lg shadow p-4">
            <h2 class="text-lg font-bold text-green-600 mb-4">Chat Window</h2>
            <div v-if="loading" class="flex items-center justify-center h-full">
                <span class="material-icons animate-spin mr-2">autorenew</span>
                <span class="text-green-600 font-semibold">Loading...</span>
            </div>
            <div v-else>
                <!-- Chat messages will go here -->
                <div class="text-gray-500">
                    Select a chat to start messaging.
                </div>
            </div>
        </div>
    </main>
</template>

<script setup>
import { ref } from "vue";
const loading = ref(true);

// Simulate loading for demo
setTimeout(() => {
    loading.value = false;
}, 1500);
</script>

<style scoped>
.material-icons.animate-spin {
    animation: spin 1s linear infinite;
}
@keyframes spin {
    100% {
        transform: rotate(360deg);
    }
}
</style>
