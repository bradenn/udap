<script>

import Range from "./Range.vue";
import Toggle from "./Toggle.vue";
import Text from "./Text.vue";

export default {
  name: "Attribute",
  components: {Toggle, Range, Text},
  data() {
    return {}
  },
  props: {
    attribute: Object,
    entity: String,
    small: Boolean,
  },
  methods: {
    commitChange(attribute) {
      this.$root.requestId("attribute", "request", attribute, this.entity)
    }
  }
}

</script>

<template>
  <div>
    <Range v-if="attribute.type === 'range'" :attribute="attribute" :commit="commitChange"
           :small="small"></Range>
    <Toggle v-else-if="attribute.type === 'toggle'" :attribute="attribute" :commit="commitChange"
            :small="small"></Toggle>
  </div>
</template>

<style scoped>

</style>
