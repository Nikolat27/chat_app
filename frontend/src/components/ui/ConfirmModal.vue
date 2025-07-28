<template>
    <div
        v-if="isVisible"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm"
        @click.self="$emit('close')"
    >
        <div class="bg-white rounded-2xl shadow-2xl border border-gray-100 max-w-md w-full mx-4 p-6">
            <!-- Header -->
            <div class="flex items-center gap-3 mb-4">
                <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
                    <span class="material-icons text-red-600 text-xl">warning</span>
                </div>
                <div>
                    <h3 class="text-lg font-bold text-gray-800">{{ title }}</h3>
                    <p class="text-sm text-gray-600">{{ subtitle }}</p>
                </div>
            </div>

            <!-- Content -->
            <div class="mb-6">
                <p class="text-gray-700 leading-relaxed">{{ message }}</p>
            </div>

            <!-- Actions -->
            <div class="flex gap-3">
                <button
                    @click="$emit('close')"
                    class="flex-1 px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg font-medium transition-colors duration-200 cursor-pointer"
                >
                    Cancel
                </button>
                <button
                    @click="handleConfirm"
                    :disabled="isLoading"
                    class="flex-1 px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg font-medium transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 cursor-pointer"
                >
                    <div v-if="isLoading" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                    <span v-else>{{ confirmText }}</span>
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
const props = defineProps({
    isVisible: {
        type: Boolean,
        default: false,
    },
    title: {
        type: String,
        default: "Confirm Action",
    },
    subtitle: {
        type: String,
        default: "Are you sure?",
    },
    message: {
        type: String,
        required: true,
    },
    confirmText: {
        type: String,
        default: "Confirm",
    },
    isLoading: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits(["close", "confirm"]);

const handleConfirm = () => {
    emit("confirm");
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style> 