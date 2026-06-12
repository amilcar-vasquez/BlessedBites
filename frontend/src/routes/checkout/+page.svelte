<script lang="ts">
  import { goto } from '$app/navigation';
  import { createOrder } from '$lib/api/orders';
  import { auth } from '$lib/stores/auth';
  import { cart, cartTotal, clearCart, setQty } from '$lib/stores/cart';
  import { showToast } from '$lib/stores/toast';
  import QuantitySelector from '$lib/components/QuantitySelector.svelte';

  let fullName = $state('');
  let phoneNo = $state('');
  let deliveryOption = $state<'pickup' | 'delivery'>('pickup');
  let paymentOption = $state<'cash' | 'card'>('cash');
  let submitting = $state(false);
  let errorMessage = $state<string | null>(null);

  $effect(() => {
    if (!fullName && $auth.user?.full_name) {
      fullName = $auth.user.full_name;
    }
  });

  const canSubmit = $derived(
    $cart.length > 0 && fullName.trim().length > 1 && phoneNo.trim().length >= 7 && !submitting
  );

  async function placeOrder(e: Event) {
    e.preventDefault();
    if (!canSubmit) return;
    submitting = true;
    errorMessage = null;
    try {
      const result = await createOrder({
        items: $cart.map((i) => ({ id: i.id, qty: i.qty })),
        full_name: fullName.trim(),
        phone_no: phoneNo.trim()
      });
      // Stash a snapshot so the success page can offer per-item ratings.
      sessionStorage.setItem(
        'bb_last_order',
        JSON.stringify({
          order_id: result.order_id,
          total: $cartTotal,
          items: $cart.map((i) => ({ id: i.id, name: i.name, qty: i.qty }))
        })
      );
      clearCart();
      showToast('Order placed — thank you!', 'success');
      goto(`/order-success?id=${result.order_id}`);
    } catch {
      errorMessage = 'We could not place your order. Please check your details and try again.';
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>Checkout — Blessed Bites</title>
</svelte:head>

<div class="page">
  <header class="page-head">
    <h1 class="headline-lg">Checkout</h1>
    <p class="body-md muted">Almost there — confirm your order details below.</p>
  </header>

  {#if $cart.length === 0}
    <div class="bb-card empty">
      <span class="material-symbols-outlined" aria-hidden="true">shopping_cart_off</span>
      <p class="body-lg">Your cart is empty.</p>
      <a class="primary-btn label-lg" href="/menu">Browse the menu</a>
    </div>
  {:else}
    <div class="layout">
      <form class="form" onsubmit={placeOrder}>
        <!-- Contact -->
        <section class="bb-card block">
          <h2 class="title-lg"><span class="step label-sm">1</span> Contact details</h2>
          {#if !$auth.user}
            <p class="body-md muted">
              Checking out as a guest. <a href="/login">Log in</a> to track your orders.
            </p>
          {/if}
          <label class="field">
            <span class="label-lg">Full name</span>
            <input
              type="text"
              required
              minlength="2"
              autocomplete="name"
              placeholder="Jane Doe"
              bind:value={fullName}
            />
          </label>
          <label class="field">
            <span class="label-lg">Phone number</span>
            <input
              type="tel"
              required
              minlength="7"
              autocomplete="tel"
              placeholder="+1 555 000 1234"
              bind:value={phoneNo}
            />
          </label>
        </section>

        <!-- Fulfilment -->
        <section class="bb-card block">
          <h2 class="title-lg"><span class="step label-sm">2</span> Fulfilment</h2>
          <div class="options">
            <label class="option" class:selected={deliveryOption === 'pickup'}>
              <input type="radio" name="delivery" value="pickup" bind:group={deliveryOption} />
              <span class="material-symbols-outlined" aria-hidden="true">storefront</span>
              <span class="option-text">
                <span class="title-md">Pickup</span>
                <span class="label-sm muted">Ready in ~20 min</span>
              </span>
            </label>
            <label class="option" class:selected={deliveryOption === 'delivery'}>
              <input type="radio" name="delivery" value="delivery" bind:group={deliveryOption} />
              <span class="material-symbols-outlined" aria-hidden="true">local_shipping</span>
              <span class="option-text">
                <span class="title-md">Delivery</span>
                <span class="label-sm muted">Coming soon</span>
              </span>
            </label>
          </div>
        </section>

        <!-- Payment -->
        <section class="bb-card block">
          <h2 class="title-lg"><span class="step label-sm">3</span> Payment</h2>
          <div class="options">
            <label class="option" class:selected={paymentOption === 'cash'}>
              <input type="radio" name="payment" value="cash" bind:group={paymentOption} />
              <span class="material-symbols-outlined" aria-hidden="true">payments</span>
              <span class="option-text">
                <span class="title-md">Pay at counter</span>
                <span class="label-sm muted">Cash or card on pickup</span>
              </span>
            </label>
            <label class="option" class:selected={paymentOption === 'card'}>
              <input type="radio" name="payment" value="card" bind:group={paymentOption} />
              <span class="material-symbols-outlined" aria-hidden="true">credit_card</span>
              <span class="option-text">
                <span class="title-md">Card online</span>
                <span class="label-sm muted">Coming soon</span>
              </span>
            </label>
          </div>
        </section>

        {#if errorMessage}
          <p class="error body-md" role="alert">{errorMessage}</p>
        {/if}

        <button type="submit" class="primary-btn submit label-lg" disabled={!canSubmit}>
          {#if submitting}
            Placing order…
          {:else}
            Place order — ${$cartTotal.toFixed(2)}
          {/if}
        </button>
      </form>

      <!-- Summary -->
      <aside class="bb-card summary">
        <h2 class="title-lg">Order summary</h2>
        <ul class="items">
          {#each $cart as item (item.id)}
            <li>
              {#if item.image_url}
                <img src={item.image_url} alt="" loading="lazy" />
              {:else}
                <span class="thumb-placeholder material-symbols-outlined" aria-hidden="true">restaurant</span>
              {/if}
              <div class="item-info">
                <span class="title-md">{item.name}</span>
                <span class="label-sm muted">${item.price.toFixed(2)} each</span>
                <QuantitySelector qty={item.qty} onchange={(q) => setQty(item.id, q)} />
              </div>
              <span class="title-md line-total">${(item.price * item.qty).toFixed(2)}</span>
            </li>
          {/each}
        </ul>
        <div class="totals">
          <div class="row body-md">
            <span>Subtotal</span>
            <span>${$cartTotal.toFixed(2)}</span>
          </div>
          <div class="row body-md muted">
            <span>Taxes &amp; fees</span>
            <span>Included</span>
          </div>
          <div class="row total title-lg">
            <span>Total</span>
            <span>${$cartTotal.toFixed(2)}</span>
          </div>
        </div>
      </aside>
    </div>
  {/if}
</div>

<style>
  .page {
    width: 100%;
    max-width: 1100px;
    margin: 0 auto;
    padding: var(--bb-space-lg) var(--bb-margin-mobile) var(--bb-space-xl);
  }

  .page-head {
    margin-bottom: var(--bb-space-lg);
  }

  .page-head h1 {
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  .page-head p {
    margin: var(--bb-space-xs) 0 0;
  }

  .muted {
    color: var(--md-sys-color-on-surface-variant);
  }

  .layout {
    display: grid;
    grid-template-columns: 1fr;
    gap: var(--bb-space-lg);
    align-items: start;
  }

  @media (min-width: 900px) {
    .layout {
      grid-template-columns: 1fr 380px;
    }
  }

  .form {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-lg);
  }

  .block {
    padding: var(--bb-space-lg);
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }

  .block h2 {
    display: flex;
    align-items: center;
    gap: var(--bb-space-sm);
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  .step {
    width: 24px;
    height: 24px;
    border-radius: var(--bb-shape-full);
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-xs);
  }

  .field span {
    color: var(--md-sys-color-on-surface-variant);
  }

  .field input {
    border: 1px solid var(--md-sys-color-outline);
    border-radius: var(--bb-shape-sm);
    background: var(--md-sys-color-surface);
    color: var(--md-sys-color-on-surface);
    padding: 14px 16px;
    font-family: var(--md-ref-typeface-plain);
    font-size: 16px;
    outline: none;
    transition: border-color 150ms ease, box-shadow 150ms ease;
  }

  .field input:focus {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 1px var(--md-sys-color-primary);
  }

  .options {
    display: grid;
    grid-template-columns: 1fr;
    gap: var(--bb-space-md);
  }

  @media (min-width: 600px) {
    .options {
      grid-template-columns: 1fr 1fr;
    }
  }

  .option {
    display: flex;
    align-items: center;
    gap: var(--bb-space-md);
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: var(--bb-shape-md);
    padding: var(--bb-space-md);
    cursor: pointer;
    transition: border-color 150ms ease, background-color 150ms ease;
  }

  .option input {
    position: absolute;
    opacity: 0;
    pointer-events: none;
  }

  .option .material-symbols-outlined {
    color: var(--md-sys-color-primary);
    font-size: 28px;
  }

  .option.selected {
    border-color: var(--md-sys-color-primary);
    background: color-mix(in srgb, var(--md-sys-color-primary-container) 35%, transparent);
  }

  .option-text {
    display: flex;
    flex-direction: column;
  }

  .error {
    margin: 0;
    color: var(--md-sys-color-error);
  }

  .primary-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    padding: 16px 32px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
    text-decoration: none;
    box-shadow: var(--bb-elev-1);
    transition: box-shadow 200ms ease, opacity 150ms ease;
  }

  .primary-btn:hover:not(:disabled) {
    box-shadow: var(--bb-elev-2);
  }

  .primary-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .summary {
    padding: var(--bb-space-lg);
    position: sticky;
    top: 96px;
  }

  .summary h2 {
    margin: 0 0 var(--bb-space-md);
    color: var(--md-sys-color-on-surface);
  }

  .items {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }

  .items li {
    display: flex;
    gap: var(--bb-space-md);
    align-items: flex-start;
  }

  .items img,
  .thumb-placeholder {
    width: 56px;
    height: 56px;
    border-radius: var(--bb-shape-sm);
    object-fit: cover;
    flex-shrink: 0;
  }

  .thumb-placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-outline);
  }

  .item-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-xs);
    min-width: 0;
  }

  .line-total {
    color: var(--md-sys-color-primary);
    white-space: nowrap;
  }

  .totals {
    margin-top: var(--bb-space-lg);
    border-top: 1px dashed var(--md-sys-color-outline-variant);
    padding-top: var(--bb-space-md);
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-sm);
  }

  .row {
    display: flex;
    justify-content: space-between;
    color: var(--md-sys-color-on-surface);
  }

  .row.total {
    color: var(--md-sys-color-primary);
    font-weight: 800;
  }

  .empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--bb-space-md);
    padding: var(--bb-space-xl);
    text-align: center;
  }

  .empty .material-symbols-outlined {
    font-size: 48px;
    color: var(--md-sys-color-outline);
  }

  .empty p {
    margin: 0;
    color: var(--md-sys-color-on-surface-variant);
  }
</style>
