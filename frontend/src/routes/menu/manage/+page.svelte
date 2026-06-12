<script lang="ts">
  import { onMount } from 'svelte';
  import {
    createAdminCategory,
    createAdminMenu,
    deleteAdminCategory,
    deleteAdminMenu,
    listAdminMenu,
    updateAdminMenu
  } from '$lib/api/admin';
  import { fetchCategories, type Category } from '$lib/api/categories';
  import { normalizeImageUrl, type MenuItem } from '$lib/api/menu';
  import { showToast } from '$lib/stores/toast';
  import Skeleton from '$lib/components/Skeleton.svelte';

  let items = $state<MenuItem[]>([]);
  let categories = $state<Category[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  // Dialog state
  let dialogOpen = $state(false);
  let editing = $state<MenuItem | null>(null);
  let saving = $state(false);
  let form = $state({
    name: '',
    description: '',
    price: '',
    category_id: '',
    image_url: '',
    is_active: true
  });

  // Delete confirm
  let deleting = $state<MenuItem | null>(null);

  // Category panel
  let newCategoryName = $state('');
  let categoryBusy = $state(false);

  async function refresh() {
    try {
      const [m, c] = await Promise.all([listAdminMenu(), fetchCategories()]);
      items = m;
      categories = c;
      error = null;
    } catch {
      error = 'Could not load menu data.';
    } finally {
      loading = false;
    }
  }

  onMount(refresh);

  function categoryName(id: number): string {
    return categories.find((c) => c.id === id)?.name ?? '—';
  }

  function openCreate() {
    editing = null;
    form = {
      name: '',
      description: '',
      price: '',
      category_id: categories[0] ? String(categories[0].id) : '',
      image_url: '',
      is_active: true
    };
    dialogOpen = true;
  }

  function openEdit(item: MenuItem) {
    editing = item;
    form = {
      name: item.name,
      description: item.description,
      price: String(item.price),
      category_id: String(item.category_id),
      image_url: item.image_url ?? '',
      is_active: item.is_active !== false
    };
    dialogOpen = true;
  }

  async function save(e: Event) {
    e.preventDefault();
    if (saving) return;
    const price = Number(form.price);
    const categoryId = Number(form.category_id);
    if (!form.name.trim() || !Number.isFinite(price) || price <= 0 || !Number.isFinite(categoryId)) {
      showToast('Please fill in name, a valid price, and a category.', 'error');
      return;
    }
    saving = true;
    const payload = {
      name: form.name.trim(),
      description: form.description.trim(),
      price,
      category_id: categoryId,
      image_url: form.image_url.trim(),
      is_active: form.is_active
    };
    try {
      if (editing) {
        await updateAdminMenu(editing.id, payload);
        showToast(`${payload.name} updated`, 'success');
      } else {
        await createAdminMenu(payload);
        showToast(`${payload.name} added to the menu`, 'success');
      }
      dialogOpen = false;
      await refresh();
    } catch {
      showToast('Could not save the dish. Please try again.', 'error');
    } finally {
      saving = false;
    }
  }

  async function confirmDelete() {
    if (!deleting) return;
    try {
      await deleteAdminMenu(deleting.id);
      showToast(`${deleting.name} deleted`, 'success');
      deleting = null;
      await refresh();
    } catch {
      showToast('Could not delete the dish.', 'error');
    }
  }

  async function addCategory(e: Event) {
    e.preventDefault();
    const name = newCategoryName.trim();
    if (!name || categoryBusy) return;
    categoryBusy = true;
    try {
      await createAdminCategory(name);
      newCategoryName = '';
      showToast(`Category “${name}” created`, 'success');
      await refresh();
    } catch {
      showToast('Could not create the category.', 'error');
    } finally {
      categoryBusy = false;
    }
  }

  async function removeCategory(cat: Category) {
    if (categoryBusy) return;
    categoryBusy = true;
    try {
      await deleteAdminCategory(cat.id);
      showToast(`Category “${cat.name}” deleted`, 'success');
      await refresh();
    } catch {
      showToast('Could not delete the category — it may still have dishes.', 'error');
    } finally {
      categoryBusy = false;
    }
  }
</script>

<svelte:head>
  <title>Menu Management — Blessed Admin</title>
</svelte:head>

<header class="head">
  <div>
    <h1 class="headline-lg">Menu &amp; Categories</h1>
    <p class="body-md muted">Manage dishes, pricing, and availability.</p>
  </div>
  <button type="button" class="bb-btn-primary" onclick={openCreate}>
    <span class="material-symbols-outlined" aria-hidden="true">add</span>
    New Dish
  </button>
</header>

{#if error}
  <p class="bb-form-error" role="alert">{error}</p>
{/if}

<div class="panels">
  <!-- Menu table -->
  <section class="bb-card panel">
    {#if loading}
      <Skeleton width="100%" height="320px" radius="8px" />
    {:else if items.length === 0}
      <p class="body-md muted">No dishes yet — create your first one.</p>
    {:else}
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th class="label-sm">Dish</th>
              <th class="label-sm">Category</th>
              <th class="label-sm num">Price</th>
              <th class="label-sm">Status</th>
              <th class="label-sm actions-col">Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each items as item (item.id)}
              <tr>
                <td>
                  <div class="dish">
                    {#if item.image_url}
                      <img src={normalizeImageUrl(item.image_url)} alt="" loading="lazy" />
                    {:else}
                      <span class="thumb material-symbols-outlined" aria-hidden="true">restaurant</span>
                    {/if}
                    <div class="dish-text">
                      <span class="title-md">{item.name}</span>
                      <span class="label-sm muted clamp">{item.description}</span>
                    </div>
                  </div>
                </td>
                <td class="body-md">{categoryName(item.category_id)}</td>
                <td class="body-md num strong">${item.price.toFixed(2)}</td>
                <td>
                  <span class="chip label-sm" class:active-chip={item.is_active !== false}>
                    {item.is_active !== false ? 'Active' : 'Hidden'}
                  </span>
                </td>
                <td>
                  <div class="row-actions">
                    <button type="button" class="icon-btn" aria-label={`Edit ${item.name}`} onclick={() => openEdit(item)}>
                      <span class="material-symbols-outlined">edit</span>
                    </button>
                    <button type="button" class="icon-btn danger" aria-label={`Delete ${item.name}`} onclick={() => (deleting = item)}>
                      <span class="material-symbols-outlined">delete</span>
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </section>

  <!-- Categories panel -->
  <section class="bb-card panel">
    <h2 class="title-lg">Categories</h2>
    <form class="cat-form" onsubmit={addCategory}>
      <input
        type="text"
        placeholder="New category name"
        aria-label="New category name"
        bind:value={newCategoryName}
      />
      <button type="submit" class="bb-btn-primary" disabled={categoryBusy || !newCategoryName.trim()}>
        Add
      </button>
    </form>
    {#if loading}
      <Skeleton width="100%" height="120px" radius="8px" />
    {:else if categories.length === 0}
      <p class="body-md muted">No categories yet.</p>
    {:else}
      <ul class="cat-list">
        {#each categories as cat (cat.id)}
          <li>
            <span class="body-lg">{cat.name}</span>
            <span class="label-sm muted count">
              {items.filter((i) => i.category_id === cat.id).length} dishes
            </span>
            <button type="button" class="icon-btn danger" aria-label={`Delete category ${cat.name}`} onclick={() => removeCategory(cat)}>
              <span class="material-symbols-outlined">delete</span>
            </button>
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</div>

<!-- Create/Edit dialog -->
{#if dialogOpen}
  <div class="overlay" role="presentation" onclick={(e) => e.target === e.currentTarget && (dialogOpen = false)}>
    <div class="dialog bb-card" role="dialog" aria-modal="true" aria-label={editing ? 'Edit dish' : 'New dish'}>
      <h2 class="headline-md">{editing ? `Edit ${editing.name}` : 'New Dish'}</h2>
      <form onsubmit={save}>
        <label class="bb-field">
          <span>Name</span>
          <input type="text" required bind:value={form.name} />
        </label>
        <label class="bb-field">
          <span>Description</span>
          <textarea rows="3" bind:value={form.description}></textarea>
        </label>
        <div class="field-row">
          <label class="bb-field">
            <span>Price ($)</span>
            <input type="number" step="0.01" min="0.01" required bind:value={form.price} />
          </label>
          <label class="bb-field">
            <span>Category</span>
            <select required bind:value={form.category_id}>
              {#each categories as cat (cat.id)}
                <option value={String(cat.id)}>{cat.name}</option>
              {/each}
            </select>
          </label>
        </div>
        <label class="bb-field">
          <span>Image URL</span>
          <input type="text" placeholder="/uploads/dish.jpg" bind:value={form.image_url} />
        </label>
        <label class="toggle">
          <input type="checkbox" bind:checked={form.is_active} />
          <span class="label-lg">Visible on the customer menu</span>
        </label>
        <div class="dialog-actions">
          <button type="button" class="text-btn label-lg" onclick={() => (dialogOpen = false)}>Cancel</button>
          <button type="submit" class="bb-btn-primary" disabled={saving}>
            {saving ? 'Saving…' : editing ? 'Save changes' : 'Create dish'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Delete confirm -->
{#if deleting}
  <div class="overlay" role="presentation" onclick={(e) => e.target === e.currentTarget && (deleting = null)}>
    <div class="dialog bb-card" role="alertdialog" aria-modal="true" aria-label="Confirm deletion">
      <h2 class="headline-md">Delete “{deleting.name}”?</h2>
      <p class="body-md muted">This removes the dish permanently. This action cannot be undone.</p>
      <div class="dialog-actions">
        <button type="button" class="text-btn label-lg" onclick={() => (deleting = null)}>Cancel</button>
        <button type="button" class="danger-btn label-lg" onclick={confirmDelete}>Delete</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--bb-space-md);
    margin-bottom: var(--bb-space-lg);
    flex-wrap: wrap;
  }

  .head h1 {
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  .head p {
    margin: var(--bb-space-xs) 0 0;
  }

  .muted {
    color: var(--md-sys-color-on-surface-variant);
  }

  .panels {
    display: grid;
    grid-template-columns: 1fr;
    gap: var(--bb-space-lg);
    align-items: start;
  }

  @media (min-width: 1200px) {
    .panels {
      grid-template-columns: 3fr 1fr;
    }
  }

  .panel {
    padding: var(--bb-space-lg);
  }

  .panel h2 {
    margin: 0 0 var(--bb-space-md);
    color: var(--md-sys-color-on-surface);
  }

  .table-wrap {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th {
    text-align: left;
    color: var(--md-sys-color-on-surface-variant);
    text-transform: uppercase;
    padding: 10px 12px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  td {
    padding: 12px;
    color: var(--md-sys-color-on-surface);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    vertical-align: middle;
  }

  tbody tr:nth-child(even) {
    background: var(--md-sys-color-surface-container-low);
  }

  tbody tr:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .num {
    text-align: right;
  }

  .strong {
    font-weight: 700;
  }

  .actions-col {
    width: 96px;
  }

  .dish {
    display: flex;
    align-items: center;
    gap: var(--bb-space-md);
    min-width: 220px;
  }

  .dish img,
  .thumb {
    width: 48px;
    height: 48px;
    border-radius: var(--bb-shape-sm);
    object-fit: cover;
    flex-shrink: 0;
  }

  .thumb {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-outline);
  }

  .dish-text {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
  }

  .clamp {
    display: -webkit-box;
    -webkit-line-clamp: 1;
    line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .chip {
    display: inline-block;
    padding: 4px 12px;
    border-radius: var(--bb-shape-full);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface-variant);
  }

  .chip.active-chip {
    border-color: transparent;
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
  }

  .row-actions {
    display: flex;
    gap: var(--bb-space-xs);
  }

  .icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border: none;
    background: transparent;
    border-radius: var(--bb-shape-full);
    color: var(--md-sys-color-on-surface-variant);
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  .icon-btn:hover {
    background: var(--md-sys-color-surface-container-highest);
  }

  .icon-btn.danger:hover {
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-on-error-container);
  }

  .icon-btn .material-symbols-outlined {
    font-size: 20px;
  }

  .cat-form {
    display: flex;
    gap: var(--bb-space-sm);
    margin-bottom: var(--bb-space-md);
  }

  .cat-form input {
    flex: 1;
    min-width: 0;
    border: 1px solid var(--md-sys-color-outline);
    border-radius: var(--bb-shape-sm);
    background: var(--md-sys-color-surface);
    color: var(--md-sys-color-on-surface);
    padding: 10px 14px;
    font-family: var(--md-ref-typeface-plain);
    font-size: 14px;
    outline: none;
  }

  .cat-form input:focus {
    border-color: var(--md-sys-color-primary);
  }

  .cat-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
  }

  .cat-list li {
    display: flex;
    align-items: center;
    gap: var(--bb-space-sm);
    padding: 10px 0;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  .cat-list li:last-child {
    border-bottom: none;
  }

  .cat-list .body-lg {
    flex: 1;
    color: var(--md-sys-color-on-surface);
  }

  .count {
    white-space: nowrap;
  }

  .overlay {
    position: fixed;
    inset: 0;
    z-index: 100;
    background: color-mix(in srgb, var(--md-sys-color-scrim, #000) 45%, transparent);
    backdrop-filter: blur(2px);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--bb-space-md);
  }

  .dialog {
    width: 100%;
    max-width: 520px;
    max-height: calc(100dvh - 48px);
    overflow-y: auto;
    padding: var(--bb-space-lg);
    border-radius: var(--bb-shape-lg);
  }

  .dialog h2 {
    margin: 0 0 var(--bb-space-md);
    color: var(--md-sys-color-on-surface);
  }

  .dialog form {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }

  .field-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--bb-space-md);
  }

  .toggle {
    display: flex;
    align-items: center;
    gap: var(--bb-space-sm);
    cursor: pointer;
    color: var(--md-sys-color-on-surface);
  }

  .toggle input {
    width: 18px;
    height: 18px;
    accent-color: var(--md-sys-color-primary);
  }

  .dialog-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--bb-space-sm);
    margin-top: var(--bb-space-sm);
  }

  .text-btn {
    border: none;
    background: transparent;
    color: var(--md-sys-color-primary);
    padding: 12px 20px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
  }

  .text-btn:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .danger-btn {
    border: none;
    background: var(--md-sys-color-error);
    color: var(--md-sys-color-on-error);
    padding: 12px 24px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
  }

  .danger-btn:hover {
    opacity: 0.9;
  }
</style>
