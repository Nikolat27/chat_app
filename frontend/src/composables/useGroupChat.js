import { ref } from 'vue';
import { showError } from '../utils/toast';
import axiosInstance from '../axiosInstance';
import { useGroupStore } from '../stores/groups';

let groupSocket = null;
let groupUsers = ref({});

export function useGroupChat() {
    const groupStore = useGroupStore();
    const isGroupConnected = ref(false);
    const groupMessages = ref([]);
    const newGroupMessage = ref('');

    // Establish group WebSocket connection
    const establishGroupConnection = (groupData, onMessageCallback) => {
        console.log("🔌 Establishing group WebSocket connection with data:", groupData);
        
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

        const wsUrl = `${backendBaseUrl.replace(/^http/, "ws")}/api/websocket/group/add/${groupId}?sender_id=${senderId}`;
        console.log("🔌 Creating group WebSocket connection to:", wsUrl);
        groupSocket = new WebSocket(wsUrl);

        groupSocket.onopen = () => {
            console.log("🔌 Group WebSocket connected for group:", groupId);
            isGroupConnected.value = true;
        };

        groupSocket.onmessage = (event) => {
            console.log("📨 Received group WebSocket message:", event.data);
            try {
                const data = JSON.parse(event.data);
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
            isGroupConnected.value = false;
            groupSocket = null;
        };

        groupSocket.onerror = (error) => {
            console.error("🔌 Group WebSocket error:", error);
            isGroupConnected.value = false;
        };
    };

    // Send group message
    const sendGroupMessage = (messageData) => {
        console.log("📤 Attempting to send group message:", messageData);
        console.log("🔌 Group WebSocket state:", groupSocket ? groupSocket.readyState : "null");
        
        if (!groupSocket || groupSocket.readyState !== WebSocket.OPEN) {
            console.error("🔌 Group WebSocket is not connected. State:", groupSocket ? groupSocket.readyState : "null");
            return false;
        }

        try {
            console.log("📤 Sending group WebSocket message:", messageData);
            groupSocket.send(messageData);
            console.log("✅ Group message sent successfully");
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

    // Get group connection status
    const getGroupConnectionStatus = () => {
        return {
            isConnected: isGroupConnected.value,
            readyState: groupSocket
                ? groupSocket.readyState
                : WebSocket.CLOSED,
        };
    };

    // Add group message
    const addGroupMessage = (message) => {
        groupMessages.value.push(message);
    };

    // Clear group messages
    const clearGroupMessages = () => {
        groupMessages.value = [];
    };

    // Load group users from API
    const loadGroupUsers = async (groupId) => {
        try {
            console.log('👥 Loading group users for group:', groupId);
            console.log('👥 Making API call to:', `/api/group/get/${groupId}/users`);
            
            const response = await axiosInstance.get(`/api/group/get/${groupId}/users`);
            console.log('👥 Group users response:', response.data);
            
            if (response.data && typeof response.data === 'object') {
                groupStore.setGroupUsers(response.data);
                console.log('✅ Loaded', Object.keys(response.data).length, 'group users');
            } else {
                console.log('👥 No group users found or invalid response format');
                groupStore.setGroupUsers({});
            }
        } catch (error) {
            console.error('❌ Failed to load group users:', error);
            console.error('❌ Error details:', error.response?.data);
            showError('Failed to load group users. Please try again.');
            groupStore.setGroupUsers({});
        }
    };

    // Get username by sender ID
    const getUsernameBySenderId = (senderId) => {
        const user = groupStore.getGroupUsers()[senderId];
        return user?.username || 'Unknown User';
    };

    // Get avatar by sender ID
    const getAvatarBySenderId = (senderId) => {
        const user = groupStore.getGroupUsers()[senderId];
        if (!user?.profile_url) return null;
        
        // Construct full avatar URL
        const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL;
        return `${backendBaseUrl}/static/${user.profile_url}`;
    };

    // Load group messages from API
    const loadGroupMessages = async (groupId) => {
        try {
            console.log('📥 Loading group messages for group:', groupId);
            const response = await axiosInstance.get(`/api/group/get/${groupId}/messages`);
            console.log('📥 Group messages response:', response.data);
            
            // Handle the response structure: { messages: [...] }
            const messagesArray = response.data?.messages || response.data || [];
            
            if (Array.isArray(messagesArray)) {
                // Transform API messages to our format
                const transformedMessages = messagesArray.map(msg => ({
                    id: msg.id || msg._id,
                    sender_id: msg.sender_id,
                    content: msg.content,
                    message_type: msg.type === 'text' ? 1 : msg.message_type || 1,
                    created_at: msg.created_at,
                    sender_name: getUsernameBySenderId(msg.sender_id),
                    sender_avatar: getAvatarBySenderId(msg.sender_id)
                }));
                
                groupMessages.value = transformedMessages;
                console.log('✅ Loaded', transformedMessages.length, 'group messages');
            } else {
                console.log('📥 No group messages found or invalid response format');
                groupMessages.value = [];
            }
        } catch (error) {
            console.error('❌ Failed to load group messages:', error);
            showError('Failed to load group messages. Please try again.');
            groupMessages.value = [];
        }
    };

    // Handle incoming group message
    const handleIncomingGroupMessage = (data) => {
        console.log('📨 Received group WebSocket message:', data);
        
        // Parse the group message according to your backend structure
        const groupMessage = parseIncomingGroupMessage(data);
        console.log('📨 Parsed group message:', groupMessage);

        // Create a message object compatible with the group messages
        const message = {
            id: groupMessage.id || `temp-${Date.now()}-${Math.random()}`,
            sender_id: groupMessage.sender_id,
            content: groupMessage.content,
            message_type: groupMessage.message_type || 1, // Default to text message
            created_at: groupMessage.created_at || new Date().toISOString(),
            // Group-specific fields - use group users data
            sender_name: getUsernameBySenderId(groupMessage.sender_id),
            sender_avatar: getAvatarBySenderId(groupMessage.sender_id)
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
            created_at: message.created_at
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
                sender_avatar: data.sender_avatar
            };
        }

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

        // Default case
        return data;
    };

    return {
        isGroupConnected,
        groupMessages,
        newGroupMessage,
        establishGroupConnection,
        sendGroupMessage,
        closeGroupConnection,
        getGroupConnectionStatus,
        addGroupMessage,
        clearGroupMessages,
        loadGroupUsers,
        loadGroupMessages,
        getUsernameBySenderId,
        getAvatarBySenderId,
        handleIncomingGroupMessage
    };
} 