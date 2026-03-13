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
  let submitting = false;

  onMount(() => {
    const raw = localStorage.getItem('bb_cart');
    cart = raw ? (JSON.parse(raw) as CartItem[]) : [];
    total = cart.reduce((sum, item) => sum + item.price * item.qty, 0);
  });

  async function submitOrder() {
    error = '';
    if (cart.length === 0) {
      error = 'Your cart is empty.';
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
      window.location.href = '/order-success';
    } catch (e) {
      console.error(e);
      error = 'Could not place order. Check your details and try again.';
    } finally {
      submitting = false;
    }
  }
</script>

<main>
  <h1>Checkout</h1>
  <p>Guest checkout is supported. This page will submit to `POST /api/v1/orders`.</p>
  {#if cart.length === 0}
    <p>Your cart is empty.</p>
  {:else}
    <ul>
      {#each cart as item}
        <li>{item.name} x {item.qty} - ${(item.price * item.qty).toFixed(2)}</li>
      {/each}
    </ul>
    <strong>Total: ${total.toFixed(2)}</strong>

    <div>
      <label>
        Full name
        <input bind:value={fullName} placeholder="Your full name" />
      </label>
      <label>
        Phone
        <input bind:value={phoneNo} placeholder="Your phone number" />
      </label>
      <button type="button" disabled={submitting} on:click={submitOrder}>
        {submitting ? 'Submitting...' : 'Place order'}
      </button>
      {#if error}
        <p>{error}</p>
      {/if}
    </div>
  {/if}
</main>
