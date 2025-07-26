<template>
    <div class="flex flex-col items-center justify-center h-full">
        <h2 class="text-xl font-bold text-green-600 mb-4">Upload Avatar</h2>
        <label class="mb-4 w-full">
            <span class="block mb-2 text-green-700 font-semibold"
                >Browse Image</span
            >
            <input
                type="file"
                accept="image/*"
                @change="onFileChange"
                class="w-full border-2 border-green-400 rounded px-2 py-2 focus:outline-none focus:ring-2 focus:ring-green-500 cursor-pointer bg-green-50 hover:bg-green-100 transition"
            />
        </label>
        <div class="mt-4">
            <img
                :src="userStore.avatar_url ? `${backendBaseUrl}/static/${userStore.avatar_url}` : defaultAvatar"
                alt="Avatar Preview"
                class="w-24 h-24 rounded-full object-cover border"
            />
        </div>
        <button
            v-if="avatarUrl"
            type="button"
            class="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 cursor-pointer mt-4"
            @click="handleUpload"
        >
            Submit
        </button>
        <div v-if="error" class="text-red-600 mt-2">{{ error }}</div>
    </div>
</template>
<script setup>
import axiosInstance from "@/axiosInstance";
import { useUserStore } from "@/stores/users";
import { ref, onMounted } from "vue";
import { showMessage, showError } from "@/utils/toast";
import defaultAvatar from '../assets/default-avatar.jpg'

const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;

const file = ref(null);
const error = ref("");
const avatarUrl = ref("");

function onFileChange(e) {
    error.value = "";
    const selected = e.target.files[0];
    if (!selected) return;
    const validTypes = [
        "image/jpeg",
        "image/png",
        "image/gif",
        "image/webp",
        "image/jpg",
    ];
    if (!validTypes.includes(selected.type)) {
        error.value = "Only image formats (JPG, PNG, GIF, WEBP) are supported.";
        file.value = null;
        avatarUrl.value = "";
        return;
    }
    if (selected.size > 5 * 1024 * 1024) {
        error.value = "File size must be less than 5MB.";
        file.value = null;
        avatarUrl.value = "";
        return;
    }
    file.value = selected;
    avatarUrl.value = URL.createObjectURL(selected);
}

const userStore = useUserStore();

onMounted(() => {
    if (userStore.avatar_url) {
        avatarUrl.value = userStore.avatar_url;
    }
});

function handleUpload() {
    if (!file.value) {
        error.value = "Please select an image.";
        return;
    }

    const formData = new FormData();
    formData.append("file", file.value);

    axiosInstance
        .post("/api/user/upload-avatar", formData, {
            headers: {
                "Content-Type": "multipart/form-data",
            },
        })
        .then((resp) => {
            error.value = "";
            resp.data.avatar_url
                ? (avatarUrl.value = `${backendBaseUrl}/static/${resp.data.avatar_url}`)
                : (avatarUrl.value = "");

            file.value = null; // Clear the file input after successful upload
            showMessage("Avatar uploaded successfully!");
            userStore.avatar_url = resp.data.avatar_url;
        })
        .catch((err) => {
            showError("Upload failed. Please try again.");
            error.value = "";
            console.error(err);
        });
}
</script>

<style scoped></style>
