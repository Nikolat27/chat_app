import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useE2EEStore = defineStore('e2ee', () => {
    // Global symmetric keys storage
    const symmetricKeys = ref(new Map());

    // Store symmetric key
    const storeSymmetricKey = (chatId, key) => {
        console.log('🔐 [Store] Storing symmetric key for chat:', chatId);
        symmetricKeys.value.set(chatId, key);
        console.log('🔍 [Store] Keys in memory after storing:', Array.from(symmetricKeys.value.keys()));
    };

    // Get symmetric key
    const getSymmetricKey = (chatId) => {
        console.log('🔍 [Store] Checking symmetric key for chat:', chatId);
        console.log('🔍 [Store] Keys in memory:', Array.from(symmetricKeys.value.keys()));
        
        if (symmetricKeys.value.has(chatId)) {
            console.log('✅ [Store] Found symmetric key in memory for chat:', chatId);
            return symmetricKeys.value.get(chatId);
        }
        
        console.log('❌ [Store] No symmetric key found in memory for chat:', chatId);
        return null;
    };

    // Check if symmetric key exists
    const hasSymmetricKey = (chatId) => {
        return symmetricKeys.value.has(chatId);
    };

    // Clear symmetric key
    const clearSymmetricKey = (chatId) => {
        console.log('🗑️ [Store] Clearing symmetric key for chat:', chatId);
        symmetricKeys.value.delete(chatId);
        console.log('🔍 [Store] Keys in memory after clearing:', Array.from(symmetricKeys.value.keys()));
    };

    // Clear all symmetric keys
    const clearAllSymmetricKeys = () => {
        console.log('🗑️ [Store] Clearing all symmetric keys');
        symmetricKeys.value.clear();
        console.log('🔍 [Store] Keys in memory after clearing all:', Array.from(symmetricKeys.value.keys()));
    };

    return {
        storeSymmetricKey,
        getSymmetricKey,
        hasSymmetricKey,
        clearSymmetricKey,
        clearAllSymmetricKeys
    };
}); 