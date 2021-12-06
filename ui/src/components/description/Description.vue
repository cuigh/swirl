<script lang="ts">
import {
  h,
  defineComponent,
  PropType,
  VNode,
  Component,
} from 'vue'
import { NGi, NGrid, useThemeVars } from 'naive-ui'
import { ThemeCommonVars } from 'naive-ui/lib/_styles/common'

const descriptionProps = {
  bordered: {
    type: Boolean,
    default: false,
  },
  cols: {
    type: [Number, String] as PropType<number | string>,
    default: 2,
  },
  labelWidth: {
    type: Number,
    default: 80,
  },
  labelAlign: {
    type: String as PropType<'left' | 'right' | 'center'>,
    default: "right",
  },
  labelPosition: {
    type: String as PropType<'top' | 'left'>,
    default: "left",
  }
} as const

function renderItem(n: VNode, props: any, themeVars: ThemeCommonVars) {
  if (n.type && (n.type as Component).name !== "DescriptionItem") {
    return null
  }

  const itemProps: any = {
    class: `desc-item ${props.labelPosition}`,
    span: n.props?.span || 1,
  }
  if (props.bordered) {
    itemProps.style = {
      border: `1px solid ${themeVars.borderColor}`
    }
  }

  return h(
    NGi,
    itemProps,
    {
      default: () => [
        h("div", {
          class: 'desc-label',
          style: {
            minWidth: `${props.labelWidth}px`,
            textAlign: props.labelAlign,
          },
        }, { default: () => n.props?.label }),
        h("div", {
          class: "desc-value",
        }, { default: () => n.children }),
      ]
    }
  )
}

export default defineComponent({
  name: 'Description',
  props: descriptionProps,
  setup() {
    return {
      themeVars: useThemeVars(),
    }
  },
  render() {
    return h(
      NGrid,
      {
        xGap: 4,
        yGap: 4,
        cols: this.$props.cols,
      },
      {
        default: () => {
          if (!this.$slots.default) {
            return null;
          }
          const children = this.$slots.default()
          return children.map(n => renderItem(n, this.$props, this.themeVars))
        }
      }
    )
  }
})
</script>
  
<style scoped>
.desc-item.left {
  display: flex;
  line-height: 30px;
}
.desc-label {
  font-weight: 400;
}
.desc-item.left > .desc-label {
  display: inline-block;
}
.desc-item.left > .desc-label:after {
  content: ":";
  margin-right: 8px;
}
.desc-value {
  display: inline-block;
  word-break: break-word;
  overflow-wrap: break-word;
  width: 100%;
}
</style>