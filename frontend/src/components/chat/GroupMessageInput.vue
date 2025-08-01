<template>
    <div class="p-4 border-t border-gray-200 bg-white flex gap-2">
        <input
            :value="modelValue"
            type="text"
            placeholder="Type a message..."
            class="flex-1 border border-gray-200 rounded-lg px-4 py-2 text-[16px] font-medium focus:outline-none focus:ring-1 focus:ring-gray-200 transition"
            @input="$emit('update:modelValue', $event.target.value)"
            @keyup.enter="handleSend"
        />
        <button
            v-if="!isSecretGroup"
            @click="handleImageUpload"
            class="bg-blue-500 hover:bg-blue-600 text-white font-semibold w-10 h-10 rounded-full shadow-sm transition cursor-pointer flex items-center justify-center"
            title="Upload Image"
        >
            <span class="material-icons text-lg">image</span>
        </button>
        <button
            @click="handleSend"
            class="bg-green-500 hover:bg-green-600 text-white font-semibold w-10 h-10 rounded-full shadow-sm transition cursor-pointer flex items-center justify-center"
        >
            <span class="material-icons text-lg">send</span>
        </button>
    </div>

    <!-- Image Preview Modal -->
    <div
        v-if="showImagePreview"
        class="fixed inset-0 z-50 flex items-center justify-center"
    >
        <!-- Backdrop -->
        <div
            class="absolute inset-0 bg-gray-900 bg-opacity-75 backdrop-blur-sm"
            @click="closeImagePreview"
        ></div>

        <!-- Modal -->
        <div
            class="relative bg-white rounded-2xl shadow-2xl max-w-2xl w-full mx-4 overflow-hidden"
        >
            <!-- Header -->
            <div
                class="flex items-center justify-between p-6 border-b border-gray-200 bg-gray-50"
            >
                <div class="flex items-center space-x-3">
                    <div
                        class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center"
                    >
                        <span class="material-icons text-blue-600"
                            >image</span
                        >
                    </div>
                    <div>
                        <h3 class="text-lg font-semibold text-gray-900">
                            Send Image
                        </h3>
                        <p class="text-sm text-gray-500">
                            Preview and send your image
                        </p>
                    </div>
                </div>
                <button
                    @click="closeImagePreview"
                    class="text-gray-400 hover:text-gray-600 transition-colors"
                >
                    <span class="material-icons">close</span>
                </button>
            </div>

            <!-- Image Preview -->
            <div class="p-6">
                <div class="mb-6">
                    <div class="relative bg-gray-100 rounded-xl overflow-hidden">
                        <img
                            :src="imagePreviewUrl"
                            alt="Image Preview"
                            class="w-full h-64 object-contain"
                        />
                    </div>
                </div>

                <!-- Image Info -->
                <div class="mb-6 p-4 bg-gray-50 rounded-lg">
                    <div class="flex items-center justify-between text-sm text-gray-600">
                        <span class="flex items-center">
                            <span class="material-icons text-xs mr-1">info</span>
                            File: {{ selectedImage?.name }}
                        </span>
                        <span class="flex items-center">
                            <span class="material-icons text-xs mr-1">storage</span>
                            {{ selectedImage ? (selectedImage.size / 1024 / 1024).toFixed(2) : '0' }} MB
                        </span>
                    </div>
                </div>

                <!-- Action Buttons -->
                <div class="flex items-center justify-end space-x-3">
                    <button
                        @click="closeImagePreview"
                        class="px-6 py-3 text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-lg font-medium transition-colors cursor-pointer"
                    >
                        Cancel
                    </button>
                    <button
                        @click="handleSendImage"
                        class="px-6 py-3 text-white bg-blue-500 hover:bg-blue-600 rounded-lg font-medium transition-colors cursor-pointer flex items-center space-x-2"
                    >
                        <span class="material-icons text-sm">send</span>
                        Send Image
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from "vue";
import { showError } from "../../utils/toast";

const props = defineProps({
    modelValue: {
        type: String,
        default: "",
    },
    isSecretGroup: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits(["update:modelValue", "send", "image-upload"]);

// Reactive data for image preview
const showImagePreview = ref(false);
const selectedImage = ref(null);
const imagePreviewUrl = ref('');

const handleSend = () => {
    emit("send");
};

const handleImageUpload = () => {
    // Create a hidden file input
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = 'image/*';
    fileInput.style.display = 'none';
    
    fileInput.addEventListener('change', (event) => {
        const file = event.target.files[0];
        if (file) {
            // Check file size (20MB = 20 * 1024 * 1024 bytes)
            const maxSize = 20 * 1024 * 1024; // 20MB in bytes
            if (file.size > maxSize) {
                showError('Image size must be less than 20MB. Please choose a smaller image.');
                return;
            }
            
            selectedImage.value = file;
            imagePreviewUrl.value = URL.createObjectURL(file);
            showImagePreview.value = true;
        }
    });
    
    // Trigger file selection
    fileInput.click();
};

const handleSendImage = () => {
    if (selectedImage.value) {
        emit('image-upload', selectedImage.value);
        closeImagePreview();
    }
};

const closeImagePreview = () => {
    showImagePreview.value = false;
    selectedImage.value = null;
    if (imagePreviewUrl.value) {
        URL.revokeObjectURL(imagePreviewUrl.value);
        imagePreviewUrl.value = '';
    }
};
</script> 