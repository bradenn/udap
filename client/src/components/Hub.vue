<script>
import Dock from "./Dock.vue";
import App from "./App.vue";
import Loading from "./Loading.vue";
import Media from "./apps/Media.vue";
import Proc from "./Proc.vue";

export default {
  components: {Proc, Media, Loading, Dock, App},
  data() {
    return {
      uptime: "",
      modClock: 1,
      modDir: false,
      open: false,
      reloading: false,
      focus: "idle",
      progress: 0,
      rotation: 0,
    }
  },
  watch: {},
  created() {
    setInterval(this.tick, 100)
  },
  methods: {
    update() {
      this.recalculateFocus()

      let n = Math.round(((new Date() - new Date(this.$root.state.last)) / 2000) * 100)
      if (n >= 150) {
        this.focus = "lapse"
      } else if (Math.abs(n - this.progress) >= 25) {
        this.focus = "active"
        setTimeout(() => {
          this.focus = "idle"
        }, 512)
        this.progress = n
      } else {
        this.progress = n
      }

    },
    hardReload() {
      this.$root.state.waiting = true
      window.location.reload()
    },
    recalculateFocus() {
      /*  this.focus = `transform: scale(calc(0.50 + ${this.$root.state.accepting ? '0.50' : (Math.random() * 0.125)})) perspective(2rem);`*/
    },
    state() {
      if (this.focus === 'lapse' || !this.$root.connection.connected) {
        return 'focus-animate-lapse'
      } else if (this.$root.state.waiting) {
        return 'focus-animate-active'
      } else {
        return 'focus-animate-idle'
      }
    },
    tick() {
      this.update()
      this.modClock += this.modDir ? 1 : -1;
      if (this.modClock >= 4) {
        this.modDir = false
      } else if (this.modClock <= 1) {
        this.modDir = true
      }
    },
  },
}


</script>


<template>
  <div class="top d-flex justify-content-start align-items-center gap-1">
    <!--    <Media small>-->
    <!--    </Media>-->
    <div class="element d-flex align-items-center justify-content-center">
      <div>
        <Proc class="" style=" z-index: -1; opacity: 0.25;"></Proc>
      </div>
      <div class="v-sep"></div>
      <div class="focus-container px-1 " @click="hardReload">
        <div :class="`focus-animate-${this.focus}`" class="focus-inner label-o2">􀝝</div>
        <div class="focus-outer label-o4">􀝝</div>
      </div>
    </div>
  </div>


</template>

<style lang="scss" scoped>
.v-sep {
  height: 1.5rem;
  width: 1px;
  margin-inline: 0.5rem;
  border-radius: 1rem !important;
  background-color: rgba(255, 255, 255, 0.125)
}

</style>
