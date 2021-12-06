<template>
  <x-page-header :subtitle="model.name">
    <template #action>
      <n-button secondary size="small" @click="$router.push('/swarm/secrets')">
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
        <n-form-item-gi :label="t('fields.name')" path="name">
          <n-input
            :placeholder="t('fields.name')"
            v-model:value="model.name"
            :disabled="Boolean(model.id)"
          />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.content')" path="data" span="2" v-if="!model.id">
          <n-input
            type="textarea"
            :placeholder="t('fields.content')"
            v-model:value="model.data"
            :autosize="{ minRows: 5, maxRows: 30 }"
          />
        </n-form-item-gi>
        <n-form-item-gi span="3" :label="t('fields.labels')" path="labels">
          <n-dynamic-input v-model:value="model.labels" #="{ index, value }" :on-create="newPair">
            <n-input :placeholder="t('fields.name')" v-model:value="value.name" />
            <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
            <n-input :placeholder="t('fields.value')" v-model:value="value.value" />
          </n-dynamic-input>
        </n-form-item-gi>
      </n-grid>
      <x-panel :title="t('fields.engine')" v-if="!model.id">
        <n-grid cols="1 640:2" :x-gap="24">
          <n-form-item-gi :label="t('fields.name')" path="driver.name">
            <n-input :placeholder="t('fields.name')" v-model:value="model.driver.name" />
          </n-form-item-gi>
          <n-form-item-gi span="3" :label="t('fields.options')" path="driver.options">
            <n-dynamic-input
              v-model:value="model.driver.options"
              #="{ index, value }"
              :on-create="newPair"
            >
              <n-input :placeholder="t('fields.name')" v-model:value="value.name" />
              <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
              <n-input :placeholder="t('fields.value')" v-model:value="value.value" />
            </n-dynamic-input>
          </n-form-item-gi>
        </n-grid>
      </x-panel>
      <x-panel :title="t('fields.template')" :subtitle="t('tips.template')" v-if="!model.id">
        <n-grid cols="1 640:2" :x-gap="24">
          <n-form-item-gi :label="t('fields.engine')" path="templating.name">
            <n-radio-group v-model:value="model.templating.name">
              <n-radio key="none" value="none">None</n-radio>
              <n-radio key="golang" value="golang">Golang</n-radio>
            </n-radio-group>
          </n-form-item-gi>
          <n-form-item-gi span="3" :label="t('fields.options')" path="templating.options">
            <n-dynamic-input
              v-model:value="model.templating.options"
              #="{ index, value }"
              :on-create="newPair"
            >
              <n-input :placeholder="t('fields.name')" v-model:value="value.name" />
              <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
              <n-input :placeholder="t('fields.value')" v-model:value="value.value" />
            </n-dynamic-input>
          </n-form-item-gi>
        </n-grid>
      </x-panel>
      <n-button @click.prevent="submit" type="primary" :disabled="submiting" :loading="submiting">
        <template #icon>
          <n-icon>
            <save-icon />
          </n-icon>
        </template>
        {{ t('buttons.save') }}
      </n-button>
    </n-form>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NButton,
  NSpace,
  NInput,
  NIcon,
  NForm,
  NGrid,
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
import XPanel from "@/components/Panel.vue";
import secretApi from "@/api/secret";
import type { Secret } from "@/api/secret";
import { useRoute } from "vue-router";
import { router } from "@/router/router";
import { useForm, requiredRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const name = route.params.name as string || ''
const model = ref({
  driver: {},
  templating: { name: 'none' },
} as Secret);
const rules: any = {
  name: requiredRule(),
  data: requiredRule(),
};
const form = ref();
const { submit, submiting } = useForm(form, () => secretApi.save(model.value), () => {
  window.message.info(t('texts.action_success'));
  router.push("/swarm/secrets")
})

function newPair() {
  return { name: '', value: '' }
}

async function fetchData() {
  const id = route.params.id as string
  if (id) {
    let tr = await secretApi.find(id);
    model.value = tr.data?.secret as Secret;
  }
}

onMounted(fetchData);
</script>