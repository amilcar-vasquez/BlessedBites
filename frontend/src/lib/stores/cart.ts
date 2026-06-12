import { browser } from '$app/environment';
import { derived, writable } from 'svelte/store';

export type CartItem = {
  id: number;
  name: string;
  price: number;
  qty: number;
  image_url?: string;
};

const CART_KEY = 'bb_cart';

function readCart(): CartItem[] {
  if (!browser) return [];
  try {
    const raw = localStorage.getItem(CART_KEY);
    return raw ? (JSON.parse(raw) as CartItem[]) : [];
  } catch {
    return [];
  }
}

export const cart = writable<CartItem[]>(readCart());

if (browser) {
  cart.subscribe((items) => {
    localStorage.setItem(CART_KEY, JSON.stringify(items));
  });
}

export const cartCount = derived(cart, ($cart) => $cart.reduce((n, i) => n + i.qty, 0));

export const cartTotal = derived(cart, ($cart) =>
  $cart.reduce((sum, i) => sum + i.price * i.qty, 0)
);

export function addToCart(item: { id: number; name: string; price: number; image_url?: string }) {
  cart.update((items) => {
    const existing = items.find((i) => i.id === item.id);
    if (existing) {
      return items.map((i) => (i.id === item.id ? { ...i, qty: i.qty + 1 } : i));
    }
    return [...items, { ...item, qty: 1 }];
  });
}

export function setQty(id: number, qty: number) {
  cart.update((items) => {
    if (qty <= 0) return items.filter((i) => i.id !== id);
    return items.map((i) => (i.id === id ? { ...i, qty } : i));
  });
}

export function removeFromCart(id: number) {
  cart.update((items) => items.filter((i) => i.id !== id));
}

export function clearCart() {
  cart.set([]);
}
