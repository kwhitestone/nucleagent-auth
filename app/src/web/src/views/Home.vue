<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { ElMessage } from "element-plus";
import { useUserStore } from "@/store/user";
import LanguageSwitcher from "@/components/LanguageSwitcher.vue";
import ApiKeyPanel from "@/components/ApiKeyPanel.vue";

const { t } = useI18n();
const router = useRouter();
const userStore = useUserStore();

const loading = ref(false);

async function loadUser(): Promise<void> {
  loading.value = true;
  try {
    await userStore.fetchUser();
  } catch {
    // The http interceptor already handles 401 -> redirect to login.
    // Any other failure is surfaced but does not block the page shell.
  } finally {
    loading.value = false;
  }
}

async function handleLogout(): Promise<void> {
  try {
    userStore.logout();
    ElMessage.success(t("home.logoutSuccess"));
    router.push("/login");
  } catch {
    ElMessage.error(t("home.logoutFailed"));
  }
}

onMounted(loadUser);
</script>

<template>
  <div v-loading="loading" class="home">
    <header class="home__header">
      <span class="home__title">{{ t("home.title") }}</span>
      <div class="home__header-actions">
        <LanguageSwitcher />
        <el-button type="danger" @click="handleLogout">
          {{ t("home.logout") }}
        </el-button>
      </div>
    </header>

    <main class="home__body">
      <section class="home__card">
        <h2 class="home__card-title">{{ t("home.userInfoTitle") }}</h2>
        <dl class="home__info">
          <div class="home__info-row">
            <dt>{{ t("common.nickname") }}</dt>
            <dd>{{ userStore.displayName || "-" }}</dd>
          </div>
          <div class="home__info-row">
            <dt>{{ t("common.username") }}</dt>
            <dd>{{ userStore.user?.username || "-" }}</dd>
          </div>
          <div class="home__info-row">
            <dt>{{ t("home.roleId") }}</dt>
            <dd>{{ userStore.roleId ?? "-" }}</dd>
          </div>
          <div v-if="userStore.roles.length" class="home__info-row">
            <dt>{{ t("apiKey.enabled") }}</dt>
            <dd>
              <el-tag
                v-for="role in userStore.roles"
                :key="role"
                size="small"
                class="home__role-tag"
              >
                {{ role }}
              </el-tag>
            </dd>
          </div>
        </dl>
      </section>

      <section class="home__card">
        <ApiKeyPanel />
      </section>
    </main>
  </div>
</template>

<style scoped>
.home {
  min-height: 100vh;
  background: #f5f6f8;
}

.home__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 32px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.home__title {
  font-size: 18px;
  font-weight: 600;
}

.home__header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.home__body {
  display: flex;
  flex-direction: column;
  gap: 24px;
  max-width: 960px;
  margin: 0 auto;
  padding: 32px 24px;
}

.home__card {
  padding: 24px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.home__card-title {
  margin: 0 0 16px;
  font-size: 18px;
  font-weight: 600;
}

.home__info {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 0;
}

.home__info-row {
  display: flex;
  align-items: center;
}

.home__info-row dt {
  width: 120px;
  margin: 0;
  color: #6b7280;
  font-size: 14px;
}

.home__info-row dd {
  margin: 0;
  font-size: 14px;
}

.home__role-tag {
  margin-right: 8px;
}
</style>
