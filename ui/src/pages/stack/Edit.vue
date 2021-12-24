<template>
  <x-page-header :subtitle="model.name">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'stack_list' })">
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
        <n-form-item-gi :label="t('fields.name')" path="name" span="2">
          <n-input
            :placeholder="t('fields.name')"
            v-model:value="model.name"
            :disabled="Boolean(model.id)"
          />
        </n-form-item-gi>
        <n-form-item-gi :show-label="false">
          <n-radio-group v-model:value="type">
            <n-radio key="input" value="input">{{ t('enums.input') }}</n-radio>
            <n-radio key="upload" value="upload">{{ t('enums.upload') }}</n-radio>
          </n-radio-group>
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.content')" path="content" span="2">
          <!-- <n-input
            type="textarea"
            :placeholder="t('fields.content')"
            v-model:value="model.content"
            :autosize="{ minRows: 10, maxRows: 30 }"
            v-if="type === 'input'"
          /> -->
          <x-code-mirror v-model="model.content" v-show="type === 'input'"/>
          <n-upload
            action="/api/stack/upload"
            :default-upload="true"
            :show-file-list="false"
            :multiple="false"
            :max="1"
            name="content"
            @finish="finishUpload"
            @before-upload="beforeUpload"
            v-show="type === 'upload'"
          >
            <n-upload-dragger style="padding: 12px">
              <div>
                <n-icon size="48" :depth="3">
                  <document-icon />
                </n-icon>
              </div>
              <n-text style="font-size: 14px" depth="3">{{ t('tips.upload') }}</n-text>
              <n-space v-if="fileName" :size="0" justify="center">
                <n-icon size="24" :depth="3">
                  <attach-icon />
                </n-icon>
                <n-text style="font-size: 16px">{{ fileName }}</n-text>
              </n-space>
            </n-upload-dragger>
          </n-upload>
        </n-form-item-gi>
      </n-grid>
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
  NRadioGroup,
  NRadio,
  NUpload,
  NUploadDragger,
  NText,
} from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
  DocumentTextOutline as DocumentIcon,
  AttachOutline as AttachIcon,
} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XCodeMirror from "@/components/CodeMirror.vue";
import stackApi from "@/api/stack";
import type { Stack } from "@/api/stack";
import { useRoute } from "vue-router";
import { router } from "@/router/router";
import { requiredRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({} as Stack);
const type = ref("input")
const rules: any = {
  name: requiredRule(),
  content: requiredRule(),
};
const form = ref();
const fileName = ref('')
const submiting = ref(false)

async function submit(e: Event) {
  e.preventDefault();
  form.value.validate(async (errors: any) => {
    if (errors) {
      return
    }

    submiting.value = true;
    try {
      await stackApi.save(model.value)
      window.message.info(t('texts.action_success'));
      router.push({ name: 'stack_list' })
    } finally {
      submiting.value = false;
    }
  });
}

function beforeUpload({ file, files }: any): any {
  fileName.value = ''
  model.value.content = ''
}

function finishUpload({ file, event }: any): any {
  fileName.value = file.name
  model.value.content = event.currentTarget.responseText
}

async function fetchData() {
  const name = route.params.name as string
  if (name) {
    let tr = await stackApi.find(name);
    model.value = tr.data as Stack;
    model.value.id = model.value.name
  }
}

onMounted(fetchData);
</script>