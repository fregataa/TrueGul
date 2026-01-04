"use client";

import { create } from "zustand";
import { authApi, type User } from "@/lib/api/auth";
import { ApiClientError } from "@/lib/api/client";

interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  error: string | null;

  signup: (email: string, password: string) => Promise<boolean>;
  login: (email: string, password: string) => Promise<boolean>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isLoading: true,
  isAuthenticated: false,
  error: null,

  signup: async (email: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      await authApi.signup({ email, password });
      set({ isLoading: false });
      return true;
    } catch (error) {
      const message = error instanceof ApiClientError ? error.message : "Failed to sign up";
      set({ isLoading: false, error: message });
      return false;
    }
  },

  login: async (email: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await authApi.login({ email, password });
      set({
        user: response.user,
        isAuthenticated: true,
        isLoading: false,
      });
      return true;
    } catch (error) {
      const message = error instanceof ApiClientError ? error.message : "Failed to log in";
      set({ isLoading: false, error: message });
      return false;
    }
  },

  logout: async () => {
    try {
      await authApi.logout();
    } catch {
      // Ignore logout errors
    }
    set({ user: null, isAuthenticated: false });
  },

  checkAuth: async () => {
    set({ isLoading: true });
    try {
      const user = await authApi.me();
      set({ user, isAuthenticated: true, isLoading: false });
    } catch {
      set({ user: null, isAuthenticated: false, isLoading: false });
    }
  },

  clearError: () => set({ error: null }),
}));
