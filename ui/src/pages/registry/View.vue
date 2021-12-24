<template>
  <x-page-header :subtitle="model.name">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'registry_list' })">
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
        @click="$router.push({ name: 'registry_edit', params: { id: model.id } })"
      >{{ t('buttons.edit') }}</n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="16">
    <x-description label-placement="left" label-align="right" :label-width="90">
      <x-description-item :label="t('fields.id')">{{ model.id }}</x-description-item>
      <x-description-item :label="t('fields.name')">{{ model.name }}</x-description-item>
      <x-description-item :label="t('fields.url')">{{ model.url }}</x-description-item>
      <x-description-item :label="t('fields.login_name')">{{ model.username }}</x-description-item>
      <x-description-item :label="t('fields.created_by')">
        <x-anchor
          :url="{ name: 'user_detail', params: { id: model.createdBy?.id } }"
          v-if="model.createdBy?.id"
        >{{ model.createdBy?.name }}</x-anchor>
      </x-description-item>
      <x-description-item :label="t('fields.created_at')">
        <n-time :time="model.createdAt" format="y-MM-dd HH:mm:ss" />
      </x-description-item>
      <x-description-item :label="t('fields.updated_by')">
        <x-anchor
          :url="{ name: 'user_detail', params: { id: model.updatedBy?.id } }"
          v-if="model.updatedBy?.id"
        >{{ model.updatedBy?.name }}</x-anchor>
      </x-description-item>
      <x-description-item :label="t('fields.updated_at')">
        <n-time :time="model.updatedAt" format="y-MM-dd HH:mm:ss" />
      </x-description-item>
    </x-description>
  </n-space>
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
import XAnchor from "@/components/Anchor.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import registryApi from "@/api/registry";
import type { Registry } from "@/api/registry";
import { useRoute } from "vue-router";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({} as Registry);

async function fetchData() {
  let r = await registryApi.find(route.params.id as string)
  model.value = r.data as Registry;
}

onMounted(fetchData);
</script>