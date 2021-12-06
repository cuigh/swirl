<template>
    <div :class="['statistic', store.state.preference.theme === 'dark' ? 'dark' : 'light']">
        <n-space :size="6" justify="end" v-if="$slots.icon">
            <n-icon size="48" style="margin-top: 6px" :depth="4">
                <slot name="icon"></slot>
            </n-icon>
        </n-space>
        <div style="text-align: center; margin-bottom: -8px">
            <div class="title">{{ title }}</div>
            <slot>
                <n-text class="value" depth="3">{{ value }}</n-text>
            </slot>
        </div>
    </div>
</template>

<script setup lang="ts">
import {
    NSpace,
    NText,
    NIcon,
} from "naive-ui";
import type { PropType } from "vue";
import { useStore } from "vuex";

const props = defineProps({
    title: {
        type: String as PropType<unknown | string>,
    },
    value: {
        type: String,
    }
})
const store = useStore()
</script>

<style lang="scss" scoped>
.statistic {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px;
    border-radius: 5px;
    box-shadow: 0 0 4px 0 rgb(10 10 10 / 10%);
    .title {
        font-size: 14px;
    }
    ::v-deep(.n-a) {
        font-size: 24px;
    }
}
.statistic.light {
    border: 1px solid rgb(239, 239, 245);
    .title {
        color: rgb(158, 164, 170);
    }
}
.statistic.dark {
    border: 1px solid rgba(255, 255, 255, 0.09);
    .title {
        color: rgba(255, 255, 255, 0.52);
    }
}
</style>