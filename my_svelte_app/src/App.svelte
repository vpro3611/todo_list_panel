<script>
  import Login from "./pages/Login.svelte"
  import Signup from "./pages/Signup.svelte"

  // какое окно сейчас показываем
  let page = "login" // "login" | "signup" | "app"

  // проверяем токен при старте приложения
  const token = localStorage.getItem("token")
  if (token) {
    page = "app"
  }

  function logout() {
    localStorage.removeItem("token")
    page = "login"
  }
</script>

<style>
  nav {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
  }

  button {
    padding: 6px 12px;
  }
</style>

<!-- Навигация -->
<nav>
  {#if page !== "app"}
    <button on:click={() => page = "login"}>Login</button>
    <button on:click={() => page = "signup"}>Signup</button>
  {:else}
    <button on:click={logout}>Logout</button>
  {/if}
</nav>

<!-- Рендер страниц -->
{#if page === "login"}
  <Login onSuccess={() => page = "app"} />

{:else if page === "signup"}
  <Signup on:done={() => page = "login"} />

{:else if page === "app"}
  <h1>You are logged in ✅</h1>
  <p>Welcome to the app.</p>
{/if}
