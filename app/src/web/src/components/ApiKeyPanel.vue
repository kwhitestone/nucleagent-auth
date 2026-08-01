<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import {
  ElMessage,
  ElMessageBox,
  type FormInstance,
  type FormRules,
} from "element-plus";
import {
  createApiKey,
  deleteApiKey,
  listApiKeys,
} from "@/api/auth";
import type { ApiKey, ApiKeyWithSecret } from "@/api/types";

const { t } = useI18n();

const keys = ref<ApiKey[]>([]);
const loading = ref(false);
const createLoading = ref(false);

const createDialogVisible = ref(false);
const createFormRef = ref<FormInstance>();
const createForm = ref({ name: "" });
const createRules: FormRules<{ name: string }> = {
  name: [{ required: true, message: t("apiKey.createPlaceholder"), trigger: "blur" }],
};

const plaintextDialogVisible = ref(false);
const createdKey = ref<ApiKeyWithSecret | null>(null);

async function loadKeys(): Promise<void> {
  loading.value = true;
  try {
    keys.value = await listApiKeys();
  } catch (error) {
    const message = error instanceof Error ? error.message : t("common.networkError");
    ElMessage.error(message);
  } finally {
    loading.value = false;
  }
}

function openCreateDialog(): void {
  createForm.value = { name: "" };
  createDialogVisible.value = true;
}

async function handleCreate(): Promise<void> {
  const valid = await createFormRef.value?.validate().catch(() => false);
  if (!valid) return;

  createLoading.value = true;
  try {
    const created = await createApiKey({ name: createForm.value.name.trim() });
    createdKey.value = created;
    createDialogVisible.value = false;
    plaintextDialogVisible.value = true;
    ElMessage.success(t("apiKey.createSuccess"));
    await loadKeys();
  } catch (error) {
    const message = error instanceof Error ? error.message : t("apiKey.createFailed");
    ElMessage.error(message);
  } finally {
    createLoading.value = false;
  }
}

async function handleDelete(key: ApiKey): Promise<void> {
  try {
    await ElMessageBox.confirm(t("apiKey.deleteConfirm"), t("apiKey.deleteConfirmTitle"), {
      type: "warning",
      confirmButtonText: t("common.confirm"),
      cancelButtonText: t("common.cancel"),
    });
  } catch {
    return; // user cancelled
  }

  try {
    await deleteApiKey(key.id);
    ElMessage.success(t("apiKey.deleteSuccess"));
    await loadKeys();
  } catch (error) {
    const message = error instanceof Error ? error.message : t("apiKey.deleteFailed");
    ElMessage.error(message);
  }
}

async function copyPlaintext(): Promise<void> {
  if (!createdKey.value) return;
  try {
    await navigator.clipboard.writeText(createdKey.value.plaintext);
    ElMessage.success(t("common.copied"));
  } catch {
    ElMessage.warning(t("common.copy"));
  }
}

onMounted(loadKeys);
</script>

<template>
  <div class="api-key-panel">
    <div class="api-key-panel__header">
      <h2 class="api-key-panel__title">{{ t("home.apiKeysTitle") }}</h2>
      <el-button type="primary" @click="openCreateDialog">
        {{ t("apiKey.create") }}
      </el-button>
    </div>

    <el-table v-loading="loading" :data="keys" border stripe>
      <el-table-column :label="t('apiKey.name')" prop="name" min-width="160" />
      <el-table-column :label="t('apiKey.prefix')" prop="prefix" min-width="140" />
      <el-table-column :label="t('apiKey.enabled')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.enable ? 'success' : 'info'" size="small">
            {{ row.enable ? t("apiKey.enabledText") : t("apiKey.disabledText") }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('apiKey.lastUsed')" min-width="160">
        <template #default="{ row }">
          {{ row.lastUsed ? row.lastUsed : t("apiKey.neverUsed") }}
        </template>
      </el-table-column>
      <el-table-column :label="t('apiKey.createdAt')" prop="createdAt" min-width="160" />
      <el-table-column :label="t('apiKey.actions')" width="100" fixed="right">
        <template #default="{ row }">
          <el-button type="danger" link @click="handleDelete(row as ApiKey)">
            {{ t("common.delete") }}
          </el-button>
        </template>
      </el-table-column>
      <template #empty>
        <span>{{ t("apiKey.empty") }}</span>
      </template>
    </el-table>

    <!-- Create dialog -->
    <el-dialog
      v-model="createDialogVisible"
      :title="t('apiKey.create')"
      width="420px"
    >
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-position="top"
      >
        <el-form-item :label="t('apiKey.name')" prop="name">
          <el-input
            v-model="createForm.name"
            :placeholder="t('apiKey.createPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">
          {{ t("common.cancel") }}
        </el-button>
        <el-button type="primary" :loading="createLoading" @click="handleCreate">
          {{ t("common.confirm") }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Plaintext dialog (shown once) -->
    <el-dialog
      v-model="plaintextDialogVisible"
      :title="t('apiKey.plaintextTitle')"
      width="520px"
      :close-on-click-modal="false"
    >
      <el-alert
        type="warning"
        :title="t('apiKey.plaintextTip')"
        show-icon
        :closable="false"
      />
      <div class="api-key-panel__plaintext">
        <code>{{ createdKey?.plaintext }}</code>
        <el-button link type="primary" @click="copyPlaintext">
          {{ t("common.copy") }}
        </el-button>
      </div>
      <template #footer>
        <el-button type="primary" @click="plaintextDialogVisible = false">
          {{ t("common.close") }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.api-key-panel {
  width: 100%;
}

.api-key-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.api-key-panel__title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.api-key-panel__plaintext {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 6px;
  word-break: break-all;
}

.api-key-panel__plaintext code {
  flex: 1;
  font-size: 13px;
  color: #1f2329;
}
</style>
