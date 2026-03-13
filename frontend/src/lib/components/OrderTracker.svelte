<script lang="ts">
  import { onMount } from 'svelte';

  type OrderEvent = {
    order_id: number;
    user_id: number;
    total_cost: number;
  };

  const base = import.meta.env.PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';
  const streamURL = `${base}/orders/stream`;

  let connected = false;
  let events: OrderEvent[] = [];

  onMount(() => {
    const source = new EventSource(streamURL);

    source.addEventListener('open', () => {
      connected = true;
    });

    source.addEventListener('error', () => {
      connected = false;
    });

    source.addEventListener('order.created', (evt) => {
      const payload = JSON.parse((evt as MessageEvent).data) as OrderEvent;
      events = [payload, ...events].slice(0, 20);
    });

    return () => source.close();
  });
</script>

<section>
  <p>Status: {connected ? 'Connected' : 'Disconnected'}</p>
  <h2>Incoming Orders</h2>
  {#if events.length === 0}
    <p>No live orders yet.</p>
  {:else}
    <ul>
      {#each events as event}
        <li>Order #{event.order_id} | User {event.user_id} | ${event.total_cost.toFixed(2)}</li>
      {/each}
    </ul>
  {/if}
</section>
