<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'role_new' })">
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
    <n-space :size="12">
      <n-input size="small" v-model:value="model.name" :placeholder="t('fields.name')" clearable />
      <n-button size="small" type="primary" @click="fetchData">{{ t('buttons.search') }}</n-button>
    </n-space>
    <n-table size="small" :bordered="true" :single-line="false">
      <thead>
        <tr>
          <th>{{ t('fields.id') }}</th>
          <th>{{ t('fields.name') }}</th>
          <th>{{ t('fields.desc') }}</th>
          <th>{{ t('fields.updated_at') }}</th>
          <th>{{ t('fields.actions') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(r, index) of model.roles" :key="r.id">
          <td>
            <x-anchor :url="`/system/roles/${r.id}`">{{ r.id }}</x-anchor>
          </td>
          <td>{{ r.name }}</td>
          <td>{{ r.desc }}</td>
          <td>{{ r.updatedAt }}</td>
          <td>
            <n-popconfirm :show-icon="false" @positive-click="deleteRole(r, index)">
              <template #trigger>
                <n-button size="tiny" quaternary type="error">{{ t('buttons.delete') }}</n-button>
              </template>
              {{ t('prompts.delete') }}
            </n-popconfirm>
            <n-button
              size="tiny"
              quaternary
              type="warning"
              @click="$router.push({ name: 'role_edit', params: { id: r.id } })"
            >{{ t('buttons.edit') }}</n-button>
          </td>
        </tr>
      </tbody>
    </n-table>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, reactive } from "vue";
import {
  NSpace,
  NInput,
  NButton,
  NIcon,
  NTable,
  NPopconfirm,
} from "naive-ui";
import {
  AddOutline as AddIcon,
} from "@vicons/ionicons5";
import XAnchor from "@/components/Anchor.vue";
import XPageHeader from "@/components/PageHeader.vue";
import roleApi from "@/api/role";
import type { Role } from "@/api/role";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const model = reactive({
  name: "",
  roles: [] as Role[],
});

async function deleteRole(r: Role, index: number) {
  await roleApi.delete(r.id, r.name);
  model.roles.splice(index, 1)
}

async function fetchData() {
  let r = await roleApi.search(model.name);
  model.roles = r.data || [];
}

onMounted(fetchData);
</script>