<template>
  <div :class="['container', isMobile ? '' : 'pc']">
    <div class="form">
      <h1 class="title">{{ t('texts.init_admin') }}</h1>
      <n-form :model="model" ref="form" :rules="rules" label-placement="top">
        <n-form-item path="loginName" :label="t('fields.login_name')">
          <n-input
            round
            v-model:value="model.loginName"
            :placeholder="t('fields.login_name')"
            clearable
          >
            <template #prefix>
              <n-icon>
                <person-outline />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
        <n-form-item path="password" :label="t('fields.password')" first>
          <n-input
            round
            v-model:value="model.password"
            type="password"
            :placeholder="t('fields.password')"
            clearable
          >
            <template #prefix>
              <n-icon>
                <lock-closed-outline />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
        <n-form-item path="name" :label="t('fields.username')">
          <n-input round v-model:value="model.name" :placeholder="t('fields.username')" clearable>
            <template #prefix>
              <n-icon>
                <person-outline />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
        <n-form-item path="email" :label="t('fields.email')">
          <n-input round v-model:value="model.email" :placeholder="t('fields.email')" clearable>
            <template #prefix>
              <n-icon>
                <mail-outline />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
        <n-button
          round
          block
          type="primary"
          :disabled="submiting"
          :loading="submiting"
          @click.prevent="submit"
        >{{ t('buttons.submit') }}</n-button>
      </n-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import { NForm, NFormItem, NInput, NButton, NIcon } from "naive-ui";
import { PersonOutline, LockClosedOutline, MailOutline } from "@vicons/ionicons5";
import systemApi from "@/api/system";
import { useIsMobile } from "@/utils";
import { useForm, requiredRule, emailRule, passwordRule, lengthRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter();
const isMobile = useIsMobile()
const form = ref();
const model = reactive({
  loginName: '',
  password: '',
  name: '',
  email: '',
});
const rules = {
  loginName: requiredRule(),
  password: [requiredRule(), passwordRule(), lengthRule(6, 15)],
  name: requiredRule(),
  email: [requiredRule(), emailRule()],
};
const { submit, submiting } = useForm(form, () => systemApi.createAdmin(model), () => {
  router.push({ name: 'login' });
})

async function checkState() {
  const r = await systemApi.checkState();
  if (!r.data?.fresh) {
    router.push({ name: 'login' });
  }
}

onMounted(checkState);
</script>

<style lang="scss" scoped>
.container {
  border-radius: 5px;
  box-shadow: 1px 1px 10px #ddd;
  display: flex;
  justify-content: center;
  align-items: center;
  .form {
    flex: 60%;
    padding: 20px;
    .title {
      margin-top: 0;
      text-align: center;
    }
  }
}
.pc {
  width: 500px;
  margin: 100px auto;
}
</style>