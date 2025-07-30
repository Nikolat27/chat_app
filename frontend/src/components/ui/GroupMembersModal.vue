<template>
    <div
        v-if="isVisible && group"
        class="fixed inset-0 bg-gray-300 bg-opacity-30 backdrop-blur-md flex items-center justify-center z-50 p-4"
    >
        <div
            class="bg-white rounded-3xl shadow-2xl border border-green-100 p-8 w-96 max-w-[90vw] max-h-[90vh] overflow-y-auto transform transition-all duration-300 scale-100 hover:shadow-3xl"
        >
            <!-- Header -->
            <div class="flex items-center justify-between mb-6">
                <h2 class="text-2xl font-bold text-gray-800 flex items-center">
                    <span class="material-icons text-green-600 mr-2"
                        >group</span
                    >
                    Group Users
                </h2>
                <button
                    @click="closeModal"
                    class="text-gray-400 hover:text-gray-600 transition-colors cursor-pointer"
                >
                    <span class="material-icons text-2xl">close</span>
                </button>
            </div>

            <!-- Group Info -->
            <div class="mb-6 p-4 bg-gray-50 rounded-xl">
                <div class="flex items-center gap-3">
                    <img
                        :src="
                            group.avatar_url
                                ? `${backendBaseUrl}/static/${group.avatar_url}`
                                : '/src/assets/default-avatar.jpg'
                        "
                        :alt="group.name"
                        class="w-12 h-12 rounded-full object-cover border-2 border-green-300"
                    />
                    <div>
                        <h3 class="font-semibold text-gray-800">
                            {{ group.name }}
                        </h3>
                        <p class="text-sm text-gray-600">
                            {{ group.member_count || 0 }} members
                        </p>
                    </div>
                </div>
            </div>

            <!-- Loading State -->
            <div v-if="isLoading" class="text-center py-8">
                <div
                    class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-500 mx-auto mb-3"
                ></div>
                <p class="text-gray-600">Loading members...</p>
            </div>

            <!-- Members List -->
            <div v-else-if="members.length > 0" class="space-y-3">
                <h4 class="text-sm font-semibold text-gray-700 mb-3">
                    Members ({{ members.length }})
                </h4>

                <div class="space-y-2 max-h-96 overflow-y-auto">
                    <div
                        v-for="member in members"
                        :key="member.user_id"
                        class="flex items-center justify-between p-3 bg-gray-50 rounded-xl transition-colors"
                        :class="
                            member.user_id === currentUserId
                                ? 'bg-gray-700'
                                : 'bg-gray-50 hover:bg-gray-100'
                        "
                    >
                        <!-- Member Info -->
                        <div class="flex items-center gap-3 flex-1">
                            <img
                                :src="getAvatarUrl(member.avatar_url)"
                                :alt="member.username"
                                class="w-10 h-10 rounded-full object-cover border-2 border-white shadow-sm"
                            />
                            <div class="flex-1 min-w-0">
                                <div class="flex items-center gap-2">
                                    <span
                                        class="truncate"
                                        :class="
                                            member.user_id === currentUserId
                                                ? 'text-white font-semibold text-lg'
                                                : 'text-gray-800'
                                        "
                                    >
                                        {{ member.username }}
                                    </span>
                                    <template
                                        v-if="
                                            member.user_id === group.owner_id &&
                                            member.is_admin
                                        "
                                        ><span
                                            class="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded-full font-medium"
                                        >
                                            Owner
                                        </span>
                                        <span
                                            class="px-2 py-1 bg-green-100 text-green-700 text-xs rounded-full font-medium"
                                        >
                                            Admin
                                        </span>
                                    </template>
                                    <span
                                        v-else-if="member.is_admin"
                                        class="px-2 py-1 bg-green-100 text-green-700 text-xs rounded-full font-medium"
                                    >
                                        Admin
                                    </span>
                                    <span
                                        v-if="member.is_banned"
                                        class="px-2 py-1 bg-red-100 text-red-700 text-xs rounded-full font-medium"
                                    >
                                        Banned
                                    </span>
                                    <span
                                        v-else-if="
                                            !member.is_admin &&
                                            !member.is_banned
                                        "
                                        class="px-2 py-1 bg-yellow-100 text-yellow-700 text-xs rounded-full font-medium"
                                    >
                                        Regular User
                                    </span>
                                </div>
                                <p
                                    class="text-xs"
                                    :class="
                                        member.user_id === currentUserId
                                            ? 'text-white'
                                            : 'text-gray-500'
                                    "
                                >
                                    {{ member.is_admin || "Member" }}
                                </p>
                            </div>
                        </div>

                        <!-- Action Buttons -->
                        <div class="flex items-center gap-2">
                            <!-- Start Chat Button (only show if not current user) -->
                            <button
                                v-if="member.user_id !== currentUserId"
                                @click="handleStartChat(member)"
                                class="w-8 h-8 text-blue-500 hover:bg-blue-50 rounded-full hover:text-blue-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 flex items-center justify-center"
                                title="Start chat with user"
                            >
                                <span class="material-icons text-sm">chat</span>
                            </button>

                            <!-- Ban Button (only show if user is not banned) -->
                            <button
                                v-if="
                                    canManageMember(member) && !member.is_banned
                                "
                                @click="handleBan(member)"
                                :disabled="isBanning === member.user_id"
                                class="px-3 py-1 rounded-lg text-xs font-medium transition-all duration-200 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1 bg-red-500 hover:bg-red-600 text-white"
                                title="Ban user"
                            >
                                <span
                                    v-if="isBanning === member.user_id"
                                    class="material-icons animate-spin text-xs"
                                    >refresh</span
                                >
                                <span v-else class="material-icons text-xs"
                                    >block</span
                                >
                                <span>Ban</span>
                            </button>

                            <!-- Unban Button (only show if user is banned) -->
                            <button
                                v-if="
                                    canManageMember(member) && member.is_banned
                                "
                                @click="handleUnban(member)"
                                :disabled="isBanning === member.user_id"
                                class="px-3 py-1 rounded-lg text-xs font-medium transition-all duration-200 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1 bg-green-500 hover:bg-green-600 text-white"
                                title="Unban user"
                            >
                                <span
                                    v-if="isBanning === member.user_id"
                                    class="material-icons animate-spin text-xs"
                                    >refresh</span
                                >
                                <span v-else class="material-icons text-xs"
                                    >person_add</span
                                >
                                <span>Unban</span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Empty State -->
            <div v-else class="text-center py-8">
                <div class="text-gray-500">
                    <span class="material-icons text-6xl mb-4 block"
                        >group_off</span
                    >
                    <p class="text-lg font-medium">No members found</p>
                    <p class="text-sm">This group has no members</p>
                </div>
            </div>

            <!-- Error State -->
            <div
                v-if="error"
                class="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg"
            >
                <p class="text-red-700 text-sm">{{ error }}</p>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import { showMessage, showError } from "../../utils/toast";
import axiosInstance from "../../axiosInstance";

const props = defineProps({
    isVisible: {
        type: Boolean,
        default: false,
    },
    group: {
        type: Object,
        default: null,
    },
    backendBaseUrl: {
        type: String,
        required: true,
    },
    currentUserId: {
        type: String,
        required: true,
    },
});

const emit = defineEmits(["close", "member-updated", "start-chat"]);

const isLoading = ref(false);
const isBanning = ref(null);
const members = ref([]);
const error = ref("");

const getAvatarUrl = (avatarUrl) => {
    if (!avatarUrl) return "/src/assets/default-avatar.jpg";
    if (avatarUrl.startsWith("http")) return avatarUrl;
    return `${props.backendBaseUrl}/static/${avatarUrl}`;
};

const canManageMember = (member) => {
    // Only admins (including owner) can ban/unban members, and admins can't ban themselves or owners
    const isCurrentUserAdmin =
        props.group.admins.includes(props.currentUserId) ||
        props.currentUserId === props.group.owner_id;
    const isCurrentUser = member.user_id === props.currentUserId;
    const isMemberOwner = member.user_id === props.group.owner_id;

    console.log("🔍 Checking permissions:", {
        groupOwnerId: props.group.owner_id,
        memberUserId: member.user_id,
        currentUserId: props.currentUserId,
        isCurrentUserAdmin,
        isCurrentUser,
        isMemberOwner,
        admins: props.group.admins,
    });

    // Can manage if current user is admin/owner AND not trying to ban themselves AND not trying to ban the owner
    return isCurrentUserAdmin && !isCurrentUser && !isMemberOwner;
};

const loadMembers = async () => {
    if (!props.group?.id) return;

    try {
        isLoading.value = true;
        error.value = "";

        console.log("👥 Loading group members for group:", props.group.id);
        
        // Check if this is a secret group
        const isSecretGroup = props.group.type === 'secret';
        const endpoint = isSecretGroup 
            ? `/api/secret-group/get/${props.group.id}/members`
            : `/api/group/get/${props.group.id}/members`;
            
        console.log("👥 Making API call to:", endpoint);
        const response = await axiosInstance.get(endpoint);
        console.log("👥 Group members response:", response.data);

        // Transform the response to include user details
        let membersData = response.data.members || response.data || {};
        // Get banned members from the group details instead of the members endpoint
        console.log(props.group)

        const bannedMembers = props.group.banned_members || [];
        const admins = props.group.admins || [];
        
        console.log("🚫 Banned members from group details:", bannedMembers);
        console.log("👑 Admins from group:", admins);

        // Handle object format where keys are user IDs
        if (typeof membersData === "object" && !Array.isArray(membersData)) {
            console.log("🔄 Converting object format to array:", membersData);

            const membersArray = Object.entries(membersData).map(
                ([user_id, userData]) => {
                    const member = {
                        user_id: user_id,
                        username: userData.username,
                        avatar_url: userData.avatar_url,
                        is_banned: bannedMembers.includes(user_id),
                        is_admin:
                            admins.includes(user_id) ||
                            user_id === props.group.owner_id,
                        current_user_id: userData.current_user_id, // For permission checking
                    };
                    console.log("👤 Member processed:", {
                        user_id: member.user_id,
                        username: member.username,
                        is_admin: member.is_admin,
                        is_owner: member.user_id === props.group.owner_id,
                        is_banned: member.is_banned,
                    });
                    return member;
                }
            );
            
            // Handle banned members - they might be in the banned_members array but not in active members
            const bannedMembersNotInActive = bannedMembers.filter(
                bannedId => !Object.keys(membersData).includes(bannedId)
            );
            
            console.log("🚫 Banned members not in active list:", bannedMembersNotInActive);
            
            if (bannedMembersNotInActive.length > 0) {
                // Try to get user details for banned members
                try {
                    // Try to get user details for each banned member
                    const bannedMembersWithDetails = await Promise.all(
                        bannedMembersNotInActive.map(async (bannedId) => {
                            try {
                                // Try to get user details from user API
                                const userResponse = await axiosInstance.get(`/api/user/get/${bannedId}`);
                                console.log(`👤 User details for banned member ${bannedId}:`, userResponse.data);
                                
                                return {
                                    user_id: bannedId,
                                    username: userResponse.data.username || `User ${bannedId.slice(-6)}`,
                                    avatar_url: userResponse.data.avatar_url || null,
                                    is_banned: true,
                                    is_admin: admins.includes(bannedId) || bannedId === props.group.owner_id,
                                    current_user_id: null,
                                };
                            } catch (userError) {
                                console.error(`❌ Failed to get user details for ${bannedId}:`, userError);
                                // Fallback to placeholder
                                return {
                                    user_id: bannedId,
                                    username: `User ${bannedId.slice(-6)}`,
                                    avatar_url: null,
                                    is_banned: true,
                                    is_admin: admins.includes(bannedId) || bannedId === props.group.owner_id,
                                    current_user_id: null,
                                };
                            }
                        })
                    );
                    
                    // Combine active and banned members
                    members.value = [...membersArray, ...bannedMembersWithDetails];
                } catch (bannedError) {
                    console.error("❌ Failed to load banned members details:", bannedError);
                    
                    // Create placeholder entries for banned members if we can't get their details
                    const bannedMembersArray = bannedMembersNotInActive.map(bannedId => ({
                        user_id: bannedId,
                        username: `User ${bannedId.slice(-6)}`, // Use last 6 chars of ID as username
                        avatar_url: null,
                        is_banned: true,
                        is_admin: admins.includes(bannedId) || bannedId === props.group.owner_id,
                        current_user_id: null,
                    }));
                    
                    // Combine active and banned members
                    members.value = [...membersArray, ...bannedMembersArray];
                }
            } else {
                members.value = membersArray;
            }
        } else if (Array.isArray(membersData)) {
            // Handle array format
            const activeMembers = membersData.map((member) => ({
                user_id: member.user_id,
                username: member.username,
                avatar_url: member.avatar_url,
                is_banned: bannedMembers.includes(member.user_id),
                is_admin:
                    admins.includes(member.user_id) ||
                    member.user_id === props.group.owner_id,
                current_user_id: member.current_user_id, // For permission checking
            }));
            
            // Handle banned members not in the active list
            const bannedMembersNotInActive = bannedMembers.filter(
                bannedId => !membersData.some(member => member.user_id === bannedId)
            );
            
            if (bannedMembersNotInActive.length > 0) {
                // Try to get user details for each banned member
                try {
                    const bannedMembersWithDetails = await Promise.all(
                        bannedMembersNotInActive.map(async (bannedId) => {
                            try {
                                // Try to get user details from user API
                                const userResponse = await axiosInstance.get(`/api/user/get/${bannedId}`);
                                console.log(`👤 User details for banned member ${bannedId}:`, userResponse.data);
                                
                                return {
                                    user_id: bannedId,
                                    username: userResponse.data.username || `User ${bannedId.slice(-6)}`,
                                    avatar_url: userResponse.data.avatar_url || null,
                                    is_banned: true,
                                    is_admin: admins.includes(bannedId) || bannedId === props.group.owner_id,
                                    current_user_id: null,
                                };
                            } catch (userError) {
                                console.error(`❌ Failed to get user details for ${bannedId}:`, userError);
                                // Fallback to placeholder
                                return {
                                    user_id: bannedId,
                                    username: `User ${bannedId.slice(-6)}`,
                                    avatar_url: null,
                                    is_banned: true,
                                    is_admin: admins.includes(bannedId) || bannedId === props.group.owner_id,
                                    current_user_id: null,
                                };
                            }
                        })
                    );
                    
                    // Combine active and banned members
                    members.value = [...activeMembers, ...bannedMembersWithDetails];
                } catch (bannedError) {
                    console.error("❌ Failed to load banned members details:", bannedError);
                    
                    // Create placeholder entries for banned members if we can't get their details
                    const bannedMembersArray = bannedMembersNotInActive.map(bannedId => ({
                        user_id: bannedId,
                        username: `User ${bannedId.slice(-6)}`, // Use last 6 chars of ID as username
                        avatar_url: null,
                        is_banned: true,
                        is_admin: admins.includes(bannedId) || bannedId === props.group.owner_id,
                        current_user_id: null,
                    }));
                    
                    // Combine active and banned members
                    members.value = [...activeMembers, ...bannedMembersArray];
                }
            } else {
                members.value = activeMembers;
            }
        } else {
            console.log("⚠️ Unexpected membersData format:", membersData);
            members.value = [];
        }

        console.log("✅ Loaded", members.value.length, "group members (including banned)");
    } catch (error) {
        console.error("❌ Failed to load group members:", error);
        error.value = "Failed to load group members. Please try again.";
    } finally {
        isLoading.value = false;
    }
};

const handleBan = async (member) => {
    if (!props.group?.id) return;

    try {
        isBanning.value = member.user_id;
        console.log("🔄 Banning user:", member.user_id);

        const response = await axiosInstance.post(
            `/api/group/ban/${props.group.id}`,
            {
                target_user: member.user_id,
            }
        );

        console.log("✅ User banned successfully:", response.data);
        showMessage("User banned successfully!");

        // Update the member's banned status
        member.is_banned = true;

        // Emit event to parent
        emit("member-updated", { member, action: "ban" });
    } catch (error) {
        console.error("❌ Failed to ban user:", error);
        showError(
            error.response?.data?.message ||
                "Failed to ban user. Please try again."
        );
    } finally {
        isBanning.value = null;
    }
};

const handleUnban = async (member) => {
    if (!props.group?.id) return;

    try {
        isBanning.value = member.user_id;
        console.log("🔄 Unbanning user:", member.user_id);

        const response = await axiosInstance.post(
            `/api/group/unban/${props.group.id}`,
            {
                target_user: member.user_id,
            }
        );

        console.log("✅ User unbanned successfully:", response.data);
        showMessage("User unbanned successfully!");

        // Update the member's banned status
        member.is_banned = false;

        // Emit event to parent
        emit("member-updated", { member, action: "unban" });
    } catch (error) {
        console.error("❌ Failed to unban user:", error);
        showError(
            error.response?.data?.message ||
                "Failed to unban user. Please try again."
        );
    } finally {
        isBanning.value = null;
    }
};

const handleStartChat = (member) => {
    console.log("💬 Starting chat with user:", member);
    emit("start-chat", {
        user_id: member.user_id,
        username: member.username,
        avatar_url: member.avatar_url
    });
    closeModal();
};

const closeModal = () => {
    emit("close");
};

// Load members when modal becomes visible
watch(
    () => props.isVisible,
    (newValue) => {
        if (newValue && props.group) {
            loadMembers();
        } else {
            // Reset state when modal closes
            members.value = [];
            error.value = "";
            isLoading.value = false;
            isBanning.value = null;
        }
    }
);

// Also load members when group changes
watch(
    () => props.group,
    (newGroup) => {
        if (props.isVisible && newGroup) {
            loadMembers();
        }
    }
);
</script>
