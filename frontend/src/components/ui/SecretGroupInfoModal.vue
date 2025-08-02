<template>
    <div
        v-if="isVisible && group"
        class="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
        @click.self="$emit('close')"
    >
        <div class="bg-white rounded-3xl shadow-2xl w-full max-w-lg relative font-sans flex flex-col max-h-[90vh]">
            <!-- Close Button -->
            <button
                class="absolute top-4 right-4 text-gray-400 hover:text-purple-700 hover:bg-purple-50 w-10 h-10 rounded-full transition-all duration-200 cursor-pointer flex items-center justify-center z-10"
                @click="$emit('close')"
                aria-label="Close"
            >
                <span class="material-icons text-xl">close</span>
            </button>

            <!-- Header -->
            <div class="text-center p-8 pb-4 flex-shrink-0">
                <div class="mb-4">
                    <span class="material-icons text-5xl text-purple-500 mb-3">security</span>
                </div>
                <h3 class="text-2xl font-bold text-gray-800 mb-2">Secret Group Info</h3>
                <p class="text-sm text-gray-600">
                    End-to-end encrypted group details
                </p>
            </div>

            <!-- Scrollable Content -->
            <div class="flex-1 overflow-y-auto px-8">
            <!-- Group Info -->
                <div class="space-y-6 pb-6">
                <!-- Group Basic Info -->
                <div class="bg-gray-50 rounded-xl p-4">
                    <div class="flex items-center gap-3 mb-3">
                        <img
                            v-if="group.avatar_url"
                            :src="`${backendBaseUrl}/static/${group.avatar_url}`"
                            class="w-12 h-12 rounded-full object-cover border-2 border-purple-300 select-none pointer-events-none"
                            alt="Group Avatar"
                        />
                        <img
                            v-else
                            src="/src/assets/default-avatar.jpg"
                            class="w-12 h-12 rounded-full object-cover border-2 border-purple-300 select-none pointer-events-none"
                            alt="Default Avatar"
                        />
                        <div>
                            <span class="font-semibold text-gray-800">{{ group.name }}</span>
                            <div class="text-xs text-purple-600 flex items-center gap-1">
                                <span class="material-icons text-xs">verified</span>
                                Secret group
                            </div>
                        </div>
                    </div>
                    <p v-if="group.description" class="text-sm text-gray-600 mt-2">
                        {{ group.description }}
                    </p>
                </div>

                <!-- Security Features -->
                <div class="space-y-4">
                    <h4 class="font-semibold text-gray-800 flex items-center gap-2">
                        <span class="material-icons text-purple-600">shield</span>
                        Security Features
                    </h4>
                    
                    <div class="space-y-3">
                        <div class="flex items-start gap-3 p-3 bg-green-50 rounded-lg border border-green-200">
                            <span class="material-icons text-green-600 text-sm mt-0.5">lock</span>
                            <div>
                                <div class="font-medium text-green-800 text-sm">End-to-End Encryption</div>
                                <div class="text-xs text-green-700 mt-1">
                                    All messages are encrypted with unique key pairs for each member.
                                </div>
                            </div>
                        </div>
                        
                        <div class="flex items-start gap-3 p-3 bg-blue-50 rounded-lg border border-blue-200">
                            <span class="material-icons text-blue-600 text-sm mt-0.5">group</span>
                            <div>
                                <div class="font-medium text-blue-800 text-sm">Per-Member Keys</div>
                                <div class="text-xs text-blue-700 mt-1">
                                    Each member has their own public/private key pair for secure communication.
                                </div>
                            </div>
                        </div>
                        
                        <div class="flex items-start gap-3 p-3 bg-purple-50 rounded-lg border border-purple-200">
                            <span class="material-icons text-purple-600 text-sm mt-0.5">key</span>
                            <div>
                                <div class="font-medium text-purple-800 text-sm">Symmetric Key Distribution</div>
                                <div class="text-xs text-purple-700 mt-1">
                                    Message content is encrypted with symmetric keys distributed to each member.
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Group Details -->
                <div class="space-y-4">
                    <h4 class="font-semibold text-gray-800 flex items-center gap-2">
                        <span class="material-icons text-gray-600">info</span>
                        Group Details
                    </h4>
                    
                    <div class="space-y-3">
                        <!-- Member Count -->
                        <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                            <div class="flex items-center gap-2">
                                <span class="material-icons text-gray-500 text-sm">group</span>
                                <span class="text-sm font-medium text-gray-700">Members</span>
                            </div>
                            <span class="text-sm text-gray-600">{{ group.member_count || 0 }}</span>
                        </div>

                        <!-- Created Date -->
                        <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                            <div class="flex items-center gap-2">
                                <span class="material-icons text-gray-500 text-sm">schedule</span>
                                <span class="text-sm font-medium text-gray-700">Created</span>
                            </div>
                            <span class="text-sm text-gray-600">{{ formatDate(group.created_at) }}</span>
                        </div>

                        <!-- Your Role -->
                        <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                            <div class="flex items-center gap-2">
                                <span class="material-icons text-gray-500 text-sm">person</span>
                                <span class="text-sm font-medium text-gray-700">Your Role</span>
                            </div>
                            <span class="text-sm text-gray-600 capitalize">{{ group.role || 'member' }}</span>
                        </div>

                        <!-- Invite Link -->
                        <div v-if="group.invite_link" class="p-3 bg-gray-50 rounded-lg">
                            <div class="flex items-center gap-2 mb-2">
                                <span class="material-icons text-gray-500 text-sm">link</span>
                                <span class="text-sm font-medium text-gray-700">Invite Link</span>
                            </div>
                            <div class="flex items-center gap-2">
                                <input
                                    :value="group.invite_link"
                                    readonly
                                    class="flex-1 px-2 py-1 text-xs bg-white border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-purple-500"
                                />
                                <button
                                    @click="copyInviteLink"
                                    class="px-2 py-1 bg-purple-500 hover:bg-purple-600 text-white text-xs rounded transition-colors duration-200 flex items-center gap-1"
                                    title="Copy invite link"
                                >
                                    <span class="material-icons text-xs">content_copy</span>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Encryption Status -->
                <div class="bg-purple-50 rounded-xl p-4 border border-purple-200">
                    <div class="flex items-center gap-2 mb-2">
                        <span class="material-icons text-purple-600 text-sm">security</span>
                        <span class="font-semibold text-purple-800 text-sm">Encryption Status</span>
                    </div>
                    <div class="text-xs text-purple-700 space-y-1">
                        <div class="flex items-center gap-2">
                            <span class="material-icons text-xs">check_circle</span>
                            <span>Your keys are generated and stored</span>
                        </div>
                        <div class="flex items-center gap-2">
                            <span class="material-icons text-xs">verified</span>
                            <span>Messages are end-to-end encrypted</span>
                        </div>
                        <div class="flex items-center gap-2">
                            <span class="material-icons text-xs">lock</span>
                            <span>Only group members can read messages</span>
                        </div>
                    </div>
                </div>

                    <!-- Secret Key Management -->
                    <div class="bg-blue-50 rounded-xl p-4 border border-blue-200">
                        <div class="flex items-center gap-2 mb-3">
                            <span class="material-icons text-blue-600 text-sm">key</span>
                            <span class="font-semibold text-blue-800 text-sm">Secret Key Management</span>
                        </div>
                        
                        <!-- For Group Owner -->
                        <div v-if="group.owner_id === currentUserId" class="space-y-3">
                            <div class="text-xs text-blue-700 mb-3">
                                As the group owner, you have the secret key. Share it with other members so they can join the group.
                            </div>
                            <button
                                @click="showKeyModal = true"
                                class="w-full bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200 text-sm flex items-center justify-center gap-2"
                            >
                                <span class="material-icons text-sm">content_copy</span>
                                Copy Secret Key
                            </button>
                        </div>
                        
                        <!-- For Group Members -->
                        <div v-else class="space-y-3">
                            <div class="text-xs text-blue-700 mb-3">
                                Ask the group owner for the secret key to join this encrypted group.
                            </div>
                            <button
                                @click="showKeyModal = true"
                                class="w-full bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200 text-sm flex items-center justify-center gap-2"
                            >
                                <span class="material-icons text-sm">key</span>
                                Enter Secret Key
                            </button>
                        </div>
                    </div>

                <!-- Warning -->
                <div class="bg-orange-50 rounded-xl p-4 border border-orange-200">
                    <div class="flex items-center gap-2 mb-2">
                        <span class="material-icons text-orange-600 text-sm">warning</span>
                        <span class="font-semibold text-orange-800 text-sm">Important</span>
                    </div>
                    <div class="text-xs text-orange-700 leading-relaxed">
                        Keep your device secure and don't share your private keys. If you leave the group, 
                        your keys will be automatically deleted for security.
                        </div>
                    </div>
                </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-3 mt-8 px-8 pb-8 flex-shrink-0">
                <button
                    class="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 font-medium py-3 px-4 rounded-xl transition-colors duration-200"
                    @click="$emit('close')"
                >
                    Close
                </button>
                <button
                    v-if="group.owner_id === currentUserId"
                    class="flex-1 bg-purple-500 hover:bg-purple-600 text-white font-medium py-3 px-4 rounded-xl transition-all duration-200"
                    @click="handleEditGroup"
                >
                    <span class="material-icons text-sm mr-2">edit</span>
                    Edit Group
                </button>
                <button
                    v-else
                    class="flex-1 bg-red-500 hover:bg-red-600 text-white font-medium py-3 px-4 rounded-xl transition-all duration-200"
                    @click="handleLeaveGroup"
                >
                    <span class="material-icons text-sm mr-2">exit_to_app</span>
                    Leave Group
                </button>
            </div>
            
            <!-- Debug Button (only in development) -->
            <div v-if="isDevelopment" class="mt-4 px-8 pb-4 flex-shrink-0">
                <button
                    class="w-full bg-red-100 hover:bg-red-200 text-red-700 font-medium py-2 px-4 rounded-lg transition-colors duration-200 text-sm"
                    @click="clearGroupKeys"
                >
                    <span class="material-icons text-sm mr-2">delete_forever</span>
                    Clear Group Keys (Debug)
                </button>
            </div>
        </div>
    </div>
    
    <!-- Secret Group Key Modal -->
    <SecretGroupKeyModal
        :is-visible="showKeyModal"
        :group-id="group?.id"
        :is-group-owner="group?.owner_id === currentUserId"
        @close="showKeyModal = false"
        @key-entered="handleKeyEntered"
    />
</template>

<script setup>
import { ref } from 'vue';
import { useKeyPair } from "../../composables/useKeyPair";
import { useSecretGroupE2EE } from "../../composables/useSecretGroupE2EE";
import { showMessage, showError } from "../../utils/toast";
import SecretGroupKeyModal from "../chat/SecretGroupKeyModal.vue";

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
        type: [String, Number],
        required: true,
    },
});

const emit = defineEmits(["close", "edit-group", "leave-group"]);

const isDevelopment = import.meta.env.DEV;
const showKeyModal = ref(false);

const { clearSecretGroupKeys } = useKeyPair();
const { clearGroupSecretKey } = useSecretGroupE2EE();

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
        showError('Failed to copy invite link');
    }
};

const handleEditGroup = () => {
    emit('edit-group', props.group);
};

const handleLeaveGroup = () => {
    emit('leave-group', props.group);
};

const handleKeyEntered = () => {
    showMessage('Secret key entered successfully! You can now read encrypted messages.');
};

const clearGroupKeys = async () => {
    try {
        await clearSecretGroupKeys(props.group.id);
        await clearGroupSecretKey(props.group.id);
        showMessage("Group keys cleared successfully!");
    } catch (error) {
        showError("Failed to clear group keys. Please try again.");
    }
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style> 