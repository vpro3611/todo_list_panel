<!--<script>-->
<!--    import { login } from "../api/http.js"-->

<!--    export let onSuccess-->

<!--    let username = ""-->
<!--    let password = ""-->
<!--    let error = ""-->

<!--    async function submit() {-->
<!--        error = ""-->

<!--        try {-->
<!--            const { token } = await login(username, password)-->
<!--            localStorage.setItem("token", token)-->
<!--            onSuccess()-->
<!--        } catch (err) {-->
<!--            error = err.message-->
<!--        }-->
<!--    }-->
<!--</script>-->

<!--<input placeholder="username" bind:value={username} />-->
<!--<input type="password" placeholder="password" bind:value={password} />-->
<!--<button on:click={submit}>Login</button>-->

<!--{#if error}-->
<!--    <p style="color:red">{error}</p>-->
<!--{/if}-->


<script>
    import { login } from "../api/http.js"

    export let onSuccess

    let username = ""
    let password = ""
    let error = ""
    let loading = false

    async function submit() {
        error = ""
        loading = true

        try {
            const { token } = await login(username, password)
            localStorage.setItem("token", token)
            onSuccess()
        } catch (err) {
            error = err.message
        } finally {
            loading = false
        }
    }
</script>

<div class="screen">
    <div class="card">
        <h2>Welcome back 👋</h2>

        <input
                placeholder="Username"
                bind:value={username}
        />

        <input
                type="password"
                placeholder="Password"
                bind:value={password}
        />

        <button on:click={submit} disabled={loading}>
            {loading ? "Logging in..." : "Login"}
        </button>

        {#if error}
            <p class="error">{error}</p>
        {/if}
    </div>
</div>

<style>
    /* === FULL SCREEN === */
    .screen {
        min-height: 100vh;
        width: 100%;
        padding: 16px;

        display: flex;
        align-items: center;
        justify-content: center;

        background:
                radial-gradient(circle at top left, #6366f1, transparent 40%),
                radial-gradient(circle at bottom right, #22d3ee, transparent 40%),
                #0f172a;
    }

    /* === CARD === */
    .card {
        width: 100%;
        max-width: 360px;

        padding: 24px;
        border-radius: 16px;

        background: rgba(255, 255, 255, 0.96);
        color: #111827; /* 🔑 ВАЖНО */

        box-shadow:
                0 20px 40px rgba(0, 0, 0, 0.25),
                0 0 0 1px rgba(255, 255, 255, 0.4) inset;

        display: flex;
        flex-direction: column;
        gap: 14px;

        animation: float 4s ease-in-out infinite;
    }

    @keyframes float {
        0%, 100% { transform: translateY(0); }
        50% { transform: translateY(-6px); }
    }

    /* === TEXT === */
    .card h2 {
        text-align: center;
        font-weight: 700;
        margin-bottom: 6px;
    }

    /* === INPUTS === */
    .card input {
        padding: 11px 12px;
        font-size: 14px;

        border-radius: 10px;
        border: 1px solid #c7c7c7;
        outline: none;

        transition: all 0.2s ease;
    }

    .card input:focus {
        border-color: #6366f1;
        box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.25);
        transform: scale(1.02);
    }

    /* === BUTTON === */
    .card button {
        margin-top: 6px;
        padding: 11px;

        border-radius: 10px;
        border: none;

        background: linear-gradient(135deg, #6366f1, #22d3ee);
        color: white;

        font-weight: 700;
        cursor: pointer;

        transition: all 0.2s ease;
    }

    .card button:hover:not(:disabled) {
        transform: translateY(-2px);
        box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
    }

    .card button:disabled {
        opacity: 0.7;
        cursor: default;
    }

    /* === ERROR === */
    .error {
        text-align: center;
        font-size: 14px;
        color: #dc2626;
        animation: fadeIn 0.3s ease;
    }

    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(4px); }
        to   { opacity: 1; transform: translateY(0); }
    }

    /* === MOBILE === */
    @media (max-width: 480px) {
        .card {
            padding: 20px;
            border-radius: 14px;
        }

        .card h2 {
            font-size: 18px;
        }
    }
</style>
