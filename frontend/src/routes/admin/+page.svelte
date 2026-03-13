<script lang="ts">
  import OrderTracker from '$lib/components/OrderTracker.svelte';
</script>

<main>
  <h1>Kitchen Dashboard</h1>
  <p>Live stream for incoming orders.</p>
  <nav class="links">
    <a href="/dashboard">Dashboard</a>
    <a href="/menu/manage">Manage Menu</a>
    <a href="/users">Users</a>
  </nav>
  <OrderTracker />
</main>

<style>
  main { max-width: 980px; margin: 1.2rem auto; padding: 1rem; }
  p { color: #6f4744; }
  .links { margin: 0.8rem 0 1rem; display: flex; gap: 0.7rem; flex-wrap: wrap; }
  .links a { text-decoration: none; background: #f6dfb7; color: #492316; padding: 0.45rem 0.75rem; border-radius: 999px; font-weight: 700; }
</style>
