<script lang="ts">
  import { goto } from '$app/navigation';
  import { signup } from '$lib/api/auth';

  let fullName = '';
  let email = '';
  let phoneNo = '';
  let password = '';
  let confirmPassword = '';
  let error = '';
  let submitting = false;

  async function submit() {
    error = '';

    if (password !== confirmPassword) {
      error = 'Passwords do not match.';
      return;
    }

    submitting = true;
    try {
      await signup({
        full_name: fullName,
        email,
        phone_no: phoneNo,
        password
      });
      await goto('/signup-thanks');
    } catch (e) {
      console.error(e);
      error = 'Could not create account. Please verify details and try again.';
    } finally {
      submitting = false;
    }
  }
</script>

<main class="auth-shell">
  <section class="card">
    <h1>Sign Up</h1>
    <p>Create your BlessedBites account.</p>

    <form on:submit|preventDefault={submit}>
      <label for="fullName">Full name</label>
      <input id="fullName" type="text" bind:value={fullName} required />

      <label for="email">Email</label>
      <input id="email" type="email" bind:value={email} required />

      <label for="phoneNo">Phone number</label>
      <input id="phoneNo" type="tel" bind:value={phoneNo} required />

      <label for="password">Password</label>
      <input id="password" type="password" bind:value={password} required minlength="8" />

      <label for="confirm">Confirm password</label>
      <input id="confirm" type="password" bind:value={confirmPassword} required minlength="8" />

      {#if error}<p class="error">{error}</p>{/if}

      <button class="btn" type="submit" disabled={submitting}>
        {submitting ? 'Creating account...' : 'Create Account'}
      </button>
    </form>
  </section>
</main>

<style>
  .auth-shell { max-width: 560px; margin: 2rem auto; padding: 1rem; }
  .card { background: #fffaf5; border: 1px solid #e7d3c9; border-radius: 20px; padding: 1rem 1.2rem; }
  h1 { margin: 0 0 0.4rem; }
  p { margin: 0 0 1rem; color: #6f4744; }
  form { display: grid; gap: 0.55rem; }
  label { font-weight: 700; font-size: 0.9rem; }
  input { border: 1px solid #d9c5bc; border-radius: 12px; padding: 0.65rem 0.8rem; }
  .btn { margin-top: 0.35rem; border: none; border-radius: 999px; padding: 0.65rem 1rem; background: #7f1d2d; color: #fff; font-weight: 700; }
  .error { color: #8a1732; margin: 0.15rem 0; }
</style>
