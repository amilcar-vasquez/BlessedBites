<script lang="ts">
  import CartPanel from './CartPanel.svelte';

  let { open = $bindable(false) }: { open?: boolean } = $props();

  function close() {
    open = false;
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') close();
  }
</script>

<svelte:window onkeydown={onKeydown} />

{#if open}
  <div class="backdrop" onclick={close} aria-hidden="true"></div>
  <aside class="drawer" aria-label="Your order">
    <button type="button" class="close-btn" aria-label="Close cart" onclick={close}>
      <span class="material-symbols-outlined">close</span>
    </button>
    <CartPanel onnavigate={close} />
  </aside>
{/if}

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    background: color-mix(in srgb, var(--md-sys-color-scrim) 40%, transparent);
    backdrop-filter: blur(4px);
    -webkit-backdrop-filter: blur(4px);
    z-index: 90;
  }

  .drawer {
    position: fixed;
    z-index: 100;
    background: var(--md-sys-color-surface-container-lowest);
    box-shadow: var(--bb-elev-3);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  /* Desktop: right side sheet */
  @media (min-width: 769px) {
    .drawer {
      top: 0;
      right: 0;
      height: 100dvh;
      width: 400px;
      max-width: 90vw;
      border-radius: var(--bb-shape-lg) 0 0 var(--bb-shape-lg);
      animation: slide-in-right 250ms ease-out;
    }
  }

  /* Mobile: bottom sheet */
  @media (max-width: 768px) {
    .drawer {
      left: 0;
      right: 0;
      bottom: 0;
      max-height: 85dvh;
      height: 70dvh;
      border-radius: var(--bb-shape-lg) var(--bb-shape-lg) 0 0;
      animation: slide-in-up 250ms ease-out;
    }
  }

  .close-btn {
    position: absolute;
    top: 12px;
    right: 12px;
    z-index: 2;
    border: none;
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-on-surface);
    border-radius: var(--bb-shape-full);
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }

  .close-btn .material-symbols-outlined {
    font-size: 20px;
  }

  @keyframes slide-in-right {
    from {
      transform: translateX(100%);
    }
    to {
      transform: translateX(0);
    }
  }

  @keyframes slide-in-up {
    from {
      transform: translateY(100%);
    }
    to {
      transform: translateY(0);
    }
  }
</style>
