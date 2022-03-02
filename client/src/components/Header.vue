<script>

export default {
  name: "Header",
  data() {
    return {}
  },
  props: {
    target: Object
  }
}

</script>

<template>
  <div class="label-xxl label-w600 label-o5 mb-2 mt-1">
    <span class="label-xl label-w600 label-o4">{{ target.icon }}</span>
    <span class="label-xl label-w600 label-o6 px-2">{{ target.name }}</span>
  </div>
</template>

<style lang="scss">


</style>
