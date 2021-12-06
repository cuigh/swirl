<template>
  <div :style="{ border: `1px solid ${themeVars.borderColor}`, width: '100%' }">
    <textarea ref="editorRef" />
  </div>
</template>

<script lang="ts">
import { defineComponent, onBeforeUnmount, onMounted, ref, toRefs, watch } from "vue";
import { useThemeVars } from "naive-ui";
import { useStore } from "vuex";
// CodeMirror: common
import CodeMirror from "codemirror";
import "codemirror/mode/yaml/yaml.js";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/seti.css";
// CodeMirror: fold
import "codemirror/addon/fold/foldgutter.css";
import "codemirror/addon/fold/foldcode.js";
import "codemirror/addon/fold/brace-fold.js";
import "codemirror/addon/fold/comment-fold.js";
import "codemirror/addon/fold/indent-fold.js";
import "codemirror/addon/fold/foldgutter.js";
// CodeMirror: search
import "codemirror/addon/scroll/annotatescrollbar.js";
import "codemirror/addon/search/matchesonscrollbar.js";
import "codemirror/addon/search/match-highlighter.js";
import "codemirror/addon/search/jump-to-line.js";
import "codemirror/addon/dialog/dialog.js";
import "codemirror/addon/dialog/dialog.css";
import "codemirror/addon/search/searchcursor.js";
import "codemirror/addon/search/search.js";

export default defineComponent({
  props: {
    modelValue: String,
    defaultValue: {
      type: String,
      default: '',
    },
    readonly: {
      type: Boolean,
      default: false
    }
  },
  setup(props, context) {
    const themeVars = useThemeVars()
    const store = useStore();
    const { modelValue, defaultValue, readonly } = toRefs(props);
    const editorRef = ref();
    let editor: CodeMirror.EditorFromTextArea | null;

    watch(modelValue, () => {
      if (null != editor && modelValue.value && modelValue.value !== editor.getValue()) {
        editor.setValue(modelValue.value);
      }
    });
    watch(readonly, () => {
      if (null != editor) {
        editor.setOption("readOnly", readonly.value);
      }
    });
    onMounted(() => {
      editor = CodeMirror.fromTextArea(editorRef.value, {
        value: modelValue.value,
        indentWithTabs: false,
        smartIndent: true,
        lineNumbers: true,
        readOnly: readonly.value,
        foldGutter: true,
        lineWrapping: true,
        gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter", "CodeMirror-lint-markers"],
        theme: store.state.preference.theme === 'dark' ? 'seti' : 'default',
      });
      editor.on("change", () => {
        context.emit("update:modelValue", editor?.getValue());
      });
      if (defaultValue.value) {
        editor.setValue(defaultValue.value);
      }
    });
    onBeforeUnmount(() => {
      if (null !== editor) {
        editor.toTextArea();
        editor = null;
      }
    });
    return { themeVars, editorRef };
  }
});
</script>

