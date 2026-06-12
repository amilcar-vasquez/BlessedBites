<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchCategories, type Category } from '$lib/api/categories';
  import { fetchMenu, type MenuItem } from '$lib/api/menu';
  import MenuCard from '$lib/components/MenuCard.svelte';
  import Skeleton from '$lib/components/Skeleton.svelte';

  let categories = $state<Category[]>([]);
  let popularItems = $state<MenuItem[]>([]);
  let loading = $state(true);

  const CATEGORY_ICONS: Record<string, string> = {
    breakfast: 'egg_alt',
    lunch: 'lunch_dining',
    dinner: 'dinner_dining',
    dessert: 'cake',
    desserts: 'cake',
    drinks: 'local_cafe',
    beverages: 'local_cafe',
    sides: 'tapas',
    salads: 'eco',
    soups: 'soup_kitchen',
    specials: 'star',
    bowls: 'rice_bowl',
    burgers: 'lunch_dining',
    pizza: 'local_pizza',
    seafood: 'set_meal',
    grill: 'outdoor_grill',
    vegan: 'eco'
  };

  function iconFor(name: string): string {
    return CATEGORY_ICONS[name.trim().toLowerCase()] ?? 'restaurant';
  }

  onMount(async () => {
    try {
      const [cats, items] = await Promise.all([fetchCategories(), fetchMenu(true)]);
      categories = cats;
      popularItems = items.slice(0, 3);
    } catch {
      categories = [];
      popularItems = [];
    } finally {
      loading = false;
    }
  });

  const TESTIMONIALS = [
    {
      quote:
        'The Golden Harvest Bowl is genuinely the best lunch in the neighborhood. Fresh, fast, and the flavors are incredible.',
      name: 'Sarah M.',
      role: 'Regular since 2023'
    },
    {
      quote:
        'Ordered for the whole office — everything arrived hot and exactly as described. Blessed Bites never misses.',
      name: 'James K.',
      role: 'Office lunch hero'
    },
    {
      quote:
        'Beautiful food, honest portions, and the seasonal specials keep me coming back every single week.',
      name: 'Amara O.',
      role: 'Weekly visitor'
    }
  ];
</script>

<svelte:head>
  <title>Blessed Bites — Wholesome Food, Made with Grace</title>
  <meta
    name="description"
    content="Premium fast-casual dining. Browse the Blessed Bites menu and order fresh, seasonal dishes online."
  />
</svelte:head>

<!-- Hero -->
<section class="hero">
  <div class="hero-inner">
    <span class="hero-eyebrow label-lg">Premium Fast-Casual Dining</span>
    <h1 class="display-lg">Wholesome food,<br />made with <em>grace</em>.</h1>
    <p class="body-lg">
      Seasonal ingredients, honest cooking, and dishes prepared fresh the moment you order.
    </p>
    <div class="hero-actions">
      <a class="cta label-lg" href="/menu">
        <span class="material-symbols-outlined" aria-hidden="true">restaurant_menu</span>
        Order Now
      </a>
      <a class="cta-secondary label-lg" href="#popular">Explore Popular</a>
    </div>
  </div>
</section>

<!-- Categories -->
<section class="section">
  <h2 class="headline-lg">Browse by craving</h2>
  <div class="categories hide-scrollbar">
    {#if loading}
      {#each Array(5) as _, i (i)}
        <div class="category">
          <Skeleton width="96px" height="96px" radius="50%" />
          <Skeleton width="64px" height="14px" />
        </div>
      {/each}
    {:else}
      {#each categories as cat (cat.id)}
        <a class="category" href={`/menu?category=${cat.id}`}>
          <span class="category-circle">
            <span class="material-symbols-outlined" aria-hidden="true">{iconFor(cat.name)}</span>
          </span>
          <span class="label-lg">{cat.name}</span>
        </a>
      {/each}
    {/if}
  </div>
</section>

<!-- Popular items -->
<section class="section" id="popular">
  <div class="section-head">
    <h2 class="headline-lg">Popular right now</h2>
    <a class="see-all label-lg" href="/menu">
      Full menu
      <span class="material-symbols-outlined" aria-hidden="true">arrow_forward</span>
    </a>
  </div>
  <div class="grid">
    {#if loading}
      {#each Array(3) as _, i (i)}
        <Skeleton width="100%" height="320px" radius="16px" />
      {/each}
    {:else}
      {#each popularItems as item (item.id)}
        <MenuCard {item} popular />
      {/each}
    {/if}
  </div>
</section>

<!-- About / testimonials -->
<section class="section about" id="about">
  <h2 class="headline-lg">Loved by the neighborhood</h2>
  <div class="grid">
    {#each TESTIMONIALS as t (t.name)}
      <figure class="bb-card testimonial">
        <span class="quote-mark material-symbols-outlined fill" aria-hidden="true">format_quote</span>
        <blockquote class="body-lg">{t.quote}</blockquote>
        <figcaption>
          <span class="title-md">{t.name}</span>
          <span class="label-sm muted">{t.role}</span>
        </figcaption>
      </figure>
    {/each}
  </div>
</section>

<style>
  .hero {
    display: flex;
    align-items: center;
    min-height: clamp(440px, 70dvh, 716px);
    padding: var(--bb-space-xl) var(--bb-margin-mobile);
    background:
      radial-gradient(ellipse 80% 60% at 80% 20%, color-mix(in srgb, var(--md-sys-color-tertiary-container) 55%, transparent), transparent),
      radial-gradient(ellipse 70% 70% at 10% 90%, color-mix(in srgb, var(--md-sys-color-secondary-container) 45%, transparent), transparent),
      var(--md-sys-color-surface-container-low);
  }

  @media (min-width: 1024px) {
    .hero {
      padding: var(--bb-space-xl) var(--bb-margin-desktop);
    }
  }

  .hero-inner {
    max-width: 720px;
  }

  .hero-eyebrow {
    display: inline-block;
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
    padding: 6px 16px;
    border-radius: var(--bb-shape-full);
    margin-bottom: var(--bb-space-md);
  }

  .hero h1 {
    margin: 0 0 var(--bb-space-md);
    color: var(--md-sys-color-on-surface);
  }

  .hero h1 em {
    color: var(--md-sys-color-primary);
    font-style: italic;
  }

  .hero p {
    margin: 0 0 var(--bb-space-xl);
    color: var(--md-sys-color-on-surface-variant);
    max-width: 480px;
  }

  .hero-actions {
    display: flex;
    flex-wrap: wrap;
    gap: var(--bb-space-md);
  }

  .cta {
    display: inline-flex;
    align-items: center;
    gap: var(--bb-space-sm);
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    padding: 14px 28px;
    border-radius: var(--bb-shape-full);
    text-decoration: none;
    box-shadow: var(--bb-elev-2);
    transition: transform 200ms ease, box-shadow 200ms ease;
  }

  .cta:hover {
    transform: translateY(-2px);
    box-shadow: var(--bb-elev-3);
  }

  .cta-secondary {
    display: inline-flex;
    align-items: center;
    padding: 14px 28px;
    border-radius: var(--bb-shape-full);
    border: 1px solid var(--md-sys-color-outline);
    color: var(--md-sys-color-primary);
    text-decoration: none;
    transition: background-color 150ms ease;
  }

  .cta-secondary:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .section {
    padding: var(--bb-space-xl) var(--bb-margin-mobile);
  }

  @media (min-width: 1024px) {
    .section {
      padding: var(--bb-space-xl) var(--bb-margin-desktop);
    }
  }

  .section h2 {
    margin: 0 0 var(--bb-space-lg);
    color: var(--md-sys-color-on-surface);
  }

  .section-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--bb-space-md);
  }

  .see-all {
    display: inline-flex;
    align-items: center;
    gap: var(--bb-space-xs);
    color: var(--md-sys-color-primary);
    text-decoration: none;
    white-space: nowrap;
  }

  .see-all:hover {
    text-decoration: underline;
  }

  .categories {
    display: flex;
    gap: var(--bb-space-lg);
    overflow-x: auto;
    padding-bottom: var(--bb-space-sm);
  }

  .category {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--bb-space-sm);
    text-decoration: none;
    color: var(--md-sys-color-on-surface);
    flex-shrink: 0;
  }

  .category-circle {
    width: 96px;
    height: 96px;
    border-radius: var(--bb-shape-full);
    background: var(--md-sys-color-surface-container-high);
    border: 1px solid var(--md-sys-color-outline-variant);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 150ms ease, transform 200ms ease;
  }

  .category-circle .material-symbols-outlined {
    font-size: 36px;
    color: var(--md-sys-color-primary);
  }

  .category:hover .category-circle {
    background: var(--md-sys-color-secondary-container);
    transform: translateY(-4px);
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

  @media (min-width: 1024px) {
    .grid {
      grid-template-columns: repeat(3, 1fr);
    }
  }

  .about {
    background: var(--md-sys-color-surface-container-low);
  }

  .testimonial {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
    padding: var(--bb-space-lg);
    margin: 0;
  }

  .quote-mark {
    font-size: 40px;
    color: var(--md-sys-color-tertiary);
  }

  blockquote {
    margin: 0;
    color: var(--md-sys-color-on-surface);
    flex: 1;
  }

  figcaption {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .muted {
    color: var(--md-sys-color-on-surface-variant);
  }
</style>
