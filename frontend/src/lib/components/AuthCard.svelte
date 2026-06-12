<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    title,
    subtitle,
    children,
    footer
  }: {
    title: string;
    subtitle?: string;
    children: Snippet;
    footer?: Snippet;
  } = $props();
</script>

<div class="wrap">
  <div class="bb-card card">
    <a class="brand" href="/">Blessed Bites</a>
    <h1 class="headline-md">{title}</h1>
    {#if subtitle}
      <p class="body-md subtitle">{subtitle}</p>
    {/if}
    {@render children()}
    {#if footer}
      <div class="footer body-md">
        {@render footer()}
      </div>
    {/if}
  </div>
</div>

<style>
  .wrap {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--bb-space-xl) var(--bb-margin-mobile);
    background:
      radial-gradient(ellipse 70% 50% at 85% 10%, color-mix(in srgb, var(--md-sys-color-tertiary-container) 40%, transparent), transparent),
      var(--md-sys-color-background);
  }

  .card {
    width: 100%;
    max-width: 440px;
    padding: var(--bb-space-xl);
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }

  .brand {
    font-family: var(--md-ref-typeface-brand);
    font-weight: 700;
    font-size: 20px;
    color: var(--md-sys-color-primary);
    text-decoration: none;
  }

  h1 {
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  .subtitle {
    margin: calc(-1 * var(--bb-space-sm)) 0 0;
    color: var(--md-sys-color-on-surface-variant);
  }

  .footer {
    text-align: center;
    color: var(--md-sys-color-on-surface-variant);
    border-top: 1px solid var(--md-sys-color-outline-variant);
    padding-top: var(--bb-space-md);
  }

  .footer :global(a) {
    color: var(--md-sys-color-primary);
    font-weight: 600;
  }
</style>
