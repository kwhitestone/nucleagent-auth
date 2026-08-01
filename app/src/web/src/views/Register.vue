<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { register } from "@/api/auth";
import LanguageSwitcher from "@/components/LanguageSwitcher.vue";

const { t } = useI18n();
const router = useRouter();

const formRef = ref<FormInstance>();
const loading = ref(false);
const form = reactive({
  username: "",
  password: "",
  nickName: "",
});

const rules: FormRules<typeof form> = {
  username: [{ required: true, message: t("register.enterUsername"), trigger: "blur" }],
  password: [{ required: true, message: t("register.enterPassword"), trigger: "blur" }],
  nickName: [{ required: true, message: t("register.enterNickname"), trigger: "blur" }],
};

async function handleSubmit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  loading.value = true;
  try {
    await register({
      username: form.username.trim(),
      password: form.password,
      nickName: form.nickName.trim(),
    });
    ElMessage.success(t("register.success"));
    router.push("/login");
  } catch (error) {
    const message = error instanceof Error ? error.message : t("register.failed");
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
      <h1 class="auth-card__title">{{ t("register.title") }}</h1>
      <p class="auth-card__subtitle">{{ t("register.subtitle") }}</p>

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
            :placeholder="t('register.usernamePlaceholder')"
            autocomplete="username"
          />
        </el-form-item>
        <el-form-item :label="t('common.nickname')" prop="nickName">
          <el-input
            v-model="form.nickName"
            :placeholder="t('register.nicknamePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('common.password')" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            :placeholder="t('register.passwordPlaceholder')"
            show-password
            autocomplete="new-password"
            @keyup.enter="handleSubmit"
          />
        </el-form-item>
        <el-button
          type="primary"
          class="auth-card__submit"
          :loading="loading"
          @click="handleSubmit"
        >
          {{ loading ? t("register.submitting") : t("register.submit") }}
        </el-button>
      </el-form>

      <p class="auth-card__footer">
        {{ t("register.hasAccount") }}
        <router-link to="/login">{{ t("register.loginLink") }}</router-link>
      </p>
    </div>
  </div>
</template>

<style scoped>
/* auth 页样式复用全局 .auth-page / .auth-card（见 styles/global.css）*/
</style>
