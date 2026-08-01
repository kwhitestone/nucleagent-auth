<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
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

const avatarLetter = computed(() => {
  const name = userStore.displayName || userStore.user?.username || "U";
  return name.charAt(0).toUpperCase();
});

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
      <div class="home__brand">
        <span class="home__brand-mark">N</span>
        <span class="home__title">{{ t("home.title") }}</span>
      </div>
      <div class="home__header-actions">
        <LanguageSwitcher />
        <el-button class="home__logout" @click="handleLogout">
          {{ t("home.logout") }}
        </el-button>
      </div>
    </header>

    <main class="home__body">
      <section class="home__card home__card--user">
        <div class="home__avatar">{{ avatarLetter }}</div>
        <div class="home__user-main">
          <h2 class="home__user-name">{{ userStore.displayName || "-" }}</h2>
          <p class="home__user-meta">@{{ userStore.user?.username || "-" }}</p>
          <div v-if="userStore.roles.length" class="home__roles">
            <el-tag
              v-for="role in userStore.roles"
              :key="role"
              size="small"
              class="home__role-tag"
            >
              {{ role }}
            </el-tag>
          </div>
        </div>
      </section>

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
  position: relative;
  min-height: 100vh;
  background-color: var(--na-bg);
}

.home::before {
  content: "";
  position: fixed;
  inset: 0;
  background: var(--na-grad-mesh);
  pointer-events: none;
  z-index: 0;
}

.home__header {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 32px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--na-border);
}

.home__brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.home__brand-mark {
  display: flex;
  width: 32px;
  height: 32px;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  background: var(--na-grad-brand);
  color: #fff;
  font-family: var(--na-font-display);
  font-size: 18px;
  box-shadow: var(--na-shadow-glow-teal);
}

.home__title {
  font-family: var(--na-font-display);
  font-size: 20px;
  color: var(--na-text);
}

.home__header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.home__logout {
  border: 1px solid var(--na-border);
  background: #fff;
  color: var(--na-text-secondary);
  border-radius: var(--na-r-md);
}

.home__logout:hover {
  border-color: var(--rose-400, #fb7185);
  color: #e11d48;
  background: #fff;
}

.home__body {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 880px;
  margin: 0 auto;
  padding: 32px 24px 64px;
}

.home__card {
  padding: 28px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(16px);
  border: 1px solid var(--na-border);
  border-radius: var(--na-r-xl);
  box-shadow: var(--na-shadow-md);
}

.home__card--user {
  display: flex;
  align-items: center;
  gap: 20px;
}

.home__avatar {
  display: flex;
  width: 64px;
  height: 64px;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 18px;
  background: var(--na-grad-teal-indigo);
  color: #fff;
  font-family: var(--na-font-display);
  font-size: 30px;
  box-shadow: var(--na-shadow-glow-teal);
}

.home__user-name {
  margin: 0;
  font-family: var(--na-font-display);
  font-size: 24px;
  font-weight: 400;
  color: var(--na-text);
}

.home__user-meta {
  margin: 2px 0 0;
  color: var(--na-text-tertiary);
  font-size: 13px;
  font-family: var(--na-font-mono);
}

.home__roles {
  margin-top: 8px;
  display: flex;
  gap: 6px;
}

.home__card-title {
  margin: 0 0 18px;
  font-size: 15px;
  font-weight: 600;
  color: var(--na-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.home__info {
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin: 0;
}

.home__info-row {
  display: flex;
  align-items: center;
}

.home__info-row dt {
  width: 120px;
  margin: 0;
  color: var(--na-text-tertiary);
  font-size: 13px;
}

.home__info-row dd {
  margin: 0;
  font-size: 14px;
  color: var(--na-text);
}

.home__role-tag {
  margin-right: 8px;
}
</style>
