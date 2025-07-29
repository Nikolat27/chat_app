<template>
    <div v-if="isVisible && group" class="fixed inset-0 z-50 flex items-center justify-center">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black bg-opacity-50 backdrop-blur-sm"></div>
        
        <!-- Modal -->
        <div class="relative bg-white rounded-2xl shadow-2xl max-w-md w-full mx-4 transform transition-all duration-300 scale-100">
            <!-- Header -->
            <div class="flex items-center justify-between p-6 border-b border-gray-100">
                <div class="flex items-center space-x-3">
                    <div class="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center">
                        <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
                        </svg>
                    </div>
                    <div>
                        <h3 class="text-lg font-semibold text-gray-900">
                            {{ isOwner ? 'Delete Group' : 'Leave Group' }}
                        </h3>
                        <p class="text-sm text-gray-500">
                            {{ isOwner ? 'This action cannot be undone' : 'You can rejoin later if needed' }}
                        </p>
                    </div>
                </div>
                <button 
                    @click="closeModal"
                    class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
                >
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                    </svg>
                </button>
            </div>
            
            <!-- Content -->
            <div class="p-6">
                <!-- Group Info -->
                <div v-if="group" class="flex items-center space-x-4 mb-6 p-4 bg-gray-50 rounded-xl">
                    <img
                        :src="group.avatar_url ? `${backendBaseUrl}/static/${group.avatar_url}` : '/src/assets/default-avatar.jpg'"
                        alt="Group Avatar"
                        class="w-12 h-12 rounded-full object-cover border-2 border-gray-200"
                    />
                    <div class="flex-1">
                        <h4 class="font-medium text-gray-900">{{ group.name }}</h4>
                        <p class="text-sm text-gray-500">{{ group.description }}</p>
                        <div class="flex items-center space-x-2 mt-1">
                            <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                                  :class="group.type === 'public' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'">
                                <svg v-if="group.type === 'public'" class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                                    <path d="M10 12a2 2 0 100-4 2 2 0 000 4z"></path>
                                    <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd"></path>
                                </svg>
                                <svg v-else class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                                    <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd"></path>
                                </svg>
                                {{ group.type === 'public' ? 'Public' : 'Private' }}
                            </span>
                            <span class="text-xs text-gray-400">{{ group.member_count }} members</span>
                        </div>
                    </div>
                </div>
                
                <!-- Warning Message -->
                <div class="mb-6 p-4 rounded-xl"
                     :class="isOwner ? 'bg-red-50 border border-red-200' : 'bg-yellow-50 border border-yellow-200'">
                    <div class="flex items-start space-x-3">
                        <div class="flex-shrink-0">
                            <svg v-if="isOwner" class="w-5 h-5 text-red-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
                            </svg>
                            <svg v-else class="w-5 h-5 text-yellow-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
                            </svg>
                        </div>
                        <div>
                            <h4 class="font-medium"
                                :class="isOwner ? 'text-red-800' : 'text-yellow-800'">
                                {{ isOwner ? 'Delete Group Permanently' : 'Leave Group' }}
                            </h4>
                            <p class="text-sm mt-1"
                               :class="isOwner ? 'text-red-700' : 'text-yellow-700'">
                                {{ isOwner 
                                    ? 'This will permanently delete the group and remove all members. This action cannot be undone.' 
                                    : 'You will be removed from this group. You can rejoin later if the group is public or you have an invite link.' 
                                }}
                            </p>
                        </div>
                    </div>
                </div>
                
                <!-- Confirmation Input for Owner -->
                <div v-if="isOwner" class="mb-6">
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                        Type "DELETE" to confirm
                    </label>
                    <input
                        v-model="deleteConfirmation"
                        type="text"
                        placeholder="DELETE"
                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-colors duration-200"
                        :class="{ 'border-red-500': deleteConfirmation && deleteConfirmation !== 'DELETE' }"
                    />
                </div>
            </div>
            
            <!-- Footer -->
            <div class="flex items-center justify-end space-x-3 p-6 border-t border-gray-100">
                <button
                    @click="closeModal"
                    class="px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg font-medium transition-colors duration-200 cursor-pointer"
                >
                    Cancel
                </button>
                <button
                    @click="handleAction"
                    :disabled="isOwner && deleteConfirmation !== 'DELETE'"
                    class="px-4 py-2 text-white font-medium rounded-lg transition-all duration-200 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    :class="isOwner 
                        ? 'bg-red-600 hover:bg-red-700 focus:ring-2 focus:ring-red-500 focus:ring-offset-2' 
                        : 'bg-red-500 hover:bg-red-600 focus:ring-2 focus:ring-red-500 focus:ring-offset-2'"
                >
                    <span v-if="isLoading" class="flex items-center space-x-2">
                        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        {{ isOwner ? 'Deleting...' : 'Leaving...' }}
                    </span>
                    <span v-else>
                        {{ isOwner ? 'Delete Group' : 'Leave Group' }}
                    </span>
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useGroupStore } from '../../stores/groups';
import { useUserStore } from '../../stores/users';
import { showMessage, showError } from '../../utils/toast';

// Props
const props = defineProps({
    isVisible: {
        type: Boolean,
        default: false
    },
    group: {
        type: Object,
        required: false,
        default: null
    },
    backendBaseUrl: {
        type: String,
        default: import.meta.env.VITE_BACKEND_BASE_URL
    }
});

// Emits
const emit = defineEmits(['close', 'action-completed']);

// Stores
const groupStore = useGroupStore();
const userStore = useUserStore();

// Reactive data
const isLoading = ref(false);
const deleteConfirmation = ref('');

// Computed
const isOwner = computed(() => {
    return userStore.user_id && props.group && props.group.owner_id === userStore.user_id;
});

// Methods
const closeModal = () => {
    deleteConfirmation.value = '';
    emit('close');
};

const handleAction = async () => {
    if (isOwner.value && deleteConfirmation.value !== 'DELETE') {
        showError('Please type "DELETE" to confirm');
        return;
    }
    
    isLoading.value = true;
    
    try {
        if (isOwner.value) {
            // Delete group
            await groupStore.deleteGroup(props.group.id);
            showMessage('Group deleted successfully');
        } else {
            // Leave group
            await groupStore.leaveGroup(props.group.id);
            showMessage('Left group successfully');
        }
        
        emit('action-completed');
        closeModal();
    } catch (error) {
        console.error('Action failed:', error);
        showError(error.response?.data?.message || 'Action failed. Please try again.');
    } finally {
        isLoading.value = false;
    }
};
</script> 