<template>
    <div :class="['pair-tag', type]">
        <div class="label">{{ label }}</div>
        <div class="value" v-if="Array.isArray(value)" v-for="v in value">{{ v }}</div>
        <div class="value" v-else>{{ value }}</div>
    </div>
</template>

<script setup lang="ts">
import { useThemeVars } from "naive-ui";
import type { PropType } from "vue";

const props = defineProps({
    label: {
        type: String as PropType<unknown | string>,
    },
    value: {
        type: [String, Array] as PropType<string | string[]>,
    },
    type: {
        type: String as PropType<'primary' | 'info' | 'success' | 'warning' | 'error'>,
        default: 'primary',
    },
})
const themeVars: any = useThemeVars().value
const color = themeVars[props.type + 'Color']
const borderColor = convertColor(color, 0.3)
const labelColor = convertColor(color, 0.2)
const valueColor = themeVars.tagColor

function convertColor(hex: string, opacity: number = 1): string {
    let h = hex.replace('#', '');
    if (h.length === 3) {
        h = `${h[0]}${h[0]}${h[1]}${h[1]}${h[2]}${h[2]}`;
    }

    const r = parseInt(h.substring(0, 2), 16);
    const g = parseInt(h.substring(2, 4), 16);
    const b = parseInt(h.substring(4, 6), 16);

    /* Backward compatibility for whole number based opacity values. */
    if (opacity > 1 && opacity <= 100) {
        opacity = opacity / 100;
    }

    return `rgba(${r},${g},${b},${opacity})`;
};
</script>

<style lang="scss" scoped>
.pair-tag {
    font-size: 12px;
    line-height: 20px;
    display: inline-block;
    .label,
    .value {
        display: inherit;
        white-space: nowrap;
        padding: 0 6px;
        border: 1px solid v-bind(borderColor);
    }
    .label {
        padding-left: 8px;
        border-radius: 10px 0 0 10px;
        // font-weight: 500;
        // color: #fff;
        // background-color: v-bind(color);
        color: v-bind(color);
        background-color: v-bind(labelColor);
    }
    .value {
        border-left: none;
        // color: v-bind(color);
        // background-color: v-bind(backgroundColor);
        background-color: v-bind(valueColor);
    }
    .value:last-child {
        padding-right: 8px;
        border-radius: 0 10px 10px 0;
    }
}
</style>