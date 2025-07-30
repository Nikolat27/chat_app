import { ref } from 'vue';
import { showError } from '../utils/toast';
import axiosInstance from '../axiosInstance';
import { useGroupStore } from '../stores/groups';
import { useUserStore } from '../stores/users';
import { useSecretGroupE2EE } from './useSecretGroupE2EE';
import { useKeyPair } from './useKeyPair';

let groupSocket = null;
let groupUsers = ref({});

export function useGroupChat() {
    const groupStore = useGroupStore();
    const userStore = useUserStore();
    const { 
        encryptGroupMessage, 
        decryptGroupMessage, 
        loadSecretGroupSymmetricKey,
        generateAndUploadGroupSymmetricKeys,
        hasGroupSymmetricKey
    } = useSecretGroupE2EE();
    
    const isGroupConnected = ref(false);
    const groupMessages = ref([]);
    const newGroupMessage = ref('');
    
    // Pagination state
    const currentPage = ref(1);
    const pageLimit = ref(20);
    const hasMoreMessages = ref(true);
    const isLoadingMessages = ref(false);

    // Establish group WebSocket connection
    const establishGroupConnection = (groupData, onMessageCallback, isSecretGroup = false) => {
        console.log("🔌 Establishing group WebSocket connection with data:", groupData, "isSecretGroup:", isSecretGroup);
        
        // Close existing group connection if any
        if (groupSocket) {
            console.log("🔌 Closing existing group WebSocket connection");
            groupSocket.close();
            groupSocket = null;
        }

        const { groupId, senderId, backendBaseUrl } = groupData;

        if (!groupId || !senderId || !backendBaseUrl) {
            console.error("Missing required data for group WebSocket connection:", { groupId, senderId, backendBaseUrl });
            return;
        }

        // Use different WebSocket URL based on group type
        let wsUrl;
        if (isSecretGroup) {
            wsUrl = `${backendBaseUrl.replace(/^http/, "ws")}/api/websocket/secret-group/add/${groupId}?sender_id=${senderId}`;
            console.log("🔐 Creating secret group WebSocket connection to:", wsUrl);
        } else {
            wsUrl = `${backendBaseUrl.replace(/^http/, "ws")}/api/websocket/group/add/${groupId}?sender_id=${senderId}`;
            console.log("🔌 Creating regular group WebSocket connection to:", wsUrl);
        }
        groupSocket = new WebSocket(wsUrl);

        groupSocket.onopen = () => {
            console.log("🔌 Group WebSocket connected for group:", groupId);
            isGroupConnected.value = true;
        };

        groupSocket.onmessage = (event) => {
            console.log("📨 Received group WebSocket message:", event.data);
            try {
                const data = JSON.parse(event.data);
                console.log("📨 Parsed group WebSocket message:", data);
                console.log("📨 Message type:", data.type || 'unknown');
                console.log("📨 Message content:", data.content ? 'present' : 'missing');
                if (onMessageCallback) {
                    onMessageCallback(data);
                }
            } catch (error) {
                console.error("Error parsing group WebSocket message:", error);
                // Try to handle as plain text if JSON parsing fails
                if (onMessageCallback) {
                    onMessageCallback({ content: event.data });
                }
            }
        };

        groupSocket.onclose = (event) => {
            console.log("🔌 Group WebSocket closed for group:", groupId, "Code:", event.code, "Reason:", event.reason);
            console.log("🔌 WebSocket close details:", {
                code: event.code,
                reason: event.reason,
                wasClean: event.wasClean,
                target: event.target
            });
            isGroupConnected.value = false;
            groupSocket = null;
        };

        groupSocket.onerror = (error) => {
            console.error("🔌 Group WebSocket error:", error);
            isGroupConnected.value = false;
        };
    };

    // Send group message with encryption for secret groups
    const sendGroupMessage = async (messageData, groupId, isSecretGroup = false) => {
        console.log("📤 Attempting to send group message:", messageData);
        console.log("🔌 Group WebSocket state:", groupSocket ? groupSocket.readyState : "null");
        console.log("🔐 Is secret group:", isSecretGroup);
        
        if (!groupSocket || groupSocket.readyState !== WebSocket.OPEN) {
            console.error("🔌 Group WebSocket is not connected. State:", groupSocket ? groupSocket.readyState : "null");
            return false;
        }

        try {
            let finalMessageData = messageData;
            console.log("📤 Original message data:", messageData);
            
            // Handle secret group messages with symmetric keys for each user
            console.log('🔐 Checking if message should be encrypted. isSecretGroup:', isSecretGroup);
            if (isSecretGroup) {
                console.log('🔐 Processing secret group message:', groupId);
                console.log('🔐 Original message content:', messageData.content);
                
                try {
                    // Use the new sendSecretGroupMessage function
                    const { sendSecretGroupMessage } = useSecretGroupE2EE();
                    console.log('🔐 About to call sendSecretGroupMessage...');
                    const secretMessagePayload = await sendSecretGroupMessage(messageData.content, groupId);
                    console.log('🔐 sendSecretGroupMessage completed successfully');
                    
                    console.log('🔐 Secret message payload:', {
                        content_length: secretMessagePayload.content.length,
                        users_count: Object.keys(secretMessagePayload.users_symmetric_keys).length,
                        users: Object.keys(secretMessagePayload.users_symmetric_keys)
                    });
                    
                    finalMessageData = {
                        ...messageData,
                        content: secretMessagePayload.content,
                        users_symmetric_keys: secretMessagePayload.users_symmetric_keys,
                        is_encrypted: true
                    };
                    console.log('✅ Secret group message prepared with symmetric keys');
                } catch (error) {
                    console.error('❌ Error in sendSecretGroupMessage:', error);
                    console.error('❌ Error details:', {
                        message: error.message,
                        stack: error.stack
                    });
                    throw error;
                }
            } else {
                console.log('🔐 Not a secret group, sending regular message');
            }

            console.log("📤 Sending group WebSocket message:", finalMessageData);
            console.log("📤 Raw message being sent:", JSON.stringify(finalMessageData));
            console.log("📤 WebSocket ready state:", groupSocket.readyState);
            console.log("📤 WebSocket URL:", groupSocket.url);
            
            try {
                groupSocket.send(JSON.stringify(finalMessageData));
                console.log("✅ Group message sent successfully");
            } catch (sendError) {
                console.error("❌ Error sending WebSocket message:", sendError);
                throw sendError;
            }
            return true;
        } catch (error) {
            console.error("❌ Error sending group message:", error);
            return false;
        }
    };

    // Close group connection
    const closeGroupConnection = () => {
        if (groupSocket) {
            groupSocket.close();
            groupSocket = null;
            isGroupConnected.value = false;
        }
    };

    const getGroupConnectionStatus = () => {
        return {
            isConnected: groupSocket ? groupSocket.readyState === WebSocket.OPEN : false,
            readyState: groupSocket ? groupSocket.readyState : null,
            socket: groupSocket
        };
    };

    const addGroupMessage = (message) => {
        groupMessages.value.push(message);
    };

    const clearGroupMessages = () => {
        groupMessages.value = [];
    };

    // Load group users from API
    const loadGroupUsers = async (groupId) => {
        try {
            console.log('👥 Loading group users for group:', groupId);
            
            // Check if this is a secret group
            const isSecretGroup = groupStore.currentGroup?.type === 'secret';
            
            let response;
            if (isSecretGroup) {
                response = await axiosInstance.get(`/api/secret-group/get/${groupId}/members`);
            } else {
                response = await axiosInstance.get(`/api/group/${groupId}/members`);
            }
            
            console.log('👥 Group users response:', response.data);
            
            // Handle different response structures
            let users = [];
            if (response.data && Array.isArray(response.data)) {
                // Direct array response
                users = response.data;
            } else if (response.data && response.data.members && Array.isArray(response.data.members)) {
                // Nested members array
                users = response.data.members;
            } else if (response.data && response.data.users && Array.isArray(response.data.users)) {
                // Nested users array
                users = response.data.users;
            } else if (response.data && typeof response.data === 'object' && !Array.isArray(response.data)) {
                // Object with user IDs as keys
                users = Object.values(response.data);
                console.log('✅ Converted object response to array with', users.length, 'users');
            } else {
                console.warn('⚠️ Unexpected response structure for group users:', response.data);
                users = [];
            }
            
            // Ensure users is an array and has the expected structure
            if (!Array.isArray(users)) {
                console.error('❌ Users is not an array:', users);
                groupUsers.value = {};
                return {};
            }
            
            groupUsers.value = users.reduce((acc, user) => {
                const userId = user.user_id || user.id || user._id;
                if (userId) {
                    acc[userId] = user;
                }
                return acc;
            }, {});
            
            console.log('✅ Loaded', Object.keys(groupUsers.value).length, 'group users');
            return groupUsers.value;
        } catch (error) {
            console.error('❌ Failed to load group users:', error);
            console.error('❌ Error details:', {
                message: error.message,
                response: error.response?.data,
                status: error.response?.status
            });
            groupUsers.value = {};
            return {};
        }
    };

    const getUsernameBySenderId = (senderId) => {
        const user = groupUsers.value[senderId];
        return user ? user.username : 'Unknown User';
    };

    const getAvatarBySenderId = (senderId) => {
        const user = groupUsers.value[senderId];
        return user ? user.avatar_url : null;
    };

    // Load group messages from API with pagination and decryption for secret groups
    const loadGroupMessages = async (groupId, page = 1, limit = 20, isSecretGroup = false) => {
        if (isLoadingMessages.value) return;

        try {
            isLoadingMessages.value = true;
            console.log('📥 Loading group messages for group:', groupId, 'page:', page, 'limit:', limit, 'isSecret:', isSecretGroup);
            
            let response;
            if (isSecretGroup) {
                response = await axiosInstance.get(`/api/secret-group/get/${groupId}/messages`, {
                    params: { page, limit }
                });
            } else {
                response = await axiosInstance.get(`/api/group/get/${groupId}/messages`, {
                params: { page, limit }
            });
            }
            
            console.log('📥 Group messages response:', response.data);
            
            // Handle the response structure: { messages: [...] }
            const messagesArray = response.data?.messages || response.data || [];
            
            if (Array.isArray(messagesArray)) {
                console.log('📥 Processing', messagesArray.length, 'messages');
                // Transform API messages to our format
                const transformedMessages = await Promise.all(messagesArray.map(async (msg, index) => {
                    console.log(`📥 Processing message ${index + 1}:`, {
                        id: msg.id,
                        content: msg.content,
                        sender_id: msg.sender_id,
                        is_encrypted: msg.is_encrypted,
                        has_symmetric_keys: !!msg.users_symmetric_keys,
                        symmetric_keys_count: msg.users_symmetric_keys ? Object.keys(msg.users_symmetric_keys).length : 0,
                        symmetric_keys_users: msg.users_symmetric_keys ? Object.keys(msg.users_symmetric_keys) : []
                    });
                    let content = msg.content;
                    
                    // Handle empty content in secret groups (old messages or failed encryption)
                    if (isSecretGroup && !msg.content && !msg.is_encrypted) {
                        console.log('⚠️ Found empty message in secret group - likely old message or failed encryption');
                        content = '[Old message - no content available]';
                    }
                    
                    // For secret groups, check if we have our private key
                    if (isSecretGroup) {
                        const { hasSecretGroupKeys } = useKeyPair();
                        const hasKeys = await hasSecretGroupKeys(groupId);
                        console.log('🔐 User has private keys for secret group:', hasKeys);
                    }
                    
                    // Decrypt message for secret groups
                    if (isSecretGroup && msg.is_encrypted) {
                        try {
                            console.log('🔐 Decrypting message for secret group:', groupId);
                            console.log('🔐 Message structure:', {
                                id: msg.id,
                                content: msg.content,
                                is_encrypted: msg.is_encrypted,
                                users_symmetric_keys: msg.users_symmetric_keys ? Object.keys(msg.users_symmetric_keys) : null
                            });
                            
                            // Check if this message has symmetric keys for users
                            if (msg.users_symmetric_keys) {
                                console.log('🔐 Message has symmetric keys for users, processing...');
                                
                                // Get current user ID to find their encrypted symmetric key
                                const currentUserId = userStore.user_id;
                                console.log('🔐 Current user ID:', currentUserId);
                                console.log('🔐 Available user IDs in symmetric keys:', Object.keys(msg.users_symmetric_keys));
                                
                                const userEncryptedKey = msg.users_symmetric_keys[currentUserId];
                                
                                if (userEncryptedKey) {
                                    console.log('🔐 Found encrypted symmetric key for current user');
                                    
                                    // Decrypt the symmetric key using our private key
                                    const { decryptGroupSymmetricKey } = useSecretGroupE2EE();
                                    const { getSecretGroupPrivateKey } = useKeyPair();
                                    
                                    const privateKey = await getSecretGroupPrivateKey(groupId);
                                    if (privateKey) {
                                        const symmetricKey = await decryptGroupSymmetricKey(userEncryptedKey, privateKey);
                                        
                                        // Now decrypt the message content with the symmetric key
                                        content = await decryptGroupMessage(msg.content, symmetricKey);
                                        console.log('✅ Message decrypted using symmetric key');
                                    } else {
                                        console.error('❌ No private key available for group');
                                        content = '[Encrypted message - no private key]';
                                    }
                                } else {
                                    console.warn('⚠️ No encrypted symmetric key found for current user');
                                    console.warn('⚠️ Current user ID:', currentUserId);
                                    console.warn('⚠️ Available keys:', Object.keys(msg.users_symmetric_keys));
                                    content = '[Encrypted message - no symmetric key]';
                                }
                            } else {
                                // Fallback to old method (should not happen with new architecture)
                                console.warn('⚠️ Message does not have symmetric keys, using fallback');
                                content = '[Encrypted message - no symmetric keys]';
                            }
                        } catch (decryptError) {
                            console.error('❌ Failed to decrypt message:', decryptError);
                            content = '[Encrypted message - unable to decrypt]';
                        }
                    }
                    
                    return {
                    id: msg.id || msg._id,
                    sender_id: msg.sender_id,
                        content: content,
                    message_type: msg.type === 'text' ? 1 : msg.message_type || 1,
                    created_at: msg.created_at,
                    sender_name: getUsernameBySenderId(msg.sender_id),
                        sender_avatar: getAvatarBySenderId(msg.sender_id),
                        is_encrypted: msg.is_encrypted || false,
                        users_symmetric_keys: msg.users_symmetric_keys
                    };
                }));
                
                // For pagination, we'll assume there are more messages if we got a full page
                const hasMore = messagesArray.length >= limit;
                
                // Update pagination state
                currentPage.value = page;
                hasMoreMessages.value = hasMore;
                
                // If it's the first page, replace messages; otherwise, prepend to existing messages
                if (page === 1) {
                    groupMessages.value = transformedMessages;
                } else {
                    // Prepend older messages to the beginning (for pagination)
                    groupMessages.value.unshift(...transformedMessages);
                }
                
                console.log('✅ Loaded', transformedMessages.length, 'group messages (page', page, ')');
            } else {
                console.log('📥 No group messages found or invalid response format');
                if (page === 1) {
                    groupMessages.value = [];
                }
                hasMoreMessages.value = false;
            }
        } catch (error) {
            console.error('❌ Failed to load group messages:', error);
            showError('Failed to load group messages. Please try again.');
            if (page === 1) {
                groupMessages.value = [];
            }
        } finally {
            isLoadingMessages.value = false;
        }
    };

    // Load next page of group messages
    const loadNextGroupPage = async (groupId, isSecretGroup = false) => {
        if (!hasMoreMessages.value || isLoadingMessages.value) return;
        
        const nextPage = currentPage.value + 1;
        return await loadGroupMessages(groupId, nextPage, pageLimit.value, isSecretGroup);
    };

    // Load initial group messages with encryption setup for secret groups
    const loadInitialGroupMessages = async (groupId, isSecretGroup = false) => {
        // Reset pagination state
        currentPage.value = 1;
        hasMoreMessages.value = true;
        isLoadingMessages.value = false;
        
        // For secret groups, verify we have private keys (symmetric keys are loaded from messages)
        if (isSecretGroup) {
            try {
                console.log('🔐 Verifying secret group encryption setup:', groupId);
                const { hasSecretGroupKeys } = useKeyPair();
                const hasKeys = await hasSecretGroupKeys(groupId);
                
                if (!hasKeys) {
                    console.warn('⚠️ No private keys available for secret group, encryption may not work');
                } else {
                    console.log('✅ Private keys available for secret group');
                }
            } catch (error) {
                console.error('❌ Failed to verify secret group encryption setup:', error);
                // Continue without encryption
            }
        }
        
        return await loadGroupMessages(groupId, 1, pageLimit.value, isSecretGroup);
    };

    // Handle incoming group message with decryption for secret groups
    const handleIncomingGroupMessage = async (data, groupId, isSecretGroup = false) => {
        console.log('📨 Received group WebSocket message:', data);
        
        // Parse the group message according to your backend structure
        const groupMessage = parseIncomingGroupMessage(data);
        console.log('📨 Parsed group message:', groupMessage);

        let content = groupMessage.content;
        
        // Decrypt message for secret groups
        if (isSecretGroup && groupMessage.is_encrypted) {
            try {
                console.log('🔐 Decrypting incoming message for secret group:', groupId);
                
                // Check if this message has symmetric keys for users
                if (groupMessage.users_symmetric_keys) {
                    console.log('🔐 Message has symmetric keys for users, processing...');
                    
                    // Get current user ID to find their encrypted symmetric key
                    const currentUserId = groupStore.currentGroup?.current_user_id || 'current_user';
                    const userEncryptedKey = groupMessage.users_symmetric_keys[currentUserId];
                    
                    if (userEncryptedKey) {
                        console.log('🔐 Found encrypted symmetric key for current user');
                        
                        // Decrypt the symmetric key using our private key
                        const { decryptGroupSymmetricKey } = useSecretGroupE2EE();
                        const { getSecretGroupPrivateKey } = useKeyPair();
                        
                        const privateKey = await getSecretGroupPrivateKey(groupId);
                        if (privateKey) {
                            const symmetricKey = await decryptGroupSymmetricKey(userEncryptedKey, privateKey);
                            
                            // Store the symmetric key for this message
                            const { storeGroupSymmetricKey } = useSecretGroupE2EE();
                            await storeGroupSymmetricKey(groupId, symmetricKey);
                            
                            // Now decrypt the message content with the symmetric key
                            content = await decryptGroupMessage(groupMessage.content, symmetricKey);
                            console.log('✅ Message decrypted using symmetric key');
                        } else {
                            console.error('❌ No private key available for group');
                            content = '[Encrypted message - no private key]';
                        }
                    } else {
                        console.warn('⚠️ No encrypted symmetric key found for current user');
                        content = '[Encrypted message - no symmetric key]';
                    }
                } else {
                    // Fallback to old method (should not happen with new architecture)
                    console.warn('⚠️ Message does not have symmetric keys, using fallback');
                    content = '[Encrypted message - no symmetric keys]';
                }
            } catch (decryptError) {
                console.error('❌ Failed to decrypt incoming message:', decryptError);
                content = '[Encrypted message - unable to decrypt]';
            }
        }

        // Create a message object compatible with the group messages
        const message = {
            id: groupMessage.id || `temp-${Date.now()}-${Math.random()}`,
            sender_id: groupMessage.sender_id,
            content: content,
            message_type: groupMessage.message_type || 1, // Default to text message
            created_at: groupMessage.created_at || new Date().toISOString(),
            // Group-specific fields - use group users data
            sender_name: getUsernameBySenderId(groupMessage.sender_id),
            sender_avatar: getAvatarBySenderId(groupMessage.sender_id),
            is_encrypted: groupMessage.is_encrypted || false,
            users_symmetric_keys: groupMessage.users_symmetric_keys
        };

        // Check for duplicate messages
        const duplicateMessage = groupMessages.value.find(
            (msg) =>
                msg.content === message.content &&
                msg.sender_id === message.sender_id &&
                (message.id ? msg.id === message.id : 
                 Math.abs(new Date(msg.created_at) - new Date(message.created_at)) < 100)
        );

        if (duplicateMessage) {
            console.log('🔄 Duplicate group message detected, skipping:', message.content);
            return;
        }

        // Add the message to the group messages
        console.log('➕ Adding new group message:', {
            id: message.id,
            content: message.content,
            sender_id: message.sender_id,
            sender_name: message.sender_name,
            created_at: message.created_at,
            is_encrypted: message.is_encrypted
        });
        addGroupMessage(message);
    };

    // Parse incoming group message data
    const parseIncomingGroupMessage = (data) => {
        // Handle the GroupMessage struct format from your backend
        if (data.message_type && data.sender_id && data.content) {
            return {
                message_type: data.message_type,
                sender_id: data.sender_id,
                content: data.content,
                id: data.id,
                created_at: data.created_at,
                sender_name: data.sender_name,
                sender_avatar: data.sender_avatar,
                is_encrypted: data.is_encrypted,
                users_symmetric_keys: data.users_symmetric_keys
            };
        }

        // If the backend sends the message object directly
        if (data.content && typeof data.content === "object") {
            return {
                ...data.content,
                users_symmetric_keys: data.users_symmetric_keys
            };
        }

        // If the backend sends content as a JSON string
        if (typeof data.content === "string" && data.content.startsWith("{")) {
            try {
                const parsed = JSON.parse(data.content);
                return {
                    ...parsed,
                    users_symmetric_keys: data.users_symmetric_keys
                };
            } catch (e) {
                return { 
                    content: data.content,
                    users_symmetric_keys: data.users_symmetric_keys
                };
            }
        }

        // Default case
        return data;
    };

    return {
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
    };
} 