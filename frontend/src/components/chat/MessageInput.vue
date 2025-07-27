<template>
    <div class="p-4 border-t bg-white flex gap-2">
        <input
            :value="modelValue"
            type="text"
            placeholder="Type a message..."
            class="flex-1 border rounded px-3 py-2 focus:outline-none"
            @input="$emit('update:modelValue', $event.target.value)"
            @keyup.enter="handleSend"
        />
        <button
            @click="handleSend"
            class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
        >
            Send
        </button>
    </div>
</template>

<script setup>
defineProps({
    modelValue: {
        type: String,
        default: ""
    }
});

const emit = defineEmits(['update:modelValue', 'send']);

const handleSend = () => {
    emit('send');
};
</script> 