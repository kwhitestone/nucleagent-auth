import http from "./http";
import type {
  ApiKey,
  ApiKeyWithSecret,
  CreateApiKeyRequest,
  LoginRequest,
  RegisterRequest,
  TokenData,
  UserInfo,
} from "./types";

const BASE = "/api/v1/addons/auth";

/**
 * POST /register
 * Body: { username, password, nickName } -> { code, message }
 */
export async function register(payload: RegisterRequest): Promise<void> {
  await http.post(`${BASE}/register`, payload);
}

/**
 * POST /login
 * Body: { username, password } -> { code, message, data: TokenData }
 */
export async function login(payload: LoginRequest): Promise<TokenData> {
  const response = await http.post<TokenData>(`${BASE}/login`, payload);
  return response.data;
}

/**
 * POST /refresh-token
 * Body: { refreshToken } -> same shape as login
 */
export async function refreshToken(refreshToken: string): Promise<TokenData> {
  const response = await http.post<TokenData>(`${BASE}/refresh-token`, {
    refreshToken,
  });
  return response.data;
}

/**
 * GET /user-info
 * Requires Authorization: Bearer <token>.
 */
export async function fetchUserInfo(): Promise<UserInfo> {
  const response = await http.get<UserInfo>(`${BASE}/user-info`);
  return response.data;
}

/**
 * POST /api-keys
 * Requires Authorization. Plaintext is returned only once.
 */
export async function createApiKey(
  payload: CreateApiKeyRequest,
): Promise<ApiKeyWithSecret> {
  const response = await http.post<ApiKeyWithSecret>(`${BASE}/api-keys`, payload);
  return response.data;
}

/**
 * GET /api-keys
 * Requires Authorization.
 */
export async function listApiKeys(): Promise<ApiKey[]> {
  const response = await http.get<ApiKey[]>(`${BASE}/api-keys`);
  return response.data;
}

/**
 * DELETE /api-keys/:id
 * Requires Authorization.
 */
export async function deleteApiKey(id: number | string): Promise<void> {
  await http.delete(`${BASE}/api-keys/${id}`);
}
