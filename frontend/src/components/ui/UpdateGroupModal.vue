<template>
    <div v-if="isVisible && group" class="fixed inset-0 bg-gray-500 bg-opacity-30 backdrop-blur-md flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-3xl shadow-2xl border border-green-100 p-8 w-96 max-w-[90vw] max-h-[90vh] overflow-y-auto transform transition-all duration-300 scale-100 hover:shadow-3xl">
            <!-- Header -->
            <div class="flex items-center justify-between mb-6">
                <h2 class="text-2xl font-bold text-gray-800 flex items-center">
                    <span class="material-icons text-green-600 mr-2">edit</span>
                    Update Group
                </h2>
                <button
                    @click="closeModal"
                    class="text-gray-400 hover:text-gray-600 transition-colors cursor-pointer"
                >
                    <span class="material-icons text-2xl">close</span>
                </button>
            </div>

            <!-- Form -->
            <form @submit.prevent="handleSubmit" class="space-y-6">
                <!-- Group Avatar Upload -->
                <div>
                    <label class="block mb-3 text-green-700 font-semibold text-sm">
                        <span class="material-icons align-middle mr-1 text-base">photo_camera</span>
                        Group Avatar (Optional)
                    </label>
                    <div class="flex items-center gap-4">
                        <div class="relative">
                            <img
                                :src="avatarPreview || (group.avatar_url ? `${backendBaseUrl}/static/${group.avatar_url}` : '/src/assets/default-avatar.jpg')"
                                alt="Group Avatar Preview"
                                class="w-16 h-16 rounded-full object-cover border-2 border-green-300 shadow-sm"
                            />
                            <div class="absolute -bottom-1 -right-1 w-4 h-4 bg-green-500 rounded-full border-2 border-white flex items-center justify-center">
                                <span class="material-icons text-white text-xs">group</span>
                            </div>
                        </div>
                        <div class="flex-1">
                            <input
                                ref="avatarInput"
                                type="file"
                                accept=".jpg,.jpeg,.png,.webp"
                                @change="handleAvatarUpload"
                                class="hidden"
                            />
                            <button
                                type="button"
                                @click="$refs.avatarInput.click()"
                                class="w-full px-4 py-2 border-2 border-green-200 rounded-xl hover:border-green-300 transition-all duration-200 text-green-700 font-medium cursor-pointer"
                            >
                                <span class="material-icons align-middle mr-1 text-sm">upload</span>
                                Choose Avatar
                            </button>
                        </div>
                    </div>
                    <p class="text-xs text-gray-500 mt-2">Leave empty to keep current avatar</p>
                </div>

                <!-- Group Name -->
                <div>
                    <label class="block mb-3 text-green-700 font-semibold text-sm">
                        <span class="material-icons align-middle mr-1 text-base">label</span>
                        Group Name
                    </label>
                    <input
                        v-model="formData.name"
                        type="text"
                        placeholder="Enter group name"
                        class="w-full px-4 py-3 border-2 border-green-200 rounded-xl focus:border-green-400 focus:outline-none transition-all duration-200 text-gray-800 font-medium"
                        :class="{ 'border-red-300': errors.name }"
                    />
                    <p v-if="errors.name" class="text-red-500 text-xs mt-1">{{ errors.name }}</p>
                </div>

                <!-- Group Description -->
                <div>
                    <label class="block mb-3 text-green-700 font-semibold text-sm">
                        <span class="material-icons align-middle mr-1 text-base">description</span>
                        Description (Optional)
                    </label>
                    <textarea
                        v-model="formData.description"
                        placeholder="Enter group description"
                        rows="3"
                        class="w-full px-4 py-3 border-2 border-green-200 rounded-xl focus:border-green-400 focus:outline-none transition-all duration-200 text-gray-800 font-medium resize-none"
                    ></textarea>
                </div>

                <!-- Group Type -->
                <div>
                    <label class="block mb-3 text-green-700 font-semibold text-sm">
                        <span class="material-icons align-middle mr-1 text-base">security</span>
                        Group Type
                    </label>
                    <select
                        v-model="formData.group_type"
                        class="w-full px-4 py-3 border-2 border-green-200 rounded-xl focus:border-green-400 focus:outline-none transition-all duration-200 text-gray-800 font-medium cursor-pointer"
                    >
                        <option value="public">Public</option>
                        <option value="private">Private</option>
                    </select>
                </div>

                <!-- Action Buttons -->
                <div class="flex gap-3 pt-4">
                    <button
                        type="button"
                        @click="closeModal"
                        class="flex-1 px-6 py-3 border-2 border-gray-300 rounded-xl hover:border-gray-400 transition-all duration-200 text-gray-700 font-semibold cursor-pointer"
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        :disabled="isLoading"
                        class="flex-1 px-6 py-3 bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white rounded-xl font-semibold transition-all duration-200 shadow-lg hover:shadow-xl cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        <span v-if="isLoading" class="material-icons animate-spin mr-2">refresh</span>
                        <span v-else class="material-icons mr-2">save</span>
                        {{ isLoading ? 'Updating...' : 'Update Group' }}
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue';
import { showMessage, showError } from '../../utils/toast';
import axiosInstance from '../../axiosInstance';

const props = defineProps({
    isVisible: {
        type: Boolean,
        default: false
    },
    group: {
        type: Object,
        default: null
    },
    backendBaseUrl: {
        type: String,
        required: true
    }
});

const emit = defineEmits(['close', 'group-updated']);

const avatarInput = ref(null);
const isLoading = ref(false);
const avatarPreview = ref(null);
const selectedFile = ref(null);

const formData = reactive({
    name: '',
    description: '',
    group_type: 'public'
});

const errors = reactive({
    name: ''
});

// Initialize form data when group changes
watch(() => props.group, (newGroup) => {
    if (newGroup) {
        formData.name = newGroup.name || '';
        formData.description = newGroup.description || '';
        formData.group_type = newGroup.type || 'public';
        avatarPreview.value = null;
        selectedFile.value = null;
        errors.name = '';
    }
}, { immediate: true });

const handleAvatarUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
        // Validate file size (max 5MB)
        if (file.size > 5 * 1024 * 1024) {
            showError('File size must be less than 5MB');
            return;
        }

        // Validate file type
        const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp'];
        if (!allowedTypes.includes(file.type)) {
            showError('Please select a valid image file (JPG, PNG, WEBP)');
            return;
        }

        selectedFile.value = file;
        
        // Create preview
        const reader = new FileReader();
        reader.onload = (e) => {
            avatarPreview.value = e.target.result;
        };
        reader.readAsDataURL(file);
    }
};

const validateForm = () => {
    errors.name = '';
    
    if (!formData.name.trim()) {
        errors.name = 'Group name is required';
        return false;
    }
    
    if (formData.name.trim().length < 3) {
        errors.name = 'Group name must be at least 3 characters';
        return false;
    }
    
    if (formData.name.trim().length > 50) {
        errors.name = 'Group name must be less than 50 characters';
        return false;
    }
    
    return true;
};

const handleSubmit = async () => {
    if (!validateForm()) return;
    
    try {
        isLoading.value = true;
        
        const formDataToSend = new FormData();
        
        // Add name (use current name if empty)
        formDataToSend.append('name', formData.name.trim() || props.group.name);
        
        // Add description (only if provided)
        if (formData.description.trim()) {
            formDataToSend.append('description', formData.description.trim());
        }
        
        // Add group type
        formDataToSend.append('group_type', formData.group_type);
        
        // Add file only if a new file was selected
        if (selectedFile.value) {
            formDataToSend.append('file', selectedFile.value);
        }
        
        console.log('📝 Updating group:', props.group.id);
        console.log('📝 Form data:', {
            name: formDataToSend.get('name'),
            description: formDataToSend.get('description'),
            group_type: formDataToSend.get('group_type'),
            hasFile: !!formDataToSend.get('file')
        });
        
        const response = await axiosInstance.put(`/api/group/update/${props.group.id}`, formDataToSend, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        
        console.log('✅ Group updated successfully:', response.data);
        showMessage('Group updated successfully!');
        
        // Emit the updated group data
        emit('group-updated', response.data);
        closeModal();
        
    } catch (error) {
        console.error('❌ Failed to update group:', error);
        showError(error.response?.data?.message || 'Failed to update group. Please try again.');
    } finally {
        isLoading.value = false;
    }
};

const closeModal = () => {
    emit('close');
};
</script> 