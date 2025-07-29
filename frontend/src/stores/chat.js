import { defineStore } from "pinia";

export const useChatStore = defineStore("chat", {
    state: () => ({
        currentChatUser: null,
        messages: [],
        chats: [],
        secretChats: [],
        avatarUrls: {},
        usernames: {},
        secretUsernames: {},
        // Pagination state
        currentPage: 1,
        pageLimit: 20,
        hasMoreMessages: true,
        isLoadingMessages: false,
    }),
    actions: {
        setChatUser(user) {
            // Create a clean user object to avoid leftover properties from previous chats
            this.currentChatUser = {
                id: user.id,
                username: user.username,
                avatar_url: user.avatar_url,
                // Only include secret_chat_id if it's explicitly provided
                ...(user.secret_chat_id && { secret_chat_id: user.secret_chat_id }),
                // Only include chat_id if it's explicitly provided
                ...(user.chat_id && { chat_id: user.chat_id })
            };
        },
        clearChatUser() {
            this.currentChatUser = null;
            this.messages = []; // Clear messages when clearing chat user
        },

        setChats(chats) {
            this.chats = chats;
        },
        setAvatarUrls(urls) {
            this.avatarUrls = urls;
        },
        setUsernames(names) {
            this.usernames = names;
        },
        updateUserData(userId, username, avatarUrl) {
            if (userId) {
                this.usernames[userId] = username;
                this.avatarUrls[userId] = avatarUrl;
            }
        },
        updateChatData(chatId, username, avatarUrl) {
            if (chatId) {
                this.usernames[chatId] = username;
                this.avatarUrls[chatId] = avatarUrl;
            }
        },

        setSecretChats(chats) {
            this.secretChats = chats;
        },
        setSecretUsernames(names) {
            this.secretUsernames = names;
        },

        // Remove a specific secret chat from the store
        removeSecretChat(chatId) {
            console.log('Removing secret chat from store:', chatId);
            console.log('Current secret chats:', this.secretChats);
            
            const chatIndex = this.secretChats.findIndex(chat => chat.id === chatId);
            console.log('Found chat at index:', chatIndex);
            
            if (chatIndex !== -1) {
                this.secretChats.splice(chatIndex, 1);
                console.log('Secret chat removed from store');
            } else {
                console.log('Secret chat not found in store');
            }
            
            // Remove associated username
            delete this.secretUsernames[chatId];
            console.log('Secret chat username removed');
        },

        // Pagination methods
        setMessages(messages, reset = false) {
            if (reset) {
                this.messages = messages;
                this.currentPage = 1;
                this.hasMoreMessages = true;
            } else {
                // Add older messages to the beginning (for infinite scroll)
                this.messages = [...messages, ...this.messages];
            }
        },

        clearMessages() {
            this.messages = [];
            this.currentPage = 1;
            this.hasMoreMessages = true;
            this.isLoadingMessages = false;
        },

        addMessage(msg) {
            this.messages.push(msg);
        },

        // Message deletion methods
        deleteMessage(messageId) {
            const index = this.messages.findIndex(
                (msg) => msg.id === messageId
            );
            if (index !== -1) {
                this.messages.splice(index, 1);
                return true;
            }
            return false;
        },

        // Update temporary message ID with real ID from backend
        updateMessageId(tempId, realId) {
            const message = this.messages.find((msg) => msg.id === tempId);
            if (message) {
                message.id = realId;
                return true;
            }
            return false;
        },

        // Check if message can be deleted (user's own message and has real ID)
        canDeleteMessage(messageId, currentUserId) {
            const message = this.messages.find((msg) => msg.id === messageId);
            if (!message) return false;

            // Can only delete own messages
            if (message.sender_id !== currentUserId) return false;

            // Can only delete messages with real IDs (not temp IDs)
            if (message.id && message.id.startsWith("temp-")) return false;

            // Message must have an ID to be deletable
            if (!message.id) return false;

            return true;
        },

        setPaginationState(page, hasMore, isLoading = false) {
            this.currentPage = page;
            this.hasMoreMessages = hasMore;
            this.isLoadingMessages = isLoading;
        },

        resetPagination() {
            this.currentPage = 1;
            this.hasMoreMessages = true;
            this.isLoadingMessages = false;
        },

        setLoadingState(isLoading) {
            this.isLoadingMessages = isLoading;
        },
    },
});
