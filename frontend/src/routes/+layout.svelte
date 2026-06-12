<script lang="ts">
  import { browser } from '$app/environment';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { role } from '$lib/stores/auth';
  import AppHeader from '$lib/components/AppHeader.svelte';
  import AppFooter from '$lib/components/AppFooter.svelte';
  import AdminShell from '$lib/components/AdminShell.svelte';
  import CartDrawer from '$lib/components/CartDrawer.svelte';
  import ToastHost from '$lib/components/ToastHost.svelte';
  import '$lib/theme/tokens.css';
  import type { Snippet } from 'svelte';

  let { children }: { children: Snippet } = $props();

  let cartOpen = $state(false);

  const ADMIN_PREFIXES = ['/dashboard', '/admin', '/menu/manage', '/users'];

  const path = $derived($page.url.pathname);
  const isAdminRoute = $derived(
    ADMIN_PREFIXES.some((p) => path === p || path.startsWith(`${p}/`))
  );

  // Route protection: admin areas require the admin role.
  $effect(() => {
    if (!browser) return;
    if (isAdminRoute && $role !== 'admin') {
      goto($role === 'guest' ? '/login' : '/');
    }
  });

  onMount(() => {
    // Register @material/web custom elements on the client only.
    import('$lib/theme/material');
  });
</script>

{#if isAdminRoute && $role === 'admin'}
  <AdminShell>
    {@render children()}
  </AdminShell>
{:else if isAdminRoute}
  <!-- Redirecting unauthorized visitor; render nothing to avoid flashes. -->
{:else}
  <div class="app">
    <AppHeader oncartopen={() => (cartOpen = true)} />
    <main>
      {@render children()}
    </main>
    <AppFooter />
  </div>
  <CartDrawer bind:open={cartOpen} />
{/if}

<ToastHost />

<style>
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100dvh;
  }

  main {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
</style>
