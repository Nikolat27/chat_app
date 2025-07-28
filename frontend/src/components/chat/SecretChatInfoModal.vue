<template>
    <div
        v-if="isVisible"
        class="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
        @click.self="$emit('close')"
    >
        <div class="bg-white rounded-3xl shadow-2xl p-8 w-full max-w-md relative font-sans">
            <!-- Close Button -->
            <button
                class="absolute top-4 right-4 text-gray-400 hover:text-purple-700 hover:bg-purple-50 w-10 h-10 rounded-full transition-all duration-200 cursor-pointer flex items-center justify-center"
                @click="$emit('close')"
                aria-label="Close"
            >
                <span class="material-icons text-xl">close</span>
            </button>

            <!-- Header -->
            <div class="text-center mb-8">
                <div class="mb-4">
                    <span class="material-icons text-5xl text-purple-500 mb-3">security</span>
                </div>
                <h3 class="text-2xl font-bold text-gray-800 mb-2">Secret Chat Info</h3>
                <p class="text-sm text-gray-600">
                    End-to-end encrypted conversation details
                </p>
            </div>

            <!-- Chat Info -->
            <div class="space-y-6">
                <!-- User Info -->
                <div class="bg-gray-50 rounded-xl p-4">
                    <div class="flex items-center gap-3 mb-3">
                        <img
                            v-if="user?.avatar_url"
                            :src="`${backendBaseUrl}/static/${user.avatar_url}`"
                            class="w-10 h-10 rounded-full object-cover border-2 border-purple-300 select-none pointer-events-none"
                            alt="Avatar"
                        />
                        <img
                            v-else
                            src="/src/assets/default-avatar.jpg"
                            class="w-10 h-10 rounded-full object-cover border-2 border-purple-300 select-none pointer-events-none"
                            alt="Default Avatar"
                        />
                        <div>
                            <span class="font-semibold text-gray-800">{{ user?.username }}</span>
                            <div class="text-xs text-purple-600 flex items-center gap-1">
                                <span class="material-icons text-xs">verified</span>
                                Secret chat participant
                            </div>
                        </div>
                    </div>
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
                                    Messages are encrypted with unique key pairs and can only be read by you and the other participant.
                                </div>
                            </div>
                        </div>
                        
                        <div class="flex items-start gap-3 p-3 bg-blue-50 rounded-lg border border-blue-200">
                            <span class="material-icons text-blue-600 text-sm mt-0.5">storage</span>
                            <div>
                                <div class="font-medium text-blue-800 text-sm">No Server Storage</div>
                                <div class="text-xs text-blue-700 mt-1">
                                    Encrypted messages are not stored on our servers, ensuring maximum privacy.
                                </div>
                            </div>
                        </div>
                        
                        <div class="flex items-start gap-3 p-3 bg-purple-50 rounded-lg border border-purple-200">
                            <span class="material-icons text-purple-600 text-sm mt-0.5">key</span>
                            <div>
                                <div class="font-medium text-purple-800 text-sm">Unique Key Pairs</div>
                                <div class="text-xs text-purple-700 mt-1">
                                    Each secret chat uses a unique ECDH key pair for enhanced security.
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Chat Status -->
                <div class="bg-purple-50 rounded-xl p-4 border border-purple-200">
                    <div class="flex items-center gap-2 mb-2">
                        <span class="material-icons text-purple-600 text-sm">info</span>
                        <span class="font-semibold text-purple-800 text-sm">Chat Status</span>
                    </div>
                    <div class="text-xs text-purple-700 space-y-1">
                        <div class="flex items-center gap-2">
                            <span class="material-icons text-xs">check_circle</span>
                            <span>Chat is active and secure</span>
                        </div>
                        <div class="flex items-center gap-2">
                            <span class="material-icons text-xs">verified</span>
                            <span>Both participants approved</span>
                        </div>
                        <div class="flex items-center gap-2">
                            <span class="material-icons text-xs">security</span>
                            <span>Encryption keys established</span>
                        </div>
                    </div>
                </div>

                <!-- Warning -->
                <div class="bg-orange-50 rounded-xl p-4 border border-orange-200">
                    <div class="flex items-center gap-2 mb-2">
                        <span class="material-icons text-orange-600 text-sm">warning</span>
                        <span class="font-semibold text-orange-800 text-sm">Important</span>
                    </div>
                    <div class="text-xs text-orange-700 leading-relaxed">
                        Keep your device secure and don't share your private keys. If you lose access to your device, 
                        you may lose access to your secret chat messages.
                    </div>
                </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-3 mt-8">
                <button
                    class="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 font-medium py-3 px-4 rounded-xl transition-colors duration-200"
                    @click="$emit('close')"
                >
                    Close
                </button>
                <button
                    class="flex-1 bg-gradient-to-r from-purple-500 to-pink-500 text-white font-medium py-3 px-4 rounded-xl hover:from-purple-600 hover:to-pink-600 transition-all duration-200"
                    @click="exportKeys"
                >
                    <span class="material-icons text-sm mr-2">download</span>
                    Export Keys
                </button>
            </div>
            
            <!-- Debug Button (only in development) -->
            <div v-if="isDevelopment" class="mt-4">
                <button
                    class="w-full bg-red-100 hover:bg-red-200 text-red-700 font-medium py-2 px-4 rounded-lg transition-colors duration-200 text-sm"
                    @click="clearAllKeys"
                >
                    <span class="material-icons text-sm mr-2">delete_forever</span>
                    Clear All Keys (Debug)
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { useKeyPair } from "../../composables/useKeyPair";
import { showMessage, showError } from "../../utils/toast";

const props = defineProps({
    isVisible: {
        type: Boolean,
        default: false,
    },
    user: {
        type: Object,
        default: null,
    },
    backendBaseUrl: {
        type: String,
        required: true,
    },
});

const emit = defineEmits(["close"]);

const isDevelopment = import.meta.env.DEV;

const { getPrivateKey, getSecretChatPrivateKey, exportSecretChatPublicKey, clearAllKeys: clearAllKeysFromComposable } = useKeyPair();

const exportKeys = async () => {
    try {
        // For now, we'll export the general private key
        // In the future, this could be specific to the current secret chat
        const privateKey = await getPrivateKey();
        if (privateKey) {
            // Create a downloadable file
            const blob = new Blob([JSON.stringify(privateKey, null, 2)], {
                type: 'application/json'
            });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `secret-chat-keys-${Date.now()}.json`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
            
            showMessage("Keys exported successfully. Keep them safe!");
        } else {
            showError("No keys found to export.");
        }
    } catch (error) {
        showError("Failed to export keys. Please try again.");
    }
};

const clearAllKeys = async () => {
    try {
        await clearAllKeysFromComposable();
        showMessage("All keys cleared successfully!");
    } catch (error) {
        showError("Failed to clear keys. Please try again.");
    }
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style> 