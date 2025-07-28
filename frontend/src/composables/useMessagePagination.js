import { ref, computed } from "vue";
import { useChatStore } from "../stores/chat";
import axiosInstance from "../axiosInstance";

export function useMessagePagination() {
    const chatStore = useChatStore();

    // Reactive state
    const isLoading = ref(false);
    const error = ref(null);

    // Computed properties
    const currentPage = computed(() => chatStore.currentPage);
    const pageLimit = computed(() => chatStore.pageLimit);
    const hasMoreMessages = computed(() => chatStore.hasMoreMessages);
    const isLoadingMessages = computed(() => chatStore.isLoadingMessages);

    // Load messages function (you'll implement the API call)
    const loadMessages = async (chatId, page = 1, limit = 20) => {
        if (isLoading.value) return;

        try {
            isLoading.value = true;
            error.value = null;
            chatStore.setLoadingState(true);

            const response = await axiosInstance.get(
                `/api/chat/get/${chatId}/messages`,
                {
                    params: { page, limit },
                }
            );

            // Handle the response structure
            // The response is an array of JSON strings that need to be parsed
            const rawMessages = response.data || [];

            // Parse each message from JSON string to object
            const messages = rawMessages
                .map((msg) => {
                    if (typeof msg === "string") {
                        try {
                            const parsed = JSON.parse(msg);
                            return parsed;
                        } catch (e) {
                            console.error("Failed to parse message:", msg, e);
                            return null;
                        }
                    }
                    return msg; // If it's already an object, return as is
                })
                .filter((msg) => msg !== null); // Remove any failed parses

            // For pagination, we'll assume there are more messages if we got a full page
            // You might need to adjust this based on your backend's actual pagination response
            const hasMore = rawMessages.length >= limit;
            const totalPages = Math.ceil(rawMessages.length / limit) || 1;

            // Update store with new messages
            chatStore.setMessages(messages, page === 1);
            chatStore.setPaginationState(page, hasMore, false);

            return { messages, hasMore, totalPages };
        } catch (err) {
            error.value = err.message || "Failed to load messages";
            console.error("Error loading messages:", err);
        } finally {
            isLoading.value = false;
            chatStore.setLoadingState(false);
        }
    };

    // Fetch secret chats
    const loadSecretChats = async (userId, page = 1, limit = 20) => {
        if (isLoading.value) return;
        try {
            isLoading.value = true;
            error.value = null;
            chatStore.setLoadingState(true);

            // Adjust the endpoint as per your backend
            const response = await axiosInstance.get(
                `/api/user/get-secret-chats`,
                {
                    params: { user_id: userId, page, limit },
                }
            );

            console.log(response);

            const rawChats = response.data || [];
            // Assign default avatar if not present
            const chats = rawChats.map((chat) => {
                if (!chat.avatar_url) {
                    return {
                        ...chat,
                        avatar_url: "default-secret-chat-avatar.jpg",
                    };
                }
                return chat;
            });

            // You may want to update the store or return the chats directly
            // chatStore.setSecretChats(chats); // If you have such a method
            return chats;
        } catch (err) {
            error.value = err.message || "Failed to load secret chats";
            console.error("Error loading secret chats:", err);
        } finally {
            isLoading.value = false;
            chatStore.setLoadingState(false);
        }
    };

    // Load next page
    const loadNextPage = async (chatId) => {
        if (!hasMoreMessages.value || isLoadingMessages.value) return;

        const nextPage = currentPage.value + 1;
        return await loadMessages(chatId, nextPage, pageLimit.value);
    };

    // Load initial messages
    const loadInitialMessages = async (chatId) => {
        chatStore.resetPagination();
        return await loadMessages(chatId, 1, pageLimit.value);
    };

    // Load secret chat messages
    const loadSecretChatMessages = async (secretChatId, page = 1, limit = 20) => {
        if (isLoading.value) return;

        try {
            isLoading.value = true;
            error.value = null;
            chatStore.setLoadingState(true);

            const response = await axiosInstance.get(
                `/api/secret-chat/get/${secretChatId}/messages`,
                {
                    params: { page, limit },
                }
            );

            // Handle the response structure for secret chat messages
            const rawMessages = response.data || [];

            // Parse each message from JSON string to object
            const messages = rawMessages
                .map((msg) => {
                    if (typeof msg === "string") {
                        try {
                            const parsed = JSON.parse(msg);
                            return parsed;
                        } catch (e) {
                            console.error("Failed to parse secret chat message:", msg, e);
                            return null;
                        }
                    }
                    return msg; // If it's already an object, return as is
                })
                .filter((msg) => msg !== null);

            // For pagination, we'll assume there are more messages if we got a full page
            const hasMore = rawMessages.length >= limit;
            const totalPages = Math.ceil(rawMessages.length / limit) || 1;

            // Update store with new messages
            chatStore.setMessages(messages, page === 1);
            chatStore.setPaginationState(page, hasMore, false);

            return { messages, hasMore, totalPages };
        } catch (err) {
            error.value = err.message || "Failed to load secret chat messages";
            console.error("Error loading secret chat messages:", err);
        } finally {
            isLoading.value = false;
            chatStore.setLoadingState(false);
        }
    };

    // Load initial secret chat messages
    const loadInitialSecretChatMessages = async (secretChatId) => {
        chatStore.resetPagination();
        return await loadSecretChatMessages(secretChatId, 1, pageLimit.value);
    };

    // Refresh messages (reset and load first page)
    const refreshMessages = async (chatId) => {
        chatStore.resetPagination();
        return await loadMessages(chatId, 1, pageLimit.value);
    };

    // Check if we should load more messages (for infinite scroll)
    const shouldLoadMore = (
        scrollTop,
        scrollHeight,
        clientHeight,
        threshold = 100
    ) => {
        return (
            scrollTop <= threshold &&
            hasMoreMessages.value &&
            !isLoadingMessages.value
        );
    };

    return {
        // State
        isLoading,
        error,
        currentPage,
        pageLimit,
        hasMoreMessages,
        isLoadingMessages,

        // Methods
        loadMessages,
        loadNextPage,
        loadInitialMessages,
        loadSecretChatMessages,
        loadInitialSecretChatMessages,
        refreshMessages,
        shouldLoadMore,
        loadSecretChats,

        // Store methods (for direct access if needed)
        setMessages: chatStore.setMessages,
        addMessage: chatStore.addMessage,
        resetPagination: chatStore.resetPagination,
    };
}
