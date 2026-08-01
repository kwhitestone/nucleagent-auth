import { defineStore } from "pinia";
import { computed, ref } from "vue";
import {
  login as loginApi,
  fetchUserInfo,
} from "@/api/auth";
import type { LoginRequest } from "@/api/types";
import type { UserInfo } from "@/api/types";
import {
  clearTokens,
  getAccessToken,
  setTokens,
} from "@/utils/token";

/**
 * User store.
 *
 * State is mutated only through actions (immutable-style: every action
 * reassigns the refs rather than mutating nested fields). The access token is
 * persisted to localStorage via the token utils so it survives reloads; the
 * `token` ref is hydrated from localStorage on store creation.
 */
export const useUserStore = defineStore("user", () => {
  const token = ref<string>(getAccessToken());
  const user = ref<UserInfo | null>(null);

  const isAuthenticated = computed(() => !!token.value);
  const displayName = computed(
    () => user.value?.nickName || user.value?.username || "",
  );
  const roleId = computed(() => user.value?.roleId ?? null);
  const roles = computed(() => user.value?.roles ?? []);

  function setToken(accessToken: string, refreshToken: string): void {
    setTokens(accessToken, refreshToken);
    token.value = accessToken;
  }

  async function login(payload: LoginRequest): Promise<UserInfo> {
    const data = await loginApi(payload);
    setToken(data.accessToken, data.refreshToken);
    // The login response already carries the user; mirror it into state so the
    // UI can render immediately without a second round-trip.
    user.value = data.user;
    return data.user;
  }

  async function fetchUser(): Promise<UserInfo> {
    const info = await fetchUserInfo();
    user.value = info;
    return info;
  }

  function logout(): void {
    user.value = null;
    token.value = "";
    clearTokens();
  }

  return {
    // state
    token,
    user,
    // getters
    isAuthenticated,
    displayName,
    roleId,
    roles,
    // actions
    setToken,
    login,
    fetchUser,
    logout,
  };
});
