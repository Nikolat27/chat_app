<template>
    <div class="h-full">
        <!-- Header -->
        <div class="mb-6">
            <div class="flex items-center gap-3 mb-2">
                <span class="material-icons text-green-600 text-xl">group</span>
                <h3 class="text-lg font-bold text-gray-800">Groups</h3>
                <div class="flex-1"></div>
                <div class="flex items-center space-x-3">
                    <button
                        @click="handleOpenApprovals"
                        class="px-3 py-1 text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-lg text-sm font-medium transition-all duration-200 cursor-pointer flex items-center space-x-1"
                        title="View pending approvals"
                    >
                        <span class="material-icons text-sm"
                            >pending_actions</span
                        >
                        <span>Approvals</span>
                    </button>
                    <span
                        class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-full"
                    >
                    {{ groups ? groups.length : 0 }} groups
                </span>
                </div>
            </div>
            <p class="text-sm text-gray-600">
                Create and join group conversations with multiple participants
            </p>
        </div>

        <!-- Action Buttons -->
        <div class="space-y-3 mb-6">
            <!-- Join Group Button -->
            <button
                @click="showJoinGroupModal = true"
                class="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3 px-4 rounded-xl transition-all duration-200 flex items-center justify-center gap-3 shadow-lg hover:shadow-xl transform hover:scale-[1.02] cursor-pointer"
            >
                <span class="material-icons text-lg">group_add</span>
                <span>Join Group</span>
            </button>

            <!-- Create Group Button -->
            <button
                @click="showCreateGroupModal = true"
                class="w-full bg-green-500 hover:bg-green-600 text-white font-semibold py-3 px-4 rounded-xl transition-all duration-200 flex items-center justify-center gap-3 shadow-lg hover:shadow-xl transform hover:scale-[1.02] cursor-pointer"
            >
                <span class="material-icons text-lg">add_circle</span>
                <span>Create Group</span>
            </button>

            <!-- Create Secret Group Button -->
            <button
                @click="showCreateSecretGroupModal = true"
                class="w-full bg-purple-500 hover:bg-purple-600 text-white font-semibold py-3 px-4 rounded-xl transition-all duration-200 flex items-center justify-center gap-3 shadow-lg hover:shadow-xl transform hover:scale-[1.02] cursor-pointer"
            >
                <span class="material-icons text-lg">lock</span>
                <span>Create Secret Group</span>
            </button>
        </div>

        <!-- Loading State -->
        <div v-if="isLoading" class="flex items-center justify-center py-12">
            <div
                class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-500"
            ></div>
            <span class="ml-3 text-gray-600">Loading groups...</span>
        </div>

        <!-- Groups List -->
        <div v-else-if="groups && groups.length > 0" class="space-y-3">
            <h4 class="text-sm font-semibold text-gray-700 mb-3">
                Your Groups ({{ groups.length }})
            </h4>
            <div
                v-for="group in groups"
                :key="group.id"
                class="group bg-white rounded-2xl shadow-sm border border-gray-200 cursor-pointer hover:shadow-lg hover:border-green-300 transition-all duration-300 overflow-hidden"
            >
                <div class="p-4">
                    <div class="flex items-center gap-4">
                        <!-- Group Avatar -->
                        <div class="relative">
                            <img
                                :src="
                                    group.avatar_url
                                        ? `${backendBaseUrl}/static/${group.avatar_url}`
                                        : '/src/assets/default-avatar.jpg'
                                "
                                alt="Group Avatar"
                                class="w-12 h-12 rounded-full object-cover border-2 border-green-300 shadow-sm group-hover:border-green-400 transition-colors duration-200 select-none pointer-events-none"
                            />
                            <div
                                class="absolute -bottom-1 -right-1 w-5 h-5 bg-green-500 rounded-full flex items-center justify-center"
                            >
                                <span class="material-icons text-white text-xs"
                                    >group</span
                                >
                            </div>
                        </div>

                        <!-- Group Info -->
                        <div class="flex-1 min-w-0">
                            <div class="flex items-center gap-2 mb-1">
                                <span
                                    class="font-semibold text-gray-800 flex-shrink-0"
                                >
                                    {{ group.name }}
                                </span>
                                <div 
                                    :class="
                                        (group.is_secret || group.type === 'secret')
                                        ? 'bg-purple-100 text-purple-700 border-purple-200' 
                                            : group.type === 'private'
                                            ? 'bg-orange-100 text-orange-700 border-orange-200'
                                            : 'bg-green-100 text-green-700 border-green-200'
                                    "
                                    class="px-2 py-1 rounded-full text-xs font-medium border flex items-center gap-1"
                                >
                                    <span class="material-icons text-xs">
                                        {{
                                            (group.is_secret || group.type === 'secret')
                                                ? "lock"
                                                : group.type === "private"
                                                ? "lock_outline"
                                                : "group"
                                        }}
                                    </span>
                                    {{
                                        (group.is_secret || group.type === 'secret')
                                            ? "Secret"
                                            : group.type === "private"
                                            ? "Private"
                                            : "Public"
                                    }}
                                </div>
                                <!-- Owner Badge -->
                                <span
                                    v-if="group.owner_id === userStore.user_id"
                                    class="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded-full font-medium border border-blue-200"
                                >
                                    Owner
                                </span>
                                <!-- Admin Badge -->
                                <span
                                    v-else-if="group.admins && group.admins.includes(userStore.user_id)"
                                    class="px-2 py-1 bg-green-100 text-green-700 text-xs rounded-full font-medium border border-green-200"
                                >
                                    Admin
                                </span>
                            </div>
                            <p class="text-sm text-gray-600 truncate">
                                {{ group.description || "No description" }}
                            </p>
                            <div class="flex items-center gap-2 mt-1">
                                <span class="text-xs text-gray-500">
                                    {{ group.member_count || 0 }} members
                                </span>
                                <span class="text-xs text-gray-400">•</span>
                                <span class="text-xs text-gray-500">
                                    {{ group.role || "Member" }}
                                </span>
                            </div>
                        </div>

                        <!-- Action Buttons -->
                        <div
                            class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-all duration-200"
                        >
                            <button
                                @click.stop="handleGroupClick(group)"
                                class="w-8 h-8 text-blue-500 hover:bg-blue-50 rounded-full hover:text-blue-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 flex items-center justify-center"
                                title="Open group chat"
                            >
                                <span class="material-icons text-sm">chat</span>
                            </button>
                            <button
                                v-if="group.owner_id === userStore.user_id"
                                @click.stop="handleEditGroup(group)"
                                class="w-8 h-8 text-green-500 hover:bg-green-50 rounded-full hover:text-green-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 flex items-center justify-center"
                                title="Edit group"
                            >
                                <span class="material-icons text-sm">edit</span>
                            </button>
                            <button
                                @click.stop="handleManageMembers(group)"
                                class="w-8 h-8 text-blue-500 hover:bg-blue-50 rounded-full hover:text-blue-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 flex items-center justify-center"
                                title="Group users"
                            >
                                <span class="material-icons text-sm">people</span>
                            </button>
                            <button
                                v-if="group.invite_link"
                                @click.stop="handleCopyInviteLink(group)"
                                class="w-8 h-8 text-blue-600 hover:bg-blue-50 rounded-full hover:text-blue-700 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 flex items-center justify-center"
                                title="Copy invite link"
                            >
                                <span class="material-icons text-sm"
                                    >content_copy</span
                                >
                            </button>
                            <button
                                v-if="group.owner_id !== userStore.user_id"
                                @click.stop="handleLeaveGroup(group)"
                                class="w-8 h-8 text-red-500 hover:bg-red-50 rounded-full hover:text-red-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 flex items-center justify-center"
                                title="Leave group"
                            >
                                <span class="material-icons text-sm"
                                    >exit_to_app</span
                                >
                            </button>
                            <button
                                v-if="group.owner_id === userStore.user_id"
                                @click.stop="handleDeleteGroup(group)"
                                class="w-8 h-8 text-red-600 hover:bg-red-50 rounded-full hover:text-red-700 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 flex items-center justify-center"
                                title="Delete group"
                            >
                                <span class="material-icons text-sm"
                                    >delete</span
                                >
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Empty State -->
        <div v-else class="text-center py-12 px-6">
            <div class="mb-4">
                <span class="material-icons text-6xl text-gray-300">group</span>
            </div>
            <h4 class="text-lg font-semibold text-gray-600 mb-2">
                No Groups Yet
            </h4>
            <p class="text-sm text-gray-500 mb-6 leading-relaxed">
                Start by creating a new group or joining an existing one to
                begin group conversations.
            </p>
            <div class="bg-green-50 rounded-xl p-4 border border-green-200">
                <div class="flex items-center gap-2 mb-2">
                    <span class="material-icons text-green-600 text-sm"
                        >info</span
                    >
                    <span class="text-sm font-semibold text-green-700"
                        >Group Features</span
                    >
                </div>
                <ul class="text-xs text-green-600 space-y-1">
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">check_circle</span>
                        Multiple participants
                    </li>
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">check_circle</span>
                        Public and private groups
                    </li>
                    <li class="flex items-center gap-2">
                        <span class="material-icons text-xs">check_circle</span>
                        Secret groups with E2EE
                    </li>
                </ul>
            </div>
        </div>

        <!-- Join Group Modal -->
        <div
            v-if="showJoinGroupModal"
            class="fixed inset-0 bg-gray-300 bg-opacity-30 backdrop-blur-md flex items-center justify-center z-50 p-4"
        >
            <div
                class="bg-white rounded-3xl shadow-2xl border border-blue-100 p-8 w-96 max-w-[90vw] transform transition-all duration-300 scale-100 hover:shadow-3xl"
            >
                <!-- Close Button -->
                <button
                    class="absolute top-4 right-4 text-gray-400 hover:text-blue-700 hover:bg-blue-50 w-10 h-10 rounded-full transition-all duration-200 cursor-pointer flex items-center justify-center shadow-sm hover:shadow-md"
                    @click="showJoinGroupModal = false"
                    aria-label="Close"
                >
                    <span class="material-icons text-xl">close</span>
                </button>

                <!-- Header -->
                <div class="mb-8 text-center">
                    <div class="mb-6">
                        <div
                            class="w-16 h-16 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-2xl flex items-center justify-center mx-auto shadow-lg"
                        >
                            <span class="material-icons text-3xl text-white"
                                >group_add</span
                            >
                        </div>
                    </div>
                    <h3
                        class="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-500 mb-3 tracking-tight"
                    >
                        Join Group
                    </h3>
                    <p class="text-sm text-gray-600 leading-relaxed">
                        Enter the invite link to join an existing group.
                    </p>
                </div>

                <!-- Join Form -->
                <div class="space-y-6">
                    <div>
                        <label
                            class="block mb-3 text-blue-700 font-semibold text-sm"
                        >
                            <span
                                class="material-icons align-middle mr-1 text-base"
                                >link</span
                            >
                            Invite Link
                        </label>
                        <div class="relative">
                <input
                    v-model="joinGroupCode"
                    type="text"
                                placeholder="Enter invite link..."
                                class="w-full border-2 border-blue-200 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-300 focus:border-blue-400 transition-all duration-200 text-gray-700 shadow-sm hover:shadow-md"
                                @keyup.enter="handleJoinGroup"
                                :disabled="isJoiningGroup"
                            />
                            <div
                                v-if="isJoiningGroup"
                                class="absolute right-3 top-1/2 transform -translate-y-1/2"
                            >
                                <div
                                    class="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-500"
                                ></div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Action Buttons -->
                <div class="flex gap-3 mt-8">
                    <button
                        @click="showJoinGroupModal = false"
                        class="flex-1 px-6 py-3 text-gray-600 border-2 border-gray-300 rounded-xl hover:bg-gray-50 transition-all duration-200 font-semibold"
                    >
                        Cancel
                    </button>
                    <button
                        @click="handleJoinGroup"
                        :disabled="isJoiningGroup || !joinGroupCode.trim()"
                        class="flex-1 px-6 py-3 bg-gradient-to-r from-blue-500 to-indigo-500 text-white rounded-xl hover:from-blue-600 hover:to-indigo-600 transition-all duration-200 font-semibold shadow-lg hover:shadow-xl transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none border border-blue-400"
                    >
                        <span
                            v-if="!isJoiningGroup"
                            class="material-icons align-middle mr-2 text-base"
                            >group_add</span
                        >
                        <span
                            v-else
                            class="material-icons align-middle mr-2 text-base animate-spin"
                            >hourglass_empty</span
                        >
                        {{ isJoiningGroup ? "Joining..." : "Join Group" }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Create Group Modal -->
        <div
            v-if="showCreateGroupModal"
            class="fixed inset-0 bg-gray-300 bg-opacity-30 backdrop-blur-md flex items-center justify-center z-50 p-4"
        >
            <div
                class="bg-white rounded-3xl shadow-2xl border border-green-100 p-8 w-96 max-w-[90vw] max-h-[90vh] overflow-y-auto transform transition-all duration-300 scale-100 hover:shadow-3xl"
            >
                <!-- Close Button -->
                <button
                    class="absolute top-4 right-4 text-gray-400 hover:text-green-700 hover:bg-green-50 w-10 h-10 rounded-full transition-all duration-200 cursor-pointer flex items-center justify-center"
                    @click="showCreateGroupModal = false"
                    aria-label="Close"
                >
                    <span class="material-icons text-xl">close</span>
                </button>

                <!-- Header -->
                <div class="mb-8 text-center">
                    <div class="mb-6">
                        <div
                            class="w-16 h-16 bg-gradient-to-br from-green-500 to-emerald-600 rounded-2xl flex items-center justify-center mx-auto shadow-lg"
                        >
                            <span class="material-icons text-3xl text-white"
                                >add_circle</span
                            >
                        </div>
                    </div>
                    <h3
                        class="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-green-600 to-emerald-500 mb-3 tracking-tight"
                    >
                        Create Group
                    </h3>
                    <p class="text-sm text-gray-600 leading-relaxed">
                        Create a new group for collaborative conversations.
                    </p>
                </div>

                <!-- Create Group Form -->
                <div class="space-y-6">
                    <!-- Group Avatar Upload -->
                    <div>
                        <label
                            class="block mb-3 text-green-700 font-semibold text-sm"
                        >
                            <span
                                class="material-icons align-middle mr-1 text-base"
                                >photo_camera</span
                            >
                            Group Avatar (Optional)
                        </label>
                        <div class="flex items-center gap-4">
                            <div class="relative">
                                <img
                                    :src="
                                        newGroup.avatar_url ||
                                        '/src/assets/default-avatar.jpg'
                                    "
                                    alt="Group Avatar Preview"
                                    class="w-16 h-16 rounded-full object-cover border-2 border-green-300 shadow-sm"
                                />
                                <div
                                    class="absolute -bottom-1 -right-1 w-6 h-6 bg-green-500 rounded-full flex items-center justify-center"
                                >
                                    <span
                                        class="material-icons text-white text-xs"
                                        >group</span
                                    >
                                </div>
                            </div>
                            <div class="flex-1">
                                <input
                                    ref="avatarInput"
                                    type="file"
                                    accept=".jpg,.jpeg,.png,.webp"
                                    @change="handleAvatarUpload"
                                    class="hidden"
                                />
                                <button
                                    @click="$refs.avatarInput.click()"
                                    class="w-full px-4 py-2 border-2 border-green-200 rounded-xl hover:border-green-300 transition-all duration-200 text-green-700 font-medium"
                                >
                                    <span
                                        class="material-icons align-middle mr-1 text-sm"
                                        >upload</span
                                    >
                                    Choose Avatar
                                </button>
                            </div>
                        </div>
                    </div>

                    <!-- Group Name -->
                    <div>
                        <label
                            class="block mb-3 text-green-700 font-semibold text-sm"
                        >
                            <span
                                class="material-icons align-middle mr-1 text-base"
                                >edit</span
                            >
                            Group Name *
                        </label>
                        <input
                            v-model="newGroup.name"
                            type="text"
                            placeholder="Enter group name..."
                            class="w-full border-2 border-green-200 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-green-300 focus:border-green-400 transition-all duration-200 text-gray-700"
                        />
                    </div>

                    <!-- Group Description -->
                    <div>
                        <label
                            class="block mb-3 text-green-700 font-semibold text-sm"
                        >
                            <span
                                class="material-icons align-middle mr-1 text-base"
                                >description</span
                            >
                            Description (Optional)
                        </label>
                        <textarea
                            v-model="newGroup.description"
                            placeholder="Enter group description..."
                            rows="3"
                            class="w-full border-2 border-green-200 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-green-300 focus:border-green-400 transition-all duration-200 text-gray-700 resize-none"
                        ></textarea>
                    </div>

                    <!-- Group Type -->
                    <div>
                        <label
                            class="block mb-3 text-green-700 font-semibold text-sm"
                        >
                            <span
                                class="material-icons align-middle mr-1 text-base"
                                >security</span
                            >
                            Group Type
                        </label>
                        <select
                            v-model="newGroup.type"
                            class="w-full border-2 border-green-200 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-green-300 focus:border-green-400 transition-all duration-200 text-gray-700"
                        >
                            <option value="public">
                                Public (Anyone can join)
                            </option>
                            <option value="private">
                                Private (Invite only)
                            </option>
                        </select>
                        <p class="text-xs text-gray-500 mt-2">
                            <span v-if="newGroup.type === 'public'">
                                <span
                                    class="material-icons text-green-500 text-xs align-middle mr-1"
                                    >public</span
                                >
                                Public groups are visible to everyone and can be
                                joined without invitation.
                            </span>
                            <span v-else>
                                <span
                                    class="material-icons text-orange-500 text-xs align-middle mr-1"
                                    >lock_outline</span
                                >
                                Private groups require an invite code to join
                                and are not publicly visible.
                            </span>
                        </p>
                    </div>
                </div>

                <!-- Action Buttons -->
                <div class="flex gap-3 mt-8">
                    <button
                        @click="showCreateGroupModal = false"
                        class="flex-1 px-6 py-3 text-gray-600 border-2 border-gray-300 rounded-xl hover:bg-gray-50 transition-all duration-200 font-semibold"
                    >
                        Cancel
                    </button>
                    <button
                        @click="handleCreateGroup"
                        :disabled="isCreatingGroup || !newGroup.name.trim()"
                        class="flex-1 px-6 py-3 bg-gradient-to-r from-green-500 to-emerald-500 text-white rounded-xl hover:from-green-600 hover:to-emerald-600 transition-all duration-200 font-semibold shadow-lg hover:shadow-xl transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                    >
                        <span
                            v-if="!isCreatingGroup"
                            class="material-icons align-middle mr-2 text-base"
                            >add_circle</span
                        >
                        <span
                            v-else
                            class="material-icons align-middle mr-2 text-base animate-spin"
                            >hourglass_empty</span
                        >
                        {{ isCreatingGroup ? "Creating..." : "Create Group" }}
                    </button>
                </div>
            </div>
        </div>

        <!-- Create Secret Group Modal -->
        <div
            v-if="showCreateSecretGroupModal"
            class="fixed inset-0 bg-gray-300 bg-opacity-30 backdrop-blur-md flex items-center justify-center z-50 p-4"
        >
            <div
                class="bg-white rounded-3xl shadow-2xl border border-purple-100 p-8 w-96 max-w-[90vw] max-h-[90vh] overflow-y-auto transform transition-all duration-300 scale-100 hover:shadow-3xl"
            >
                <!-- Close Button -->
                <button
                    class="absolute top-4 right-4 text-gray-400 hover:text-purple-700 hover:bg-purple-50 w-10 h-10 rounded-full transition-all duration-200 cursor-pointer flex items-center justify-center"
                    @click="showCreateSecretGroupModal = false"
                    aria-label="Close"
                >
                    <span class="material-icons text-xl">close</span>
                </button>

                <!-- Header -->
                <div class="mb-8 text-center">
                    <div class="mb-6">
                        <div
                            class="w-16 h-16 bg-gradient-to-br from-purple-500 to-pink-600 rounded-2xl flex items-center justify-center mx-auto shadow-lg"
                        >
                            <span class="material-icons text-3xl text-white"
                                >lock</span
                            >
                        </div>
                    </div>
                    <h3
                        class="text-3xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-purple-600 to-pink-500 mb-3 tracking-tight"
                    >
                        Create Secret Group
                    </h3>
                    <p class="text-sm text-gray-600 leading-relaxed">
                        Create an end-to-end encrypted group chat for maximum
                        privacy.
                    </p>
                </div>

                <!-- Create Secret Group Form -->
                <div class="space-y-6">
                    <!-- Group Name -->
                    <div>
                        <label
                            class="block mb-3 text-purple-700 font-semibold text-sm"
                        >
                            <span
                                class="material-icons align-middle mr-1 text-base"
                                >edit</span
                            >
                            Group Name *
                        </label>
                        <input
                            v-model="newSecretGroup.name"
                            type="text"
                            placeholder="Enter secret group name..."
                            class="w-full border-2 border-purple-200 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-purple-300 focus:border-purple-400 transition-all duration-200 text-gray-700"
                        />
                    </div>

                    <!-- Group Description -->
                    <div>
                        <label
                            class="block mb-3 text-purple-700 font-semibold text-sm"
                        >
                            <span
                                class="material-icons align-middle mr-1 text-base"
                                >description</span
                            >
                            Description (Optional)
                        </label>
                        <textarea
                            v-model="newSecretGroup.description"
                            placeholder="Enter group description..."
                            rows="3"
                            class="w-full border-2 border-purple-200 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-purple-300 focus:border-purple-400 transition-all duration-200 text-gray-700 resize-none"
                        ></textarea>
                    </div>

                    <!-- Security Info -->
                    <div
                        class="bg-purple-50 rounded-xl p-4 border border-purple-200"
                    >
                        <div class="flex items-center gap-2 mb-2">
                            <span class="material-icons text-purple-600 text-sm"
                                >security</span
                            >
                            <span class="text-sm font-semibold text-purple-700"
                                >End-to-End Encrypted</span
                            >
                        </div>
                        <p class="text-xs text-purple-600 leading-relaxed">
                            All messages in this group will be encrypted with
                            unique keys for each member. Only group members can
                            decrypt the messages.
                        </p>
                    </div>
                </div>

                <!-- Action Buttons -->
                <div class="flex gap-3 mt-8">
                    <button
                        @click="showCreateSecretGroupModal = false"
                        class="flex-1 px-6 py-3 text-gray-600 border-2 border-gray-300 rounded-xl hover:bg-gray-50 transition-all duration-200 font-semibold"
                    >
                        Cancel
                    </button>
                    <button
                        @click="handleCreateSecretGroup"
                        :disabled="
                            isCreatingSecretGroup || !newSecretGroup.name.trim()
                        "
                        class="flex-1 px-6 py-3 bg-gradient-to-r from-purple-500 to-pink-500 text-white rounded-xl hover:from-purple-600 hover:to-pink-600 transition-all duration-200 font-semibold shadow-lg hover:shadow-xl transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                    >
                        <span
                            v-if="!isCreatingSecretGroup"
                            class="material-icons align-middle mr-2 text-base"
                            >lock</span
                        >
                        <span
                            v-else
                            class="material-icons align-middle mr-2 text-base animate-spin"
                            >hourglass_empty</span
                        >
                        {{
                            isCreatingSecretGroup
                                ? "Creating..."
                                : "Create Secret Group"
                        }}
                    </button>
                </div>
            </div>
        </div>
    </div>

    <!-- Leave Group Modal -->
    <LeaveGroupModal
        v-if="selectedGroup"
        :is-visible="showLeaveModal"
        :group="selectedGroup"
        :backend-base-url="backendBaseUrl"
        @close="handleLeaveModalClose"
        @action-completed="handleLeaveModalActionCompleted"
    />

    <!-- Delete Group Modal -->
    <LeaveGroupModal
        v-if="selectedGroup"
        :is-visible="showDeleteModal"
        :group="selectedGroup"
        :backend-base-url="backendBaseUrl"
        @close="handleDeleteModalClose"
        @action-completed="handleDeleteModalActionCompleted"
    />

    <!-- Approval Modal -->
    <ApprovalModal
        :is-visible="showApprovalModal"
        :invite-link="approvalInviteLink"
        @close="handleApprovalModalClose"
        @approval-submitted="handleApprovalSubmitted"
    />

    <!-- Approvals Modal -->
    <ApprovalsModal
        :is-visible="showApprovalsModal"
        @close="handleApprovalsModalClose"
        @approval-updated="handleApprovalUpdated"
    />

    <!-- Update Group Modal -->
    <UpdateGroupModal
        :is-visible="showUpdateGroupModal"
        :group="selectedGroup"
        :backend-base-url="backendBaseUrl"
        @close="handleUpdateGroupModalClose"
        @group-updated="handleGroupUpdated"
    />

    <!-- Group Members Modal -->
    <GroupMembersModal
        :is-visible="showGroupMembersModal"
        :group="selectedGroup"
        :backend-base-url="backendBaseUrl"
        :current-user-id="userStore.user_id"
        @close="handleGroupMembersModalClose"
        @member-updated="handleMemberUpdated"
        @start-chat="handleStartChat"
    />
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { showMessage, showError } from "../../utils/toast";
import { useGroupStore } from "../../stores/groups";
import { useUserStore } from "../../stores/users";
import LeaveGroupModal from "../ui/LeaveGroupModal.vue";
import ApprovalModal from "../ui/ApprovalModal.vue";
import ApprovalsModal from "../ui/ApprovalsModal.vue";
import UpdateGroupModal from "../ui/UpdateGroupModal.vue";
import GroupMembersModal from "../ui/GroupMembersModal.vue";

// Props
const props = defineProps({
    backendBaseUrl: {
        type: String,
        default: import.meta.env.VITE_BACKEND_BASE_URL,
    },
});

// Stores
const groupStore = useGroupStore();
const userStore = useUserStore();

// Emits
const emit = defineEmits(["group-clicked", "switch-to-chat"]);

// Reactive state
const isLoading = computed(() => groupStore.isLoading);
const groups = computed(() => groupStore.getUserGroups);
const showJoinGroupModal = ref(false);
const showCreateGroupModal = ref(false);
const showCreateSecretGroupModal = ref(false);
const showLeaveModal = ref(false);
const showDeleteModal = ref(false);
const showApprovalModal = ref(false);
const showApprovalsModal = ref(false);
const showUpdateGroupModal = ref(false);
const showGroupMembersModal = ref(false);
const selectedGroup = ref(null);
const approvalInviteLink = ref(null);
const joinGroupCode = ref("");
const isJoiningGroup = ref(false);
const isCreatingGroup = ref(false);
const isCreatingSecretGroup = ref(false);

// Form data
const newGroup = reactive({
    name: "",
    description: "",
    type: "public",
    avatar_url: "",
    avatar_file: null,
});

const newSecretGroup = reactive({
    name: "",
    description: "",
});

// Load user groups on component mount
onMounted(async () => {
    await loadUserGroups();
});

// API Functions
const loadUserGroups = async () => {
    try {
        await groupStore.loadUserGroups();
    } catch (error) {
        console.error("Failed to load groups:", error);
        showError("Failed to load groups. Please try again.");
    }
};

const handleGroupClick = (group) => {
    console.log("Opening group:", group);
    // Set the current group in the store
    groupStore.setCurrentGroup(group);
    // Emit event to parent to open group chat
    emit("group-clicked", group);
};

const handleLeaveGroup = (group) => {
    selectedGroup.value = group;
    showLeaveModal.value = true;
};

const handleLeaveModalClose = () => {
    showLeaveModal.value = false;
    selectedGroup.value = null;
};

const handleLeaveModalActionCompleted = () => {
    // Group was successfully left/deleted, no need to do anything else
    // The store already updated the groups list
};

const handleDeleteGroup = (group) => {
    selectedGroup.value = group;
    showDeleteModal.value = true;
};

const handleDeleteModalClose = () => {
    showDeleteModal.value = false;
    selectedGroup.value = null;
};

const handleDeleteModalActionCompleted = () => {
    // Group was successfully deleted, no need to do anything else
    // The store already updated the groups list
};

const handleApprovalModalClose = () => {
    showApprovalModal.value = false;
    approvalInviteLink.value = null;
};

const handleApprovalSubmitted = () => {
    // This function is no longer needed since we show the toast directly in the modal
};

const handleOpenApprovals = () => {
    showApprovalsModal.value = true;
};

const handleEditGroup = (group) => {
    selectedGroup.value = group;
    showUpdateGroupModal.value = true;
};

const handleUpdateGroupModalClose = () => {
    showUpdateGroupModal.value = false;
    selectedGroup.value = null;
};

const handleGroupUpdated = async (updatedGroup) => {
    console.log('Group updated:', updatedGroup);
    
    try {
        // Fetch updated groups from the backend
        await loadUserGroups();
        console.log('✅ Groups refreshed after update');
    } catch (error) {
        console.error('❌ Failed to refresh groups after update:', error);
        showError('Group updated but failed to refresh groups list');
    }
    
    showUpdateGroupModal.value = false;
    selectedGroup.value = null;
};

const handleManageMembers = (group) => {
    selectedGroup.value = group;
    showGroupMembersModal.value = true;
};

const handleGroupMembersModalClose = () => {
    showGroupMembersModal.value = false;
    selectedGroup.value = null;
};

const handleMemberUpdated = (data) => {
    console.log('Member updated:', data);
    // Optionally refresh groups if needed
    // await loadUserGroups();
};

const handleStartChat = (userData) => {
    console.log('💬 Starting chat with user from group:', userData);
    
    // Create a user object for the chat
    const user = {
        id: userData.user_id,
        username: userData.username,
        avatar_url: userData.avatar_url
    };
    
    // Emit the event to parent component to switch to chats tab and open the chat
    emit('switch-to-chat', user);
};

const handleApprovalsModalClose = () => {
    showApprovalsModal.value = false;
};

const handleApprovalUpdated = () => {
    // Refresh groups list after approval action
    loadUserGroups();
};

const handleCopyInviteLink = async (group) => {
    try {
        await navigator.clipboard.writeText(group.invite_link);
        showMessage("Invite link copied to clipboard!");
    } catch (error) {
        console.error("Failed to copy invite link:", error);
        showError("Failed to copy invite link. Please try again.");
    }
};

const handleJoinGroup = async () => {
    try {
        if (!joinGroupCode.value.trim()) {
            showError("Please enter an invite link");
            return;
        }
        
        isJoiningGroup.value = true;

        // Extract invite link from input
        const inviteLink = joinGroupCode.value.trim();

        // Check if this is a secret group invite (you might need to detect this based on your invite link format)
        // For now, we'll try regular group first, then secret group if it fails
        try {
            // Try joining as regular group first
            const response = await groupStore.joinGroup(inviteLink);
            showMessage("Successfully joined group!");
        } catch (regularGroupError) {
            // If regular group join fails, try as secret group
            if (regularGroupError.response?.status === 404 || regularGroupError.response?.status === 400) {
                try {
                    const response = await groupStore.joinSecretGroup(inviteLink);
                    showMessage("Successfully joined secret group!");
                } catch (secretGroupError) {
                    // Handle approval-specific errors for secret groups
                    const errorType = secretGroupError.response?.data?.type;
                    const errorDetail = secretGroupError.response?.data?.detail;

                    switch (errorType) {
                        case "userApprovalNotFound":
                            // Store the invite link for approval submission
                            approvalInviteLink.value = joinGroupCode.value.trim();
                            showApprovalModal.value = true;
                            break;
                        case "userApprovalStatus":
                            if (errorDetail?.includes("pending")) {
                                showError(
                                    "Your approval is pending. Please wait for admin approval."
                                );
                            } else if (errorDetail?.includes("rejected")) {
                                showError(
                                    "Your approval has been rejected. Please contact an administrator."
                                );
                            } else {
                                showError(errorDetail || "Approval status error occurred.");
                            }
                            break;
                        case "getUserApproval":
                            showError("Failed to check approval status. Please try again.");
                            break;
                        default:
                            showError(
                                errorDetail ||
                                    "Failed to join secret group. Please check the invite link and try again."
                            );
                    }
                    return;
                }
            } else {
                // Handle approval-specific errors for regular groups
                const errorType = regularGroupError.response?.data?.type;
                const errorDetail = regularGroupError.response?.data?.detail;

                switch (errorType) {
                    case "userApprovalNotFound":
                        // Store the invite link for approval submission
                        approvalInviteLink.value = joinGroupCode.value.trim();
                        showApprovalModal.value = true;
                        break;
                    case "userApprovalStatus":
                        if (errorDetail?.includes("pending")) {
                            showError(
                                "Your approval is pending. Please wait for admin approval."
                            );
                        } else if (errorDetail?.includes("rejected")) {
                            showError(
                                "Your approval has been rejected. Please contact an administrator."
                            );
                        } else {
                            showError(errorDetail || "Approval status error occurred.");
                        }
                        break;
                    case "getUserApproval":
                        showError("Failed to check approval status. Please try again.");
                        break;
                    default:
                        showError(
                            errorDetail ||
                                "Failed to join group. Please check the invite link and try again."
                        );
                }
                return;
            }
        }

        showJoinGroupModal.value = false;
        joinGroupCode.value = "";
    } catch (error) {
        console.error("Failed to join group:", error);
        showError("Failed to join group. Please check the invite link and try again.");
    } finally {
        isJoiningGroup.value = false;
    }
};

const handleCreateGroup = async () => {
    try {
        if (!newGroup.name.trim()) {
            showError("Please enter a group name");
            return;
        }
        
        isCreatingGroup.value = true;

        // Create FormData for multipart/form-data
        const formData = new FormData();
        formData.append("name", newGroup.name.trim());
        formData.append("description", newGroup.description.trim());
        formData.append("group_type", newGroup.type);

        // Add file if selected
        if (newGroup.avatar_file) {
            formData.append("file", newGroup.avatar_file);
        }

        console.log("Creating group with FormData:", {
            name: newGroup.name.trim(),
            description: newGroup.description.trim(),
            group_type: newGroup.type,
            hasFile: !!newGroup.avatar_file,
        });

        const createdGroup = await groupStore.createGroup(formData);

        // Show success message from backend or default
        const successMessage =
            createdGroup.message || "Group created successfully!";
        showMessage(successMessage);

        // If there's an invite link, show it
        if (createdGroup.invite_link) {
            console.log("Group invite link:", createdGroup.invite_link);
            // You could show this in a modal or copy to clipboard
        }

        showCreateGroupModal.value = false;
        
        // Reset form
        newGroup.name = "";
        newGroup.description = "";
        newGroup.type = "public";
        newGroup.avatar_file = null;
        newGroup.avatar_url = "";
    } catch (error) {
        console.error("Failed to create group:", error);
        let errorMessage = "Failed to create group. Please try again.";

        if (error.response?.data?.detail) {
            errorMessage = error.response.data.detail;
        } else if (error.response?.data?.message) {
            errorMessage = error.response.data.message;
        } else if (error.message) {
            errorMessage = error.message;
        }

        showError(errorMessage);
    } finally {
        isCreatingGroup.value = false;
    }
};

const handleCreateSecretGroup = async () => {
    try {
        if (!newSecretGroup.name.trim()) {
            showError("Please enter a group name");
            return;
        }
        
        isCreatingSecretGroup.value = true;

        // Create FormData for multipart/form-data
        const formData = new FormData();
        formData.append("name", newSecretGroup.name.trim());
        formData.append("description", newSecretGroup.description.trim());
        formData.append("group_type", "secret");

        console.log("Creating secret group with FormData:", {
            name: newSecretGroup.name.trim(),
            description: newSecretGroup.description.trim(),
            group_type: "secret",
        });

        const createdGroup = await groupStore.createSecretGroup(formData);

        // Show success message from backend or default
        const successMessage =
            createdGroup.message || "Secret group created successfully!";
        showMessage(successMessage);

        // If there's an invite link, show it
        if (createdGroup.invite_link) {
            console.log("Secret group invite link:", createdGroup.invite_link);
            // You could show this in a modal or copy to clipboard
        }

        showCreateSecretGroupModal.value = false;
        
        // Reset form
        newSecretGroup.name = "";
        newSecretGroup.description = "";
    } catch (error) {
        console.error("Failed to create secret group:", error);
        let errorMessage = "Failed to create secret group. Please try again.";

        if (error.response?.data?.detail) {
            errorMessage = error.response.data.detail;
        } else if (error.response?.data?.message) {
            errorMessage = error.response.data.message;
        } else if (error.message) {
            errorMessage = error.message;
        }

        showError(errorMessage);
    } finally {
        isCreatingSecretGroup.value = false;
    }
};

// Avatar upload functions
const handleAvatarUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
        // Check file size (5MB limit)
        if (file.size > 5 * 1024 * 1024) {
            showError("Avatar file size must be less than 5MB");
            return;
        }

        // Check file format
        const allowedFormats = [".jpg", ".jpeg", ".png", ".webp"];
        const fileExtension = "." + file.name.split(".").pop().toLowerCase();

        if (!allowedFormats.includes(fileExtension)) {
            showError("Only JPG, JPEG, PNG, and WebP formats are allowed");
            return;
        }

        // Store the file for FormData
        newGroup.avatar_file = file;

        // Create preview
        const reader = new FileReader();
        reader.onload = (e) => {
            newGroup.avatar_url = e.target.result;
        };
        reader.readAsDataURL(file);
    }
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/icon?family=Material+Icons");
</style> 
