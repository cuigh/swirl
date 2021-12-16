<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="$router.push('/system/users/new')">
        <template #icon>
          <n-icon>
            <add-icon />
          </n-icon>
        </template>
        {{ t('buttons.new') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <x-tab>
      <x-tab-pane href="/system/users" :active="!$route.query.filter">{{ t('fields.all') }}</x-tab-pane>
      <x-tab-pane
        href="/system/users?filter=admins"
        :active="$route.query.filter === 'admins'"
      >{{ t('fields.admins') }}</x-tab-pane>
      <x-tab-pane
        href="/system/users?filter=active"
        :active="$route.query.filter === 'active'"
      >{{ t('fields.active') }}</x-tab-pane>
      <x-tab-pane
        href="/system/users?filter=blocked"
        :active="$route.query.filter === 'blocked'"
      >{{ t('fields.blocked') }}</x-tab-pane>
    </x-tab>
    <n-space :size="12">
      <n-input size="small" v-model:value="args.name" :placeholder="t('fields.name')" clearable />
      <n-input
        size="small"
        v-model:value="args.loginName"
        :placeholder="t('fields.login_name')"
        clearable
      />
      <n-button size="small" type="primary" @click="() => fetchData()">{{ t('buttons.search') }}</n-button>
    </n-space>
    <n-data-table
      remote
      :row-key="row => row.name"
      size="small"
      :columns="columns"
      :data="state.data"
      :pagination="pagination"
      :loading="state.loading"
      @update:page="fetchData"
      @update-page-size="changePageSize"
      scroll-x="max-content"
    />
  </n-space>
</template>

<script setup lang="ts">
import { reactive, watch } from "vue";
import {
  NSpace,
  NInput,
  NButton,
  NIcon,
  NDataTable,
} from "naive-ui";
import {
  AddOutline as AddIcon,
} from "@vicons/ionicons5";
import { useRoute, useRouter } from "vue-router";
import { XTab, XTabPane } from "@/components/tab";
import XPageHeader from "@/components/PageHeader.vue";
import userApi from "@/api/user";
import type { User } from "@/api/user";
import { useDataTable } from "@/utils/data-table";
import { renderButtons, renderLink, renderTag } from "@/utils/render";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const router = useRouter();
const args = reactive({
  name: "",
  loginName: "",
});
const columns = [
  {
    title: t('fields.id'),
    key: "id",
    render: (row: User) => renderLink(`/system/users/${row.id}`, row.id),
  },
  {
    title: t('fields.name'),
    key: "name",
  },
  {
    title: t('fields.login_name'),
    key: "loginName",
  },
  {
    title: t('fields.email'),
    key: "email",
  },
  {
    title: t('fields.admin'),
    key: "admin",
    render: (row: User) => t(row.admin ? 'enums.yes' : 'enums.no'),
  },
  {
    title: t('fields.status'),
    key: "status",
    render: (row: User) => renderTag(
      row.status ? t('enums.normal') : t('enums.blocked'),
      row.status ? "success" : "warning"
    ),
  },
  {
    title: t('fields.updated_at'),
    key: "updatedAt",
  },
  {
    title: t('fields.actions'),
    key: "actions",
    render(row: User, index: number) {
      return renderButtons([
        row.status ?
          { type: 'warning', text: t('buttons.block'), action: () => setStatus(row, 0), prompt: t('prompts.block'), } :
          { type: 'success', text: t('buttons.enable'), action: () => setStatus(row, 1) },
        { type: 'warning', text: t('buttons.edit'), action: () => router.push(`/system/users/${row.id}/edit`) },
        { type: 'error', text: t('buttons.delete'), action: () => remove(row, index), prompt: t('prompts.delete') },
      ])
    },
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(userApi.search, () => {
  return { ...args, filter: route.query.filter }
})

async function setStatus(u: User, status: number) {
  await userApi.setStatus({ id: u.id, status });
  u.status = status
}

async function remove(u: User, index: number) {
  await userApi.delete(u.id, u.name);
  state.data.splice(index, 1)
}

watch(() => route.query.filter, (newValue: any, oldValue: any) => {
  fetchData()
})
</script>