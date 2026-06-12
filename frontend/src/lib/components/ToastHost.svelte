<script lang="ts">
  import { toasts } from '$lib/stores/toast';
</script>

<div class="toast-host" role="status" aria-live="polite">
  {#each $toasts as toast (toast.id)}
    <div class="toast label-lg" class:error={toast.kind === 'error'} class:success={toast.kind === 'success'}>
      <span class="material-symbols-outlined" aria-hidden="true">
        {toast.kind === 'error' ? 'error' : toast.kind === 'success' ? 'check_circle' : 'info'}
      </span>
      {toast.message}
    </div>
  {/each}
</div>

<style>
  .toast-host {
    position: fixed;
    bottom: var(--bb-space-lg);
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-sm);
    z-index: 1000;
    pointer-events: none;
  }

  .toast {
    display: flex;
    align-items: center;
    gap: var(--bb-space-sm);
    background: var(--md-sys-color-inverse-surface);
    color: var(--md-sys-color-inverse-on-surface);
    padding: 12px 20px;
    border-radius: var(--bb-shape-sm);
    box-shadow: var(--bb-elev-3);
    animation: slide-up 200ms ease-out;
  }

  .toast.error {
    background: var(--md-sys-color-error);
    color: var(--md-sys-color-on-error);
  }

  .toast.success {
    background: var(--md-sys-color-secondary);
    color: var(--md-sys-color-on-secondary);
  }

  .toast .material-symbols-outlined {
    font-size: 18px;
  }

  @keyframes slide-up {
    from {
      opacity: 0;
      transform: translateY(8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
