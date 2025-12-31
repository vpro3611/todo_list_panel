<!--<script>-->
<!--    import { onMount } from "svelte"-->
<!--    import {-->
<!--        getMe,-->
<!--        renameMe,-->
<!--        changeMyPassword,-->
<!--        deleteMe,-->
<!--        getMyTasks,-->
<!--        createNewTask,-->
<!--        deleteTask,-->
<!--        switchTaskStatus,-->
<!--        changeTaskTitle,-->
<!--        changeTaskDescription-->
<!--    } from "../api/http.js"-->

<!--    export let onDeleted-->
<!--    export let onLoaded-->


<!--    let me = null-->
<!--    let error = ""-->
<!--    let message = ""-->

<!--    // &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45; TASKS &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;-->
<!--    let tasks = []-->

<!--    // &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45; DERIVED &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;-->
<!--    $: activeTasks = tasks.filter(t => !t.done)-->
<!--    $: completedTasks = tasks.filter(t => t.done)-->

<!--    // &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45; FORMS &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;-->
<!--    let newName = ""-->
<!--    let oldPassword = ""-->
<!--    let newPassword = ""-->

<!--    let newTaskTitle = ""-->
<!--    let newTaskDescription = ""-->

<!--    // &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45; EDIT STATE &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;-->
<!--    let editingTaskId = null-->
<!--    let editingTitle = ""-->
<!--    let editingDescription = ""-->

<!--    // &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45; HELPERS &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;-->
<!--    function normalizeTask(task) {-->
<!--        return {-->
<!--            ...task,-->
<!--            done: task.is_completed-->
<!--        }-->
<!--    }-->

<!--    function startEditTask(task) {-->
<!--        editingTaskId = task.id-->
<!--        editingTitle = task.title-->
<!--        editingDescription = task.description-->
<!--    }-->

<!--    function cancelEdit() {-->
<!--        editingTaskId = null-->
<!--        editingTitle = ""-->
<!--        editingDescription = ""-->
<!--    }-->

<!--    // &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45; LIFECYCLE &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;-->
<!--    onMount(async () => {-->
<!--        try {-->
<!--            me = await getMe()-->
<!--            newName = me.name-->

<!--            const rawTasks = await getMyTasks()-->
<!--            tasks = rawTasks.map(normalizeTask)-->

<!--            onLoaded(me.role)-->

<!--        } catch (err) {-->
<!--            error = err.message-->
<!--        }-->
<!--    })-->

<!--    // &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45; ACTIONS &#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;&#45;-->
<!--    async function submitRename() {-->
<!--        error = ""-->
<!--        message = ""-->
<!--        try {-->
<!--            const updated = await renameMe(newName)-->
<!--            me.name = updated.name-->
<!--            message = "Name updated"-->
<!--        } catch (err) {-->
<!--            error = err.message-->
<!--        }-->
<!--    }-->

<!--    async function submitPasswordChange() {-->
<!--        error = ""-->
<!--        message = ""-->
<!--        try {-->
<!--            await changeMyPassword(oldPassword, newPassword)-->
<!--            oldPassword = ""-->
<!--            newPassword = ""-->
<!--            message = "Password changed"-->
<!--        } catch (err) {-->
<!--            error = err.message-->
<!--        }-->
<!--    }-->

<!--    async function submitCreateTask() {-->
<!--        error = ""-->
<!--        message = ""-->
<!--        try {-->
<!--            const created = await createNewTask(-->
<!--                newTaskTitle,-->
<!--                newTaskDescription-->
<!--            )-->

<!--            tasks = [...tasks, normalizeTask(created)]-->
<!--            newTaskTitle = ""-->
<!--            newTaskDescription = ""-->
<!--            message = "Task created"-->
<!--        } catch (err) {-->
<!--            error = err.message-->
<!--        }-->
<!--    }-->

<!--    async function submitDeleteTask(taskId) {-->
<!--        if (!confirm("Delete this task?")) return-->

<!--        error = ""-->
<!--        message = ""-->

<!--        try {-->
<!--            await deleteTask(taskId)-->
<!--            tasks = tasks.filter(t => t.id !== taskId)-->
<!--            message = "Task deleted"-->
<!--        } catch (err) {-->
<!--            error = err.message-->
<!--        }-->
<!--    }-->

<!--    async function submitSwitchTask(taskId) {-->
<!--        error = ""-->

<!--        try {-->
<!--            const updated = await switchTaskStatus(taskId)-->
<!--            const normalized = normalizeTask(updated)-->

<!--            tasks = tasks.map(t =>-->
<!--                t.id === taskId ? normalized : t-->
<!--            )-->
<!--        } catch (err) {-->
<!--            error = err.message-->
<!--        }-->
<!--    }-->

<!--    async function submitEditTask(taskId) {-->
<!--        error = ""-->
<!--        message = ""-->

<!--        try {-->
<!--            await changeTaskTitle(taskId, editingTitle)-->
<!--            const updated = await changeTaskDescription(taskId, editingDescription)-->

<!--            const normalized = normalizeTask(updated)-->

<!--            tasks = tasks.map(t =>-->
<!--                t.id === taskId ? normalized : t-->
<!--            )-->

<!--            cancelEdit()-->
<!--            message = "Task updated"-->
<!--        } catch (err) {-->
<!--            error = err.message-->
<!--        }-->
<!--    }-->

<!--    async function submitDeleteAccount() {-->
<!--        if (!confirm("Are you sure you want to delete your account?")) return-->

<!--        error = ""-->
<!--        message = ""-->

<!--        try {-->
<!--            await deleteMe()-->
<!--            onDeleted()-->
<!--        } catch (err) {-->
<!--            error = err.message-->
<!--        }-->
<!--    }-->
<!--</script>-->

<!--{#if error}-->
<!--    <p style="color:red">{error}</p>-->
<!--{/if}-->

<!--{#if !me}-->
<!--    <p>Loading profile...</p>-->
<!--{:else}-->

<!--    <h1>Hello, {me.name} 👋</h1>-->
<!--    <p><b>Role:</b> {me.role}</p>-->
<!--    <p><b>ID:</b> {me.id}</p>-->

<!--    <hr />-->

<!--    <h2>Rename account</h2>-->
<!--    <input bind:value={newName} />-->
<!--    <button on:click={submitRename}>Rename</button>-->

<!--    <hr />-->

<!--    <h2>Change password</h2>-->
<!--    <input type="password" placeholder="old" bind:value={oldPassword} />-->
<!--    <input type="password" placeholder="new" bind:value={newPassword} />-->
<!--    <button on:click={submitPasswordChange}>Change</button>-->

<!--    <hr />-->

<!--    <h2>Active tasks</h2>-->

<!--    {#each activeTasks as task (task.id)}-->
<!--        <div style="margin-bottom:12px">-->
<!--            <button on:click={() => submitSwitchTask(task.id)}>⬜</button>-->

<!--            {#if editingTaskId === task.id}-->
<!--                <input bind:value={editingTitle} />-->
<!--                <textarea bind:value={editingDescription} />-->

<!--                <button on:click={() => submitEditTask(task.id)}>💾</button>-->
<!--                <button on:click={cancelEdit}>✖</button>-->
<!--            {:else}-->
<!--                <b>{task.title}</b>-->
<!--                <p style="margin:4px 0">{task.description}</p>-->
<!--                <button on:click={() => startEditTask(task)}>✏️</button>-->
<!--            {/if}-->

<!--            <button on:click={() => submitDeleteTask(task.id)}>❌</button>-->
<!--        </div>-->
<!--    {/each}-->

<!--    <hr />-->

<!--    <h2>Completed tasks</h2>-->

<!--    {#each completedTasks as task (task.id)}-->
<!--        <div style="margin-bottom:12px">-->
<!--            <button on:click={() => submitSwitchTask(task.id)}>✅</button>-->

<!--            {#if editingTaskId === task.id}-->
<!--                <input bind:value={editingTitle} />-->
<!--                <textarea bind:value={editingDescription} />-->

<!--                <button on:click={() => submitEditTask(task.id)}>💾</button>-->
<!--                <button on:click={cancelEdit}>✖</button>-->
<!--            {:else}-->
<!--                <s><b>{task.title}</b></s>-->
<!--                <p style="margin:4px 0">{task.description}</p>-->
<!--                <button on:click={() => startEditTask(task)}>✏️</button>-->
<!--            {/if}-->

<!--            <button on:click={() => submitDeleteTask(task.id)}>❌</button>-->
<!--        </div>-->
<!--    {/each}-->

<!--    <hr />-->

<!--    <h3>Create new task</h3>-->
<!--    <input placeholder="title" bind:value={newTaskTitle} />-->
<!--    <textarea placeholder="description" bind:value={newTaskDescription} />-->
<!--    <button on:click={submitCreateTask}>Create</button>-->

<!--    <hr />-->

<!--    <h2 style="color:red">Danger zone</h2>-->
<!--    <button style="background:red;color:white" on:click={submitDeleteAccount}>-->
<!--        Delete my account-->
<!--    </button>-->

<!--    {#if message}-->
<!--        <p style="color:green">{message}</p>-->
<!--    {/if}-->

<!--{/if}-->


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
    export let onLoaded

    let me = null
    let error = ""
    let message = ""

    let tasks = []

    $: activeTasks = tasks.filter(t => !t.done)
    $: completedTasks = tasks.filter(t => t.done)

    let newName = ""
    let oldPassword = ""
    let newPassword = ""

    let newTaskTitle = ""
    let newTaskDescription = ""

    let editingTaskId = null
    let editingTitle = ""
    let editingDescription = ""

    function normalizeTask(task) {
        return { ...task, done: task.is_completed }
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

    onMount(async () => {
        try {
            me = await getMe()
            newName = me.name

            const rawTasks = await getMyTasks()
            tasks = rawTasks.map(normalizeTask)

            onLoaded(me.role)
        } catch (err) {
            error = err.message
        }
    })

    async function submitRename() {
        try {
            const updated = await renameMe(newName)
            me.name = updated.name
            message = "Name updated"
        } catch (err) {
            error = err.message
        }
    }

    async function submitPasswordChange() {
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
        try {
            const created = await createNewTask(newTaskTitle, newTaskDescription)
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
        try {
            await deleteTask(taskId)
            tasks = tasks.filter(t => t.id !== taskId)
        } catch (err) {
            error = err.message
        }
    }

    async function submitSwitchTask(taskId) {
        try {
            const updated = await switchTaskStatus(taskId)
            const normalized = normalizeTask(updated)
            tasks = tasks.map(t => t.id === taskId ? normalized : t)
        } catch (err) {
            error = err.message
        }
    }

    async function submitEditTask(taskId) {
        try {
            await changeTaskTitle(taskId, editingTitle)
            const updated = await changeTaskDescription(taskId, editingDescription)
            const normalized = normalizeTask(updated)
            tasks = tasks.map(t => t.id === taskId ? normalized : t)
            cancelEdit()
            message = "Task updated"
        } catch (err) {
            error = err.message
        }
    }

    async function submitDeleteAccount() {
        if (!confirm("Are you sure you want to delete your account?")) return
        try {
            await deleteMe()
            onDeleted()
        } catch (err) {
            error = err.message
        }
    }
</script>

<div class="page">
    {#if error}<p class="error">{error}</p>{/if}

    {#if !me}
        <p>Loading profile…</p>
    {:else}
        <section class="card">
            <h1>Hello, {me.name} 👋</h1>
            <p class="muted">Role: {me.role} · ID: {me.id}</p>
        </section>

        <section class="card">
            <h2>Account</h2>
            <input bind:value={newName} />
            <button on:click={submitRename}>Rename</button>

            <input type="password" placeholder="Old password" bind:value={oldPassword} />
            <input type="password" placeholder="New password" bind:value={newPassword} />
            <button on:click={submitPasswordChange}>Change password</button>
        </section>

        <section class="card">
            <h2>Create new task</h2>
            <input placeholder="Task title" bind:value={newTaskTitle} />
            <textarea placeholder="Task description" bind:value={newTaskDescription}></textarea>
            <button on:click={submitCreateTask}>Create task</button>
        </section>

        <section class="card">
            <h2>Active tasks</h2>

            {#each activeTasks as task (task.id)}
                <div class="task">
                    <button class="icon-btn toggle" on:click={() => submitSwitchTask(task.id)}>⬜</button>

                    {#if editingTaskId === task.id}
                        <input bind:value={editingTitle} />
                        <textarea bind:value={editingDescription}></textarea>
                        <button class="icon-btn edit" on:click={() => submitEditTask(task.id)}>💾</button>
                        <button class="icon-btn delete" on:click={cancelEdit}>✖</button>
                    {:else}
                        <div class="task-body">
                            <b>{task.title}</b>
                            <p>{task.description || "No description"}</p>
                        </div>
                        <button class="icon-btn edit" on:click={() => startEditTask(task)}>✏️</button>
                    {/if}

                    <button class="icon-btn delete" on:click={() => submitDeleteTask(task.id)}>✖</button>
                </div>
            {/each}
        </section>

        <section class="card">
            <h2>Completed tasks</h2>

            {#each completedTasks as task (task.id)}
                <div class="task">
                    <button class="icon-btn toggle done" on:click={() => submitSwitchTask(task.id)}>✓</button>

                    <div class="task-body">
                        <s><b>{task.title}</b></s>
                        <p>{task.description}</p>
                    </div>

                    <button class="icon-btn delete" on:click={() => submitDeleteTask(task.id)}>✖</button>
                </div>
            {/each}
        </section>

        <section class="card danger">
            <h2>Danger zone</h2>
            <button on:click={submitDeleteAccount}>Delete my account</button>
        </section>

        {#if message}<p class="success">{message}</p>{/if}
    {/if}
</div>

<style>
    .page {
        max-width: 900px;
        margin: 0 auto;
        padding: 2rem;
        color: #1f2937;
    }

    .card {
        background: #f3e0dc; /* тёплый бежево-розовый */
        border-radius: 12px;
        padding: 1.5rem;
        margin-bottom: 1.5rem;
        box-shadow: 0 10px 25px rgba(0,0,0,.06);
    }

    input,
    textarea {
        width: 100%;
        padding: 10px 12px;
        margin: 8px 0;

        background: #fdfbf8;
        border: 1px solid #e5e1db;
        border-radius: 8px;

        font-size: 14px;
        color: #1f2937;

        transition: border-color .15s, box-shadow .15s;
    }

    input::placeholder,
    textarea::placeholder {
        color: #9ca3af;
    }

    input:focus,
    textarea:focus {
        outline: none;
        border-color: #2563eb;
        box-shadow: 0 0 0 2px rgba(37,99,235,.15);
    }

    button {
        margin-top: 8px;
        padding: 10px 16px;
        border-radius: 8px;
        border: none;

        background: #2563eb;
        color: white;
        font-size: 14px;
        font-weight: 500;

        cursor: pointer;
        transition: background .15s, transform .1s;
    }

    button:hover {
        background: #1dd820;
        transform: translateY(-1px);
    }

    .task {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 12px;
    }

    .task-body p {
        margin: 4px 0 0;
        opacity: .7;
    }

    .icon-btn {
        width: 36px;
        height: 36px;
        border-radius: 8px;
        border: none;

        display: flex;
        align-items: center;
        justify-content: center;

        font-weight: 600;
        cursor: pointer;
        color: white;

        transition: all .15s ease;
    }

    .icon-btn.toggle { background: #374151; }
    .icon-btn.toggle.done { background: #16a34a; }
    .icon-btn.edit { background: #2563eb; }
    .icon-btn.delete { background: #dc2626; }

    .icon-btn:hover {
        transform: translateY(-1px);
        opacity: .9;
    }

    .muted { opacity: .6; }
    .error { color: #dc2626; }
    .success { color: #16a34a; }

    .danger {
        border: 1px solid #fecaca;
    }

    .danger button {
        background: #dc2626;
    }

    .danger button:hover {
        background: #b91c1c;
    }
</style>



