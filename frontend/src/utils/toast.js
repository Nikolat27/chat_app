import { useToast } from "vue-toastification";

let toastInstance;
export function getToast() {
    if (!toastInstance) {
        toastInstance = useToast();
    }
    return toastInstance;
}

export function showError(message) {
    getToast().error(message, {
        timeout: 4000,
        position: "top-right",
    });
}

export function showMessage(message) {
    getToast().success(message, {
        timeout: 3000,
        position: "top-right",
    });
}

export function showInfo(message) {
    getToast().info(message, {
        timeout: 3500,
        position: "top-right",
    });
}
