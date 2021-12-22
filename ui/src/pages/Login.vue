<template>
  <div :class="['container', isMobile ? '' : 'pc']">
    <img src="../assets/login.jpg" height="300" style="border-radius: 5px 0 0 5px" v-if="!isMobile" />
    <div class="form">
      <h1 class="title">Swirl</h1>
      <n-form
        :model="model"
        ref="form"
        :rules="rules"
        label-placement="left"
        @keydown.enter="submit"
      >
        <n-form-item path="name">
          <n-input round v-model:value="model.name" :placeholder="t('fields.login_name')" clearable>
            <template #prefix>
              <n-icon>
                <person-outline />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
        <n-form-item path="password">
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
        <n-button
          round
          block
          type="primary"
          :disabled="submiting"
          :loading="submiting"
          @click.prevent="submit"
        >{{ t('buttons.sign_in') }}</n-button>
      </n-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { NForm, NFormItem, NInput, NButton, NIcon } from "naive-ui";
import { PersonOutline, LockClosedOutline } from "@vicons/ionicons5";
import userApi from "@/api/user";
import type { AuthUser } from "@/api/user";
import systemApi from "@/api/system";
import type { LoginArgs } from "@/api/user";
import { useStore } from "vuex";
import { useIsMobile } from "@/utils";
import { Mutations } from "@/store/mutations";
import { useForm, requiredRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter();
const route = useRoute();
const store = useStore();
const isMobile = useIsMobile()
const form = ref();
const model = reactive({} as LoginArgs);
const rules = {
  name: requiredRule(),
  password: requiredRule(),
};
const { submit, submiting } = useForm<AuthUser>(form, () => userApi.login(model), (user: AuthUser) => {
  store.commit(Mutations.SetUser, user);
  let redirect = decodeURIComponent(<string>route.query.redirect || "/");
  router.push({ path: redirect });
})

async function checkState() {
  const r = await systemApi.checkState();
  if (r.data?.fresh) {
    router.push("/init")
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
      margin-top: -10px;
      text-align: center;
    }
  }
}
.pc {
  width: 500px;
  margin: 100px auto;
}
</style>