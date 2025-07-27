import { ref } from 'vue';
import { useChatStore } from '../stores/chat';
import axiosInstance from '../axiosInstance';
import { showError, showMessage } from '../utils/toast';

export function useMessageDeletion() {
    const chatStore = useChatStore();
    const isDeleting = ref(false);
    const error = ref(null);

    // Delete message from backend and update store
    const deleteMessage = async (messageId, chatId) => {
        if (isDeleting.value) return false;
        
        try {
            isDeleting.value = true;
            error.value = null;

            // Call backend API to delete message
            const response = await axiosInstance.delete(`/api/chat/message/${messageId}`, {
                params: { chat_id: chatId }
            });

            // Remove message from store
            const success = chatStore.deleteMessage(messageId);
            
            if (success) {
                showMessage('Message deleted successfully');
                return true;
            } else {
                showError('Failed to delete message from local store');
                return false;
            }

        } catch (err) {
            error.value = err.message || 'Failed to delete message';
            showError(error.value);
            console.error('Error deleting message:', err);
            return false;
        } finally {
            isDeleting.value = false;
        }
    };

    // Check if user can delete a specific message
    const canDeleteMessage = (messageId, currentUserId) => {
        return chatStore.canDeleteMessage(messageId, currentUserId);
    };

    // Update temporary message ID with real ID (called when backend confirms message)
    const updateMessageId = (tempId, realId) => {
        return chatStore.updateMessageId(tempId, realId);
    };

    return {
        // State
        isDeleting,
        error,
        
        // Methods
        deleteMessage,
        canDeleteMessage,
        updateMessageId,
    };
} 