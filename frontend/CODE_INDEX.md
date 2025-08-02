# ApprovalsModal.vue Code Index

## Component Overview
The `ApprovalsModal.vue` component is a Vue 3 modal dialog for managing group approval requests. It displays pending approval requests and allows users to approve or reject them.

## File Structure

### Template Section (Lines 1-390)

#### Modal Container (Lines 1-15)
- **Backdrop**: Semi-transparent overlay with blur effect
- **Modal Container**: Centered modal with max-width and responsive design
- **Z-index**: 50 for proper layering

#### Header Section (Lines 16-40)
- **Icon**: Material Icons "pending_actions" in blue circle
- **Title**: "Pending Approvals" with subtitle
- **Close Button**: X icon with hover effects

#### Content Section (Lines 41-390)

##### Loading State (Lines 42-52)
- **Spinner**: Animated loading indicator
- **Text**: "Loading approvals..." message

##### Empty State (Lines 54-66)
- **Icon**: Material Icons "check_circle" in gray
- **Message**: "No Pending Approvals" with description

##### Approvals List (Lines 68-390)
Each approval item contains:

###### Approval Header (Lines 70-105)
- **User Avatar**: Gray circle with person icon
- **User Info**: User ID and request date
- **Status Badge**: Color-coded status indicator
  - Yellow: pending
  - Green: approved  
  - Red: rejected

###### Group Information (Lines 107-125)
- **Group Icon**: Lock for secret groups, group for regular groups
- **Group ID**: Displayed with secret group indicator
- **Secret Badge**: Purple badge for secret groups

###### Reason Section (Lines 127-140)
- **Label**: "Reason for Request:"
- **Content**: White background box with request reason

###### Action Buttons (Lines 142-220)
- **Reject Button**: Red styling with loading state
- **Approve Button**: Green styling with loading state
- **Loading Indicators**: SVG spinners during processing

### Script Section (Lines 221-390)

#### Imports (Lines 222-224)
```javascript
import { ref, watch, onMounted } from "vue";
import { showInfo, showError } from "../../utils/toast";
import axiosInstance from "../../axiosInstance";
```

#### Props (Lines 226-231)
- `isVisible`: Boolean controlling modal visibility

#### Emits (Lines 233-235)
- `close`: Emitted when modal is closed
- `approval-updated`: Emitted when approval status changes

#### Reactive Data (Lines 237-240)
- `approvals`: Array of approval objects
- `isLoading`: Boolean for loading state
- `isProcessingApproval`: ID of approval being processed

#### Methods

##### closeModal (Lines 242-244)
- Emits close event to parent component

##### watch (Lines 246-252)
- Watches for `isVisible` prop changes
- Loads approvals when modal becomes visible

##### loadApprovals (Lines 254-268)
- **Endpoint**: `/api/received-approvals/get/`
- **Error Handling**: Shows error toast on failure
- **Loading State**: Manages loading indicator

##### handleApproveApproval (Lines 270-300)
- **Parameters**: `approvalId`
- **Endpoint**: `/api/approvals/edit-status/${approvalId}`
- **Secret Groups**: Adds `?is_secret=true` query parameter
- **Success**: Shows success message and removes from list
- **Error Handling**: Shows error toast on failure

##### handleRejectApproval (Lines 302-332)
- **Parameters**: `approvalId`
- **Endpoint**: `/api/approvals/edit-status/${approvalId}`
- **Secret Groups**: Adds `?is_secret=true` query parameter
- **Success**: Shows success message and removes from list
- **Error Handling**: Shows error toast on failure

##### formatDate (Lines 334-337)
- **Input**: Date string
- **Output**: Formatted date and time string

##### onMounted (Lines 339-343)
- Loads approvals if modal is visible on mount

## Key Features

### Visual Design
- **Modern UI**: Rounded corners, shadows, and clean typography
- **Color Coding**: Status-based color schemes
- **Responsive**: Mobile-friendly design with max-width constraints
- **Loading States**: Spinners and disabled states during processing

### User Experience
- **Real-time Updates**: Removes processed approvals from list
- **Error Handling**: Toast notifications for success/error states
- **Accessibility**: Proper button states and loading indicators
- **Modal Behavior**: Backdrop click to close, proper z-indexing

### Data Management
- **API Integration**: RESTful endpoints for CRUD operations
- **State Management**: Reactive Vue 3 composition API
- **Secret Group Support**: Special handling for secret group approvals
- **Optimistic Updates**: Immediate UI updates with API confirmation

## API Endpoints Used
- `GET /api/received-approvals/get/` - Fetch pending approvals
- `PUT /api/approvals/edit-status/{id}` - Update approval status
- `PUT /api/approvals/edit-status/{id}?is_secret=true` - Update secret group approval

## Dependencies
- **Vue 3**: Composition API with ref, watch, onMounted
- **Material Icons**: Icon library for UI elements
- **Axios**: HTTP client for API requests
- **Toast Utils**: Notification system for user feedback
- **Tailwind CSS**: Utility classes for styling

## Component Lifecycle
1. **Mount**: Loads approvals if modal is visible
2. **Visibility Change**: Watches for prop changes to load data
3. **User Actions**: Handles approve/reject with loading states
4. **Cleanup**: Emits events to parent for state management 