<template>
  <x-page-header :subtitle="model.name ?? model.hostname">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'node_list' })">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-form :model="model" ref="form" label-placement="top">
      <n-grid cols="1 640:2" :x-gap="24">
        <n-form-item-gi :label="t('fields.name')" span="2" path="name" label-placement="left">
          <n-input
            :placeholder="t('fields.name')"
            v-model:value="model.name"
            style="max-width: 400px"
          />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.role')" path="role" label-placement="left">
          <n-radio-group v-model:value="model.role">
            <n-radio key="worker" value="worker">Worker</n-radio>
            <n-radio key="manager" value="manager">Manager</n-radio>
          </n-radio-group>
        </n-form-item-gi>
        <n-form-item-gi
          :label="t('fields.availability')"
          path="availability"
          label-placement="left"
        >
          <n-radio-group v-model:value="model.availability">
            <n-radio key="active" value="active">Active</n-radio>
            <n-radio key="pause" value="pause">Pause</n-radio>
            <n-radio key="drain" value="drain">Drain</n-radio>
          </n-radio-group>
        </n-form-item-gi>
        <n-form-item-gi span="3" :label="t('fields.tags')" path="labels">
          <n-dynamic-input v-model:value="model.labels" #="{ index, value }" :on-create="newLabel">
            <n-input :placeholder="t('fields.name')" v-model:value="value.name" />
            <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
            <n-input :placeholder="t('fields.value')" v-model:value="value.value" />
          </n-dynamic-input>
        </n-form-item-gi>
        <n-gi :span="2">
          <n-button
            @click.prevent="submit"
            type="primary"
            :disabled="submiting"
            :loading="submiting"
          >
            <template #icon>
              <n-icon>
                <save-icon />
              </n-icon>
            </template>
            {{ t('buttons.save') }}
          </n-button>
        </n-gi>
      </n-grid>
    </n-form>
  </n-space>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from "vue";
import {
  NButton,
  NSpace,
  NInput,
  NIcon,
  NForm,
  NGrid,
  NGi,
  NFormItemGi,
  NDynamicInput,
  NRadioGroup,
  NRadio,
} from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import nodeApi from "@/api/node";
import type { Node } from "@/api/node";
import { useRoute } from "vue-router";
import { router } from "@/router/router";
import { useForm } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({} as Node);
const form = ref();
const { submit, submiting } = useForm(form, () => nodeApi.save(model.value), () => {
  window.message.info(t('texts.action_success'));
  router.push({ name: 'node_list' })
})

function newLabel() {
  return {
    name: '',
    value: ''
  }
}

async function fetchData() {
  const id = route.params.id as string
  let r = await nodeApi.find(id);
  model.value = r.data?.node as Node;
}

onMounted(fetchData);
</script>
