<template>
    <div
        v-if="isVisible && group"
        class="fixed inset-0 bg-gray-300 bg-opacity-50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
        @click.self="closeModal"
    >
        <div
            class="bg-white rounded-3xl shadow-2xl border border-gray-100 p-8 w-96 max-w-[90vw] max-h-[90vh] overflow-y-auto transform transition-all duration-300 scale-100 hover:shadow-3xl"
        >
            <!-- Header -->
            <div class="flex items-center justify-between mb-6">
                <h2 class="text-2xl font-bold text-gray-800 flex items-center">
                    <span class="material-icons text-blue-600 mr-2">info</span>
                    Group Information
                </h2>
                <button
                    @click="closeModal"
                    class="text-gray-400 hover:text-gray-600 transition-colors cursor-pointer"
                >
                    <span class="material-icons text-2xl">close</span>
                </button>
            </div>

            <!-- Group Avatar and Basic Info -->
            <div class="mb-6 p-6 bg-gradient-to-br from-blue-50 to-indigo-50 rounded-2xl border border-blue-100">
                <div class="flex items-center gap-4">
                    <img
                        :src="getAvatarUrl(group.avatar_url)"
                        :alt="group.name"
                        class="w-16 h-16 rounded-full object-cover border-4 border-white shadow-lg"
                    />
                    <div class="flex-1">
                        <h3 class="text-xl font-bold text-gray-800 mb-1">{{ group.name }}</h3>
                        <div class="flex items-center gap-2">
                            <span
                                :class="
                                    group.type === 'public'
                                        ? 'bg-green-100 text-green-700 border-green-200'
                                        : group.type === 'private'
                                        ? 'bg-orange-100 text-orange-700 border-orange-200'
                                        : 'bg-purple-100 text-purple-700 border-purple-200'
                                "
                                class="px-3 py-1 text-xs rounded-full font-medium border"
                            >
                                {{ group.type.charAt(0).toUpperCase() + group.type.slice(1) }}
                            </span>
                            <span class="text-sm text-gray-500">
                                {{ group.member_count || 0 }} members
                            </span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Group Details -->
            <div class="space-y-4">
                <!-- Description -->
                <div class="bg-gray-50 rounded-xl p-4">
                    <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                        <span class="material-icons text-gray-500 mr-2 text-sm">description</span>
                        Description
                    </h4>
                    <p class="text-gray-600 text-sm leading-relaxed">
                        {{ group.description || 'No description available' }}
                    </p>
                </div>

                <!-- Created Date -->
                <div class="bg-gray-50 rounded-xl p-4">
                    <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                        <span class="material-icons text-gray-500 mr-2 text-sm">schedule</span>
                        Created
                    </h4>
                    <p class="text-gray-600 text-sm">
                        {{ formatDate(group.created_at) }}
                    </p>
                </div>

                <!-- Invite Link -->
                <div v-if="group.invite_link" class="bg-gray-50 rounded-xl p-4">
                    <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                        <span class="material-icons text-gray-500 mr-2 text-sm">link</span>
                        Invite Link
                    </h4>
                    <div class="flex items-center gap-2">
                        <input
                            :value="group.invite_link"
                            readonly
                            class="flex-1 px-3 py-2 text-xs bg-white border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                        <button
                            @click="copyInviteLink"
                            class="px-3 py-2 bg-blue-500 hover:bg-blue-600 text-white text-xs rounded-lg transition-colors duration-200 flex items-center gap-1"
                            title="Copy invite link"
                        >
                            <span class="material-icons text-sm">content_copy</span>
                        </button>
                    </div>
                </div>

                <!-- Owner Info -->
                <div class="bg-gray-50 rounded-xl p-4">
                    <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                        <span class="material-icons text-gray-500 mr-2 text-sm">person</span>
                        Owner
                    </h4>
                    <div class="flex items-center gap-3">
                        <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center">
                            <span class="material-icons text-white text-sm">person</span>
                        </div>
                        <span class="text-gray-600 text-sm">Group Owner</span>
                    </div>
                </div>
            </div>

            <!-- Action Buttons -->
            <div class="mt-8 flex gap-3">
                <button
                    @click="closeModal"
                    class="flex-1 px-4 py-3 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-xl font-medium transition-colors duration-200"
                >
                    Close
                </button>
                <button
                    v-if="group.owner_id === currentUserId"
                    @click="handleEditGroup"
                    class="flex-1 px-4 py-3 bg-blue-500 hover:bg-blue-600 text-white rounded-xl font-medium transition-colors duration-200"
                >
                    Edit Group
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { showMessage } from '../../utils/toast';

const props = defineProps({
    isVisible: {
        type: Boolean,
        default: false,
    },
    group: {
        type: Object,
        default: null,
    },
    backendBaseUrl: {
        type: String,
        required: true,
    },
    currentUserId: {
        type: String,
        required: true,
    },
});

const emit = defineEmits(['close', 'edit-group']);

const getAvatarUrl = (avatarUrl) => {
    if (!avatarUrl) return '/src/assets/default-avatar.jpg';
    if (avatarUrl.startsWith('http')) return avatarUrl;
    return `${props.backendBaseUrl}/static/${avatarUrl}`;
};

const formatDate = (timestamp) => {
    if (!timestamp) return 'Unknown';
    const date = new Date(timestamp);
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
};

const copyInviteLink = async () => {
    try {
        await navigator.clipboard.writeText(props.group.invite_link);
        showMessage('Invite link copied to clipboard!');
    } catch (error) {
        console.error('Failed to copy invite link:', error);
        showMessage('Failed to copy invite link');
    }
};

const handleEditGroup = () => {
    emit('edit-group', props.group);
    closeModal();
};

const closeModal = () => {
    emit('close');
};
</script>

<style scoped>
/* Custom scrollbar for the modal */
.overflow-y-auto::-webkit-scrollbar {
    width: 4px;
}

.overflow-y-auto::-webkit-scrollbar-track {
    background: transparent;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
    background: rgba(156, 163, 175, 0.3);
    border-radius: 2px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
    background: rgba(156, 163, 175, 0.5);
}
</style> 