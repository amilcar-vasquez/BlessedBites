<script lang="ts">
  import { onMount } from 'svelte';

  onMount(async () => {
    // In local dev, clear stale SW/cache artifacts that can serve old bundles.
    if (window.location.hostname !== 'localhost') return;

    try {
      if ('serviceWorker' in navigator) {
        const registrations = await navigator.serviceWorker.getRegistrations();
        await Promise.all(registrations.map((registration) => registration.unregister()));
      }

      if ('caches' in window) {
        const keys = await caches.keys();
        await Promise.all(keys.map((key) => caches.delete(key)));
      }
    } catch (err) {
      console.warn('Local cache purge skipped:', err);
    }
  });
</script>

<svelte:head>
  <title>BlessedBites | Food in San Ignacio</title>
  <meta name="description" content="BlessedBites restaurant ordering platform in San Ignacio, Belize." />
  <meta name="keywords" content="food in San Ignacio,best food in Belize,restaurant menu Belize" />
  <script type="application/ld+json">
    {JSON.stringify({
      '@context': 'https://schema.org',
      '@type': 'Restaurant',
      name: 'BlessedBites',
      areaServed: 'San Ignacio, Belize',
      servesCuisine: 'Belizean'
    })}
  </script>
</svelte:head>

<slot />
