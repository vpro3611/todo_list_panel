<script>
  import { login } from "../api/http.js"

  export let onSuccess

  let username = ""
  let password = ""
  let message = ""

  async function submit() {
    try {
      const { token } = await login(username, password)
      localStorage.setItem("token", token)

      onSuccess()   
    } catch (err) {
      message = err.message
    }
  }
</script>

<input placeholder="username" bind:value={username} />

<input type="password" placeholder="password" bind:value={password} />


<button on:click={submit}>login</button>

<p>{message}</p>



