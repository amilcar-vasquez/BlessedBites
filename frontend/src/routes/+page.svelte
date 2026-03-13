<script lang="ts">
  import { goto } from '$app/navigation';
  import { onDestroy, onMount } from 'svelte';
  import { cubicOut } from 'svelte/easing';
  import { fade, fly, scale } from 'svelte/transition';
  import MenuCard from '$lib/components/MenuCard.svelte';
  import { fetchCategories, type Category } from '$lib/api/categories';
  import { fetchMenu, type MenuItem } from '$lib/api/menu';

  type CartItem = MenuItem & { qty: number };

  let items: MenuItem[] = [];
  let categories: Category[] = [];
  let selectedCategory = 'all';
  let searchTerm = '';

  let cart: CartItem[] = [];
  let cartOpen = false;
  let cartPulse = false;
  let addFeedback = '';
  let feedbackTimer: ReturnType<typeof setTimeout> | null = null;
  let pulseTimer: ReturnType<typeof setTimeout> | null = null;

  let loading = true;
  let categoriesLoading = true;
  let error = '';

  $: normalizedSearch = searchTerm.trim().toLowerCase();
  $: categoryLookup = new Map<number, string>(categories.map((cat) => [cat.id, cat.name]));
  $: filteredItems = items.filter((item) => {
    const categoryMatches = selectedCategory === 'all' || item.category_id === Number(selectedCategory);
    if (!categoryMatches) {
      return false;
    }

    if (!normalizedSearch) {
      return true;
    }

    return (
      item.name.toLowerCase().includes(normalizedSearch) ||
      item.description.toLowerCase().includes(normalizedSearch)
    );
  });
  $: cartCount = cart.reduce((acc, item) => acc + item.qty, 0);
  $: cartTotal = cart.reduce((acc, item) => acc + item.price * item.qty, 0);

  function loadCart() {
    const raw = localStorage.getItem('bb_cart');
    cart = raw ? (JSON.parse(raw) as CartItem[]) : [];
  }

  function persistCart(next: CartItem[]) {
    localStorage.setItem('bb_cart', JSON.stringify(next));
  }

  function updateCart(updater: (current: CartItem[]) => CartItem[]) {
    const next = updater(cart);
    cart = next;
    persistCart(next);
  }

  function showAddedFeedback(itemName: string) {
    addFeedback = `${itemName} added to order`;
    if (feedbackTimer) {
      clearTimeout(feedbackTimer);
    }
    feedbackTimer = setTimeout(() => {
      addFeedback = '';
    }, 1600);
  }

  function triggerCartPulse() {
    cartPulse = false;
    if (pulseTimer) {
      clearTimeout(pulseTimer);
    }

    requestAnimationFrame(() => {
      cartPulse = true;
      pulseTimer = setTimeout(() => {
        cartPulse = false;
      }, 320);
    });
  }

  function addToCart(item: MenuItem) {
    updateCart((current) => {
      const idx = current.findIndex((x) => x.id === item.id);
      if (idx >= 0) {
        return current.map((entry, i) => (i === idx ? { ...entry, qty: entry.qty + 1 } : entry));
      }

      return [...current, { ...item, qty: 1 }];
    });

    if (window.matchMedia('(max-width: 819px)').matches) {
      cartOpen = true;
    }

    triggerCartPulse();
    showAddedFeedback(item.name);
  }

  function decrementFromCart(id: number) {
    updateCart((current) => {
      const idx = current.findIndex((entry) => entry.id === id);
      if (idx < 0) {
        return current;
      }

      if (current[idx].qty > 1) {
        return current.map((entry, i) => (i === idx ? { ...entry, qty: entry.qty - 1 } : entry));
      }

      return current.filter((entry) => entry.id !== id);
    });
  }

  function incrementCart(id: number) {
    updateCart((current) => {
      const idx = current.findIndex((entry) => entry.id === id);
      if (idx < 0) {
        return current;
      }
      return current.map((entry, i) => (i === idx ? { ...entry, qty: entry.qty + 1 } : entry));
    });
  }

  function clearCart() {
    cart = [];
    localStorage.removeItem('bb_cart');
  }

  async function moveToCheckout() {
    await goto('/checkout');
  }

  onMount(async () => {
    loadCart();
    try {
      const [menuPayload, categoriesPayload] = await Promise.all([fetchMenu(true), fetchCategories()]);
      items = menuPayload;
      categories = categoriesPayload;
    } catch (e) {
      error = 'Could not load menu. Please try again.';
      console.error(e);
    } finally {
      loading = false;
      categoriesLoading = false;
    }
  });

  onDestroy(() => {
    if (feedbackTimer) {
      clearTimeout(feedbackTimer);
    }
    if (pulseTimer) {
      clearTimeout(pulseTimer);
    }
  });
</script>

<main class="app-shell">
  <header class="top-app-bar" in:fly={{ y: -14, duration: 280, easing: cubicOut }}>
    <div class="brand-wrap">
      <span class="brand-mark" aria-hidden="true">BB</span>
      <div>
        <h1>BlessedBites</h1>
        <p>San Ignacio pickup menu</p>
      </div>
    </div>
    <div class="bar-actions">
      <button type="button" class="tonal" on:click={() => (cartOpen = !cartOpen)}>
        Cart
        {#if cartCount > 0}
          <span class:badge-pulse={cartPulse} class="badge">{cartCount}</span>
        {/if}
      </button>
      <button type="button" class="filled" on:click={moveToCheckout}>Checkout</button>
    </div>
  </header>

  <section class="hero" in:fade={{ duration: 260 }}>
    <h2>Order quickly, eat happily.</h2>
    <p>Browse by category, search your comfort food, and send your order in minutes.</p>

    <div class="search-row">
      <input
        id="search"
        type="search"
        bind:value={searchTerm}
        placeholder="Search menu items..."
        aria-label="Search menu items"
      />
      <button
        type="button"
        class="outlined"
        on:click={() => {
          searchTerm = '';
        }}
      >
        Clear
      </button>
    </div>
  </section>

  {#if addFeedback}
    <p class="snackbar" role="status" aria-live="polite" transition:fly={{ y: 10, duration: 220, easing: cubicOut }}>
      {addFeedback}
    </p>
  {/if}

  <section class="category-strip" aria-label="Menu categories" in:fade={{ duration: 320, delay: 40 }}>
    <button
      type="button"
      class:active-chip={selectedCategory === 'all'}
      class="chip"
      on:click={() => {
        selectedCategory = 'all';
      }}
    >
      All
    </button>

    {#if categoriesLoading}
      <span class="chip muted">Loading categories...</span>
    {:else}
      {#each categories as category}
        <button
          type="button"
          class:active-chip={selectedCategory === String(category.id)}
          class="chip"
          on:click={() => {
            selectedCategory = String(category.id);
          }}
        >
          {category.name}
        </button>
      {/each}
    {/if}
  </section>

  <section class="content-grid">
    <div>
      {#if loading}
        <p class="status">Loading menu...</p>
      {:else if error}
        <p class="status error">{error}</p>
      {:else if filteredItems.length === 0}
        <p class="status">No menu items match your filters.</p>
      {:else}
        <p class="results-copy">Showing {filteredItems.length} dishes</p>
        <section class="grid">
          {#each filteredItems as item, index (item.id)}
            <div in:fly={{ y: 16, duration: 280, delay: Math.min(index * 35, 280), easing: cubicOut }}>
              <MenuCard {item} categoryName={categoryLookup.get(item.category_id) || ''} onAdd={addToCart} />
            </div>
          {/each}
        </section>
      {/if}
    </div>

    <aside class:open={cartOpen} class="cart-panel" aria-label="Cart summary">
      <div class="cart-head">
        <h3>Your Order</h3>
        {#if cart.length > 0}
          <button type="button" class="clear-link" on:click={clearCart}>Clear</button>
        {/if}
      </div>

      {#if cart.length === 0}
        <p class="status">Your cart is empty.</p>
      {:else}
        <ul>
          {#each cart as line (line.id)}
            <li transition:scale={{ duration: 180, start: 0.95 }}>
              <div>
                <p>{line.name}</p>
                <small>${line.price.toFixed(2)} each</small>
              </div>
              <div class="qty-controls">
                <button type="button" on:click={() => decrementFromCart(line.id)} aria-label={`Remove one ${line.name}`}>
                  -
                </button>
                <span>{line.qty}</span>
                <button type="button" on:click={() => incrementCart(line.id)} aria-label={`Add one ${line.name}`}>
                  +
                </button>
              </div>
            </li>
          {/each}
        </ul>
        <p class="cart-total">Total: ${cartTotal.toFixed(2)}</p>
      {/if}

      <button type="button" class="filled wide" disabled={cart.length === 0} on:click={moveToCheckout}>
        Continue to Checkout
      </button>
    </aside>
  </section>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: "Nunito Sans", "Avenir Next", "Segoe UI", sans-serif;
    background:
      radial-gradient(circle at 12% 12%, rgb(203 34 61 / 8%), transparent 30%),
      radial-gradient(circle at 88% 8%, rgb(180 130 8 / 10%), transparent 32%),
      linear-gradient(180deg, #fff8ee 0%, #fffdf8 100%);
    color: #352220;
  }

  :global(:root) {
    --primary: #7f1d2d;
    --on-primary: #fff;
    --primary-container: #ffd9de;
    --secondary: #8d6500;
    --secondary-container: #f6dfb7;
    --on-secondary-container: #492316;
    --surface: #fffaf5;
    --surface-container: #f6ece2;
    --surface-container-high: #f2e6da;
    --on-surface: #2f1d1c;
    --on-surface-variant: #6f4744;
    --outline-variant: #e7d3c9;
  }

  .app-shell {
    max-width: 1200px;
    margin: 0 auto;
    padding: 16px 16px 28px;
  }

  .top-app-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 10px 14px;
    border-radius: 22px;
    background: color-mix(in srgb, var(--surface) 84%, white 16%);
    border: 1px solid var(--outline-variant);
    box-shadow: 0 2px 8px rgb(0 0 0 / 6%);
    position: sticky;
    top: 10px;
    z-index: 20;
  }

  .brand-wrap {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .brand-mark {
    width: 40px;
    height: 40px;
    border-radius: 12px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    color: #fff9f1;
    font-weight: 900;
    background: linear-gradient(145deg, #7f1d2d 0%, #651625 100%);
  }

  h1 {
    margin: 0;
    color: var(--on-surface);
    font-size: 1.1rem;
    line-height: 1.1;
    letter-spacing: 0.01em;
  }

  .brand-wrap p {
    margin: 0;
    color: var(--on-surface-variant);
    font-size: 0.82rem;
  }

  .bar-actions {
    display: flex;
    gap: 8px;
  }

  button {
    border: none;
    border-radius: 999px;
    padding: 0.55rem 1rem;
    font-weight: 700;
    cursor: pointer;
  }

  .filled {
    background: var(--primary);
    color: var(--on-primary);
  }

  .tonal {
    background: var(--secondary-container);
    color: var(--on-secondary-container);
    position: relative;
  }

  .outlined {
    background: transparent;
    color: var(--on-surface-variant);
    border: 1px solid var(--outline-variant);
  }

  .badge {
    margin-left: 0.45rem;
    border-radius: 999px;
    padding: 0.1rem 0.45rem;
    background: var(--primary);
    color: var(--on-primary);
    font-size: 0.74rem;
  }

  .badge-pulse {
    animation: badge-bump 320ms cubic-bezier(0.2, 0.8, 0.2, 1);
  }

  .hero {
    margin-top: 16px;
    background: linear-gradient(150deg, #fef4dc 0%, #f9e5ec 100%);
    border: 1px solid var(--outline-variant);
    border-radius: 24px;
    padding: 20px;
  }

  .hero h2 {
    margin: 0;
    color: #532122;
    font-size: 1.35rem;
  }

  .hero p {
    margin: 8px 0 0;
    color: #6f4340;
  }

  .search-row {
    margin-top: 14px;
    display: flex;
    gap: 10px;
  }

  input {
    flex: 1;
    border-radius: 16px;
    border: 1px solid var(--outline-variant);
    padding: 0.75rem 0.9rem;
    background: #fffdfa;
    color: var(--on-surface);
  }

  .category-strip {
    margin-top: 14px;
    display: flex;
    gap: 8px;
    overflow-x: auto;
    padding-bottom: 2px;
  }

  .chip {
    border-radius: 999px;
    border: 1px solid var(--outline-variant);
    background: #fff;
    color: var(--on-surface-variant);
    white-space: nowrap;
  }

  .chip.muted {
    opacity: 0.7;
    pointer-events: none;
  }

  .active-chip {
    background: var(--primary-container);
    color: #5c1524;
    border-color: #e8b6bf;
  }

  .content-grid {
    margin-top: 16px;
    display: grid;
    gap: 16px;
    align-items: start;
  }

  .results-copy {
    margin: 0 0 10px;
    color: var(--on-surface-variant);
    font-size: 0.88rem;
  }

  .grid {
    display: grid;
    gap: 12px;
  }

  .cart-panel {
    background: var(--surface-container);
    border: 1px solid var(--outline-variant);
    border-radius: 24px;
    padding: 14px;
  }

  .cart-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .cart-head h3 {
    margin: 0;
    color: var(--on-surface);
  }

  .clear-link {
    background: none;
    border: none;
    color: #7f1d2d;
    padding: 0.1rem 0.2rem;
    border-radius: 8px;
  }

  .cart-panel ul {
    list-style: none;
    margin: 12px 0;
    padding: 0;
    display: grid;
    gap: 10px;
  }

  .cart-panel li {
    display: flex;
    justify-content: space-between;
    gap: 10px;
    align-items: center;
    background: #fff;
    border-radius: 14px;
    border: 1px solid var(--outline-variant);
    padding: 10px;
  }

  .cart-panel li p {
    margin: 0;
    color: var(--on-surface);
    font-weight: 700;
    font-size: 0.9rem;
  }

  .cart-panel li small {
    color: var(--on-surface-variant);
  }

  .qty-controls {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .qty-controls button {
    width: 30px;
    height: 30px;
    border-radius: 50%;
    padding: 0;
    background: var(--secondary-container);
    color: var(--on-secondary-container);
  }

  .cart-total {
    margin: 2px 0 14px;
    color: var(--on-surface);
    font-weight: 800;
  }

  .wide {
    width: 100%;
  }

  .status {
    margin: 0;
    color: var(--on-surface-variant);
    background: var(--surface);
    border: 1px solid var(--outline-variant);
    border-radius: 16px;
    padding: 0.9rem;
  }

  .status.error {
    color: #7f1d2d;
    border-color: #f0b8c1;
    background: #ffeef1;
  }

  .snackbar {
    margin: 12px 0 0;
    width: fit-content;
    max-width: 100%;
    background: #4f2f00;
    color: #fff8e9;
    border-radius: 999px;
    padding: 0.45rem 0.8rem;
    font-size: 0.85rem;
    box-shadow: 0 2px 8px rgb(0 0 0 / 16%);
  }

  @keyframes badge-bump {
    0% {
      transform: scale(1);
    }
    40% {
      transform: scale(1.22);
    }
    100% {
      transform: scale(1);
    }
  }

  @media (min-width: 820px) {
    .content-grid {
      grid-template-columns: minmax(0, 1fr) 320px;
    }

    .cart-panel {
      position: sticky;
      top: 88px;
    }

    .grid {
      grid-template-columns: 1fr 1fr;
      gap: 14px;
    }
  }

  @media (max-width: 819px) {
    .top-app-bar {
      flex-wrap: wrap;
      position: static;
    }

    .bar-actions {
      width: 100%;
    }

    .bar-actions button {
      flex: 1;
      text-align: center;
    }

    .search-row {
      flex-direction: column;
    }

    .cart-panel {
      display: none;
    }

    .cart-panel.open {
      display: block;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .badge-pulse {
      animation: none;
    }
  }
</style>
