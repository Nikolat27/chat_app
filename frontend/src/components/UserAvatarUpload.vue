<template>
    <div class="flex flex-col items-center justify-center h-full px-6">
        <h2 class="text-2xl font-bold text-green-600 mb-6">Upload Avatar</h2>

        <!-- File Input -->
        <label class="w-full max-w-md mb-6">
            <span class="block mb-2 text-green-700 font-medium text-base">
                Browse Image
            </span>
            <input
                type="file"
                accept="image/*"
                @change="onFileChange"
                class="w-full border-2 border-gray-200 rounded-lg px-3 py-2 bg-green-50 hover:bg-green-100 text-sm transition focus:outline-none focus:ring-2 focus:ring-green-500 cursor-pointer"
            />
        </label>

        <!-- Avatar Preview -->
        <div class="mt-4">
            <img
                :src="avatarUrl || defaultAvatar"
                alt="Avatar Preview"
                class="w-28 h-28 rounded-full object-cover border border-gray-200 shadow-sm select-none pointer-events-none"
            />
        </div>

        <!-- Submit Button -->
        <button
            v-if="avatarUrl"
            type="button"
            class="w-full max-w-md bg-green-500 hover:bg-green-600 text-white font-semibold py-2 rounded-lg shadow-sm mt-6 transition duration-150 cursor-pointer"
            @click="handleUpload"
        >
            Submit
        </button>

        <!-- Error Message -->
        <div v-if="error" class="text-red-600 mt-3 text-sm font-medium">
            {{ error }}
        </div>
    </div>
</template>

<script setup>
import axiosInstance from "@/axiosInstance";
import { useUserStore } from "@/stores/users";
import { ref, onMounted } from "vue";
import { showMessage, showError } from "@/utils/toast";
import defaultAvatar from "../assets/default-avatar.jpg";

const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;

const file = ref(null);
const error = ref("");
const avatarUrl = ref("");

const userStore = useUserStore();

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
    avatarUrl.value = URL.createObjectURL(selected); // live preview
}

onMounted(() => {
    if (userStore.avatar_url) {
        avatarUrl.value = `${backendBaseUrl}/static/${userStore.avatar_url}`;
    }
});

async function handleUpload() {
    if (!file.value) {
        error.value = "Please select an image.";
        return;
    }

    const formData = new FormData();
    formData.append("file", file.value);

    try {
        const resp = await axiosInstance.post(
            "/api/user/upload-avatar",
            formData,
            {
                headers: {
                    "Content-Type": "multipart/form-data",
                },
            }
        );

        error.value = "";
        const uploadedUrl = resp.data.avatar_url;

        if (uploadedUrl) {
            avatarUrl.value = `${backendBaseUrl}/static/${uploadedUrl}`;
            userStore.avatar_url = uploadedUrl;
            showMessage("Avatar uploaded successfully!");
        } else {
            avatarUrl.value = "";
            showError("Upload succeeded but no URL returned.");
        }

        file.value = null; // Clear input after success
    } catch (err) {
        showError("Upload failed. Please try again.");
        error.value = "";
        console.error(err);
    }
}
</script>
<style scoped></style>
