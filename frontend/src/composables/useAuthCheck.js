import { ref } from "vue";
import axiosInstance from "@/axiosInstance";

export function useAuthCheck() {
    const isChecking = ref(false);
    const isAuthenticated = ref(false);

    const checkAuth = async () => {
        isChecking.value = true;
        try {
            const response = await axiosInstance.get("/api/auth-check");
            isAuthenticated.value = true;
            return response.data;
        } catch (error) {
            isAuthenticated.value = false;
            // The axios interceptor will handle the redirect for noAuthCookie errors
            throw error;
        } finally {
            isChecking.value = false;
        }
    };

    return {
        isChecking,
        isAuthenticated,
        checkAuth,
    };
} 