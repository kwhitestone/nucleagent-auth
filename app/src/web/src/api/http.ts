import axios from "axios";
import type { AxiosResponse, InternalAxiosRequestConfig } from "axios";
import { getAccessToken, clearTokens } from "@/utils/token";
import type { ApiEnvelope } from "./types";

/**
 * Shared axios instance.
 *
 * - baseURL is empty: requests use relative `/api/...` URLs which the Vite dev
 *   server proxies to the prism-fusion backend on :6670. In production the
 *   reverse proxy handles the same path.
 * - The request interceptor attaches the JWT from localStorage.
 * - The response interceptor unwraps the unified `{ code, message, data }`
 *   envelope: code === 0 means success and we resolve with `data`; any other
 *   code is rejected with an `ApiError`. HTTP 401 clears tokens and bounces
 *   the user back to /login.
 */
const http = axios.create({
  baseURL: "",
  timeout: 15000,
  headers: {
    "Content-Type": "application/json",
  },
});

/** Error thrown when the envelope reports a business failure (code !== 0). */
export class ApiError extends Error {
  code: number;
  status: number;

  constructor(message: string, code: number, status: number) {
    super(message);
    this.name = "ApiError";
    this.code = code;
    this.status = status;
  }
}

http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

function redirectToLogin(): void {
  clearTokens();
  // 在 micro-app 子应用模式下，不做 window.location 硬跳转（会劫持整个壳的 URL）。
  // 只清 token，让壳应用自行决定何时切到登录页。
  const w = globalThis as Record<string, unknown>;
  if (w.__MICRO_APP_ENVIRONMENT__) return;
  // 独立运行时跳 /login。
  if (typeof window !== "undefined" && !window.location.pathname.startsWith("/login")) {
    window.location.href = "/login";
  }
}

http.interceptors.response.use(
  (response: AxiosResponse<ApiEnvelope<unknown>>) => {
    const envelope = response.data;
    // Some endpoints (e.g. raw passthrough) may not follow the envelope;
    // treat a missing `code` as a pass-through success.
    if (envelope === null || typeof envelope !== "object" || !("code" in envelope)) {
      return response;
    }

    if (envelope.code === 0) {
      // Resolve with the unwrapped payload so callers deal with `data` directly.
      return { ...response, data: envelope.data };
    }

    // Business failure: reject so callers can `.catch()` it.
    throw new ApiError(
      envelope.message || "Request failed",
      envelope.code,
      response.status,
    );
  },
  (error: unknown) => {
    // Network or HTTP-level error.
    if (axios.isAxiosError(error)) {
      const status = error.response?.status ?? 0;
      if (status === 401) {
        redirectToLogin();
      }
      const envelope = error.response?.data as ApiEnvelope<unknown> | undefined;
      const message =
        envelope?.message || error.message || "Network error";
      return Promise.reject(new ApiError(message, envelope?.code ?? -1, status));
    }
    return Promise.reject(error);
  },
);

export default http;
