<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchCategories, type Category } from '$lib/api/categories';
  import {
    createAdminCategory,
    createAdminMenu,
    deleteAdminCategory,
    deleteAdminMenu,
    listAdminMenu,
    updateAdminMenu
  } from '$lib/api/admin';
  import type { MenuItem } from '$lib/api/menu';

  type FormState = {
    id?: number;
    name: string;
    description: string;
    price: string;
    category_id: string;
    image_url: string;
    is_active: boolean;
  };

  let items: MenuItem[] = [];
  let categories: Category[] = [];
  let loading = true;
  let error = '';
  let categoryError = '';
  let categoryName = '';
  let submitting = false;

  let form: FormState = {
    name: '',
    description: '',
    price: '',
    category_id: '',
    image_url: '',
    is_active: true
  };

  async function load() {
    loading = true;
    error = '';
    try {
      const [menuItems, categoryItems] = await Promise.all([listAdminMenu(), fetchCategories()]);
      items = menuItems;
      categories = categoryItems;
    } catch (e) {
      console.error(e);
      error = 'Could not load admin menu data. Admin login required.';
    } finally {
      loading = false;
    }
  }

  function resetForm() {
    form = {
      name: '',
      description: '',
      price: '',
      category_id: categories.length > 0 ? String(categories[0].id) : '',
      image_url: '',
      is_active: true
    };
  }

  function startEdit(item: MenuItem) {
    form = {
      id: item.id,
      name: item.name,
      description: item.description,
      price: String(item.price),
      category_id: String(item.category_id),
      image_url: item.image_url || '',
      is_active: item.is_active ?? true
    };
  }

  async function submitMenu() {
    error = '';
    submitting = true;

    try {
      const payload = {
        name: form.name,
        description: form.description,
        price: Number(form.price),
        category_id: Number(form.category_id),
        image_url: form.image_url,
        is_active: form.is_active
      };

      if (form.id) {
        await updateAdminMenu(form.id, payload);
      } else {
        await createAdminMenu(payload);
      }

      await load();
      resetForm();
    } catch (e) {
      console.error(e);
      error = 'Could not save menu item.';
    } finally {
      submitting = false;
    }
  }

  async function removeMenuItem(id: number) {
    if (!confirm('Delete this menu item?')) return;
    try {
      await deleteAdminMenu(id);
      await load();
    } catch (e) {
      console.error(e);
      error = 'Could not delete menu item.';
    }
  }

  async function submitCategory() {
    categoryError = '';
    try {
      await createAdminCategory(categoryName.trim());
      categoryName = '';
      await load();
    } catch (e) {
      console.error(e);
      categoryError = 'Could not create category.';
    }
  }

  async function removeCategory(id: number) {
    if (!confirm('Delete this category?')) return;
    try {
      await deleteAdminCategory(id);
      await load();
    } catch (e) {
      console.error(e);
      categoryError = 'Could not delete category.';
    }
  }

  onMount(async () => {
    await load();
    resetForm();
  });
</script>

<main class="shell">
  <h1>Menu Management</h1>
  <p class="muted">Add categories and manage menu items.</p>

  <section class="layout">
    <aside class="card">
      <h2>Categories</h2>
      <form on:submit|preventDefault={submitCategory} class="stack">
        <label for="category">New category</label>
        <input id="category" bind:value={categoryName} required />
        <button class="btn mustard" type="submit">Add Category</button>
      </form>
      {#if categoryError}<p class="error">{categoryError}</p>{/if}

      <ul class="stack list">
        {#each categories as cat (cat.id)}
          <li>
            <span>{cat.name}</span>
            <button class="icon" type="button" on:click={() => removeCategory(cat.id)}>Delete</button>
          </li>
        {/each}
      </ul>
    </aside>

    <section>
      <article class="card">
        <h2>{form.id ? 'Edit Menu Item' : 'Add Menu Item'}</h2>
        <form on:submit|preventDefault={submitMenu} class="stack">
          <label for="name">Name</label>
          <input id="name" bind:value={form.name} required />

          <label for="description">Description</label>
          <textarea id="description" bind:value={form.description} rows="3" required></textarea>

          <label for="price">Price (BZD)</label>
          <input id="price" type="number" step="0.01" min="0" bind:value={form.price} required />

          <label for="catSelect">Category</label>
          <select id="catSelect" bind:value={form.category_id} required>
            {#each categories as cat (cat.id)}
              <option value={String(cat.id)}>{cat.name}</option>
            {/each}
          </select>

          <label for="image">Image URL</label>
          <input id="image" bind:value={form.image_url} placeholder="/uploads/filename.jpg" />

          <label class="toggle">
            <input type="checkbox" bind:checked={form.is_active} /> Active
          </label>

          <div class="row">
            <button class="btn" type="submit" disabled={submitting}>{submitting ? 'Saving...' : 'Save Item'}</button>
            {#if form.id}
              <button class="btn ghost" type="button" on:click={resetForm}>Cancel Edit</button>
            {/if}
          </div>
        </form>
      </article>

      <article class="card">
        <h2>Existing Items</h2>
        {#if loading}
          <p>Loading menu...</p>
        {:else}
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Category</th>
                  <th>Price</th>
                  <th>State</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {#each items as item (item.id)}
                  <tr>
                    <td>{item.name}</td>
                    <td>{categories.find((c) => c.id === item.category_id)?.name || '-'}</td>
                    <td>${item.price.toFixed(2)}</td>
                    <td>{item.is_active ? 'Active' : 'Inactive'}</td>
                    <td>
                      <button class="icon" type="button" on:click={() => startEdit(item)}>Edit</button>
                      <button class="icon danger" type="button" on:click={() => removeMenuItem(item.id)}>Delete</button>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
        {#if error}<p class="error">{error}</p>{/if}
      </article>
    </section>
  </section>
</main>

<style>
  .shell { max-width: 1200px; margin: 1.2rem auto; padding: 1rem; }
  .muted { color: #6f4744; }
  .layout { display: grid; grid-template-columns: 280px 1fr; gap: 1rem; }
  .card { background: #fffaf5; border: 1px solid #e7d3c9; border-radius: 20px; padding: 1rem; }
  .stack { display: grid; gap: 0.55rem; }
  .row { display: flex; gap: 0.6rem; flex-wrap: wrap; }
  input, textarea, select { border: 1px solid #d9c5bc; border-radius: 12px; padding: 0.6rem 0.8rem; font: inherit; }
  .btn { border: none; border-radius: 999px; padding: 0.6rem 0.95rem; background: #7f1d2d; color: #fff; font-weight: 700; }
  .btn.mustard { background: #8d6500; }
  .btn.ghost { background: #eee0d5; color: #482826; }
  .toggle { display: inline-flex; align-items: center; gap: 0.45rem; font-weight: 700; }
  .list { list-style: none; margin: 0.9rem 0 0; padding: 0; }
  .list li { display: flex; justify-content: space-between; gap: 0.7rem; align-items: center; }
  .icon { border: none; border-radius: 999px; padding: 0.32rem 0.7rem; background: #f6dfb7; color: #492316; font-weight: 700; }
  .icon.danger { background: #ffd9de; color: #5c1524; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { text-align: left; padding: 0.55rem 0.4rem; border-bottom: 1px solid #f0dfd6; }
  .error { color: #8a1732; }
  @media (max-width: 900px) {
    .layout { grid-template-columns: 1fr; }
  }
</style>
