const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

interface ApiError {
  error_code: string;
  message: string;
}

class ApiClient {
  private csrfToken: string | null = null;

  setCSRFToken(token: string) {
    this.csrfToken = token;
  }

  getCSRFToken(): string | null {
    return this.csrfToken;
  }

  async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    if (this.csrfToken && options.method && options.method !== "GET") {
      (headers as Record<string, string>)["X-CSRF-Token"] = this.csrfToken;
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
      credentials: "include",
    });

    if (!response.ok) {
      const error: ApiError = await response.json();
      throw new ApiClientError(error.error_code, error.message, response.status);
    }

    return response.json();
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "GET" });
  }

  async post<T>(endpoint: string, data?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      body: data ? JSON.stringify(data) : undefined,
    });
  }
}

export class ApiClientError extends Error {
  constructor(
    public errorCode: string,
    message: string,
    public status: number
  ) {
    super(message);
    this.name = "ApiClientError";
  }
}

export const apiClient = new ApiClient();
