<template>
  <n-space vertical :size="12">
    <n-input-group style="width: 100%">
      <n-input-group-label style="min-width: 52px">{{ t('fields.command') }}</n-input-group-label>
      <n-input :placeholder="t('fields.command')" v-model:value="command" :readonly="active" />
      <n-button type="primary" @click="connect" v-if="!active">{{ t('buttons.execute') }}</n-button>
      <n-button type="error" @click="disconnect" v-else>{{ t('buttons.disconnect') }}</n-button>
    </n-input-group>
    <div id="xterm" class="xterm" />
  </n-space>
</template>

<script setup lang="ts">
import { onUnmounted, ref } from "vue";
import {
  NSpace,
  NButton,
  NInput,
  NInputGroup,
  NInputGroupLabel,
} from "naive-ui";
import 'xterm/css/xterm.css'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { AttachAddon } from 'xterm-addon-attach'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const props = defineProps({
  node: {
    type: String,
    required: true,
  },
  id: {
    type: String,
    required: true,
  },
})
const command = ref('/bin/sh');
const active = ref(false);
var socket: null | WebSocket;
var term: null | Terminal;

function connect() {
  if (!command.value) {
    window.message.error(t('tips.command_empty'))
    return
  }

  active.value = true
  let protocol = (location.protocol === "https:") ? "wss://" : "ws://";
  let host = import.meta.env.DEV ? 'localhost:8002' : location.host;
  let cmd = encodeURIComponent(command.value)
  socket = new WebSocket(`${protocol}${host}/api/container/connect?node=${props.node}&id=${props.id}&cmd=${cmd}`);
  socket.onopen = () => {
    const fit = new FitAddon();
    term = new Terminal({ fontSize: 14, cursorBlink: true });
    term.loadAddon(new AttachAddon(socket as WebSocket));
    term.loadAddon(fit);
    term.open(document.getElementById('xterm') as HTMLElement);
    fit.fit();
    term.focus();
  };
  // socket.onclose = () => {
  //   console.log('close socket')
  // };
  socket.onerror = (e) => {
    console.log('socket error: ' + e)
  }
};

function disconnect() {
  if (socket) {
    socket.close()
    socket = null
  }
  if (term) {
    term.dispose()
    term = null
  }
  active.value = false
}

onUnmounted(disconnect);
</script>
