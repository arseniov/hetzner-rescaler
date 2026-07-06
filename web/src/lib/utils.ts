import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * Concatenate class names with conflict-resolution. shadcn-svelte
 * convention: lets callers layer conditional classes (e.g.
 * `cn('p-4', isActive && 'bg-accent')`) without producing duplicate /
 * conflicting Tailwind utilities.
 */
export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}
