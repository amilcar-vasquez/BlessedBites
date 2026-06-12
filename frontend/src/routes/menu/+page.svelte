<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { fetchCategories, type Category } from '$lib/api/categories';
  import { fetchMenu, searchMenu, type MenuItem } from '$lib/api/menu';
  import CartPanel from '$lib/components/CartPanel.svelte';
  import CategoryChips from '$lib/components/CategoryChips.svelte';
  import MenuCard from '$lib/components/MenuCard.svelte';
  import Skeleton from '$lib/components/Skeleton.svelte';

  let categories = $state<Category[]>([]);
  let items = $state<MenuItem[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  const query = $derived($page.url.searchParams.get('q')?.trim() ?? '');
  const selectedCategory = $derived.by(() => {
    const raw = $page.url.searchParams.get('category');
    const id = raw ? Number(raw) : NaN;
    return Number.isFinite(id) ? id : null;
  });

  $effect(() => {
    fetchCategories()
      .then((cats) => (categories = cats))
      .catch(() => (categories = []));
  });

  $effect(() => {
    loading = true;
    error = null;
    const load = query
      ? searchMenu(query)
      : fetchMenu(true, selectedCategory ?? undefined);
    load
      .then((result) => (items = result))
      .catch(() => {
        items = [];
        error = 'We could not load the menu. Please try again.';
      })
      .finally(() => (loading = false));
  });

  function selectCategory(id: number | null) {
    const url = new URL($page.url);
    url.searchParams.delete('q');
    if (id === null) url.searchParams.delete('category');
    else url.searchParams.set('category', String(id));
    goto(`${url.pathname}${url.search}`, { keepFocus: true, noScroll: true });
  }

  const heading = $derived.by(() => {
    if (query) return `Results for “${query}”`;
    if (selectedCategory !== null) {
      return categories.find((c) => c.id === selectedCategory)?.name ?? 'Menu';
    }
    return 'All Day Menu';
  });
</script>

<svelte:head>
  <title>Menu — Blessed Bites</title>
</svelte:head>

<div class="layout">
  <!-- Desktop category sidebar -->
  <aside class="sidebar">
    <h2 class="title-lg">Categories</h2>
    <nav aria-label="Categories">
      <button
        type="button"
        class="side-link label-lg"
        class:active={selectedCategory === null && !query}
        onclick={() => selectCategory(null)}
      >
        All Day Menu
      </button>
      {#each categories as cat (cat.id)}
        <button
          type="button"
          class="side-link label-lg"
          class:active={selectedCategory === cat.id}
          onclick={() => selectCategory(cat.id)}
        >
          {cat.name}
        </button>
      {/each}
    </nav>
  </aside>

  <!-- Main content -->
  <section class="content">
    <header class="content-head">
      <h1 class="headline-lg">{heading}</h1>
      {#if query}
        <button type="button" class="clear-search label-lg" onclick={() => goto('/menu')}>
          <span class="material-symbols-outlined" aria-hidden="true">close</span>
          Clear search
        </button>
      {/if}
    </header>

    <div class="chips-row">
      <CategoryChips
        {categories}
        selected={query ? null : selectedCategory}
        onselect={selectCategory}
      />
    </div>

    {#if loading}
      <div class="grid">
        {#each Array(6) as _, i (i)}
          <Skeleton width="100%" height="320px" radius="16px" />
        {/each}
      </div>
    {:else if error}
      <div class="empty bb-card">
        <span class="material-symbols-outlined" aria-hidden="true">error</span>
        <p class="body-lg">{error}</p>
      </div>
    {:else if items.length === 0}
      <div class="empty bb-card">
        <span class="material-symbols-outlined" aria-hidden="true">search_off</span>
        <p class="body-lg">
          {query ? `Nothing matched “${query}”. Try a different craving.` : 'No dishes here yet — check back soon.'}
        </p>
      </div>
    {:else}
      <div class="grid">
        {#each items as item (item.id)}
          <MenuCard {item} />
        {/each}
      </div>
    {/if}
  </section>

  <!-- Sticky cart panel on wide screens -->
  <aside class="cart-rail">
    <div class="cart-sticky bb-card">
      <CartPanel />
    </div>
  </aside>
</div>

<style>
  .layout {
    display: grid;
    grid-template-columns: 1fr;
    gap: var(--bb-space-lg);
    padding: var(--bb-space-lg) var(--bb-margin-mobile) var(--bb-space-xl);
    width: 100%;
    max-width: 1600px;
    margin: 0 auto;
  }

  @media (min-width: 1024px) {
    .layout {
      grid-template-columns: 220px 1fr;
      padding: var(--bb-space-lg) var(--bb-margin-desktop) var(--bb-space-xl);
    }
  }

  @media (min-width: 1280px) {
    .layout {
      grid-template-columns: 220px 1fr 360px;
    }
  }

  .sidebar {
    display: none;
  }

  @media (min-width: 1024px) {
    .sidebar {
      display: block;
    }
  }

  .sidebar h2 {
    margin: 0 0 var(--bb-space-md);
    color: var(--md-sys-color-on-surface);
  }

  .sidebar nav {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-xs);
    position: sticky;
    top: 96px;
  }

  .side-link {
    border: none;
    background: transparent;
    text-align: left;
    color: var(--md-sys-color-on-surface-variant);
    padding: 10px 16px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
    transition: background-color 150ms ease, color 150ms ease;
  }

  .side-link:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .side-link.active {
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
    font-weight: 700;
  }

  .content {
    min-width: 0;
  }

  .content-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--bb-space-md);
    margin-bottom: var(--bb-space-md);
  }

  .content-head h1 {
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  .clear-search {
    display: inline-flex;
    align-items: center;
    gap: var(--bb-space-xs);
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface);
    color: var(--md-sys-color-on-surface-variant);
    padding: 8px 16px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
    white-space: nowrap;
  }

  .clear-search:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .chips-row {
    margin-bottom: var(--bb-space-lg);
  }

  @media (min-width: 1024px) {
    .chips-row {
      display: none;
    }
  }

  .grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: var(--bb-space-lg);
  }

  @media (min-width: 640px) {
    .grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (min-width: 1440px) {
    .grid {
      grid-template-columns: repeat(3, 1fr);
    }
  }

  .empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--bb-space-md);
    padding: var(--bb-space-xl);
    text-align: center;
    color: var(--md-sys-color-on-surface-variant);
  }

  .empty .material-symbols-outlined {
    font-size: 48px;
    color: var(--md-sys-color-outline);
  }

  .empty p {
    margin: 0;
  }

  .cart-rail {
    display: none;
  }

  @media (min-width: 1280px) {
    .cart-rail {
      display: block;
    }
  }

  .cart-sticky {
    position: sticky;
    top: 96px;
    padding: var(--bb-space-lg);
    max-height: calc(100dvh - 120px);
    overflow-y: auto;
  }
</style>
