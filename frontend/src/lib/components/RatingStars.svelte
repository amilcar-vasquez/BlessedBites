<script lang="ts">
  let {
    value = 0,
    interactive = false,
    onrate
  }: { value?: number; interactive?: boolean; onrate?: (rating: number) => void } = $props();

  let hovered = $state(0);

  const display = $derived(hovered || value);

  function starIcon(index: number): string {
    if (display >= index) return 'star';
    if (display >= index - 0.5) return 'star_half';
    return 'star';
  }

  function starFilled(index: number): boolean {
    return display >= index - 0.5;
  }
</script>

<div class="stars" class:interactive role={interactive ? 'radiogroup' : 'img'} aria-label="Rating: {value} out of 5">
  {#each [1, 2, 3, 4, 5] as i (i)}
    {#if interactive}
      <button
        type="button"
        aria-label="Rate {i} star{i > 1 ? 's' : ''}"
        onmouseenter={() => (hovered = i)}
        onmouseleave={() => (hovered = 0)}
        onclick={() => onrate?.(i)}
      >
        <span class="material-symbols-outlined" class:fill={starFilled(i)}>{starIcon(i)}</span>
      </button>
    {:else}
      <span class="material-symbols-outlined" class:fill={starFilled(i)}>{starIcon(i)}</span>
    {/if}
  {/each}
</div>

<style>
  .stars {
    display: inline-flex;
    gap: 2px;
    color: var(--md-sys-color-tertiary);
  }

  .material-symbols-outlined {
    font-size: 20px;
  }

  button {
    border: none;
    background: transparent;
    padding: 2px;
    cursor: pointer;
    color: inherit;
    border-radius: var(--bb-shape-sm);
  }

  button:hover {
    background: var(--md-sys-color-surface-container-high);
  }
</style>
