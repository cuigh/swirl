<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="$router.push('/swarm/networks')">
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
        <n-form-item-gi :label="t('fields.name')" span="2" path="name" label-placement="left">
          <n-input
            :placeholder="t('fields.name')"
            v-model:value="model.name"
            style="max-width: 400px"
          />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.driver')" path="driver" span="1" label-placement="left">
          <n-radio-group v-model:value="model.driver">
            <n-radio key="overlay" value="overlay">Overlay</n-radio>
            <n-radio key="bridge" value="bridge">Bridge</n-radio>
            <n-radio key="ipvlan" value="ipvlan">IPvlan</n-radio>
            <n-radio key="macvlan" value="macvlan">Macvlan</n-radio>
            <n-radio key="host" value="host">Host</n-radio>
          </n-radio-group>
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.scope')" path="scope" span="1" label-placement="left">
          <n-radio-group v-model:value="model.scope">
            <n-radio key="swarm" value="swarm">Swarm</n-radio>
            <n-radio key="local" value="local">Local</n-radio>
          </n-radio-group>
        </n-form-item-gi>
        <n-form-item-gi
          :label="t('fields.internal')"
          path="internal"
          label-placement="left"
          label-align="right"
        >
          <n-switch v-model:value="model.internal" />
        </n-form-item-gi>
        <n-form-item-gi
          :label="t('fields.attachable')"
          path="attachable"
          label-placement="left"
          label-align="right"
        >
          <n-switch v-model:value="model.attachable" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.ingress')" path="ingress" label-placement="left">
          <n-switch v-model:value="model.ingress" />
        </n-form-item-gi>
        <n-form-item-gi label="IPv6" path="ipv6" label-placement="left" label-align="right">
          <n-switch v-model:value="model.ipv6">
            <template #checked>{{ t('enums.enabled') }}</template>
            <template #unchecked>{{ t('enums.disabled') }}</template>
          </n-switch>
        </n-form-item-gi>
        <n-form-item-gi span="3" label="IP" path="ipam.config">
          <n-dynamic-input
            v-model:value="model.ipam.config"
            #="{ index, value }"
            :on-create="newConfig"
          >
            <n-input-group>
              <n-input-group-label>{{ t('fields.subnet') }}</n-input-group-label>
              <n-input placeholder="e.g. 172.20.0.0/16" v-model:value="value.subnet" />
            </n-input-group>
            <n-input-group style="margin: auto 10px">
              <n-input-group-label>{{ t('fields.gateway') }}</n-input-group-label>
              <n-input placeholder="e.g. 172.20.10.11" v-model:value="value.gateway" />
            </n-input-group>
            <n-input-group>
              <n-input-group-label>{{ t('fields.range') }}</n-input-group-label>
              <n-input placeholder="e.g. 172.20.10.0/24" v-model:value="value.ipRange" />
            </n-input-group>
          </n-dynamic-input>
        </n-form-item-gi>
        <n-form-item-gi span="3" :label="t('fields.options')" path="options">
          <n-dynamic-input v-model:value="model.options" #="{ index, value }" :on-create="newPair">
            <n-input :placeholder="t('fields.name')" v-model:value="value.name" />
            <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
            <n-input :placeholder="t('fields.value')" v-model:value="value.value" />
          </n-dynamic-input>
        </n-form-item-gi>
        <n-form-item-gi span="3" :label="t('fields.tags')" path="labels">
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
  NSwitch,
  NDynamicInput,
  NRadioGroup,
  NRadio,
  NInputGroup,
  NInputGroupLabel,
} from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import networkApi from "@/api/network";
import type { Network } from "@/api/network";
import { router } from "@/router/router";
import { useForm, requiredRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const model = ref({ driver: 'overlay', scope: 'swarm', ipam: {} } as Network);
const rules: any = {
  name: requiredRule(),
};
const form = ref();
const { submit, submiting } = useForm(form, () => networkApi.save(model.value), () => {
  window.message.info(t('texts.action_success'));
  router.push("/swarm/networks")
})

function newConfig() {
  return { subnet: '', gateway: '', range: '' }
}

function newPair() {
  return { name: '', value: '' }
}
</script>