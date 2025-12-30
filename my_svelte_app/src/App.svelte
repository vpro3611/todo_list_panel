<script>
    import Login from "./pages/Login.svelte"
    import Signup from "./pages/Signup.svelte"
    import Me from "./pages/Me.svelte"
    import AdminUsers from "./pages/AdminUsers.svelte"
    import { getMe } from "./api/http.js"

    // какое представление показываем
    // "login" | "signup" | "me" | "admin"
    let page = "login"

    // текущий пользователь (глобальное состояние)
    let me = null
    let loading = true

    // при старте приложения
    const token = localStorage.getItem("token")

    if (token) {
        loadMe()
    } else {
        loading = false
    }

    async function loadMe() {
        try {
            me = await getMe()
            page = "me"
        } catch (err) {
            // токен невалиден или истёк
            localStorage.removeItem("token")
            me = null
            page = "login"
        } finally {
            loading = false
        }
    }

    function logout() {
        localStorage.removeItem("token")
        me = null
        page = "login"
    }
</script>

<style>
    nav {
        display: flex;
        gap: 10px;
        padding: 10px;
        border-bottom: 1px solid #ccc;
        margin-bottom: 20px;
    }

    button {
        padding: 6px 12px;
        cursor: pointer;
    }
</style>

{#if loading}
    <p>Loading...</p>

{:else}

    <!-- NAVIGATION -->
    <nav>
        {#if !me}
            <button on:click={() => page = "login"}>Login</button>
            <button on:click={() => page = "signup"}>Signup</button>
        {:else}
            <button on:click={() => page = "me"}>Home</button>

            {#if me.role === "admin"}
                <button on:click={() => page = "admin"}>Admin</button>
            {/if}

            <button on:click={logout}>Logout</button>
        {/if}
    </nav>

    <!-- PAGES -->
    {#if page === "login"}
        <Login onSuccess={loadMe} />

    {:else if page === "signup"}
        <Signup on:done={() => page = "login"} />

    {:else if page === "me"}
        <Me onDeleted="{logout}"/>

    {:else if page === "admin"}
        <AdminUsers />
    {/if}

{/if}
