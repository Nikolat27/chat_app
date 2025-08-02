import axios from "axios";
import router from "./router";

const axiosInstance = axios.create({
    baseURL: import.meta.env.VITE_BACKEND_BASE_URL,
    timeout: 10000,
    withCredentials: true,
    headers: {
        "Content-Type": "application/json",
    },
});

// Response interceptor to handle authentication errors
axiosInstance.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        // Handle 401 Unauthorized or noAuthCookie errors
        if (error.response?.status === 401 || error.response?.data?.error === "noAuthCookie") {
            // Redirect to auth page if unauthorized or no auth cookie is found
            router.replace("/auth");
        }
        return Promise.reject(error);
    }
);

export default axiosInstance;
