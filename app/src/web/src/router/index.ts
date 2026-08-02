import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router";
import { useUserStore } from "@/store/user";

// micro-app 子应用环境探测：被壳应用加载时 window.__MICRO_APP_ENVIRONMENT__ 为 true。
// 子应用路由 base 需对齐壳应用分配的路径（/auth），独立运行时用根路径。
const isMicroApp = (globalThis as Record<string, unknown>).__MICRO_APP_ENVIRONMENT__ === true;
const routerBase = isMicroApp ? "/auth" : "/";

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/Login.vue"),
    meta: { public: true },
  },
  {
    path: "/register",
    name: "register",
    component: () => import("@/views/Register.vue"),
    meta: { public: true },
  },
  {
    path: "/",
    name: "home",
    component: () => import("@/views/Home.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/",
  },
];

const router = createRouter({
  history: createWebHistory(routerBase),
  routes,
});

router.beforeEach((to) => {
  const userStore = useUserStore();
  // Unauthenticated users hitting a protected route are sent to /login; the
  // redirect query param lets us bounce back after a successful sign-in.
  if (to.meta.requiresAuth && !userStore.isAuthenticated) {
    return { name: "login", query: { redirect: to.fullPath } };
  }
  // Authenticated users visiting login/register are forwarded to the home page
  // to avoid re-authenticating unnecessarily.
  if (to.meta.public && userStore.isAuthenticated && to.name !== "register") {
    return { name: "home" };
  }
  return true;
});

export default router;
