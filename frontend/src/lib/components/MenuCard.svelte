<script lang="ts">
  import type { MenuItem } from '$lib/api/menu';

  export let item: MenuItem;
  export let onAdd: (item: MenuItem) => void;
  export let categoryName = '';
</script>

<article class="menu-card">
  <img src={item.image_url || '/icons/placeholder.svg'} alt={item.name} loading="lazy" />
  <div class="content">
    <div class="meta-row">
      <h3>{item.name}</h3>
      {#if item.popular}
        <span class="pill">Popular</span>
      {/if}
    </div>
    {#if categoryName}
      <p class="category">{categoryName}</p>
    {/if}
    <p>{item.description}</p>
    <div class="footer">
      <strong>${item.price.toFixed(2)}</strong>
      <button type="button" on:click={() => onAdd(item)} aria-label={`Add ${item.name} to cart`}>Add</button>
    </div>
  </div>
</article>

<style>
  .menu-card {
    display: grid;
    grid-template-columns: 92px 1fr;
    gap: 12px;
    background: var(--surface, #fffaf5);
    border: 1px solid var(--outline-variant, #e5d4c9);
    border-radius: 20px;
    padding: 12px;
    box-shadow: 0 1px 2px rgb(0 0 0 / 10%);
  }

  img {
    width: 92px;
    height: 92px;
    border-radius: 16px;
    object-fit: cover;
    background: #efe1d4;
  }

  .meta-row {
    display: flex;
    align-items: center;
    gap: 8px;
    justify-content: space-between;
  }

  h3 {
    margin: 0;
    font-size: 1.03rem;
    color: var(--on-surface, #2e1d1c);
    font-weight: 700;
    line-height: 1.25;
  }

  .pill {
    font-size: 0.72rem;
    font-weight: 700;
    border-radius: 999px;
    background: var(--secondary-container, #f6dfb7);
    color: var(--on-secondary-container, #492316);
    padding: 0.2rem 0.55rem;
    white-space: nowrap;
  }

  .category {
    margin: 4px 0 2px;
    color: var(--on-surface-variant, #684240);
    font-size: 0.78rem;
    text-transform: uppercase;
    letter-spacing: 0.03em;
    font-weight: 600;
  }

  p {
    margin: 4px 0 10px;
    color: var(--on-surface-variant, #684240);
    font-size: 0.92rem;
  }

  .footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  button {
    min-width: 82px;
    min-height: 40px;
    border-radius: 999px;
    border: none;
    color: var(--on-primary, #fff);
    background: var(--primary, #7f1d2d);
    font-weight: 700;
    cursor: pointer;
  }
</style>
