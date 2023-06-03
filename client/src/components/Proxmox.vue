<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {onMounted, reactive, watchEffect} from "vue";
import core from "@/core";
import moment from "moment";
import {formatByteSize} from "../types";
import Utilization from "@/components/Utilization.vue";

declare interface Network {
  netout: number;
  netin: number;
}

declare interface Utilization {
  mem: number;
  maxmem: number;
  cpu: number;
  cpus?: number;
  maxcpu: number;
  disk: number;
  maxdisk: number;
}

declare interface VM {
  vmid: number;
  status: string;
  name: string;
  pid: number;
  uptime: number;
  utilization: Utilization;
  network: Network;
}

declare interface Node {
  node: string;
  status: string;
  uptime: number;
  utilization: Utilization;
  virtualMachines: VM[];
}

const state = reactive({
  nodes: [] as Node[],
  tick: 0,
  updated: 0,
  loaded: false,
})

onMounted(() => {
  updateStats()
})


const remote = core.remote()

watchEffect(() => {
  updateStats()
  return remote.attributes
})

function updateStats() {
  state.tick = (state.tick + 10) % 100
  let attribute = remote.attributes.find(a => a.key == "nodes")
  if (!attribute) return

  state.nodes = JSON.parse(attribute.value).nodes as Node[]
  state.updated = Date.now().valueOf()
  state.loaded = true
}


</script>

<template>
  <div v-for="node in state.nodes" v-if="state.loaded" :key="node.node" class="element w-100 p-1">
    <div class="label-c2 d-flex justify-content-between px-1">
      <div class="label-c2 label-w600 label-o4">{{ node.node }}</div>
      <div class="label-c2 label-w600 label-o4">{{ moment(state.updated).fromNow() }}
      </div>
    </div>
    <div class="gap-1 d-flex flex-column mt-1">
      <div v-for="vm in node.virtualMachines.sort((a, b) => a.vmid - b.vmid)" :key="vm.vmid"
           class="subplot py-1 d-flex justify-content-between flex-column align-items-start"
           style="padding-left: 1rem">
        <div class="d-flex justify-content-between w-100">
          <div class="label-c2 label-w600 label-o4">{{ vm.name }}</div>
          <div class="label-c2 label-w600 label-o4">{{ vm.status }},
            {{ moment().subtract(node.uptime * 1000).fromNow(true) }}
          </div>
        </div>

        <!--        <Utilization :percent="vm.utilization.cpu"></Utilization>-->

        <div class="label-o3">{{ formatByteSize(vm.utilization.mem) }} / {{
            formatByteSize(vm.utilization.maxmem)
          }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>