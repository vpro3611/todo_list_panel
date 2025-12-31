<script>
    import Login from "./pages/Login.svelte"
    import Signup from "./pages/Signup.svelte"
    import Me from "./pages/Me.svelte"
    import AdminUsers from "./pages/AdminUsers.svelte"

    // pages: login | signup | me | admin
    let page = "login"
    let isAuth = false
    let userRole = null

    // init
    if (localStorage.getItem("token")) {
        isAuth = true
        page = "me"
    }

    function onLoginSuccess() {
        isAuth = true
        page = "me"
    }

    function logout() {
        localStorage.removeItem("token")
        isAuth = false
        userRole = null
        page = "login"
    }

    function onMeLoaded(role) {
        userRole = role
    }
</script>

<nav style="display:flex; gap:10px; margin-bottom:20px">
    {#if !isAuth}
        <button on:click={() => page = "login"}>Login</button>
        <button on:click={() => page = "signup"}>Signup</button>
    {:else}
        <button on:click={() => page = "me"}>My profile</button>

        {#if userRole === "admin"}
            <button on:click={() => page = "admin"}>Admin panel</button>
        {/if}

        <button on:click={logout}>Logout</button>
    {/if}
</nav>

{#if page === "login"}
    <Login onSuccess={onLoginSuccess} />

{:else if page === "signup"}
    <Signup />

{:else}
    <div class="app-container">
        {#if page === "me"}
            <Me onLoaded={onMeLoaded} onDeleted={logout} />
        {:else if page === "admin"}
            <AdminUsers />
        {/if}
    </div>
{/if}