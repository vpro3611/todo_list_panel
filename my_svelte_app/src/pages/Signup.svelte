<script>
  // 1
  import { signIn } from "../api/http.js"

  let username = ""
  let password = ""
  let signIn_message = ""
  let user = null

  async function submit() {
    try {
      const result = await signIn(username, password)
      user = result
      signIn_message = "sign-in success"
      console.log(result)
    } catch (err) {
      signIn_message = err.message
    }
  }
</script>


<input
  placeholder="username"
  bind:value={username}
/>

<input
  type="password"
  placeholder="password"
  bind:value={password}
/>

<button on:click={submit}>Sign in</button>

<p>{signIn_message}</p>

{#if user}
  <h3>
    <p>Created user with id: {user.id}</p>
    <p>Status of the process: {user.status}</p>
  </h3>
{/if}
