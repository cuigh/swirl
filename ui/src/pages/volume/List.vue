<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" type="warning" @click="pruneVolume">
        <template #icon>
          <n-icon>
            <close-icon />
          </n-icon>
        </template>
        {{ t('buttons.prune') }}
      </n-button>
      <n-button secondary size="small" @click="$router.push('/local/volumes/new')">
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
      <n-input size="small" v-model:value="filter.name" :placeholder="t('fields.name')" clearable />
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
import { reactive } from "vue";
import {
  NSpace,
  NButton,
  NDataTable,
  NInput,
  NIcon,
} from "naive-ui";
import { AddOutline as AddIcon, CloseOutline as CloseIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import volumeApi from "@/api/volume";
import type { Volume } from "@/api/volume";
import { useDataTable } from "@/utils/data-table";
import { formatSize, renderButton, renderLink, renderTag } from "@/utils/render";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const filter = reactive({
  name: "",
});
const columns = [
  {
    title: t('fields.name'),
    key: "name",
    fixed: "left" as const,
    render: (v: Volume) => renderLink(`/local/volumes/${v.name}`, v.name),
  },
  {
    title: t('fields.driver'),
    key: "driver",
    render(v: Volume) {
      return renderTag(v.driver)
    }
  },
  {
    title: t('fields.scope'),
    key: "scope",
    render(v: Volume) {
      return renderTag(v.scope)
    }
  },
  {
    title: t('fields.mount_point'),
    key: "mountPoint"
  },
  {
    title: t('fields.created_at'),
    key: "createdAt"
  },
  {
    title: t('fields.actions'),
    key: "actions",
    render(v: Volume, index: number) {
      return renderButton('error', t('buttons.delete'), () => deleteVolume(v.name, index), t('prompts.delete'))
    },
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(volumeApi.search, filter)

async function deleteVolume(name: string, index: number) {
  await volumeApi.delete(name);
  state.data.splice(index, 1)
}

async function pruneVolume() {
  window.dialog.warning({
    title: t('dialogs.prune_volume.title'),
    content: t('dialogs.prune_volume.body'),
    positiveText: t('buttons.confirm'),
    negativeText: t('buttons.cancel'),
    onPositiveClick: async () => {
      const r = await volumeApi.prune();
      window.message.info(t('texts.prune_volume_success', {
        count: r.data?.deletedVolumes.length,
        size: formatSize(r.data?.reclaimedSpace as number)
      }));
      fetchData();
    }
  })
}
</script>