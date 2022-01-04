<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" type="warning" @click="prune">
        <template #icon>
          <n-icon>
            <close-icon />
          </n-icon>
        </template>
        {{ t('buttons.prune') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-space :size="12">
      <n-select
        size="small"
        :placeholder="t('fields.type')"
        v-model:value="filter.type"
        :options="types"
        style="width: 140px"
        clearable
      />
      <n-input
        size="small"
        v-model:value="filter.name"
        :placeholder="t('fields.object')"
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
import { h, reactive, ref } from "vue";
import {
  NSpace,
  NButton,
  NDataTable,
  NSelect,
  NInput,
  NIcon,
  NFormItem,
  NInputNumber,
} from "naive-ui";
import { CloseOutline as CloseIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import eventApi from "@/api/event";
import type { Event } from "@/api/event";
import { useDataTable } from "@/utils/data-table";
import { renderLink, renderTag, renderTime } from "@/utils/render";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const filter = reactive({
  type: undefined,
  name: "",
});
const types: any = [
  {
    type: 'group',
    label: 'System',
    key: 'system',
    children: [
      {
        label: 'User',
        value: 'User'
      },
      {
        label: 'Role',
        value: 'Role'
      },
      {
        label: 'Chart',
        value: 'Chart'
      },
      {
        label: 'Setting',
        value: 'Setting'
      },
    ],
  },
  {
    type: 'group',
    label: 'Swarm',
    key: 'swarm',
    children: [
      {
        label: 'Registry',
        value: 'Registry'
      },
      {
        label: 'Node',
        value: 'Node'
      },
      {
        label: 'Network',
        value: 'Network'
      },
      {
        label: 'Service',
        value: 'Service'
      },
      {
        label: 'Stack',
        value: 'Stack'
      },
      {
        label: 'Secret',
        value: 'Secret'
      },
      {
        label: 'Config',
        value: 'Config'
      },
    ],
  },
  {
    type: 'group',
    label: 'Local',
    key: 'local',
    children: [
      {
        label: 'Image',
        value: 'Image'
      },
      {
        label: 'Container',
        value: 'Container'
      },
      {
        label: 'Volume',
        value: 'Volume'
      },
    ],
  },
]
const columns = [
  {
    title: t('fields.id'),
    key: "id",
    width: 210,
    fixed: "left" as const,
  },
  {
    title: t('fields.type'),
    key: "type",
    render(e: Event) {
      return renderTag(e.type)
    },
  },
  {
    title: t('fields.action'),
    key: "action",
    render(e: Event) {
      return renderTag(e.action)
    },
  },
  {
    title: t('fields.object'),
    key: "name",
    render: renderObject,
  },
  {
    title: t('fields.operator'),
    key: "name",
    render: (e: Event) => e.userId ? renderLink({ name: 'user_detail', params: { id: e.userId } }, e.username) : null,
  },
  {
    title: t('fields.time'),
    key: "time",
    render: (e: Event) => renderTime(e.time),
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(eventApi.search, filter)

function renderObject(e: Event) {
  switch (e.type) {
    case "User":
    case "Role":
    case "Chart":
    case "Registry":
    case "Node":
    case "Config":
    case "Secret":
      return renderLink({ name: e.type.toLowerCase() + '_detail', params: { id: e.args.id } }, e.args.name)
    case "Network":
    case "Service":
    case "Stack":
      return renderLink({ name: e.type.toLowerCase() + '_detail', params: { name: e.args.name } }, e.args.name)
    case "Image":
      if (e.args.id) {
        return renderLink({ name: 'image_detail', params: { node: e.args.node || '-', id: e.args.id } }, e.args.id)
      } else {
        return renderLink({ name: 'image_list' }, t('objects.image'))
      }
    case "Container":
      if (e.args.id) {
        return renderLink({ name: 'container_detail', params: { node: e.args.node || '-', id: e.args.id } }, e.args.name)
      } else {
        return renderLink({ name: 'container_list' }, t('objects.container'))
      }
    case "Volume":
      if (e.args.name) {
        return renderLink({ name: 'volume_detail', params: { node: e.args.node || '-', name: e.args.name } }, e.args.name)
      } else {
        return renderLink({ name: 'volume_list' }, t('objects.volume'))
      }
    case "Setting":
      return renderLink({ name: 'setting' }, t('objects.setting'))
  }
  return null
}

function prune() {
  const days = ref(7) as any
  window.dialog.warning({
    title: t('dialogs.prune_event.title'),
    content: () => h(
      NFormItem,
      { label: t('dialogs.prune_event.label'), labelPlacement: 'top', showFeedback: false },
      { default: () => h(NInputNumber, { min: 0, defaultValue: days, style: 'width: 100%' }) }
    ),
    positiveText: t('buttons.confirm'),
    negativeText: t('buttons.cancel'),
    onPositiveClick: async () => {
      eventApi.prune(days.value);
      window.message.success(t('texts.action_success'))
      fetchData()
    }
  })
}
</script>