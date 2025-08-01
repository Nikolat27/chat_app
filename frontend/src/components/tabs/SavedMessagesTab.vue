<template>
    <div class="bg-gray-200 min-h-screen p-6">
        <div class="max-w-6xl mx-auto">
            <!-- Header -->
            <div class="bg-white rounded-xl shadow-md p-6 mb-6">
                <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-4">
                        <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                            <span class="material-icons text-blue-600 text-2xl">bookmark</span>
                        </div>
                        <div>
                            <h2 class="text-2xl font-bold text-gray-800">Saved Messages</h2>
                            <p class="text-gray-600">Manage your important messages and notes</p>
                        </div>
                    </div>
                    <div class="text-right">
                        <div class="text-3xl font-bold text-blue-600">{{ savedMessages.length }}</div>
                        <div class="text-sm text-gray-500">Total Messages</div>
                    </div>
                </div>
            </div>

            <!-- Action Buttons -->
            <div class="bg-white rounded-xl shadow-md p-6 mb-6">
                <div class="flex flex-wrap gap-4">
                    <button
                        class="bg-blue-500 hover:bg-blue-600 text-white font-medium px-6 py-3 rounded-lg transition-colors flex items-center space-x-2 cursor-pointer"
                        @click="showCreateSavedMessage = true"
                    >
                        <span class="material-icons">add</span>
                        <span>Create Message</span>
                    </button>
                    <button
                        class="bg-green-500 hover:bg-green-600 text-white font-medium px-6 py-3 rounded-lg transition-colors flex items-center space-x-2 cursor-pointer"
                        @click="refreshSavedMessages"
                        :disabled="isLoading"
                    >
                        <span class="material-icons">refresh</span>
                        <span>Refresh</span>
                    </button>
                    <button
                        class="bg-purple-500 hover:bg-purple-600 text-white font-medium px-6 py-3 rounded-lg transition-colors flex items-center space-x-2 cursor-pointer"
                        @click="exportMessages"
                    >
                        <span class="material-icons">download</span>
                        <span>Export</span>
                    </button>
                </div>
            </div>

            <!-- Loading State -->
            <div v-if="isLoading" class="flex items-center justify-center py-12">
                <div class="bg-white rounded-xl shadow-md p-8 text-center">
                    <div class="animate-spin rounded-full h-8 w-8 border-2 border-blue-500 border-t-transparent mx-auto mb-4"></div>
                    <p class="text-gray-600">Loading saved messages...</p>
                </div>
            </div>

            <!-- Empty State -->
            <div v-else-if="!savedMessages || savedMessages.length === 0" class="bg-white rounded-xl shadow-md p-8 text-center">
                <span class="material-icons text-6xl text-gray-300 mb-4">bookmark_border</span>
                <h3 class="text-xl font-semibold text-gray-600 mb-2">No Saved Messages</h3>
                <p class="text-gray-500 mb-6">Start saving important messages for quick access.</p>
                <button
                    @click="showCreateSavedMessage = true"
                    class="bg-blue-500 hover:bg-blue-600 text-white font-medium px-6 py-3 rounded-lg transition-colors flex items-center space-x-2 mx-auto cursor-pointer"
                >
                    <span class="material-icons">add</span>
                    <span>Create Message</span>
                </button>
            </div>

            <!-- Saved Messages List -->
            <div v-else>
                <!-- Stats -->
                <div class="bg-white rounded-xl shadow-md p-6 mb-6">
                    <div class="flex items-center justify-between">
                        <div class="flex space-x-8">
                            <div class="text-center">
                                <div class="text-2xl font-bold text-blue-600">{{ savedMessages.length }}</div>
                                <div class="text-sm text-gray-500">Total</div>
                            </div>
                            <div class="text-center">
                                <div class="text-2xl font-bold text-green-600">{{ getRecentCount() }}</div>
                                <div class="text-sm text-gray-500">Recent</div>
                            </div>
                            <div class="text-center">
                                <div class="text-2xl font-bold text-purple-600">{{ getCategoryCount() }}</div>
                                <div class="text-sm text-gray-500">Categories</div>
                            </div>
                        </div>
                        <div class="flex items-center space-x-2">
                            <span class="text-sm text-gray-500">Sort:</span>
                            <select 
                                v-model="sortBy" 
                                @change="sortMessages"
                                class="text-sm border border-gray-300 rounded px-3 py-1"
                            >
                                <option value="date">Date</option>
                                <option value="title">Title</option>
                                <option value="category">Category</option>
                            </select>
                        </div>
                    </div>
                </div>

                <!-- Messages Grid -->
                <div class="space-y-6">
                    <div
                        v-for="message in savedMessages"
                        :key="message.id"
                        class="bg-white rounded-xl shadow-md p-6 border border-gray-200 hover:shadow-lg transition-shadow"
                    >
                        <!-- Category Badge -->
                        <div v-if="message.category" class="mb-3">
                            <span class="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-purple-100 text-purple-800">
                                <span class="material-icons text-xs mr-1">category</span>
                                {{ message.category }}
                            </span>
                        </div>
                                        <!-- Message Header -->
                        <div class="flex items-center space-x-3 mb-4">
                            <div class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                                <span class="material-icons text-blue-600">bookmark</span>
                            </div>
                            <div class="flex-1">
                                <h4 class="font-semibold text-gray-900">
                                    {{ message.title || 'Untitled Message' }}
                                </h4>
                                <p class="text-sm text-gray-500">
                                    {{ formatDate(message.created_at) }}
                                </p>
                            </div>
                        </div>

                                        <!-- Message Content -->
                        <div class="mb-4">
                            <div class="bg-gray-50 rounded-lg p-4">
                                <p class="text-gray-700 whitespace-pre-wrap">
                                    {{ message.content }}
                                </p>
                            </div>
                        </div>

                        <!-- Action Buttons -->
                        <div class="flex items-center justify-end space-x-2">
                            <button
                                @click="copyMessage(message.content)"
                                class="px-3 py-1 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded text-sm transition-colors cursor-pointer"
                            >
                                <span class="material-icons text-xs mr-1">content_copy</span>
                                Copy
                            </button>
                            <button
                                @click="editMessage(message)"
                                class="px-3 py-1 text-green-600 bg-green-50 hover:bg-green-100 rounded text-sm transition-colors cursor-pointer"
                            >
                                <span class="material-icons text-xs mr-1">edit</span>
                                Edit
                            </button>
                            <button
                                @click="showDeleteModal(message)"
                                class="px-3 py-1 text-red-600 bg-red-50 hover:bg-red-100 rounded text-sm transition-colors cursor-pointer"
                            >
                                <span class="material-icons text-xs mr-1">delete</span>
                                Delete
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

                    <!-- Create Saved Message Modal -->
            <div
                v-if="showCreateSavedMessage"
                class="fixed inset-0 z-50 flex items-center justify-center"
            >
                <!-- Backdrop -->
                <div
                    class="absolute inset-0 bg-gray-500 bg-opacity-30 backdrop-blur-sm"
                    @click="showCreateSavedMessage = false"
                ></div>

                <!-- Modal -->
                <div
                    class="relative bg-gray-200 rounded-2xl shadow-2xl max-w-md w-full mx-4"
                >
                <!-- Header -->
                <div
                    class="flex items-center justify-between p-6 border-b border-gray-300 bg-gray-100"
                >
                    <div class="flex items-center space-x-3">
                        <div
                            class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center"
                        >
                            <span class="material-icons text-blue-600"
                                >bookmark</span
                            >
                        </div>
                        <div>
                            <h3 class="text-lg font-semibold text-gray-900">
                                {{ editingMessageId ? 'Edit Saved Message' : 'Create Saved Message' }}
                            </h3>
                            <p class="text-sm text-gray-500">
                                {{ editingMessageId ? 'Update your saved message' : 'Save an important message for later' }}
                            </p>
                        </div>
                    </div>
                    <button
                        @click="closeModal"
                        class="text-gray-400 hover:text-gray-600 transition-colors"
                    >
                        <span class="material-icons">close</span>
                    </button>
                </div>

                <!-- Content -->
                <div class="p-6 bg-gray-100">
                    <form @submit.prevent="createSavedMessage">
                        <div class="mb-4">
                            <label class="block text-sm font-medium text-gray-700 mb-2">
                                Title (Optional)
                            </label>
                            <input
                                v-model="newMessage.title"
                                type="text"
                                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                placeholder="Enter a title for your saved message"
                            />
                        </div>
                        <div class="mb-4">
                            <label class="block text-sm font-medium text-gray-700 mb-2">
                                Message Content
                            </label>
                            <textarea
                                v-model="newMessage.content"
                                rows="6"
                                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                                placeholder="Enter your message content here..."
                                required
                            ></textarea>
                        </div>
                        <div class="mb-6">
                            <label class="block text-sm font-medium text-gray-700 mb-2">
                                Category (Optional)
                            </label>
                            <input
                                v-model="newMessage.category"
                                type="text"
                                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                placeholder="e.g., Work, Personal, Important"
                            />
                        </div>
                        <div class="flex items-center justify-end space-x-3">
                            <button
                                type="button"
                                @click="closeModal"
                                class="px-4 py-2 text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-lg font-medium transition-colors cursor-pointer"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                :disabled="!newMessage.content.trim() || isCreating"
                                class="px-4 py-2 text-white bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                            >
                                <span v-if="isCreating" class="flex items-center space-x-2">
                                    <svg
                                        class="animate-spin h-4 w-4 text-white"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                    >
                                        <circle
                                            class="opacity-25"
                                            cx="12"
                                            cy="12"
                                            r="10"
                                            stroke="currentColor"
                                            stroke-width="4"
                                        ></circle>
                                        <path
                                            class="opacity-75"
                                            fill="currentColor"
                                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                        ></path>
                                    </svg>
                                    {{ editingMessageId ? 'Updating...' : 'Creating...' }}
                                </span>
                                <span v-else>{{ editingMessageId ? 'Update Message' : 'Save Message' }}</span>
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>

        <!-- Delete Confirmation Modal -->
        <div
            v-if="showDeleteConfirmation"
            class="fixed inset-0 z-50 flex items-center justify-center"
        >
            <!-- Backdrop -->
            <div
                class="absolute inset-0 bg-gray-500 bg-opacity-30 backdrop-blur-sm"
                @click="closeDeleteModal"
            ></div>

            <!-- Modal -->
            <div
                class="relative bg-gray-200 rounded-2xl shadow-2xl max-w-md w-full mx-4"
            >
                <!-- Header -->
                <div
                    class="flex items-center justify-between p-6 border-b border-gray-300 bg-gray-100"
                >
                    <div class="flex items-center space-x-3">
                        <div
                            class="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center"
                        >
                            <span class="material-icons text-red-600"
                                >warning</span
                            >
                        </div>
                        <div>
                            <h3 class="text-lg font-semibold text-gray-900">
                                Delete Saved Message
                            </h3>
                            <p class="text-sm text-gray-500">
                                This action cannot be undone
                            </p>
                        </div>
                    </div>
                    <button
                        @click="closeDeleteModal"
                        class="text-gray-400 hover:text-gray-600 transition-colors"
                    >
                        <span class="material-icons">close</span>
                    </button>
                </div>

                <!-- Content -->
                <div class="p-6 bg-gray-100">
                    <div class="mb-6">
                        <p class="text-gray-700 mb-4">
                            Are you sure you want to delete this saved message?
                        </p>
                        <div class="bg-red-50 rounded-lg p-4 border border-red-200">
                            <h4 class="font-medium text-red-800 mb-2">
                                {{ messageToDelete?.title || 'Untitled Message' }}
                            </h4>
                            <p class="text-sm text-red-700 line-clamp-3">
                                {{ messageToDelete?.content }}
                            </p>
                        </div>
                    </div>
                    
                    <div class="flex items-center justify-end space-x-3">
                        <button
                            @click="closeDeleteModal"
                            class="px-4 py-2 text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-lg font-medium transition-colors cursor-pointer"
                        >
                            Cancel
                        </button>
                        <button
                            @click="confirmDelete"
                            :disabled="isDeleting"
                            class="px-4 py-2 text-white bg-red-600 hover:bg-red-700 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2 cursor-pointer"
                        >
                            <span v-if="isDeleting" class="flex items-center space-x-2">
                                <svg
                                    class="animate-spin h-4 w-4 text-white"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                >
                                    <circle
                                        class="opacity-25"
                                        cx="12"
                                        cy="12"
                                        r="10"
                                        stroke="currentColor"
                                        stroke-width="4"
                                    ></circle>
                                    <path
                                        class="opacity-75"
                                        fill="currentColor"
                                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                    ></path>
                                </svg>
                                Deleting...
                            </span>
                            <span v-else class="flex items-center space-x-2">
                                <span class="material-icons text-sm">delete</span>
                                Delete Message
                            </span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { showInfo, showError } from "../../utils/toast";
import axiosInstance from "../../axiosInstance";

// Reactive data
const savedMessages = ref([]);
const isLoading = ref(false);
const isCreating = ref(false);
const showCreateSavedMessage = ref(false);
const newMessage = ref({
    title: "",
    content: "",
    category: "",
});
const editingMessageId = ref(null);
const showDeleteConfirmation = ref(false);
const messageToDelete = ref(null);
const isDeleting = ref(false);
const sortBy = ref('date');

// Methods
const loadSavedMessages = async () => {
    try {
        isLoading.value = true;
        const response = await axiosInstance.get("/api/save-message/get");
        
        // Handle response structure with "messages" key
        if (response.data && response.data.messages) {
            savedMessages.value = response.data.messages;
        } else if (response.data && Array.isArray(response.data)) {
            savedMessages.value = response.data;
        } else {
            savedMessages.value = [];
        }
    } catch (error) {
        console.error("Failed to load saved messages:", error);
        showError("Failed to load saved messages. Please try again.");
        savedMessages.value = [];
    } finally {
        isLoading.value = false;
    }
};

const refreshSavedMessages = () => {
    loadSavedMessages();
};

const createSavedMessage = async () => {
    if (!newMessage.value.content.trim()) {
        showError("Please enter message content.");
        return;
    }

    try {
        isCreating.value = true;
        
        if (editingMessageId.value) {
            // Update existing message
            const response = await axiosInstance.put(`/api/save-message/update/${editingMessageId.value}`, {
                title: newMessage.value.title,
                content: newMessage.value.content,
                category: newMessage.value.category,
            });
            showInfo("Message updated successfully!");
        } else {
            // Create new message
            const response = await axiosInstance.post("/api/save-message/create", {
                title: newMessage.value.title,
                content: newMessage.value.content,
                category: newMessage.value.category,
            });
            showInfo("Message saved successfully!");
        }
        
        // Reset form
        newMessage.value = {
            title: "",
            content: "",
            category: "",
        };
        editingMessageId.value = null;
        
        showCreateSavedMessage.value = false;
        
        // Reload messages
        await loadSavedMessages();
    } catch (error) {
        console.error("Failed to save message:", error);
        showError(
            error.response?.data?.message ||
                "Failed to save message. Please try again."
        );
    } finally {
        isCreating.value = false;
    }
};

const copyMessage = async (content) => {
    if (!content) {
        showError("No content to copy.");
        return;
    }
    
    try {
        await navigator.clipboard.writeText(content);
        showInfo("Message copied to clipboard!");
    } catch (error) {
        console.error("Failed to copy message:", error);
        showError("Failed to copy message to clipboard.");
    }
};

const editMessage = (message) => {
    // Check if message exists and has required properties
    if (!message || !message.id) {
        showError("Invalid message data. Please try again.");
        return;
    }
    
    // Populate the form with existing message data
    newMessage.value = {
        title: message.title || "",
        content: message.content || "",
        category: message.category || "",
    };
    
    // Store the message ID for editing
    editingMessageId.value = message.id;
    
    // Show the modal
    showCreateSavedMessage.value = true;
};

const showDeleteModal = (message) => {
    if (!message || !message.id) {
        showError("Invalid message data. Please try again.");
        return;
    }
    
    messageToDelete.value = message;
    showDeleteConfirmation.value = true;
};

const closeDeleteModal = () => {
    showDeleteConfirmation.value = false;
    messageToDelete.value = null;
    isDeleting.value = false;
};

const confirmDelete = async () => {
    if (!messageToDelete.value || !messageToDelete.value.id) {
        showError("Invalid message data. Please try again.");
        return;
    }

    try {
        isDeleting.value = true;
        await axiosInstance.delete(`/api/save-message/delete/${messageToDelete.value.id}`);
        showInfo("Message deleted successfully!");
        await loadSavedMessages();
        closeDeleteModal();
    } catch (error) {
        console.error("Failed to delete saved message:", error);
        showError(
            error.response?.data?.message ||
                "Failed to delete message. Please try again."
        );
    } finally {
        isDeleting.value = false;
    }
};

const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString() + " " + date.toLocaleTimeString();
};

const closeModal = () => {
    showCreateSavedMessage.value = false;
    editingMessageId.value = null;
    newMessage.value = {
        title: "",
        content: "",
        category: "",
    };
};

const getRecentCount = () => {
    const sevenDaysAgo = new Date();
    sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);
    return savedMessages.value.filter(msg => new Date(msg.created_at) > sevenDaysAgo).length;
};

const getCategoryCount = () => {
    const categories = new Set(savedMessages.value.map(msg => msg.category).filter(cat => cat));
    return categories.size;
};

const sortMessages = () => {
    const sorted = [...savedMessages.value];
    switch (sortBy.value) {
        case 'date':
            sorted.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
            break;
        case 'title':
            sorted.sort((a, b) => (a.title || '').localeCompare(b.title || ''));
            break;
        case 'category':
            sorted.sort((a, b) => (a.category || '').localeCompare(b.category || ''));
            break;
    }
    savedMessages.value = sorted;
};

const exportMessages = () => {
    try {
        const data = savedMessages.value.map(msg => ({
            title: msg.title || 'Untitled',
            content: msg.content,
            category: msg.category || 'Uncategorized',
            created_at: formatDate(msg.created_at)
        }));
        
        const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `saved-messages-${new Date().toISOString().split('T')[0]}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        
        showInfo('Messages exported successfully!');
    } catch (error) {
        console.error('Failed to export messages:', error);
        showError('Failed to export messages. Please try again.');
    }
};

// Load saved messages on mount
onMounted(() => {
    loadSavedMessages();
});
</script> 