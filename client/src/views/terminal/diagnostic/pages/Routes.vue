<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import type {Remote} from "@/remote";
import core from "@/core";

const remote: Remote = core.remote()

</script>

<template>
  <div class="">
    <div v-for="user in remote.users">
      {{ user.username }}
      {{ user.first }}
      {{ user.middle }}
      {{ user.last }}
      {{ user.classification }}
      {{ user.residency }}
    </div>
  </div>
  <div class="element d-flex">
    <div class="label-c1 px-1">
      􀝝
    </div>
    <div class="nav-bar">
      <div class="subplot d-flex align-items-center">
        <div>
          <div class="lh-1 label-r label-w600">
            Timer
          </div>
          <div class="lh-1 label-o4 label-c4 label-w400">
            12:03
          </div>
        </div>
        <div class="label-c3 label-o2 lh-1">
          􀆄
        </div>
      </div>
      <div class="subplot">

      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.nav-bar {
  height: 2rem;
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  grid-template-rows: repeat(1, 1fr);
}
</style>