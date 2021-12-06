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
      <x-description-item :span="2" :label="t('fields.perms')">
        <n-grid cols="1 480:2 960:3 1440:4" x-gap="6">
          <n-gi span="1" v-for="g in ps">
            <x-pair-tag type="warning" :label="g.group" :value="g.items" />
          </n-gi>
        </n-grid>
      </x-description-item>
    </x-description>
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
} from "naive-ui";
import { useRoute } from "vue-router";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XPairTag from "@/components/PairTag.vue";
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
    const items = g.items.filter(p => set.has(p.key)).map(p => t('perms.' + p.perm))
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