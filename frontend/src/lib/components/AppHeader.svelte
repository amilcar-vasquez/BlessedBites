<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { logout } from '$lib/api/auth';
  import { auth, role, clearSession } from '$lib/stores/auth';
  import { cartCount } from '$lib/stores/cart';
  import { darkMode, toggleTheme } from '$lib/stores/theme';

  let { oncartopen }: { oncartopen: () => void } = $props();

  let searchQuery = $state('');
  let userMenuOpen = $state(false);
  let mobileNavOpen = $state(false);

  const path = $derived($page.url.pathname);

  function submitSearch(e: Event) {
    e.preventDefault();
    const q = searchQuery.trim();
    goto(q ? `/menu?q=${encodeURIComponent(q)}` : '/menu');
  }

  async function handleLogout() {
    userMenuOpen = false;
    try {
      await logout();
    } catch {
      // clear local session regardless
    }
    clearSession();
    goto('/login');
  }

  function navTo(href: string) {
    userMenuOpen = false;
    mobileNavOpen = false;
    goto(href);
  }
</script>

<header class="bb-glass">
  <div class="left">
    <button
      type="button"
      class="mobile-menu-btn"
      aria-label="Open navigation"
      onclick={() => (mobileNavOpen = !mobileNavOpen)}
    >
      <span class="material-symbols-outlined">menu</span>
    </button>
    <a class="brand" href="/">Blessed Bites</a>
    <nav class="desktop-nav" aria-label="Main">
      <a class="label-lg" class:active={path === '/menu'} href="/menu">Menu</a>
      <a class="label-lg" href="/#popular">Offers</a>
      <a class="label-lg" href="/#about">About</a>
    </nav>
  </div>

  <form class="search" onsubmit={submitSearch} role="search">
    <span class="material-symbols-outlined search-icon" aria-hidden="true">search</span>
    <input
      type="search"
      placeholder="Search flavors..."
      aria-label="Search menu"
      bind:value={searchQuery}
    />
  </form>

  <div class="actions">
    <button
      type="button"
      class="icon-btn"
      aria-label={$darkMode ? 'Switch to light mode' : 'Switch to dark mode'}
      onclick={toggleTheme}
    >
      <span class="material-symbols-outlined">{$darkMode ? 'light_mode' : 'dark_mode'}</span>
    </button>

    {#if $role === 'guest'}
      <a class="login-btn label-lg" href="/login">Login</a>
      <a class="signup-btn label-lg" href="/signup">Sign Up</a>
    {:else}
      <div class="user-menu-wrap">
        <button
          type="button"
          class="user-btn label-lg"
          aria-haspopup="menu"
          aria-expanded={userMenuOpen}
          onclick={() => (userMenuOpen = !userMenuOpen)}
        >
          <span class="avatar label-sm">
            {($auth.user?.full_name || '?')
              .split(' ')
              .map((p) => p[0])
              .slice(0, 2)
              .join('')
              .toUpperCase()}
          </span>
          <span class="user-name">{$auth.user?.full_name || 'Account'}</span>
          <span class="material-symbols-outlined" aria-hidden="true">expand_more</span>
        </button>
        {#if userMenuOpen}
          <div class="menu bb-card" role="menu">
            {#if $role === 'admin'}
              <button type="button" class="menu-item label-lg" role="menuitem" onclick={() => navTo('/dashboard')}>
                <span class="material-symbols-outlined">dashboard</span> Dashboard
              </button>
              <button type="button" class="menu-item label-lg" role="menuitem" onclick={() => navTo('/admin')}>
                <span class="material-symbols-outlined">notifications_active</span> Live Orders
              </button>
              <button type="button" class="menu-item label-lg" role="menuitem" onclick={() => navTo('/menu/manage')}>
                <span class="material-symbols-outlined">restaurant_menu</span> Manage Menu
              </button>
              <button type="button" class="menu-item label-lg" role="menuitem" onclick={() => navTo('/users')}>
                <span class="material-symbols-outlined">group</span> Users
              </button>
            {/if}
            <button type="button" class="menu-item label-lg" role="menuitem" onclick={handleLogout}>
              <span class="material-symbols-outlined">logout</span> Logout
            </button>
          </div>
        {/if}
      </div>
    {/if}

    <button type="button" class="cart-btn label-lg" aria-label="Open cart" onclick={oncartopen}>
      <span class="material-symbols-outlined" aria-hidden="true">shopping_cart</span>
      <span class="cart-label">Cart</span>
      {#if $cartCount > 0}
        <span class="badge label-sm">{$cartCount}</span>
      {/if}
    </button>
  </div>
</header>

{#if mobileNavOpen}
  <nav class="mobile-nav bb-glass" aria-label="Mobile">
    <a class="label-lg" href="/menu" onclick={() => (mobileNavOpen = false)}>Menu</a>
    <a class="label-lg" href="/#popular" onclick={() => (mobileNavOpen = false)}>Offers</a>
    <a class="label-lg" href="/#about" onclick={() => (mobileNavOpen = false)}>About</a>
    {#if $role === 'guest'}
      <a class="label-lg" href="/login" onclick={() => (mobileNavOpen = false)}>Login</a>
      <a class="label-lg" href="/signup" onclick={() => (mobileNavOpen = false)}>Sign Up</a>
    {/if}
  </nav>
{/if}

<style>
  header {
    position: sticky;
    top: 0;
    z-index: 50;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--bb-space-md);
    height: 80px;
    padding: 0 var(--bb-margin-mobile);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    box-shadow: var(--bb-elev-1);
  }

  @media (min-width: 1024px) {
    header {
      padding: 0 var(--bb-margin-desktop);
    }
  }

  .left {
    display: flex;
    align-items: center;
    gap: var(--bb-space-lg);
    min-width: 0;
  }

  .brand {
    font-family: var(--md-ref-typeface-brand);
    font-size: 24px;
    line-height: 32px;
    font-weight: 700;
    color: var(--md-sys-color-primary);
    text-decoration: none;
    white-space: nowrap;
  }

  .desktop-nav {
    display: none;
    gap: var(--bb-space-md);
  }

  @media (min-width: 900px) {
    .desktop-nav {
      display: flex;
    }
  }

  .desktop-nav a {
    color: var(--md-sys-color-on-surface-variant);
    text-decoration: none;
    padding: 6px 8px;
    border-radius: 6px;
    transition: color 150ms ease, background-color 150ms ease;
  }

  .desktop-nav a:hover {
    color: var(--md-sys-color-primary);
    background: var(--md-sys-color-surface-container-high);
  }

  .desktop-nav a.active {
    color: var(--md-sys-color-primary);
    border-bottom: 2px solid var(--md-sys-color-primary);
    border-radius: 6px 6px 0 0;
  }

  .search {
    position: relative;
    flex: 1;
    max-width: 420px;
    display: none;
  }

  @media (min-width: 768px) {
    .search {
      display: block;
    }
  }

  .search-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--md-sys-color-tertiary);
    font-size: 20px;
    pointer-events: none;
  }

  .search input {
    width: 100%;
    border: none;
    border-radius: var(--bb-shape-full);
    background: var(--md-sys-color-surface-container-highest);
    color: var(--md-sys-color-on-surface);
    padding: 10px 16px 10px 40px;
    font-family: var(--md-ref-typeface-plain);
    font-size: 14px;
    outline: none;
    box-shadow: var(--bb-elev-1);
    transition: box-shadow 150ms ease, background-color 150ms ease;
  }

  .search input:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .search input:focus {
    box-shadow: 0 0 0 2px var(--md-sys-color-primary);
  }

  .actions {
    display: flex;
    align-items: center;
    gap: var(--bb-space-sm);
  }

  .icon-btn,
  .mobile-menu-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    color: var(--md-sys-color-primary);
    width: 40px;
    height: 40px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  .icon-btn:hover,
  .mobile-menu-btn:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .mobile-menu-btn {
    display: flex;
  }

  @media (min-width: 900px) {
    .mobile-menu-btn {
      display: none;
    }
  }

  .login-btn,
  .signup-btn {
    text-decoration: none;
    padding: 8px 16px;
    border-radius: var(--bb-shape-full);
    transition: background-color 150ms ease;
    white-space: nowrap;
  }

  .login-btn {
    color: var(--md-sys-color-primary);
  }

  .login-btn:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .signup-btn {
    display: none;
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
  }

  @media (min-width: 600px) {
    .signup-btn {
      display: inline-block;
    }
  }

  .user-menu-wrap {
    position: relative;
  }

  .user-btn {
    display: flex;
    align-items: center;
    gap: var(--bb-space-sm);
    border: none;
    background: transparent;
    color: var(--md-sys-color-on-surface);
    padding: 4px 8px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
  }

  .user-btn:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .avatar {
    width: 32px;
    height: 32px;
    border-radius: var(--bb-shape-full);
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
  }

  .user-name {
    display: none;
    max-width: 140px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  @media (min-width: 768px) {
    .user-name {
      display: inline;
    }
  }

  .menu {
    position: absolute;
    right: 0;
    top: calc(100% + 8px);
    min-width: 200px;
    padding: var(--bb-space-sm);
    display: flex;
    flex-direction: column;
    gap: 2px;
    box-shadow: var(--bb-elev-3);
    z-index: 60;
  }

  .menu-item {
    display: flex;
    align-items: center;
    gap: var(--bb-space-sm);
    border: none;
    background: transparent;
    color: var(--md-sys-color-on-surface);
    padding: 10px 12px;
    border-radius: var(--bb-shape-sm);
    cursor: pointer;
    text-align: left;
  }

  .menu-item:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .menu-item .material-symbols-outlined {
    font-size: 18px;
    color: var(--md-sys-color-on-surface-variant);
  }

  .cart-btn {
    position: relative;
    display: flex;
    align-items: center;
    gap: var(--bb-space-xs);
    border: none;
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    padding: 8px 16px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
    box-shadow: var(--bb-elev-1);
    transition: background-color 150ms ease;
  }

  .cart-btn:hover {
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-on-primary-container);
  }

  .cart-label {
    display: none;
  }

  @media (min-width: 600px) {
    .cart-label {
      display: inline;
    }
  }

  .badge {
    position: absolute;
    top: -4px;
    right: -4px;
    background: var(--md-sys-color-tertiary);
    color: var(--md-sys-color-on-tertiary);
    min-width: 18px;
    height: 18px;
    border-radius: var(--bb-shape-full);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 4px;
    font-weight: 700;
  }

  .mobile-nav {
    position: sticky;
    top: 80px;
    z-index: 49;
    display: flex;
    flex-direction: column;
    padding: var(--bb-space-sm) var(--bb-margin-mobile);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  @media (min-width: 900px) {
    .mobile-nav {
      display: none;
    }
  }

  .mobile-nav a {
    text-decoration: none;
    color: var(--md-sys-color-on-surface);
    padding: 12px 8px;
    border-radius: var(--bb-shape-sm);
  }

  .mobile-nav a:hover {
    background: var(--md-sys-color-surface-container-high);
  }
</style>
