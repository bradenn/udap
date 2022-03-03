<script>

export default {
  name: "Dock",
  data() {
    return {}
  },
  props: {
    small: Boolean,
    os: Boolean,
  }
}

</script>

<template>
  <div class="element dock" v-bind:class="`${small?'dock-small':''} ${os?'dock-os':''}`">
    <slot></slot>
  </div>
</template>

<style lang="scss">

</style>
