<script setup lang="ts">
/**
 * 登录页 —— 独立运行模式（非 micro-app）的入口。
 *
 * 在主壳里登录由壳的 LoginModal 处理（壳直连 auth 后端），本页只在 auth 子应用
 * 独立 dev 时使用。已移除 Element Plus：用原生表单 + useToast，校验改为简单的
 * 必填检查（原 EP 规则只有 required）。
 */
import { reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { useUserStore } from "@/store/user";
import { toast } from "@/composables/useToast";

const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const loading = ref(false);
const errorMsg = ref("");
const form = reactive({ username: "", password: "" });

function redirectTarget(): string {
  const redirect = route.query.redirect;
  return typeof redirect === "string" && redirect.startsWith("/") ? redirect : "/";
}

async function handleSubmit(): Promise<void> {
  if (loading.value) return;
  if (!form.username.trim() || !form.password) {
    errorMsg.value = t("login.enterUsername");
    return;
  }
  loading.value = true;
  errorMsg.value = "";
  try {
    await userStore.login({ username: form.username.trim(), password: form.password });
    toast.success(t("login.success"));
    router.push(redirectTarget());
  } catch (error) {
    errorMsg.value = error instanceof Error ? error.message : t("login.failed");
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1 class="auth-card__title">{{ t("login.title") }}</h1>
      <p class="auth-card__subtitle">{{ t("login.subtitle") }}</p>

      <form class="auth-form" @submit.prevent="handleSubmit">
        <label class="auth-field">
          <span class="auth-field__label">{{ t("common.username") }}</span>
          <input
            v-model="form.username"
            type="text"
            class="auth-field__input"
            :placeholder="t('login.usernamePlaceholder')"
            autocomplete="username"
          />
        </label>
        <label class="auth-field">
          <span class="auth-field__label">{{ t("common.password") }}</span>
          <input
            v-model="form.password"
            type="password"
            class="auth-field__input"
            :placeholder="t('login.passwordPlaceholder')"
            autocomplete="current-password"
            @keyup.enter="handleSubmit"
          />
        </label>

        <p v-if="errorMsg" class="auth-error">{{ errorMsg }}</p>

        <button type="submit" class="auth-card__submit" :disabled="loading">
          {{ loading ? t("common.sending") : t("login.title") }}
        </button>
      </form>

      <p class="auth-card__footer">
        {{ t("register.hasAccount") }}
        <router-link to="/register">{{ t("register.loginLink") }}</router-link>
      </p>
    </div>
  </div>
</template>

<style scoped>
.auth-form { display: flex; flex-direction: column; gap: 14px; }
.auth-field { display: flex; flex-direction: column; gap: 6px; }
.auth-field__label { font-size: 13px; font-weight: 600; color: var(--text-primary); }
.auth-field__input {
  width: 100%; padding: 10px 14px;
  border: 1.5px solid var(--border); border-radius: var(--r-md);
  font-size: 14px; color: var(--text-primary); background: var(--bg-card);
  outline: none; transition: all 0.2s var(--ease);
}
.auth-field__input:focus { border-color: var(--teal-400); box-shadow: 0 0 0 4px rgba(20, 184, 166, 0.08); }
.auth-error { color: var(--rose-500); font-size: 12.5px; margin: 0; }
.auth-card__submit {
  margin-top: 4px; padding: 12px; border: none; border-radius: var(--r-md);
  background: var(--grad-teal-indigo); background-size: 200% 200%;
  animation: gradient-flow 5s var(--ease) infinite;
  color: white; font-weight: 600; cursor: pointer;
  box-shadow: var(--shadow-glow-teal); transition: all 0.2s var(--ease);
}
.auth-card__submit:hover:not(:disabled) { transform: translateY(-1px); box-shadow: var(--shadow-glow-indigo); }
.auth-card__submit:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
