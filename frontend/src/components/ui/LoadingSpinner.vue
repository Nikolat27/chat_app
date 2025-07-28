<template>
    <div class="flex items-center justify-center">
        <div 
            class="animate-spin rounded-full border-2 border-gray-300 border-t-purple-500"
            :class="sizeClasses[size]"
        ></div>
        <span v-if="text" class="ml-3 text-sm text-gray-600">{{ text }}</span>
    </div>
</template>

<script setup>
const props = defineProps({
    size: {
        type: String,
        default: 'md',
        validator: (value) => ['sm', 'md', 'lg'].includes(value)
    },
    text: {
        type: String,
        default: ''
    }
});

const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8'
};
</script> 