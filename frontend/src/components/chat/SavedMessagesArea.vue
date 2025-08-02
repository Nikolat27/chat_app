<template>
    <div class="flex flex-col h-full">
        <!-- Header -->
        <div class="p-4 border-b border-gray-200 bg-white">
            <div class="flex items-center justify-between">
                <h2 class="text-lg font-semibold text-gray-800">Saved Messages</h2>
                <button
                    @click="showCreateModal = true"
                    class="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded-lg text-sm font-medium transition"
                >
                    <span class="material-icons text-sm mr-1">add</span>
                    New Message
                </button>
            </div>
        </div>

        <!-- Messages List -->
        <div class="flex-1 overflow-y-auto p-4 space-y-3">
            <div v-if="isLoading" class="flex justify-center py-8">
                <div class="loading-spinner"></div>
            </div>
            
            <div v-else-if="savedMessages.length === 0" class="text-center py-8 text-gray-500">
                <span class="material-icons text-4xl mb-2 text-gray-300">bookmark_border</span>
                <p>No saved messages yet</p>
                <p class="text-sm">Click "New Message" to create your first saved message</p>
            </div>
            
            <div v-else class="space-y-3">
                <div
                    v-for="message in savedMessages"
                    :key="message?.id || Math.random()"
                    v-if="message && message.id && message.content"
                    class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-sm transition"
                >
                    <div class="flex items-start justify-between">
                        <div class="flex-1">
                            <div class="flex items-center gap-2 mb-2">
                                <span class="material-icons text-sm text-gray-400">bookmark</span>
                                <span class="text-xs text-gray-500">
                                    {{ formatDate(message?.created_at) }}
                                </span>
                                <span v-if="message?.updated_at && message?.updated_at !== message?.created_at" 
                                      class="text-xs text-gray-400">
                                    (edited)
                                </span>
                            </div>
                            
                            <div v-if="!message?.isEditing" class="text-gray-800 whitespace-pre-wrap">
                                {{ message?.content || '' }}
                            </div>
                            
                            <div v-else class="space-y-2">
                                <textarea
                                    v-model="message?.editContent"
                                    class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-green-500"
                                    rows="3"
                                    placeholder="Edit your message..."
                                ></textarea>
                                <div class="flex gap-2">
                                    <button
                                        @click="saveEdit(message)"
                                        class="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded text-sm font-medium transition"
                                    >
                                        Save
                                    </button>
                                    <button
                                        @click="cancelEdit(message)"
                                        class="bg-gray-300 hover:bg-gray-400 text-gray-700 px-3 py-1 rounded text-sm font-medium transition"
                                    >
                                        Cancel
                                    </button>
                                </div>
                            </div>
                        </div>
                        
                        <div class="flex gap-1 ml-2">
                            <button
                                v-if="!message?.isEditing"
                                @click="startEdit(message)"
                                class="text-gray-400 hover:text-blue-500 p-1 rounded transition"
                                title="Edit message"
                            >
                                <span class="material-icons text-sm">edit</span>
                            </button>
                            <button
                                @click="deleteMessage(message?.id)"
                                class="text-gray-400 hover:text-red-500 p-1 rounded transition"
                                title="Delete message"
                            >
                                <span class="material-icons text-sm">delete</span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Create Message Modal -->
        <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
                <h3 class="text-lg font-semibold mb-4">Create New Saved Message</h3>
                
                <textarea
                    v-model="newMessageContent"
                    class="w-full border border-gray-300 rounded-lg px-3 py-2 mb-4 focus:outline-none focus:ring-1 focus:ring-green-500"
                    rows="4"
                    placeholder="Enter your message..."
                ></textarea>
                
                <div class="flex gap-2 justify-end">
                    <button
                        @click="showCreateModal = false"
                        class="bg-gray-300 hover:bg-gray-400 text-gray-700 px-4 py-2 rounded font-medium transition"
                    >
                        Cancel
                    </button>
                    <button
                        @click="createMessage"
                        :disabled="!newMessageContent.trim() || isLoading"
                        class="bg-green-500 hover:bg-green-600 disabled:bg-gray-300 text-white px-4 py-2 rounded font-medium transition"
                    >
                        {{ isLoading ? 'Creating...' : 'Create' }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useSaveMessage } from '../../composables/useSaveMessage';

const {
    isLoading,
    error,
    savedMessages,
    getSavedMessages,
    createSavedMessage,
    updateSavedMessage,
    deleteSavedMessage
} = useSaveMessage();

const showCreateModal = ref(false);
const newMessageContent = ref('');

// Load saved messages on component mount
onMounted(async () => {
    await getSavedMessages();
});

// Create new message
const createMessage = async () => {
    if (!newMessageContent.value.trim()) return;
    
    const result = await createSavedMessage(newMessageContent.value.trim());
    if (result) {
        newMessageContent.value = '';
        showCreateModal.value = false;
    }
};

// Start editing a message
const startEdit = (message) => {
    if (!message || !message.id) return;
    message.isEditing = true;
    message.editContent = message.content || '';
};

// Cancel editing
const cancelEdit = (message) => {
    message.isEditing = false;
    message.editContent = '';
};

// Save edited message
const saveEdit = async (message) => {
    if (!message || !message.id || !message.editContent.trim()) return;
    
    const result = await updateSavedMessage(message.id, message.editContent.trim());
    if (result) {
        message.isEditing = false;
        message.editContent = '';
    }
};

// Delete message
const deleteMessage = async (messageId) => {
    if (!messageId) return;
    if (confirm('Are you sure you want to delete this saved message?')) {
        await deleteSavedMessage(messageId);
    }
};

// Format date
const formatDate = (dateString) => {
    if (!dateString) return '';
    
    const date = new Date(dateString);
    const now = new Date();
    const diffInHours = (now - date) / (1000 * 60 * 60);
    
    if (diffInHours < 24) {
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    } else if (diffInHours < 48) {
        return 'Yesterday';
    } else {
        return date.toLocaleDateString();
    }
};
</script>

<style scoped>
.loading-spinner {
    width: 24px;
    height: 24px;
    border: 2px solid #f3f4f6;
    border-top: 2px solid #10b981;
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}
</style> 