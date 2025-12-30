<script>
    import { onMount } from "svelte"
    import {
        getMe,
        renameMe,
        changeMyPassword,
        deleteMe,
        getMyTasks,
        createNewTask,
        deleteTask,
        switchTaskStatus,
        changeTaskTitle,
        changeTaskDescription
    } from "../api/http.js"

    export let onDeleted

    let me = null
    let error = ""
    let message = ""

    // ---------- TASKS ----------
    let tasks = []

    // ---------- DERIVED ----------
    $: activeTasks = tasks.filter(t => !t.done)
    $: completedTasks = tasks.filter(t => t.done)

    // ---------- FORMS ----------
    let newName = ""
    let oldPassword = ""
    let newPassword = ""

    let newTaskTitle = ""
    let newTaskDescription = ""

    // ---------- EDIT STATE ----------
    let editingTaskId = null
    let editingTitle = ""
    let editingDescription = ""

    // ---------- HELPERS ----------
    function normalizeTask(task) {
        return {
            ...task,
            done: task.is_completed
        }
    }

    function startEditTask(task) {
        editingTaskId = task.id
        editingTitle = task.title
        editingDescription = task.description
    }

    function cancelEdit() {
        editingTaskId = null
        editingTitle = ""
        editingDescription = ""
    }

    // ---------- LIFECYCLE ----------
    onMount(async () => {
        try {
            me = await getMe()
            newName = me.name

            const rawTasks = await getMyTasks()
            tasks = rawTasks.map(normalizeTask)
        } catch (err) {
            error = err.message
        }
    })

    // ---------- ACTIONS ----------
    async function submitRename() {
        error = ""
        message = ""
        try {
            const updated = await renameMe(newName)
            me.name = updated.name
            message = "Name updated"
        } catch (err) {
            error = err.message
        }
    }

    async function submitPasswordChange() {
        error = ""
        message = ""
        try {
            await changeMyPassword(oldPassword, newPassword)
            oldPassword = ""
            newPassword = ""
            message = "Password changed"
        } catch (err) {
            error = err.message
        }
    }

    async function submitCreateTask() {
        error = ""
        message = ""
        try {
            const created = await createNewTask(
                newTaskTitle,
                newTaskDescription
            )

            tasks = [...tasks, normalizeTask(created)]
            newTaskTitle = ""
            newTaskDescription = ""
            message = "Task created"
        } catch (err) {
            error = err.message
        }
    }

    async function submitDeleteTask(taskId) {
        if (!confirm("Delete this task?")) return

        error = ""
        message = ""

        try {
            await deleteTask(taskId)
            tasks = tasks.filter(t => t.id !== taskId)
            message = "Task deleted"
        } catch (err) {
            error = err.message
        }
    }

    async function submitSwitchTask(taskId) {
        error = ""

        try {
            const updated = await switchTaskStatus(taskId)
            const normalized = normalizeTask(updated)

            tasks = tasks.map(t =>
                t.id === taskId ? normalized : t
            )
        } catch (err) {
            error = err.message
        }
    }

    async function submitEditTask(taskId) {
        error = ""
        message = ""

        try {
            await changeTaskTitle(taskId, editingTitle)
            const updated = await changeTaskDescription(taskId, editingDescription)

            const normalized = normalizeTask(updated)

            tasks = tasks.map(t =>
                t.id === taskId ? normalized : t
            )

            cancelEdit()
            message = "Task updated"
        } catch (err) {
            error = err.message
        }
    }

    async function submitDeleteAccount() {
        if (!confirm("Are you sure you want to delete your account?")) return

        error = ""
        message = ""

        try {
            await deleteMe()
            onDeleted()
        } catch (err) {
            error = err.message
        }
    }
</script>

{#if error}
    <p style="color:red">{error}</p>
{/if}

{#if !me}
    <p>Loading profile...</p>
{:else}

    <h1>Hello, {me.name} 👋</h1>
    <p><b>Role:</b> {me.role}</p>
    <p><b>ID:</b> {me.id}</p>

    <hr />

    <h2>Rename account</h2>
    <input bind:value={newName} />
    <button on:click={submitRename}>Rename</button>

    <hr />

    <h2>Change password</h2>
    <input type="password" placeholder="old" bind:value={oldPassword} />
    <input type="password" placeholder="new" bind:value={newPassword} />
    <button on:click={submitPasswordChange}>Change</button>

    <hr />

    <h2>Active tasks</h2>

    {#each activeTasks as task (task.id)}
        <div style="margin-bottom:12px">
            <button on:click={() => submitSwitchTask(task.id)}>⬜</button>

            {#if editingTaskId === task.id}
                <input bind:value={editingTitle} />
                <textarea bind:value={editingDescription} />

                <button on:click={() => submitEditTask(task.id)}>💾</button>
                <button on:click={cancelEdit}>✖</button>
            {:else}
                <b>{task.title}</b>
                <p style="margin:4px 0">{task.description}</p>
                <button on:click={() => startEditTask(task)}>✏️</button>
            {/if}

            <button on:click={() => submitDeleteTask(task.id)}>❌</button>
        </div>
    {/each}

    <hr />

    <h2>Completed tasks</h2>

    {#each completedTasks as task (task.id)}
        <div style="margin-bottom:12px">
            <button on:click={() => submitSwitchTask(task.id)}>✅</button>

            {#if editingTaskId === task.id}
                <input bind:value={editingTitle} />
                <textarea bind:value={editingDescription} />

                <button on:click={() => submitEditTask(task.id)}>💾</button>
                <button on:click={cancelEdit}>✖</button>
            {:else}
                <s><b>{task.title}</b></s>
                <p style="margin:4px 0">{task.description}</p>
                <button on:click={() => startEditTask(task)}>✏️</button>
            {/if}

            <button on:click={() => submitDeleteTask(task.id)}>❌</button>
        </div>
    {/each}

    <hr />

    <h3>Create new task</h3>
    <input placeholder="title" bind:value={newTaskTitle} />
    <textarea placeholder="description" bind:value={newTaskDescription} />
    <button on:click={submitCreateTask}>Create</button>

    <hr />

    <h2 style="color:red">Danger zone</h2>
    <button style="background:red;color:white" on:click={submitDeleteAccount}>
        Delete my account
    </button>

    {#if message}
        <p style="color:green">{message}</p>
    {/if}

{/if}
