<script lang="ts">
  import { onMount } from 'svelte';
  import { listAdminMenu, listAdminOrders, type AdminOrder } from '$lib/api/admin';
  import type { MenuItem } from '$lib/api/menu';
  import Skeleton from '$lib/components/Skeleton.svelte';
  import StatusChip from '$lib/components/StatusChip.svelte';

  let orders = $state<AdminOrder[]>([]);
  let menu = $state<MenuItem[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(async () => {
    try {
      const [o, m] = await Promise.all([listAdminOrders(), listAdminMenu()]);
      orders = o;
      menu = m;
    } catch {
      error = 'Could not load dashboard data.';
    } finally {
      loading = false;
    }
  });

  const totalRevenue = $derived(orders.reduce((sum, o) => sum + o.total_cost, 0));

  const ordersToday = $derived.by(() => {
    const today = new Date().toDateString();
    return orders.filter((o) => new Date(o.created_at).toDateString() === today).length;
  });

  const activeItems = $derived(menu.filter((m) => m.is_active !== false).length);

  const recentOrders = $derived(
    [...orders]
      .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
      .slice(0, 8)
  );

  // Revenue per day for the last 7 days (CSS bar chart).
  const weekly = $derived.by(() => {
    const days: { label: string; total: number }[] = [];
    for (let i = 6; i >= 0; i--) {
      const d = new Date();
      d.setDate(d.getDate() - i);
      const key = d.toDateString();
      const total = orders
        .filter((o) => new Date(o.created_at).toDateString() === key)
        .reduce((sum, o) => sum + o.total_cost, 0);
      days.push({ label: d.toLocaleDateString(undefined, { weekday: 'short' }), total });
    }
    const max = Math.max(...days.map((d) => d.total), 1);
    return days.map((d) => ({ ...d, pct: Math.round((d.total / max) * 100) }));
  });

  function formatMoney(n: number): string {
    return `$${n.toFixed(2)}`;
  }

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
  <title>Dashboard — Blessed Admin</title>
</svelte:head>

<header class="head">
  <h1 class="headline-lg">Dashboard</h1>
  <p class="body-md muted">An overview of how the kitchen is doing.</p>
</header>

{#if error}
  <p class="bb-form-error" role="alert">{error}</p>
{/if}

<!-- Metric cards -->
<div class="metrics">
  {#if loading}
    {#each Array(4) as _, i (i)}
      <Skeleton width="100%" height="110px" radius="16px" />
    {/each}
  {:else}
    <div class="bb-card metric">
      <span class="metric-icon material-symbols-outlined" aria-hidden="true">payments</span>
      <span class="label-lg muted">Total Revenue</span>
      <span class="headline-md">{formatMoney(totalRevenue)}</span>
    </div>
    <div class="bb-card metric">
      <span class="metric-icon material-symbols-outlined" aria-hidden="true">today</span>
      <span class="label-lg muted">Orders Today</span>
      <span class="headline-md">{ordersToday}</span>
    </div>
    <div class="bb-card metric">
      <span class="metric-icon material-symbols-outlined" aria-hidden="true">receipt_long</span>
      <span class="label-lg muted">Total Orders</span>
      <span class="headline-md">{orders.length}</span>
    </div>
    <div class="bb-card metric">
      <span class="metric-icon material-symbols-outlined" aria-hidden="true">restaurant_menu</span>
      <span class="label-lg muted">Active Dishes</span>
      <span class="headline-md">{activeItems}</span>
    </div>
  {/if}
</div>

<div class="panels">
  <!-- Recent orders -->
  <section class="bb-card panel">
    <div class="panel-head">
      <h2 class="title-lg">Recent Orders</h2>
      <a class="label-lg link" href="/admin">Live view</a>
    </div>
    {#if loading}
      <Skeleton width="100%" height="240px" radius="8px" />
    {:else if recentOrders.length === 0}
      <p class="body-md muted">No orders yet.</p>
    {:else}
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
            {#each recentOrders as order (order.id)}
              <tr>
                <td class="body-md mono">#{order.id}</td>
                <td class="body-md">{order.full_name || `User ${order.user_id}`}</td>
                <td class="body-md muted">{formatTime(order.created_at)}</td>
                <td><StatusChip status={order.status} /></td>
                <td class="body-md num strong">{formatMoney(order.total_cost)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </section>

  <!-- Weekly revenue chart -->
  <section class="bb-card panel">
    <div class="panel-head">
      <h2 class="title-lg">Weekly Revenue</h2>
    </div>
    {#if loading}
      <Skeleton width="100%" height="240px" radius="8px" />
    {:else}
      <div class="chart" role="img" aria-label="Bar chart of revenue for the last seven days">
        {#each weekly as day (day.label + day.total)}
          <div class="col">
            <span class="label-sm amount">{day.total > 0 ? formatMoney(day.total) : ''}</span>
            <div class="bar" style={`height: ${Math.max(day.pct, 4)}%`}></div>
            <span class="label-sm muted">{day.label}</span>
          </div>
        {/each}
      </div>
    {/if}
  </section>
</div>

<style>
  .head {
    margin-bottom: var(--bb-space-lg);
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

  .metrics {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: var(--bb-space-md);
    margin-bottom: var(--bb-space-lg);
  }

  .metric {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-xs);
    padding: var(--bb-space-lg);
    position: relative;
  }

  .metric .headline-md {
    color: var(--md-sys-color-on-surface);
    font-weight: 800;
  }

  .metric-icon {
    position: absolute;
    top: var(--bb-space-md);
    right: var(--bb-space-md);
    color: var(--md-sys-color-primary);
    background: color-mix(in srgb, var(--md-sys-color-primary-container) 50%, transparent);
    border-radius: var(--bb-shape-full);
    padding: 8px;
  }

  .panels {
    display: grid;
    grid-template-columns: 1fr;
    gap: var(--bb-space-lg);
  }

  @media (min-width: 1200px) {
    .panels {
      grid-template-columns: 3fr 2fr;
    }
  }

  .panel {
    padding: var(--bb-space-lg);
  }

  .panel-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--bb-space-md);
  }

  .panel-head h2 {
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  .link {
    color: var(--md-sys-color-primary);
    text-decoration: none;
  }

  .link:hover {
    text-decoration: underline;
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

  .num {
    text-align: right;
  }

  .strong {
    font-weight: 700;
  }

  .mono {
    font-variant-numeric: tabular-nums;
  }

  .chart {
    display: flex;
    align-items: stretch;
    gap: var(--bb-space-sm);
    height: 240px;
  }

  .col {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-end;
    gap: var(--bb-space-xs);
    min-width: 0;
  }

  .amount {
    color: var(--md-sys-color-on-surface-variant);
    white-space: nowrap;
    font-size: 10px;
  }

  .bar {
    width: 70%;
    max-width: 48px;
    border-radius: 8px 8px 2px 2px;
    background: linear-gradient(to top, var(--md-sys-color-primary), color-mix(in srgb, var(--md-sys-color-primary) 65%, var(--md-sys-color-tertiary)));
    transition: height 400ms ease;
  }
</style>
