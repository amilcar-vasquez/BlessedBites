<script lang="ts">
  import { goto } from '$app/navigation';
  import { login } from '$lib/api/auth';

  let email = '';
  let password = '';
  let error = '';
  let submitting = false;

  async function submit() {
    error = '';
    submitting = true;
    try {
      const result = await login(email, password);
      localStorage.setItem('bb_access_token', result.token);
      localStorage.setItem('bb_user', JSON.stringify(result.user));
      await goto('/');
    } catch (e) {
      console.error(e);
      error = 'Invalid credentials. Please try again.';
    } finally {
      submitting = false;
    }
  }
</script>

<main class="auth-shell">
  <section class="card">
    <h1>Login</h1>
    <p>Welcome back to BlessedBites.</p>

    <form on:submit|preventDefault={submit}>
      <label for="email">Email</label>
      <input id="email" type="email" bind:value={email} required />

      <label for="password">Password</label>
      <input id="password" type="password" bind:value={password} required />

      {#if error}<p class="error">{error}</p>{/if}

      <button class="btn" type="submit" disabled={submitting}>
        {submitting ? 'Logging in...' : 'Login'}
      </button>
    </form>

    <div class="links">
      <a href="/reset-password-request">Forgot password?</a>
      <a href="/signup">Create account</a>
    </div>
  </section>
</main>

<style>
  .auth-shell { max-width: 520px; margin: 2rem auto; padding: 1rem; }
  .card { background: #fffaf5; border: 1px solid #e7d3c9; border-radius: 20px; padding: 1rem 1.2rem; }
  h1 { margin: 0 0 0.4rem; }
  p { margin: 0 0 1rem; color: #6f4744; }
  form { display: grid; gap: 0.55rem; }
  label { font-weight: 700; font-size: 0.9rem; }
  input { border: 1px solid #d9c5bc; border-radius: 12px; padding: 0.65rem 0.8rem; }
  .btn { margin-top: 0.35rem; border: none; border-radius: 999px; padding: 0.65rem 1rem; background: #7f1d2d; color: #fff; font-weight: 700; }
  .error { color: #8a1732; margin: 0.15rem 0; }
  .links { margin-top: 0.9rem; display: flex; justify-content: space-between; }
</style>
