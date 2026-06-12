<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { resetPassword } from '$lib/api/auth';
  import { showToast } from '$lib/stores/toast';
  import AuthCard from '$lib/components/AuthCard.svelte';

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
      await resetPassword($page.params.token ?? '', password);
      showToast('Password updated — log in with your new password.', 'success');
      goto('/login');
    } catch {
      errorMessage = 'This reset link is invalid or has expired. Request a new one.';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Choose a New Password — Blessed Bites</title>
</svelte:head>

<AuthCard title="Choose a new password" subtitle="Make it strong — at least 8 characters.">
  <form onsubmit={handleSubmit}>
    <label class="bb-field">
      <span>New password</span>
      <input type="password" required minlength="8" autocomplete="new-password" placeholder="At least 8 characters" bind:value={password} />
    </label>
    <label class="bb-field">
      <span>Confirm new password</span>
      <input type="password" required autocomplete="new-password" placeholder="Repeat your password" bind:value={confirmPassword} />
    </label>
    {#if errorMessage}
      <p class="bb-form-error" role="alert">{errorMessage}</p>
    {/if}
    <button type="submit" class="bb-btn-primary" disabled={submitting}>
      {submitting ? 'Updating…' : 'Update password'}
    </button>
  </form>
  {#snippet footer()}
    Need a new link? <a href="/reset-password-request">Request reset</a>
  {/snippet}
</AuthCard>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }
</style>
