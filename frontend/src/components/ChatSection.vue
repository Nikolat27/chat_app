<template>
    <section class="flex flex-col h-full w-full bg-gray-50 font-roboto">
        <!-- Chat Header -->
        <ChatHeader
            v-if="hasActiveChat && chatStore.currentChatUser && !isCurrentChatSecret && !isCurrentChatGroup"
            :user="chatStore.currentChatUser"
            :backend-base-url="backendBaseUrl"
            :is-secret-chat="false"
            @delete-chat="handleDeleteChat"
        />
        
        <!-- Secret Chat Header -->
        <SecretChatHeader
            v-if="hasActiveChat && currentSecretChat && isCurrentChatSecret"
            :secret-chat="currentSecretChat"
            :secret-usernames="chatStore.secretUsernames"
            :backend-base-url="backendBaseUrl"
            :current-user-id="userStore.user_id"
        />

        <!-- Group Chat Header -->
        <GroupChatHeader
            v-if="hasActiveChat && currentGroup && isCurrentChatGroup"
            :group="currentGroup"
            :backend-base-url="backendBaseUrl"
            :current-user-id="userStore.user_id"
            @leave-group="handleLeaveGroup"
            @delete-group="handleDeleteGroup"
            @show-group-info="handleShowGroupInfo"
        />

        <!-- Regular Chat Messages Area -->
        <MessagesArea
            v-if="hasActiveChat && chatStore.currentChatUser && !isCurrentChatGroup"
            :messages="chatStore.messages"
            :current-user-id="userStore.user_id"
            :backend-base-url="backendBaseUrl"
            :user-avatar="userStore.avatar_url"
            :other-user-avatar="chatStore.currentChatUser?.avatar_url"
            :chat-id="getCurrentChatId()"
            :is-secret-chat="isCurrentChatSecret"
            :is-secret-chat-approved="isSecretChatApproved"
            @load-more-messages="handleLoadMoreMessages"
        />
        


        <!-- Group Chat Messages Area -->
        <GroupMessagesArea
            v-if="hasActiveChat && currentGroup && isCurrentChatGroup"
            :messages="groupMessages"
            :current-user-id="userStore.user_id"
            :backend-base-url="backendBaseUrl"
            :user-avatar="userStore.avatar_url ? `${backendBaseUrl}/static/${userStore.avatar_url}` : null"
            :other-user-avatar="null"
            :chat-id="currentGroup?.id"
            :is-loading-messages="isLoadingMessages"
            :is-secret-group="isCurrentChatSecret"
            @load-more-messages="handleLoadMoreGroupMessages"
            @open-secret-key-modal="openSecretKeyModal"
        />
        

        <!-- Debug user avatar -->
        <div v-if="currentGroup && isCurrentChatGroup" style="display: none;">
            Debug - User Avatar: {{ userStore.avatar_url }}
        </div>

        <!-- Regular Chat Message Input -->
        <MessageInput
            v-if="hasActiveChat && chatStore.currentChatUser && !isCurrentChatGroup"
            v-model="newMessage"
            :is-secret-chat="isCurrentChatSecret"
            :is-secret-chat-approved="isSecretChatApproved"
            @send="sendMessage"
            @image-upload="handleImageUpload"
        />

        <!-- Group Chat Message Input -->
        <GroupMessageInput
            v-if="hasActiveChat && currentGroup && isCurrentChatGroup"
            v-model="newGroupMessage"
            :is-secret-group="currentGroup?.type === 'secret'"
            :group-id="currentGroup?.id || ''"
            :is-group-owner="currentGroup?.owner_id === userStore.user_id"
            :key-status="secretKeyStatus"
            @send="sendGroupMessageHandler"
            @image-upload="handleGroupImageUpload"
            @open-secret-key-modal="handleOpenSecretKeyModal"
        />

        <!-- No Chat Selected State -->
        <NoChatSelected v-if="!hasActiveChat" />
        


        <!-- Group Info Modal -->
        <GroupInfoModal
            v-if="!selectedGroupForInfo?.type || selectedGroupForInfo?.type !== 'secret'"
            :is-visible="showGroupInfoModal"
            :group="selectedGroupForInfo"
            :backend-base-url="backendBaseUrl"
            :current-user-id="userStore.user_id"
            @close="handleGroupInfoModalClose"
            @edit-group="handleEditGroupFromInfo"
        />

        <!-- Secret Group Info Modal -->
        <SecretGroupInfoModal
            v-if="selectedGroupForInfo?.type === 'secret'"
            :is-visible="showGroupInfoModal"
            :group="selectedGroupForInfo"
            :backend-base-url="backendBaseUrl"
            :current-user-id="userStore.user_id"
            @close="handleGroupInfoModalClose"
            @edit-group="handleEditGroupFromInfo"
            @leave-group="handleLeaveGroup"
        />

        <!-- Update Group Modal -->
        <UpdateGroupModal
            :is-visible="showUpdateGroupModal"
            :group="selectedGroupForUpdate"
            :backend-base-url="backendBaseUrl"
            @close="handleUpdateGroupModalClose"
            @group-updated="handleGroupUpdated"
        />

        <!-- Secret Group Key Modal -->
        <SecretGroupKeyModal
            :is-visible="showSecretKeyModal"
            :group-id="secretChatIdForKey || ''"
            :group-name="currentGroup?.name || ''"
            :is-group-owner="currentGroup?.owner_id === userStore.user_id"
            @close="closeSecretKeyModal"
            @key-entered="handleSecretKeyLoaded"
        />
    </section>
</template>

<script setup>
import { ref, watch, computed, nextTick } from "vue";
import { useChatStore } from "../stores/chat";
import { useUserStore } from "../stores/users";
import { useGroupStore } from "../stores/groups";
import ChatHeader from "./chat/ChatHeader.vue";
import SecretChatHeader from "./chat/SecretChatHeader.vue";
import GroupChatHeader from "./chat/GroupChatHeader.vue";
import NoChatSelected from "./chat/NoChatSelected.vue";
import MessagesArea from "./chat/MessagesArea.vue";
import MessageInput from "./chat/MessageInput.vue";
import GroupMessagesArea from "./chat/GroupMessagesArea.vue";
import GroupMessageInput from "./chat/GroupMessageInput.vue";
import GroupInfoModal from "./ui/GroupInfoModal.vue";
import SecretGroupInfoModal from "./ui/SecretGroupInfoModal.vue";
import UpdateGroupModal from "./ui/UpdateGroupModal.vue";
import SecretGroupKeyModal from "./chat/SecretGroupKeyModal.vue";
import { useWebSocket } from "../composables/useWebSocket";
import { useGroupChat } from "../composables/useGroupChat";
import { useMessagePagination } from "../composables/useMessagePagination";
import { useMessageDeletion } from "../composables/useMessageDeletion";
import { useE2EE } from "../composables/useE2EE";
import { useSecretChatEncryption } from "../composables/useSecretChatEncryption";
import axiosInstance from "../axiosInstance";
import { showMessage, showError } from "../utils/toast";

const chatStore = useChatStore();
const userStore = useUserStore();
const groupStore = useGroupStore();
const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;
const newMessage = ref("");
const secretKeyStatus = ref('not-entered');

// Group chat composable
const {
    isGroupConnected,
    groupMessages,
    newGroupMessage,
    // Pagination state
    currentPage,
    pageLimit,
    hasMoreMessages,
    isLoadingMessages,
    establishGroupConnection,
    sendGroupMessage,
    closeGroupConnection,
    getGroupConnectionStatus,
    addGroupMessage,
    clearGroupMessages,
    loadGroupUsers,
    loadGroupMessages,
    loadNextGroupPage,
    loadInitialGroupMessages,
    getUsernameBySenderId,
    getAvatarBySenderId,
    handleIncomingGroupMessage
} = useGroupChat();

// Check if current chat is a secret chat
const isCurrentChatSecret = computed(() => {
    if (!chatStore.currentChatUser?.id) return false;
    
    // If the current chat user has a secret chat ID, it's a secret chat
    if (chatStore.currentChatUser.secret_chat_id) {
        return true;
    }
    
    // If no secret_chat_id is present, it's definitely not a secret chat
    return false;
});

// Get the current secret chat object
const currentSecretChat = computed(() => {
    if (!isCurrentChatSecret.value || !chatStore.currentChatUser?.secret_chat_id) return null;
    
    // Find the secret chat by its ID
    return chatStore.secretChats?.find(chat => chat.id === chatStore.currentChatUser.secret_chat_id) || null;
});

// Check if the current secret chat is approved
const isSecretChatApproved = computed(() => {
    if (!isCurrentChatSecret.value) return true; // Not a secret chat, so "approved"
    
    const secretChat = currentSecretChat.value;
    if (!secretChat) return false; // No secret chat found, not approved
    
    return secretChat.user_2_accepted === true;
});

// Check if current chat is a group chat
const isCurrentChatGroup = computed(() => {
    return groupStore.currentGroup !== null;
});

// Get the current group object
const currentGroup = computed(() => {
    return groupStore.currentGroup;
});

// Check if there's an active chat (either regular chat or group chat)
const hasActiveChat = computed(() => {
    return chatStore.currentChatUser || currentGroup.value;
});



// WebSocket management
const { establishConnection, sendMessage: sendWebSocketMessage, closeConnection, getConnectionStatus } =
    useWebSocket();

// Message pagination
const { loadNextPage, loadInitialMessages, loadInitialSecretChatMessages } = useMessagePagination();

// Message deletion
const { updateMessageId } = useMessageDeletion();

// E2EE
const { encryptMessage, decryptMessage, loadChatSymmetricKey } = useE2EE();
const { loadSecretChatSymmetricKey, validateSecretChatForEncryption, encryptMessageForSending } = useSecretChatEncryption();

// Watch for chat user changes to manage WebSocket connections
watch(
    () => chatStore.currentChatUser?.id,
    async (newUserId, oldUserId) => {
        console.log('🔄 Chat user changed:', { 
            oldUserId, 
            newUserId, 
            currentChatUser: chatStore.currentChatUser,
            isSecretChat: isCurrentChatSecret.value 
        });
        
        // Only close previous connection if we're actually switching to a different user
        if (oldUserId && oldUserId !== newUserId) {
            console.log('🔌 Closing previous WebSocket connection due to user change');
            closeConnection();
            // Wait for connection to close before establishing new one
            await new Promise(resolve => setTimeout(resolve, 300));
        }

        if (newUserId) {
            // Add a small delay to ensure chat is properly added to store
            await nextTick();
            // Try to get chat data with retries
            let chatData = null;
            let retries = 0;
            const maxRetries = 5; // Increased retries
            while (!chatData && retries < maxRetries) {
                chatData = getChatData(newUserId);
                if (!chatData) {
                    await new Promise(resolve => setTimeout(resolve, 200)); // Increased wait time
                    retries++;
                }
            }
            if (chatData) {
                console.log('🔌 Establishing WebSocket connection for chat:', chatData);
                establishConnection(chatData, handleIncomingMessage);
                
                // Load initial messages for the chat
                const chatId = getCurrentChatId();
                if (chatId) {
                    console.log('📥 Loading initial messages for chat:', chatId);
                    if (isCurrentChatSecret.value) {
                        await loadInitialSecretChatMessages(chatId);
                    } else {
                        await loadInitialMessages(chatId);
                    }
                }
                
                // Wait a bit for connection to establish
                await new Promise(resolve => setTimeout(resolve, 500)); // Increased wait time
            } else {
                console.error('❌ Failed to get chat data after retries');
            }
        }
    }
);

// Watch for when currentChatUser becomes null to close WebSocket
watch(
    () => chatStore.currentChatUser,
    (newChatUser, oldChatUser) => {
        if (oldChatUser && !newChatUser) {
            console.log('🔄 Chat user cleared, closing WebSocket connection');
            closeConnection();
        }
    }
);

// Watch for group changes to manage WebSocket connections
watch(
    () => groupStore.currentGroup?.id,
    async (newGroupId, oldGroupId) => {
        console.log('🔄 Group changed:', { oldGroupId, newGroupId, currentGroup: groupStore.currentGroup });
        
        // Clear group messages when switching groups
        if (oldGroupId && oldGroupId !== newGroupId) {
            console.log('🔌 Closing previous group WebSocket connection');
            closeGroupConnection();
            clearGroupMessages();
            // Wait for connection to close before establishing new one
            await new Promise(resolve => setTimeout(resolve, 300));
        }

        if (newGroupId) {
            console.log('🔄 Starting group initialization for:', newGroupId);
            
            // Check if this is a secret group
            const isSecretGroup = groupStore.currentGroup.type === 'secret';
            console.log('🔐 Group type:', groupStore.currentGroup.type, 'isSecretGroup:', isSecretGroup);
            console.log('🔐 Full group data:', groupStore.currentGroup);
            
            // Add a small delay to ensure group is properly set
            await nextTick();
            
            // Load group users first (needed for usernames)
            console.log('👥 About to load group users...');
            console.log('👥 loadGroupUsers function:', typeof loadGroupUsers);
            console.log('👥 Calling loadGroupUsers with groupId:', newGroupId);
            await loadGroupUsers(newGroupId);
            console.log('👥 Group users loaded successfully');
            
            // Load existing group messages with encryption for secret groups
            console.log('📥 About to load group messages...');
            await loadInitialGroupMessages(newGroupId, isSecretGroup);
            console.log('📥 Group messages loaded successfully');
            
            const groupData = getGroupChatData(newGroupId);
            if (groupData) {
                // Add group type information to the data
                const enhancedGroupData = {
                    ...groupData,
                    type: groupStore.currentGroup.type
                };
                console.log('🔌 Establishing group WebSocket connection:', enhancedGroupData);
                establishGroupConnection(enhancedGroupData, (data) => handleIncomingGroupMessage(data, newGroupId, isSecretGroup), isSecretGroup);
                // Wait a bit for connection to establish
                await new Promise(resolve => setTimeout(resolve, 500));
            } else {
                console.error('❌ Failed to get group chat data');
            }
        }
    }
);

// Watch for when currentGroup becomes null to close WebSocket
watch(
    () => groupStore.currentGroup,
    (newGroup, oldGroup) => {
        if (oldGroup && !newGroup) {
            console.log('🔄 Group cleared, closing WebSocket connection');
            closeConnection();
        }
    }
);

// Get chat data for WebSocket connection
const getChatData = (targetUserId) => {
    // Check if this is a secret chat
    if (chatStore.currentChatUser?.secret_chat_id) {
        const senderId = userStore.user_id;
        const receiverId = targetUserId;
        return {
            chatId: chatStore.currentChatUser.secret_chat_id,
            senderId,
            receiverId,
            backendBaseUrl,
            isSecretChat: true,
        };
    }

    let chat = null;
    // Try to find by chat_id if available in currentChatUser
    if (chatStore.currentChatUser?.chat_id) {
        chat = chatStore.chats?.find((c) => c.id === chatStore.currentChatUser.chat_id);
        if (chat) {
            const senderId = userStore.user_id;
            // Find the other participant
            const receiverId = chat.participants?.find((id) => id !== senderId) || targetUserId;
            return {
                chatId: chat.id,
                senderId,
                receiverId,
                backendBaseUrl,
                isSecretChat: false,
            };
        }
    }
    // Try to find by participants
    chat = chatStore.chats?.find(
        (c) =>
            c.participants &&
            c.participants.includes(targetUserId) &&
            c.participants.includes(userStore.user_id)
    );
    if (!chat) {
        return null;
    }
    if (!chat.participants || chat.participants.length < 2) {
        return null;
    }
    const senderId = userStore.user_id;
    const receiverId = targetUserId;
    return {
        chatId: chat.id,
        senderId,
        receiverId,
        backendBaseUrl,
        isSecretChat: false,
    };
};

// Get group chat data for WebSocket connection
const getGroupChatData = (groupId) => {
    if (!groupId || !userStore.user_id) {
        return null;
    }
    
    return {
        chatId: groupId, // Use groupId as chatId for group chats
        senderId: userStore.user_id,
        receiverId: null, // Not needed for group chats
        backendBaseUrl,
        isSecretChat: false,
        isGroupChat: true,
        groupId: groupId,
    };
};

// Handle incoming messages
const handleIncomingMessage = async (data) => {
    console.log('📨 Received WebSocket message:', data);
    const message = parseIncomingMessage(data);
    console.log('📨 Parsed message:', message);

    // Decrypt message if this is a secret chat
    let decryptedContent = message.content;
    if (isCurrentChatSecret.value && message.content) {
        try {
            // For secret chats, use the secret_chat_id instead of chat_id
            const chatId = chatStore.currentChatUser?.secret_chat_id || message.chat_id;
            console.log('🔐 Decrypting incoming message for secret chat:', {
                message_chat_id: message.chat_id,
                secret_chat_id: chatStore.currentChatUser?.secret_chat_id,
                used_chat_id: chatId,
                content_length: message.content.length,
                is_secret: message.is_secret
            });
            
            // For WebSocket messages in secret chats, always try to decrypt if content looks encrypted
            // WebSocket messages might not have the is_secret flag, so we rely on the chat context
            const shouldDecrypt = message.content.length > 20 && 
                /^[A-Za-z0-9+/=]+$/.test(message.content) && // Base64 pattern
                message.content.length % 4 === 0; // Base64 length check
            
            console.log('🔍 Incoming WebSocket message encryption check:', {
                shouldDecrypt,
                length: message.content.length,
                isSecret: message.is_secret,
                isBase64: /^[A-Za-z0-9+/=]+$/.test(message.content),
                isSecretChat: isCurrentChatSecret.value
            });
            
            if (shouldDecrypt) {
                // Check if we have the symmetric key, if not try to load it
                const { hasSymmetricKey } = useE2EE();
                const { loadSecretChatSymmetricKey } = useSecretChatEncryption();
                const keyAvailable = await hasSymmetricKey(chatId);
                
                if (!keyAvailable) {
                    console.log('🔐 Symmetric key not found, attempting to load...');
                    await loadSecretChatSymmetricKey(chatId);
                }
                
                decryptedContent = await decryptMessage(message.content, chatId);
                console.log('✅ Successfully decrypted incoming WebSocket message');
            } else {
                console.log('🔍 Incoming WebSocket message appears to be plaintext or not encrypted');
            }
        } catch (error) {
            console.error('Error decrypting WebSocket message:', error);
            // If decryption fails, show encrypted content or error message
            decryptedContent = '[Encrypted message - decryption failed]';
        }
    }

    const decryptedMessage = {
        ...message,
        content: decryptedContent,
        // Add timestamp if not present (for WebSocket messages)
        created_at: message.created_at || new Date().toISOString(),
        // Ensure type field is set for image messages
        type: message.content_type === 'image' ? 'image' : message.type
    };

    // Check if this is a confirmation of a sent message (same content and sender)
    const existingMessage = chatStore.messages.find(
        (msg) =>
            msg.content === decryptedMessage.content &&
            msg.sender_id === decryptedMessage.sender_id &&
            msg.id &&
            msg.id.startsWith("temp-")
    );

    // Check for duplicate messages based on content and sender
    // This handles both WebSocket messages (no ID) and API messages (with ID)
    const duplicateMessage = chatStore.messages.find(
        (msg) =>
            msg.content === decryptedMessage.content &&
            msg.sender_id === decryptedMessage.sender_id &&
            // If the incoming message has an ID, check for exact ID match
            // If not, check if we have a very recent message with same content and sender (within 100ms)
            (decryptedMessage.id ? msg.id === decryptedMessage.id : 
             // For WebSocket messages without ID, only skip if it's the exact same message received twice very quickly
             Math.abs(new Date(msg.created_at) - new Date(decryptedMessage.created_at)) < 100) // Within 100ms
    );

    if (duplicateMessage) {
        console.log('🔄 Duplicate message detected, skipping:', decryptedMessage.content);
        console.log('🔄 Duplicate details:', {
            existingMessage: duplicateMessage,
            incomingMessage: decryptedMessage,
            timeDiff: Math.abs(new Date(duplicateMessage.created_at) - new Date(decryptedMessage.created_at))
        });
        return; // Skip adding duplicate message
    }

    if (existingMessage && decryptedMessage.id && !decryptedMessage.id.startsWith("temp-")) {
        // Update the temp ID with the real ID from backend
        console.log('🔄 Updating temp message with real ID:', existingMessage.id, '->', decryptedMessage.id);
        chatStore.updateMessageId(existingMessage.id, decryptedMessage.id);
    } else {
        // This is a new message from someone else
        console.log('➕ Adding new message:', {
            id: decryptedMessage.id,
            content: decryptedMessage.content,
            sender_id: decryptedMessage.sender_id,
            created_at: decryptedMessage.created_at
        });
        chatStore.addMessage(decryptedMessage);
    }
};



// Parse incoming message data
const parseIncomingMessage = (data) => {
    // If the backend sends the message object directly
    if (data.content && typeof data.content === "object") {
        return data.content;
    }

    // If the backend sends content as a JSON string
    if (typeof data.content === "string" && data.content.startsWith("{")) {
        try {
            return JSON.parse(data.content);
        } catch (e) {
            return { content: data.content };
        }
    }

    // If the backend sends the message object as the root data
    if (data.chat_id && data.sender_id && data.content) {
        return data;
    }

    // Default case
    return data;
};

// Send message
const sendMessage = async () => {
    if (!newMessage.value.trim()) return;

    // Check if this is a secret chat that's not approved
    if (isCurrentChatSecret.value && !isSecretChatApproved.value) {
        showError('Cannot send messages in unapproved secret chat');
        return;
    }

    // Handle group chat
    if (isCurrentChatGroup.value) {
        await sendGroupMessage();
        return;
    }

    const targetUserId = chatStore.currentChatUser?.id;
    
    const chatData = getChatData(targetUserId);

    if (!chatData) {
        return;
    }

    // Check WebSocket connection status
    const connectionStatus = getConnectionStatus();
    
    if (!connectionStatus.isConnected) {
        console.log('🔌 WebSocket not connected, establishing connection...');
        establishConnection(chatData, handleIncomingMessage);
        
        // Wait for connection with retries
        let retries = 0;
        const maxRetries = 5;
        while (!getConnectionStatus().isConnected && retries < maxRetries) {
            await new Promise(resolve => setTimeout(resolve, 200));
            retries++;
        }
        
        if (!getConnectionStatus().isConnected) {
            showError('Failed to establish WebSocket connection. Please try again.');
            return;
        }
    }

    // Create temporary ID for immediate display
    const tempId = `temp-${Date.now()}-${Math.random()
        .toString(36)
        .substr(2, 9)}`;

    let messageContent = newMessage.value;
    
    // Encrypt message if this is a secret chat
    if (isCurrentChatSecret.value) {
        try {
            // Validate that the chat is ready for messaging
            const validation = await validateSecretChatForEncryption(chatData.chatId);
            if (!validation.valid) {
                showError(validation.message);
                return;
            }
            
            // Now encrypt the message
            messageContent = await encryptMessageForSending(newMessage.value, chatData.chatId);
        } catch (error) {
            console.error('Error encrypting message:', error);
            showError('Failed to encrypt message: ' + error.message);
            return;
        }
    }

    const messageData = {
        id: tempId,
        chat_id: chatData.chatId,
        sender_id: chatData.senderId,
        receiver_id: chatData.receiverId,
        content: messageContent,
        created_at: new Date().toISOString(),
    };

    // Add message to store immediately with temp ID (store decrypted content for display)
    chatStore.addMessage({
        ...messageData,
        content: newMessage.value, // Store decrypted content for display
    });

    // Send message via WebSocket with new structure
    const messagePayload = {
        sender_id: chatData.senderId.toString(),
        receiver_id: chatData.receiverId.toString(),
        content: messageContent,
        content_address: "",
        content_type: "text"
    };
    const success = sendWebSocketMessage(JSON.stringify(messagePayload));
    
    newMessage.value = "";
};

// Handle image upload for regular chats
const handleImageUpload = async (file) => {
    try {
        console.log('📸 Uploading image for regular chat:', file);
        
        // Create FormData for multipart upload
        const formData = new FormData();
        formData.append('file', file);
        
        // Upload image to server with chat_id parameter
        const chatId = getCurrentChatId();
        const response = await axiosInstance.post(`/api/message/upload-chat-image/${chatId}`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });
        
        const imageAddress = response.data.image_address;
        console.log('📸 Image uploaded successfully:', imageAddress);
        
        // Send image message via WebSocket
        await sendImageMessage(imageAddress, false);
        
    } catch (error) {
        console.error('❌ Image upload failed:', error);
        showError('Failed to upload image. Please try again.');
    }
};

// Handle image upload for group chats
const handleGroupImageUpload = async (file) => {
    try {
        console.log('📸 Uploading image for group chat:', file);
        
        // Create FormData for multipart upload
        const formData = new FormData();
        formData.append('file', file);
        
        // Upload image to server using group-specific endpoint with group_id parameter
        const response = await axiosInstance.post(`/api/message/upload-group-image/${currentGroup.value.id}`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });
        
        const imageAddress = response.data.image_address;
        console.log('📸 Group image uploaded successfully:', imageAddress);
        
        // Send image message via WebSocket for group
        await sendGroupImageMessage(imageAddress);
        
    } catch (error) {
        console.error('❌ Group image upload failed:', error);
        showError('Failed to upload image. Please try again.');
    }
};

// Send image message via WebSocket
const sendImageMessage = async (imageAddress, isGroup = false) => {
    if (isGroup) {
        // Handle group image message
        if (!currentGroup.value) {
            showError('No group selected');
            return;
        }
        
        // Create temporary ID for immediate display
        const tempId = `temp-${Date.now()}-${Math.random()
            .toString(36)
            .substr(2, 9)}`;

        const messageData = {
            id: tempId,
            sender_id: userStore.user_id,
            content: '',
            content_address: imageAddress,
            type: 'image',
            created_at: new Date().toISOString(),
            sender_name: userStore.username || 'Unknown User',
            sender_avatar: userStore.avatar_url || null
        };

        // Add message to group messages immediately with temp ID
        addGroupMessage(messageData);
        
        const messagePayload = {
            sender_id: userStore.user_id.toString(),
            receiver_id: currentGroup.value.id.toString(),
            content: '',
            content_address: imageAddress,
            content_type: 'image'
        };
        
        // Send via group WebSocket
        const success = sendWebSocketMessage(JSON.stringify(messagePayload));
        if (success) {
            showMessage('Image sent successfully!');
        } else {
            showError('Failed to send image message');
        }
        
    } else {
        // Handle regular chat image message
        if (!chatStore.currentChatUser) {
            showError('No chat selected');
            return;
        }
        
        const targetUserId = chatStore.currentChatUser.id;
        const chatData = getChatData(targetUserId);
        
        if (!chatData) {
            showError('Invalid chat data');
            return;
        }
        
        // Create temporary ID for immediate display
        const tempId = `temp-${Date.now()}-${Math.random()
            .toString(36)
            .substr(2, 9)}`;

        const messageData = {
            id: tempId,
            chat_id: chatData.chatId,
            sender_id: chatData.senderId,
            receiver_id: chatData.receiverId,
            content: '',
            content_address: imageAddress,
            type: 'image',
            created_at: new Date().toISOString(),
        };

        // Add message to store immediately with temp ID
        chatStore.addMessage(messageData);
        
        const messagePayload = {
            sender_id: chatData.senderId.toString(),
            receiver_id: chatData.receiverId.toString(),
            content: '',
            content_address: imageAddress,
            content_type: 'image'
        };
        
        // Send via regular WebSocket
        const success = sendWebSocketMessage(JSON.stringify(messagePayload));
        if (success) {
            showMessage('Image sent successfully!');
        } else {
            showError('Failed to send image message');
        }
    }
};

// Send group image message via WebSocket
const sendGroupImageMessage = async (imageAddress) => {
    if (!currentGroup.value) {
        showError('No group selected');
        return;
    }
    
    // Check if this is a secret group (secret groups can't use image upload)
    if (currentGroup.value.type === 'secret') {
        showError('Image upload is not available for secret groups');
        return;
    }
    
    // Create temporary ID for immediate display
    const tempId = `temp-${Date.now()}-${Math.random()
        .toString(36)
        .substr(2, 9)}`;

    const messageData = {
        id: tempId,
        sender_id: userStore.user_id,
        content: '',
        content_address: imageAddress,
        type: 'image',
        created_at: new Date().toISOString(),
        sender_name: userStore.username || 'Unknown User',
        sender_avatar: userStore.avatar_url || null
    };

    // Add message to group messages immediately with temp ID
    addGroupMessage(messageData);
    
    // Create GroupMessage payload structure
    const groupMessagePayload = {
        sender_id: userStore.user_id.toString(),
        content: '',
        content_address: imageAddress,
        content_type: 'image'
    };
    
    // Send via group WebSocket using the existing sendGroupMessage function
    // For image messages, we need to send the GroupMessage structure
    const imageMessageData = {
        ...messageData,
        content: '', // Empty content for image messages
        content_address: imageAddress,
        content_type: 'image'
    };
    
    const success = await sendGroupMessage(imageMessageData, currentGroup.value.id, false);
    if (success) {
        showMessage('Image sent successfully!');
    } else {
        showError('Failed to send image message');
    }
};

// Send group message handler
const sendGroupMessageHandler = async () => {
    if (!groupStore.currentGroup?.id) {
        showError('No group selected');
        return;
    }

    if (!newGroupMessage.value.trim()) return;

    // Check if this is a secret group
    const isSecretGroup = groupStore.currentGroup.type === 'secret';
    console.log('🔐 Sending group message, is secret group:', isSecretGroup);

    // Check group WebSocket connection status
    const connectionStatus = getGroupConnectionStatus();
    
    if (!connectionStatus.isConnected) {
        console.log('🔌 Group WebSocket not connected, establishing connection...');
        const groupData = getGroupChatData(groupStore.currentGroup.id);
        if (groupData) {
            establishGroupConnection(groupData, (data) => handleIncomingGroupMessage(data, groupStore.currentGroup.id, isSecretGroup), isSecretGroup);
            
            // Wait for connection with retries and better error handling
            let retries = 0;
            const maxRetries = 10; // Increased retries
            const retryDelay = 300; // Increased delay
            
            while (retries < maxRetries) {
                await new Promise(resolve => setTimeout(resolve, retryDelay));
                const currentStatus = getGroupConnectionStatus();
                
                if (currentStatus.isConnected) {
                    console.log('✅ Group WebSocket connection established successfully');
                    break;
                }
                
                retries++;
                console.log(`🔄 Connection attempt ${retries}/${maxRetries} failed, retrying...`);
            }
            
            const finalStatus = getGroupConnectionStatus();
            if (!finalStatus.isConnected) {
                console.error('❌ Failed to establish group WebSocket connection after all retries');
                showError('Failed to establish group WebSocket connection. Please try again.');
                return;
            }
        } else {
            console.error('❌ Failed to get group chat data for WebSocket connection');
            showError('Failed to get group data. Please try again.');
            return;
        }
    } else {
        console.log('✅ Group WebSocket already connected');
    }

    // Create temporary ID for immediate display
    const tempId = `temp-${Date.now()}-${Math.random()
        .toString(36)
        .substr(2, 9)}`;

    const messageData = {
        id: tempId,
        sender_id: userStore.user_id,
        content: newGroupMessage.value,
        message_type: 1, // Text message
        created_at: new Date().toISOString(),
        sender_name: userStore.username || 'Unknown User',
        sender_avatar: userStore.avatar_url || null
    };

    // Add message to group messages immediately with temp ID
    addGroupMessage(messageData);

    // Send message via group WebSocket with encryption for secret groups
    let retryCount = 0;
    const maxRetries = 3;
    let success = false;
    
    while (!success && retryCount < maxRetries) {
        success = await sendGroupMessage(messageData, groupStore.currentGroup.id, isSecretGroup);
        
        if (!success) {
            retryCount++;
            console.log(`🔌 Retry ${retryCount}/${maxRetries} for sending group message`);
            
            if (retryCount < maxRetries) {
                // Wait before retrying
                await new Promise(resolve => setTimeout(resolve, 1000 * retryCount));
            }
        }
    }
    
    if (!success) {
        console.error("Failed to send group message after all retries");
        showError('Failed to send group message. Please try again.');
        return;
    }
    
    newGroupMessage.value = "";
};



// Get current chat ID
const getCurrentChatId = () => {
    // Check if this is a secret chat
    if (chatStore.currentChatUser?.secret_chat_id) {
        return chatStore.currentChatUser.secret_chat_id;
    }
    
    const targetUserId = chatStore.currentChatUser?.id;
    const chat = chatStore.chats?.find(
        (c) =>
            c.participants &&
            c.participants.includes(targetUserId) &&
            c.participants.includes(userStore.user_id)
    );
    return chat?.id || null;
};

// Handle loading more messages
const handleLoadMoreMessages = async () => {
    const chatId = getCurrentChatId();
    if (chatId) {
        // Check if this is a secret chat
        if (chatStore.currentChatUser?.secret_chat_id) {
            // For secret chats, we might need a different approach for pagination
            // For now, we'll use the regular loadNextPage but with the secret chat ID
            await loadNextPage(chatId);
        } else {
            await loadNextPage(chatId);
        }
    }
};

// Handle loading more group messages
const handleLoadMoreGroupMessages = async () => {
    if (currentGroup.value?.id) {
        console.log('📥 Loading more group messages for group:', currentGroup.value.id);
        await loadNextGroupPage(currentGroup.value.id);
    }
};

// Handle chat deletion
const handleDeleteChat = async (user) => {
    try {
        // Find the chat ID for this user
        const chatId = getCurrentChatId();
        if (!chatId) {
            showError('No chat ID found for this user');
            return;
        }

        // Delete the chat from backend
        const response = await axiosInstance.delete(`/api/chat/delete/${chatId}`);
        
        if (response.status === 200) {
            // Remove chat from store
            const chatIndex = chatStore.chats.findIndex(chat => chat.id === chatId);
            if (chatIndex !== -1) {
                chatStore.chats.splice(chatIndex, 1);
            }
            
            // Clear current chat user
            chatStore.clearChatUser();
            // Close WebSocket connection
            closeConnection();
            
            showMessage('Chat deleted successfully');
        }
    } catch (error) {
        console.error('Error deleting chat:', error);
        showError('Failed to delete chat. Please try again.');
    }
};

// Handle group leave
const handleLeaveGroup = async (group) => {
    try {
        const confirmed = confirm(`Are you sure you want to leave "${group.name}"?`);
        if (!confirmed) return;

        // Use the correct leave method based on group type
        if (group.type === 'secret') {
            await groupStore.leaveSecretGroup(group.id);
            showMessage('Successfully left secret group');
        } else {
            await groupStore.leaveGroup(group.id);
            showMessage('Successfully left group');
        }
        
        // Clear current group and chat user
        groupStore.clearCurrentGroup();
        chatStore.clearChatUser();
        closeConnection();
    } catch (error) {
        console.error('Failed to leave group:', error);
        showError('Failed to leave group. Please try again.');
    }
};

// Handle group delete
const handleDeleteGroup = async (group) => {
    try {
        const confirmed = confirm(`Are you sure you want to delete "${group.name}"? This action cannot be undone.`);
        if (!confirmed) return;

        // Use the correct deletion method based on group type
        if (group.type === 'secret') {
            await groupStore.deleteSecretGroup(group.id);
            showMessage('Secret group deleted successfully');
        } else {
            await groupStore.deleteGroup(group.id);
            showMessage('Group deleted successfully');
        }
        
        // Clear current group and chat user
        groupStore.clearCurrentGroup();
        chatStore.clearChatUser();
        closeConnection();
    } catch (error) {
        console.error('Failed to delete group:', error);
        showError('Failed to delete group. Please try again.');
    }
};

// Handle show group info
const showGroupInfoModal = ref(false);
const selectedGroupForInfo = ref(null);

// Handle edit group
const showUpdateGroupModal = ref(false);
const selectedGroupForUpdate = ref(null);

const handleShowGroupInfo = (group) => {
    selectedGroupForInfo.value = group;
    showGroupInfoModal.value = true;
};

const handleGroupInfoModalClose = () => {
    showGroupInfoModal.value = false;
    selectedGroupForInfo.value = null;
};

const handleEditGroupFromInfo = (group) => {
    // Close info modal and open update modal
    handleGroupInfoModalClose();
    selectedGroupForUpdate.value = group;
    showUpdateGroupModal.value = true;
};

const handleUpdateGroupModalClose = () => {
    showUpdateGroupModal.value = false;
    selectedGroupForUpdate.value = null;
};

const handleGroupUpdated = (updatedGroup) => {
    console.log('Group updated:', updatedGroup);
    // Update the current group if it's the same one
    if (currentGroup.value && currentGroup.value.id === updatedGroup.id) {
        currentGroup.value = { ...currentGroup.value, ...updatedGroup };
    }
    handleUpdateGroupModalClose();
};

// Secret key modal handling
const showSecretKeyModal = ref(false);
const secretChatIdForKey = ref(null);

const openSecretKeyModal = (chatId) => {
    secretChatIdForKey.value = chatId;
    showSecretKeyModal.value = true;
};

const closeSecretKeyModal = () => {
    showSecretKeyModal.value = false;
    secretChatIdForKey.value = null;
};

const handleSecretKeyLoaded = async (groupId) => {
    console.log('Secret key loaded for group:', groupId);
    secretKeyStatus.value = 'entered';
    showMessage('Secret key entered successfully! You can now send messages in this secret group.');
    closeSecretKeyModal();
};

const handleOpenSecretKeyModal = () => {
    if (currentGroup.value?.id) {
        secretChatIdForKey.value = currentGroup.value.id;
        showSecretKeyModal.value = true;
    }
};
</script>

<style scoped>
.font-roboto {
    font-family: "Roboto", Arial, sans-serif;
}
</style>
