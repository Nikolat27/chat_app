<template>
    <section class="w-1/3 bg-white border-r p-4 overflow-y-auto">
        <!-- Chats Tab -->
        <div v-if="activeTab === 'chats'">
            <ChatsTab
                :chat-store="chatStore"
                :user-store="userStore"
                :backend-base-url="backendBaseUrl"
                :secret-chats="chatStore.secretChats"
                :secret-usernames="chatStore.secretUsernames"
                @open-chat="handleOpenChat"
            />
        </div>

        <!-- Groups Tab -->
        <div v-else-if="activeTab === 'groups'">
            <GroupsTab
                :backend-base-url="backendBaseUrl"
                @group-clicked="handleGroupClick"
                @switch-to-chat="handleSwitchToChat"
            />
        </div>

        <!-- Settings Tab -->
        <div v-else-if="activeTab === 'settings'">
            <SettingsTab />
        </div>
    </section>
</template>

<script setup>
import { defineProps, onMounted, ref, nextTick } from "vue";
import { useChatStore } from "../stores/chat";
import { useUserStore } from "../stores/users";
import { useGroupStore } from "../stores/groups";
import axiosInstance from "@/axiosInstance";
import { useMessagePagination } from "../composables/useMessagePagination";
import { useE2EE } from "../composables/useE2EE";
import { useSecretChatEncryption } from "../composables/useSecretChatEncryption";
import { showMessage } from "../utils/toast";
import ChatsTab from "./tabs/ChatsTab.vue";
import GroupsTab from "./tabs/GroupsTab.vue";
import SettingsTab from "./tabs/SettingsTab.vue";

const props = defineProps({ activeTab: String });
const emit = defineEmits(['switch-to-chats-tab']);
const chatStore = useChatStore();
const userStore = useUserStore();
const groupStore = useGroupStore();
const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;

const { loadInitialMessages, loadInitialSecretChatMessages } =
    useMessagePagination();
const { loadChatSymmetricKey } = useE2EE();
const { loadSecretChatSymmetricKey } = useSecretChatEncryption();
const chatsLoaded = ref(false);

onMounted(async () => {
    try {
        await fetchUserChats();
        await fetchUserSecretChats();
    } catch (error) {
        // console.error("Failed to fetch chats:", error);
    } finally {
        chatsLoaded.value = true;
    }
});

async function fetchUserChats() {
    try {
        const response = await axiosInstance.get("/api/user/get-chats");
        chatStore.setChats(response.data.chats);
        chatStore.setAvatarUrls(response.data.avatar_urls);
        chatStore.setUsernames(response.data.usernames);
    } catch (error) {
        // console.error("Failed to fetch user chats:", error);
    }
}

async function fetchUserSecretChats() {
    try {
        const response = await axiosInstance.get("/api/user/get-secret-chats");
        chatStore.setSecretChats(response.data.secret_chats);
        chatStore.setSecretUsernames(response.data.secret_usernames);
    } catch (error) {
        // console.error("Failed to fetch user secret chats:", error);
    }
}

const handleOpenChat = async (user) => {
    // Clear previous messages before switching to new chat
    chatStore.clearMessages();

    // Clear any existing group
    groupStore.clearCurrentGroup();

    chatStore.setChatUser(user);

    if (user.id && user.username) {
        chatStore.updateUserData(user.id, user.username, user.avatar_url);
    }

    // Check if this is a secret chat
    if (user.secret_chat_id) {
        console.log("🔐 Opening secret chat with ID:", user.secret_chat_id);
        // For secret chats, load the symmetric key first
        try {
            console.log("🔐 Attempting to load symmetric key...");
            const symmetricKeyLoaded = await loadSecretChatSymmetricKey(
                user.secret_chat_id
            );
            console.log("🔐 Symmetric key loading result:", symmetricKeyLoaded);
            if (!symmetricKeyLoaded) {
                console.log(
                    "⚠️ Symmetric key not available for secret chat, encryption disabled"
                );
                // Continue without E2EE if key loading fails
            } else {
                console.log("✅ Symmetric key loaded successfully");
            }

            // Check if key is in memory after loading
            const { hasSymmetricKey } = useE2EE();
            const keyInMemory = await hasSymmetricKey(user.secret_chat_id);
            console.log("🔍 Key in memory after loading:", keyInMemory);
        } catch (error) {
            console.error("❌ Error loading symmetric key:", error);
            // Continue without E2EE if key loading fails
        }

        // Load secret chat messages
        await loadInitialSecretChatMessages(user.secret_chat_id);
        return;
    }

    const existingChat = findExistingChat(user.id);

    if (existingChat) {
        await loadInitialMessages(existingChat.id);
    } else {
        await createNewChat(user);
    }
};

const findExistingChat = (targetUserId) => {
    const currentUserId = userStore.user_id;
    return chatStore.chats?.find(
        (chat) =>
            chat.participants &&
            chat.participants.includes(targetUserId) &&
            chat.participants.includes(currentUserId)
    );
};

// Helper function to add new chat to store and set up UI
const addNewChatToStore = async (newChat, user) => {
    // Force reactivity
    const updatedChats = [...chatStore.chats, newChat];
    chatStore.setChats(updatedChats);
    chatStore.updateChatData(newChat.id, user.username, user.avatar_url);
    await setupNewChat(newChat, user);
};

// Helper function to set up UI for new chat (without adding to store)
const setupNewChat = async (newChat, user) => {
    // Update chat metadata
    chatStore.updateChatData(newChat.id, user.username, user.avatar_url);
    await nextTick();
    chatStore.setChatUser({
        id: user.id,
        username: user.username,
        avatar_url: user.avatar_url,
        chat_id: newChat.id,
    });
    await loadInitialMessages(newChat.id);
};

const createNewChat = async (user) => {
    try {
        const response = await axiosInstance.post("/api/chat/create", {
            target_user: user.id,
        });

        // Check if response contains a chat object or if we need to fetch chats
        if (response.data?.chat) {
            // Backend returned the chat object directly
            const newChat = response.data.chat;
            // Add to store and continue
            await addNewChatToStore(newChat, user);
        } else if (
            response.data === "chat created successfully" ||
            response.status === 200
        ) {
            // Backend only returned success message, need to fetch the new chat
            // Re-fetch all chats to get the new one
            await fetchUserChats();
            // Find the new chat (it should be the one with the target user)
            const newChat = findExistingChat(user.id);
            if (newChat) {
                // Don't add to store again (already added by fetchUserChats)
                // Just set up the current chat user and load messages
                await setupNewChat(newChat, user);
            }
        }
    } catch (error) {
        // Optionally, show a toast error here
    }
};

const handleSwitchToChat = async (user) => {
    console.log("💬 Switching to chat with user:", user);
    
    // Emit event to parent to switch to chats tab
    emit('switch-to-chats-tab', user);
    
    // Also handle opening the chat directly
    await handleOpenChat(user);
};

const handleGroupClick = async (group) => {
    console.log("🎯 handleGroupClick called with group:", group);
    console.log("🎯 Group ID:", group.id);

    // Clear any existing chat user and group
    chatStore.clearChatUser();
    groupStore.clearCurrentGroup();

    // Set the current group
    groupStore.setCurrentGroup(group);

    // Clear messages for the new group chat
    chatStore.clearMessages();

    // Load group users and messages
    try {
        console.log("👥 Loading group users for group:", group.id);
        console.log(
            "👥 Making API call to:",
            `/api/group/get/${group.id}/members`
        );

        const response = await axiosInstance.get(
            `/api/group/get/${group.id}/members`
        );
        console.log("👥 Group users response:", response.data);

        // Store group users in the group store or a reactive variable
        // For now, we'll store it in the group store
        console.log(response.data);
        groupStore.setGroupUsers(response.data);

        console.log(
            "✅ Loaded",
            Object.keys(response.data).length,
            "group users"
        );
    } catch (error) {
        console.error("❌ Failed to load group users:", error);
        console.error("❌ Error details:", error.response?.data);
        console.error("❌ Full error:", error);
    }
};
</script>
