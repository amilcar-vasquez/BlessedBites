<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { submitRating } from '$lib/api/ratings';
  import { auth } from '$lib/stores/auth';
  import { showToast } from '$lib/stores/toast';
  import RatingStars from '$lib/components/RatingStars.svelte';

  type LastOrder = {
    order_id: number;
    total: number;
    items: { id: number; name: string; qty: number }[];
  };

  let lastOrder = $state<LastOrder | null>(null);
  let rated = $state<Record<number, number>>({});

  const orderId = $derived($page.url.searchParams.get('id'));

  onMount(() => {
    try {
      const raw = sessionStorage.getItem('bb_last_order');
      if (raw) lastOrder = JSON.parse(raw) as LastOrder;
    } catch {
      lastOrder = null;
    }
  });

  async function rate(itemId: number, value: number) {
    try {
      await submitRating(itemId, value, $auth.user?.id);
      rated = { ...rated, [itemId]: value };
      showToast('Thanks for the feedback!', 'success', 2500);
    } catch {
      showToast('Could not save your rating. Please try again.', 'error');
    }
  }
</script>

<svelte:head>
  <title>Order confirmed — Blessed Bites</title>
</svelte:head>

<div class="page">
  <div class="bb-card hero">
    <span class="check material-symbols-outlined fill" aria-hidden="true">check_circle</span>
    <h1 class="headline-lg">Order confirmed!</h1>
    <p class="body-lg">
      {#if orderId}
        Your order <strong>#{orderId}</strong> is in the kitchen.
      {:else}
        Your order is in the kitchen.
      {/if}
      We'll have it ready shortly.
    </p>

    <div class="timeline" aria-hidden="true">
      <div class="step done">
        <span class="dot"></span>
        <span class="label-sm">Received</span>
      </div>
      <div class="bar"></div>
      <div class="step active">
        <span class="dot"></span>
        <span class="label-sm">Preparing</span>
      </div>
      <div class="bar dim"></div>
      <div class="step">
        <span class="dot"></span>
        <span class="label-sm">Ready</span>
      </div>
    </div>

    <div class="actions">
      <a class="primary-btn label-lg" href="/menu">Order more</a>
      <a class="text-btn label-lg" href="/">Back to home</a>
    </div>
  </div>

  {#if lastOrder && lastOrder.items.length > 0}
    <section class="bb-card ratings">
      <h2 class="title-lg">How was everything?</h2>
      <p class="body-md muted">Rate your dishes to help the kitchen keep improving.</p>
      <ul>
        {#each lastOrder.items as item (item.id)}
          <li>
            <span class="title-md">{item.name}</span>
            {#if rated[item.id]}
              <span class="thanks label-lg">
                <span class="material-symbols-outlined fill" aria-hidden="true">favorite</span>
                Rated {rated[item.id]}/5
              </span>
            {:else}
              <RatingStars value={0} interactive onrate={(v) => rate(item.id, v)} />
            {/if}
          </li>
        {/each}
      </ul>
    </section>
  {/if}
</div>

<style>
  .page {
    width: 100%;
    max-width: 640px;
    margin: 0 auto;
    padding: var(--bb-space-xl) var(--bb-margin-mobile);
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-lg);
  }

  .hero {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: var(--bb-space-xl);
    gap: var(--bb-space-md);
  }

  .check {
    font-size: 72px;
    color: var(--md-sys-color-secondary);
  }

  .hero h1 {
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  .hero p {
    margin: 0;
    color: var(--md-sys-color-on-surface-variant);
  }

  .timeline {
    display: flex;
    align-items: center;
    gap: var(--bb-space-sm);
    margin: var(--bb-space-md) 0;
    width: 100%;
    max-width: 400px;
  }

  .step {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--bb-space-xs);
    color: var(--md-sys-color-outline);
  }

  .step .dot {
    width: 14px;
    height: 14px;
    border-radius: var(--bb-shape-full);
    background: var(--md-sys-color-outline-variant);
  }

  .step.done {
    color: var(--md-sys-color-secondary);
  }

  .step.done .dot {
    background: var(--md-sys-color-secondary);
  }

  .step.active {
    color: var(--md-sys-color-primary);
  }

  .step.active .dot {
    background: var(--md-sys-color-primary);
    animation: bb-pulse 1.5s ease-in-out infinite;
  }

  .bar {
    flex: 1;
    height: 3px;
    border-radius: 2px;
    background: var(--md-sys-color-secondary);
    margin-bottom: 18px;
  }

  .bar.dim {
    background: var(--md-sys-color-outline-variant);
  }

  .actions {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: var(--bb-space-md);
  }

  .primary-btn {
    display: inline-flex;
    align-items: center;
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    padding: 12px 28px;
    border-radius: var(--bb-shape-full);
    text-decoration: none;
    box-shadow: var(--bb-elev-1);
  }

  .text-btn {
    display: inline-flex;
    align-items: center;
    color: var(--md-sys-color-primary);
    padding: 12px 20px;
    border-radius: var(--bb-shape-full);
    text-decoration: none;
  }

  .text-btn:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .ratings {
    padding: var(--bb-space-lg);
  }

  .ratings h2 {
    margin: 0 0 var(--bb-space-xs);
    color: var(--md-sys-color-on-surface);
  }

  .ratings p {
    margin: 0 0 var(--bb-space-md);
  }

  .muted {
    color: var(--md-sys-color-on-surface-variant);
  }

  .ratings ul {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }

  .ratings li {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--bb-space-md);
    flex-wrap: wrap;
  }

  .thanks {
    display: inline-flex;
    align-items: center;
    gap: var(--bb-space-xs);
    color: var(--md-sys-color-tertiary);
  }
</style>
