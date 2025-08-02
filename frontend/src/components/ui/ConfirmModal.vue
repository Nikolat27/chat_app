<template>
    <transition name="modal-fade">
        <div
            v-if="isVisible"
            class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-80 backdrop-blur-md"
            @click.self="$emit('close')"
        >
            <div class="bg-white rounded-2xl shadow-2xl border-2 border-red-400 max-w-md w-full mx-4 p-6 animate-modal-pop relative">
                <!-- Header -->
                <div class="flex items-center gap-3 mb-4">
                    <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center border-2 border-red-300">
                        <span class="material-icons text-red-600 text-2xl">warning</span>
                    </div>
                    <div>
                        <h3 class="text-lg font-bold text-red-700">{{ title }}</h3>
                        <p class="text-sm text-gray-600">{{ subtitle }}</p>
                    </div>
                </div>

                <!-- Content -->
                <div class="mb-6">
                    <p class="text-gray-700 leading-relaxed font-medium">{{ message }}</p>
                </div>

                <!-- Actions -->
                <div class="flex gap-3">
                    <button
                        @click="$emit('close')"
                        class="flex-1 px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg font-medium transition-colors duration-200 cursor-pointer"
                        aria-label="Cancel deletion"
                    >
                        Cancel
                    </button>
                    <button
                        @click="handleConfirm"
                        :disabled="isLoading"
                        class="flex-1 px-4 py-2 bg-gradient-to-r from-red-500 to-red-700 hover:from-red-600 hover:to-red-800 text-white rounded-lg font-bold text-base transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 cursor-pointer shadow-md focus:outline-none focus:ring-2 focus:ring-red-400"
                        aria-label="Confirm deletion"
                    >
                        <div v-if="isLoading" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                        <span v-else>{{ confirmText }}</span>
                    </button>
                </div>
            </div>
        </div>
    </transition>
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

.modal-fade-enter-active, .modal-fade-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-fade-enter-from, .modal-fade-leave-to {
  opacity: 0;
  transform: scale(0.96);
}

.animate-modal-pop {
  animation: modal-pop 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
@keyframes modal-pop {
  0% {
    opacity: 0;
    transform: scale(0.96);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}
</style> 