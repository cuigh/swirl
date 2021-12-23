<template>
  <n-card
    size="small"
    :segmented="{ content: true, action: 'soft' }"
    :style="{ height: info.height + 'px' }"
  >
    <div ref="container" style="width: 100%; height: 100%"></div>
    <template #header>
      <n-icon class="drag-handle" style="cursor: move; vertical-align: middle" depth="3">
        <menu-outline />
      </n-icon>
      {{ info.title }}
    </template>
    <template #header-extra>
      <n-button quaternary circle type="error" size="tiny" #icon @click="remove">
        <n-icon>
          <close-outline />
        </n-icon>
      </n-button>
    </template>
  </n-card>
</template>

<script setup lang="ts">
import {
  NCard,
  NIcon,
  NButton,
} from "naive-ui";
import { ref, onMounted } from "vue";
import { CloseOutline, MenuOutline } from "@vicons/ionicons5";
import { useResizeObserver } from '@vueuse/core'
import { ChartInfo } from "@/api/dashboard";
import { createChart } from "./chart";
import type { Chart } from "./chart";

interface Props {
  info: ChartInfo;
  data?: any;
}

// const props = withDefaults(defineProps<Props>(), {
//   data: () => null,
// })
const props = defineProps<Props>()
const emits = defineEmits(['remove'])

const container = ref()
var chart: Chart
var resizeTimer: NodeJS.Timeout

useResizeObserver(container, () => {
  clearTimeout(resizeTimer)
  resizeTimer = setTimeout(() => chart.resize(), 500)
})

function initChart() {
  chart = createChart(container.value as HTMLElement, props.info)
  props.data && chart.setData(props.data)
  setTimeout(() => chart.resize(), 100)
}

function setData(d: any) {
  chart.setData(d)
}

function remove() {
  emits('remove', props.info.id)
}

defineExpose({ id: props.info.id, setData })

onMounted(initChart)
</script>

<style lang="scss" scoped>
.n-card {
  box-shadow: 0 0 5px 0 rgb(10 10 10 / 10%);
  ::v-deep(.n-card-header) {
    padding: 6px 8px;
  }
  ::v-deep(.n-card__content) {
    padding: 6px 8px;
  }
}
// .chart.light {
//     border: 1px solid rgb(239, 239, 245);
//     .title {
//         color: rgb(158, 164, 170);
//     }
// }
// .chart.dark {
//     border: 1px solid rgba(255, 255, 255, 0.09);
//     .title {
//         color: rgba(255, 255, 255, 0.52);
//     }
// }
</style>