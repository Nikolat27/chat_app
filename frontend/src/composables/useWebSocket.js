import { ref } from "vue";

let currentSocket = null;

export function useWebSocket() {
    const isConnected = ref(false);

    // Establish WebSocket connection
    const establishConnection = (chatData, onMessageCallback) => {
        console.log("Establishing WebSocket connection with data:", chatData);
        
        // Close existing connection if any
        if (currentSocket) {
            console.log("Closing existing WebSocket connection");
            currentSocket.close();
            currentSocket = null;
            // Wait a bit for the connection to close
            setTimeout(() => {
                createNewConnection(chatData, onMessageCallback);
            }, 100);
        } else {
            createNewConnection(chatData, onMessageCallback);
        }
    };

    const createNewConnection = (chatData, onMessageCallback) => {
        const { chatId, senderId, receiverId, backendBaseUrl, isSecretChat, isGroupChat, groupId } = chatData;

        if (!chatId || !senderId || !backendBaseUrl) {
            console.error("Missing required data for WebSocket connection:", { chatId, senderId, backendBaseUrl });
            return;
        }

        // Use different WebSocket URL based on chat type
        let wsUrl;
        if (isGroupChat && groupId) {
            wsUrl = `${backendBaseUrl.replace(/^http/, "ws")}/api/websocket/group/add/${groupId}?sender_id=${senderId}`;
        } else if (isSecretChat) {
            wsUrl = `${backendBaseUrl.replace(/^http/, "ws")}/api/websocket/secret-chat/add/${chatId}?sender_id=${senderId}&receiver_id=${receiverId}`;
        } else {
            wsUrl = `${backendBaseUrl.replace(/^http/, "ws")}/api/websocket/chat/add/${chatId}?sender_id=${senderId}&receiver_id=${receiverId}`;
        }

        console.log("Creating WebSocket connection to:", wsUrl);
        currentSocket = new WebSocket(wsUrl);

        currentSocket.onopen = () => {
            console.log("WebSocket connected for chat:", chatId);
            isConnected.value = true;
        };

        currentSocket.onmessage = (event) => {
            console.log("📨 Received WebSocket message:", event.data);
            try {
                const data = JSON.parse(event.data);
                if (onMessageCallback) {
                    onMessageCallback(data);
                }
            } catch (error) {
                console.error("Error parsing WebSocket message:", error);
                // Try to handle as plain text if JSON parsing fails
                if (onMessageCallback) {
                    onMessageCallback({ content: event.data });
                }
            }
        };

        currentSocket.onclose = (event) => {
            console.log("WebSocket closed for chat:", chatId, "Code:", event.code, "Reason:", event.reason);
            isConnected.value = false;
            currentSocket = null;
        };

        currentSocket.onerror = (error) => {
            console.error("WebSocket error:", error);
            isConnected.value = false;
        };
    };

    // Send message through WebSocket
    const sendMessage = (messageData) => {
        console.log("Attempting to send message:", messageData);
        console.log("WebSocket state:", currentSocket ? currentSocket.readyState : "null");
        
        if (!currentSocket || currentSocket.readyState !== WebSocket.OPEN) {
            console.error("WebSocket is not connected. State:", currentSocket ? currentSocket.readyState : "null");
            return false;
        }

        try {
            console.log("📤 Sending WebSocket message:", messageData);
            currentSocket.send(messageData);
            console.log("Message sent successfully");
            return true;
        } catch (error) {
            console.error("Error sending message:", error);
            return false;
        }
    };

    // Close current connection
    const closeConnection = () => {
        if (currentSocket) {
            currentSocket.close();
            currentSocket = null;
            isConnected.value = false;
        }
    };

    // Get connection status
    const getConnectionStatus = () => {
        return {
            isConnected: isConnected.value,
            readyState: currentSocket
                ? currentSocket.readyState
                : WebSocket.CLOSED,
        };
    };

    return {
        isConnected,
        establishConnection,
        sendMessage,
        closeConnection,
        getConnectionStatus,
    };
}
