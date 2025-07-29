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
});

const emit = defineEmits(["update:modelValue", "send"]);

const handleSend = () => {
    emit("send");
};
</script> 