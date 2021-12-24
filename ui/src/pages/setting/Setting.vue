<template>
  <x-page-header />
  <n-space class="page-body" vertical :size="12">
    <x-panel title="LDAP" :subtitle="t('tips.ldap')" divider="bottom" :collapsed="panel !== 'ldap'">
      <template #action>
        <n-button
          secondary
          strong
          class="toggle"
          size="small"
          @click="togglePanel('ldap')"
        >{{ panel === 'ldap' ? t('buttons.collapse') : t('buttons.expand') }}</n-button>
      </template>
      <n-form
        :model="setting"
        ref="formLdap"
        label-placement="left"
        style="padding: 4px 0 0 12px"
        label-width="100"
      >
        <n-form-item :label="t('fields.enabled')" path="ldap.enabled" label-align="right">
          <n-switch v-model:value="setting.ldap.enabled" />
        </n-form-item>
        <n-form-item :label="t('fields.address')" path="ldap.address" label-align="right">
          <n-input :placeholder="t('tips.ldap_address')" v-model:value="setting.ldap.address" />
        </n-form-item>
        <n-form-item :label="t('fields.security')" path="ldap.security">
          <n-radio-group v-model:value="setting.ldap.security">
            <n-radio :value="0">None</n-radio>
            <n-radio :value="1">TLS</n-radio>
            <n-radio :value="2">StartTLS</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item :label="t('fields.authentication')" path="ldap.auth">
          <n-radio-group v-model:value="setting.ldap.auth">
            <n-radio value="simple">{{ t('enums.simple') }}</n-radio>
            <n-radio value="bind">{{ t('enums.bind') }}</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item
          :label="t('fields.user_dn')"
          path="ldap.user_dn"
          label-align="right"
          v-show="setting.ldap.auth === 'simple'"
        >
          <n-input :placeholder="t('tips.ldap_user_dn')" v-model:value="setting.ldap.user_dn" />
        </n-form-item>
        <n-form-item
          :label="t('fields.bind_dn')"
          label-align="right"
          :show-feedback="false"
          v-show="setting.ldap.auth === 'bind'"
        >
          <n-grid :cols="2" :x-gap="24">
            <n-form-item-gi path="ldap.bind_dn">
              <n-input-group>
                <n-input-group-label style="min-width: 60px">{{ t('fields.dn') }}</n-input-group-label>
                <n-input
                  :placeholder="t('tips.ldap_bind_dn')"
                  v-model:value="setting.ldap.bind_dn"
                />
              </n-input-group>
            </n-form-item-gi>
            <n-form-item-gi path="ldap.bind_pwd">
              <n-input-group>
                <n-input-group-label style="min-width: 60px">{{ t('fields.password') }}</n-input-group-label>
                <n-input
                  type="password"
                  :placeholder="t('tips.ldap_bind_pwd')"
                  v-model:value="setting.ldap.bind_pwd"
                />
              </n-input-group>
            </n-form-item-gi>
          </n-grid>
        </n-form-item>
        <n-form-item :label="t('fields.base_dn')" path="ldap.base_dn" label-align="right">
          <n-input :placeholder="t('tips.ldap_base_dn')" v-model:value="setting.ldap.base_dn" />
        </n-form-item>
        <n-form-item :label="t('fields.user_filter')" path="ldap.user_filter" label-align="right">
          <n-input
            :placeholder="t('tips.ldap_user_filter')"
            v-model:value="setting.ldap.user_filter"
          />
        </n-form-item>
        <n-form-item :label="t('fields.attr_map')" label-align="right" :show-feedback="false">
          <n-grid :cols="2" :x-gap="24">
            <n-form-item-gi path="ldap.name_attr">
              <n-input-group>
                <n-input-group-label style="min-width: 80px">{{ t('fields.username') }}</n-input-group-label>
                <n-input placeholder="e.g. displayName" v-model:value="setting.ldap.name_attr" />
              </n-input-group>
            </n-form-item-gi>
            <n-form-item-gi path="ldap.email_attr">
              <n-input-group>
                <n-input-group-label style="min-width: 80px">{{ t('fields.email') }}</n-input-group-label>
                <n-input placeholder="e.g. mail" v-model:value="setting.ldap.email_attr" />
              </n-input-group>
            </n-form-item-gi>
          </n-grid>
        </n-form-item>
        <n-button type="primary" @click="() => save('ldap', setting.ldap)">{{ t('buttons.save') }}</n-button>
      </n-form>
    </x-panel>
    <x-panel
      :title="t('fields.monitor')"
      :subtitle="t('tips.monitor')"
      :collapsed="panel !== 'metric'"
    >
      <template #action>
        <n-button
          secondary
          strong
          class="toggle"
          size="small"
          @click="togglePanel('metric')"
        >{{ panel === 'metric' ? t('buttons.collapse') : t('buttons.expand') }}</n-button>
      </template>
      <n-form
        :model="setting"
        ref="formMetrics"
        label-placement="left"
        style="padding: 4px 0 0 12px"
      >
        <n-form-item label="Prometheus" path="metric.prometheus" label-align="right">
          <n-input :placeholder="t('tips.prometheus')" v-model:value="setting.metric.prometheus" />
        </n-form-item>
        <n-button
          type="primary"
          @click="() => save('metric', setting.metric)"
        >{{ t('buttons.save') }}</n-button>
      </n-form>
    </x-panel>
    <n-alert type="info">{{ t('texts.setting_notice') }}</n-alert>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NGrid,
  NButton,
  NSpace,
  NInput,
  NInputGroup,
  NInputGroupLabel,
  NForm,
  NFormItem,
  NFormItemGi,
  NRadioGroup,
  NRadio,
  NSwitch,
  NAlert,
} from "naive-ui";
import XPageHeader from "@/components/PageHeader.vue";
import XPanel from "@/components/Panel.vue";
import settingApi from "@/api/setting";
import type { Setting } from "@/api/setting";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const setting = ref({
  ldap: {
    security: 0,
    auth: 'simple',
  },
  metric: {},
  deploy: {},
} as Setting);
const panel = ref('')

function togglePanel(name: string) {
  if (panel.value === name) {
    panel.value = ''
  } else {
    panel.value = name
  }
}

async function save(id: string, options: any) {
  await settingApi.save(id, options)
  window.message.info(t('texts.action_success'));
}

async function fetchData() {
  let r = (await settingApi.load()).data as Setting;
  setting.value = Object.assign(setting.value, r)
}

onMounted(fetchData);
</script>

<style scoped>
.toggle {
  width: 75px;
}
</style>
