<template>
  <x-page-header :subtitle="service.name">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'service_list' })">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
      <n-button
        secondary
        size="small"
        @click="scale"
        v-if="service.mode === 'replicated' || service.mode === 'replicated-job'"
      >{{ t('buttons.scale') }}</n-button>
      <n-button secondary size="small" @click="restart">{{ t('buttons.restart') }}</n-button>
      <n-button secondary size="small" @click="rollback" type="warning">{{ t('buttons.rollback') }}</n-button>
      <n-button secondary size="small" @click="remove" type="error">{{ t('buttons.delete') }}</n-button>
      <n-button
        secondary
        size="small"
        @click="$router.push({ name: 'service_edit', params: { name: service.name } })"
      >{{ t('buttons.edit') }}</n-button>
    </template>
  </x-page-header>
  <div class="page-body">
    <n-tabs type="line" style="margin-top: -12px">
      <n-tab-pane name="detail" :tab="t('fields.detail')">
        <n-space vertical :size="16">
          <x-description label-placement="left" label-align="right" :label-width="100">
            <x-description-item :label="t('fields.id')">{{ service.id }}</x-description-item>
            <x-description-item :label="t('fields.name')">{{ service.name }}</x-description-item>
            <x-description-item :label="t('objects.image')" :span="2">{{ service.image }}</x-description-item>
            <x-description-item :label="t('fields.mode')">
              <n-tag
                round
                size="small"
                :type="['replicated', 'replicated-job'].includes(service.mode) ? 'info' : 'error'"
              >{{ service.mode }}</n-tag>
            </x-description-item>
            <x-description-item :label="t('objects.task', 2)">
              <n-tag
                round
                size="small"
                :type="service.runningTasks === service.desiredTasks ? 'success' : 'error'"
              >{{ service.runningTasks + '/' + service.desiredTasks }}</n-tag>
            </x-description-item>
            <x-description-item
              :label="t('fields.command')"
              :span="2"
              v-if="service.command"
            >{{ service.command }}</x-description-item>
            <x-description-item
              :label="t('fields.args')"
              :span="2"
              v-if="service.args"
            >{{ service.args }}</x-description-item>
            <x-description-item :label="t('fields.work_dir')" v-if="service.dir">{{ service.dir }}</x-description-item>
            <x-description-item :label="t('fields.user')" v-if="service.user">{{ service.user }}</x-description-item>
            <x-description-item :label="t('fields.created_at')">{{ service.createdAt }}</x-description-item>
            <x-description-item :label="t('fields.updated_at')">{{ service.updatedAt }}</x-description-item>
            <x-description-item :label="t('fields.update_state')" v-if="service.update?.state">
              <n-tag
                round
                size="small"
                :type="service.update.state === 'completed' ? 'success' : 'error'"
              >{{ service.update.state }}</n-tag>
            </x-description-item>
          </x-description>
          <x-panel :title="t('fields.cli')" :collapsed="!showCli">
            <template #action>
              <n-button
                secondary
                size="small"
                @click="showCli = !showCli"
              >{{ showCli ? t('buttons.collapse') : t('buttons.expand') }}</n-button>
            </template>
            <x-code :code="cli" language="bash" />
          </x-panel>
          <x-panel :title="t('fields.env')" v-if="!isEmpty(service.env)">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>{{ t('fields.value') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="e in service.env">
                  <td>{{ e.name }}</td>
                  <td>{{ e.value }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
          <x-panel :title="t('fields.labels')" v-if="!isEmpty(service.labels)">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>{{ t('fields.value') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="label in service.labels">
                  <td>{{ label.name }}</td>
                  <td>{{ label.value }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
          <x-panel :title="t('objects.network', 2)">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>{{ t('fields.address') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="label in service.endpoint?.vips">
                  <td>
                    <x-anchor
                      :url="{ name: 'network_detail', params: { name: label.name || label.id } }"
                    >{{ label.name || label.id }}</x-anchor>
                  </td>
                  <td>{{ label.ip }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
          <x-panel :title="t('objects.config', 2)" v-if="!isEmpty(service.configs)">
            <n-table size="small" :bordered="true" :single-line="false" style="width: 100%">
              <thead>
                <tr>
                  <th>{{ t('fields.filename') }}</th>
                  <th>{{ t('fields.target_path') }}</th>
                  <th>UID</th>
                  <th>GID</th>
                  <th>{{ t('fields.mode') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="f in service.configs">
                  <td>{{ f.key.substring(f.key.indexOf(':') + 1) }}</td>
                  <td>{{ f.path }}</td>
                  <td>{{ f.uid }}</td>
                  <td>{{ f.gid }}</td>
                  <td>{{ f.mode }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
          <x-panel :title="t('objects.secret', 2)" v-if="!isEmpty(service.secrets)">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.filename') }}</th>
                  <th>{{ t('fields.target_path') }}</th>
                  <th>UID</th>
                  <th>GID</th>
                  <th>{{ t('fields.mode') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="f in service.secrets">
                  <td>{{ f.key.substring(f.key.indexOf(':') + 1) }}</td>
                  <td>{{ f.path }}</td>
                  <td>{{ f.uid }}</td>
                  <td>{{ f.gid }}</td>
                  <td>{{ f.mode }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
          <x-panel :title="t('fields.mounts')" v-if="!isEmpty(service.mounts)">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.type') }}</th>
                  <th>{{ t('fields.source_path') }}</th>
                  <th>{{ t('fields.target_path') }}</th>
                  <th>{{ t('fields.readonly') }}</th>
                  <th>{{ t('fields.consistency') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="m in service.mounts">
                  <td>{{ m.type }}</td>
                  <td>{{ m.source }}</td>
                  <td>{{ m.target }}</td>
                  <td>
                    <n-tag round size="small">{{ m.readonly ? '是' : '否' }}</n-tag>
                  </td>
                  <td>{{ m.consistency }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
          <x-panel :title="t('fields.resources')">
            <x-description label-align="left" :label-width="40">
              <x-description-item
                :label="t('fields.limit')"
                v-if="service.resource.limit.cpu && service.resource.limit.memory"
              >
                <n-space :size="4">
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.cpu')"
                    :value="service.resource.limit.cpu.toString()"
                  />
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.memory')"
                    :value="service.resource.limit.memory"
                  />
                </n-space>
              </x-description-item>
              <x-description-item
                :label="t('fields.reserve')"
                v-if="service.resource.reserve.cpu && service.resource.reserve.memory"
              >
                <n-space :size="4">
                  <x-pair-tag
                    type="success"
                    :label="t('fields.cpu')"
                    :value="service.resource.reserve.cpu.toString()"
                  />
                  <x-pair-tag
                    type="success"
                    :label="t('fields.memory')"
                    :value="service.resource.reserve.memory"
                  />
                </n-space>
              </x-description-item>
            </x-description>
          </x-panel>
          <x-panel :title="t('fields.placement')">
            <x-description label-align="left" cols="1" :label-width="40">
              <x-description-item
                :label="t('fields.constraints')"
                v-if="service.placement.constraints && service.placement.constraints.length"
              >
                <n-space :size="4">
                  <n-tag
                    round
                    size="small"
                    v-for="c in service.placement.constraints"
                  >{{ c.name + ' ' + c.op + ' ' + c.value }}</n-tag>
                </n-space>
              </x-description-item>
              <x-description-item
                :label="t('fields.preferences')"
                v-if="service.placement.preferences && service.placement.preferences.length"
              >
                <n-space :size="4">
                  <n-tag round size="small" v-for="p in service.placement.preferences">{{ p }}</n-tag>
                </n-space>
              </x-description-item>
            </x-description>
          </x-panel>
          <x-panel :title="t('fields.schedule')">
            <x-description label-align="left" cols="1" :label-width="40">
              <x-description-item :label="t('fields.update')" v-if="hasUpdatePolicy">
                <n-space :size="4">
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.parallelism')"
                    :value="service.updatePolicy.parallelism.toString()"
                    v-if="service.updatePolicy.parallelism"
                  />
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.delay')"
                    :value="service.updatePolicy.delay"
                    v-if="service.updatePolicy.delay"
                  />
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.failure_action')"
                    :value="service.updatePolicy.failureAction"
                    v-if="service.updatePolicy.failureAction"
                  />
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.order')"
                    :value="service.updatePolicy.order"
                    v-if="service.updatePolicy.order"
                  />
                </n-space>
              </x-description-item>
              <x-description-item :label="t('fields.rollback')" v-if="hasRollbackPolicy">
                <n-space :size="4">
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.parallelism')"
                    :value="service.rollbackPolicy.parallelism.toString()"
                    v-if="service.rollbackPolicy.parallelism"
                  />
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.delay')"
                    :value="service.rollbackPolicy.delay"
                    v-if="service.rollbackPolicy.delay"
                  />
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.failure_action')"
                    :value="service.rollbackPolicy.failureAction"
                    v-if="service.rollbackPolicy.failureAction"
                  />
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.order')"
                    :value="service.rollbackPolicy.order"
                    v-if="service.rollbackPolicy.order"
                  />
                </n-space>
              </x-description-item>
              <x-description-item :label="t('fields.restart')" v-if="hasRestartPolicy">
                <n-space :size="4">
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.condition')"
                    :value="service.restartPolicy.condition"
                    v-if="service.restartPolicy.condition"
                  />
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.max_attempts')"
                    :value="service.restartPolicy.maxAttempts.toString()"
                    v-if="service.restartPolicy.maxAttempts"
                  />
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.delay')"
                    :value="service.restartPolicy.delay"
                    v-if="service.restartPolicy.delay"
                  />
                  <x-pair-tag
                    type="warning"
                    :label="t('fields.window')"
                    :value="service.restartPolicy.window"
                    v-if="service.restartPolicy.window"
                  />
                </n-space>
              </x-description-item>
            </x-description>
          </x-panel>
          <x-panel :title="t('fields.log_driver')" v-if="service.logDriver.name">
            <x-description label-align="left" cols="1" :label-width="40">
              <x-description-item :label="t('fields.name')">{{ service.logDriver.name }}</x-description-item>
            </x-description>
            <n-table
              size="small"
              :bordered="true"
              :single-line="false"
              v-if="!isEmpty(service.logDriver.options)"
            >
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>{{ t('fields.value') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="label in service.logDriver.options">
                  <td>{{ label.name }}</td>
                  <td>{{ label.value }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
          <x-panel :title="t('fields.container_labels')" v-if="!isEmpty(service.containerLabels)">
            <n-table size="small" :bordered="true" :single-line="false">
              <thead>
                <tr>
                  <th>{{ t('fields.name') }}</th>
                  <th>{{ t('fields.value') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="label in service.containerLabels">
                  <td>{{ label.name }}</td>
                  <td>{{ label.value }}</td>
                </tr>
              </tbody>
            </n-table>
          </x-panel>
          <x-panel title="DNS" v-if="hasDns">
            <x-description label-align="left" cols="1" :label-width="40">
              <x-description-item
                :label="t('fields.hosts')"
                v-if="service.hosts && service.hosts.length"
              >
                <n-space :size="4">
                  <n-tag round size="small" v-for="h in service.hosts">{{ h }}</n-tag>
                </n-space>
              </x-description-item>
              <x-description-item
                :label="t('fields.dns_servers')"
                v-if="service.dns.servers && service.dns.servers.length"
              >
                <n-space :size="4">
                  <n-tag round size="small" v-for="s in service.dns.servers">{{ s }}</n-tag>
                </n-space>
              </x-description-item>
              <x-description-item
                :label="t('fields.dns_search')"
                v-if="service.dns.search && service.dns.search.length"
              >
                <n-space :size="4">
                  <n-tag round size="small" v-for="s in service.dns.search">{{ s }}</n-tag>
                </n-space>
              </x-description-item>
              <x-description-item
                :label="t('fields.options')"
                v-if="service.dns.options && service.dns.options.length"
              >
                <n-space :size="4">
                  <n-tag round size="small" v-for="s in service.dns.options">{{ s }}</n-tag>
                </n-space>
              </x-description-item>
            </x-description>
          </x-panel>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="raw" :tab="t('fields.raw')">
        <x-code :code="raw" language="json" />
      </n-tab-pane>
      <n-tab-pane name="task" :tab="t('objects.task', 2)">
        <n-table size="small" :bordered="true" :single-line="false">
          <thead>
            <tr>
              <th>{{ t('fields.id') }}</th>
              <th>{{ t('fields.state') }}</th>
              <th>{{ t('objects.node') }}</th>
              <th>{{ t('objects.network', 2) }}</th>
              <th>{{ t('fields.updated_at') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in tasks">
              <td>
                <x-anchor :url="{ name: 'task_detail', params: { id: t.id } }">{{ t.id }}</x-anchor>
              </td>
              <td>
                <n-tag
                  round
                  size="small"
                  :type="t.state === 'running' ? 'success' : 'error'"
                >{{ t.state }}</n-tag>
              </td>
              <td>
                <x-anchor :url="{ name: 'node_detail', params: { id: t.nodeId } }">{{ t.nodeName }}</x-anchor>
              </td>
              <td>
                <n-space :size="4">
                  <n-tag
                    round
                    size="small"
                    v-for="n in t.networks"
                  >{{ isEmpty(n.ips) ? n.name : (n.name + ": " + n.ips?.join(',')) }}</n-tag>
                </n-space>
              </td>
              <td>{{ t.updatedAt }}</td>
            </tr>
          </tbody>
        </n-table>
      </n-tab-pane>
      <n-tab-pane name="logs" :tab="t('fields.logs')" display-directive="show:lazy">
        <x-logs type="service" :id="service.name" v-if="store.getters.allow('service.logs')"></x-logs>
        <n-alert type="info" v-else>{{ t('texts.403') }}</n-alert>
      </n-tab-pane>
      <n-tab-pane name="status" :tab="t('fields.status')">
        <x-dashboard type="service" :name="service.name" />
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, h, computed } from "vue";
import {
  NButton,
  NTag,
  NSpace,
  NIcon,
  NTable,
  NTabs,
  NTabPane,
  NInputNumber,
  NAlert,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import { useStore } from "vuex";
import XPageHeader from "@/components/PageHeader.vue";
import XAnchor from "@/components/Anchor.vue";
import XPairTag from "@/components/PairTag.vue";
import XLogs from "@/components/Logs.vue";
import XDashboard from "@/components/Dashboard.vue";
import XPanel from "@/components/Panel.vue";
import XCode from "@/components/Code.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import serviceApi from "@/api/service";
import type { Service } from "@/api/service";
import taskApi from "@/api/task";
import type { Task } from "@/api/task";
import { useRoute } from "vue-router";
import { router } from "@/router/router";
import { isEmpty } from "@/utils";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const store = useStore();
const service = ref({
  resource: {
    limit: {},
    reserve: {},
  },
  placement: {},
  updatePolicy: {},
  rollbackPolicy: {},
  restartPolicy: {},
  logDriver: {},
  dns: {},
} as Service);
const tasks = ref([] as Task[]);
const raw = ref('')
const showCli = ref(false);
const cli = ref('')
const hasUpdatePolicy = computed(() => {
  const p = service.value.updatePolicy
  return p.parallelism || p.delay || p.failureAction || p.order
})
const hasRollbackPolicy = computed(() => {
  const p = service.value.rollbackPolicy
  return p.parallelism || p.delay || p.failureAction || p.order
})
const hasRestartPolicy = computed(() => {
  const p = service.value.restartPolicy
  return p.condition || p.delay || p.maxAttempts || p.window
})
const hasDns = computed(() => {
  const dns = service.value.dns
  return !isEmpty(service.value.hosts, dns.servers, dns.search, dns.options)
})

function restart() {
  window.dialog.warning({
    title: t('dialogs.restart_service.title'),
    content: t('dialogs.restart_service.body'),
    positiveText: t('buttons.confirm'),
    negativeText: t('buttons.cancel'),
    onPositiveClick: async () => {
      serviceApi.restart(service.value.name);
      window.message.success(t('texts.action_success'))
      fetchData()
    }
  })
}

function rollback() {
  window.dialog.warning({
    title: t('dialogs.rollback_service.title'),
    content: t('dialogs.rollback_service.body'),
    positiveText: t('buttons.confirm'),
    negativeText: t('buttons.cancel'),
    onPositiveClick: async () => {
      serviceApi.rollback(service.value.name);
      window.message.success(t('texts.action_success'))
      fetchData()
    }
  })
}

function scale() {
  const count = ref(service.value.replicas) as any
  window.dialog.warning({
    title: t('dialogs.scale_service.title'),
    content: () => h(NInputNumber, { min: 0, defaultValue: count }),
    positiveText: t('buttons.confirm'),
    negativeText: t('buttons.cancel'),
    onPositiveClick: async () => {
      serviceApi.scale(service.value.name, count.value, service.value.version);
      window.message.success(t('texts.action_success'))
      fetchData()
    }
  })
}

function remove() {
  window.dialog.warning({
    title: t('dialogs.remove_service.title'),
    content: t('dialogs.remove_service.body'),
    positiveText: t('buttons.confirm'),
    negativeText: t('buttons.cancel'),
    onPositiveClick: async () => {
      serviceApi.delete(service.value.name);
      window.message.success(t('texts.action_success'))
      router.push({ name: 'service_list' })
    }
  })
}

function generateCli(s: Service) {
  const arr = new Array<string>()

  arr.push("docker service create")
  arr.push(`--name ${s.name}`)
  s.dir && arr.push(`--workdir ${s.dir}`)
  s.user && arr.push(`--workdir ${s.user}`)
  s.hostname && arr.push(`--hostname ${s.hostname}`)
  s.mode !== "replicated" && arr.push(`--mode ${s.mode}`)
  s.replicas > 1 && arr.push(`--replicas ${s.replicas}`)
  // network
  s.networks.forEach(n => {
    const vip = s.endpoint.vips.find(v => v.id === n)
    vip && arr.push(`--network ${vip.name}`)
  })
  // ports
  if (s.endpoint.ports && s.endpoint.ports.length) {
    s.endpoint.mode !== 'vip' && arr.push(`--endpoint-mode ${s.endpoint.mode}`)
    s.endpoint.ports.forEach(p => {
      arr.push(`--publish mode=${p.mode},target=${p.targetPort},published=${p.publishedPort},protocol=${p.protocol}`)
    })
  }
  // placement
  s.placement.constraints?.forEach(c => arr.push(`--constraint '${c.name} ${c.op} ${c.value}'`))
  s.placement.preferences?.forEach(p => arr.push(`--placement-pref '${p}'`))
  // env
  s.env?.forEach(e => arr.push(`--env '${e.name}=${e.value}'`))
  // labels
  s.labels?.forEach(l => arr.push(`--label '${l.name}=${l.value}'`))
  // container labels
  s.containerLabels?.forEach(l => arr.push(`--container-label '${l.name}=${l.value}'`))
  // mounts
  s.mounts?.forEach(m => {
    let arg = `--mount type=${m.type},source=${m.source},destination=${m.target}`
    if (m.readonly) {
      arg += ',ro=1'
    }
    arr.push(arg)
  })
  // configs
  s.configs?.forEach(c => {
    const source = c.key.substring(c.key.indexOf(':') + 1)
    // TODO: UID / GID
    arr.push(`--config source=${source},target=${c.path},mode=0${c.mode}`)
  })
  // secrets
  s.secrets?.forEach(c => {
    const source = c.key.substring(c.key.indexOf(':') + 1)
    // TODO: UID / GID
    arr.push(`--secret source=${source},target=${c.path},mode=0${c.mode}`)
  })
  // resource
  s.resource.limit?.cpu > 0 && arr.push(`--limit-cpu ${s.resource.limit.cpu}`)
  s.resource.limit?.memory && arr.push(`--limit-memory ${s.resource.limit.memory.replace(/\s*/g, "")}`)
  s.resource.reserve?.cpu > 0 && arr.push(`--reserve-cpu ${s.resource.reserve.cpu}`)
  s.resource.reserve?.memory && arr.push(`--reserve-memory ${s.resource.reserve.memory.replace(/\s*/g, "")}`)
  // log driver
  if (s.logDriver.name) {
    arr.push(`--log-driver ${s.logDriver.name}`)
    s.logDriver.options?.forEach(opt => arr.push(`--log-opt '${opt.name}=${opt.value}'`))
  }
  // update policy
  s.updatePolicy.parallelism && arr.push(`--update-parallelism ${s.updatePolicy.parallelism}`)
  s.updatePolicy.delay && arr.push(`--update-delay ${s.updatePolicy.delay}`)
  s.updatePolicy.failureAction && arr.push(`--update-failure-action ${s.updatePolicy.failureAction}`)
  s.updatePolicy.order && arr.push(`--update-order ${s.updatePolicy.order}`)
  // rollback policy
  s.rollbackPolicy.parallelism && arr.push(`--rollback-parallelism ${s.rollbackPolicy.parallelism}`)
  s.rollbackPolicy.delay && arr.push(`--rollback-delay ${s.rollbackPolicy.delay}`)
  s.rollbackPolicy.failureAction && arr.push(`--rollback-failure-action ${s.rollbackPolicy.failureAction}`)
  s.rollbackPolicy.order && arr.push(`--rollback-order ${s.rollbackPolicy.order}`)
  // restart policy
  s.restartPolicy.condition && arr.push(`--restart-condition ${s.restartPolicy.condition}`)
  s.restartPolicy.delay && arr.push(`--restart-delay ${s.restartPolicy.delay}`)
  s.restartPolicy.maxAttempts && arr.push(`--restart-max-attempts ${s.restartPolicy.maxAttempts}`)
  s.restartPolicy.window && arr.push(`--restart-window ${s.restartPolicy.window}`)
  // hosts
  s.hosts?.forEach(h => arr.push(`--host ${h}`))
  // dns
  s.dns?.servers?.forEach(s => arr.push(`--dns ${s}`))
  s.dns?.search?.forEach(s => arr.push(`--dns-search ${s}`))
  s.dns?.options?.forEach(s => arr.push(`--dns-option ${s}`))
  // image
  arr.push(s.image)
  // command
  s.command && arr.push(s.command)
  // args
  s.args && arr.push(s.args)
  return arr.join(' \\\n    ')
}

async function fetchData() {
  const name = route.params.name as string
  let results = await Promise.all([
    serviceApi.find(name, true),
    taskApi.search({ service: name, pageIndex: 1, pageSize: 100 }),
  ])

  service.value = results[0].data?.service as Service
  raw.value = results[0].data?.raw as string;
  tasks.value = results[1].data?.items as Task[];
  cli.value = generateCli(service.value)
}

onMounted(fetchData);
</script>