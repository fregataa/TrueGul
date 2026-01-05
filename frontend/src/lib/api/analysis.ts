import { apiClient } from "./client";

export type AnalysisStatus = "pending" | "processing" | "completed" | "failed";

export interface Analysis {
  id: string;
  writing_id: string;
  status: AnalysisStatus;
  ai_score?: number;
  feedback?: string;
  error_code?: string;
  error_message?: string;
  latency_ms?: number;
  created_at: string;
  updated_at: string;
}

export interface SubmitResponse {
  message: string;
  analysis_id: string;
}

export const analysisApi = {
  submit: async (writingId: string): Promise<SubmitResponse> => {
    return apiClient.post<SubmitResponse>(`/writings/${writingId}/submit`);
  },

  getAnalysis: async (writingId: string): Promise<Analysis> => {
    return apiClient.get<Analysis>(`/writings/${writingId}/analysis`);
  },
};
