<template>
    <div class="bg-white border-b border-gray-200 px-6 py-4">
        <div class="flex items-center justify-between">
            <!-- Group Info -->
            <div class="flex items-center gap-4">
                <!-- Group Avatar -->
                <div class="relative">
                    <img
                        :src="group.avatar_url ? `${backendBaseUrl}/static/${group.avatar_url}` : '/src/assets/default-avatar.jpg'"
                        alt="Group Avatar"
                        class="w-12 h-12 rounded-full object-cover border-2 border-green-300 shadow-sm"
                    />
                    <div class="absolute -bottom-1 -right-1 w-5 h-5 bg-green-500 rounded-full flex items-center justify-center">
                        <span class="material-icons text-white text-xs">group</span>
                    </div>
                </div>

                <!-- Group Details -->
                <div class="flex flex-col">
                    <div class="flex items-center gap-2">
                        <h2 class="text-lg font-semibold text-gray-800">{{ group.name }}</h2>
                        <div 
                            :class="group.is_secret 
                                ? 'bg-purple-100 text-purple-700 border-purple-200' 
                                : group.type === 'private'
                                ? 'bg-orange-100 text-orange-700 border-orange-200'
                                : 'bg-green-100 text-green-700 border-green-200'"
                            class="px-2 py-1 rounded-full text-xs font-medium border flex items-center gap-1"
                        >
                            <span class="material-icons text-xs">
                                {{ group.is_secret ? 'lock' : group.type === 'private' ? 'lock_outline' : 'group' }}
                            </span>
                            {{ group.is_secret ? 'Secret' : group.type === 'private' ? 'Private' : 'Public' }}
                        </div>
                    </div>
                    <p class="text-sm text-gray-600">{{ group.description || 'No description' }}</p>
                    <div class="flex items-center gap-2 mt-1">
                        <span class="text-xs text-gray-500">{{ group.member_count || 0 }} members</span>
                        <span class="text-xs text-gray-400">•</span>
                        <span class="text-xs text-gray-500">Created {{ formatDate(group.created_at) }}</span>
                    </div>
                </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex items-center gap-2">
                <button
                    @click="showGroupInfo = true"
                    class="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full transition-colors"
                    title="Group info"
                >
                    <span class="material-icons text-lg">info</span>
                </button>
                <button
                    @click="handleLeaveGroup"
                    class="p-2 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-full transition-colors"
                    title="Leave group"
                >
                    <span class="material-icons text-lg">exit_to_app</span>
                </button>
            </div>
        </div>

        <!-- Group Info Modal -->
        <div v-if="showGroupInfo" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div class="bg-white rounded-2xl p-6 w-96 max-w-[90vw] max-h-[80vh] overflow-y-auto">
                <div class="flex items-center justify-between mb-4">
                    <h3 class="text-lg font-bold text-gray-800">Group Information</h3>
                    <button
                        @click="showGroupInfo = false"
                        class="text-gray-400 hover:text-gray-600"
                    >
                        <span class="material-icons">close</span>
                    </button>
                </div>

                <!-- Group Details -->
                <div class="space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Group Name</label>
                        <p class="text-gray-900">{{ group.name }}</p>
                    </div>
                    
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                        <p class="text-gray-900">{{ group.description || 'No description' }}</p>
                    </div>
                    
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
                        <div class="flex items-center gap-2">
                            <span 
                                :class="group.is_secret 
                                    ? 'bg-purple-100 text-purple-700 border-purple-200' 
                                    : group.type === 'private'
                                    ? 'bg-orange-100 text-orange-700 border-orange-200'
                                    : 'bg-green-100 text-green-700 border-green-200'"
                                class="px-2 py-1 rounded-full text-xs font-medium border flex items-center gap-1"
                            >
                                <span class="material-icons text-xs">
                                    {{ group.is_secret ? 'lock' : group.type === 'private' ? 'lock_outline' : 'group' }}
                                </span>
                                {{ group.is_secret ? 'Secret' : group.type === 'private' ? 'Private' : 'Public' }}
                            </span>
                        </div>
                    </div>
                    
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Created</label>
                        <p class="text-gray-900">{{ formatDate(group.created_at) }}</p>
                    </div>
                    
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Members</label>
                        <p class="text-gray-900">{{ group.member_count || 0 }} members</p>
                    </div>
                    
                    <div v-if="group.invite_code">
                        <label class="block text-sm font-medium text-gray-700 mb-1">Invite Code</label>
                        <div class="flex items-center gap-2">
                            <input
                                :value="group.invite_code"
                                readonly
                                class="flex-1 p-2 border border-gray-300 rounded-lg bg-gray-50 text-sm"
                            />
                            <button
                                @click="copyInviteCode"
                                class="p-2 text-blue-500 hover:text-blue-700 hover:bg-blue-50 rounded-lg transition-colors"
                                title="Copy invite code"
                            >
                                <span class="material-icons text-sm">content_copy</span>
                            </button>
                        </div>
                    </div>
                </div>

                <div class="flex gap-3 mt-6">
                    <button
                        @click="showGroupInfo = false"
                        class="flex-1 px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                    >
                        Close
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { showMessage, showError } from '../../utils/toast';
import { useGroupStore } from '../../stores/groups';

const props = defineProps({
    group: {
        type: Object,
        required: true
    },
    backendBaseUrl: {
        type: String,
        default: import.meta.env.VITE_BACKEND_BASE_URL
    }
});

const emit = defineEmits(['leave-group']);

const groupStore = useGroupStore();
const showGroupInfo = ref(false);

const formatDate = (dateString) => {
    if (!dateString) return 'Unknown';
    return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
};

const handleLeaveGroup = async () => {
    try {
        const confirmed = confirm(`Are you sure you want to leave "${props.group.name}"?`);
        if (!confirmed) return;

        await groupStore.leaveGroup(props.group.id);
        showMessage('Successfully left group');
        emit('leave-group');
    } catch (error) {
        console.error('Failed to leave group:', error);
        showError('Failed to leave group. Please try again.');
    }
};

const copyInviteCode = async () => {
    try {
        await navigator.clipboard.writeText(props.group.invite_code);
        showMessage('Invite code copied to clipboard!');
    } catch (error) {
        console.error('Failed to copy invite code:', error);
        showError('Failed to copy invite code');
    }
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style> 