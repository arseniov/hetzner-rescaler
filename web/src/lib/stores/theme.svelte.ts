import { browser } from '$app/environment';

export const STORAGE_KEY = 'hetzner-rescaler-theme';
export type Theme = 'light' | 'dark';

function readInitial(): Theme {
  if (!browser) return 'light';
  const stored = localStorage.getItem(STORAGE_KEY) as Theme | null;
  if (stored === 'light' || stored === 'dark') return stored;
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

class ThemeStore {
  current = $state<Theme>(readInitial());

  constructor() {
    if (browser) {
      $effect.root(() => {
        $effect(() => {
          const root = document.documentElement;
          root.classList.toggle('dark', this.current === 'dark');
          if (this.current) localStorage.setItem(STORAGE_KEY, this.current);
        });
      });
    }
  }

  toggle() {
    this.current = this.current === 'dark' ? 'light' : 'dark';
  }

  set(theme: Theme) {
    this.current = theme;
  }
}

export const theme = new ThemeStore();