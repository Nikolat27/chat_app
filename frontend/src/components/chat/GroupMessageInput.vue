<template>
    <!-- Show secret key input if it's a secret group and user hasn't entered the key -->
    <div v-if="isSecretGroup && !hasSecretKey" class="p-4 border-t border-gray-200 bg-white">
        <div class="w-full p-4 bg-purple-50 border border-purple-200 rounded-lg">
            <div class="flex items-center gap-2 text-purple-700 mb-3">
                <span class="material-icons text-sm">lock</span>
                <span class="text-sm font-medium">Secret Key Required</span>
            </div>
            <p class="text-xs text-purple-600 mb-4">
                You need to enter the secret key to send messages in this secret group.
            </p>
            <div class="flex gap-2">
                <input
                    v-model="secretKeyInput"
                    type="password"
                    placeholder="Enter secret key..."
                    class="flex-1 border border-purple-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-purple-300 transition"
                    @keyup.enter="handleEnterSecretKey"
                />
                <button
                    @click="handleEnterSecretKey"
                    :disabled="!secretKeyInput.trim() || isEnteringKey"
                    class="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg font-medium transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed text-sm"
                >
                    <span v-if="isEnteringKey" class="flex items-center gap-1">
                        <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Entering...
                    </span>
                    <span v-else class="flex items-center gap-1">
                        <span class="material-icons text-sm">key</span>
                        Enter Key
                    </span>
                </button>
                <button
                    v-if="isGroupOwner"
                    @click="handleOpenSecretKeyModal"
                    class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors cursor-pointer text-sm"
                    title="Copy Secret Key"
                >
                    <span class="material-icons text-sm">content_copy</span>
                </button>
            </div>
        </div>
    </div>
    
    <!-- Show message input only if it's not a secret group or user has entered the key -->
    <div v-else class="p-4 border-t border-gray-200 bg-white flex gap-2">
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
import { ref, computed, onMounted, watch } from "vue";
import { showError, showMessage } from "../../utils/toast";
import { useSecretGroupE2EE } from "../../composables/useSecretGroupE2EE";

const props = defineProps({
    modelValue: {
        type: String,
        default: "",
    },
    isSecretGroup: {
        type: Boolean,
        default: false,
    },
    groupId: {
        type: String,
        default: "",
    },
    isGroupOwner: {
        type: Boolean,
        default: false,
    },
    keyStatus: {
        type: String,
        default: 'not-entered', // 'not-entered', 'entering', 'entered'
    }
});

const emit = defineEmits(["update:modelValue", "send", "image-upload", "open-secret-key-modal"]);

// Secret group key management
const { hasGroupSecretKey, enterSecretKey } = useSecretGroupE2EE();
const hasSecretKey = ref(false);
const secretKeyInput = ref("");
const isEnteringKey = ref(false);

// Reactive data for image preview
const showImagePreview = ref(false);
const selectedImage = ref(null);
const imagePreviewUrl = ref('');

// Check if user has secret key on mount and when groupId changes
const checkSecretKey = async () => {
    if (props.isSecretGroup && props.groupId) {
        hasSecretKey.value = await hasGroupSecretKey(props.groupId);
    }
};

onMounted(checkSecretKey);

// Watch for groupId changes
watch(() => props.groupId, checkSecretKey);

// Watch for external key status updates (when key is entered via modal)
watch(() => props.keyStatus, (newStatus) => {
    if (newStatus === 'entered') {
        hasSecretKey.value = true;
    }
});

const handleSend = () => {
    // Only allow sending if not a secret group or user has entered the key
    if (props.isSecretGroup && !hasSecretKey.value) {
        showError("Please enter the secret key first");
        return;
    }
    emit("send");
};

const handleEnterSecretKey = async () => {
    if (!secretKeyInput.value.trim()) {
        showError("Please enter a secret key");
        return;
    }

    isEnteringKey.value = true;
    
    try {
        await enterSecretKey(props.groupId, secretKeyInput.value.trim());
        hasSecretKey.value = true;
        secretKeyInput.value = "";
        showMessage("Secret key entered successfully! You can now send messages.");
    } catch (error) {
        console.error("Failed to enter secret key:", error);
        showError(error.message || "Invalid secret key. Please try again.");
    } finally {
        isEnteringKey.value = false;
    }
};

const handleImageUpload = () => {
    // Only allow image upload if not a secret group
    if (props.isSecretGroup) {
        showError("Image upload is not available for secret groups");
        return;
    }
    
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

const handleOpenSecretKeyModal = () => {
    emit('open-secret-key-modal');
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