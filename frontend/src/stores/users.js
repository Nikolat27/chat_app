import { defineStore } from 'pinia';

export const useUserStore = defineStore('user', {
  state: () => ({
    token: null,
    username: null,
    user_id: null,
  }),
  actions: {
    setUser({ token, username, user_id }) {
      this.token = token;
      this.username = username;
      this.user_id = user_id;
    },
    clearUser() {
      this.token = null;
      this.username = null;
      this.user_id = null;
    },
  },
});
