<script>
import Loading from "../Loading.vue";
import Context from "../Context.vue";
import History from "../History.vue";
import Attribute from "./Attribute.vue";
import Range from "./Range.vue";


function map_range(value, low1, high1, low2, high2) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function cctToRgb(cct) {
  return [map_range(cct, 2000, 8000, 255, 213),
    map_range(cct, 2000, 8000, 160, 230),
    map_range(cct, 2000, 8000, 75, 255)]
}

/**
 * Converts an RGB color value to HSL. Conversion formula
 * adapted from http://en.wikipedia.org/wiki/HSL_color_space.
 * Assumes r, g, and b are contained in the set [0, 255] and
 * returns h, s, and l in the set [0, 1].
 *
 * @param   Number  r       The red color value
 * @param   Number  g       The green color value
 * @param   Number  b       The blue color value
 * @return  Array           The HSL representation
 */
function rgbToHsl(r, g, b) {
  r /= 255, g /= 255, b /= 255;

  var max = Math.max(r, g, b), min = Math.min(r, g, b);
  var h, s, l = (max + min) / 2;

  if (max == min) {
    h = s = 0; // achromatic
  } else {
    var d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

    switch (max) {
      case r:
        h = (g - b) / d + (g < b ? 6 : 0);
        break;
      case g:
        h = (b - r) / d + 2;
        break;
      case b:
        h = (r - g) / d + 4;
        break;
    }

    h /= 6;
  }

  return [h, s, l];
}

/**
 * Converts an HSL color value to RGB. Conversion formula
 * adapted from http://en.wikipedia.org/wiki/HSL_color_space.
 * Assumes h, s, and l are contained in the set [0, 1] and
 * returns r, g, and b in the set [0, 255].
 *
 * @param   Number  h       The hue
 * @param   Number  s       The saturation
 * @param   Number  l       The lightness
 * @return  Array           The RGB representation
 */
function hslToRgb(h, s, l) {
  var r, g, b;

  if (s == 0) {
    r = g = b = l; // achromatic
  } else {
    function hue2rgb(p, q, t) {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1 / 6) return p + (q - p) * 6 * t;
      if (t < 1 / 2) return q;
      if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
      return p;
    }

    var q = l < 0.5 ? l * (1 + s) : l + s - l * s;
    var p = 2 * l - q;

    r = hue2rgb(p, q, h + 1 / 3);
    g = hue2rgb(p, q, h);
    b = hue2rgb(p, q, h - 1 / 3);
  }

  return [r * 255, g * 255, b * 255];
}


export default {
  components: {Range, Attribute, History, Context, Loading},
  data() {
    return {
      values: {
        level: 0,
        cct: 0,
        color: 0,
      },
      lastChange: new Date(),
      double: false,
      queue: [],
      waiting: false,
      localCCT: 4500,
      localBrightness: 0,
      waitState: {},
      holding: false,
      holdValue: 0,
      context: false,
      advanced: false,
      sliding: false,
      mode: "",
      state: {},
      attributes: {
        dim: 0,
        cct: 0,
        hue: 0,
        on: false
      },
    }
  },
  props: {
    entity: Object,
    manage: Boolean,
    selected: Boolean,
    small: Boolean,
    bulb: Boolean,
    diagnostic: Boolean,
    powerLevel: Boolean,
    bulk: Boolean,
    level: Number,
    group: Array,
  },
  computed: {
    attributeManifest: function () {
      function cmp(a, b) {
        return a.order - b.order
      }

      let attrs = this.$root.attributes.filter(a => a.entity === this.entity.id && a.key !== 'on')

      return attrs.sort(cmp)
    },
    client: function () {
      return {
        neural: {
          name: "Neural",
          icon: "􀴿",
          key: "neural",
          operation: "neural",
          value: this.entity.neural,
          type: "select",
          options: [{name: "Control", value: "control"}, {name: "Suggest", value: "suggest"}, {
            name: "Inactive",
            value: "inactive"
          }]
        },
        brightness: {
          name: "Brightness",
          translate: 'dim',
          icon: "􀇯",
          key: "level",
          mode: "brightness",
          operation: "state",
          value: this.attributes.dim,
          type: "range",
          range: {
            cls: '',
            unit: '%',
            value: 0,
            min: 0,
            max: 100,
            step: 5
          }
        },
        cct: {
          name: "Temperature",
          icon: "􀍽",
          translate: 'cct',
          key: "cct",
          mode: "cct",
          value: this.attributes.cct,
          type: "range",
          range: {
            cls: 'slider-cct',
            unit: 'K',
            min: 2000,
            max: 8000,
            step: 100
          }
        },
        color: {
          name: "Color",
          icon: "􀎗",
          translate: 'hue',
          key: "color",
          mode: "color",
          value: this.attributes.hue,
          type: "range",
          range: {
            cls: 'slider-color',
            unit: '°',
            min: 0,
            value: 0,
            max: 360,
            step: 1
          }
        },
        power: {
          shown: "settings",
          name: "Power",
          icon: "􀡷",
          key: "power",
          operation: "state",
          value: this.attributes.on,
          type: "select",
          options: [{name: "On", value: true}, {name: "Off", value: false}]
        },
        locked: {
          shown: "settings",
          name: "Locked",
          icon: "􀎡",
          key: "locked",
          operation: "locked",
          value: this.entity.locked,
          type: "select",
          options: [{name: "􀎥", value: false}, {name: "􀎡", value: true}]
        },
        icon: {
          shown: "settings",
          name: "Icon",
          icon: "􀎡",
          key: "icon",
          operation: "icon",
          value: this.entity.icon,
          type: "select",
          options: [{name: "􀛮", value: "􀛮"},
            {name: "􀡸", value: "􀡸"},
            {name: "􁁌", value: "􁁌"},
            {name: "􀇯", value: "􀇯"},
            {name: "􀇮", value: "􀇮"}]
        },
        frequency: {
          shown: "settings",
          name: "Frequency",
          icon: "􀎡",
          key: "frequency",
          operation: "frequency",
          value: this.entity.frequency,
          type: "select",
          options: [{name: "250ms", value: 250}, {name: "1s", value: 1000}, {name: "3s", value: 3000}, {
            name: "10s",
            value: 10000
          }]
        },
      }
    },
    slideStart: function () {
      this.sliding = true;

    },
    slideStop: function () {
      this.sliding = false;

    },
    currentColor: function () {
      let defaultColor = cctToRgb(6000)
      if (this.entity.type !== "spectrum") return `rgba(${defaultColor[0]},${defaultColor[1]},${defaultColor[2]},${this.active ? 0.5 : 0})`
      let rgb = [], intensity = 0;
      let attrs = this.attributeManifest.sort((a, b) => new Date(a.requested).valueOf() - new Date(b.requested).valueOf())
      let latest = 0
      for (const attr of attrs) {
        let value = Number(attr.request)
        let ts = new Date(attr.updated).getTime()
        switch (attr.key) {
          case 'cct':
            if (value <= 0) return
            if (ts >= latest) {
              rgb = cctToRgb(value)
              latest = ts
            }
            break
          case 'hue':
            if (value <= 0) return
            if (ts >= latest) {
              rgb = hslToRgb(value / 360, 1, 0.5)
              latest = ts
            }
            break
          case 'dim':
            intensity = value > 0 ? value / 100 : 0
            break
          default:
            break
        }
      }
      return `rgba(${rgb[0]},${rgb[1]},${rgb[2]},${this.active ? intensity : 0})`
    },
    isOn: function () {
      let c;
      c = this.$root.attributes.find(a => a.entity === this.entity.id && a.key === 'on')
      return c ? c.value === 'true' : false;

    },
    active: function () {
      return this.$root.attributes.find(a => a.entity === this.entity.id && a.key === 'on')
    },
  },
  methods: {
    mouseDown() {
      this.double++
      if (this.double >= 2) {
        this.context = true
      }
      setTimeout(function () {
        this.double--
      }, 500)
    },
    pushStop: function (e) {
      this.sliding = false
    },
    pushStart: function (e) {
      this.sliding = true
    },
    toggleMenu() {
      this.context = !this.context
    },
    toggleAdvanced() {
      this.advanced = !this.advanced
    },
    computeClasses() {
      let cls = [];
      if (this.active) cls.push("element")
      if (this.manage) cls.push("manage")
      return [].join(" ")
    },
    goTo() {
      this.$router.replace(`/terminal/settings/entities/${this.entity}`)
    },
    commitChange(attribute) {
      this.$root.requestId("attribute", "request", attribute, this.entity.id)
    }

  }

}


</script>

<template>
  <div v-if="bulb">
    <span :style="active?`text-shadow: 0 0 4px ${currentColor};`:''" class="label-sm"><span
        v-if="selected">􀛮</span><span v-else>􀛭</span></span>
    <span class="label-xxs px-1 label-w500">{{ entity.name }}</span>
  </div>
  <div v-else-if="manage">
    <div :class="isOn?'active':''" class="entity-small">
      <div class="entity-header">
        <div :style="`text-shadow: 0 0 4px ${currentColor};`" class="icon">
          {{ entity.icon ? entity.icon : '􀛮' }}
        </div>
        <div class="label-xxs label-w400 label-o5">
          {{ entity.name }}
        </div>
      </div>
      <div class="fill"></div>
      <div class="d-flex align-items-center">
        <div class="label-xxs label-o3 label-w300 px-1">
          <span v-if="selected">􀝜</span><span v-else>􀀀</span>
        </div>
      </div>

    </div>
  </div>
  <div v-else-if="small">
    <div :class="isOn?'active':''" class="entity-small" @click="toggleMenu()">
      <div class="entity-header">
        <div class="icon">
          {{ entity.icon ? entity.icon : '􀛮' }}
        </div>
        <div class="label-xxs label-w400 label-o4 px-1">
          {{ entity.name }}
        </div>
      </div>
      <div class="fill"></div>
      <div class="label-xxs label-o3 label-w400 px-2">
        <div v-if="isOn">ON</div>
        <div v-else>OFF</div>
      </div>
    </div>
    <div v-if="context" class="context"></div>
    <div v-if="context" @click="toggleMenu">
      <div class="entity-context top d-flex">
        <div class="d-flex flex-column gap px-3 top ">
          <div class="d-flex justify-content-start align-items-end align-content-end" v-on:click.stop>
            <div>
              <span :style="`text-shadow: 0 0 8px ${currentColor};`"
                    class="label-md label-w600 label-o3">{{ this.entity.icon }}</span>
              <span class="label-md label-w600 label-o6 px-2">{{ this.entity.name }}</span>
            </div>
            <div class="fill"></div>
            <div class="h-bar">
              <Attribute v-if="this.active" :attribute="this.active" small></Attribute>
              <div class="mx-1 my-1"
                   style="width: 0.0625rem; border-radius: 1rem; background-color: rgba(255,255,255,0.2);"></div>
              <div class="label-sm label-w600 label-o3 px-2 text-uppercase" @click="toggleAdvanced()">
                <span v-if="advanced">􁅦</span> <span v-else>􀍟</span>
              </div>
            </div>
          </div>

          <div class="context-container gap v-bar">
            <div v-for="attribute in attributeManifest.sort((a, b) => a.order - b.order)" :key="attribute.id">
              <Attribute :key="attribute.id" :attribute="attribute" :entity="this.entity.id" primitive></Attribute>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-else>
    <div>
      <div :class="`${this.active?'router-link-active':''}`" class="element entity"
           v-on:click="">
        <div class="toolbar">
          <div class="toolbar-icon">{{ entity.icon ? entity.icon : "􀛭" }}</div>
          <div class="toolbar-title">&nbsp;&nbsp;{{ entity.name }}</div>
          <div class="fill"></div>
          <span v-if="waiting" class="working text-muted mx-1">􀋧</span>
          <span v-else-if="!entity.live" class="text-warning">􀇿</span>
          <span v-else class="text-uppercase mx-1">{{ active ? "On" : "Off" }}</span>
        </div>
      </div>
    </div>

  </div>
</template>

<style lang="scss" scoped>

</style>
