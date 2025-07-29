import { defineStore } from 'pinia';
import axiosInstance from '../axiosInstance';

export const useGroupStore = defineStore('groups', {
    state: () => ({
        groups: [],
        isLoading: false,
        currentGroup: null,
        groupMembers: [],
        groupMessages: []
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
                const response = await axiosInstance.get('/api/user/get-groups');
                this.groups = response.data.groups || response.data || [];
                return this.groups;
            } catch (error) {
                console.error('Failed to load user groups:', error);
                // If endpoint doesn't exist yet, return empty array
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
                    id: Date.now(), // Temporary ID for frontend
                    name: formData.get('name'),
                    description: formData.get('description'),
                    type: formData.get('group_type'),
                    avatar_url: responseData.avatar_url || '',
                    invite_link: responseData.invite_link || '',
                    message: responseData.message || '',
                    created_at: new Date().toISOString(),
                    member_count: 1,
                    role: 'admin'
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
                const response = await axiosInstance.post('/api/group/create-secret', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
                
                console.log('Secret group creation response:', response.data);
                
                // Handle your backend response format
                const responseData = response.data;
                
                // Create a group object from the response
                const newGroup = {
                    id: Date.now(), // Temporary ID for frontend
                    name: formData.get('name'),
                    description: formData.get('description'),
                    type: 'secret',
                    avatar_url: responseData.avatar_url || '',
                    invite_link: responseData.invite_link || '',
                    message: responseData.message || '',
                    created_at: new Date().toISOString(),
                    member_count: 1,
                    role: 'admin',
                    is_secret: true
                };
                
                // Ensure groups is an array
                if (!Array.isArray(this.groups)) {
                    this.groups = [];
                }
                
                // Add the new group to the list
                this.groups.push(newGroup);
                
                return newGroup;
            } catch (error) {
                console.error('Failed to create secret group:', error);
                throw error;
            }
        },

        async joinGroup(joinData) {
            try {
                const response = await axiosInstance.post('/api/group/join', joinData);
                const joinedGroup = response.data.group || response.data;
                
                // Check if group already exists in our list
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

        async leaveGroup(groupId) {
            try {
                await axiosInstance.post(`/api/group/leave/${groupId}`);
                
                // Remove group from local state
                this.groups = this.groups.filter(g => g.id !== groupId);
                
                // If this was the current group, clear it
                if (this.currentGroup && this.currentGroup.id === groupId) {
                    this.currentGroup = null;
                    this.groupMembers = [];
                    this.groupMessages = [];
                }
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
                    return;
                }
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

        async updateGroup(groupId, updateData) {
            try {
                const response = await axiosInstance.put(`/api/group/${groupId}`, updateData);
                const updatedGroup = response.data.group;
                
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

        async deleteGroup(groupId) {
            try {
                await axiosInstance.delete(`/api/group/${groupId}`);
                
                // Remove from local state
                this.groups = this.groups.filter(g => g.id !== groupId);
                
                // If this was the current group, clear it
                if (this.currentGroup && this.currentGroup.id === groupId) {
                    this.currentGroup = null;
                    this.groupMembers = [];
                    this.groupMessages = [];
                }
            } catch (error) {
                console.error('Failed to delete group:', error);
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
        }
    }
}); 