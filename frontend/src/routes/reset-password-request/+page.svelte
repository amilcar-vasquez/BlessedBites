<script lang="ts">
  import { requestPasswordReset } from '$lib/api/auth';
  import AuthCard from '$lib/components/AuthCard.svelte';

  let email = $state('');
  let submitting = $state(false);
  let sent = $state(false);
  let errorMessage = $state<string | null>(null);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (submitting) return;
    submitting = true;
    errorMessage = null;
    try {
      await requestPasswordReset(email.trim());
      sent = true;
    } catch {
      errorMessage = 'Something went wrong. Please try again.';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Reset Password — Blessed Bites</title>
</svelte:head>

<AuthCard
  title="Reset your password"
  subtitle="Enter your email and we'll send you a reset link."
>
  {#if sent}
    <div class="sent">
      <span class="material-symbols-outlined icon" aria-hidden="true">mark_email_read</span>
      <p class="body-lg">
        If an account exists for <strong>{email}</strong>, a reset link is on its way. Check your inbox.
      </p>
    </div>
  {:else}
    <form onsubmit={handleSubmit}>
      <label class="bb-field">
        <span>Email</span>
        <input type="email" required autocomplete="email" placeholder="you@example.com" bind:value={email} />
      </label>
      {#if errorMessage}
        <p class="bb-form-error" role="alert">{errorMessage}</p>
      {/if}
      <button type="submit" class="bb-btn-primary" disabled={submitting}>
        {submitting ? 'Sending…' : 'Send reset link'}
      </button>
    </form>
  {/if}
  {#snippet footer()}
    Remembered it? <a href="/login">Back to Login</a>
  {/snippet}
</AuthCard>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }

  .sent {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: var(--bb-space-md);
  }

  .icon {
    font-size: 48px;
    color: var(--md-sys-color-secondary);
  }

  .sent p {
    margin: 0;
    color: var(--md-sys-color-on-surface-variant);
  }
</style>
