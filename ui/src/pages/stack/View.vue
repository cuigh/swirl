<template>
  <x-page-header :subtitle="model.name || model.id">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'stack_list' })">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
      <n-button
        secondary
        size="small"
        @click="$router.push({ name: 'stack_edit', params: { name: model.name } })"
      >{{ t('buttons.edit') }}</n-button>
    </template>
  </x-page-header>
  <div class="page-body">
    <n-space vertical :size="16">
      <x-description :label-width="90">
        <x-description-item :label="t('fields.name')" :span="2">{{ model.name }}</x-description-item>
        <x-description-item :label="t('fields.created_by')">
          <x-anchor :url="`/system/users/${model.createdBy?.id}`">{{ model.createdBy?.name }}</x-anchor>
        </x-description-item>
        <x-description-item :label="t('fields.created_at')">
          <n-time :time="model.createdAt" format="y-MM-dd HH:mm:ss" />
        </x-description-item>
        <x-description-item :label="t('fields.updated_by')">
          <x-anchor :url="`/system/users/${model.updatedBy?.id}`">{{ model.updatedBy?.name }}</x-anchor>
        </x-description-item>
        <x-description-item :label="t('fields.updated_at')">
          <n-time :time="model.updatedAt" format="y-MM-dd HH:mm:ss" />
        </x-description-item>
      </x-description>
      <x-panel :title="t('fields.content')">
        <x-code :code="model.content" language="yaml" />
      </x-panel>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NButton,
  NSpace,
  NIcon,
  NTime,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XCode from "@/components/Code.vue";
import XPanel from "@/components/Panel.vue";
import XAnchor from "@/components/Anchor.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import stackApi from "@/api/stack";
import type { Stack } from "@/api/stack";
import { useRoute } from "vue-router";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({} as Stack);

async function fetchData() {
  const name = route.params.name as string;
  let r = await stackApi.find(name);
  model.value = r.data as Stack;
}

onMounted(fetchData);
</script>