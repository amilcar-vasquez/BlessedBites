<script lang="ts">
  let {
    width = '100%',
    height = '16px',
    radius = 'var(--bb-shape-sm)'
  }: { width?: string; height?: string; radius?: string } = $props();
</script>

<div class="bb-skeleton" style:width style:height style:border-radius={radius} aria-hidden="true"></div>
