<template>
    <div v-if="isVisible" class="fixed inset-0 z-50 flex items-center justify-center">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-gray-500 bg-opacity-30 backdrop-blur-sm"></div>
        
        <!-- Modal -->
        <div class="relative bg-white rounded-2xl shadow-2xl max-w-md w-full mx-4 transform transition-all duration-300 scale-100 hover:shadow-3xl">
            <!-- Header -->
            <div class="flex items-center justify-between p-6 border-b border-gray-100">
                <div class="flex items-center space-x-3">
                    <div class="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                        <svg class="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                        </svg>
                    </div>
                    <div>
                        <h3 class="text-lg font-semibold text-gray-900">
                            Submit Approval Request
                        </h3>
                        <p class="text-sm text-gray-500">
                            Request permission to join groups
                        </p>
                    </div>
                </div>
                <button 
                    @click="closeModal"
                    class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
                >
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                    </svg>
                </button>
            </div>
            
            <!-- Content -->
            <div class="p-6">
                <!-- Info Box -->
                <div class="mb-6 p-4 bg-orange-50 rounded-xl border border-orange-200">
                    <div class="flex items-start space-x-3">
                        <div class="flex-shrink-0">
                            <svg class="w-5 h-5 text-orange-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                            </svg>
                        </div>
                        <div>
                            <h4 class="font-medium text-orange-800">
                                Approval Required
                            </h4>
                            <p class="text-sm mt-1 text-orange-700">
                                You need admin approval to join groups. Please provide a reason for your request.
                            </p>
                        </div>
                    </div>
                </div>
                
                <!-- Reason Field -->
                <div class="mb-6">
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                        Reason for Approval Request *
                    </label>
                    <textarea
                        v-model="approvalReason"
                        rows="4"
                        placeholder="Please explain why you want to join groups..."
                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-colors duration-200 resize-none"
                        :disabled="isSubmitting"
                    ></textarea>
                    <p class="text-xs text-gray-500 mt-1">
                        Be specific about your intentions and how you plan to use the group features.
                    </p>
                </div>
            </div>
            
            <!-- Footer -->
            <div class="flex items-center justify-end space-x-3 p-6 border-t border-gray-100">
                <button
                    @click="closeModal"
                    class="px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg font-medium transition-colors duration-200 cursor-pointer"
                    :disabled="isSubmitting"
                >
                    Cancel
                </button>
                <button
                    @click="handleSubmitApproval"
                    :disabled="isSubmitting || !approvalReason.trim()"
                    class="px-4 py-2 text-white font-medium rounded-lg transition-all duration-200 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed bg-orange-600 hover:bg-orange-700 focus:ring-2 focus:ring-orange-500 focus:ring-offset-2"
                >
                    <span v-if="isSubmitting" class="flex items-center space-x-2">
                        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Submitting...
                    </span>
                    <span v-else>
                        Submit Request
                    </span>
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { showInfo, showError } from '../../utils/toast';
import axiosInstance from '../../axiosInstance';

// Props
const props = defineProps({
    isVisible: {
        type: Boolean,
        default: false
    },
    inviteLink: {
        type: String,
        default: null
    }
});

// Emits
const emit = defineEmits(['close', 'approval-submitted']);

// Reactive data
const approvalReason = ref('');
const isSubmitting = ref(false);

// Methods
const closeModal = () => {
    approvalReason.value = '';
    emit('close');
};

const handleSubmitApproval = async () => {
    if (!approvalReason.value.trim()) {
        showError('Please provide a reason for your approval request');
        return;
    }
    
    if (!props.inviteLink) {
        showError('Invite link is required for approval submission');
        return;
    }
    
    isSubmitting.value = true;
    
    try {
        // Check if this is a secret group invite
        // Secret group invite links are typically longer UUIDs
        // Regular group invite links might be shorter or have a different format
        const isSecretGroup = props.inviteLink.length >= 36; // UUID length
        
        // Clean the invite link and ensure proper encoding
        const cleanInviteLink = props.inviteLink.trim();
        const encodedInviteLink = encodeURIComponent(cleanInviteLink);
        
        console.log('🔐 Submitting approval for:', isSecretGroup ? 'secret group' : 'regular group');
        console.log('🔐 Original invite link:', cleanInviteLink);
        console.log('🔐 Encoded invite link:', encodedInviteLink);
        console.log('🔐 Request payload:', { reason: approvalReason.value.trim() });
        
        let response;
        let endpoint = `/api/approvals/submit/${encodedInviteLink}`;
        if (isSecretGroup) {
            endpoint += '?is_secret=true';
        }
        
        console.log('🔐 Endpoint:', endpoint);
        console.log('🔐 Full request details:', {
            url: endpoint,
            method: 'POST',
            payload: { reason: approvalReason.value.trim() },
            reasonLength: approvalReason.value.trim().length
        });
        
        try {
            // Try with simple JSON payload first
            const payload = { reason: approvalReason.value.trim() };
            console.log('🔐 Sending JSON payload:', payload);
            
            response = await axiosInstance.post(endpoint, payload);
        } catch (firstError) {
            console.log('🔐 JSON attempt failed, trying FormData...');
            console.log('🔐 First error:', firstError);
            
            // Try with FormData
            const formData = new FormData();
            formData.append('reason', approvalReason.value.trim());
            
            console.log('🔐 FormData payload:', formData);
            response = await axiosInstance.post(endpoint, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            });
        }
        
        console.log('Approval submission response:', response.data);
        
        const groupType = isSecretGroup ? 'secret group' : 'group';
        showInfo(`Your approval request for the ${groupType} has been submitted. Please wait for admin approval.`);
        closeModal();
    } catch (error) {
        console.error('Failed to submit approval:', error);
        console.error('Error details:', {
            message: error.message,
            code: error.code,
            status: error.response?.status,
            statusText: error.response?.statusText,
            data: error.response?.data,
            url: error.config?.url
        });
        
        let errorMessage = 'Failed to submit approval request. Please try again.';
        
        if (error.code === 'ERR_NETWORK') {
            errorMessage = 'Network error. Please check your connection and try again.';
        } else if (error.response?.status === 404) {
            errorMessage = 'Invalid invite link. Please check the link and try again.';
        } else if (error.response?.status === 400) {
            if (error.response?.data?.detail?.includes('parseJson')) {
                errorMessage = 'Backend parsing error. Please try again or contact support.';
            } else {
                errorMessage = error.response?.data?.message || 'Invalid request. Please check your input and try again.';
            }
        } else if (error.response?.status === 500) {
            errorMessage = 'Server error. Please try again later.';
        }
        
        showError(errorMessage);
    } finally {
        isSubmitting.value = false;
    }
};
</script> 