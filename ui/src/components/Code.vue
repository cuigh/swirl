<template>
  <div class="code" :style="{ backgroundColor: themeVars.codeColor }" v-if="code">
    <n-button
      tertiary
      size="tiny"
      class="copy"
      :type="copied ? 'success' : 'default'"
      #icon
      @click="() => copy()"
      v-if="isSupported"
    >
      <n-icon>
        <checkmark-outline v-if="copied" />
        <copy-outline v-else />
      </n-icon>
    </n-button>
    <n-code :code="code" :language="language" />
  </div>
</template>
  
<script setup lang="ts">
import {
  NCode,
  NButton,
  NIcon,
  useThemeVars,
} from "naive-ui";
import { CopyOutline, CheckmarkOutline } from "@vicons/ionicons5";
import { useClipboard } from '@vueuse/core'

const props = defineProps({
  code: {
    type: String,
  },
  language: {
    type: String,
  },
})

const themeVars = useThemeVars()
const { copy, copied, isSupported } = useClipboard({ source: props.code })
</script>

<style lang="scss" scoped>
.code {
  border-radius: 3px;
  position: relative;
  .copy {
    position: absolute;
    top: 0;
    right: 0;
    border-radius: 0 3px 0 3px;
  }
  ::v-deep(pre) {
    padding: 6px;
    overflow-x: scroll;
    font-size: 12px;
  }
}
</style>
  