<template>
  <x-page-header :subtitle="model.name">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'role_list' })">
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
        @click="$router.push({ name: 'role_edit', params: { id: model.id } })"
      >{{ t('buttons.edit') }}</n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="16">
    <x-description cols="1 640:2" label-position="left" label-align="right" :label-width="85">
      <x-description-item :label="t('fields.id')">{{ model.id }}</x-description-item>
      <x-description-item :label="t('fields.name')">{{ model.name }}</x-description-item>
      <x-description-item :span="2" :label="t('fields.desc')">{{ model.desc }}</x-description-item>
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
    <x-panel :title="t('fields.perms')">
      <n-grid cols="1 640:2 960:3 1440:4" x-gap="6" y-gap="6">
        <n-gi span="1" v-for="g in ps">
          <x-pair-tag type="warning" :label="g.group" :value="g.items" />
        </n-gi>
      </n-grid>
    </x-panel>
  </n-space>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import {
  NButton,
  NSpace,
  NIcon,
  NGrid,
  NGi,
  NTime,
} from "naive-ui";
import { useRoute } from "vue-router";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XPanel from "@/components/Panel.vue";
import XPairTag from "@/components/PairTag.vue";
import XAnchor from "@/components/Anchor.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import roleApi from "@/api/role";
import type { Role } from "@/api/role";
import { perms } from "@/utils/perm";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({} as Role);
const ps = computed(() => {
  const set = new Set(model.value.perms)
  const arr: any = []
  perms.forEach(g => {
    const items = g.actions.filter(action => set.has(g.key + '.' + action)).map(p => t('perms.' + p))
    items.length && arr.push({ group: t('objects.' + g.key), items: items })
  })
  return arr
})

async function fetchData() {
  let r = await roleApi.find(route.params.id as string);
  model.value = r.data as Role;
}

onMounted(fetchData);
</script>