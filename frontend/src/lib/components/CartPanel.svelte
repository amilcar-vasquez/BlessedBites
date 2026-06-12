<script lang="ts">
  import { goto } from '$app/navigation';
  import { cart, cartCount, cartTotal, setQty } from '$lib/stores/cart';
  import QuantitySelector from './QuantitySelector.svelte';

  let { onnavigate }: { onnavigate?: () => void } = $props();

  function checkout() {
    onnavigate?.();
    goto('/checkout');
  }
</script>

<div class="cart-panel">
  <div class="panel-header">
    <h2 class="title-lg">
      <span class="material-symbols-outlined" aria-hidden="true">shopping_basket</span>
      Your Order
    </h2>
    <span class="count-badge label-sm">{$cartCount} item{$cartCount === 1 ? '' : 's'}</span>
  </div>

  <div class="items">
    {#if $cart.length === 0}
      <div class="empty">
        <span class="material-symbols-outlined" aria-hidden="true">shopping_basket</span>
        <p class="body-md">Your order is empty.</p>
        <p class="label-sm">Add something delicious from the menu.</p>
      </div>
    {:else}
      {#each $cart as item (item.id)}
        <div class="cart-item">
          {#if item.image_url}
            <img src={item.image_url} alt={item.name} />
          {:else}
            <div class="thumb-placeholder">
              <span class="material-symbols-outlined" aria-hidden="true">restaurant</span>
            </div>
          {/if}
          <div class="item-body">
            <div class="item-top">
              <h4 class="label-lg name">{item.name}</h4>
              <span class="label-lg price">${(item.price * item.qty).toFixed(2)}</span>
            </div>
            <div class="item-actions">
              <QuantitySelector qty={item.qty} onchange={(q) => setQty(item.id, q)} />
            </div>
          </div>
        </div>
      {/each}
    {/if}
  </div>

  <div class="panel-footer">
    <div class="row body-md">
      <span class="muted">Subtotal</span>
      <span>${$cartTotal.toFixed(2)}</span>
    </div>
    <div class="row body-md">
      <span class="muted">Taxes &amp; Fees</span>
      <span>Calculated next</span>
    </div>
    <div class="row total">
      <span class="title-md">Total</span>
      <span class="title-lg total-price">${$cartTotal.toFixed(2)}</span>
    </div>
    <button
      type="button"
      class="checkout-btn label-lg"
      disabled={$cart.length === 0}
      onclick={checkout}
    >
      Checkout
      <span class="material-symbols-outlined" aria-hidden="true">arrow_forward</span>
    </button>
  </div>
</div>

<style>
  .cart-panel {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--md-sys-color-surface-container-lowest);
    overflow: hidden;
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: var(--bb-space-md);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface-container-low);
  }

  .panel-header h2 {
    margin: 0;
    display: flex;
    align-items: center;
    gap: var(--bb-space-sm);
    color: var(--md-sys-color-on-surface);
  }

  .count-badge {
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-on-primary-container);
    padding: 4px 8px;
    border-radius: var(--bb-shape-full);
  }

  .items {
    flex: 1;
    overflow-y: auto;
    padding: var(--bb-space-md);
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }

  .empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    color: var(--md-sys-color-on-surface-variant);
    gap: var(--bb-space-xs);
    padding: var(--bb-space-xl) var(--bb-space-md);
  }

  .empty .material-symbols-outlined {
    font-size: 40px;
    opacity: 0.5;
  }

  .empty p {
    margin: 0;
  }

  .cart-item {
    display: flex;
    gap: var(--bb-space-sm);
  }

  .cart-item img,
  .thumb-placeholder {
    width: 64px;
    height: 64px;
    border-radius: 0.75rem;
    object-fit: cover;
    flex-shrink: 0;
  }

  .thumb-placeholder {
    background: var(--md-sys-color-surface-container-high);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--md-sys-color-on-surface-variant);
  }

  .item-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    min-width: 0;
  }

  .item-top {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: var(--bb-space-sm);
  }

  .name {
    margin: 0;
    color: var(--md-sys-color-on-surface);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .price {
    color: var(--md-sys-color-primary);
    white-space: nowrap;
  }

  .item-actions {
    margin-top: var(--bb-space-xs);
  }

  .panel-footer {
    padding: var(--bb-space-md);
    background: var(--md-sys-color-surface-container-low);
    border-top: 1px solid var(--md-sys-color-outline-variant);
  }

  .row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--bb-space-sm);
  }

  .muted {
    color: var(--md-sys-color-on-surface-variant);
  }

  .row.total {
    margin: var(--bb-space-md) 0 var(--bb-space-lg);
  }

  .total-price {
    color: var(--md-sys-color-primary);
  }

  .checkout-btn {
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: var(--bb-space-sm);
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    border: none;
    border-radius: var(--bb-shape-sm);
    padding: 12px;
    cursor: pointer;
    box-shadow: var(--bb-elev-2);
    transition: background-color 150ms ease;
  }

  .checkout-btn:hover:not(:disabled) {
    background: var(--md-sys-color-surface-tint);
  }

  .checkout-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .checkout-btn .material-symbols-outlined {
    font-size: 20px;
  }
</style>
