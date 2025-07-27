<template>
    <div class="p-4 border-t border-gray-200 bg-white flex gap-2">
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
            class="bg-green-500 hover:bg-green-600 text-white font-semibold px-4 py-2 rounded-lg shadow-sm transition cursor-pointer"
        >
            Send
        </button>
    </div>
</template>

<script setup>
const { modelValue } = defineProps({
    modelValue: {
        type: String,
        default: "",
    },
});

const emit = defineEmits(["update:modelValue", "send"]);

const handleSend = () => {
    emit("send");
};
</script>
