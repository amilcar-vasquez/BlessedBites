<script lang="ts">
  import type { MenuItem } from '$lib/api/menu';
  import { getAverageRating } from '$lib/api/ratings';
  import { addToCart } from '$lib/stores/cart';
  import { showToast } from '$lib/stores/toast';

  let { item, popular = false }: { item: MenuItem; popular?: boolean } = $props();

  let average = $state<number | null>(null);

  $effect(() => {
    getAverageRating(item.id)
      .then((r) => {
        average = r.average > 0 ? r.average : null;
      })
      .catch(() => {
        average = null;
      });
  });

  function handleAdd() {
    addToCart({ id: item.id, name: item.name, price: item.price, image_url: item.image_url });
    showToast(`${item.name} added to your order`, 'success', 2500);
  }
</script>

<article class="card">
  <div class="media">
    {#if item.image_url}
      <img src={item.image_url} alt={item.name} loading="lazy" />
    {:else}
      <div class="placeholder">
        <span class="material-symbols-outlined" aria-hidden="true">restaurant</span>
      </div>
    {/if}
    {#if popular}
      <span class="popular-badge label-sm">Popular</span>
    {/if}
    {#if average !== null}
      <span class="rating-badge label-sm">
        <span class="material-symbols-outlined fill" aria-hidden="true">star</span>
        {average.toFixed(1)}
      </span>
    {/if}
  </div>

  <div class="content">
    <div class="header-row">
      <h3 class="title-md name">{item.name}</h3>
      <span class="title-md price">${item.price.toFixed(2)}</span>
    </div>
    <p class="body-md description">{item.description}</p>
    <div class="footer-row">
      {#if average !== null}
        <span class="rating-inline label-sm">
          <span class="material-symbols-outlined fill" aria-hidden="true">star</span>
          {average.toFixed(1)}
        </span>
      {:else}
        <span class="rating-inline label-sm muted">No ratings yet</span>
      {/if}
      <button type="button" class="add-btn label-lg" onclick={handleAdd}>
        <span class="material-symbols-outlined" aria-hidden="true">add</span>
        Add
      </button>
    </div>
  </div>
</article>

<style>
  .card {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--md-sys-color-surface-container-lowest);
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: var(--bb-shape-md);
    overflow: hidden;
    box-shadow: var(--bb-elev-1);
    transition: box-shadow 200ms ease;
  }

  .card:hover {
    box-shadow: var(--bb-elev-2);
  }

  .media {
    position: relative;
    height: 192px;
    overflow: hidden;
    background: var(--md-sys-color-surface-container-high);
  }

  .media img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 500ms ease;
  }

  .card:hover .media img {
    transform: scale(1.05);
  }

  .placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--md-sys-color-on-surface-variant);
  }

  .placeholder .material-symbols-outlined {
    font-size: 48px;
  }

  .popular-badge {
    position: absolute;
    top: 12px;
    left: 12px;
    background: var(--md-sys-color-tertiary);
    color: var(--md-sys-color-on-tertiary);
    padding: 4px 8px;
    border-radius: 6px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .rating-badge {
    position: absolute;
    top: 12px;
    right: 12px;
    display: inline-flex;
    align-items: center;
    gap: 4px;
    background: color-mix(in srgb, var(--md-sys-color-surface) 90%, transparent);
    backdrop-filter: blur(4px);
    color: var(--md-sys-color-on-surface);
    font-weight: 700;
    padding: 4px 8px;
    border-radius: 6px;
    box-shadow: var(--bb-elev-1);
  }

  .rating-badge .material-symbols-outlined {
    font-size: 14px;
    color: var(--md-sys-color-tertiary);
  }

  .content {
    display: flex;
    flex-direction: column;
    flex: 1;
    padding: var(--bb-space-md);
  }

  .header-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: var(--bb-space-sm);
    margin-bottom: var(--bb-space-sm);
  }

  .name {
    margin: 0;
    color: var(--md-sys-color-on-surface);
    display: -webkit-box;
    -webkit-line-clamp: 1;
    line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .price {
    color: var(--md-sys-color-primary);
    white-space: nowrap;
  }

  .description {
    margin: 0 0 var(--bb-space-md);
    flex: 1;
    color: var(--md-sys-color-on-surface-variant);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .footer-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: auto;
  }

  .rating-inline {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    color: var(--md-sys-color-tertiary);
    font-weight: 700;
  }

  .rating-inline .material-symbols-outlined {
    font-size: 16px;
  }

  .rating-inline.muted {
    color: var(--md-sys-color-on-surface-variant);
    font-weight: 400;
  }

  .add-btn {
    display: inline-flex;
    align-items: center;
    gap: var(--bb-space-sm);
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    border: none;
    border-radius: var(--bb-shape-full);
    padding: 8px 16px;
    cursor: pointer;
    box-shadow: var(--bb-elev-2);
    transition: background-color 150ms ease, transform 100ms ease;
  }

  .add-btn:hover {
    background: var(--md-sys-color-surface-tint);
  }

  .add-btn:active {
    transform: scale(0.96);
  }

  .add-btn .material-symbols-outlined {
    font-size: 18px;
  }
</style>
