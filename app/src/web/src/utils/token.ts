/**
 * Token persistence helpers.
 *
 * The backend issues a JWT (accessToken) plus a long-lived refreshToken.
 * We persist both to localStorage so the page survives reloads, and the
 * refresh flow can recover an expired access token without forcing a
 * re-login.
 */

const ACCESS_TOKEN_KEY = "nucleagent_access_token";
const REFRESH_TOKEN_KEY = "nucleagent_refresh_token";

export function getAccessToken(): string {
  return localStorage.getItem(ACCESS_TOKEN_KEY) ?? "";
}

export function getRefreshToken(): string {
  return localStorage.getItem(REFRESH_TOKEN_KEY) ?? "";
}

export function setTokens(accessToken: string, refreshToken: string): void {
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
}

export function clearTokens(): void {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
}
