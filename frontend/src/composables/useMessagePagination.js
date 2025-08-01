import { ref, computed } from "vue";
import { useChatStore } from "../stores/chat";
import { useE2EE } from "./useE2EE";
import { useSecretChatEncryption } from "./useSecretChatEncryption";
import axiosInstance from "../axiosInstance";

export function useMessagePagination() {
    const chatStore = useChatStore();
    const { decryptMessage } = useE2EE();

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

            console.log(rawMessages);

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
    const loadSecretChatMessages = async (
        secretChatId,
        page = 1,
        limit = 20
    ) => {
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

            console.log(response.data);

            // Parse each message from JSON string to object and decrypt if needed
            const messages = await Promise.all(
                rawMessages
                    .map(async (msg) => {
                        if (typeof msg === "string") {
                            try {
                                const parsed = JSON.parse(msg);
                                return parsed;
                            } catch (e) {
                                console.error(
                                    "Failed to parse secret chat message:",
                                    msg,
                                    e
                                );
                                return null;
                            }
                        }
                        return msg; // If it's already an object, return as is
                    })
                    .filter((msg) => msg !== null)
            );

            // Decrypt messages if they are encrypted
            const decryptedMessages = await Promise.all(
                messages.map(async (msg) => {
                    if (msg.content && typeof msg.content === "string") {
                        console.log(
                            "🔍 Checking message for encryption:",
                            msg.id,
                            "Content length:",
                            msg.content.length,
                            "Is secret:",
                            msg.is_secret
                        );

                        // For secret chat messages, always try to decrypt if they have encrypted content
                        const shouldDecrypt =
                            msg.is_secret === true &&
                            msg.content.length > 20 &&
                            /^[A-Za-z0-9+/=]+$/.test(msg.content) && // Base64 pattern
                            msg.content.length % 4 === 0; // Base64 length check

                        console.log("🔍 Message encryption check:", {
                            id: msg.id,
                            length: msg.content.length,
                            isSecret: msg.is_secret,
                            isBase64: /^[A-Za-z0-9+/=]+$/.test(msg.content),
                            lengthMod4: msg.content.length % 4 === 0,
                            shouldDecrypt,
                        });

                        if (shouldDecrypt) {
                            try {
                                console.log(
                                    "🔐 Attempting to decrypt message:",
                                    msg.id
                                );

                                // Ensure we have the symmetric key loaded
                                const { hasSymmetricKey } = useE2EE();
                                const { loadSecretChatSymmetricKey } =
                                    useSecretChatEncryption();
                                const keyAvailable = await hasSymmetricKey(
                                    secretChatId
                                );

                                if (!keyAvailable) {
                                    console.log(
                                        "🔐 Symmetric key not found, attempting to load..."
                                    );
                                    await loadSecretChatSymmetricKey(
                                        secretChatId
                                    );
                                }

                                const decryptedContent = await decryptMessage(
                                    msg.content,
                                    secretChatId
                                );
                                console.log(
                                    "✅ Successfully decrypted message:",
                                    msg.id,
                                    "Content:",
                                    decryptedContent
                                );
                                return {
                                    ...msg,
                                    content: decryptedContent,
                                };
                            } catch (error) {
                                console.error(
                                    "❌ Error decrypting message:",
                                    msg.id,
                                    error
                                );
                                // If decryption fails, return the encrypted content with a note
                                return {
                                    ...msg,
                                    content:
                                        "[Encrypted message - decryption failed]",
                                };
                            }
                        } else {
                            console.log(
                                "🔍 Message appears to be plaintext or not encrypted:",
                                msg.id
                            );
                        }
                    }
                    return msg; // Not encrypted, return as is
                })
            );

            // For pagination, we'll assume there are more messages if we got a full page
            const hasMore = rawMessages.length >= limit;
            const totalPages = Math.ceil(rawMessages.length / limit) || 1;

            // Update store with decrypted messages
            chatStore.setMessages(decryptedMessages, page === 1);
            chatStore.setPaginationState(page, hasMore, false);

            return { messages: decryptedMessages, hasMore, totalPages };
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
