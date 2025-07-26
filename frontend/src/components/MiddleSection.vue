<template>
  <section class="w-1/3 bg-white border-r p-4 overflow-y-auto">
    <div v-if="activeTab === 'chats'">
      <h2 class="text-lg font-bold text-green-600 mb-4">Chats</h2>
      <!-- List of chats will go here -->
      <div class="text-gray-500">No chats yet.</div>
    </div>
    <div v-else-if="activeTab === 'groups'">
      <h2 class="text-lg font-bold text-green-600 mb-4">Groups</h2>
      <!-- List of groups will go here -->
      <div class="text-gray-500">No groups yet.</div>
    </div>
    <div v-else-if="activeTab === 'settings'">
      <h2 class="text-lg font-bold text-green-600 mb-4">User Settings</h2>
      <!-- User settings content will go here -->
      <div class="text-gray-500">Coming soon.</div>
    </div>
  </section>
</template>

<script setup>
import { defineProps } from 'vue';
const props = defineProps({ activeTab: String });
</script>
