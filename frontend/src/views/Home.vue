<template>
  <div class="flex h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50">
    <Sidebar :activeTab="activeTab" @changeTab="activeTab = $event" @logout="logout" />
    <MiddleSection :activeTab="activeTab" @switch-to-chats-tab="handleSwitchToChatsTab" />
    <ChatSection />
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useUserStore } from "../stores/users";
import { useRouter } from "vue-router";
import { useKeyPair } from "../composables/useKeyPair";
import { useE2EE } from "../composables/useE2EE";
import Sidebar from "../components/Sidebar.vue";
import MiddleSection from "../components/MiddleSection.vue";
import ChatSection from "../components/ChatSection.vue";

const activeTab = ref("chats");
const userStore = useUserStore();
const router = useRouter();
const { clearAllKeys } = useKeyPair();
const { clearAllSymmetricKeys } = useE2EE();

onMounted(() => {
  if (!userStore.token || userStore.isTokenExpired()) {
    router.replace("/auth");
  }
});

async function logout() {
  try {
    // Clear all E2EE keys
    await clearAllKeys();
    clearAllSymmetricKeys();
  } catch (error) {
    console.error("Error clearing E2EE keys on logout:", error);
  }
  
  userStore.$reset();
  router.replace("/auth");
}

function handleSwitchToChatsTab(user) {
  console.log("💬 Switching to chats tab for user:", user);
  activeTab.value = "chats";
}
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style>
