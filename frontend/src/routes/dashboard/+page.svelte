<script lang="ts">
  import { onMount } from 'svelte';
  import { listAdminOrders, type AdminOrder } from '$lib/api/admin';

  let orders: AdminOrder[] = [];
  let loading = true;
  let error = '';

  $: totalOrders = orders.length;
  $: todaySales = orders
    .filter((order) => {
      const created = new Date(order.created_at);
      const now = new Date();
      return (
        created.getFullYear() === now.getFullYear() &&
        created.getMonth() === now.getMonth() &&
        created.getDate() === now.getDate()
      );
    })
    .map((order) => ({ full_name: order.full_name, total_cost: order.total_cost }));

  onMount(async () => {
    try {
      orders = await listAdminOrders();
    } catch (e) {
      console.error(e);
      error = 'Could not load dashboard data. Ensure you are logged in as admin.';
    } finally {
      loading = false;
    }
  });
</script>

<main class="shell">
  <h1>Dashboard</h1>
  <p class="muted">Manage store performance and operations.</p>

  <nav class="links">
    <a href="/menu/manage">Menu Management</a>
    <a href="/users">Users</a>
    <a href="/admin">Live Orders</a>
  </nav>

  {#if loading}
    <p>Loading dashboard...</p>
  {:else if error}
    <p class="error">{error}</p>
  {:else}
    <section class="cards">
      <article class="card">
        <h2>Total Orders</h2>
        <p class="metric">{totalOrders}</p>
      </article>
      <article class="card">
        <h2>Sales for Today</h2>
        {#if todaySales.length === 0}
          <p class="muted">No sales yet today.</p>
        {:else}
          <ul>
            {#each todaySales as sale}
              <li><span>{sale.full_name}</span><strong>${sale.total_cost.toFixed(2)}</strong></li>
            {/each}
          </ul>
        {/if}
      </article>
    </section>
  {/if}
</main>

<style>
  .shell { max-width: 1080px; margin: 1.2rem auto; padding: 1rem; }
  h1 { margin: 0; }
  .muted { color: #6f4744; }
  .links { margin: 0.8rem 0 1rem; display: flex; gap: 0.7rem; flex-wrap: wrap; }
  .links a { text-decoration: none; background: #f6dfb7; color: #492316; padding: 0.45rem 0.75rem; border-radius: 999px; font-weight: 700; }
  .cards { display: grid; gap: 1rem; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); }
  .card { background: #fffaf5; border: 1px solid #e7d3c9; border-radius: 20px; padding: 1rem; }
  .metric { font-size: 2rem; margin: 0.35rem 0; color: #7f1d2d; font-weight: 800; }
  ul { list-style: none; margin: 0; padding: 0; display: grid; gap: 0.5rem; }
  li { display: flex; justify-content: space-between; gap: 1rem; }
  .error { color: #8a1732; }
</style>
