<template>
    <!-- Don't render anything if secret chat is not approved -->
    <div v-if="isSecretChat && !isSecretChatApproved" class="p-4 border-t border-gray-200 bg-white">
        <div class="w-full p-3 bg-orange-50 border border-orange-200 rounded-lg">
            <div class="flex items-center gap-2 text-orange-700">
                <span class="material-icons text-sm">warning</span>
                <span class="text-sm font-medium">Secret chat not approved yet</span>
            </div>
            <p class="text-xs text-orange-600 mt-1">
                You can only send messages after the other user approves this secret chat.
            </p>
        </div>
    </div>
    
    <!-- Show message input only if secret chat is approved or it's a regular chat -->
    <div v-else class="p-4 border-t border-gray-200 bg-white flex gap-2">
        <input
            :value="modelValue"
            type="text"
            placeholder="Type a message..."
            class="flex-1 border border-gray-200 rounded-lg px-4 py-2 text-[16px] font-medium focus:outline-none focus:ring-1 focus:ring-gray-200 transition"
            @input="$emit('update:modelValue', $event.target.value)"
            @keyup.enter="handleSend"
        />
        <button
            @click="handleSend"
            class="bg-green-500 hover:bg-green-600 text-white font-semibold w-10 h-10 rounded-full shadow-sm transition cursor-pointer flex items-center justify-center"
        >
            <span class="material-icons text-lg">send</span>
        </button>
    </div>
</template>

<script setup>
const props = defineProps({
    modelValue: {
        type: String,
        default: "",
    },
    isSecretChat: {
        type: Boolean,
        default: false,
    },
    isSecretChatApproved: {
        type: Boolean,
        default: true,
    },
});

const emit = defineEmits(["update:modelValue", "send"]);

const handleSend = () => {
    emit("send");
};
</script>
