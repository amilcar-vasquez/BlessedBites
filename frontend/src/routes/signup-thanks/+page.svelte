<script lang="ts">
  import AuthCard from '$lib/components/AuthCard.svelte';
</script>

<svelte:head>
  <title>Welcome — Blessed Bites</title>
</svelte:head>

<AuthCard title="You're in!" subtitle="Your account has been created successfully.">
  <div class="content">
    <span class="material-symbols-outlined fill icon" aria-hidden="true">celebration</span>
    <p class="body-lg">Welcome to the Blessed Bites family. Log in to start ordering.</p>
    <a class="bb-btn-primary" href="/login">Continue to Login</a>
  </div>
</AuthCard>

<style>
  .content {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: var(--bb-space-md);
  }

  .icon {
    font-size: 56px;
    color: var(--md-sys-color-tertiary);
  }

  p {
    margin: 0;
    color: var(--md-sys-color-on-surface-variant);
  }
</style>
