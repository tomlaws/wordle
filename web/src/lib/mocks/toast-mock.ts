// toast-mock.ts
import type { ToastAPI } from '$lib/context/toast-context';
import { vi } from 'vitest';

// Typed spy functions
export const toastSpies: ToastAPI = {
  success: vi.fn<(msg: string) => void>(),
  error: vi.fn<(msg: string) => void>(),
  info: vi.fn<(msg: string) => void>(),
  warning: vi.fn<(msg: string) => void>()
};

export const mockToast: ToastAPI = toastSpies;
