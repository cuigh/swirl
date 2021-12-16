<template>
  <x-page-header :subtitle="user.id">
    <template #action>
      <n-button secondary size="small" @click="$router.push('/system/users')">
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
    <n-form :model="user" :rules="rules" ref="form" label-placement="top">
      <n-grid cols="1 640:2" :x-gap="24">
        <n-form-item-gi :label="t('fields.username')" path="name">
          <n-input :placeholder="t('fields.username')" v-model:value="user.name" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.login_name')" path="loginName">
          <n-input :placeholder="t('fields.login_name')"  v-model:value="user.loginName" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.password')" path="password" v-if="!user.id">
          <n-input type="password" :placeholder="t('fields.password')"  v-model:value="user.password" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.password_confirm')" path="passwordConfirm" v-if="!user.id">
          <n-input type="password" :placeholder="t('fields.password_confirm')"  v-model:value="user.passwordConfirm" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.email')" path="email">
          <n-input :placeholder="t('fields.email')"  v-model:value="user.email" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.admin')" path="admin">
          <n-switch v-model:value="user.admin">
            <template #checked>{{ t('enums.yes') }}</template>
            <template #unchecked>{{ t('enums.no') }}</template>
          </n-switch>
        </n-form-item-gi>
        <n-form-item-gi
          :label="t('fields.type')"
          path="type"
          label-placement="left"
          label-width="41"
          v-if="user.id"
        >
          <n-radio-group v-model:value="user.type">
            <n-radio key="internal" value="internal">Internal</n-radio>
            <n-radio key="ldap" value="ldap">LDAP</n-radio>
          </n-radio-group>
        </n-form-item-gi>
        <n-form-item-gi :label="t('objects.role', 2)" span="2" path="roles" v-if="roles && roles.length">
          <n-checkbox-group v-model:value="user.roles">
            <n-space item-style="display: flex;">
              <n-checkbox :value="r.id" :label="r.name" v-for="r of roles" />
            </n-space>
          </n-checkbox-group>
        </n-form-item-gi>
        <n-gi :span="2">
          <n-button
            :disabled="submiting"
            :loading="submiting"
            @click.prevent="submit"
            type="primary"
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
import { onMounted, ref } from "vue";
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
  NCheckboxGroup,
  NCheckbox,
  NRadioGroup,
  NRadio,
} from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import { useRoute } from "vue-router";
import { router } from "@/router/router";
import userApi from "@/api/user";
import roleApi from "@/api/role";
import type { User } from "@/api/user";
import type { Role } from "@/api/role";
import { useForm, emailRule, requiredRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const user = ref({ type: 'internal', admin: false } as User)
const roles = ref([] as Role[]);
const rules: any = {
  name: requiredRule(),
  loginName: requiredRule(),
  email: [requiredRule(), emailRule()],
  password: requiredRule(),
  passwordConfirm: requiredRule(),
};
const form = ref();
const { submit, submiting } = useForm(form, () => userApi.save(user.value), () => {
  window.message.info(t('texts.action_success'));
  router.push("/system/users")
})

async function fetchData() {
  const id = route.params.id as string || ''
  if (id) {
    let r = await userApi.find(id);
    user.value = r.data as User;
  }
  let r = await roleApi.search()
  roles.value = r.data as Role[]
}

onMounted(fetchData);
</script>