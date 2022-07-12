<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Allocation} from "@/components/AllocationBar.vue";
import AllocationBar from "@/components/AllocationBar.vue";
import Header from "@/components/Header.vue";
import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";
import type {Remote} from "@/types";

let s1: Allocation[] = [
  {
    name: 'bmni',
    source: "sq2",
    allocated: 0.325,
    active: 0.21,
    battery: true
  },
  {
    name: 'nuc',
    source: "sq2",
    allocated: 1.5,
    active: 1.17,
    battery: true
  }, {
    name: 'ap',
    source: "sq2",
    allocated: 0.15,
    active: 0.14,
    battery: true
  },
]

let s2: Allocation[] = [

  {
    name: 'bmbp',
    source: "sq1",
    allocated: 2,
    active: 1.56
  },
  {
    name: 'mon1',
    source: "sq1",
    allocated: 1.8,
    active: 1.1
  },
  {
    name: 'mon2',
    source: "sq1",
    allocated: 1.2,
    active: 0.8
  },
  {
    name: 'hp1',
    source: "sq1",
    allocated: 0.5,
    active: 0.45
  }, {
    name: 'hp2',
    source: "sq1",
    allocated: 0.5,
    active: 0.45
  }, {
    name: 'tv',
    source: "sq1",
    allocated: 1.21,
    active: 0.11
  },
  {
    name: 'at',
    source: "sq1",
    allocated: 1,
    active: 0.12
  }
]

let ups: Allocation[] = [
  {
    name: 'bmni',
    source: "sq2",
    allocated: 0.325,
    active: 0.22
  },
  {
    name: 'nuc',
    source: "sq2",
    allocated: 1.5,
    active: 1.34
  }, {
    name: 'ap',
    source: "sq2",
    allocated: 0.15,
    active: 0.14
  },
]

let grid: Allocation[] = []


let state = reactive({
  side1: s1,
  side2: s2,
  grid: grid,
  isFixed: true,
  attributes: [] as Allocation[]
})


let remote = inject("remote") as Remote

onMounted(() => {

  updateDynamic()
})


watchEffect(() => {
  updateDynamic()
  return remote.attributes
})


function updateDynamic() {
  let entities = remote.entities.filter(e => e.type === "dimmer" || e.type === "spectrum" || e.type === "toggle")
  for (let i = 0; i < entities.length; i++) {
    let attrs = remote.attributes.find(a => a.key === 'dim' && a.entity == entities[i].id)
    if (!attrs) return
    state.grid = state.grid.filter(g => g.name !== entities[i].name)
    state.grid.push({
      name: entities[i].name,
      source: "grid",
      allocated: 9 / 120,
      active: (9 / 120) * (parseInt(attrs.value, 10) / 100),
    })
  }
}

function setFixed(isFixed: boolean) {
  state.isFixed = isFixed
}

</script>

<template>

  <div class="d-flex flex-column gap-1 justify-content-start">
    <div class="d-flex justify-content-between">
      <Header icon="diagram-project" name="Allocations" title="Allocations"></Header>
      <Plot :cols="2" :rows="1" style="width: 13rem">
        <Subplot :active="state.isFixed" :fn="() => setFixed(true)" name="Fixed"></Subplot>
        <Subplot :active="!state.isFixed" :fn="() => setFixed(false)" name="Dynamic"></Subplot>
      </Plot>
    </div>

    <AllocationBar :allocatable="15" :allocations="state.side1" :is-fixed="state.isFixed"
                   name="Squid 1-8"></AllocationBar>
    <AllocationBar :allocatable="15" :allocations="state.side2" :is-fixed="state.isFixed"
                   name="Squid 9-16"></AllocationBar>
    <AllocationBar :allocatable="0.5" :allocations="state.grid" :is-fixed="state.isFixed"
                   name="Lighting"></AllocationBar>
  </div>
</template>

<style lang="scss" scoped>
$baseFontSize: 16;
$green: #1abc9c;
$yellow: #f1c40f;
$red: #c0392b;
$blue: #3498db;
$grey: #f2f2f2;

@function rem($val) {
  @return #{calc($val / $baseFontSize)}rem;
}

// Gauge
.mask {
  position: relative;
  overflow: hidden;

  display: block;
  width: rem(200);
  height: rem(100);
  margin: rem(20);
}

.semi-circle {
  position: relative;

  display: block;
  width: rem(200);
  height: rem(100);

  background: linear-gradient(to right, $red 0%, $yellow 50%, $green 100%);

  border-radius: 50% 50% 50% 50% / 100% 100% 0% 0%;

  &::before {
    content: "";

    position: absolute;
    bottom: 0;
    left: 50%;
    z-index: 2;

    display: block;
    width: rem(140);
    height: rem(70);
    margin-left: rem(-70);

    background: #fff;

    border-radius: 50% 50% 50% 50% / 100% 100% 0% 0%;
  }
}

.semi-circle--mask {
  position: absolute;
  top: 0;
  left: 0;

  width: rem(200);
  height: rem(200);

  background: transparent;

  transform: rotate(120deg) translate3d(0, 0, 0);
  transform-origin: center center;
  backface-visibility: hidden;
  transition: all .3s ease-in-out;

  &::before {
    content: "";

    position: absolute;
    top: 0;
    left: 0%;
    z-index: 2;

    display: block;
    width: rem(202);
    height: rem(102);
    margin: -1px 0 0 -1px;

    background: #f2f2f2;

    border-radius: 50% 50% 50% 50% / 100% 100% 0% 0%;
  }
}


.gauge--1 {
  .semi-circle {
    background: $green;
  }
}
</style>
