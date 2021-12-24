<template>
  <x-page-header :subtitle="model.name">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'role_list' })">
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
    <n-form :model="model" :rules="rules" ref="form" label-placement="top">
      <n-grid cols="1 640:2" :x-gap="24">
        <n-form-item-gi :label="t('fields.name')" path="name">
          <n-input :placeholder="t('fields.name')" v-model:value="model.name" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.desc')" path="desc" span="2">
          <n-input :placeholder="t('fields.desc')" v-model:value="model.desc" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.perms')" span="2" path="perms">
          <n-checkbox-group v-model:value="model.perms" style="display: contents">
            <n-table :single-line="false" size="small">
              <tr v-for="g in perms">
                <td width="75" style="font-weight: 500">{{ t('objects.' + g.key) }}</td>
                <td>
                  <n-space item-style="display: flex;">
                    <n-checkbox
                      :value="g.key + '.' + action"
                      :label="t('perms.' + action)"
                      v-for="action in g.actions"
                    />
                  </n-space>
                </td>
                <td width="100">
                  <n-space :size="4">
                    <n-button
                      quaternary
                      type="info"
                      size="tiny"
                      @click="() => checkGroup(g.key)"
                    >{{ t('buttons.check') }}</n-button>
                    <n-button
                      quaternary
                      type="info"
                      size="tiny"
                      @click="() => checkGroup(g.key, false)"
                    >{{ t('buttons.cancel') }}</n-button>
                  </n-space>
                </td>
              </tr>
            </n-table>
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
  NCheckboxGroup,
  NCheckbox,
  NTable,
} from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
} from "@vicons/ionicons5";
import { useRoute } from "vue-router";
import XPageHeader from "@/components/PageHeader.vue";
import { router } from "@/router/router";
import roleApi from "@/api/role";
import type { Role } from "@/api/role";
import { perms } from "@/utils/perm";
import { useForm, requiredRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({ perms: [] as string[] } as Role);
const rules: any = {
  name: requiredRule(),
};
const form = ref();
const { submit, submiting } = useForm(form, () => roleApi.save(model.value), () => {
  window.message.info(t('texts.action_success'));
  router.push({ name: 'role_list' })
})

function checkGroup(key: string, checked: boolean = true) {
  const g = perms.find(g => g.key === key)
  if (!g) {
    return
  }

  g.actions.forEach(action => {
    const perm = g.key + '.' + action
    if (checked) {
      !model.value.perms.includes(perm) && model.value.perms.push(perm)
    } else {
      const index = model.value.perms.indexOf(perm)
      if (index !== -1) {
        model.value.perms.splice(index, 1)
      }
    }
  })
}

async function fetchData() {
  let id = route.params.id as string
  if (id) {
    let r = await roleApi.find(id);
    model.value = r.data as Role
  }
}

onMounted(fetchData);
</script>
