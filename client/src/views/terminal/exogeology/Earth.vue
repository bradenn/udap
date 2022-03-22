<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import Sidebar from "@/components/sidebar/Sidebar.vue";
import SidebarItem from "@/components/sidebar/SidebarItem.vue";
import moment from 'moment'
import {reactive} from "vue";
import axios from "axios";

let state = reactive({
  satellite: "GOES17",
  mode: "GEOCOLOR",
  section: "FD",
  currentImage: "",
  lastUpdated: "",
  lastDate: 0,
  nextUpdate: "",
  loading: false,
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
    state.currentImage = url
  }).catch(err => {

  })
}

function downloadSha(url: string) {
  axios.get(url).then(res => {

    state.lastUpdated = moment(new Date(res.headers['last-modified'])).fromNow(false)
    state.nextUpdate = moment(new Date(res.headers['last-modified'])).add(1000 * 60 * 15).fromNow(false)
  })
}

function buildURL(satellite: string, section: string, mode: string) {
  // Mode ABI, Advanced Baseline Imager
  // Section, FD, etc
  // Mode GEOCOLOR, etc
  downloadSha(`https://cdn.star.nesdis.noaa.gov/${satellite}/ABI/${section}/${mode}/${satellite}-ABI-${section}-${mode}-10848x10848.tif.sha256`)
  downloadImage(`https://cdn.star.nesdis.noaa.gov/${satellite}/ABI/${section}/${mode}/${section === "FD" ? '1808x1808' : '1200x1200'}.jpg`)
  return state.currentImage
}


</script>

<template>
  <div class="d-flex flex-column w-100">
    <div class="d-flex justify-content-start p-2 px-1">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-earth-americas fa-fw`"></i></div>
      <div class="label-w500  label-xxl px-2">Earth</div>
    </div>
    <div class="">
      <div class=" d-flex flex-row gap-4">
        <div class="d-flex flex-column gap w-25">
          <!-- Satellite -->
          <Sidebar class="w-100"
                   icon="satellite"
                   name="Satellite"
                   small>
            <div v-for="sat in satellites"
                 :key="sat.key"
                 :class="`${state.satellite === sat.key?'router-link-active':''}`"
                 class="macro-icon-default">
              <SidebarItem :name="`${sat.name} (${sat.alt})`"
                           @click="selectSatellite(sat.key)"></SidebarItem>
            </div>
          </Sidebar>
          <!-- Perspective -->
          <Sidebar class="w-100"
                   icon="crop-simple"
                   name="Perspective"
                   small>
            <div v-for="sect in sections" :key="sect.key"
                 :class="`${state.section === sect.key?'router-link-active':''}`"
                 class="macro-icon-default">
              <SidebarItem :name="sect.name"
                           @click="selectSection(sect.key)"></SidebarItem>
            </div>
          </Sidebar>
          <!-- Wavelengths -->
          <Sidebar class="w-100"
                   icon="swatchbook"
                   name="Wavelength"
                   small>
            <div v-for="mode in viewModes" :key="mode.key" :class="`${state.mode === mode.key?'router-link-active':''}`"
                 class="macro-icon-default">
              <SidebarItem :name="mode.name"
                           @click="selectMode(mode.key)"></SidebarItem>
            </div>
          </Sidebar>
        </div>
        <div class="d-flex justify-content-between align-items-center flex-column flex-grow-1 w-100">
          <div class="d-flex justify-content-center align-items-center align-content-center h-100">
            <div v-if="state.section === 'FD'"
                 :style="`background-image: url('${buildURL(state.satellite, state.section, state.mode)}');`"
                 class="earth-full-disk p-5">
            </div>
            <div v-else :style="`background-image: url('${buildURL(state.satellite, state.section, state.mode)}');`"
                 class="preview p-5">
            </div>
          </div>
          <div class="element d-inline-block d-flex flex-row justify-content-start align-items-center gap p-2 mt-3">
            <img alt="noaa" class="noaa mx-1" src="/noaa.svg"/>

            <div class="flex-shrink-0">
              <div class="label-c2 label-w500 label-o5 lh-sm">Last Update</div>
              <div class="label-c2 label-w300 label-o4 lh-1">{{ state.lastUpdated }}</div>
            </div>
            <div class="flex-shrink-0">
              <div class="label-c2 label-w500 label-o5 lh-sm">Next Update</div>
              <div class="label-c2 label-w300 label-o4 lh-1">{{ state.nextUpdate }}</div>
            </div>
            <div class="d-flex gap-2 align-items-center justify-content-center flex-grow-1 px-2">

              <div class="d-flex gap">
                <div class="label-c2 label-w500 label-o5 lh-1">{{ state.satellite }}</div>
                <div class="label-c2 lh-1"><i class="fa-solid fa-caret-right"></i></div>
              </div>
              <div class="d-flex gap">
                <div class="label-c2 label-w500 label-o5 lh-1">{{ state.section }}</div>
                <div class="label-c2 lh-1"><i class="fa-solid fa-caret-right"></i></div>
              </div>
              <div class="d-flex gap">
                <div class="label-c2 label-w500 label-o5 lh-1">{{ state.mode }}</div>
              </div>
            </div>
            <div class="dock-icon lh-1 button-icon" @click="buildURL(state.satellite, state.section, state.mode)">
              <i class="fa-solid fa-arrow-rotate-right fa-fw"></i>
            </div>
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

.outline {
  border-radius: 100%;
}

.preview {
  width: 24rem !important;
  height: 24rem;
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