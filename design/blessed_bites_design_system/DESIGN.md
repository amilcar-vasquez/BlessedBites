---
name: Blessed Bites Design System
colors:
  surface: '#fff8f6'
  surface-dim: '#e8d6d1'
  surface-bright: '#fff8f6'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#fff1ec'
  surface-container: '#fceae4'
  surface-container-high: '#f6e4df'
  surface-container-highest: '#f1dfd9'
  on-surface: '#231a16'
  on-surface-variant: '#55433c'
  inverse-surface: '#382e2b'
  inverse-on-surface: '#ffede7'
  outline: '#88726b'
  outline-variant: '#dbc1b8'
  surface-tint: '#9a4523'
  primary: '#994422'
  on-primary: '#ffffff'
  primary-container: '#b85c38'
  on-primary-container: '#ffffff'
  inverse-primary: '#ffb59a'
  secondary: '#486550'
  on-secondary: '#ffffff'
  secondary-container: '#caebd0'
  on-secondary-container: '#4e6b56'
  tertiary: '#7c5700'
  on-tertiary: '#ffffff'
  tertiary-container: '#9b6e06'
  on-tertiary-container: '#ffffff'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#ffdbcf'
  primary-fixed-dim: '#ffb59a'
  on-primary-fixed: '#380d00'
  on-primary-fixed-variant: '#7b2f0e'
  secondary-fixed: '#caebd0'
  secondary-fixed-dim: '#aeceb5'
  on-secondary-fixed: '#042010'
  on-secondary-fixed-variant: '#314d39'
  tertiary-fixed: '#ffdeaa'
  tertiary-fixed-dim: '#f5bd58'
  on-tertiary-fixed: '#271900'
  on-tertiary-fixed-variant: '#5f4100'
  background: '#fff8f6'
  on-background: '#231a16'
  surface-variant: '#f1dfd9'
typography:
  display-lg:
    fontFamily: Montserrat
    fontSize: 57px
    fontWeight: '700'
    lineHeight: 64px
    letterSpacing: -0.25px
  headline-lg:
    fontFamily: Montserrat
    fontSize: 32px
    fontWeight: '600'
    lineHeight: 40px
  headline-lg-mobile:
    fontFamily: Montserrat
    fontSize: 28px
    fontWeight: '600'
    lineHeight: 36px
  headline-md:
    fontFamily: Montserrat
    fontSize: 28px
    fontWeight: '600'
    lineHeight: 36px
  title-lg:
    fontFamily: Montserrat
    fontSize: 22px
    fontWeight: '500'
    lineHeight: 28px
  title-md:
    fontFamily: Inter
    fontSize: 16px
    fontWeight: '600'
    lineHeight: 24px
    letterSpacing: 0.15px
  body-lg:
    fontFamily: Inter
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
    letterSpacing: 0.5px
  body-md:
    fontFamily: Inter
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
    letterSpacing: 0.25px
  label-lg:
    fontFamily: Inter
    fontSize: 14px
    fontWeight: '500'
    lineHeight: 20px
    letterSpacing: 0.1px
  label-sm:
    fontFamily: Inter
    fontSize: 11px
    fontWeight: '500'
    lineHeight: 16px
    letterSpacing: 0.5px
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  base: 8px
  xs: 4px
  sm: 8px
  md: 16px
  lg: 24px
  xl: 32px
  gutter: 16px
  margin-mobile: 16px
  margin-desktop: 64px
---

## Brand & Style
The design system for Blessed Bites bridges the gap between high-end hospitality and modern digital efficiency. It is rooted in **Corporate Modern** principles with a heavy influence from **Material Design 3 (MD3)**, focusing on accessibility, warmth, and professional clarity. 

The target audience ranges from busy professionals seeking a healthy lunch to families ordering dinner. The UI must evoke a sense of "culinary craft" through generous whitespace and high-quality food photography, while maintaining the systematic reliability of a top-tier SaaS platform. The aesthetic is "Premium Fast-Casual"—clean, structured, and inviting, avoiding the cluttered "fast food" look in favor of a curated, editorial experience.

## Colors
The palette is inspired by natural Caribbean ingredients and earth tones. 
- **Primary (Terracotta):** Used for key actions, brand moments, and active states. It provides the "warmth" of the brand.
- **Secondary (Sage):** Used for health-focused labels, dietary tags (vegan, gluten-free), and success states.
- **Tertiary (Gold):** Reserved for highlights, ratings, and loyalty program elements.
- **Neutrals:** The background is a warm off-white to prevent screen fatigue and maintain a "paper" feel. Typography uses a Charcoal gray rather than pure black to keep the interface feeling premium and soft.

## Typography
The system uses a pairing of **Montserrat** for impactful storytelling and **Inter** for functional density.
- **Headlines:** Montserrat's geometric clarity provides a modern, professional look for food categories and section titles.
- **Body & Labels:** Inter is used for all functional text, nutrition facts, and admin data to ensure maximum readability at small sizes.
- **Hierarchy:** Use `title-md` for food item names in cards and `label-lg` for button text.

## Layout & Spacing
This design system follows a **12-column fluid grid** for desktop and a **4-column grid** for mobile.
- **Spacing Scale:** Built on an 8px base unit. Use 16px (md) for most internal card padding and 24px (lg) for vertical section spacing.
- **Mobile:** Elements should reach the 16px side margins. Horizontal scrolling is used for Category Chips and "Featured" collections to maintain a compact vertical footprint.
- **Admin Dashboard:** Uses a "contained" layout with a 240px fixed left navigation rail and a fluid content area for data tables.

## Elevation & Depth
Consistent with MD3, depth is communicated through **Tonal Layers** and subtle shadows.
- **Level 0 (Flat):** Background surface.
- **Level 1 (Card):** Used for standard food cards. Surface color is a slightly lighter tint than the background with a 1px soft border (#E0DDD7) instead of a heavy shadow.
- **Level 2 (Hover/Floating):** Used for buttons and active cards. Features a soft, diffused shadow (0px 4px 12px rgba(45, 41, 38, 0.08)).
- **Level 3 (Modals/Drawers):** Highest elevation. Uses a 20% backdrop blur (glassmorphism Lite) on the overlay to maintain context of the restaurant menu behind the cart.

## Shapes
We use a **Rounded** (Level 2) approach to mirror the friendly, organic nature of food.
- **Small Components:** Buttons and Input fields use 8px (0.5rem) radius.
- **Medium Components:** Food Cards and Search Bars use 16px (1rem) radius.
- **Large Components:** Cart Drawers and Checkout Modals use 24px (1.5rem) top-corner radius for a "sheet" effect.
- **Chips:** Fully pill-shaped for easy tap targets.

## Components
- **Buttons:** Primary buttons use the Terracotta fill with white text. Secondary buttons use a Sage outline.
- **Food Cards:** Images should have a subtle zoom-in effect on hover. Pricing is placed in the bottom right using `title-md`.
- **Search Bar:** Elevated at Level 1, featuring a "Search flavors..." placeholder and the Tertiary Gold for the search icon.
- **Cart Drawer:** Enters from the right on desktop, or slides up as a bottom sheet on mobile. Use a sticky "Checkout" button at the bottom.
- **Admin Tables:** Clean, zebra-striped rows using the Neutral Surface color. Status indicators (e.g., "Out of Stock") use small, filled chips.
- **Quantity Selector:** A pill-shaped component with an outlined border. The minus/plus icons are Primary Terracotta.
- **Loading Skeletons:** Use a soft pulse animation moving from Warm Gray to a slightly lighter tint.