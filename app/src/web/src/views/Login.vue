<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { useUserStore } from "@/store/user";
import LanguageSwitcher from "@/components/LanguageSwitcher.vue";

const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const formRef = ref<FormInstance>();
const loading = ref(false);
const form = reactive({
  username: "",
  password: "",
});

const rules: FormRules<typeof form> = {
  username: [{ required: true, message: t("login.enterUsername"), trigger: "blur" }],
  password: [{ required: true, message: t("login.enterPassword"), trigger: "blur" }],
};

function redirectTarget(): string {
  const redirect = route.query.redirect;
  return typeof redirect === "string" && redirect.startsWith("/") ? redirect : "/";
}

async function handleSubmit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  loading.value = true;
  try {
    await userStore.login({
      username: form.username.trim(),
      password: form.password,
    });
    ElMessage.success(t("login.success"));
    router.push(redirectTarget());
  } catch (error) {
    const message = error instanceof Error ? error.message : t("login.failed");
    ElMessage.error(message);
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-card__header">
        <LanguageSwitcher />
      </div>
      <h1 class="auth-card__title">{{ t("login.title") }}</h1>
      <p class="auth-card__subtitle">{{ t("login.subtitle") }}</p>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        @submit.prevent="handleSubmit"
      >
        <el-form-item :label="t('common.username')" prop="username">
          <el-input
            v-model="form.username"
            :placeholder="t('login.usernamePlaceholder')"
            autocomplete="username"
          />
        </el-form-item>
        <el-form-item :label="t('common.password')" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            :placeholder="t('login.passwordPlaceholder')"
            show-password
            autocomplete="current-password"
            @keyup.enter="handleSubmit"
          />
        </el-form-item>
        <el-button
          type="primary"
          class="auth-card__submit"
          :loading="loading"
          @click="handleSubmit"
        >
          {{ loading ? t("login.submitting") : t("login.submit") }}
        </el-button>
      </el-form>

      <p class="auth-card__footer">
        {{ t("login.noAccount") }}
        <router-link to="/register">{{ t("login.registerLink") }}</router-link>
      </p>
    </div>
  </div>
</template>

<style scoped>
/* auth 页样式复用全局 .auth-page / .auth-card（见 styles/global.css）*/
</style>
