<script>
function blurBlurable() {

}

export default {
  name: "Switch",
  data() {
    return {
      current: {},
      waiting: false,
    }
  },
  props: {
    entity: Object,
  },
  beforeCreate() {
    blurBlurable()
  },
  created() {
  },
  watch: {
    'entity': function (d) {
      this.waiting = false
      this.current = this.entity.state
    },
    newState: function (d) {

    }
  },
  methods: {
    setState(s) {
      this.waiting = true
      this.$root.requestId("entity", "state", s, this.entity.id)
    },
    toObject(b) {

      return JSON.parse(atob(b))
    }
  }
}

</script>

<template>

  <div v-if="current" :class="current.on?'light-active':''"
       class="element no-select entity d-flex flex-column justify-content-between" v-on:click="setState(current)">
    <div :class="`${current.on?'active':''}`" class="entity-small">

      <div class="fill"></div>
      <div class="label-xxs label-o3 label-w400 px-2">
        {{ current.on ? this.entity.neural === 'control' ? '􀴿 On' : 'On' : 'Off' }}
      </div>
    </div>
    <div class="d-flex justify-content-between align-items-baseline">
      <div class="toggle-status  small text-uppercase text-muted">{{ current.on ? "ON" : "OFF" }}</div>
      <div class="toggle-meta font-monospace text-muted small"><span
          class="border-bottom opacity-50">{{ entity.module }} {{ entity.live }}</span></div>
    </div>
  </div>
  <div v-else :class="current.on > 0?'light-active':''"
       class="element entity d-flex flex-column justify-content-between">
    <div class="d-flex justify-content-between align-items-center align-content-center">
      <div class="d-flex">
        <i class="bi toggle-icon" style="font-size: 0.8em; line-height: 1.9em;"><span
            v-if="!current.on">􀡷</span><span v-else class="energy">􀡸</span></i>&nbsp;&nbsp;
        <div class="toggle-title">{{ entity.name.charAt(0).toUpperCase() + entity.name.slice(1) }}</div>
      </div>
      <div>


      </div>
    </div>
    <div class="d-flex justify-content-between align-items-baseline">
      <div class="toggle-status  small text-uppercase text-muted">%</div>
      <div class="toggle-meta font-monospace text-muted small"><span
          class="border-bottom opacity-50">{{ entity.module }} {{ entity.live }}</span></div>
    </div>
  </div>

</template>

<style scoped>
.working {
  font-size: 1em;
  opacity: 1;
  animation: surge 250ms;
}

@keyframes surge {
  0% {
    opacity: 0;
  }
  20% {
    opacity: 0;
  }
  50% {
    opacity: 0;
  }
  100% {
    opacity: 0.8;
  }
}

.light-active {
  background-color: rgba(255, 255, 255, 0.08);
}

.entity {
  color: rgba(255, 255, 255, 0.70);

}


.play {
  font-size: 3em !important;
  line-height: 1em;
}

.icon > .bi {
  text-shadow: 0 0 8px rgba(0, 0, 0, 0.125);
}


.draw {
  color: rgba(255, 255, 255, 0.5);
  text-shadow: 0 0 16px rgba(255, 103, 0, 0.8);
}

.control {

}

@keyframes controller {
  0% {

  }
  100% {

  }
}

.seek {
  flex: 1 1 auto;
  height: 4px;
  width: 100%;
  border-radius: 2px;
  background-color: rgba(255, 255, 255, 0.10);
  backdrop-filter: blur(64px);
}

.thumbnail {
  height: 100%;
  width: 48px;
  border-radius: 10px;
  aspect-ratio: 1 / 1;
  border: 1px solid;
}

.media {
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  align-items: start;
}

.controls > i {
  margin-inline: 0.25em;
  line-height: 1em;
  font-size: 2em;
}

.controls {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.icon {
  min-height: 58px;
  flex: 1 1 auto;
  height: 100%;
  border-radius: 12px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 24px;
  /*border: 1px solid rgba(255, 255, 255, 0.16);*/
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  background-color: rgba(255, 255, 255, 0.10);
  backdrop-filter: blur(64px);
  transition: scale 2s;
}

.multiline {
  flex-wrap: wrap;
  height: auto !important;
}

.dock {
  display: flex;
  flex-direction: row;
  justify-content: space-around;
  align-items: center;
  width: 100%;
  height: 76px;
  padding: 0.5em 0.5em;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.16);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12), 0 8px 24px rgba(0, 0, 0, 0.11);
  background-color: rgba(255, 255, 255, 0.10);
  backdrop-filter: blur(64px);
  opacity: 1;
  gap: 0.5em;
}
</style>
