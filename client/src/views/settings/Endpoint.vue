<script>
import Dock from "../../components/Dock.vue";


export default {
  components: {Dock},
  name: "Endpoint",
  data() {
    return {
      loaded: false,
      advanced: false,
      blurs: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
      padding: [1, 2, 3, 4, 5, 6, 7, 8],
      backgrounds: ["bells", "viridian", "nerfdart", "nasa", "milk", "hole"],
      colors: [
        {name: "Slate", cls: "accent-slate"},
        {name: "Red", cls: "accent-red"},
        {name: "Blue", cls: "accent-blue"},
        {name: "Green", cls: "accent-green"},
        {name: "Orange", cls: "accent-orange"}],
      inputs: [
        {name: "Touch", value: "touchscreen"},
        {name: "Mouse", value: "mouse"}
      ],
      preferences: [{
        name: "Input Style",
        key: "input",
        icon: "􀇳",
        default: "hid",
      }],
    }
  },
  computed: {
    session: function () {
      return {
        theme: {

          name: "Material Theme",
          icon: "􀣶",
          key: "theme",
          value: this.$root.preferences.theme,
          type: "select",
          options: [{name: "Dark", operation: "", value: "dark"}, {name: "Light", value: "light"}]
        },
        accent: {
          advanced: true,
          name: "Accent Color",
          icon: "􀟗",
          key: "accent",
          value: this.$root.preferences.accent,
          type: "select",
          options: [{name: "Slate", value: "slate"},
            {name: "Red", value: "red"},
            {name: "Blue", value: "blue"},
            {name: "Green", value: "green"},
            {name: "Orange", value: "orange"}]
        },
        padding: {
          name: "Padding",
          advanced: true,
          icon: "􀫼",
          key: "padding",
          value: this.$root.preferences.padding,
          type: "range",
          range: {
            cls: 'slider-mono',
            unit: 'pt',
            min: 1,
            value: 0,
            max: 8,
            step: 1
          }
        },
        picture: {
          name: "Background",
          icon: "􀞲",
          key: "background",
          value: this.$root.preferences.background,
          type: "select",
          options: [{name: "Bella", value: "bells"},
            {name: "Viridian", value: "viridian"},
            {name: "Milk", value: "milk"},
            {name: "Nerfdart", value: "nerfdart"},
            {name: "Daemon", value: "daemon"},
            {name: "Black Hole", value: "cblue"}]
        },
        blur: {
          name: "Blur",
          icon: "􀯲",
          key: "blur",
          value: this.$root.preferences.blur,
          type: "range",
          range: {
            cls: 'slider-mono',
            unit: 'pt',
            min: 1,
            value: 0,
            max: 10,
            step: 1
          }
        },
        input: {
          advanced: true,
          name: "Input Style",
          icon: "􀞲",
          key: "input",
          value: this.$root.preferences.input,
          type: "select",
          options: [{name: "Pseudotouch", value: "touchscreen"},
            {name: "Touch", value: "touch"},
            {name: "Mouse", value: "mouse"}]
        }
      }
    }
  },
  methods: {
    toggleAdvanced() {
      this.advanced = !this.advanced
    },
    set: function (key, value) {
      this.$root.preferences[key] = value
    },
    load: function (e) {
      e.target.style.backdropFilter = `blur(${this.$root.preferences.blur ** 2}px)`
    }
  },
}
</script>

<template>
  <div class="d-flex flex-column context-container w-100">
    <div class="context-container">
      <div class="d-flex justify-content-between align-items-end align-content-end mb-1" v-on:click.stop>
        <div class="label-xxl label-w600 label-o5">
          <span class="label-xl label-w600 label-o4">􀏞</span>
          <span class="label-xl label-w600 label-o6 px-2">Endpoint</span>
        </div>
      </div>
    </div>
    <div class="frame-container">
      <div v-for="attribute in Object.values(this.session)"
           :key="attribute.key"
           class="element" v-on:click.stop>
        <div class="h-bar justify-content-start align-items-center align-content-center pb-1">
          <div class="label-xxs label-o2 label-w600">{{ attribute.icon }}</div>
          <div class="label-xxs label-o4 label-w500">&nbsp;&nbsp;{{ attribute.name }}</div>
          <div class="fill"></div>
          <div v-if="attribute.range" class="h-bar gap label-xxs label-o3 label-w400 px-2">
            <div class="label-xxs label-o3">{{ attribute.value }} {{ attribute.range.unit }}</div>
          </div>
        </div>
        <div class="d-flex gap ">
          <input v-if="attribute.type==='range'" :key="attribute.key"
                 v-model="attribute.value"
                 :class="`${attribute.range.cls}`"
                 :max="attribute.range.max"
                 :min="attribute.range.min" :step="attribute.range.step"
                 :style="attribute.key==='level'?`background: linear-gradient(90deg, rgba(255, 255, 255, 0) 0%, 100%);`:''"
                 class=" slider element surface"
                 type="range"
                 @mousemove="this.$root.preferences[attribute.key] = attribute.value">
          <Dock v-else-if="attribute.type==='select'" class="surface w-100" small>
            <div v-for="option in attribute.options"
                 :class="`${attribute.value === option.value?'router-link-active':''}`"
                 class="dock-link"
                 href="#" @click="this.$root.preferences[attribute.key] = option.value">
              {{ option.name }}
            </div>
          </Dock>

        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.frame-container .element {

}

.frame-container {
  display: grid;
  width: 100%;
  gap: 0.375rem;
  grid-template-columns: repeat(auto-fit, minmax(22rem, 1fr));
}

.option-grid {
  display: grid;
  max-width: calc(100%);
  grid-template-columns: repeat(auto-fill, minmax(26rem, 1fr));
  grid-auto-flow: row;
  flex-wrap: wrap;
}

.bg-active {

}

.toolbar {
  margin-bottom: 8px;
}

</style>
