<template>
  <aside class="w-20 bg-white border-r flex flex-col items-center py-4">
    <button class="mb-6 p-2 rounded hover:bg-green-100 cursor-pointer" @click="$emit('changeTab', 'chats')" :class="{ 'bg-green-200': activeTab === 'chats' }">
      <span class="material-icons">chat</span>
    </button>
    <button class="mb-6 p-2 rounded hover:bg-green-100 cursor-pointer" @click="$emit('changeTab', 'groups')" :class="{ 'bg-green-200': activeTab === 'groups' }">
      <span class="material-icons">group</span>
    </button>
    <button class="mb-6 p-2 rounded hover:bg-green-100 cursor-pointer" @click="$emit('changeTab', 'settings')" :class="{ 'bg-green-200': activeTab === 'settings' }">
      <span class="material-icons">settings</span>
    </button>
    <button v-if="userStore.token" class="mt-auto p-2 rounded hover:bg-green-100 cursor-pointer" @click="$emit('logout')">
      <span class="material-icons">logout</span>
    </button>
  </aside>
</template>

<script setup>
import { useUserStore } from '../stores/users';
import { defineProps } from 'vue';
const props = defineProps({ activeTab: String });
const userStore = useUserStore();
</script>

<style scoped>
@import url('https://fonts.googleapis.com/icon?family=Material+Icons');
</style>
