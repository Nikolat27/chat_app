import { defineStore } from "pinia";

export const useUserStore = defineStore("user", {
    state: () => ({
        token: null,
        username: null,
        user_id: null,
        avatar_url: null,
        tokenTimestamp: null,
    }),
    actions: {
        setUser({ token, username, user_id, avatar_url }) {
            this.token = token;
            this.username = username;
            this.user_id = user_id;
            this.avatar_url = avatar_url;
            this.tokenTimestamp = Date.now();
        },
        clearUser() {
            this.token = null;
            this.username = null;
            this.user_id = null;
            this.avatar_url = null;
            this.tokenTimestamp = null;
        },
        isTokenExpired() {
            if (!this.tokenTimestamp) return true;
            // 12 hours in ms
            return Date.now() - this.tokenTimestamp > 12 * 60 * 60 * 1000;
        },
    },
    persist: true,
});
