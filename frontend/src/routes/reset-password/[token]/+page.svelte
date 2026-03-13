<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { resetPassword } from '$lib/api/auth';

  let password = '';
  let confirm = '';
  let error = '';
  let submitting = false;

  async function submit() {
    error = '';

    if (password !== confirm) {
      error = 'Passwords do not match.';
      return;
    }

    submitting = true;
    try {
      await resetPassword($page.params.token, password);
      await goto('/login');
    } catch (e) {
      console.error(e);
      error = 'Reset failed. Token may be expired or invalid.';
    } finally {
      submitting = false;
    }
  }
</script>

<main class="auth-shell">
  <section class="card">
    <h1>Set New Password</h1>
    <p>Choose a secure password for your account.</p>

    <form on:submit|preventDefault={submit}>
      <label for="password">New password</label>
      <input id="password" type="password" bind:value={password} required minlength="8" />

      <label for="confirm">Confirm password</label>
      <input id="confirm" type="password" bind:value={confirm} required minlength="8" />

      {#if error}<p class="error">{error}</p>{/if}

      <button class="btn" type="submit" disabled={submitting}>
        {submitting ? 'Updating...' : 'Reset Password'}
      </button>
    </form>
  </section>
</main>

<style>
  .auth-shell { max-width: 520px; margin: 2rem auto; padding: 1rem; }
  .card { background: #fffaf5; border: 1px solid #e7d3c9; border-radius: 20px; padding: 1rem 1.2rem; }
  form { display: grid; gap: 0.55rem; }
  label { font-weight: 700; font-size: 0.9rem; }
  input { border: 1px solid #d9c5bc; border-radius: 12px; padding: 0.65rem 0.8rem; }
  .btn { margin-top: 0.35rem; border: none; border-radius: 999px; padding: 0.65rem 1rem; background: #7f1d2d; color: #fff; font-weight: 700; }
  .error { color: #8a1732; }
</style>
