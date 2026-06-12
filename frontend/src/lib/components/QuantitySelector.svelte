<script lang="ts">
  let {
    qty,
    onchange
  }: { qty: number; onchange: (qty: number) => void } = $props();
</script>

<div class="qty-selector">
  <button
    type="button"
    aria-label="Decrease quantity"
    onclick={() => onchange(qty - 1)}
  >
    <span class="material-symbols-outlined">remove</span>
  </button>
  <span class="label-lg count">{qty}</span>
  <button type="button" aria-label="Increase quantity" onclick={() => onchange(qty + 1)}>
    <span class="material-symbols-outlined">add</span>
  </button>
</div>

<style>
  .qty-selector {
    display: inline-flex;
    align-items: center;
    border: 1px solid var(--md-sys-color-outline);
    border-radius: var(--bb-shape-full);
    background: var(--md-sys-color-surface);
    overflow: hidden;
  }

  button {
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    color: var(--md-sys-color-primary);
    padding: 4px 8px;
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  button:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .material-symbols-outlined {
    font-size: 16px;
  }

  .count {
    min-width: 24px;
    text-align: center;
    color: var(--md-sys-color-on-surface);
  }
</style>
