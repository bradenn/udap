<script>
import moment from "moment";
import Dock from "../Dock.vue";


export default {
  name: "Media",
  components: {Dock, Range},
  data() {
    return {
      playing: false,
      lastChange: new Date(),
      menu: false,
      entity: "",
      attribute: {},
      diff: 0,
      track: {
        title: "",
        artist: "",
        album: "",
        cover: "",
        thumbnail: "",
        explicit: false,
        popularity: 0,
      },
      device: {
        name: "",
        volume: 0
      },
      progress: {
        min: 0,
        max: 0,
        current: 0,
        buffered: 0,
        timer: null,
      }
    }

  },
  beforeCreate() {

  },
  created() {

    setInterval(this.tick, 50)
  },
  computed: {
    playback: function () {
      if (!this.playing) return -1
      return this.progress.buffered
    },
  },
  props: {
    small: Boolean,
  },
  watch: {
    '$root.media': {
      immediate: true,
      handler(media, s) {
        if (!media.entity) return
        let latest = media.attributes.find(v => v.key === "current")
        if (!latest) return
        this.entity = latest.entity
        let current = JSON.parse(latest.value)
        this.lastChange = new Date(current.updated)
        this.playing = current.playing
        this.playingAttribute = media.attributes.find(v => v.key === "playing")
        this.skipAttribute = media.attributes.find(v => v.key === "cmd")
        this.progress.max = current.duration
        this.progress.current = current.progress
        this.track.album = current.album
        this.track.artist = current.artists
        this.track.title = current.title
        this.track.cover = current.cover
        this.track.thumbnail = current.thumbnail
        this.track.popularity = current.popularity
        this.device.name = current.device
        this.device.volume = current.volume
      }, deep: true
    },
  },
  beforeDestroy() {

  },
  methods: {
    tick() {
      let diff = new Date().valueOf() - this.lastChange.valueOf()
      if (this.playing) {
        this.progress.buffered = this.progress.current + (diff % this.progress.max)
      }
    },
    msFormat(ms) {
      return moment(0).add(moment.duration(ms)).format("mm:ss")
    },
    togglePlay: function (e) {
      let pa = this.playingAttribute
      pa.request = `${!this.playing}`
      this.$root.requestId("attribute", "request", pa, this.entity)
    },
    skipPrevious: function (e) {
      let pa = this.skipAttribute
      pa.request = 'previous'
      this.$root.requestId("attribute", "request", pa, this.entity)
    },
    skipNext: function (e) {
      let pa = this.skipAttribute
      pa.request = 'next'
      this.$root.requestId("attribute", "request", pa, this.entity)
    },
  },
}
</script>

<template>

  <div v-if="small" class="">
    <div v-if="playing" class="d-flex justify-content-between element" style="width: 13rem;" @click="menu = !menu">
      <div class="d-flex align-items-start gap">
        <div>
          <div :style="`background-image: url(${track.thumbnail});`" class="media-content-sm"></div>
        </div>
        <div class="overflow-text">
          <div v-if="track.title" class="label-xxs lh-1 overflow-target">
            {{ track.title.length >= 22 ? track.title.substring(0, 22) + `...` : track.title }}
          </div>
          <div class="label-ys label-o3 d-flex justify-content-between w-100">
            <div>{{ track.artist }}</div>
          </div>
        </div>

      </div>

      <div class="d-flex justify-content-between align-items-center gap label-o1 label-xxs px-1">
        <span v-if="menu">􀆈</span><span v-else>􀆈</span>
      </div>
    </div>
    <div v-if="menu" class="element drop-full d-flex flex-column gap-2 py-2 px-2">
      <div class="d-flex justify-content-between align-items-center gap px-1">
        <div class="label-ys label-o3 size-fixed-sm">
          {{ msFormat(progress.buffered) }}
        </div>
        <div class="u-slider">
          <div :style="`width: ${(progress.buffered/progress.max)*100}%; transition: width 50ms ease;`"
               class="u-slider-inner">
          </div>
        </div>

        <div class="label-ys label-o3 size-fixed-sm">
          {{ msFormat(progress.max) }}
        </div>
      </div>
      <div class="d-flex justify-content-center align-items-center gap-4 px-4">
        <div class="label-r label-o3" v-on:click="skipPrevious">􀊍</div>
        <div class="label-r label-o3" v-on:click="togglePlay"><span v-if="playing">􀊅</span><span v-else>􀊃</span></div>
        <div class="label-r label-o3" v-on:click="skipNext">􀊏</div>
      </div>
    </div>
  </div>
  <div v-else class="d-flex flex-column justify-content-center align-items-center gap-1 pb-5 fogged h-100">
    <div v-if="playing" class="element d-flex align-items-center gap-2 w-50">
      <div :style="`background-image: url(${track.cover});`" class="media-content"></div>
      <div v-if="track.explicit" class="tag-danger">EXPLICIT</div>
      <div class="overflow-text  px-2">
        <div class="label-xl label-w500 label-o4 lh-1">
          {{ track.title }}
        </div>
        <div class="label-xs label-w400 label-o4">
          {{ track.artist }}
        </div>

      </div>
    </div>
    <div class="d-flex gap w-50">
      <Dock class="media-controls px-2 w-100">
        <div class="d-flex justify-content-between align-items-center gap">
          <div class="label-ys label-o3 size-fixed-sm">
            {{ msFormat(progress.buffered) }}
          </div>
          <div class="u-slider">
            <div :style="`width: ${(progress.buffered/progress.max)*99}%; transition: width 50ms ease;`"
                 class="u-slider-inner">
            </div>
          </div>

          <div class="label-ys label-o3 size-fixed-sm">
            {{ msFormat(progress.max) }}
          </div>
        </div>
        <span class="mx-2"
              style="width: 0.05rem; height: 1rem; border-radius: 1rem; background-color: rgba(255,255,255,0.1);"></span>
        <div class="dock-icon w-50" v-on:click="skipPrevious">􀊍</div>
        <div class="dock-icon w-100" v-on:click="togglePlay"><span v-if="playing">􀊅</span><span v-else>􀊃</span></div>
        <div class="dock-icon w-50" v-on:click="skipNext">􀊏</div>


      </Dock>

    </div>
  </div>
</template>

<style scoped>
.media-controls .dock-icon {
  font-size: 0.7rem;
  font-weight: 500;

}

.drop-full {
  position: absolute !important;
  top: 100%;
  left: 0;
  margin-top: 0.125rem;
  width: 13rem;

}


.fogged:before {
  z-index: -1;
  content: ' ';
  position: absolute;
  width: 30rem;
  height: 12rem;
  filter: blur(120px);
  background-color: rgba(0, 0, 0, 0.33);
}

.size-fixed-sm {
  width: 2rem;
  height: 0.95rem;
  text-align: center;
}

.slider-track {
  height: 0.6rem;
}

.media-content-sm {
  width: 2rem;
  height: 2rem;
  background-position: center;
  box-shadow: 0 0 4rem 0.25rem rgba(0, 0, 0, 0.25);
  background-size: cover;
  border-radius: 0.25rem;
}

.media-content {
  width: 7rem;
  height: 7rem;
  opacity: 1;
  background-position: center;
  box-shadow: 0 0 4rem 0.25rem rgba(0, 0, 0, 0.25);
  background-size: cover;
  border-radius: 0.25rem;
}

.media-large {
  width: 100%;
  border-radius: 1rem;
  padding: 1rem;
  background-color: rgba(27, 27, 27, 0.3);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.tag-danger {
  background-color: rgba(255, 82, 69, 0.6);
  border-radius: 0.25rem;
  width: 4rem;
  font-size: 0.6rem;
  font-weight: 700;
  backdrop-filter: blur(2px);
}

.whiteboard {

  margin-top: 0.125em;
  height: 100%;
  font-family: "SF Pro Text", sans-serif;
  font-size: 15px;
  font-weight: 400;
  padding: 0.5em 1em;

  border-radius: 10px;
  box-shadow: 0 0 32px 4px rgba(0, 0, 0, 0.06);

}
</style>
