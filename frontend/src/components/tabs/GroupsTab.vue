<template>
    <div class="h-full">
        <!-- Header -->
        <div class="mb-6">
            <div class="flex items-center gap-3 mb-2">
                <span class="material-icons text-green-600 text-xl">group</span>
                <h3 class="text-lg font-bold text-gray-800">Groups</h3>
                <div class="flex-1"></div>
                <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-full">
                    {{ groups ? groups.length : 0 }} groups
                </span>
            </div>
            <p class="text-sm text-gray-600">
                Create and join group conversations with multiple participants
            </p>
        </div>

        <!-- Action Buttons -->
        <div class="space-y-3 mb-6">
            <!-- Join Group Button -->
            <button
                @click="showJoinGroupModal = true"
                class="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3 px-4 rounded-xl transition-all duration-200 flex items-center justify-center gap-3 shadow-lg hover:shadow-xl transform hover:scale-[1.02]"
            >
                <span class="material-icons text-lg">group_add</span>
                <span>Join Group</span>
            </button>

            <!-- Create Group Button -->
            <button
                @click="showCreateGroupModal = true"
                class="w-full bg-green-500 hover:bg-green-600 text-white font-semibold py-3 px-4 rounded-xl transition-all duration-200 flex items-center justify-center gap-3 shadow-lg hover:shadow-xl transform hover:scale-[1.02]"
            >
                <span class="material-icons text-lg">add_circle</span>
                <span>Create Group</span>
            </button>

            <!-- Create Secret Group Button -->
            <button
                @click="showCreateSecretGroupModal = true"
                class="w-full bg-purple-500 hover:bg-purple-600 text-white font-semibold py-3 px-4 rounded-xl transition-all duration-200 flex items-center justify-center gap-3 shadow-lg hover:shadow-xl transform hover:scale-[1.02]"
            >
                <span class="material-icons text-lg">lock</span>
                <span>Create Secret Group</span>
            </button>
        </div>

        <!-- Groups List -->
        <div v-if="groups && groups.length > 0" class="space-y-3">
            <h4 class="text-sm font-semibold text-gray-700 mb-3">Your Groups</h4>
            <div
                v-for="group in groups"
                :key="group.id"
                class="group bg-white rounded-2xl shadow-sm border border-gray-200 cursor-pointer hover:shadow-lg hover:border-green-300 transition-all duration-300 overflow-hidden"
            >
                <div class="p-4">
                    <div class="flex items-center gap-4">
                        <!-- Group Avatar -->
                        <div class="relative">
                            <img
                                :src="group.avatar_url || '/src/assets/default-group-avatar.jpg'"
                                alt="Group Avatar"
                                class="w-12 h-12 rounded-full object-cover border-2 border-green-300 shadow-sm group-hover:border-green-400 transition-colors duration-200 select-none pointer-events-none"
                            />
                            <div class="absolute -bottom-1 -right-1 w-5 h-5 bg-green-500 rounded-full flex items-center justify-center">
                                <span class="material-icons text-white text-xs">group</span>
                            </div>
                        </div>

                        <!-- Group Info -->
                        <div class="flex-1 min-w-0">
                            <div class="flex items-center gap-2 mb-1">
                                <span class="font-semibold text-gray-800 truncate">
                                    {{ group.name }}
                                </span>
                                <div 
                                    :class="group.is_secret 
                                        ? 'bg-purple-100 text-purple-700 border-purple-200' 
                                        : 'bg-green-100 text-green-700 border-green-200'"
                                    class="px-2 py-1 rounded-full text-xs font-medium border flex items-center gap-1"
                                >
                                    <span class="material-icons text-xs">
                                        {{ group.is_secret ? 'lock' : 'group' }}
                                    </span>
                                    {{ group.is_secret ? 'Secret' : 'Public' }}
                                </div>
                            </div>
                            <p class="text-sm text-gray-600 truncate">
                                {{ group.description || 'No description' }}
                            </p>
                            <div class="flex items-center gap-2 mt-1">
                                <span class="text-xs text-gray-500">
                                    {{ group.member_count || 0 }} members
                                </span>
                                <span class="text-xs text-gray-400">•</span>
                                <span class="text-xs text-gray-500">
                                    {{ group.role || 'Member' }}
                                </span>
                            </div>
                        </div>

                        <!-- Action Buttons -->
                        <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-all duration-200">
                            <button
                                @click.stop="handleGroupClick(group)"
                                class="w-8 h-8 text-blue-500 hover:bg-blue-50 rounded-full hover:text-blue-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 flex items-center justify-center"
                                title="Open group chat"
                            >
                                <span class="material-icons text-sm">chat</span>
                            </button>
                            <button
                                @click.stop="handleLeaveGroup(group)"
                                class="w-8 h-8 text-red-500 hover:bg-red-50 rounded-full hover:text-red-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 flex items-center justify-center"
                                title="Leave group"
                            >
                                <span class="material-icons text-sm">exit_to_app</span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Empty State -->
        <div
            v-else
            class="text-center py-12 px-6"
        >
            <div class="mb-4">
                <span class="material-icons text-6xl text-gray-300">group</span>
            </div>
            <h4 class="text-lg font-semibold text-gray-600 mb-2">No Groups Yet</h4>
            <p class="text-sm text-gray-500 mb-6 leading-relaxed">
                Start by creating a new group or joining an existing one to begin group conversations.
            </p>
            <div class="bg-green-50 rounded-xl p-4 border border-green-200">
                <div class="flex items-center gap-2 mb-2">
                    <span class="material-icons text-green-600 text-sm">info</span>
                    <span class="text-sm font-semibold text-green-700">Group Features</span>
                </div>
                <ul class="text-xs text-green-600 space-y-1">
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">check_circle</span>
                        Multiple participants
                    </li>
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">check_circle</span>
                        Public and private groups
                    </li>
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">check_circle</span>
                        Secret groups with E2EE
                    </li>
                </ul>
            </div>
        </div>

        <!-- Modals will be added here later -->
        <!-- Join Group Modal -->
        <div v-if="showJoinGroupModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div class="bg-white rounded-2xl p-6 w-96 max-w-[90vw]">
                <h3 class="text-lg font-bold text-gray-800 mb-4">Join Group</h3>
                <p class="text-sm text-gray-600 mb-4">Enter the group code or name to join an existing group.</p>
                <input
                    v-model="joinGroupCode"
                    type="text"
                    placeholder="Enter group code or name"
                    class="w-full p-3 border border-gray-300 rounded-lg mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <div class="flex gap-3">
                    <button
                        @click="showJoinGroupModal = false"
                        class="flex-1 px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        @click="handleJoinGroup"
                        class="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                    >
                        Join Group
                    </button>
                </div>
            </div>
        </div>

        <!-- Create Group Modal -->
        <div v-if="showCreateGroupModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div class="bg-white rounded-2xl p-6 w-96 max-w-[90vw]">
                <h3 class="text-lg font-bold text-gray-800 mb-4">Create Group</h3>
                <div class="space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Group Name</label>
                        <input
                            v-model="newGroup.name"
                            type="text"
                            placeholder="Enter group name"
                            class="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                        />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Description (Optional)</label>
                        <textarea
                            v-model="newGroup.description"
                            placeholder="Enter group description"
                            rows="3"
                            class="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                        ></textarea>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Group Type</label>
                        <select
                            v-model="newGroup.type"
                            class="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                        >
                            <option value="public">Public (Anyone can join)</option>
                            <option value="private">Private (Invite only)</option>
                        </select>
                    </div>
                </div>
                <div class="flex gap-3 mt-6">
                    <button
                        @click="showCreateGroupModal = false"
                        class="flex-1 px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        @click="handleCreateGroup"
                        class="flex-1 px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
                    >
                        Create Group
                    </button>
                </div>
            </div>
        </div>

        <!-- Create Secret Group Modal -->
        <div v-if="showCreateSecretGroupModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div class="bg-white rounded-2xl p-6 w-96 max-w-[90vw]">
                <h3 class="text-lg font-bold text-gray-800 mb-4">Create Secret Group</h3>
                <p class="text-sm text-gray-600 mb-4">Create an end-to-end encrypted group chat for maximum privacy.</p>
                <div class="space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Group Name</label>
                        <input
                            v-model="newSecretGroup.name"
                            type="text"
                            placeholder="Enter secret group name"
                            class="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                        />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Description (Optional)</label>
                        <textarea
                            v-model="newSecretGroup.description"
                            placeholder="Enter group description"
                            rows="3"
                            class="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                        ></textarea>
                    </div>
                </div>
                <div class="flex gap-3 mt-6">
                    <button
                        @click="showCreateSecretGroupModal = false"
                        class="flex-1 px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        @click="handleCreateSecretGroup"
                        class="flex-1 px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition-colors"
                    >
                        Create Secret Group
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { showMessage, showError } from '../../utils/toast';

// Reactive state
const groups = ref([]); // Will be populated from backend
const showJoinGroupModal = ref(false);
const showCreateGroupModal = ref(false);
const showCreateSecretGroupModal = ref(false);
const joinGroupCode = ref('');

// Form data
const newGroup = reactive({
    name: '',
    description: '',
    type: 'public'
});

const newSecretGroup = reactive({
    name: '',
    description: ''
});

// Event handlers
const handleGroupClick = (group) => {
    console.log('Opening group:', group);
    // TODO: Implement group chat opening
    showMessage('Group chat functionality coming soon!');
};

const handleLeaveGroup = async (group) => {
    try {
        console.log('Leaving group:', group);
        // TODO: Implement leave group API call
        showMessage('Leave group functionality coming soon!');
    } catch (error) {
        showError('Failed to leave group');
    }
};

const handleJoinGroup = async () => {
    try {
        if (!joinGroupCode.value.trim()) {
            showError('Please enter a group code or name');
            return;
        }
        
        console.log('Joining group with code:', joinGroupCode.value);
        // TODO: Implement join group API call
        showMessage('Join group functionality coming soon!');
        showJoinGroupModal.value = false;
        joinGroupCode.value = '';
    } catch (error) {
        showError('Failed to join group');
    }
};

const handleCreateGroup = async () => {
    try {
        if (!newGroup.name.trim()) {
            showError('Please enter a group name');
            return;
        }
        
        console.log('Creating group:', newGroup);
        // TODO: Implement create group API call
        showMessage('Create group functionality coming soon!');
        showCreateGroupModal.value = false;
        
        // Reset form
        newGroup.name = '';
        newGroup.description = '';
        newGroup.type = 'public';
    } catch (error) {
        showError('Failed to create group');
    }
};

const handleCreateSecretGroup = async () => {
    try {
        if (!newSecretGroup.name.trim()) {
            showError('Please enter a group name');
            return;
        }
        
        console.log('Creating secret group:', newSecretGroup);
        // TODO: Implement create secret group API call
        showMessage('Create secret group functionality coming soon!');
        showCreateSecretGroupModal.value = false;
        
        // Reset form
        newSecretGroup.name = '';
        newSecretGroup.description = '';
    } catch (error) {
        showError('Failed to create secret group');
    }
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style> 