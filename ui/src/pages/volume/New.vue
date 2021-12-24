<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'volume_list' })">
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
    <n-form :model="model" :rules="rules" ref="form" label-placement="top" label-width="90">
      <n-grid cols="1 640:2" :x-gap="24">
        <n-form-item-gi :label="t('fields.name')" span="1" path="name" label-placement="left">
          <n-input :placeholder="t('fields.name')" v-model:value="model.name" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.driver')" path="driver" span="2" label-placement="left">
          <n-radio-group v-model:value="model.driver">
            <n-radio key="local" value="local">Local</n-radio>
            <n-radio key="other" value="other">Other</n-radio>
          </n-radio-group>
          <n-input
            :placeholder="t('fields.driver_name')"
            v-model:value="model.customDriver"
            :disabled="model.driver === 'local'"
            style="max-width: 150px"
          />
        </n-form-item-gi>
        <n-form-item-gi span="3" :label="t('fields.options')" path="options">
          <n-dynamic-input v-model:value="model.options" #="{ index, value }" :on-create="newPair">
            <n-input :placeholder="t('fields.name')" v-model:value="value.name" />
            <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
            <n-input :placeholder="t('fields.value')" v-model:value="value.value" />
          </n-dynamic-input>
        </n-form-item-gi>
        <n-form-item-gi span="3" :label="t('fields.labels')" path="labels">
          <n-dynamic-input v-model:value="model.labels" #="{ index, value }" :on-create="newPair">
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
import { ref } from "vue";
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
import volumeApi from "@/api/volume";
import type { Volume } from "@/api/volume";
import { router } from "@/router/router";
import { useForm, requiredRule, customRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const model = ref({ driver: 'local' } as Volume);
const rules: any = {
  name: requiredRule(),
  driver: customRule((rule: any, value: string): boolean => {
    if (value === 'other') {
      return Boolean(model.value.customDriver)
    }
    return true;
  }, t('tips.volume_driver_rule')),
};
const form = ref();
const { submit, submiting } = useForm(form, () => volumeApi.save(model.value), () => {
  window.message.info(t('texts.action_success'));
  router.push({ name: 'volume_list' })
})

function newPair() {
  return { name: '', value: '' }
}
</script>