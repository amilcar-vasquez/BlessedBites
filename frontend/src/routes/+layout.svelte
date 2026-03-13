<script lang="ts">
  import { browser } from '$app/environment';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { logout } from '$lib/api/auth';

  type Role = 'guest' | 'user' | 'admin';

  type NavItem = {
    href: string;
    label: string;
    access: 'public' | 'guest' | 'user' | 'admin';
  };

  const NAV_ITEMS: NavItem[] = [
    { href: '/', label: 'Home', access: 'public' },
    { href: '/checkout', label: 'Checkout', access: 'user' },
    { href: '/login', label: 'Login', access: 'guest' },
    { href: '/signup', label: 'Sign Up', access: 'guest' },
    { href: '/dashboard', label: 'Dashboard', access: 'admin' },
    { href: '/admin', label: 'Live Orders', access: 'admin' },
    { href: '/menu/manage', label: 'Manage Menu', access: 'admin' },
    { href: '/users', label: 'Users', access: 'admin' }
  ];

  const PROTECTED_ROUTE_RULES: Array<{ prefix: string; access: NavItem['access'] }> = [
    { prefix: '/dashboard', access: 'admin' },
    { prefix: '/admin', access: 'admin' },
    { prefix: '/menu/manage', access: 'admin' },
    { prefix: '/users', access: 'admin' },
    { prefix: '/checkout', access: 'user' }
  ];

  let role: Role = 'guest';
  let fullName = '';
  let authBusy = false;

  function syncAuthState() {
    if (!browser) return;

    const raw = localStorage.getItem('bb_user');
    if (!raw) {
      role = 'guest';
      fullName = '';
      return;
    }

    try {
      const user = JSON.parse(raw) as { role?: string; full_name?: string };
      const normalizedRole = (user.role || '').toLowerCase();
      role = normalizedRole === 'admin' ? 'admin' : 'user';
      fullName = user.full_name || '';
    } catch {
      role = 'guest';
      fullName = '';
    }
  }

  function canAccess(item: NavItem): boolean {
    if (item.access === 'public') return true;
    if (item.access === 'guest') return role === 'guest';
    if (item.access === 'user') return role === 'user' || role === 'admin';
    return role === 'admin';
  }

  function canAccessLevel(access: NavItem['access']): boolean {
    if (access === 'public') return true;
    if (access === 'guest') return role === 'guest';
    if (access === 'user') return role === 'user' || role === 'admin';
    return role === 'admin';
  }

  function requiredAccessForPath(pathname: string): NavItem['access'] {
    const match = PROTECTED_ROUTE_RULES.find((rule) =>
      pathname === rule.prefix || pathname.startsWith(`${rule.prefix}/`)
    );
    return match?.access || 'public';
  }

  async function enforceRouteProtection(pathname: string) {
    const required = requiredAccessForPath(pathname);
    if (canAccessLevel(required)) return;

    if (required === 'user') {
      await goto('/login');
      return;
    }

    await goto('/');
  }

  async function handleLogout() {
    authBusy = true;
    try {
      await logout();
    } catch (err) {
      console.warn('Logout request failed, clearing local auth anyway.', err);
    } finally {
      if (browser) {
        localStorage.removeItem('bb_access_token');
        localStorage.removeItem('bb_user');
      }
      syncAuthState();
      authBusy = false;
      await goto('/login');
    }
  }

  $: if (browser) {
    // Re-sync on client-side navigation.
    const pathname = $page.url.pathname;
    syncAuthState();
    enforceRouteProtection(pathname);
  }

  onMount(() => {
    syncAuthState();

    const handleStorage = () => syncAuthState();
    window.addEventListener('storage', handleStorage);

    // In local dev, clear stale SW/cache artifacts that can serve old bundles.
    if (window.location.hostname === 'localhost') {
      void (async () => {
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
      })();
    }

    return () => {
      window.removeEventListener('storage', handleStorage);
    };
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

<header class="topbar">
  <a class="brand" href="/">BlessedBites</a>

  <nav class="links" aria-label="Primary">
    {#each NAV_ITEMS as item (item.href)}
      {#if canAccess(item)}
        <a
          href={item.href}
          class:active={$page.url.pathname === item.href || (item.href !== '/' && $page.url.pathname.startsWith(item.href))}
        >
          {item.label}
        </a>
      {/if}
    {/each}
  </nav>

  <div class="session">
    {#if role !== 'guest'}
      <span>{fullName || 'Signed in'} ({role})</span>
      <button type="button" class="logout" on:click={handleLogout} disabled={authBusy}>
        {authBusy ? '...' : 'Logout'}
      </button>
    {:else}
      <span>Guest</span>
    {/if}
  </div>
</header>

<slot />

<style>
  .topbar {
    position: sticky;
    top: 0;
    z-index: 20;
    display: grid;
    grid-template-columns: auto 1fr auto;
    align-items: center;
    gap: 0.9rem;
    padding: 0.75rem 1rem;
    background: linear-gradient(120deg, #7f1d2d, #8d6500);
    box-shadow: 0 6px 20px rgba(35, 12, 8, 0.2);
  }

  .brand {
    color: #fff8e4;
    text-decoration: none;
    font-weight: 900;
    letter-spacing: 0.02em;
    font-size: 1.05rem;
  }

  .links {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .links a {
    text-decoration: none;
    color: #ffeec9;
    padding: 0.4rem 0.7rem;
    border-radius: 999px;
    font-weight: 700;
    border: 1px solid transparent;
  }

  .links a.active {
    background: #fff6df;
    color: #5e1a28;
    border-color: #f0dab5;
  }

  .session {
    display: inline-flex;
    align-items: center;
    gap: 0.6rem;
    color: #fff8e4;
    font-size: 0.85rem;
    font-weight: 600;
  }

  .logout {
    border: none;
    border-radius: 999px;
    background: #fff6df;
    color: #5e1a28;
    padding: 0.35rem 0.72rem;
    font-weight: 700;
    cursor: pointer;
  }

  .logout:disabled {
    opacity: 0.65;
    cursor: not-allowed;
  }

  @media (max-width: 860px) {
    .topbar {
      grid-template-columns: 1fr;
      gap: 0.55rem;
    }

    .session {
      justify-self: start;
    }
  }
</style>
