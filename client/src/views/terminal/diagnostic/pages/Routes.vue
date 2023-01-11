<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import type {Remote} from "@/remote";
import core from "@/core";

const remote: Remote = core.remote()

</script>

<template>
  <div>

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