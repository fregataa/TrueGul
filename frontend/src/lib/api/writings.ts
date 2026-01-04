import { apiClient } from "./client";

export type WritingType = "essay" | "cover_letter";
export type WritingStatus = "draft" | "submitted" | "analyzed";

export interface Writing {
  id: string;
  user_id: string;
  type: WritingType;
  title: string;
  content: string;
  status: WritingStatus;
  created_at: string;
  updated_at: string;
  submitted_at?: string;
}

export interface WritingListResponse {
  writings: Writing[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface CreateWritingRequest {
  type: WritingType;
  title: string;
  content: string;
}

export interface UpdateWritingRequest {
  type?: WritingType;
  title?: string;
  content?: string;
}

export const writingsApi = {
  list: async (page = 1, limit = 10): Promise<WritingListResponse> => {
    return apiClient.get<WritingListResponse>(`/writings?page=${page}&limit=${limit}`);
  },

  getById: async (id: string): Promise<Writing> => {
    return apiClient.get<Writing>(`/writings/${id}`);
  },

  create: async (data: CreateWritingRequest): Promise<Writing> => {
    return apiClient.post<Writing>("/writings", data);
  },

  update: async (id: string, data: UpdateWritingRequest): Promise<Writing> => {
    return apiClient.request<Writing>(`/writings/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.request(`/writings/${id}`, {
      method: "DELETE",
    });
  },
};
