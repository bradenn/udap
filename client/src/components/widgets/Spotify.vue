<script lang="ts" setup>
import moment from "moment";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Remote} from "@/types";
import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";

export interface Spotify {
  title: string;
  cover: string;
  thumbnail: string;
  artists: string;
  album: string;
  progress: number;
  updated: string;
  duration: number;
  explicit: boolean;
  playing: boolean;
  popularity: number;
  volume: number;
  device: string;
}

let remote = inject("remote") as Remote

interface SpotifyState {
  metadata: Spotify
  playing: Attribute
  menu: boolean
  current: number
  interval: number
}

let state = reactive<SpotifyState>({
  metadata: {} as Spotify,
  playing: {} as Attribute,
  menu: false,
  interval: 0,
  current: 0
})

onMounted(() => {
  updateMetadata(remote.attributes)
  clearInterval(state.interval)
  state.interval = setInterval(updateTime, 1000)
  updateTime()
})

watchEffect(() => updateMetadata(remote.attributes))


function updateMetadata(attributes: Attribute[]) {
  let proto = attributes.find(a => a.key === 'current')

  if (!proto) return
  state.metadata = JSON.parse(proto.value) as Spotify

  state.playing = attributes.find(a => a.key === 'playing') || {} as Attribute


  return state.metadata
}

function updateTime() {
  if (state.playing.value === 'true') {
    state.current = state.metadata.progress + (new Date().valueOf() - new Date(state.metadata.updated).valueOf())
  } else {
    state.current = state.metadata.progress
  }
  return state.current
}

// Apply changes made to an attribute
function togglePlayback() {
  // Make the request to the websocket object
  state.playing.request = `${state.playing.value === 'true' ? 'false' : 'true'}`
  remote.nexus.requestAttribute(state.playing.entity, state.playing.key, state.playing.request)
}


</script>

<template>
  <div v-if="state.metadata.playing" class="element label-c2 flex-grow-1" @mousedown="state.menu = !state.menu">

    <div class="d-flex flex-row gap-1 justify-content-start align-items-center">
      <div v-if="state.metadata" class="">
        <div :style="`background-image: url('${state.metadata.thumbnail}')`" class="album-sm ">
        </div>
      </div>

      <div v-if="!state.menu" class="d-flex flex-column justify-content-between w-100 flex-grow-1">
        <div>
          <div class="label-c2 label-w600 label-o5 overflow-ellipsis">{{ state.metadata.title }}<span
              v-if="state.metadata.explicit" class="text-danger label-c3 px-1">E</span></div>
          <div class="label-c3 label-w400 label-o4 lh-1">{{ state.metadata.artists }}</div>
        </div>
        <div>
          <div class="d-flex flex-row justify-content-between label-c3 label-o3 label-w400">
            <div>
              {{ moment(state.current).format("m:ss") }}
            </div>
            <div class="label-c3 label-o3 label-w500 d-flex flex-row">
              <div v-for="i in Array(Math.ceil((state.metadata.popularity || 50)/20.0)).keys()" :key="i">
                •
              </div>
            </div>

            <div>
              {{ moment(state.metadata.duration).format("m:ss") }}
            </div>
          </div>

          <div class="timeline-sm">

            <div :style="`width:${100*(state.current/state.metadata.duration)}%;`" class="timeline-value">
            </div>
          </div>
        </div>

      </div>

      <Plot v-else :cols="3" :rows="1" class="w-100" @mousedown.stop>
        <Subplot :active="true" :fn="() => togglePlayback()" name="􀊎"></Subplot>
        <Subplot :active="true" :fn="() => togglePlayback()" :name="state.metadata.playing?'􀊆':'􀊄'"></Subplot>
        <Subplot :active="true" :fn="() => togglePlayback()" name="􀊐"></Subplot>

      </Plot>


    </div>
  </div>
  <div v-else class=" label-c2 flex-grow-1 h-100">

  </div>
</template>

<style lang="scss" scoped>

.overflow-ellipsis {
  max-width: 9rem !important;
  white-space: nowrap;
  overflow: clip !important;
  text-overflow: ellipsis !important;
}


.earth-full-disk-k {
  height: 12rem;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 0 4px 8px rgba(0, 0, 0, 0.025);
  aspect-ratio: 1808/1778 !important;
  border-radius: 100%;
  background-color: rgba(76, 87, 101, 0.37);
  background-size: cover;
  transition: background-image 250ms ease-in-out;

}

.night {
  background-color: rgba(255, 255, 255, 0.2) !important;
}

.label-subtext {
  font-size: 0.30rem;
  width: 6px;
  opacity: 0.5;
  line-height: 0.4rem;
}

.border-box {
  border: 1px solid white;
}

.temp-bar {
  width: 6px;
  border-radius: 3px;
  background-color: rgba(255, 255, 255, 0.5);
}

.widget {
  aspect-ratio: 122/12 !important;
  width: 12rem;
}

.sky-gradient {
  height: 100%;
  border-radius: 0.4rem;
  transition: all 1s ease-in;
}

</style>
