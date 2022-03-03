<script>

export default {
  name: "Toggle",
  data() {
    return {
      waiting: false,
      loading: false,
      local: null,
      cache: {}
    }
  },
  props: {
    attribute: String,
  },
  created() {
    this.cache = this.attribute
  },
  computed: {
    'loading': function () {
      return this.cache === {}
    }
  },
  watch: {
    'attribute': function (d) {
      this.waiting = false
      this.cache = d
    }
  },
  methods: {}
}

</script>

<template>
  <div v-if="loading">Loading...</div>
  <div v-else class="element">
    <div class="h-bar justify-content-start align-items-center align-content-center pb-1">
      <div class="label-xxs label-o2 label-w600">{{ cache.icon }}</div>
      <div class="label-xxs label-o4 label-w500">&nbsp;&nbsp;{{ cache.key }}</div>
      <div class="fill"></div>
      <div v-if="attribute.type === 'int'" class="h-bar gap label-xxs label-o3 label-w400 px-2">
        <div class="label-xxs label-o3">{{ }}</div>
      </div>
    </div>
    <div class="d-flex gap">
      <input v-if="attribute.type==='range'" :key="attribute.key" v-model="local"
             :class="`${attribute.range.cls}`"
             :max="attribute.range.max"
             :min="attribute.range.min" :step="attribute.range.step"
             class="dock dock-small slider"
             type="range"
             v-on:mousedown="changeAttribute(attribute.mode, attribute.key, local)">
    </div>
  </div>
</template>

<style scoped>

</style>
