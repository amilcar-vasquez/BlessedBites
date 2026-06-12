<script lang="ts">
  import { goto } from '$app/navigation';
  import { signup } from '$lib/api/auth';
  import AuthCard from '$lib/components/AuthCard.svelte';

  let fullName = $state('');
  let email = $state('');
  let phoneNo = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let submitting = $state(false);
  let errorMessage = $state<string | null>(null);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (submitting) return;
    if (password !== confirmPassword) {
      errorMessage = 'Passwords do not match.';
      return;
    }
    submitting = true;
    errorMessage = null;
    try {
      await signup({
        email: email.trim(),
        full_name: fullName.trim(),
        phone_no: phoneNo.trim(),
        password
      });
      goto('/signup-thanks');
    } catch {
      errorMessage = 'Could not create your account. The email may already be in use.';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Sign Up — Blessed Bites</title>
</svelte:head>

<AuthCard title="Create your account" subtitle="Join Blessed Bites for faster checkout and order history.">
  <form onsubmit={handleSubmit}>
    <label class="bb-field">
      <span>Full name</span>
      <input type="text" required minlength="2" autocomplete="name" placeholder="Jane Doe" bind:value={fullName} />
    </label>
    <label class="bb-field">
      <span>Email</span>
      <input type="email" required autocomplete="email" placeholder="you@example.com" bind:value={email} />
    </label>
    <label class="bb-field">
      <span>Phone number</span>
      <input type="tel" required minlength="7" autocomplete="tel" placeholder="+1 555 000 1234" bind:value={phoneNo} />
    </label>
    <label class="bb-field">
      <span>Password</span>
      <input type="password" required minlength="8" autocomplete="new-password" placeholder="At least 8 characters" bind:value={password} />
    </label>
    <label class="bb-field">
      <span>Confirm password</span>
      <input type="password" required autocomplete="new-password" placeholder="Repeat your password" bind:value={confirmPassword} />
    </label>
    {#if errorMessage}
      <p class="bb-form-error" role="alert">{errorMessage}</p>
    {/if}
    <button type="submit" class="bb-btn-primary" disabled={submitting}>
      {submitting ? 'Creating account…' : 'Sign Up'}
    </button>
  </form>
  {#snippet footer()}
    Already have an account? <a href="/login">Login</a>
  {/snippet}
</AuthCard>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }
</style>
