<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import moment from 'moment'
import {onMounted, reactive} from "vue";
import axios from "axios";

let state = reactive({
  satellite: "GOES17",
  mode: "GEOCOLOR",
  section: "FD",
  currentImage: "",
  lastUpdated: "",
  lastDate: 0,
  progress: "",
  currentSize: 0,
  nextUpdate: "",
  toggles: {
    satellite: false,
    section: false,
    mode: false,
  },
  loading: false,
})

onMounted(() => {
  setInterval(updateData, 1000)
})

const sections = [
  {
    name: "Full Disk",
    key: "FD",
  },
  {
    name: "California",
    key: "SECTOR/PSW",
  },
  {
    name: "USA",
    key: "CONUS",
  }
]

const viewModes = [{
  name: "Color",
  key: "GEOCOLOR"
}, {
  name: "Heatmap",
  key: "FireTemperature"
}, {
  name: "Air Mass",
  key: "AirMass"
},
  {
    name: "Combined",
    key: "Sandwich"
  },
  {
    name: "Co2",
    key: "16"
  }]

const satellites = [
  {
    name: "GOES 17",
    alt: "West",
    key: "GOES17",
  },
  {
    name: "GOES 16",
    alt: "East",
    key: "GOES16",
  }
]

function selectSatellite(satellite: string) {
  state.satellite = satellite
}

function selectMode(mode: string) {
  state.mode = mode
}

function selectSection(section: string) {
  state.section = section
}

function downloadImage(url: string) {
  axios.get(url).then(res => {
    let lastModified: Date = new Date(res.headers['last-modified'])
    state.lastDate = lastModified.valueOf()
    state.currentSize = res.headers['content-length']
    state.lastUpdated = moment(lastModified).fromNow(false)
    state.nextUpdate = moment(lastModified).add(1000 * 60 * 10).fromNow(false)
    updateData()
    state.currentImage = url
  }).catch(err => {

  })
}

function updateData() {
  let lastModified: Date = new Date(state.lastDate)
  let diff = moment(lastModified).add(1000 * 60 * 10).subtract(moment().valueOf())
  if (new Date().valueOf() - lastModified.valueOf() + 1000 * 60 * 10 % 1000 === 1) {
    buildURL()
  }
  state.progress = diff.format("m[m] s[s]")
}

function buildURL() {
  // Mode ABI, Advanced Baseline Imager
  // Section, FD, etc
  // Mode GEOCOLOR, etc

  downloadImage(`https://cdn.star.nesdis.noaa.gov/${state.satellite}/ABI/${state.section}/${state.mode}/${state.section === "FD" ? '1808x1808' : '1200x1200'}.jpg?ts=${state.lastDate}`)
  return state.currentImage
}


</script>

<template>
  <div class="d-flex flex-row justify-content-start align-items-start align-items-start w-100">

    <div class="d-flex  p-2 px-1 w-100">

      <div class=" d-flex flex-row gap w-100">
        <!-- Sidebar -->
        <div class="d-flex flex-column gap w-25  flex-grow-0">
          <div class="d-flex justify-content-start px-1 pt-0">
            <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-earth-americas fa-fw`"></i></div>
            <div class="label-w500  label-xxl px-2">Earth</div>
          </div>
          <div class="element">
            <div class="d-flex justify-content-between mx-1 align-items-center"
                 @click="state.toggles.satellite = !state.toggles.satellite">
              <div class="h5">Satellite</div>
              <div class="label-c2 label-w500 label-o4">{{ satellites.length }} options <i
                  class="fa-solid fa-caret-down label-o5"></i></div>
            </div>
            <div v-for="sat in satellites"
                 :key="sat.key"

                 :class="`${state.satellite === sat.key?'router-link-active':''}`"
                 class="gap">
              <div v-if="sat.key === state.satellite || state.toggles.satellite" class="option"
                   @click="selectSatellite(sat.key)">
                {{ sat.name }}<span class="float-end label-o3 text-uppercase">{{ sat.alt }}</span>
              </div>
            </div>

          </div>
          <div class="element">
            <div class="d-flex justify-content-between mx-1 align-items-center"
                 @click="state.toggles.section = !state.toggles.section">
              <div class="h5">Perspective</div>
              <div class="label-c2 label-w500 label-o4">{{ sections.length }} options <i
                  class="fa-solid fa-caret-down label-o5"></i></div>
            </div>
            <div v-for="sect in sections" :key="sect.key"

                 :class="`${state.section === sect.key?'router-link-active':''}`"
                 class="gap">
              <div v-if="sect.key === state.section || state.toggles.section" class="option"
                   @click="selectSection(sect.key)">{{ sect.name }}
              </div>
            </div>
          </div>
          <div class="element">
            <div class="d-flex justify-content-between mx-1 align-items-center"
                 @click="state.toggles.mode = !state.toggles.mode">
              <div class="h5">Wavelength</div>
              <div class="label-c2 label-w500 label-o4">{{ viewModes.length }} options <i
                  class="fa-solid fa-caret-down label-o5"></i></div>
            </div>
            <div v-for="mode in viewModes" :key="mode.key"

                 :class="`${state.mode === mode.key?'router-link-active':''}`"
                 class="gap">
              <div v-if="mode.key === state.mode || state.toggles.mode" class="option"
                   @click="selectMode(mode.key)">{{ mode.name }}
              </div>
            </div>
          </div>


        </div>
        <!-- Earth -->
        <div
            class="d-flex justify-content-between align-items-center align-content-center flex-column flex-grow-1 h-100">

          <!-- Disk -->
          <div class="d-flex justify-content-center align-items-center align-content-center h-100">
            <div v-if="state.section === 'FD'"
                 :style="`background-image: url('${buildURL()}');`"
                 class="earth-full-disk">
            </div>
            <div v-else :style="`background-image: url('${buildURL()}');`"
                 class="preview p-5">
            </div>
          </div>

          <!-- Help Bar -->
          <div class="element d-flex flex-row justify-content-start align-items-center gap-1 p-1 mt-4">
            <!-- NOAA Logo, as per their TOS -->
            <img alt="noaa" class="noaa mx-1" src="/noaa.svg"/>

            <!-- Last update from the headers of the photo -->
            <div class="flex-shrink-0">
              <div class="label-c2 label-w500 label-o5 lh-sm">Last updated</div>
              <div class="label-c2 label-w300 label-o4 lh-1">{{ state.lastUpdated }}</div>
            </div>

            <div class="v-sep"></div>

            <!-- Last update from the headers of the photo -->
            <div class="flex-shrink-0">
              <div class="label-c2 label-w500 label-o5 lh-sm">Next update</div>
              <div class="label-c2 label-w300 label-o4 lh-1 monospace">{{ state.progress }}</div>
            </div>

            <div class="v-sep"></div>

            <!-- The path tp the current photo -->
            <div class="flex-shrink-0">
              <div class="label-c2 label-w500 label-o5 lh-sm">Path</div>
              <div class="d-flex gap-2 align-items-center justify-content-center flex-grow-1">

                <div class="d-flex gap">
                  <div class="label-c2 label-w300 label-o4 lh-1">{{ state.satellite }}</div>
                  <div class="label-c2 label-o4 lh-1"><i class="fa-solid fa-caret-right"></i></div>
                </div>
                <div class="d-flex gap">
                  <div class="label-c2 label-w300 label-o4 lh-1">{{ state.section }}</div>
                  <div class="label-c2 label-o4 lh-1"><i class="fa-solid fa-caret-right"></i></div>
                </div>
                <div class="d-flex gap">
                  <div class="label-c2 label-w300 label-o4 lh-1">{{ state.mode }}</div>
                </div>
              </div>
            </div>
            <div class="v-sep"></div>
            <!-- Last update from the headers of the photo -->
            <div class="flex-shrink-0">
              <div class="label-c2 label-w500 label-o5 lh-sm">Size</div>
              <div class="label-c2 label-w300 label-o4 lh-1 monospace">
                {{ Math.round(state.currentSize / 1000 / 1000 * 100) / 100 }} MB
              </div>
            </div>


            <div class="v-sep"></div>


            <!-- Refresh the image manually -->
            <div
                class="dock-icon surface button-icon d-flex flex-row justify-content-center align-content-center align-items-center flex-grow-0"
                @click="buildURL()">
              <div class="label-c2 label-w500 label-o5"><i class="fa-solid fa-arrow-rotate-right fa-fw"></i></div>

            </div>

          </div>
          <!-- Current url -->
          <div class="label-c3 label-o4 pt-2">
            {{ state.currentImage }}
          </div>
        </div>
      </div>
    </div>
  </div>

</template>

<style scoped>
.noaa {
  height: 1.5rem;
  padding: 1px;
  background-color: white;
  border-radius: 1rem;
}

.button-icon {
  aspect-ratio: 1/1 !important;

}

.option {
  font-size: 0.625rem;
  line-height: 0.625rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);

  border-radius: 0.375rem;
  padding: 0.4rem 0.5rem;


}

.router-link-active > .option {
  font-size: 0.625rem;
  line-height: 0.625rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
  background-color: rgba(255, 255, 255, 0.02);
  border-radius: 0.375rem;
  padding: 0.4rem 0.5rem;


}

.outline {
  border-radius: 100%;
}

.preview {
  width: 22rem !important;
  height: 22rem;
  border-radius: 1rem;
  background-size: cover;

}

.preview:hover {
  transform: scale(1.2);
}

.earth-full-disk {
  width: 24rem;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 0 4px 8px rgba(0, 0, 0, 0.025);
  aspect-ratio: 1808/1778;
  border-radius: 100%;
  background-color: rgba(76, 87, 101, 0.37);
  background-size: cover;
  transition: background-image 250ms ease-in-out;
}

.earth-full-disk:hover {


}
</style>