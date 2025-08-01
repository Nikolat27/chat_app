<template>
    <div
        v-if="isVisible"
        class="fixed inset-0 z-50 flex items-center justify-center"
    >
        <!-- Backdrop -->
        <div
            class="absolute inset-0 bg-gray-500 bg-opacity-30 backdrop-blur-sm"
        ></div>

        <!-- Modal -->
        <div
            class="relative bg-white rounded-2xl shadow-2xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-hidden"
        >
            <!-- Header -->
            <div
                class="flex items-center justify-between p-6 border-b border-gray-100"
            >
                <div class="flex items-center space-x-3">
                    <div
                        class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center"
                    >
                        <span class="material-icons text-blue-600"
                            >pending_actions</span
                        >
                    </div>
                    <div>
                        <h3 class="text-lg font-semibold text-gray-900">
                            Pending Approvals
                        </h3>
                        <p class="text-sm text-gray-500">
                            Review requests to join your groups
                        </p>
                    </div>
                </div>
                <div class="flex items-center space-x-3">
                    <button
                        @click="closeModal"
                        class="text-gray-400 hover:text-gray-600 transition-colors"
                    >
                        <span class="material-icons">close</span>
                    </button>
                </div>
            </div>

            <!-- Content -->
            <div class="p-6 overflow-y-auto max-h-[calc(90vh-140px)]">
                <!-- Loading State -->
                <div
                    v-if="isLoading"
                    class="flex items-center justify-center py-8"
                >
                    <div
                        class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"
                    ></div>
                    <span class="ml-3 text-gray-600">Loading approvals...</span>
                </div>

                <!-- Empty State -->
                <div
                    v-else-if="!approvals || approvals.length === 0"
                    class="text-center py-8"
                >
                    <span class="material-icons text-6xl text-gray-300"
                        >check_circle</span
                    >
                    <h4 class="text-lg font-semibold text-gray-600 mt-4">
                        No Pending Approvals
                    </h4>
                    <p class="text-sm text-gray-500">
                        All approval requests have been reviewed.
                    </p>
                </div>

                <!-- Approvals List -->
                <div v-else class="space-y-4">
                    <div
                        v-for="(approval, index) in approvals"
                        :key="approval?.id || `approval-${index}`"
                        class="bg-gray-50 rounded-xl p-4 border border-gray-200"
                    >
                        <!-- Approval Header -->
                        <div class="flex items-center justify-between mb-3">
                            <div class="flex items-center space-x-3">
                                <div
                                    class="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center"
                                >
                                    <span class="material-icons text-gray-600"
                                        >person</span
                                    >
                                </div>
                                <div>
                                    <h4 class="font-medium text-gray-900">
                                        User ID: {{ approval?.requester_id || 'Unknown' }}
                                    </h4>
                                    <p class="text-sm text-gray-500">
                                        Requested
                                        {{ formatDate(approval?.created_at) }}
                                    </p>
                                </div>
                            </div>
                            <span
                                class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                                :class="{
                                    'bg-yellow-100 text-yellow-800':
                                        approval?.status === 'pending',
                                    'bg-green-100 text-green-800':
                                        approval?.status === 'approved',
                                    'bg-red-100 text-red-800':
                                        approval?.status === 'rejected',
                                }"
                            >
                                <span class="material-icons text-xs mr-1">
                                    {{
                                        approval?.status === "approved"
                                            ? "check"
                                            : approval?.status === "rejected"
                                            ? "close"
                                            : "schedule"
                                    }}
                                </span>
                                {{ approval?.status || 'pending' }}
                            </span>
                        </div>

                        <!-- Group Info -->
                        <div
                            class="mb-3 p-3 bg-white rounded-lg border border-gray-200"
                        >
                            <div class="flex items-center space-x-2">
                                <span
                                    :class="
                                        approval?.is_secret
                                            ? 'text-purple-600'
                                            : 'text-blue-600'
                                    "
                                    class="material-icons text-sm"
                                >
                                    {{ approval?.is_secret ? "lock" : "group" }}
                                </span>
                                <span class="text-sm font-medium text-gray-700">
                                    Group ID: {{ approval?.group_id || 'Unknown' }}
                                    <span
                                        v-if="approval?.is_secret"
                                        class="ml-2 inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-700"
                                    >
                                        Secret Group
                                    </span>
                                </span>
                            </div>
                        </div>

                        <!-- Reason -->
                        <div class="mb-4">
                            <h5 class="text-sm font-medium text-gray-700 mb-2">
                                Reason for Request:
                            </h5>
                            <div
                                class="bg-white rounded-lg p-3 border border-gray-200"
                            >
                                <p class="text-sm text-gray-600">
                                    {{ approval?.reason || 'No reason provided' }}
                                </p>
                            </div>
                        </div>

                        <!-- Action Buttons -->
                        <div class="flex items-center justify-end space-x-3">
                            <button
                                @click="handleRejectApproval(approval?.id)"
                                :disabled="isProcessingApproval === approval?.id"
                                class="px-4 py-2 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg font-medium transition-colors cursor-pointer disabled:opacity-50 border border-red-200"
                            >
                                <span
                                    v-if="isProcessingApproval === approval?.id"
                                    class="flex items-center space-x-2"
                                >
                                    <svg
                                        class="animate-spin h-4 w-4 text-red-600"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                    >
                                        <circle
                                            class="opacity-25"
                                            cx="12"
                                            cy="12"
                                            r="10"
                                            stroke="currentColor"
                                            stroke-width="4"
                                        ></circle>
                                        <path
                                            class="opacity-75"
                                            fill="currentColor"
                                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                        ></path>
                                    </svg>
                                    Rejecting...
                                </span>
                                <span
                                    v-else
                                    class="flex items-center space-x-2"
                                >
                                    <span class="material-icons text-sm"
                                        >close</span
                                    >
                                    Reject
                                </span>
                            </button>
                            <button
                                @click="handleApproveApproval(approval?.id)"
                                :disabled="isProcessingApproval === approval?.id"
                                class="px-4 py-2 text-white bg-green-600 hover:bg-green-700 rounded-lg font-medium transition-colors cursor-pointer disabled:opacity-50"
                            >
                                <span
                                    v-if="isProcessingApproval === approval?.id"
                                    class="flex items-center space-x-2"
                                >
                                    <svg
                                        class="animate-spin h-4 w-4 text-white"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                    >
                                        <circle
                                            class="opacity-25"
                                            cx="12"
                                            cy="12"
                                            r="10"
                                            stroke="currentColor"
                                            stroke-width="4"
                                        ></circle>
                                        <path
                                            class="opacity-75"
                                            fill="currentColor"
                                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                        ></path>
                                    </svg>
                                    Approving...
                                </span>
                                <span
                                    v-else
                                    class="flex items-center space-x-2"
                                >
                                    <span class="material-icons text-sm"
                                        >check</span
                                    >
                                    Approve
                                </span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, watch, onMounted } from "vue";
import { showInfo, showError } from "../../utils/toast";
import axiosInstance from "../../axiosInstance";

// Props
const props = defineProps({
    isVisible: {
        type: Boolean,
        default: false,
    },
});

// Emits
const emit = defineEmits(["close", "approval-updated"]);

// Reactive data
const approvals = ref([]);
const isLoading = ref(false);
const isProcessingApproval = ref(null);

// Methods
const closeModal = () => {
    emit("close");
};

watch(
    () => props.isVisible,
    (visible) => {
        if (visible) {
            // Reset approvals and load fresh data
            approvals.value = [];
            isLoading.value = true;
            loadApprovals();
        }
    },
    { immediate: true }
);

const loadApprovals = async () => {
    try {
        isLoading.value = true;
        const response = await axiosInstance.get(
            "/api/received-approvals/get/"
        );

        // Handle different response formats and filter out null/undefined approvals
        let approvalsData = [];
        
        if (response.data && response.data.approvals) {
            // Backend returned { approvals: [...] }
            approvalsData = response.data.approvals;
        } else if (Array.isArray(response.data)) {
            // Backend returned array directly
            approvalsData = response.data;
        } else if (response.data && typeof response.data === 'object') {
            // Backend returned object with other structure
            approvalsData = Object.values(response.data);
        }
        
        // Ensure it's an array and filter out invalid entries
        if (Array.isArray(approvalsData)) {
            approvals.value = approvalsData.filter(approval => 
                approval && 
                typeof approval === 'object' && 
                approval.id !== null && 
                approval.id !== undefined
            );
            
            // Add default status if missing
            approvals.value.forEach(approval => {
                if (!approval.status) {
                    approval.status = 'pending'; // Default status
                }
            });
            
            console.log("✅ Loaded approvals:", approvals.value.length);
        } else {
            console.warn("Unexpected approvals response format:", response.data);
            approvals.value = [];
        }
    } catch (error) {
        console.error("Failed to load approvals:", error);
        showError("Failed to load pending approvals. Please try again.");
        approvals.value = [];
    } finally {
        isLoading.value = false;
    }
};

const handleApproveApproval = async (approvalId) => {
    try {
        isProcessingApproval.value = approvalId;

        // Find the approval to check if it's for a secret group
        const approval = approvals.value.find((a) => a && a.id === approvalId);
        
        // Check if approval exists
        if (!approval) {
            console.error("Approval not found:", approvalId);
            showError("Approval not found. Please refresh and try again.");
            return;
        }
        
        const isSecretGroup = approval.is_secret;

        let endpoint = `/api/approvals/edit-status/${approvalId}`;
        if (isSecretGroup) {
            endpoint += "?is_secret=true";
        }

        const response = await axiosInstance.put(endpoint, {
            status: "approved",
        });
        console.log("Approve approval response:", response.data);

        const groupType = isSecretGroup ? "secret group" : "group";
        showInfo(`Approval request for ${groupType} approved successfully!`);

        // Update the approval status instead of removing it
        approval.status = "approved";

        emit("approval-updated");
    } catch (error) {
        console.error("Failed to approve approval:", error);
        showError(
            error.response?.data?.message ||
                "Failed to approve request. Please try again."
        );
    } finally {
        isProcessingApproval.value = null;
    }
};

const handleRejectApproval = async (approvalId) => {
    try {
        isProcessingApproval.value = approvalId;

        // Find the approval to check if it's for a secret group
        const approval = approvals.value.find((a) => a && a.id === approvalId);
        
        // Check if approval exists
        if (!approval) {
            console.error("Approval not found:", approvalId);
            showError("Approval not found. Please refresh and try again.");
            return;
        }
        
        const isSecretGroup = approval.is_secret;

        let endpoint = `/api/approvals/edit-status/${approvalId}`;
        if (isSecretGroup) {
            endpoint += "?is_secret=true";
        }

        const response = await axiosInstance.put(endpoint, {
            status: "rejected",
        });
        console.log("Reject approval response:", response.data);

        const groupType = isSecretGroup ? "secret group" : "group";
        showInfo(`Approval request for ${groupType} rejected successfully!`);

        // Update the approval status instead of removing it
        approval.status = "rejected";

        emit("approval-updated");
    } catch (error) {
        console.error("Failed to reject approval:", error);
        showError(
            error.response?.data?.message ||
                "Failed to reject request. Please try again."
        );
    } finally {
        isProcessingApproval.value = null;
    }
};

const formatDate = (dateString) => {
    if (!dateString) {
        return 'Unknown date';
    }
    try {
        const date = new Date(dateString);
        if (isNaN(date.getTime())) {
            return 'Invalid date';
        }
        return date.toLocaleDateString() + " " + date.toLocaleTimeString();
    } catch (error) {
        console.error("Error formatting date:", error);
        return 'Invalid date';
    }
};

// Load approvals when modal opens
// onMounted(() => {
//     if (props.isVisible) {
//         loadApprovals();
//     }
// });
</script>
