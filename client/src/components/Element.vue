<script>

export default {
  name: "Element",
  data() {
    return {}
  },
  props: {}
}

</script>

<template>
  <div class="element">
    <slot></slot>
  </div>
</template>

<style lang="scss">


.root > .element {
  overflow-y: scroll;
}
</style>
