<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import moment from 'moment'
import {onMounted, reactive} from "vue";
import axios from "axios";
import Loader from "@/components/Loader.vue";
import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";

let state = reactive({
  satellite: "GOES17",
  mode: "GEOCOLOR",
  section: "FD",
  currentImage: "/custom/1808x1808.jpg",
  lowRes: "",
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
  state.loading = true
  buildURL()
  setInterval(updateData, 1000)
})

const sections = [
  {
    name: "Full Disk",
    key: "FD",
    alt: ""
  },
  {
    name: "California",
    key: "SECTOR/PSW",
    alt: ""
  },
  {
    name: "BC",
    key: "SECTOR/sea",
    alt: ""
  },
  {
    name: "USA",
    key: "CONUS",
    alt: ""
  }
]

const viewModes = [{
  name: "Color",
  key: "GEOCOLOR"
}, {
  name: "Heatmap",
  key: "FireTemperature"
},
  {
    name: "Sandwich",
    key: "Sandwich"
  },
  {
    name: "Dust",
    key: "Dust"
  },
  {
    name: "Clouds",
    key: "DayCloudPhase"
  },
  {
    name: "Fog",
    key: "NightMicrophysics"
  },
  {
    name: "Red",
    key: "02",
    alt: "640nm"
  },
  {
    name: "Green",
    key: "03",
    alt: "860nm"
  },

  {
    name: "Blue",
    key: "01",
    alt: "470nm"
  },
  {
    name: "Air Mass",
    key: "AirMass",
    alt: "Water Vapor"
  },

  {
    name: "3.9µm",
    key: "07",
    alt: "3.9µm"
  },
  {
    name: "Ozone",
    key: "12",
    alt: "3.9um"
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

function selectSatellite(satellite: string): any {
  state.satellite = satellite
  buildURL()
}

function selectMode(mode: string): any {
  state.mode = mode
  buildURL()
}

function selectSection(section: string): any {
  state.section = section
  buildURL()
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
    state.loading = false
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
  state.loading = true
  state.lowRes = `/custom/1808x1808.jpg`
  downloadImage(`https://cdn.star.nesdis.noaa.gov/${state.satellite}/ABI/${state.section}/${state.mode}/${state.section === "FD" ? '1808x1808' : '1200x1200'}.jpg?ts=${state.lastDate}`)
}


</script>

<template>
  <div class="d-flex flex-row justify-content-start align-items-start align-items-start w-100">

    <div class="d-flex  p-2 px-1 w-100">

      <div class=" d-flex flex-row gap w-100">

        <!-- Earth -->
        <div
            class="d-flex justify-content-between align-items-center align-content-center flex-row flex-grow-1 h-100">

          <!-- Help Bar -->
          <div
              class="element d-flex flex-column  justify-content-start align-items-center gap-2 p-1 mt-3 px-2 py-2"
              style="width: 13rem; ">
            <!-- NOAA Logo, as per their TOS -->
            <div class="d-flex align-items-center gap-2">


              <img alt="noaa" class="noaa mx-0" src="/noaa.svg"/>
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
            </div>
            <div class="h-sep"></div>
            <div class="d-flex justify-content-around w-100">
              <!-- Last update from the headers of the photo -->
              <div class="flex-shrink-0">
                <div class="label-c2 label-w500 label-o5 lh-sm">Last updated</div>
                <div class="label-c2 label-w300 label-o4 lh-1">{{ state.lastUpdated }}</div>
              </div>

              <div class="v-sep"></div>
              <!-- Last update from the headers of the photo -->
              <div class="flex-shrink-0">
                <div class="label-c2 label-w500 label-o5 lh-sm">Next update</div>
                <div class="label-c2 label-w300 label-o4 lh-1">{{ state.progress }}</div>
              </div>
            </div>


            <div class="h-sep"></div>
            <!-- Last update from the headers of the photo -->
            <div class="w-100 px-3">
              <div class="label-c2 label-w500 label-o5 lh-sm">Size</div>
              <div class="label-c2 label-w300 label-o4 lh-1">
                {{ Math.round(state.currentSize / 1000 / 1000 * 100) / 100 }} MB
              </div>
            </div>


            <div class="h-sep"></div>


            <!-- Refresh the image manually -->
            <div
                class="subplot p-1 px-2 d-flex flex-row justify-content-center align-content-center align-items-center flex-grow-0"
                @click="buildURL()">
              <div class="label-c2 label-w500 label-o4"><i
                  class="fa-solid fa-arrow-rotate-right fa-fw"></i> Force Reload</div>

            </div>

          </div>
          <!-- Disk -->
          <div class="d-flex justify-content-center align-items-center align-content-center earth-background element">
            <div v-if="state.section === 'FD'">
              <div v-if="state.loading"
                   :style="`background-image: url('${state.currentImage}');`"
                   class="earth-full-disk">
                <Loader v-if="state.loading"></Loader>
              </div>
              <div v-else :style="`background-image: url('${state.currentImage}');`"
                   class="earth-full-disk">
              </div>
            </div>
            <div v-else>
              <div v-if="state.loading"
                   :style="`background-image: url('${state.currentImage}');`"
                   class="preview p-5 d-flex justify-content-center align-items-center">
                <Loader v-if="state.loading"></Loader>
              </div>

              <div v-else :style="`background-image: url('${state.currentImage}');`"
                   class="preview p-5">
              </div>
            </div>

          </div>
          <div class="d-flex flex-column gap w-25  flex-grow-0">
            <div class="label-o4 label-sm label-r label-w500 lh-1">Satellite</div>
            <Plot :cols="2" :rows="1">
              <Subplot v-for="sat in satellites" :active="state.satellite===sat.key" :alt="sat.alt"
                       :fn="() => selectSatellite(sat.key)"
                       :name="sat.name"
                       icon="satellite"></Subplot>
            </Plot>
            <div class="label-o4 label-sm label-r label-w500 lh-1">Prespective</div>
            <Plot :cols="2" :rows="2">
              <Subplot v-for="sect in sections" :active="sect.key === state.section || state.toggles.section"
                       :alt="sect.alt" :fn="() => selectSection(sect.key)"
                       :name="sect.name"
                       icon="satellite"></Subplot>
            </Plot>
            <div class="label-o4 label-sm label-r label-w500 lh-1">Wavelengths</div>
            <Plot :cols="3" :rows="2">
              <Subplot v-for="mode in viewModes" :active="state.mode === mode.key"
                       :fn="() => selectMode(mode.key)"
                       :name="mode.name"></Subplot>
            </Plot>

            <!-- Current url -->

          </div>

        </div>   <!-- Sidebar -->

      </div>

    </div>

  </div>

</template>

<style scoped>

.earth-background {
  border-radius: 100%;
  padding: 0.125rem;
}

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
  width: 24.5rem;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 0 4px 8px rgba(0, 0, 0, 0.025);
  aspect-ratio: 1808/1778;
  border-radius: 100%;
  background-color: rgba(76, 87, 101, 0.37);
  background-size: cover;
  background-position: top;
  transition: background-image 250ms ease-in-out;
}

.earth-full-disk:hover {


}
</style>