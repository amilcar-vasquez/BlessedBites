<script lang="ts">
  import { onMount } from 'svelte';
  import { listAdminOrders, type AdminOrder } from '$lib/api/admin';
  import { subscribeOrdersStream } from '$lib/api/orders';
  import { showToast } from '$lib/stores/toast';
  import Skeleton from '$lib/components/Skeleton.svelte';
  import StatusChip from '$lib/components/StatusChip.svelte';

  let orders = $state<AdminOrder[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let connected = $state(false);
  let freshIds = $state<Set<number>>(new Set());

  async function refresh() {
    try {
      const list = await listAdminOrders();
      orders = list.sort(
        (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      );
      error = null;
    } catch {
      error = 'Could not load orders.';
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    refresh();
    const unsubscribe = subscribeOrdersStream({
      onOpen: () => (connected = true),
      onError: () => (connected = false),
      onOrderCreated: (evt) => {
        showToast(`New order #${evt.order_id} — $${evt.total_cost.toFixed(2)}`, 'info');
        freshIds = new Set([...freshIds, evt.order_id]);
        refresh();
      }
    });
    return unsubscribe;
  });

  function formatTime(iso: string): string {
    return new Date(iso).toLocaleString(undefined, {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit'
    });
  }
</script>

<svelte:head>
  <title>Live Orders — Blessed Admin</title>
</svelte:head>

<header class="head">
  <div>
    <h1 class="headline-lg">Live Orders</h1>
    <p class="body-md muted">New orders appear here in real time.</p>
  </div>
  <span class="conn label-lg" class:on={connected}>
    <span class="dot"></span>
    {connected ? 'Live' : 'Reconnecting…'}
  </span>
</header>

{#if error}
  <p class="bb-form-error" role="alert">{error}</p>
{/if}

{#if loading}
  <Skeleton width="100%" height="320px" radius="16px" />
{:else if orders.length === 0}
  <div class="bb-card empty">
    <span class="material-symbols-outlined" aria-hidden="true">notifications</span>
    <p class="body-lg">No orders yet. They'll show up here the moment they're placed.</p>
  </div>
{:else}
  <div class="bb-card table-card">
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th class="label-sm">Order</th>
            <th class="label-sm">Customer</th>
            <th class="label-sm">Placed</th>
            <th class="label-sm">Status</th>
            <th class="label-sm num">Total</th>
          </tr>
        </thead>
        <tbody>
          {#each orders as order (order.id)}
            <tr class:fresh={freshIds.has(order.id)}>
              <td class="body-md mono">#{order.id}</td>
              <td class="body-md">{order.full_name || `User ${order.user_id}`}</td>
              <td class="body-md muted">{formatTime(order.created_at)}</td>
              <td><StatusChip status={order.status} /></td>
              <td class="body-md num strong">${order.total_cost.toFixed(2)}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
{/if}

<style>
  .head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--bb-space-md);
    margin-bottom: var(--bb-space-lg);
    flex-wrap: wrap;
  }

  .head h1 {
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  .head p {
    margin: var(--bb-space-xs) 0 0;
  }

  .muted {
    color: var(--md-sys-color-on-surface-variant);
  }

  .conn {
    display: inline-flex;
    align-items: center;
    gap: var(--bb-space-sm);
    padding: 8px 16px;
    border-radius: var(--bb-shape-full);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface-variant);
  }

  .conn .dot {
    width: 10px;
    height: 10px;
    border-radius: var(--bb-shape-full);
    background: var(--md-sys-color-outline);
  }

  .conn.on {
    border-color: transparent;
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
  }

  .conn.on .dot {
    background: var(--md-sys-color-secondary);
    animation: bb-pulse-dot 1.5s ease-in-out infinite;
  }

  @keyframes bb-pulse-dot {
    0%,
    100% {
      opacity: 1;
    }
    50% {
      opacity: 0.4;
    }
  }

  .empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--bb-space-md);
    padding: var(--bb-space-xl);
    text-align: center;
    color: var(--md-sys-color-on-surface-variant);
  }

  .empty .material-symbols-outlined {
    font-size: 48px;
    color: var(--md-sys-color-outline);
  }

  .empty p {
    margin: 0;
  }

  .table-card {
    padding: var(--bb-space-md);
  }

  .table-wrap {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th {
    text-align: left;
    color: var(--md-sys-color-on-surface-variant);
    text-transform: uppercase;
    padding: 10px 12px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  td {
    padding: 12px;
    color: var(--md-sys-color-on-surface);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  tbody tr:nth-child(even) {
    background: var(--md-sys-color-surface-container-low);
  }

  tbody tr:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  tr.fresh td {
    animation: bb-flash 2s ease-out 1;
  }

  @keyframes bb-flash {
    from {
      background-color: color-mix(in srgb, var(--md-sys-color-tertiary-container) 70%, transparent);
    }
    to {
      background-color: transparent;
    }
  }

  .num {
    text-align: right;
  }

  .strong {
    font-weight: 700;
  }

  .mono {
    font-variant-numeric: tabular-nums;
  }
</style>
