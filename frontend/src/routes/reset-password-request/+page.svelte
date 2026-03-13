<script lang="ts">
  import { requestPasswordReset } from '$lib/api/auth';

  let email = '';
  let message = '';
  let error = '';
  let submitting = false;

  async function submit() {
    error = '';
    message = '';
    submitting = true;
    try {
      const res = await requestPasswordReset(email);
      message = res.message;
    } catch (e) {
      console.error(e);
      error = 'Could not submit reset request.';
    } finally {
      submitting = false;
    }
  }
</script>

<main class="auth-shell">
  <section class="card">
    <h1>Reset Password</h1>
    <p>Enter your account email to request a reset token.</p>

    <form on:submit|preventDefault={submit}>
      <label for="email">Email</label>
      <input id="email" type="email" bind:value={email} required />

      {#if message}<p class="ok">{message}</p>{/if}
      {#if error}<p class="error">{error}</p>{/if}

      <button class="btn" type="submit" disabled={submitting}>
        {submitting ? 'Sending...' : 'Send Reset Request'}
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
  .ok { color: #225f38; }
</style>
