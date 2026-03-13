<script lang="ts">
  import { page } from '$app/stores';
</script>

<main>
  <h1>Category: {$page.params.category}</h1>
  <p>This route is ready for category-specific filtering with `/api/v1/menu?category=` support.</p>
</main>
