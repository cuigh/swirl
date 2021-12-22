<template>
  <x-page-header :subtitle="model.user.name">
    <template #action>
      <n-button secondary size="small" @click="$router.push('/system/users')">
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
        @click="$router.push(`/system/users/${model.user.id}/edit`)"
      >{{ t('buttons.edit') }}</n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="16">
    <x-description cols="1 640:2" label-position="left" label-align="right" :label-width="100">
      <x-description-item :label="t('fields.id')">{{ model.user.id }}</x-description-item>
      <x-description-item :label="t('fields.email')">{{ model.user.email }}</x-description-item>
      <x-description-item :label="t('fields.username')">{{ model.user.name }}</x-description-item>
      <x-description-item :label="t('fields.login_name')">{{ model.user.loginName }}</x-description-item>
      <x-description-item :label="t('fields.status')">
        <n-tag
          round
          size="small"
          :type="model.user.status ? 'primary' : 'warning'"
        >{{ t(model.user.status ? 'enums.normal' : 'enums.blocked') }}</n-tag>
      </x-description-item>
      <x-description-item :label="t('fields.admin')">
        <n-tag
          size="small"
          round
          :type="model.user.admin ? 'success' : 'default'"
        >{{ t(model.user.admin ? "enums.yes" : "enums.no") }}</n-tag>
      </x-description-item>
      <x-description-item :label="t('fields.type')" :span="2">
        <n-tag
          size="small"
          round
          :type="model.user.type === 'internal' ? 'default' : 'warning'"
        >{{ model.user.type }}</n-tag>
      </x-description-item>
      <x-description-item :label="t('fields.created_by')">
        <x-anchor
          :url="`/system/users/${model.user.createdBy?.id}`"
        >{{ model.user.createdBy?.name }}</x-anchor>
      </x-description-item>
      <x-description-item :label="t('fields.created_at')">
        <n-time :time="model.user.createdAt" format="y-MM-dd HH:mm:ss" />
      </x-description-item>
      <x-description-item :label="t('fields.updated_by')">
        <x-anchor
          :url="`/system/users/${model.user.updatedBy?.id}`"
        >{{ model.user.updatedBy?.name }}</x-anchor>
      </x-description-item>
      <x-description-item :label="t('fields.updated_at')">
        <n-time :time="model.user.updatedAt" format="y-MM-dd HH:mm:ss" />
      </x-description-item>
      <x-description-item
        :label="t('objects.role', 2)"
        :span="2"
        v-if="model.user.roles && model.user.roles.length > 0"
      >
        <n-space :size="6">
          <n-tag
            round
            size="small"
            type="info"
            v-for="r in model.user.roles"
          >{{ model.roles.get(r) }}</n-tag>
        </n-space>
      </x-description-item>
    </x-description>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, reactive, watch } from "vue";
import {
  NButton,
  NTag,
  NSpace,
  NIcon,
  NTime,
} from "naive-ui";
import { useRoute } from "vue-router";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XAnchor from "@/components/Anchor.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import userApi from "@/api/user";
import roleApi from "@/api/role";
import type { User } from "@/api/user";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = reactive({
  user: {} as User,
  roles: new Map<string, string>(),
});

async function fetchData() {
  let user = await userApi.find(route.params.id as string);
  let roles = await roleApi.search()
  model.user = user.data as User
  roles.data?.forEach(r => model.roles.set(r.id, r.name))
}

watch(() => route.params.id, fetchData)

onMounted(fetchData);
</script>