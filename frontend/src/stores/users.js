import { defineStore } from "pinia";

export const useUserStore = defineStore("user", {
    state: () => ({
        username: null,
        user_id: null,
        avatar_url: null,
    }),
    actions: {
        setUser({ username, user_id, avatar_url }) {
            this.username = username;
            this.user_id = user_id;
            this.avatar_url = avatar_url;
        },
        clearUser() {
            this.username = null;
            this.user_id = null;
            this.avatar_url = null;
        },
    },
    persist: true,
});
