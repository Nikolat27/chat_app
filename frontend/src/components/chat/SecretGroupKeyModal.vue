<template>
    <div v-if="isVisible" class="fixed inset-0 bg-gray-500 bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 class="text-xl font-semibold mb-2">Secret Group Key</h3>
            <p v-if="groupName" class="text-sm text-gray-600 mb-4">{{ groupName }}</p>
            
            <!-- Show existing key -->
            <div v-if="existingKey" class="mb-4">
                <label class="block text-sm font-medium text-gray-700 mb-2">Group Secret Key:</label>
                <div class="flex items-center space-x-2">
                    <input
                        :value="existingKey"
                        type="text"
                        readonly
                        class="flex-1 border border-gray-300 rounded px-3 py-2 text-sm bg-gray-50"
                    />
                    <button
                        @click="copyToClipboard"
                        class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded text-sm cursor-pointer"
                    >
                        Copy
                    </button>
                </div>
                <p class="text-sm text-gray-500 mt-1">
                    Share this key with other group members so they can join the secret group.
                </p>
            </div>
            
            <!-- Enter key form -->
            <div v-else class="mb-4">
                <label class="block text-sm font-medium text-gray-700 mb-2">Enter Secret Key:</label>
                <input
                    v-model="keyInput"
                    type="text"
                    placeholder="Paste the secret key here..."
                    class="w-full border border-gray-300 rounded px-3 py-2 text-sm"
                />
                <p class="text-sm text-gray-500 mt-1">
                    Ask the group owner for the secret key to join this group.
                </p>
            </div>
            
            <!-- Error message -->
            <div v-if="error" class="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
                {{ error }}
            </div>
            
            <!-- Success message -->
            <div v-if="success" class="mb-4 p-3 bg-green-100 border border-green-400 text-green-700 rounded">
                {{ success }}
            </div>
            
            <!-- Action buttons -->
            <div class="flex justify-end space-x-2">
                <button
                    @click="closeModal"
                    class="px-4 py-2 text-gray-600 hover:text-gray-800 cursor-pointer"
                >
                    Cancel
                </button>
                <button
                    v-if="!existingKey"
                    @click="enterKey"
                    :disabled="!keyInput.trim()"
                    class="px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded disabled:opacity-50 cursor-pointer"
                >
                    Join Group
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useSecretGroupE2EE } from '../../composables/useSecretGroupE2EE';

const props = defineProps({
    isVisible: {
        type: Boolean,
        default: false
    },
    groupId: {
        type: String,
        required: true
    },
    groupName: {
        type: String,
        default: ''
    },
    isGroupOwner: {
        type: Boolean,
        default: false
    }
});

const emit = defineEmits(['close', 'key-entered']);

const { getGroupSecretKey, enterSecretKey, copySecretKeyToClipboard } = useSecretGroupE2EE();

const keyInput = ref('');
const error = ref('');
const success = ref('');
const existingKey = ref('');

// Check if user already has the key
const checkExistingKey = async () => {
    if (props.isGroupOwner && props.groupId) {
        existingKey.value = await getGroupSecretKey(props.groupId);
    }
};

// Copy key to clipboard
const copyToClipboard = async () => {
    try {
        if (!props.groupId) {
            throw new Error('No group ID available');
        }
        await copySecretKeyToClipboard(props.groupId);
        success.value = 'Secret key copied to clipboard!';
        setTimeout(() => {
            success.value = '';
        }, 3000);
    } catch (err) {
        error.value = 'Failed to copy key to clipboard';
    }
};

// Enter the secret key
const enterKey = async () => {
    try {
        if (!props.groupId) {
            throw new Error('No group ID available');
        }
        error.value = '';
        success.value = '';
        
        await enterSecretKey(props.groupId, keyInput.value.trim());
        success.value = 'Secret key entered successfully!';
        
        setTimeout(() => {
            emit('key-entered');
            closeModal();
        }, 1000);
    } catch (err) {
        error.value = err.message || 'Failed to enter secret key';
    }
};

// Close modal
const closeModal = () => {
    keyInput.value = '';
    error.value = '';
    success.value = '';
    emit('close');
};

// Watch for visibility changes
watch(() => props.isVisible, (newValue) => {
    if (newValue) {
        checkExistingKey();
    }
});
</script> 