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
.auth-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 32px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.auth-card__header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
}

.auth-card__title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.auth-card__subtitle {
  margin: 4px 0 24px;
  color: #6b7280;
  font-size: 14px;
}

.auth-card__submit {
  width: 100%;
  margin-top: 8px;
}

.auth-card__footer {
  margin: 20px 0 0;
  text-align: center;
  font-size: 14px;
  color: #6b7280;
}
</style>
