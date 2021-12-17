<template>
  <x-page-header :subtitle="model.name">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'service_list' })">
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
        <n-form-item-gi :label="t('objects.image')" path="image">
          <n-input :placeholder="t('objects.image')" v-model:value="model.image" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.mode')" path="mode">
          <n-radio-group v-model:value="model.mode" :disabled="Boolean(model.id)">
            <n-radio key="replicated" value="replicated">Replicated</n-radio>
            <n-radio key="global" value="global">Global</n-radio>
            <n-radio key="replicated-job" value="replicated-job">Replicated Job</n-radio>
            <n-radio key="global-job" value="global-job">Global Job</n-radio>
          </n-radio-group>
        </n-form-item-gi>
        <n-form-item-gi
          :label="t('fields.replicas')"
          path="replicas"
          v-if="['replicated', 'replicated-job'].includes(model.mode)"
        >
          <n-input-number
            :placeholder="t('fields.replicas')"
            :min="0"
            v-model:value="model.replicas"
            style="max-width: 120px"
          />
        </n-form-item-gi>
        <n-form-item-gi :label="t('objects.network', 2)" span="2" path="networks">
          <n-checkbox-group v-model:value="model.networks">
            <n-space item-style="display: flex;">
              <n-checkbox :value="n.id" :label="n.name" v-for="n of networks" />
            </n-space>
          </n-checkbox-group>
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.command')" path="command">
          <n-input :placeholder="t('tips.command')" v-model:value="model.command" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.args')" path="command">
          <n-input :placeholder="t('tips.args')" v-model:value="model.args" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.work_dir')" path="w">
          <n-input :placeholder="t('fields.work_dir')" v-model:value="model.dir" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.user')" path="command">
          <n-input :placeholder="t('fields.user')" v-model:value="model.user" />
        </n-form-item-gi>
        <n-form-item-gi span="2" :label="t('fields.env')" path="env">
          <n-dynamic-input v-model:value="model.env" #="{ index, value }" :on-create="newLabel">
            <n-input :placeholder="t('fields.name')" v-model:value="value.name" />
            <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
            <n-input :placeholder="t('fields.value')" v-model:value="value.value" />
          </n-dynamic-input>
        </n-form-item-gi>
        <n-form-item-gi span="2" :label="t('fields.service_labels')" path="labels">
          <n-dynamic-input v-model:value="model.labels" #="{ index, value }" :on-create="newLabel">
            <n-input :placeholder="t('fields.name')" v-model:value="value.name" />
            <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
            <n-input :placeholder="t('fields.name')" v-model:value="value.value" />
          </n-dynamic-input>
        </n-form-item-gi>
      </n-grid>
      <x-group :title="t('fields.endpoint')">
        <n-grid cols="1">
          <n-form-item-gi
            :label="t('fields.resolution_mode')"
            path="endpoint.mode"
            label-placement="left"
            :show-feedback="false"
          >
            <n-radio-group v-model:value="model.endpoint.mode">
              <n-radio key="vip" value="vip">VIP</n-radio>
              <n-radio key="dnsrr" value="dnsrr">DNSRR</n-radio>
            </n-radio-group>
          </n-form-item-gi>
          <n-form-item-gi
            :label="t('fields.ports')"
            :show-label="false"
            path="endpoint.ports"
            label-width="70"
          >
            <n-dynamic-input
              v-model:value="model.endpoint.ports"
              #="{ index, value }"
              :on-create="newPort"
            >
              <n-input-group>
                <n-input
                  :placeholder="t('fields.name')"
                  v-model:value="value.name"
                  style="width: 20%"
                />
                <n-select
                  :placeholder="t('fields.mode')"
                  v-model:value="value.mode"
                  :options="[{ label: 'Ingress', value: 'ingress' }, { label: 'Host', value: 'host' }]"
                  style="width: 20%"
                />
                <n-select
                  :placeholder="t('fields.protocol')"
                  v-model:value="value.protocol"
                  :options="[{ label: 'TCP', value: 'tcp' }, { label: 'UDP', value: 'udp' }, { label: 'SCTP', value: 'sctp' }]"
                  style="width: 20%"
                />
                <n-input-number
                  :show-button="false"
                  :min="0"
                  :max="65535"
                  v-model:value="value.targetPort"
                  style="width: 20%"
                >
                  <template #suffix>
                    <n-text depth="3">{{ t('fields.target_port') }}</n-text>
                  </template>
                </n-input-number>
                <n-input-number
                  :show-button="false"
                  :min="0"
                  :max="65535"
                  v-model:value="value.publishedPort"
                  style="width: 20%"
                >
                  <template #suffix>
                    <n-text depth="3">{{ t('fields.published_port') }}</n-text>
                  </template>
                </n-input-number>
              </n-input-group>
            </n-dynamic-input>
          </n-form-item-gi>
        </n-grid>
      </x-group>
      <x-panel :title="t('fields.more')" :collapsed="!showMore">
        <template #action>
          <n-button
            secondary
            size="small"
            @click="showMore = !showMore"
          >{{ showMore ? t('buttons.collapse') : t('buttons.expand') }}</n-button>
        </template>
        <n-tabs type="line" style="margin-top: -12px">
          <n-tab-pane name="mounts" :tab="t('fields.mounts')" display-directive="show">
            <n-form-item :show-label="false" path="mounts">
              <n-dynamic-input
                v-model:value="model.mounts"
                #="{ index, value }"
                :on-create="newMount"
              >
                <n-input-group>
                  <n-select
                    :placeholder="t('fields.type')"
                    v-model:value="value.type"
                    :options="mountTypeOptions"
                    style="width: 20%"
                  />
                  <n-input
                    :placeholder="t('fields.source_path')"
                    v-model:value="value.source"
                    style="width: 30%"
                  />
                  <n-input
                    :placeholder="t('fields.target_path')"
                    v-model:value="value.target"
                    style="width: 30%"
                  />
                  <n-select
                    :placeholder="t('fields.consistency')"
                    v-model:value="value.consistency"
                    :options="consistencyOptions"
                    style="width: 20%; margin-right: 10px"
                  />
                </n-input-group>
                <n-checkbox v-model:checked="value.readonly">{{ t('fields.readonly') }}</n-checkbox>
              </n-dynamic-input>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="files" :tab="t('fields.files')" display-directive="show">
            <n-form-item :label="t('objects.config', 2)" path="configs">
              <n-dynamic-input
                v-model:value="model.configs"
                #="{ index, value }"
                :on-create="newFile"
              >
                <n-input-group>
                  <n-select
                    filterable
                    :placeholder="t('fields.files')"
                    v-model:value="value.key"
                    style="width: 20%"
                    :options="configFiles"
                  />
                  <n-input
                    :placeholder="t('fields.target_path')"
                    v-model:value="value.path"
                    style="width: 35%"
                  />
                  <n-input placeholder="UID" v-model:value="value.uid" style="width: 15%">
                    <template #suffix>
                      <n-text depth="3">UID</n-text>
                    </template>
                  </n-input>
                  <n-input placeholder="GID" v-model:value="value.gid" style="width: 15%">
                    <template #suffix>
                      <n-text depth="3">GID</n-text>
                    </template>
                  </n-input>
                  <n-input-number
                    :show-button="false"
                    :placeholder="t('fields.mode')"
                    v-model:value="value.mode"
                    style="width: 15%"
                  >
                    <template #suffix>
                      <n-text depth="3">Mode</n-text>
                    </template>
                  </n-input-number>
                </n-input-group>
              </n-dynamic-input>
            </n-form-item>
            <n-form-item :label="t('objects.secret', 2)" path="secrets">
              <n-dynamic-input
                v-model:value="model.secrets"
                #="{ index, value }"
                :on-create="newFile"
              >
                <n-input-group>
                  <n-select
                    filterable
                    :placeholder="t('fields.files')"
                    v-model:value="value.key"
                    style="width: 20%"
                    :options="secretFiles"
                  />
                  <n-input
                    :placeholder="t('fields.target_path')"
                    v-model:value="value.path"
                    style="width: 35%"
                  />
                  <n-input placeholder="UID" v-model:value="value.uid" style="width: 15%">
                    <template #suffix>
                      <n-text depth="3">UID</n-text>
                    </template>
                  </n-input>
                  <n-input placeholder="GID" v-model:value="value.gid" style="width: 15%">
                    <template #suffix>
                      <n-text depth="3">GID</n-text>
                    </template>
                  </n-input>
                  <n-input-number
                    :show-button="false"
                    :placeholder="t('fields.mode')"
                    v-model:value="value.mode"
                    style="width: 15%"
                  >
                    <template #suffix>
                      <n-text depth="3">Mode</n-text>
                    </template>
                  </n-input-number>
                </n-input-group>
              </n-dynamic-input>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="resources" :tab="t('fields.resources')" display-directive="show">
            <x-group :title="t('fields.limit')">
              <n-grid :cols="2" :x-gap="12">
                <n-form-item-gi :label="t('fields.cpu')" path="resource.limit.cpu">
                  <n-input-number
                    :step="0.1"
                    :min="0"
                    placeholder="e.g. 0.5"
                    v-model:value="model.resource.limit.cpu"
                    style="width: 100%"
                  />
                </n-form-item-gi>
                <n-form-item-gi :label="t('fields.memory')" path="esource.limit.memory">
                  <n-input placeholder="e.g. 2G" v-model:value="model.resource.limit.memory" />
                </n-form-item-gi>
              </n-grid>
            </x-group>
            <x-group :title="t('fields.reserve')">
              <n-grid :cols="2" :x-gap="12">
                <n-form-item-gi :label="t('fields.cpu')" path="resource.reserve.cpu">
                  <n-input-number
                    :step="0.1"
                    :min="0"
                    placeholder="e.g. 0.1"
                    v-model:value="model.resource.reserve.cpu"
                    style="width: 100%"
                  />
                </n-form-item-gi>
                <n-form-item-gi :label="t('fields.memory')" path="esource.reserve.memory">
                  <n-input placeholder="e.g. 100M" v-model:value="model.resource.reserve.memory" />
                </n-form-item-gi>
              </n-grid>
            </x-group>
          </n-tab-pane>
          <n-tab-pane name="placement" :tab="t('fields.placement')" display-directive="show">
            <n-form-item :label="t('fields.constraints')" path="placement.constraints">
              <n-dynamic-input
                v-model:value="model.placement.constraints"
                #="{ index, value }"
                :on-create="newConstraint"
              >
                <n-input-group>
                  <n-input :placeholder="t('tips.constraint_name')" v-model:value="value.name" />
                  <n-select
                    v-model:value="value.op"
                    :options="[{ label: '==', value: '==' }, { label: '!=', value: '!=' }]"
                    style="width: 20%"
                  />
                  <n-input :placeholder="t('tips.constraint_value')" v-model:value="value.value" />
                </n-input-group>
              </n-dynamic-input>
            </n-form-item>
            <n-form-item :label="t('fields.preferences')" path="constraints.preferences">
              <n-dynamic-input
                v-model:value="model.placement.preferences"
                placeholder="e.g. engine.labels.az"
              />
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="schedule" :tab="t('fields.schedule')" display-directive="show">
            <x-group :title="t('fields.update')">
              <n-grid cols="1 640:2 960:4" :x-gap="12">
                <n-form-item-gi :label="t('fields.parallelism')" path="updatePolicy.parallelism">
                  <n-input-number
                    :min="0"
                    placeholder="Parallelism"
                    v-model:value="model.updatePolicy.parallelism"
                    style="width: 100%"
                  />
                </n-form-item-gi>
                <n-form-item-gi :label="t('fields.delay')" path="updatePolicy.delay">
                  <n-input placeholder="e.g. 30s" v-model:value="model.updatePolicy.delay" />
                </n-form-item-gi>
                <n-form-item-gi
                  :label="t('fields.failure_action')"
                  path="updatePolicy.failureAction"
                >
                  <n-radio-group v-model:value="model.updatePolicy.failureAction">
                    <n-radio key="pause" value="pause">{{ t('enums.pause') }}</n-radio>
                    <n-radio key="continue" value="continue">{{ t('enums.continue') }}</n-radio>
                    <n-radio key="rollback" value="rollback">{{ t('enums.rollback') }}</n-radio>
                  </n-radio-group>
                </n-form-item-gi>
                <n-form-item-gi :label="t('fields.order')" path="updatePolicy.order">
                  <n-radio-group v-model:value="model.updatePolicy.order">
                    <n-radio key="start-first" value="start-first">{{ t('enums.start_first') }}</n-radio>
                    <n-radio key="stop-first" value="stop-first">{{ t('enums.stop_first') }}</n-radio>
                  </n-radio-group>
                </n-form-item-gi>
              </n-grid>
            </x-group>
            <x-group :title="t('fields.rollback')">
              <n-grid cols="1 640:2 960:4" :x-gap="12">
                <n-form-item-gi :label="t('fields.parallelism')" path="rollbackPolicy.parallelism">
                  <n-input-number
                    :min="0"
                    placeholder="Parallelism"
                    v-model:value="model.rollbackPolicy.parallelism"
                    style="width: 100%"
                  />
                </n-form-item-gi>
                <n-form-item-gi :label="t('fields.delay')" path="rollbackPolicy.delay">
                  <n-input placeholder="e.g. 30s" v-model:value="model.rollbackPolicy.delay" />
                </n-form-item-gi>
                <n-form-item-gi
                  :label="t('fields.failure_action')"
                  path="rollbackPolicy.failureAction"
                >
                  <n-radio-group v-model:value="model.rollbackPolicy.failureAction">
                    <n-radio key="pause" value="pause">{{ t('enums.pause') }}</n-radio>
                    <n-radio key="continue" value="continue">{{ t('enums.continue') }}</n-radio>
                  </n-radio-group>
                </n-form-item-gi>
                <n-form-item-gi :label="t('fields.order')" path="rollbackPolicy.order">
                  <n-radio-group v-model:value="model.rollbackPolicy.order">
                    <n-radio key="start-first" value="start-first">{{ t('enums.start_first') }}</n-radio>
                    <n-radio key="stop-first" value="stop-first">{{ t('enums.stop_first') }}</n-radio>
                  </n-radio-group>
                </n-form-item-gi>
              </n-grid>
            </x-group>
            <x-group :title="t('fields.restart')">
              <n-grid cols="1 640:2 960:4" :x-gap="12">
                <n-form-item-gi :label="t('fields.max_attempts')" path="restartPolicy.maxAttempts">
                  <n-input-number
                    :min="0"
                    placeholder="Max attempts"
                    v-model:value="model.restartPolicy.maxAttempts"
                    style="width: 100%"
                  />
                </n-form-item-gi>
                <n-form-item-gi :label="t('fields.delay')" path="restartPolicy.delay">
                  <n-input placeholder="e.g. 30s" v-model:value="model.restartPolicy.delay" />
                </n-form-item-gi>
                <n-form-item-gi :label="t('fields.window')" path="restartPolicy.window">
                  <n-input placeholder="e.g. 1m" v-model:value="model.restartPolicy.window" />
                </n-form-item-gi>
                <n-form-item-gi :label="t('fields.condition')" path="restartPolicy.failure">
                  <n-radio-group v-model:value="model.restartPolicy.condition">
                    <n-radio key="any" value="any">Any</n-radio>
                    <n-radio key="on-failure" value="on-failure">On failure</n-radio>
                    <n-radio key="none" value="none">None</n-radio>
                  </n-radio-group>
                </n-form-item-gi>
              </n-grid>
            </x-group>
          </n-tab-pane>
          <n-tab-pane name="log" :tab="t('fields.log_driver')" display-directive="show">
            <n-form-item :label="t('fields.name')" path="logDriver.name">
              <n-radio-group v-model:value="model.logDriver.name">
                <n-radio :key="d.value" :value="d.value" v-for="d in logDrivers">{{ d.label }}</n-radio>
              </n-radio-group>
            </n-form-item>
            <n-form-item :label="t('fields.options')" path="logDriver.options">
              <n-dynamic-input
                v-model:value="model.logDriver.options"
                #="{ index, value }"
                :on-create="newLabel"
              >
                <n-input :placeholder="t('fields.name')" v-model:value="value.name" />
                <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
                <n-input :placeholder="t('fields.value')" v-model:value="value.value" />
              </n-dynamic-input>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="clabels" :tab="t('fields.container_labels')" display-directive="show">
            <n-form-item :show-label="false" path="containerLabels">
              <n-dynamic-input
                v-model:value="model.containerLabels"
                #="{ index, value }"
                :on-create="newLabel"
              >
                <n-input :placeholder="t('fields.name')" v-model:value="value.name" />
                <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
                <n-input :placeholder="t('fields.value')" v-model:value="value.value" />
              </n-dynamic-input>
            </n-form-item>
          </n-tab-pane>
          <n-tab-pane name="dns" tab="DNS" display-directive="show">
            <n-grid cols="1">
              <n-form-item-gi :label="t('fields.hosts')" path="hosts">
                <n-dynamic-input
                  v-model:value="model.hosts"
                  placeholder="IP_address canonical_hostname [aliases...]"
                />
              </n-form-item-gi>
              <n-form-item-gi
                label-placement="top"
                :label="t('fields.dns_servers')"
                path="dns.servers"
              >
                <n-dynamic-input v-model:value="model.dns.servers" placeholder="e.g. 10.10.10.200" />
              </n-form-item-gi>
              <n-form-item-gi
                label-placement="top"
                :label="t('fields.dns_search')"
                path="dns.search"
              >
                <n-dynamic-input v-model:value="model.dns.search" placeholder="e.g. abc.com" />
              </n-form-item-gi>
              <n-form-item-gi label-placement="top" :label="t('fields.options')" path="dns.options">
                <n-dynamic-input v-model:value="model.dns.options" />
              </n-form-item-gi>
            </n-grid>
          </n-tab-pane>
          <n-tab-pane name="others" :tab="t('fields.others')" display-directive="show">
            <n-form-item :label="t('fields.hostname')" path="hostname">
              <n-input :placeholder="t('fields.hostname')" v-model:value="model.hostname" />
            </n-form-item>
          </n-tab-pane>
        </n-tabs>
      </x-panel>
      <n-button
        :disabled="submiting"
        :loading="submiting"
        @click.prevent="submit"
        type="primary"
        :style="{ marginTop: showMore ? '0' : '12px' }"
      >
        <template #icon>
          <n-icon>
            <save-icon />
          </n-icon>
        </template>
        {{ t('buttons.save') }}
      </n-button>
    </n-form>
    <n-alert type="warning" v-if="stack">{{ t('texts.service_notice', { stack }) }}</n-alert>
  </n-space>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import {
  NButton,
  NSpace,
  NInput,
  NInputNumber,
  NDynamicInput,
  NIcon,
  NForm,
  NGrid,
  NFormItem,
  NFormItemGi,
  NCheckboxGroup,
  NCheckbox,
  NRadioGroup,
  NRadio,
  NInputGroup,
  NSelect,
  NTabs,
  NTabPane,
  NText,
  NAlert,
} from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XGroup from "@/components/Group.vue";
import XPanel from "@/components/Panel.vue";
import { useRoute } from "vue-router";
import { router } from "@/router/router";
import serviceApi from "@/api/service";
import type { Service } from "@/api/service";
import configApi from "@/api/config";
import secretApi from "@/api/secret";
import networkApi from "@/api/network";
import type { Network } from "@/api/network";
import { useForm, requiredRule, customRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const model = ref({
  replicas: 1,
  endpoint: {},
  updatePolicy: {},
  rollbackPolicy: {},
  restartPolicy: {},
  logDriver: {},
  resource: { limit: {}, reserve: {} },
  placement: {},
  dns: {}
} as Service)
const stack = computed(() => model.value.labels?.find(l => l.name === 'com.docker.stack.namespace')?.value)
const showMore = ref(false);
const rules: any = {
  name: requiredRule(),
  mounts: customRule((rule: any, value: any[]) => {
    return value?.every(v => v.source && v.target)
  }, t('tips.mounts_rule')),
  configs: customRule((rule: any, value: any[]) => {
    return value?.every(v => v.key && v.path)
  }, t('tips.files_rule')),
  secrets: customRule((rule: any, value: any[]) => {
    return value?.every(v => v.key && v.path)
  }, t('tips.files_rule')),
};
const mountTypeOptions = [
  { label: 'Bind', value: 'bind' },
  { label: 'Volume', value: 'volume' },
  { label: 'Tmpfs', value: 'tmpfs' },
  { label: 'NamedPipe', value: 'npipe' },
];
const consistencyOptions = [
  { label: 'Default', value: 'default' },
  { label: 'Full', value: 'consistent' },
  { label: 'Cached', value: 'cached' },
  { label: 'Delegated', value: 'delegated' },
];
const logDrivers = [
  { label: 'default', value: '' },
  { label: 'json-file', value: 'json-file' },
  { label: 'syslog', value: 'syslog' },
  { label: 'journald', value: 'journald' },
  { label: 'gelf', value: 'gelf' },
  { label: 'fluentd', value: 'fluentd' },
  { label: 'awslogs', value: 'awslogs' },
  { label: 'splunk', value: 'splunk' },
  { label: 'etwlogs', value: 'etwlogs' },
  { label: 'gcplogs', value: 'gcplogs' },
  { label: 'none', value: 'none' },
];
const networks = ref([] as Network[])
const configFiles = ref();
const secretFiles = ref();
const form = ref();
const { submit, submiting } = useForm(form, () => serviceApi.save(model.value), () => {
  window.message.info(t('texts.action_success'));
  router.push({ name: 'service_list' })
})

function newLabel() {
  return {
    name: '',
    value: ''
  }
}

function newPort() {
  return {
    name: '',
    mode: 'ingress',
    protocol: 'tcp',
    targetPort: 0,
    publishedPort: 0,
  }
}

function newMount() {
  return {
    type: 'bind',
    source: '',
    target: '',
    readonly: false,
    consistency: undefined,
  }
}

function newConstraint() {
  return {
    name: '',
    value: '',
    op: '==',
  }
}

function newFile() {
  return {
    key: undefined,
    path: '',
    uid: '0',
    gid: '0',
    mode: 444,
  }
}

async function fetchData() {
  const name = route.params.name as string || ''
  if (name) {
    let r = await serviceApi.find(name);
    model.value = r.data?.service as Service;
  }

  let nr = await networkApi.search();
  networks.value = nr.data as Network[];

  let cr = await configApi.search({ pageIndex: 1, pageSize: 1000 });
  configFiles.value = cr.data?.items.map(c => {
    return { label: c.name, value: `${c.id}:${c.name}` }
  })

  let sr = await secretApi.search({ pageIndex: 1, pageSize: 1000 });
  secretFiles.value = sr.data?.items.map(c => {
    return { label: c.name, value: `${c.id}:${c.name}` }
  })
}

onMounted(fetchData);
</script>
