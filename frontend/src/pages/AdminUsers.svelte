<script>
    import { onMount } from "svelte"
    import {
        getAllUsersAdmin,
        createNewUserAdmin,
        getUserByIDAdmin,
        getUserTasksAdmin,
        createNewTaskAdmin,
        renameUserAdmin,
        changeUserPasswordAdmin,
        updateUserRoleAdmin,
        deleteUserAdmin,
        deleteTaskAdmin,
        switchTaskStatusAdmin,
        changeTaskTitleAdmin,
        changeTaskDescriptionAdmin
    } from "../api/http.js"

    let users = []
    let loadingUsers = false

    let selectedUser = null
    let selectedUserTasks = []

    let actionMessage = ""
    let actionError = ""

    let newUserName = ""
    let newUserPassword = ""

    let newTaskTitle = ""
    let newTaskDescription = ""

    let renameUserName = ""
    let adminNewPassword = ""

    let editingTaskId = null
    let editingTaskTitle = ""
    let editingTaskDescription = ""

    onMount(loadUsers)

    async function loadUsers() {
        loadingUsers = true
        users = await getAllUsersAdmin()
        loadingUsers = false
    }

    async function submitCreateUser() {
        await createNewUserAdmin(newUserName, newUserPassword)
        newUserName = ""
        newUserPassword = ""
        await loadUsers()
        actionMessage = "User created"
    }

    async function selectUser(userId) {
        selectedUser = await getUserByIDAdmin(userId)
        renameUserName = selectedUser.name
        selectedUserTasks = await getUserTasksAdmin(userId)
    }

    function clearSelectedUser() {
        selectedUser = null
        selectedUserTasks = []
        editingTaskId = null
    }

    /* ===== USER ACTIONS ===== */

    async function submitRenameUser() {
        const updated = await renameUserAdmin(selectedUser.id, renameUserName)
        selectedUser = updated
        users = users.map(u => u.id === updated.id ? updated : u)
        actionMessage = "User renamed"
    }

    async function submitChangeUserPassword() {
        await changeUserPasswordAdmin(selectedUser.id, "", adminNewPassword)
        adminNewPassword = ""
        actionMessage = "Password changed"
    }

    async function submitUpdateUserRole() {
        const updated = await updateUserRoleAdmin(selectedUser.id)
        selectedUser = updated
        users = users.map(u => u.id === updated.id ? updated : u)
        actionMessage = "Role switched"
    }

    async function submitDeleteUser() {
        const ok = confirm(`Delete user "${selectedUser.name}"?`)
        if (!ok) return

        await deleteUserAdmin(selectedUser.id)
        users = users.filter(u => u.id !== selectedUser.id)
        clearSelectedUser()
        actionMessage = "User deleted"
    }

    /* ===== TASKS ===== */

    async function submitCreateTaskForUser() {
        const created = await createNewTaskAdmin(
            selectedUser.id,
            newTaskTitle,
            newTaskDescription
        )
        selectedUserTasks = [...selectedUserTasks, created]
        newTaskTitle = ""
        newTaskDescription = ""
        actionMessage = "Task created"
    }

    async function submitSwitchTask(taskId) {
        const updated = await switchTaskStatusAdmin(taskId)
        selectedUserTasks = selectedUserTasks.map(t =>
            t.id === updated.id ? updated : t
        )
    }

    async function submitDeleteTask(taskId) {
        await deleteTaskAdmin(taskId)
        selectedUserTasks = selectedUserTasks.filter(t => t.id !== taskId)
    }

    function startEditTask(task) {
        editingTaskId = task.id
        editingTaskTitle = task.title
        editingTaskDescription = task.description
    }

    async function submitEditTask(taskId) {
        await changeTaskTitleAdmin(taskId, editingTaskTitle)
        const updated = await changeTaskDescriptionAdmin(
            taskId,
            editingTaskDescription
        )

        selectedUserTasks = selectedUserTasks.map(t =>
            t.id === updated.id ? updated : t
        )

        editingTaskId = null
        actionMessage = "Task updated"
    }
</script>

<main>
    <h1>Admin panel</h1>

    {#if actionMessage}
        <p class="success">{actionMessage}</p>
    {/if}

    <!-- CREATE USER -->
    <div class="card">
        <h2>Create new user</h2>
        <div class="form-row">
            <input placeholder="Username" bind:value={newUserName} />
            <input type="password" placeholder="Password" bind:value={newUserPassword} />
            <button on:click={submitCreateUser}>Create user</button>
        </div>
    </div>

    <!-- USERS -->
    <div class="card">
        <h2>Users</h2>
        <table class="users-table">
            <thead>
            <tr>
                <th>ID</th>
                <th>Username</th>
                <th>Role</th>
                <th>Action</th>
            </tr>
            </thead>
            <tbody>
            {#each users as user}
                <tr>
                    <td>{user.id}</td>
                    <td>{user.name}</td>
                    <td><span class="role {user.role}">{user.role}</span></td>
                    <td>
                        <button on:click={() => selectUser(user.id)}>View</button>
                    </td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>

    {#if selectedUser}
        <!-- USER SETTINGS -->
        <div class="card">
            <h2>User: {selectedUser.name}</h2>

            <input bind:value={renameUserName} />
            <button on:click={submitRenameUser}>Rename</button>

            <hr />

            <input type="password" placeholder="New password" bind:value={adminNewPassword} />
            <button on:click={submitChangeUserPassword}>Change password</button>

            <hr />

            <button on:click={submitUpdateUserRole}>Switch role</button>

            <hr />

            <button class="danger" on:click={submitDeleteUser}>Delete user</button>
            <button on:click={clearSelectedUser}>Close</button>
        </div>

        <!-- CREATE TASK -->
        <div class="card">
            <h2>Create task</h2>
            <input placeholder="Task title" bind:value={newTaskTitle} />
            <textarea placeholder="Task description" bind:value={newTaskDescription}></textarea>
            <button on:click={submitCreateTaskForUser}>Create task</button>
        </div>

        <!-- TASKS -->
        <div class="card">
            <h2>User tasks</h2>

            <table class="users-table">
                <thead>
                <tr>
                    <th>Done</th>
                    <th>Title</th>
                    <th>Description</th>
                    <th>Actions</th>
                </tr>
                </thead>
                <tbody>
                {#each selectedUserTasks as task}
                    <tr>
                        <td>
                            <button on:click={() => submitSwitchTask(task.id)}>
                                {task.is_completed ? "✅" : "⬜"}
                            </button>
                        </td>

                        <td>
                            {#if editingTaskId === task.id}
                                <input bind:value={editingTaskTitle} />
                            {:else}
                                {task.title}
                            {/if}
                        </td>

                        <td>
                            {#if editingTaskId === task.id}
                                <textarea bind:value={editingTaskDescription}></textarea>
                            {:else}
                                {task.description || "NO DESCRIPTION"}
                            {/if}
                        </td>

                        <td>
                            {#if editingTaskId === task.id}
                                <button on:click={() => submitEditTask(task.id)}>Save</button>
                                <button on:click={() => editingTaskId = null}>Cancel</button>
                            {:else}
                                <button on:click={() => startEditTask(task)}>Edit</button>
                                <button class="danger" on:click={() => submitDeleteTask(task.id)}>Delete</button>
                            {/if}
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    {/if}
</main>

<!--<style>-->
<!--    /* стили можно оставить как у тебя — они уже нормальные */-->
<!--</style>-->


<style>
    /* ===== GLOBAL ===== */
    body {
        margin: 0;
        font-family: system-ui, sans-serif;
        background: linear-gradient(180deg, #0f2a3d, #2437b7);
    }

    /* ===== LAYOUT ===== */
    main {
        max-width: 1100px;
        margin: 0 auto;
        padding: 40px 20px 80px;
    }

    h1 {
        text-align: center;
        color: white;
        margin-bottom: 40px;
    }

    /* ===== CARD ===== */
    .card {
        background: #f7e6e2;
        border-radius: 14px;
        padding: 24px;
        margin-bottom: 30px;
    }

    /* ===== FORM ===== */
    .form-row {
        display: grid;
        grid-template-columns: 1fr 1fr auto;
        gap: 12px;
    }

    input,
    textarea {
        border-radius: 10px; /* ← скругление */
        background: #ffffff;
        color: #111827;
        border: 1px solid #d1d5db;
    }

    input::placeholder,
    textarea::placeholder {
        color: #9ca3af;
    }

    input:focus,
    textarea:focus {
        outline: none;
        border-color: #2563eb;
        box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.25);
    }
    /* ===== BUTTON ===== */
    button {
        padding: 10px 18px;
        border-radius: 8px;
        border: none;
        background: #2563eb;
        color: white;
        cursor: pointer;
    }

    button.danger {
        background: #dc2626;
    }

    /* ===== TABLE ===== */
    .table-wrapper {
        margin-top: 20px;
        overflow-x: auto;
    }

    .users-table {
        width: 100%;
        border-collapse: collapse;
        background: #f5e4e0;
        border-radius: 12px;
        overflow: hidden;
    }

    .users-table th {
        color: #181818;
    }
    .users-table td {
        padding: 14px 18px;
        color: #111827; /* тёмно-серый, почти чёрный */
        font-size: 14px;
    }

    .users-table thead {
        background: #0038ef;
        color: #ffffff;
    }

    .users-table tbody tr:hover {
        background: #f1f5f9;
    }

    /* ===== ROLE ===== */
    .role {
        padding: 4px 12px;
        border-radius: 999px;
        font-size: 12px;
        font-weight: 600;
    }

    .role.admin {
        background: #fee2e2;
        color: #991b1b;
    }

    .role.user {
        background: #dbeafe;
        color: #1e40af;
    }

    /* ===== MESSAGES ===== */
    .success {
        color: #22c55e;
        text-align: center;
        margin-bottom: 20px;
    }
</style>
