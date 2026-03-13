<script lang="ts">
  import { onMount } from 'svelte';
  import { apiPost } from '$lib/api/client';

  type CartItem = { id: number; name: string; price: number; qty: number };
  type OrderResponse = { order_id: number; message: string };

  let cart: CartItem[] = [];
  let total = 0;
  let fullName = '';
  let phoneNo = '';
  let error = '';
  let successHint = '';
  let submitting = false;

  onMount(() => {
    const raw = localStorage.getItem('bb_cart');
    cart = raw ? (JSON.parse(raw) as CartItem[]) : [];
    total = cart.reduce((sum, item) => sum + item.price * item.qty, 0);

    const userRaw = localStorage.getItem('bb_user');
    if (userRaw) {
      try {
        const user = JSON.parse(userRaw) as { full_name?: string; phone_no?: string };
        fullName = user.full_name || '';
        phoneNo = user.phone_no || '';
      } catch {
        // ignore malformed local user cache
      }
    }
  });

  async function submitOrder() {
    error = '';
    successHint = '';
    if (cart.length === 0) {
      error = 'Your cart is empty.';
      return;
    }

    if (!fullName.trim() || !phoneNo.trim()) {
      error = 'Please provide both full name and phone number.';
      return;
    }

    submitting = true;
    try {
      const items = cart.map((item) => ({ id: item.id, qty: item.qty }));
      await apiPost<OrderResponse>('/orders', {
        items,
        full_name: fullName,
        phone_no: phoneNo
      });

      localStorage.removeItem('bb_cart');
      successHint = 'Order placed. Redirecting...';
      window.location.href = '/order-success';
    } catch (e) {
      console.error(e);
      error = 'Could not place order. Check your details and try again.';
    } finally {
      submitting = false;
    }
  }
</script>

<main class="checkout-shell">
  <section class="hero card">
    <div>
      <p class="eyebrow">Final Step</p>
      <h1>Secure Checkout</h1>
      <p class="support">Fast pickup flow designed for San Ignacio rush hour.</p>
    </div>
    <a class="back-link" href="/">Back To Menu</a>
  </section>

  {#if cart.length === 0}
    <section class="empty card">
      <h2>Your cart is empty</h2>
      <p>Add a few favorites before checking out.</p>
      <a href="/" class="cta">Browse Menu</a>
    </section>
  {:else}
    <section class="grid">
      <article class="card summary" aria-label="Order summary">
        <h2>Order Summary</h2>
        <ul>
          {#each cart as item (item.id)}
            <li>
              <div class="item-meta">
                <span class="item-name">{item.name}</span>
                <span class="item-qty">Qty {item.qty}</span>
              </div>
              <strong>${(item.price * item.qty).toFixed(2)}</strong>
            </li>
          {/each}
        </ul>

        <div class="total-row">
          <span>Total</span>
          <strong>${total.toFixed(2)}</strong>
        </div>
      </article>

      <article class="card form-card" aria-label="Customer information">
        <h2>Pickup Details</h2>
        <p class="support">We will use this to confirm your order is ready.</p>

        <form on:submit|preventDefault={submitOrder} class="form">
          <label for="fullName">Full name</label>
          <input id="fullName" bind:value={fullName} placeholder="e.g. Ana Lopez" autocomplete="name" required />

          <label for="phone">Phone number</label>
          <input id="phone" bind:value={phoneNo} placeholder="e.g. +501 610 1234" autocomplete="tel" required />

          <button class="cta" type="submit" disabled={submitting}>
            {submitting ? 'Placing Order...' : 'Place Order'}
          </button>
        </form>

        {#if error}
          <p class="feedback error" role="alert">{error}</p>
        {/if}
        {#if successHint}
          <p class="feedback ok" role="status">{successHint}</p>
        {/if}
      </article>
    </section>
  {/if}
</main>

<style>
  :global(body) {
    background:
      radial-gradient(circle at 8% 8%, #f9e9c8 0%, transparent 42%),
      radial-gradient(circle at 90% 12%, #f1d6bf 0%, transparent 36%),
      #fff7ef;
  }

  .checkout-shell {
    --m3-primary: #7f1d2d;
    --m3-secondary: #8d6500;
    --m3-surface: #fffaf5;
    --m3-outline: #e6d2c7;
    --m3-on-surface-variant: #6f4744;

    max-width: 1120px;
    margin: 1.25rem auto 2rem;
    padding: 0 1rem;
    font-family: 'Sora', 'Nunito Sans', 'Trebuchet MS', sans-serif;
  }

  .card {
    background: var(--m3-surface);
    border: 1px solid var(--m3-outline);
    border-radius: 24px;
    box-shadow: 0 8px 30px rgba(61, 20, 16, 0.08);
  }

  .hero {
    padding: 1.2rem 1.2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .eyebrow {
    margin: 0;
    color: var(--m3-secondary);
    text-transform: uppercase;
    letter-spacing: 0.08em;
    font-size: 0.75rem;
    font-weight: 800;
  }

  h1 {
    margin: 0.2rem 0 0.4rem;
    color: #4f1a24;
    font-size: clamp(1.4rem, 2.4vw, 2rem);
  }

  h2 {
    margin: 0 0 0.5rem;
    color: #4f1a24;
  }

  .support {
    margin: 0;
    color: var(--m3-on-surface-variant);
  }

  .back-link {
    text-decoration: none;
    font-weight: 700;
    color: #5e1a28;
    background: #f6dfb7;
    border-radius: 999px;
    padding: 0.55rem 0.9rem;
    white-space: nowrap;
  }

  .grid {
    display: grid;
    gap: 1rem;
    grid-template-columns: 1fr 1fr;
  }

  .summary,
  .form-card,
  .empty {
    padding: 1rem;
  }

  ul {
    list-style: none;
    margin: 0.75rem 0;
    padding: 0;
    display: grid;
    gap: 0.65rem;
  }

  li {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    border-bottom: 1px solid #f0dfd6;
    padding-bottom: 0.55rem;
  }

  .item-meta {
    display: grid;
    gap: 0.1rem;
  }

  .item-name {
    font-weight: 700;
    color: #4f1a24;
  }

  .item-qty {
    font-size: 0.82rem;
    color: var(--m3-on-surface-variant);
  }

  .total-row {
    margin-top: 0.8rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 1.05rem;
  }

  .total-row strong {
    color: var(--m3-primary);
    font-size: 1.3rem;
  }

  .form {
    margin-top: 0.7rem;
    display: grid;
    gap: 0.45rem;
  }

  label {
    font-size: 0.88rem;
    font-weight: 700;
    color: #57322e;
  }

  input {
    border: 1px solid #d8c4b8;
    border-radius: 14px;
    padding: 0.72rem 0.84rem;
    font: inherit;
    transition: border-color 150ms ease, box-shadow 150ms ease;
  }

  input:focus {
    outline: none;
    border-color: #8d6500;
    box-shadow: 0 0 0 3px rgba(141, 101, 0, 0.15);
  }

  .cta {
    margin-top: 0.55rem;
    border: none;
    border-radius: 999px;
    padding: 0.72rem 1rem;
    background: linear-gradient(130deg, var(--m3-primary), #a2253b);
    color: #fff9ea;
    font-weight: 800;
    letter-spacing: 0.01em;
    text-decoration: none;
    text-align: center;
    cursor: pointer;
  }

  .cta:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .feedback {
    margin-top: 0.7rem;
    border-radius: 12px;
    padding: 0.65rem 0.75rem;
    font-weight: 600;
  }

  .feedback.error {
    background: #ffe5e8;
    color: #7d1228;
  }

  .feedback.ok {
    background: #e2f5e8;
    color: #1f6b3c;
  }

  .empty p {
    color: var(--m3-on-surface-variant);
    margin-bottom: 1rem;
  }

  @media (max-width: 860px) {
    .hero {
      flex-direction: column;
      align-items: flex-start;
    }

    .grid {
      grid-template-columns: 1fr;
    }
  }
</style>
