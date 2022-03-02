<script>

import Loading from "../Loading.vue";

export default {
  name: "Range",
  components: {Loading},
  data() {
    return {
      latest: new Date(),
      waiting: false,
      sliding: false,
      loading: false,
      local: {},
      config: {
        dim: {
          name: 'Intensity',
          min: 0,
          max: 100,
          step: 5,
          unit: '%',
          icon: '􀇯',
          class: 'slider-dim',
        },
        cct: {
          name: 'Warmth',
          min: 2000,
          max: 8000,
          step: 100,
          unit: 'K',
          icon: '􀍽',
          class: 'slider-cct',
        },
        hue: {
          name: 'Color',
          min: 1,
          max: 360,
          step: 1,
          unit: '°',
          icon: '􀟗',
          class: 'slider-hue',
        },
        main: {
          min: 1,
          max: 100,
          step: 5,
          unit: '%',
          icon: '􀌆',
          class: 'slider-dim',
        }
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
    primitive: Boolean,
  },
  computed: {},
  watch: {
    'attribute': {
      immediate: true,
      handler(d, a) {
        if (this.sliding) return
        this.local = d
        this.waiting = false
      }
    }
  },
  methods: {
    slideStart: function (e) {
      this.sliding = true
    },
    commitChanges: function (e) {
      this.commit(this.local)
      this.waiting = true
      this.latest = new Date()
      // Prevent updates until after we send the state
      this.sliding = false
    }
  }
}

</script>

<template>
  <div v-if="small">
    <input v-if="attribute.type==='range'" :key="attribute.id" v-model="this.local.request"
           :class="`slider-${attribute.key}`"
           :max=config[attribute.key].max :min=config[attribute.key].min
           :step=config[attribute.key].step
           class="dock dock-xsmall slider"
           type="range"
           v-on:mousedown="slideStart"
           v-on:mouseup="commitChanges">
  </div>
  <div v-else-if="primitive">
    <div class="h-bar justify-content-start align-items-center align-content-center pb-1">
      <div class="label-xxs label-o2 label-w600">{{ config[attribute.key].icon }}</div>
      <div class="label-xxs label-o4 label-w500">&nbsp;&nbsp;{{ config[attribute.key].name }}</div>
    </div>
    <input v-if="attribute.type==='range'" :key="attribute.id" v-model="this.local.request"
           :class="`slider-${attribute.key}`"
           :max=config[attribute.key].max :min=config[attribute.key].min
           :step=config[attribute.key].step
           class="slider element"
           type="range"
           v-on:mousedown="slideStart"
           v-on:mouseup="commitChanges">

  </div>
  <div v-else class="element" v-on:click.stop>

    <div class="h-bar justify-content-start align-items-center align-content-center pb-1">
      <div class="label-xxs label-o2 label-w600">{{ config[attribute.key].icon }}</div>
      <div class="label-xxs label-o4 label-w500">&nbsp;&nbsp;{{ config[attribute.key].name }}</div>
      <div class="fill"></div>
      <div class="h-bar gap label-xxs label-o3 label-w400 px-2">
        <Loading v-if="waiting"></Loading>
        <div class="label-xxs label-o3">{{ this.local.request }} {{ config[attribute.key].unit }}</div>
      </div>
    </div>
    <div class="d-flex gap">
      <!--      <div class="slider-ticks">
              <div v-for="a in [...Array(22).keys()]"></div>

            </div>-->
      <input v-if="attribute.type==='range'" :key="attribute.id" v-model="this.local.request"
             :class="`slider-${attribute.key}`"
             :max=config[attribute.key].max :min=config[attribute.key].min
             :step=config[attribute.key].step
             class="slider element"
             type="range"
             v-on:mousedown="slideStart"
             v-on:mouseup="commitChanges">

    </div>

  </div>
</template>

<style scoped>

</style>
