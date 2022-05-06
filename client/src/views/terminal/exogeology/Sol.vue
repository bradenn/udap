<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import Sidebar from "@/components/sidebar/Sidebar.vue";
import SidebarItem from "@/components/sidebar/SidebarItem.vue";
import {reactive} from "vue";

let state = reactive({
  angstroms: '171',
  view: 'disk'
})

const views = [
  {
    name: "Disk",
    key: "disk"
  },
  {
    name: "Square",
    key: "square"
  },
]

const wavelengths = [
  {
    angstroms: '094',
  },
  {
    angstroms: '131',
  },
  {
    angstroms: '171',
  },
  {
    angstroms: '195',
  },
  {
    angstroms: '284',
  },
  {
    angstroms: '304',
  }
]

function selectView(view: string) {
  state.view = view
}

function selectWavelength(angstroms: string) {
  state.angstroms = angstroms
}

function buildUrl(angstroms: string) {
  return `https://services.swpc.noaa.gov/images/animations/suvi/primary/${angstroms}/latest.png`
}

</script>

<template>
  <div>
    <div class="d-flex justify-content-start p-2 px-1">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-sun fa-fw`"></i></div>
      <div class="label-w500  label-xxl px-2">Sol</div>
    </div>
    <div class="">
      <div class=" d-flex flex-row gap-4">
        <div class="d-flex flex-column gap w-25">
          <!-- Perspective -->
          <Sidebar class="w-100" icon="crop-simple" name="Perspective" small>
            <div v-for="mode in views"
                 :key="mode.key"
                 :class="`${state.view === mode.key?'router-link-active':''}`">
              <SidebarItem :name="`${mode.name}`"
                           @click="selectView(mode.key)"></SidebarItem>
            </div>
          </Sidebar>
          <!-- Wavelengths -->
          <Sidebar class="w-100" icon="swatchbook" name="Wavelength" small>
            <div v-for="mode in wavelengths" :key="mode.angstroms"
                 :class="`${state.angstroms === mode.angstroms?'router-link-active':''}`">
              <SidebarItem :alt="`${parseFloat(mode.angstroms)/10.0}nm`"
                           :name="`${parseFloat(mode.angstroms)} Angstroms`"
                           @click="selectWavelength(mode.angstroms)"></SidebarItem>
            </div>
          </Sidebar>
        </div>

        <div class="d-flex justify-content-center align-items-center w-75">
          <div :class="`${state.view === 'disk'?'sol-disk':'sol-square'}`"
               :style="`background-image: url('${buildUrl(state.angstroms)}');`"
               class=" p-3">
          </div>

        </div>
      </div>
    </div>
  </div>

</template>

<style scoped>
.sol-disk {
  width: 22rem !important;

  box-shadow: 0 0 4px 8px rgba(0, 0, 0, 0.025);
  aspect-ratio: 1/1;
  border-radius: 100%;
  background-color: rgba(76, 87, 101, 0.37);
  background-size: 163%;
  background-position: center;
  transition: all 0.2s ease-in-out;
}

.sol-disk:hover {
  transform: scale(1.2);
}

.sol-square {
  width: 22rem !important;

  box-shadow: 0 0 4px 8px rgba(0, 0, 0, 0.025);
  aspect-ratio: 1/1;
  border-radius: 1rem;
  background-color: rgba(76, 87, 101, 0.37);
  background-position: center;
  background-size: cover;
  transition: all 0.2s ease-in-out;
}

.sol-square:hover {
  transform: scale(1.2);
}
</style>