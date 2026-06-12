<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { logout } from '$lib/api/auth';
  import { auth, clearSession } from '$lib/stores/auth';
  import { darkMode, toggleTheme } from '$lib/stores/theme';
  import type { Snippet } from 'svelte';

  let { children }: { children: Snippet } = $props();

  const NAV = [
    { href: '/dashboard', label: 'Dashboard', icon: 'dashboard' },
    { href: '/menu/manage', label: 'Menu & Categories', icon: 'restaurant_menu' },
    { href: '/admin', label: 'Live Orders', icon: 'shopping_bag' },
    { href: '/users', label: 'Users', icon: 'group' }
  ];

  const path = $derived($page.url.pathname);

  async function handleLogout() {
    try {
      await logout();
    } catch {
      // clear local session regardless
    }
    clearSession();
    goto('/login');
  }
</script>

<div class="shell">
  <nav class="rail" aria-label="Admin">
    <div class="rail-header">
      <h1 class="headline-md">Blessed Admin</h1>
      <p class="label-sm">Management Portal</p>
    </div>

    <div class="rail-nav">
      {#each NAV as item (item.href)}
        <a class="rail-link label-lg" class:active={path === item.href} href={item.href}>
          <span class="material-symbols-outlined" class:fill={path === item.href}>{item.icon}</span>
          {item.label}
        </a>
      {/each}
    </div>

    <div class="rail-footer">
      <a class="rail-link label-lg" href="/">
        <span class="material-symbols-outlined">storefront</span>
        View Store
      </a>
      <button type="button" class="rail-link label-lg" onclick={toggleTheme}>
        <span class="material-symbols-outlined">{$darkMode ? 'light_mode' : 'dark_mode'}</span>
        {$darkMode ? 'Light Mode' : 'Dark Mode'}
      </button>
      <button type="button" class="rail-link label-lg" onclick={handleLogout}>
        <span class="material-symbols-outlined">logout</span>
        Logout
      </button>
      <div class="user">
        <div class="avatar label-lg">
          {($auth.user?.full_name || '?')
            .split(' ')
            .map((p) => p[0])
            .slice(0, 2)
            .join('')
            .toUpperCase()}
        </div>
        <span class="title-md user-name">{$auth.user?.full_name || 'Admin'}</span>
      </div>
    </div>
  </nav>

  <main class="canvas">
    {@render children()}
  </main>
</div>

<nav class="mobile-bar" aria-label="Admin mobile">
  {#each NAV as item (item.href)}
    <a class="mobile-link label-sm" class:active={path === item.href} href={item.href}>
      <span class="material-symbols-outlined" class:fill={path === item.href}>{item.icon}</span>
      {item.label.split(' ')[0]}
    </a>
  {/each}
  <button type="button" class="mobile-link label-sm" onclick={handleLogout}>
    <span class="material-symbols-outlined">logout</span>
    Logout
  </button>
</nav>

<style>
  .shell {
    display: flex;
    min-height: 100dvh;
    background: var(--md-sys-color-background);
  }

  .rail {
    display: none;
    position: fixed;
    left: 0;
    top: 0;
    height: 100dvh;
    width: 240px;
    background: var(--md-sys-color-surface-container-low);
    border-right: 1px solid var(--md-sys-color-outline-variant);
    flex-direction: column;
    padding: var(--bb-space-lg) 0;
    z-index: 50;
  }

  @media (min-width: 900px) {
    .rail {
      display: flex;
    }
  }

  .rail-header {
    padding: 0 var(--bb-space-md);
    margin-bottom: var(--bb-space-xl);
  }

  .rail-header h1 {
    margin: 0;
    color: var(--md-sys-color-primary);
    font-weight: 700;
  }

  .rail-header p {
    margin: var(--bb-space-xs) 0 0;
    color: var(--md-sys-color-on-surface-variant);
  }

  .rail-nav {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-sm);
  }

  .rail-link {
    display: flex;
    align-items: center;
    gap: var(--bb-space-md);
    margin: 0 8px;
    padding: 12px 16px;
    border-radius: var(--bb-shape-sm);
    color: var(--md-sys-color-on-surface-variant);
    text-decoration: none;
    border: none;
    background: transparent;
    cursor: pointer;
    text-align: left;
    transition: background-color 150ms ease, transform 200ms ease;
  }

  .rail-link:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .rail-link.active {
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
    font-weight: 700;
    transform: translateX(4px);
  }

  .rail-footer {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-sm);
    border-top: 1px solid var(--md-sys-color-outline-variant);
    padding-top: var(--bb-space-md);
  }

  .user {
    display: flex;
    align-items: center;
    gap: var(--bb-space-md);
    margin: var(--bb-space-md) var(--bb-space-md) 0;
  }

  .avatar {
    width: 40px;
    height: 40px;
    border-radius: var(--bb-shape-full);
    background: var(--md-sys-color-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface-variant);
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    flex-shrink: 0;
  }

  .user-name {
    color: var(--md-sys-color-on-surface);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .canvas {
    flex: 1;
    min-width: 0;
    padding: var(--bb-margin-mobile);
    padding-bottom: 96px;
  }

  @media (min-width: 900px) {
    .canvas {
      margin-left: 240px;
      padding: var(--bb-margin-desktop);
    }
  }

  .mobile-bar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 50;
    display: flex;
    justify-content: space-around;
    background: var(--md-sys-color-surface-container-low);
    border-top: 1px solid var(--md-sys-color-outline-variant);
    padding: 6px 0 max(6px, env(safe-area-inset-bottom));
  }

  @media (min-width: 900px) {
    .mobile-bar {
      display: none;
    }
  }

  .mobile-link {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    color: var(--md-sys-color-on-surface-variant);
    text-decoration: none;
    border: none;
    background: transparent;
    cursor: pointer;
    padding: 6px 10px;
    border-radius: var(--bb-shape-sm);
  }

  .mobile-link.active {
    color: var(--md-sys-color-on-secondary-container);
    background: var(--md-sys-color-secondary-container);
  }
</style>
