import { apiClient } from "./client";

export interface User {
  id: string;
  email: string;
}

export interface AuthResponse {
  user: User;
  csrf_token: string;
}

export interface SignupRequest {
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export const authApi = {
  signup: async (data: SignupRequest): Promise<User> => {
    return apiClient.post<User>("/auth/signup", data);
  },

  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/login", data);
    apiClient.setCSRFToken(response.csrf_token);
    return response;
  },

  logout: async (): Promise<void> => {
    await apiClient.post("/auth/logout");
    apiClient.setCSRFToken("");
  },

  me: async (): Promise<User> => {
    return apiClient.get<User>("/auth/me");
  },
};
