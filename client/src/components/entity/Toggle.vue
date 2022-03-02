<script>

import Dock from "../Dock.vue";
import Loading from "../Loading.vue";

export default {
  name: "Toggle",
  components: {Loading, Dock},
  data() {
    return {
      latest: new Date(),
      waiting: false,
      sliding: false,
      loading: false,
      local: {},
      config: {
        on: {
          min: 0,
          max: 100,
          step: 5,
          unit: '%',
          icon: '􀇯',
          class: 'slider-dim',
        },
      }
    }
  },
  created() {
    this.latest = new Date()
    this.local = this.attribute
  },
  beforeMount() {

  },
  props: {
    attribute: Object,
    commit: Function,
    small: Boolean,
  },
  computed: {
    active: function () {
      return this.local.request === "true"
    }
  },
  watch: {
    'attribute': function (d) {
      this.waiting = false
      this.local = d
    }
  },
  methods: {
    commitChanges: function (e) {
      this.local.request = this.local.request === "false" ? "true" : "false"
      this.commit(this.local)
      this.waiting = true
      this.latest = new Date()
      // Prevent updates until after we send the state
    }
  }
}

</script>

<template>
  <div v-if="small">
    <div class="h-bar gap label-sm label-w600 text-uppercase label-o4 px-2">
      <Loading v-if="waiting" class="lh-2"></Loading>
      <div @click="commitChanges">{{ active ? "ON" : "OFF" }}</div>
    </div>
  </div>
  <div v-else v-if="attribute" class="element" v-on:click.stop>

    <div class="h-bar justify-content-start align-items-center align-content-center">
      <div class="label-xxs label-o2 label-w600">{{ config[attribute.key].icon }}</div>
      <div class="label-xxs label-o4 label-w500">&nbsp;&nbsp;{{ attribute.key }}</div>
      <div class="fill"></div>

      <Loading v-if="waiting" class="lh-2"></Loading>
      <div class="h-bar gap label-sm label-w600 text-uppercase label-o4 px-2" @click="commitChanges">
        <div>{{ active ? "ON" : "OFF" }}</div>
      </div>
    </div>


  </div>
</template>

<style scoped>

</style>
