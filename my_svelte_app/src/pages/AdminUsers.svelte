<script>
    import { getAllUsersAdmin, createNewUserAdmin } from "../api/http.js"

    let users = []
    let error = ""

    let newName = ""
    let newPassword = ""

    async function loadUsers() {
        try {
            users = await getAllUsersAdmin()
        } catch (err) {
            error = err.message
        }
    }

    async function createUser() {
        try {
            await createNewUserAdmin(newName, newPassword)
            newName = ""
            newPassword = ""
            await loadUsers() // обновляем список
        } catch (err) {
            error = err.message
        }
    }
</script>

<button on:click={loadUsers}>Load users</button>

{#if error}
    <p style="color:red">{error}</p>
{/if}

<ul>
    {#each users as user}
        <li>{user.id} — {user.name} ({user.role})</li>
    {/each}
</ul>

<h3>Create new user</h3>

<input placeholder="name" bind:value={newName} />
<input placeholder="password" bind:value={newPassword} />

<button on:click={createUser}>Create</button>
