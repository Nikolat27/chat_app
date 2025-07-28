<template>
  <div class="flex h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50">
    <Sidebar :activeTab="activeTab" @changeTab="activeTab = $event" @logout="logout" />
    <MiddleSection :activeTab="activeTab" />
    <ChatSection />
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useUserStore } from "../stores/users";
import { useRouter } from "vue-router";
import Sidebar from "../components/Sidebar.vue";
import MiddleSection from "../components/MiddleSection.vue";
import ChatSection from "../components/ChatSection.vue";

const activeTab = ref("chats");
const userStore = useUserStore();
const router = useRouter();

onMounted(() => {
  if (!userStore.token || userStore.isTokenExpired()) {
    router.replace("/auth");
  }
});

function logout() {
  userStore.$reset();
  router.replace("/auth");
}
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style>
