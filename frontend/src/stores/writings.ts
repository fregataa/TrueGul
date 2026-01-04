"use client";

import { create } from "zustand";
import { ApiClientError } from "@/lib/api/client";
import {
  type UpdateWritingRequest,
  type Writing,
  type WritingType,
  writingsApi,
} from "@/lib/api/writings";

interface WritingsState {
  writings: Writing[];
  currentWriting: Writing | null;
  total: number;
  page: number;
  limit: number;
  totalPages: number;
  isLoading: boolean;
  error: string | null;

  fetchWritings: (page?: number, limit?: number) => Promise<void>;
  fetchWriting: (id: string) => Promise<void>;
  createWriting: (type: WritingType, title: string, content: string) => Promise<Writing | null>;
  updateWriting: (id: string, data: UpdateWritingRequest) => Promise<Writing | null>;
  deleteWriting: (id: string) => Promise<boolean>;
  setCurrentWriting: (writing: Writing | null) => void;
  clearError: () => void;
}

export const useWritingsStore = create<WritingsState>((set) => ({
  writings: [],
  currentWriting: null,
  total: 0,
  page: 1,
  limit: 10,
  totalPages: 0,
  isLoading: false,
  error: null,

  fetchWritings: async (page = 1, limit = 10) => {
    set({ isLoading: true, error: null });
    try {
      const response = await writingsApi.list(page, limit);
      set({
        writings: response.writings,
        total: response.total,
        page: response.page,
        limit: response.limit,
        totalPages: response.total_pages,
        isLoading: false,
      });
    } catch (error) {
      const message = error instanceof ApiClientError ? error.message : "Failed to fetch writings";
      set({ isLoading: false, error: message });
    }
  },

  fetchWriting: async (id: string) => {
    set({ isLoading: true, error: null });
    try {
      const writing = await writingsApi.getById(id);
      set({ currentWriting: writing, isLoading: false });
    } catch (error) {
      const message = error instanceof ApiClientError ? error.message : "Failed to fetch writing";
      set({ isLoading: false, error: message, currentWriting: null });
    }
  },

  createWriting: async (type: WritingType, title: string, content: string) => {
    set({ isLoading: true, error: null });
    try {
      const writing = await writingsApi.create({ type, title, content });
      set((state) => ({
        writings: [writing, ...state.writings],
        currentWriting: writing,
        isLoading: false,
      }));
      return writing;
    } catch (error) {
      const message = error instanceof ApiClientError ? error.message : "Failed to create writing";
      set({ isLoading: false, error: message });
      return null;
    }
  },

  updateWriting: async (id: string, data: UpdateWritingRequest) => {
    set({ isLoading: true, error: null });
    try {
      const writing = await writingsApi.update(id, data);
      set((state) => ({
        writings: state.writings.map((w) => (w.id === id ? writing : w)),
        currentWriting: state.currentWriting?.id === id ? writing : state.currentWriting,
        isLoading: false,
      }));
      return writing;
    } catch (error) {
      const message = error instanceof ApiClientError ? error.message : "Failed to update writing";
      set({ isLoading: false, error: message });
      return null;
    }
  },

  deleteWriting: async (id: string) => {
    set({ isLoading: true, error: null });
    try {
      await writingsApi.delete(id);
      set((state) => ({
        writings: state.writings.filter((w) => w.id !== id),
        currentWriting: state.currentWriting?.id === id ? null : state.currentWriting,
        isLoading: false,
      }));
      return true;
    } catch (error) {
      const message = error instanceof ApiClientError ? error.message : "Failed to delete writing";
      set({ isLoading: false, error: message });
      return false;
    }
  },

  setCurrentWriting: (writing: Writing | null) => {
    set({ currentWriting: writing });
  },

  clearError: () => set({ error: null }),
}));
