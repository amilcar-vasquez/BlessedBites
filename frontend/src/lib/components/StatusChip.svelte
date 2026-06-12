<script lang="ts">
  let { status }: { status: string } = $props();

  const normalized = $derived((status || '').toLowerCase());

  const config = $derived.by(() => {
    switch (normalized) {
      case 'pending':
        return { label: 'Pending', kind: 'tertiary', icon: 'schedule' };
      case 'preparing':
        return { label: 'Preparing', kind: 'secondary', icon: '' };
      case 'ready':
        return { label: 'Ready', kind: 'tertiary', icon: '' };
      case 'completed':
      case 'delivered':
        return { label: status, kind: 'neutral', icon: 'check' };
      case 'cancelled':
        return { label: 'Cancelled', kind: 'error', icon: 'close' };
      default:
        return { label: status || 'Unknown', kind: 'neutral', icon: '' };
    }
  });
</script>

<span class="chip label-sm {config.kind}">
  {#if config.icon}
    <span class="material-symbols-outlined" aria-hidden="true">{config.icon}</span>
  {:else}
    <span class="dot" aria-hidden="true"></span>
  {/if}
  {config.label}
</span>

<style>
  .chip {
    display: inline-flex;
    align-items: center;
    gap: var(--bb-space-xs);
    padding: 4px 10px;
    border-radius: var(--bb-shape-full);
    white-space: nowrap;
    text-transform: capitalize;
  }

  .dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: currentColor;
  }

  .material-symbols-outlined {
    font-size: 12px;
  }

  .secondary {
    background: color-mix(in srgb, var(--md-sys-color-secondary) 12%, transparent);
    color: var(--md-sys-color-secondary);
  }

  .tertiary {
    background: color-mix(in srgb, var(--md-sys-color-tertiary) 12%, transparent);
    color: var(--md-sys-color-tertiary);
  }

  .neutral {
    background: var(--md-sys-color-surface-variant);
    color: var(--md-sys-color-on-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .error {
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-on-error-container);
  }
</style>
