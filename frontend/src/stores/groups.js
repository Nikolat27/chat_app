import { defineStore } from 'pinia';
import axiosInstance from '../axiosInstance';
import { useSecretGroupE2EE } from '../composables/useSecretGroupE2EE';
import { useKeyPair } from '../composables/useKeyPair';

export const useGroupStore = defineStore('groups', {
    state: () => ({
        groups: [],
        isLoading: false,
        currentGroup: null,
        groupMembers: [],
        groupMessages: [],
        groupUsers: {}
    }),

    getters: {
        getUserGroups: (state) => state.groups,
        getCurrentGroup: (state) => state.currentGroup,
        getGroupMembers: (state) => state.groupMembers,
        getGroupMessages: (state) => state.groupMessages,
        getGroupsCount: (state) => state.groups.length
    },

    actions: {
        async loadUserGroups() {
            try {
                this.isLoading = true;
                
                // Load regular groups
                const regularGroupsResponse = await axiosInstance.get('/api/user/get-groups');
                console.log('Regular groups response:', regularGroupsResponse.data);
                
                // Load secret groups with is_secret parameter
                const secretGroupsResponse = await axiosInstance.get('/api/user/get-groups?is_secret=true');
                console.log('Secret groups response:', secretGroupsResponse.data);
                
                // Handle different response formats for regular groups
                let regularGroupsData = regularGroupsResponse.data.groups || regularGroupsResponse.data || [];
                
                // Handle different response formats for secret groups
                let secretGroupsData = secretGroupsResponse.data.groups || secretGroupsResponse.data || [];
                
                // Ensure both are arrays
                if (!Array.isArray(regularGroupsData)) {
                    regularGroupsData = [];
                }
                if (!Array.isArray(secretGroupsData)) {
                    secretGroupsData = [];
                }
                
                // Process regular groups
                const regularGroups = regularGroupsData.map(group => ({
                    id: group.id || group._id,
                    name: group.name,
                    description: group.description,
                    type: group.type,
                    avatar_url: group.avatar_url,
                    invite_link: group.invite_link,
                    owner_id: group.owner_id,
                    created_at: group.created_at,
                    member_count: group.users?.length || group.member_count || 1,
                    role: group.role || 'member',
                    admins: group.admins || [],
                    banned_members: group.banned_members || []
                }));
                
                // Process secret groups
                const secretGroups = secretGroupsData.map(group => {
                    const processedGroup = {
                    id: group.id || group._id,
                    name: group.name,
                    description: group.description,
                    type: 'secret',
                    avatar_url: group.avatar_url,
                    invite_link: group.invite_link,
                    owner_id: group.owner_id,
                    created_at: group.created_at,
                    member_count: group.users?.length || group.member_count || 1,
                    role: group.role || 'member',
                    admins: group.admins || [],
                    banned_members: group.banned_members || [],
                    is_secret: true
                    };
                    console.log('🔍 Processing secret group:', processedGroup);
                    return processedGroup;
                });
                
                // Combine both types of groups
                this.groups = [...regularGroups, ...secretGroups];
                
                console.log('✅ Loaded', regularGroups.length, 'regular groups and', secretGroups.length, 'secret groups');
                console.log('🔍 All groups:', this.groups.map(g => ({ id: g.id, name: g.name, type: g.type, is_secret: g.is_secret })));
                
                return this.groups;
            } catch (error) {
                console.error('Failed to load user groups:', error);
                // If endpoints don't exist yet, return empty array
                if (error.response?.status === 404) {
                    this.groups = [];
                    return this.groups;
                }
                throw error;
            } finally {
                this.isLoading = false;
            }
        },

        async createGroup(formData) {
            try {
                const response = await axiosInstance.post('/api/group/create', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
                
                console.log('Group creation response:', response.data);
                
                // Handle your backend response format
                const responseData = response.data;
                
                // Create a group object from the response
                const newGroup = {
                    id: responseData.group_id || Date.now(), // Use backend group_id or fallback
                    name: formData.get('name'),
                    description: formData.get('description'),
                    type: formData.get('group_type'),
                    avatar_url: responseData.avatar_url || '',
                    invite_link: responseData.invite_link || '',
                    message: responseData.message || '',
                    created_at: new Date().toISOString(),
                    member_count: 1,
                    role: 'admin',
                    owner_id: responseData.owner_id || null // Use owner_id from backend
                };
                
                // Ensure groups is an array
                if (!Array.isArray(this.groups)) {
                    this.groups = [];
                }
                
                // Add the new group to the list
                this.groups.push(newGroup);
                
                return newGroup;
            } catch (error) {
                console.error('Failed to create group:', error);
                throw error;
            }
        },

        async createSecretGroup(formData) {
            try {
                console.log('Creating secret group with FormData:', formData);
                
                const response = await axiosInstance.post('/api/group/create?is_secret=true', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
                
                console.log('Secret group creation response:', response.data);
                
                // Extract the group ID from the response
                const groupId = response.data.group_id;
                const newGroup = {
                    id: groupId,
                    name: formData.get('name'),
                    description: formData.get('description'),
                    type: 'secret',
                    avatar_url: response.data.avatar_url,
                    invite_link: response.data.invite_link,
                    owner_id: response.data.owner_id,
                    created_at: new Date().toISOString(),
                    member_count: 1,
                    role: 'owner'
                };
                
                // Generate keys for the actual group ID
                try {
                    const { initializeSecretGroup } = useSecretGroupE2EE();
                    await initializeSecretGroup(groupId);
                    console.log('✅ Secret group encryption initialized');
                } catch (encryptionError) {
                    console.error('❌ Failed to initialize secret group encryption:', encryptionError);
                    // Don't fail the creation, just log the error
                }
                
                // Add to local state
                this.groups.push(newGroup);
                
                return newGroup;
            } catch (error) {
                console.error('Failed to create secret group:', error);
                throw error;
            }
        },

        async joinGroup(inviteLink) {
            try {
                // Use GET method with invite_link in URL path as per backend specification
                const response = await axiosInstance.get(`/api/group/join/${inviteLink}`);
                console.log('Join group response:', response.data);
                
                const joinedGroup = response.data.group || response.data;
                
                // After joining, fetch the complete group data to get name, avatar_url, etc.
                try {
                    // Load both regular and secret groups to find the newly joined group
                    const regularGroupsResponse = await axiosInstance.get('/api/user/get-groups');
                    const secretGroupsResponse = await axiosInstance.get('/api/user/get-groups?is_secret=true');
                    
                    let regularGroupsData = regularGroupsResponse.data.groups || regularGroupsResponse.data || [];
                    let secretGroupsData = secretGroupsResponse.data.groups || secretGroupsResponse.data || [];
                    
                    // Ensure both are arrays
                    if (!Array.isArray(regularGroupsData)) {
                        regularGroupsData = [];
                    }
                    if (!Array.isArray(secretGroupsData)) {
                        secretGroupsData = [];
                    }
                    
                    // Combine all groups to search for the newly joined one
                    const allGroups = [...regularGroupsData, ...secretGroupsData];
                    
                    // Find the newly joined group in the updated list
                    const updatedGroup = allGroups.find(group => 
                        group.id === joinedGroup.id || group._id === joinedGroup.id
                    );
                    
                    if (updatedGroup) {
                        // Replace the basic joined group with complete data
                        const completeGroup = {
                            id: updatedGroup.id || updatedGroup._id,
                            name: updatedGroup.name,
                            description: updatedGroup.description,
                            type: updatedGroup.type,
                            avatar_url: updatedGroup.avatar_url,
                            invite_link: updatedGroup.invite_link,
                            owner_id: updatedGroup.owner_id,
                            created_at: updatedGroup.created_at,
                            member_count: updatedGroup.users?.length || updatedGroup.member_count || 1,
                            role: updatedGroup.role || 'member'
                        };
                        
                        // Check if group already exists in our list
                        const existingIndex = this.groups.findIndex(g => g.id === completeGroup.id);
                        if (existingIndex === -1) {
                            this.groups.push(completeGroup);
                        } else {
                            // Update existing group with complete data
                            this.groups[existingIndex] = completeGroup;
                        }
                        
                        // Initialize encryption for secret groups
                        if (completeGroup.type === 'secret') {
                            try {
                                const { initializeSecretGroup } = useSecretGroupE2EE();
                                await initializeSecretGroup(completeGroup.id);
                                console.log('✅ Secret group encryption initialized after joining');
                            } catch (encryptionError) {
                                console.error('❌ Failed to initialize secret group encryption after joining:', encryptionError);
                                // Don't fail the join, just log the error
                            }
                        }
                        
                        return completeGroup;
                    }
                } catch (detailsError) {
                    console.error('Failed to fetch complete group details:', detailsError);
                    // Fallback to the basic joined group data
                }
                
                // Fallback: use the basic joined group data
                const existingIndex = this.groups.findIndex(g => g.id === joinedGroup.id);
                if (existingIndex === -1) {
                    this.groups.push(joinedGroup);
                }
                
                return joinedGroup;
            } catch (error) {
                console.error('Failed to join group:', error);
                throw error;
            }
        },

        async joinSecretGroup(inviteLink) {
            try {
                console.log('🔐 joinSecretGroup called with inviteLink:', inviteLink);
                
                console.log('🔐 Making GET request to join secret group');
                // Join secret group with GET method and invite_link in URL path
                const response = await axiosInstance.get(`/api/group/join/${inviteLink}?is_secret=true`);
                console.log('Join secret group response:', response.data);
                
                const joinedGroup = response.data.group || response.data;
                
                // After joining, fetch the complete group data to get name, avatar_url, etc.
                try {
                    // Load both regular and secret groups to find the newly joined group
                    const regularGroupsResponse = await axiosInstance.get('/api/user/get-groups');
                    const secretGroupsResponse = await axiosInstance.get('/api/user/get-groups?is_secret=true');
                    
                    let regularGroupsData = regularGroupsResponse.data.groups || regularGroupsResponse.data || [];
                    let secretGroupsData = secretGroupsResponse.data.groups || secretGroupsResponse.data || [];
                    
                    // Ensure both are arrays
                    if (!Array.isArray(regularGroupsData)) {
                        regularGroupsData = [];
                    }
                    if (!Array.isArray(secretGroupsData)) {
                        secretGroupsData = [];
                    }
                    
                    // Combine all groups to search for the newly joined one
                    const allGroups = [...regularGroupsData, ...secretGroupsData];
                    
                    // Find the newly joined group in the updated list
                    const updatedGroup = allGroups.find(group => 
                        group.id === joinedGroup.id || group._id === joinedGroup.id
                    );
                    
                    if (updatedGroup) {
                        // Replace the basic joined group with complete data
                        const completeGroup = {
                            id: updatedGroup.id || updatedGroup._id,
                            name: updatedGroup.name,
                            description: updatedGroup.description,
                            type: 'secret',
                            avatar_url: updatedGroup.avatar_url,
                            invite_link: updatedGroup.invite_link,
                            owner_id: updatedGroup.owner_id,
                            created_at: updatedGroup.created_at,
                            member_count: updatedGroup.users?.length || updatedGroup.member_count || 1,
                            role: updatedGroup.role || 'member',
                            is_secret: true
                        };
                        
                        // Check if group already exists in our list
                        const existingIndex = this.groups.findIndex(g => g.id === completeGroup.id);
                        if (existingIndex === -1) {
                            this.groups.push(completeGroup);
                        } else {
                            // Update existing group with complete data
                            this.groups[existingIndex] = completeGroup;
                        }
                        
                        // Initialize encryption for secret groups
                        try {
                            const { initializeSecretGroup } = useSecretGroupE2EE();
                            await initializeSecretGroup(completeGroup.id);
                            console.log('✅ Secret group encryption initialized after joining');
                        } catch (encryptionError) {
                            console.error('❌ Failed to initialize secret group encryption after joining:', encryptionError);
                            // Don't fail the join, just log the error
                        }
                        
                        return completeGroup;
                    }
                } catch (detailsError) {
                    console.error('Failed to fetch complete group details:', detailsError);
                    // Fallback to the basic joined group data
                }
                
                // Fallback: use the basic joined group data
                const existingIndex = this.groups.findIndex(g => g.id === joinedGroup.id);
                if (existingIndex === -1) {
                    this.groups.push(joinedGroup);
                }
                
                return joinedGroup;
            } catch (error) {
                console.error('Failed to join secret group:', error);
                throw error;
            }
        },

        async leaveGroup(groupId) {
            try {
                const response = await axiosInstance.delete(`/api/group/leave/${groupId}`);
                console.log('Left group response:', response.data);
                
                // Remove group from local state
                this.groups = this.groups.filter(g => g.id !== groupId);
                
                // If this was the current group, clear it
                if (this.currentGroup && this.currentGroup.id === groupId) {
                    this.currentGroup = null;
                    this.groupMembers = [];
                    this.groupMessages = [];
                }
                
                // Clear encryption keys for secret groups
                const { clearSecretGroupKeys } = useKeyPair();
                const { clearGroupSecretKey } = useSecretGroupE2EE();
                
                try {
                    clearGroupSecretKey(groupId);
                    await clearSecretGroupKeys(groupId);
                    console.log('✅ Cleared encryption keys for left group');
                } catch (keyError) {
                    console.error('❌ Failed to clear encryption keys:', keyError);
                }
                
                return response.data;
            } catch (error) {
                console.error('Failed to leave group:', error);
                // If endpoint doesn't exist yet, still remove from local state
                if (error.response?.status === 404) {
                    this.groups = this.groups.filter(g => g.id !== groupId);
                    if (this.currentGroup && this.currentGroup.id === groupId) {
                        this.currentGroup = null;
                        this.groupMembers = [];
                        this.groupMessages = [];
                    }
                    return { message: 'Left group successfully' };
                }
                throw error;
            }
        },

        async leaveSecretGroup(groupId) {
            try {
                const response = await axiosInstance.delete(`/api/group/leave/${groupId}?is_secret=true`);
                console.log('Leave secret group response:', response.data);
                
                // Remove from local state
                this.groups = this.groups.filter(g => g.id !== groupId);
                
                // If this was the current group, clear it
                if (this.currentGroup && this.currentGroup.id === groupId) {
                    this.currentGroup = null;
                    this.groupMembers = [];
                    this.groupMessages = [];
                }
                
                // Clear encryption keys for secret groups
                const { clearGroupSecretKey } = useSecretGroupE2EE();
                const { clearSecretGroupKeys } = useKeyPair();
                
                try {
                    clearGroupSecretKey(groupId);
                    await clearSecretGroupKeys(groupId);
                    console.log('✅ Cleared encryption keys for left secret group');
                } catch (keyError) {
                    console.error('❌ Failed to clear encryption keys:', keyError);
                }
                
                return response.data;
            } catch (error) {
                console.error('Failed to leave secret group:', error);
                throw error;
            }
        },

        async loadGroupDetails(groupId) {
            try {
                const response = await axiosInstance.get(`/api/group/${groupId}`);
                this.currentGroup = response.data.group;
                this.groupMembers = response.data.group.members || [];
                return this.currentGroup;
            } catch (error) {
                console.error('Failed to load group details:', error);
                throw error;
            }
        },

        async loadSecretGroupDetails(groupId) {
            try {
                const response = await axiosInstance.get(`/api/group/get/${groupId}/members?is_secret=true`);
                this.currentGroup = response.data.group;
                this.groupMembers = response.data.members || [];
                return this.currentGroup;
            } catch (error) {
                console.error('Failed to load secret group details:', error);
                throw error;
            }
        },

        async searchGroups(query) {
            try {
                const response = await axiosInstance.get(`/api/group/search?q=${query}`);
                return response.data.groups || [];
            } catch (error) {
                console.error('Failed to search groups:', error);
                // If endpoint doesn't exist yet, return empty array
                if (error.response?.status === 404) {
                    return [];
                }
                throw error;
            }
        },

        async updateGroup(groupId, formData) {
            try {
                const response = await axiosInstance.put(`/api/group/update/${groupId}`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
                
                console.log('Update group response:', response.data);
                
                const updatedGroup = response.data.group || response.data;
                
                // Update in local state
                const index = this.groups.findIndex(g => g.id === groupId);
                if (index !== -1) {
                    this.groups[index] = { ...this.groups[index], ...updatedGroup };
                }
                
                if (this.currentGroup && this.currentGroup.id === groupId) {
                    this.currentGroup = { ...this.currentGroup, ...updatedGroup };
                }
                
                return updatedGroup;
            } catch (error) {
                console.error('Failed to update group:', error);
                throw error;
            }
        },

        async updateSecretGroup(groupId, formData) {
            try {
                const response = await axiosInstance.put(`/api/group/update/${groupId}?is_secret=true`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
                
                console.log('Update secret group response:', response.data);
                
                const updatedGroup = response.data.group || response.data;
                
                // Update in local state
                const index = this.groups.findIndex(g => g.id === groupId);
                if (index !== -1) {
                    this.groups[index] = { ...this.groups[index], ...updatedGroup };
                }
                
                if (this.currentGroup && this.currentGroup.id === groupId) {
                    this.currentGroup = { ...this.currentGroup, ...updatedGroup };
                }
                
                return updatedGroup;
            } catch (error) {
                console.error('Failed to update secret group:', error);
                throw error;
            }
        },

        async deleteGroup(groupId) {
            try {
                const response = await axiosInstance.delete(`/api/group/delete/${groupId}`);
                console.log('Delete group response:', response.data);
                
                // Remove from local state
                this.groups = this.groups.filter(g => g.id !== groupId);
                
                // If this was the current group, clear it
                if (this.currentGroup && this.currentGroup.id === groupId) {
                    this.currentGroup = null;
                    this.groupMembers = [];
                    this.groupMessages = [];
                }
                
                // Clear encryption keys for secret groups
                const { clearGroupSecretKey } = useSecretGroupE2EE();
                const { clearSecretGroupKeys } = useKeyPair();
                
                try {
                    clearGroupSecretKey(groupId);
                    await clearSecretGroupKeys(groupId);
                    console.log('✅ Cleared encryption keys for deleted group');
                } catch (keyError) {
                    console.error('❌ Failed to clear encryption keys:', keyError);
                }
                
                return response.data;
            } catch (error) {
                console.error('Failed to delete group:', error);
                // If endpoint doesn't exist yet, still remove from local state
                if (error.response?.status === 404) {
                    this.groups = this.groups.filter(g => g.id !== groupId);
                    if (this.currentGroup && this.currentGroup.id === groupId) {
                        this.currentGroup = null;
                        this.groupMembers = [];
                        this.groupMessages = [];
                    }
                    return { message: 'Group deleted successfully' };
                }
                throw error;
            }
        },

        async deleteSecretGroup(groupId) {
            try {
                const response = await axiosInstance.delete(`/api/group/delete/${groupId}?is_secret=true`);
                console.log('Delete secret group response:', response.data);
                
                // Remove from local state
                this.groups = this.groups.filter(g => g.id !== groupId);
                
                // If this was the current group, clear it
                if (this.currentGroup && this.currentGroup.id === groupId) {
                    this.currentGroup = null;
                    this.groupMembers = [];
                    this.groupMessages = [];
                }
                
                // Clear encryption keys for secret groups
                const { clearGroupSecretKey } = useSecretGroupE2EE();
                const { clearSecretGroupKeys } = useKeyPair();
                
                try {
                    clearGroupSecretKey(groupId);
                    await clearSecretGroupKeys(groupId);
                    console.log('✅ Cleared encryption keys for deleted secret group');
                } catch (keyError) {
                    console.error('❌ Failed to clear encryption keys:', keyError);
                }
                
                return response.data;
            } catch (error) {
                console.error('Failed to delete secret group:', error);
                throw error;
            }
        },

        async addMember(groupId, userId) {
            try {
                await axiosInstance.post(`/api/group/${groupId}/add-member`, { user_id: userId });
                
                // Refresh group details
                await this.loadGroupDetails(groupId);
            } catch (error) {
                console.error('Failed to add member:', error);
                throw error;
            }
        },

        async removeMember(groupId, userId) {
            try {
                await axiosInstance.post(`/api/group/${groupId}/remove-member`, { user_id: userId });
                
                // Refresh group details
                await this.loadGroupDetails(groupId);
            } catch (error) {
                console.error('Failed to remove member:', error);
                throw error;
            }
        },

        async removeUserFromSecretGroup(groupId, userId) {
            try {
                const response = await axiosInstance.delete(`/api/group/remove-user/${groupId}/${userId}?is_secret=true`);
                console.log('Remove user from secret group response:', response.data);
                return response.data;
            } catch (error) {
                console.error('Failed to remove user from secret group:', error);
                throw error;
            }
        },

        async banMemberFromSecretGroup(groupId, userId) {
            try {
                const response = await axiosInstance.post(`/api/group/ban/${groupId}?is_secret=true`, { user_id: userId });
                console.log('Ban member from secret group response:', response.data);
                return response.data;
            } catch (error) {
                console.error('Failed to ban member from secret group:', error);
                throw error;
            }
        },

        async unbanMemberFromSecretGroup(groupId, userId) {
            try {
                const response = await axiosInstance.post(`/api/group/unban/${groupId}?is_secret=true`, { user_id: userId });
                console.log('Unban member from secret group response:', response.data);
                return response.data;
            } catch (error) {
                console.error('Failed to unban member from secret group:', error);
                throw error;
            }
        },

        async changeMemberRole(groupId, userId, role) {
            try {
                await axiosInstance.post(`/api/group/${groupId}/change-role`, { 
                    user_id: userId, 
                    role: role 
                });
                
                // Refresh group details
                await this.loadGroupDetails(groupId);
            } catch (error) {
                console.error('Failed to change member role:', error);
                throw error;
            }
        },

        setCurrentGroup(group) {
            this.currentGroup = group;
        },

        clearCurrentGroup() {
            this.currentGroup = null;
            this.groupMembers = [];
            this.groupMessages = [];
        },

        addGroupMessage(message) {
            this.groupMessages.push(message);
        },

        clearGroupMessages() {
            this.groupMessages = [];
        },

        setGroupUsers(users) {
            this.groupUsers = users;
        },

        getGroupUsers() {
            return this.groupUsers;
        }
    }
}); 