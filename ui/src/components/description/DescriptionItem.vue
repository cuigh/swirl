<script lang="ts">
import { defineComponent } from "vue"

const descriptionItemProps = {
  label: String,
  span: {
    type: Number,
    default: 1
  }
} as const

export default defineComponent({
  name: 'DescriptionItem',
  props: descriptionItemProps,
  render() {
    return null
  }
})
</script>

