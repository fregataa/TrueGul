"use client";

import { create } from "zustand";
import { type Analysis, analysisApi } from "@/lib/api/analysis";
import { ApiClientError } from "@/lib/api/client";

const POLLING_CONFIG = {
  initialInterval: 2000,
  maxInterval: 10000,
  backoffMultiplier: 1.5,
  maxAttempts: 60,
};

interface AnalysisState {
  currentAnalysis: Analysis | null;
  isSubmitting: boolean;
  isPolling: boolean;
  error: string | null;
  errorCode: string | null;
  isRateLimited: boolean;

  submitWriting: (writingId: string) => Promise<boolean>;
  fetchAnalysis: (writingId: string) => Promise<void>;
  startPolling: (writingId: string) => void;
  stopPolling: () => void;
  clearAnalysis: () => void;
  clearError: () => void;
}

let pollingTimeoutId: ReturnType<typeof setTimeout> | null = null;

export const useAnalysisStore = create<AnalysisState>((set, get) => ({
  currentAnalysis: null,
  isSubmitting: false,
  isPolling: false,
  error: null,
  errorCode: null,
  isRateLimited: false,

  submitWriting: async (writingId: string) => {
    set({ isSubmitting: true, error: null, errorCode: null, isRateLimited: false });
    try {
      await analysisApi.submit(writingId);
      set({ isSubmitting: false });
      get().startPolling(writingId);
      return true;
    } catch (error) {
      if (error instanceof ApiClientError) {
        if (error.status === 429 || error.errorCode === "FORBIDDEN") {
          set({
            isSubmitting: false,
            isRateLimited: true,
            error: error.message,
            errorCode: error.errorCode,
          });
        } else {
          set({
            isSubmitting: false,
            error: error.message,
            errorCode: error.errorCode,
          });
        }
      } else {
        set({
          isSubmitting: false,
          error: "제출 중 오류가 발생했습니다.",
          errorCode: null,
        });
      }
      return false;
    }
  },

  fetchAnalysis: async (writingId: string) => {
    try {
      const analysis = await analysisApi.getAnalysis(writingId);
      set({ currentAnalysis: analysis, error: null, errorCode: null });
    } catch (error) {
      if (error instanceof ApiClientError) {
        if (error.status !== 404) {
          set({ error: error.message, errorCode: error.errorCode });
        }
      }
    }
  },

  startPolling: (writingId: string) => {
    get().stopPolling();
    set({ isPolling: true });

    let attempts = 0;
    let interval = POLLING_CONFIG.initialInterval;

    const poll = async () => {
      attempts++;

      if (attempts > POLLING_CONFIG.maxAttempts) {
        get().stopPolling();
        set({ error: "분석이 예상보다 오래 걸리고 있습니다. 나중에 다시 확인해주세요." });
        return;
      }

      try {
        const analysis = await analysisApi.getAnalysis(writingId);
        set({ currentAnalysis: analysis });

        if (analysis.status === "completed" || analysis.status === "failed") {
          get().stopPolling();
          return;
        }

        interval = Math.min(
          interval * POLLING_CONFIG.backoffMultiplier,
          POLLING_CONFIG.maxInterval
        );
        pollingTimeoutId = setTimeout(poll, interval);
      } catch (error) {
        if (error instanceof ApiClientError && error.status === 401) {
          get().stopPolling();
          return;
        }
        pollingTimeoutId = setTimeout(poll, interval);
      }
    };

    poll();
  },

  stopPolling: () => {
    if (pollingTimeoutId) {
      clearTimeout(pollingTimeoutId);
      pollingTimeoutId = null;
    }
    set({ isPolling: false });
  },

  clearAnalysis: () => {
    get().stopPolling();
    set({
      currentAnalysis: null,
      error: null,
      errorCode: null,
      isRateLimited: false,
    });
  },

  clearError: () => set({ error: null, errorCode: null }),
}));
