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
import { useAuthCheck } from "../composables/useAuthCheck";
import axiosInstance from "../axiosInstance";
import Sidebar from "../components/Sidebar.vue";
import MiddleSection from "../components/MiddleSection.vue";
import ChatSection from "../components/ChatSection.vue";

const activeTab = ref("chats");
const userStore = useUserStore();
const router = useRouter();
const { clearAllKeys } = useKeyPair();
const { clearAllSymmetricKeys } = useE2EE();
const { checkAuth } = useAuthCheck();

onMounted(async () => {
  try {
    await checkAuth();
  } catch (error) {
    // The axios interceptor will handle the redirect for noAuthCookie errors
    console.error("Authentication check failed:", error);
  }
});

async function logout() {
  try {
    // Send logout request to backend
    await axiosInstance.get("/api/logout");
  } catch (error) {
    console.error("Error during logout:", error);
  }
  
  try {
    // Clear all E2EE keys
    await clearAllKeys();
    clearAllSymmetricKeys();
  } catch (error) {
    console.error("Error clearing E2EE keys on logout:", error);
  }
  
  userStore.clearUser();
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
