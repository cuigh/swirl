<template>
  <x-page-header />
  <n-space class="page-body" vertical :size="12">
    <x-panel
      :title="t('fields.profile')"
      :subtitle="t('tips.profile')"
      divider="bottom"
      :collapsed="panel !== 'profile'"
    >
      <template #action>
        <n-button
          secondary
          strong
          class="toggle"
          size="small"
          @click="togglePanel('profile')"
        >{{ panel === 'profile' ? t('buttons.collapse') : t('buttons.expand') }}</n-button>
      </template>
      <div style="padding: 4px 0 0 12px">
        <n-form inline :model="profile" :rules="profileRules" ref="profileForm">
          <n-grid cols="1 640:2" :x-gap="24">
            <n-form-item-gi :label="t('fields.username')" path="name">
              <n-input :placeholder="t('fields.username')" v-model:value="profile.name" />
            </n-form-item-gi>
            <n-form-item-gi :label="t('fields.login_name')" path="loginName">
              <n-input :placeholder="t('fields.login_name')" v-model:value="profile.loginName" />
            </n-form-item-gi>
            <n-form-item-gi :label="t('fields.email')" path="email">
              <n-input :placeholder="t('fields.email')" v-model:value="profile.email" />
            </n-form-item-gi>
            <n-form-item-gi :label="t('fields.tokens')" path="tokens" span="2">
              <n-dynamic-input
                v-model:value="profile.tokens"
                #="{ index, value }"
                :on-create="() => ({ name: '', value: guid() })"
              >
                <n-input
                  :placeholder="t('fields.name')"
                  v-model:value="value.name"
                  style="width: 300px"
                />
                <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
                <n-input-group>
                  <n-input :placeholder="t('fields.value')" v-model:value="value.value" readonly></n-input>
                  <n-tooltip trigger="hover">
                    <template #trigger>
                      <n-button
                        type="default"
                        #icon
                        @click="() => copy(value.value)"
                        v-if="isSupported"
                      >
                        <n-icon>
                          <copy-icon />
                        </n-icon>
                      </n-button>
                    </template>
                    {{ t(copied ? 'tips.copied' : 'buttons.copy') }}
                  </n-tooltip>
                </n-input-group>
              </n-dynamic-input>
            </n-form-item-gi>
          </n-grid>
        </n-form>
        <n-button
          type="primary"
          :disabled="profileSubmiting"
          :loading="profileSubmiting"
          @click="modifyProfile"
        >
          <template #icon>
            <n-icon>
              <save-icon />
            </n-icon>
          </template>
          {{ t('buttons.update') }}
        </n-button>
      </div>
    </x-panel>
    <x-panel
      :title="t('fields.password')"
      :subtitle="t('tips.password')"
      divider="bottom"
      :collapsed="panel !== 'password'"
    >
      <template #action>
        <n-button
          secondary
          strong
          size="small"
          class="toggle"
          @click="togglePanel('password')"
        >{{ panel === 'password' ? t('buttons.collapse') : t('buttons.expand') }}</n-button>
      </template>
      <div style="padding: 4px 0 0 12px" v-if="profile.type === 'ldap'">
        <n-alert type="info">{{ t('texts.password_notice') }}</n-alert>
      </div>
      <div style="padding: 4px 0 0 12px" v-else>
        <n-form :model="password" ref="passwordForm" :rules="passwordRules">
          <n-grid cols="1 640:3" :x-gap="24">
            <n-form-item-gi path="old" :label="t('fields.password_old')">
              <n-input
                v-model:value="password.old"
                type="password"
                :placeholder="t('fields.password_old')"
              />
            </n-form-item-gi>
            <n-form-item-gi first path="new" :label="t('fields.password_new')">
              <n-input
                v-model:value="password.new"
                type="password"
                :placeholder="t('fields.password_new')"
              />
            </n-form-item-gi>
            <n-form-item-gi first path="confirm" :label="t('fields.password_confirm')">
              <n-input
                :disabled="!password.new"
                v-model:value="password.confirm"
                type="password"
                :placeholder="t('fields.password_confirm')"
              />
            </n-form-item-gi>
          </n-grid>
        </n-form>
        <n-button
          type="primary"
          :disabled="passwordSubmiting"
          :loading="passwordSubmiting"
          @click="modifyPassword"
        >
          <template #icon>
            <n-icon>
              <save-icon />
            </n-icon>
          </template>
          {{ t('buttons.update') }}
        </n-button>
      </div>
    </x-panel>
    <x-panel
      :title="t('fields.preferences')"
      :subtitle="t('tips.preference')"
      :collapsed="panel !== 'preference'"
    >
      <template #action>
        <n-button
          secondary
          strong
          class="toggle"
          size="small"
          @click="togglePanel('preference')"
        >{{ panel === 'preference' ? t('buttons.collapse') : t('buttons.expand') }}</n-button>
      </template>
      <div style="padding: 4px 0 0 12px">
        <n-form inline :model="preference" ref="preferenceForm" label-placement="left">
          <n-form-item :label="t('fields.language')" path="locale">
            <n-radio-group v-model:value="preference.locale">
              <n-radio-button value="zh">中文</n-radio-button>
              <n-radio-button value="en">English</n-radio-button>
            </n-radio-group>
          </n-form-item>
          <n-form-item :label="t('fields.theme')" path="theme">
            <n-radio-group v-model:value="preference.theme">
              <n-radio-button value="light">{{ t('enums.light') }}</n-radio-button>
              <n-radio-button value="dark">{{ t('enums.dark') }}</n-radio-button>
            </n-radio-group>
          </n-form-item>
        </n-form>
        <n-button type="primary" @click="savePreference">
          <template #icon>
            <n-icon>
              <save-icon />
            </n-icon>
          </template>
          {{ t('buttons.save') }}
        </n-button>
      </div>
    </x-panel>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import {
  NButton,
  NSpace,
  NInput,
  NInputGroup,
  NIcon,
  NForm,
  NFormItem,
  NFormItemGi,
  NGrid,
  NRadioButton,
  NRadioGroup,
  NAlert,
  NDynamicInput,
  NTooltip,
} from "naive-ui";
import {
  SaveOutline as SaveIcon,
  CopyOutline as CopyIcon,
} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XPanel from "@/components/Panel.vue";
import userApi from "@/api/user";
import type { User } from "@/api/user";
import { useForm, emailRule, requiredRule, customRule, lengthRule } from "@/utils/form";
import { Mutations } from "@/store/mutations";
import { useStore } from "vuex";
import { useI18n } from 'vue-i18n'
import { useClipboard } from '@vueuse/core'
import { guid } from "@/utils";

const { t } = useI18n()
const panel = ref('')
function togglePanel(name: string) {
  if (panel.value === name) {
    panel.value = ''
  } else {
    panel.value = name
  }
}

// profile
const profile = ref({} as User)
const profileRules: any = {
  name: requiredRule(),
  loginName: requiredRule(),
  email: [requiredRule(), emailRule()],
};
const profileForm = ref();
const { submit: modifyProfile, submiting: profileSubmiting } = useForm(profileForm, () => userApi.modifyProfile(profile.value))
const { copy, copied, isSupported } = useClipboard()

// password
const password = reactive({
  showDlg: false,
  old: "",
  new: "",
  confirm: "",
})
const passwordForm = ref()
const passwordRules = {
  old: requiredRule(),
  new: [
    requiredRule(),
    lengthRule(6, 15),
    customRule((rule: any, value: string) => value !== password.old, t('tips.password_new_rule')),
  ],
  confirm: [
    requiredRule(),
    customRule((rule: any, value: string) => value === password.new, t('tips.password_confirm_rule')),
  ],
};
const { submit: modifyPassword, submiting: passwordSubmiting } = useForm(
  passwordForm,
  () => userApi.modifyPassword({ oldPwd: password.old, newPwd: password.new }),
  () => {
    password.showDlg = false
    window.message.info(t('texts.action_success'));
  }
);

// preference
const store = useStore();
const preference = reactive({
  locale: store.state.preference.locale,
  theme: store.state.preference.theme || 'light',
})
const preferenceForm = ref()
function savePreference() {
  store.commit(Mutations.SetPreference, preference)
  window.message.info(t('texts.action_success'));
  setTimeout(() => location.reload(), 100)
}

async function fetchData() {
  let r = await userApi.find('');
  profile.value = r.data as User;
}

onMounted(fetchData);
</script>

<style scoped>
.toggle {
  width: 75px;
}
</style>
