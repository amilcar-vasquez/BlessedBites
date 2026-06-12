<script lang="ts">
  import type { Category } from '$lib/api/categories';

  let {
    categories,
    selected = null,
    onselect
  }: {
    categories: Category[];
    selected?: number | null;
    onselect: (id: number | null) => void;
  } = $props();
</script>

<div class="chips hide-scrollbar" role="tablist" aria-label="Categories">
  <button
    type="button"
    role="tab"
    class="chip label-lg"
    class:active={selected === null}
    aria-selected={selected === null}
    onclick={() => onselect(null)}
  >
    All Day Menu
  </button>
  {#each categories as cat (cat.id)}
    <button
      type="button"
      role="tab"
      class="chip label-lg"
      class:active={selected === cat.id}
      aria-selected={selected === cat.id}
      onclick={() => onselect(cat.id)}
    >
      {cat.name}
    </button>
  {/each}
</div>

<style>
  .chips {
    display: flex;
    gap: var(--bb-space-sm);
    overflow-x: auto;
    padding: var(--bb-space-xs) 0;
  }

  .chip {
    flex-shrink: 0;
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface);
    color: var(--md-sys-color-on-surface-variant);
    padding: 8px 18px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
    white-space: nowrap;
    transition: background-color 150ms ease, color 150ms ease;
  }

  .chip:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .chip.active {
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
    border-color: transparent;
    font-weight: 700;
  }
</style>
