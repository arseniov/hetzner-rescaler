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

/**
 * Classes that style a form input trigger the same way regardless of
 * whether it's a bits-ui Select (base/top type pickers) or a native
 * <select> (fallback-chain add picker). Keeping these in one string
 * guarantees the two surfaces look identical: same border, same height,
 * same focus ring, same disabled treatment. Components that wrap a
 * className prop should append via `cn(SELECT_TRIGGER_CLASSES, className)`.
 */
export const SELECT_TRIGGER_CLASSES =
  'flex h-9 w-full items-center justify-between rounded-md border border-border bg-input px-3 py-1 text-sm text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-background disabled:cursor-not-allowed disabled:opacity-50';

