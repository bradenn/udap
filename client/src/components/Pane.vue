<script>

export default {
  name: "Pane",
  data() {
    return {}
  },
  props: {}
}

</script>

<template>
  <div class="element surface">
    <div class="pane mh-100">
      <slot></slot>
    </div>
  </div>
</template>

<style lang="scss" scoped>

.element {
  max-height: 100%;
}

.pane {
  max-height: 80vh !important;
  overflow-scrolling: touch !important;
  overflow-y: scroll !important;
  padding-right: 0.5rem;
}


</style>
