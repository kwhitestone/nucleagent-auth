import { createApp, type App as VueApp } from "vue";
import { createPinia } from "pinia";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css";

import App from "./App.vue";
import router from "./router";
import i18n from "./i18n";
import "./styles/global.css";

// 子应用挂载点用唯一 ID（#auth-app），避免与壳的 #app 或其他子应用冲突。
// index.html 里的 <div id="auth-app"> 对应（独立运行时），micro-app inline 模式下
// 子应用 HTML 被注入 <micro-app-body>，该 div 也随之进入。
const MOUNT_ID = "auth-app";

let app: VueApp | null = null;

function mount() {
  app = createApp(App);
  app.use(createPinia());
  app.use(router);
  app.use(i18n);
  app.use(ElementPlus);
  app.mount(`#${MOUNT_ID}`);
}

function unmount() {
  if (app) {
    app.unmount();
    app = null;
  }
}

const w = globalThis as Record<string, unknown>;
if (w.__MICRO_APP_ENVIRONMENT__) {
  w.mount = mount;
  w.unmount = unmount;
} else {
  mount();
}
