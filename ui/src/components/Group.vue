<template>
  <fieldset :class="['group', store.state.preference.theme === 'dark' ? 'dark' : 'light']">
    <legend class="title" v-if="title">{{ title }}</legend>
    <slot />
  </fieldset>
</template>

<script setup lang="ts">
import type { PropType } from "vue";
import { useStore } from "vuex";

const props = defineProps({
  title: {
    type: String as PropType<unknown | string>,
  }
})
const store = useStore()
</script>

<style lang="scss" scoped>
.group {
  // box-shadow: 0 0 2px 0 rgb(10 10 10 / 10%);
  border-radius: 5px;
  border-width: 1px;
  border-style: solid;
  margin-bottom: 12px;
  padding-bottom: 0;
  // border-left: none;
  // border-right: none;
  // border-bottom: none;
  // padding: 6px 0 0 0;
  .title {
    font-weight: 500;
    // margin: 0 auto;
    // padding: 0 6px;
  }
}
.group.light {
  border-color: rgba(224, 224, 230, 0.5);
  .title {
    color: rgb(158, 164, 170);
  }
}
.group.dark {
  border-color: rgba(255, 255, 255, 0.24);
  .title {
    color: rgba(255, 255, 255, 0.52);
  }
}
</style>