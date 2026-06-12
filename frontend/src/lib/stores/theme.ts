import { browser } from '$app/environment';
import { writable } from 'svelte/store';

function initialDark(): boolean {
  if (!browser) return false;
  return document.documentElement.classList.contains('dark');
}

export const darkMode = writable<boolean>(initialDark());

export function toggleTheme() {
  darkMode.update((dark) => {
    const next = !dark;
    if (browser) {
      document.documentElement.classList.toggle('dark', next);
      localStorage.setItem('bb_theme', next ? 'dark' : 'light');
    }
    return next;
  });
}
