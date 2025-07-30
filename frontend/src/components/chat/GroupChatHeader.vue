<template>
    <div class="flex items-center justify-between p-4 bg-white border-b border-gray-200 shadow-sm">
        <!-- Group Info -->
        <div class="flex items-center space-x-3">
            <!-- Group Avatar -->
            <div class="relative">
                <img 
                    :src="group.avatar_url ? `${backendBaseUrl}/static/${group.avatar_url}` : '/src/assets/default-avatar.jpg'"
                    :alt="group.name"
                    :class="[
                        'w-12 h-12 rounded-full object-cover border-2',
                        group.type === 'secret' ? 'border-purple-300' : 'border-gray-200'
                    ]"
                />
                <div 
                    :class="[
                        'absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-white',
                        group.type === 'secret' ? 'bg-purple-500' : 'bg-green-500'
                    ]"
                >
                    <span 
                        v-if="group.type === 'secret'"
                        class="material-icons text-white text-xs"
                    >lock</span>
                </div>
            </div>
            
            <!-- Group Details -->
            <div class="flex flex-col">
                <div class="flex items-center gap-2">
                    <h3 class="text-lg font-semibold text-gray-900">{{ group.name }}</h3>
                    <div 
                        v-if="group.type === 'secret'"
                        class="px-2 py-1 bg-purple-100 text-purple-700 rounded-full text-xs font-medium border border-purple-200 flex items-center gap-1"
                    >
                        <span class="material-icons text-xs">security</span>
                        Secret
                    </div>
                </div>
                <div class="flex items-center space-x-2 text-sm text-gray-500">
                    <span>{{ group.member_count || 0 }} members</span>
                    <span>•</span>
                    <span class="capitalize">{{ group.type }}</span>
                    <span v-if="group.owner_id === currentUserId" class="text-blue-600 font-medium">• Admin</span>
                </div>
            </div>
        </div>
        
        <!-- Actions -->
        <div class="flex items-center space-x-2">
            <!-- Copy Invite Link -->
            <button
                v-if="group.invite_link"
                @click="copyInviteLink"
                class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors cursor-pointer"
                title="Copy invite link"
            >
                <span class="material-icons text-lg">content_copy</span>
            </button>
            
            <!-- Group Actions Menu -->
            <div class="relative">
                <button
                    @click="toggleActionsMenu"
                    class="p-2 text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded-lg transition-colors cursor-pointer"
                    title="Group actions"
                >
                    <span class="material-icons text-lg">more_vert</span>
                </button>
                
                <!-- Dropdown Menu -->
                <div
                    v-if="showActionsMenu"
                    class="absolute right-0 top-full mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 z-10"
                >
                    <div class="py-1">
                        <!-- Leave Group (for non-owners) -->
                        <button
                            v-if="group.owner_id !== currentUserId"
                            @click="handleLeaveGroup"
                            class="w-full px-4 py-2 text-left text-red-600 hover:bg-red-50 flex items-center space-x-2"
                        >
                            <span class="material-icons text-sm">exit_to_app</span>
                            <span>Leave Group</span>
                        </button>
                        
                        <!-- Delete Group (for owners) -->
                        <button
                            v-if="group.owner_id === currentUserId"
                            @click="handleDeleteGroup"
                            class="w-full px-4 py-2 text-left text-red-600 hover:bg-red-50 flex items-center space-x-2"
                        >
                            <span class="material-icons text-sm">delete</span>
                            <span>Delete Group</span>
                        </button>
                        
                        <!-- Group Info -->
                        <button
                            @click="showGroupInfo"
                            class="w-full px-4 py-2 text-left text-gray-700 hover:bg-gray-50 flex items-center space-x-2"
                        >
                            <span class="material-icons text-sm">info</span>
                            <span>Group Info</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { showMessage, showError } from '../../utils/toast';

// Props
const props = defineProps({
    group: {
        type: Object,
        required: true
    },
    backendBaseUrl: {
        type: String,
        required: true
    },
    currentUserId: {
        type: String,
        required: true
    }
});

// Emits
const emit = defineEmits(['leave-group', 'delete-group', 'show-group-info']);

// Reactive data
const showActionsMenu = ref(false);

// Methods
const toggleActionsMenu = () => {
    showActionsMenu.value = !showActionsMenu.value;
};

const copyInviteLink = async () => {
    try {
        await navigator.clipboard.writeText(props.group.invite_link);
        showMessage('Invite link copied to clipboard!');
    } catch (error) {
        console.error('Failed to copy invite link:', error);
        showError('Failed to copy invite link. Please try again.');
    }
};

const handleLeaveGroup = () => {
    showActionsMenu.value = false;
    emit('leave-group', props.group);
};

const handleDeleteGroup = () => {
    showActionsMenu.value = false;
    emit('delete-group', props.group);
};

const showGroupInfo = () => {
    showActionsMenu.value = false;
    emit('show-group-info', props.group);
};

// Close menu when clicking outside
const closeMenu = (event) => {
    if (!event.target.closest('.relative')) {
        showActionsMenu.value = false;
    }
};

// Add click listener to close menu
if (typeof window !== 'undefined') {
    window.addEventListener('click', closeMenu);
}
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style> 