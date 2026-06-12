<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { login } from '$lib/api/auth';
  import { setSession } from '$lib/stores/auth';
  import { showToast } from '$lib/stores/toast';
  import AuthCard from '$lib/components/AuthCard.svelte';

  let email = $state('');
  let password = $state('');
  let submitting = $state(false);
  let errorMessage = $state<string | null>(null);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (submitting) return;
    submitting = true;
    errorMessage = null;
    try {
      const res = await login(email.trim(), password);
      setSession(res.token, res.user);
      showToast(`Welcome back, ${res.user.full_name.split(' ')[0]}!`, 'success', 2500);
      const redirect = $page.url.searchParams.get('redirect');
      goto(redirect && redirect.startsWith('/') ? redirect : res.user.role === 'admin' ? '/dashboard' : '/menu');
    } catch {
      errorMessage = 'Invalid email or password.';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Login — Blessed Bites</title>
</svelte:head>

<AuthCard title="Welcome back" subtitle="Log in to order faster and track your meals.">
  <form onsubmit={handleSubmit}>
    <label class="bb-field">
      <span>Email</span>
      <input type="email" required autocomplete="email" placeholder="you@example.com" bind:value={email} />
    </label>
    <label class="bb-field">
      <span>Password</span>
      <input type="password" required autocomplete="current-password" placeholder="••••••••" bind:value={password} />
    </label>
    <a class="forgot label-lg" href="/reset-password-request">Forgot password?</a>
    {#if errorMessage}
      <p class="bb-form-error" role="alert">{errorMessage}</p>
    {/if}
    <button type="submit" class="bb-btn-primary" disabled={submitting}>
      {submitting ? 'Logging in…' : 'Login'}
    </button>
  </form>
  {#snippet footer()}
    New here? <a href="/signup">Create an account</a>
  {/snippet}
</AuthCard>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }

  .forgot {
    align-self: flex-end;
    color: var(--md-sys-color-primary);
    text-decoration: none;
  }

  .forgot:hover {
    text-decoration: underline;
  }
</style>
