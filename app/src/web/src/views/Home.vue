<script setup lang="ts">
/**
 * 账户首页 —— micro-app 模式下 auth 子应用的唯一视图。
 *
 * 去掉自带头部（chrome 已在壳），仅保留用户卡片 + API Key 面板。
 * 已移除 Element Plus：el-tag/el-button/ElMessage 改为原生 + useToast。
 */
import { computed, onMounted, ref } from "vue";
import { useUserStore } from "@/store/user";
import { toast } from "@/composables/useToast";
import ApiKeyPanel from "@/components/ApiKeyPanel.vue";

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
    // 401 由 http 拦截器处理；其他失败不阻塞页面骨架。
  } finally {
    loading.value = false;
  }
}

function handleLogout(): void {
  userStore.logout();
  toast.info("已登出");
}

onMounted(loadUser);
</script>

<template>
  <div class="account-view">
    <main class="account-body">
      <section class="account-card account-card--user">
        <div class="account-avatar">{{ avatarLetter }}</div>
        <div class="account-user-main">
          <h2 class="account-user-name">{{ userStore.displayName || "-" }}</h2>
          <p class="account-user-meta">@{{ userStore.user?.username || "-" }}</p>
          <div v-if="userStore.roles.length" class="account-roles">
            <span v-for="role in userStore.roles" :key="role" class="account-role-tag">{{ role }}</span>
          </div>
        </div>
        <button class="account-logout" type="button" @click="handleLogout">登出</button>
      </section>

      <section class="account-card">
        <h2 class="account-card-title">用户信息</h2>
        <dl class="account-info">
          <div class="account-info-row"><dt>昵称</dt><dd>{{ userStore.displayName || "-" }}</dd></div>
          <div class="account-info-row"><dt>用户名</dt><dd>{{ userStore.user?.username || "-" }}</dd></div>
          <div class="account-info-row"><dt>角色 ID</dt><dd>{{ userStore.roleId ?? "-" }}</dd></div>
        </dl>
      </section>

      <section class="account-card">
        <ApiKeyPanel />
      </section>
    </main>
  </div>
</template>

<style scoped>
.account-view { padding: 32px 24px 64px; max-width: 880px; margin: 0 auto; }

.account-body { display: flex; flex-direction: column; gap: 20px; }

.account-card {
  padding: 28px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(16px);
  border: 1px solid var(--border);
  border-radius: var(--r-xl);
  box-shadow: var(--shadow-md);
}

.account-card--user { display: flex; align-items: center; gap: 20px; }

.account-avatar {
  display: flex; width: 64px; height: 64px; flex-shrink: 0;
  align-items: center; justify-content: center;
  border-radius: 18px;
  background: var(--grad-teal-indigo);
  color: #fff; font-family: var(--font-display); font-size: 30px;
  box-shadow: var(--shadow-glow-teal);
}

.account-user-main { flex: 1; min-width: 0; }

.account-user-name {
  margin: 0; font-family: var(--font-display); font-size: 24px; font-weight: 400;
  color: var(--text-primary);
}

.account-user-meta { margin: 2px 0 0; color: var(--text-tertiary); font-size: 13px; font-family: var(--font-mono); }

.account-roles { margin-top: 8px; display: flex; gap: 6px; flex-wrap: wrap; }

.account-role-tag {
  padding: 2px 10px; border-radius: var(--r-full);
  background: var(--grad-brand-soft); color: var(--indigo-600);
  font-size: 11px; font-weight: 600;
}

.account-logout {
  padding: 8px 16px; border-radius: var(--r-md);
  border: 1px solid var(--border); background: var(--bg-card);
  color: var(--text-secondary); font-size: 13px; cursor: pointer;
  transition: all 0.2s var(--ease);
}

.account-logout:hover { border-color: var(--rose-400); color: var(--rose-500); }

.account-card-title {
  margin: 0 0 18px; font-size: 15px; font-weight: 600;
  color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.05em;
}

.account-info { display: flex; flex-direction: column; gap: 14px; margin: 0; }

.account-info-row { display: flex; align-items: center; }
.account-info-row dt { width: 120px; margin: 0; color: var(--text-tertiary); font-size: 13px; }
.account-info-row dd { margin: 0; font-size: 14px; color: var(--text-primary); }
</style>
