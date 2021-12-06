<template>
    <div class="panel">
        <n-hr v-if="divider === 'top' || divider === 'all'" style="margin-top: 0" />
        <n-element class="title" v-if="title" style>
            <div>
                <n-h4 prefix="bar">{{ title }}</n-h4>
                <n-text depth="3" v-if="subtitle">{{ subtitle }}</n-text>
            </div>
            <n-space :size="6" justify="end" v-if="$slots.action">
                <slot name="action" />
            </n-space>
        </n-element>
        <n-hr v-if="divider === 'middle' || divider === 'all'" style="margin-top: 0" />
        <n-collapse-transition :show="!collapsed">
            <div class="body">
                <slot />
            </div>
        </n-collapse-transition>
        <n-hr v-if="divider === 'bottom' || divider === 'all'" style="margin-bottom: 0" />
    </div>
</template>
  
<script setup lang="ts">
import {
    NElement,
    NSpace,
    NH4,
    NHr,
    NText,
    NCollapseTransition,
} from "naive-ui";

const props = defineProps({
    title: {
        type: String,
    },
    subtitle: {
        type: String,
    },
    divider: {
        type: String,
    },
    collapsed: {
        type: Boolean,
        default: false,
    }
})
</script>

<style lang="scss" scoped>
// .panel {
//     padding: 6px;
//     box-shadow: 0 0px 3px 1px rgb(10 10 10 / 10%);
// }
.title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    // padding-bottom: 6px;
    .n-h4 {
        display: inline;
        margin-right: 8px;
    }
}
.body {
    margin-top: 10px;
}
</style>
  