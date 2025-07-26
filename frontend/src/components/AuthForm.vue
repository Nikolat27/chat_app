<template>
    <div class="min-h-screen flex items-center justify-center bg-green-50">
        <div class="w-full max-w-md p-8 bg-white rounded-lg shadow-md">
            <h2 class="text-2xl font-bold text-green-600 mb-6 text-center">
                Register
            </h2>
            <form @submit.prevent="register">
                <input
                    v-model="registerForm.username"
                    type="text"
                    placeholder="Username"
                    class="w-full mb-4 px-4 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-green-500"
                />
                <input
                    v-model="registerForm.password"
                    type="password"
                    placeholder="Password"
                    class="w-full mb-4 px-4 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-green-500"
                />
                <button
                    type="submit"
                    class="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 cursor-pointer"
                >
                    Register
                </button>
            </form>
            <div class="my-6 border-t pt-6">
                <h2 class="text-2xl font-bold text-green-600 mb-6 text-center">
                    Login
                </h2>
                <form @submit.prevent="login">
                    <input
                        v-model="loginForm.username"
                        type="text"
                        placeholder="Username"
                        class="w-full mb-4 px-4 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-green-500"
                    />
                    <input
                        v-model="loginForm.password"
                        type="password"
                        placeholder="Password"
                        class="w-full mb-4 px-4 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-green-500"
                    />
                    <button
                        type="submit"
                        class="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 cursor-pointer"
                    >
                        Login
                    </button>
                </form>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from "vue";
import axiosInstance from "../axiosInstance";
import { showMessage, showError } from "../utils/toast";
import { useUserStore } from "../stores/users";
import { useRouter } from "vue-router";

const registerForm = ref({ username: "", password: "" });
const loginForm = ref({ username: "", password: "" });
const message = ref("");
const router = useRouter();

const register = async () => {
    try {
        const res = await axiosInstance.post(
            "/api/register",
            registerForm.value
        );
        showMessage("Registration successful!");
        message.value = "";
    } catch (err) {
        showError(err.response?.data?.detail || "Registration failed.");
        message.value = "";
    }
};

const userStore = useUserStore();

const login = async () => {
    try {
        const res = await axiosInstance.post("/api/login", loginForm.value);
        userStore.setUser({
            token: res.data.token,
            username: res.data.username,
            user_id: res.data.user_id,
        });
        showMessage("Login successful!");
        message.value = "";
        router.push("/");
    } catch (err) {
        showError(err.response?.data?.detail || "Login failed.");
        message.value = "";
    }
};
</script>
