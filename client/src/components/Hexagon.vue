<script>
import Clock from "./Clock.vue";
import Dock from "./Dock.vue";
import App from "./App.vue";

export default {
  name: "Hexagon",
  components: {Clock, Dock, App},
  data() {
    return {}
  },
  watch: {},
  created() {

  },
  methods: {
    genColors: function () {
      return `rgba(${64 + Math.random() * 128}, ${64 + Math.random() * 128}, ${64 + Math.random() * 128}, 0.8)`
    }
  },
}


</script>


<template>
  <div class="hex-container element">
    <div class="hex-grid">
      <div v-for="(h, a) in 'abcdefgh'" class="hex-frame">
        <div :style="`background-color:${genColors()};`" class="hex-body flex-column">
          <div :key="a" class="label-o3 label-xxs label-w500">{{ h }}</div>
        </div>
      </div>

    </div>
  </div>
</template>

<style lang="scss" scoped>


</style>
