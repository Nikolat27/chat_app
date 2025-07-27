import { defineStore } from "pinia";

export const useChatStore = defineStore("chat", {
    state: () => ({
        currentChatUser: null,
        messages: [],
        chats: [],
        avatarUrls: {},
        usernames: {},
    }),
    actions: {
        setChatUser(user) {
            this.currentChatUser = user;
        },
        clearChatUser() {
            this.currentChatUser = null;
        },
        setMessages(msgs) {
            this.messages = msgs;
        },
        addMessage(msg) {
            this.messages.push(msg);
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
    },
});
