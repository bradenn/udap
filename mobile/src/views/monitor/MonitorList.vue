<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>


import {onMounted, reactive, watchEffect} from "vue";

import type {Entity} from "@/types"
import {formatByteSize} from "@/types";
import core from "@/core";
import Utilization from "@/components/Utilization.vue";
import ToggleView from "@/components/elements/ToggleView.vue";

interface Network {
  netout: number;
  netin: number;
}

interface Utilization {
  mem: number;
  maxmem: number;
  cpu: number;
  cpus?: number;
  maxcpu: number;
  disk: number;
  maxdisk: number;
}

interface VM {
  vmid: number;
  status: string;
  name: string;
  pid: number;
  uptime: number;
  utilization: Utilization;
  network: Network;
}

interface Node {
  node: string;
  status: string;
  uptime: number;
  utilization: Utilization;
  virtualMachines: VM[];
}

const remote = core.remote()

const state = reactive({
  proxmox: {} as Entity,
  data: [] as Node[],
  loaded: false
})

onMounted(() => {
  updateEntities()
})

watchEffect(() => {
  updateEntities()
  return remote.attributes
})

function updateEntities() {
  let entity = remote.entities.find(e => e.name == "nodes" && e.module == "proxmox")
  if (!entity) return
  state.proxmox = entity
  let attr = remote.attributes.find(a => a.key == "nodes" && a.entity == state.proxmox.id)
  if (!attr) return;
  let out = JSON.parse(attr.value)
  if (!out) return;
  state.data = out.nodes as Node[]
  state.loaded = true
}

</script>

<template>
  <div class="px-2">
    <!--    <div class="label-c2 label-w600">Monitoring</div>-->
    <div v-if="state.loaded" class="d-flex flex-column gap-1">
      <div v-for="node in state.data" class=" element px-2 py-2 d-flex flex-column gap-2">
        <div class="d-flex justify-content-between">
          <div class="label-c3 label-o6 label-w600 d-flex align-items-center gap-1 ">
            <div class="text-success label-w800 label-c1 lh-1">•</div>
            <div>
              {{
                node.node
              }}
            </div>
          </div>
          <div class="label-c3 label-o6 label-w600">
            <div v-if="node.status == 'online'" class="text-success">{{ node.status }}</div>
            <div v-else class="text-warning">{{ node.status }}</div>
          </div>
        </div>
        <div class="d-flex w-100 justify-content-between gap-1">
          <Utilization :max="100" :min="0" :value="node.utilization.cpu*100" class="flex-grow-1"
                       name="CPU">{{ Math.round(node.utilization.cpu * 1000) / 10 }}%
          </Utilization>
          <Utilization :max="node.utilization.maxdisk" :min="0" :value="node.utilization.disk" class="flex-grow-1"
                       name="SSD">
            {{ Math.round((node.utilization.disk / node.utilization.maxdisk) * 1000) / 10 }}%
          </Utilization>

          <Utilization :max="node.utilization.maxmem" :min="0" :value="node.utilization.mem" class="flex-grow-1"
                       name="Memory">
            {{ Math.round((node.utilization.mem / node.utilization.maxmem) * 1000) / 10 }}%
          </Utilization>
        </div>

        <!--        <div class="label-c6 label-o4 label-w600 px-1 d-flex mt-1">Virtual Machines</div>-->
        <div class="d-flex flex-column gap-1">
          <div v-for="vm in node.virtualMachines.sort((a, b) => a.vmid - b.vmid)" :key="vm.vmid">
            <ToggleView :subtitle="`pve ${vm.vmid}`" :title="vm.name" class="subplot px-1 py-2">

              <template v-slot:description>
                <div class="label-c5 label-o6 label-w600">
                  <div v-if="vm.status == 'running'" class="text-success">{{ vm.status }}</div>
                  <div v-else class="text-warning">{{ vm.status }}</div>
                </div>
              </template>
              <template v-slot:primary>
                <div class="d-flex w-100 justify-content-between gap-1 px-1">
                  <Utilization :max="100" :min="0" :value="vm.utilization.cpu*100" class="flex-grow-1 w-25"
                               name="CPU">{{ Math.round(vm.utilization.cpu * 1000) / 10 }}%
                  </Utilization>

                  <Utilization :max="vm.utilization.maxmem" :min="0" :value="vm.utilization.mem"
                               class="flex-grow-1 w-75"
                               name="Memory">
                    {{ formatByteSize(vm.utilization.mem) }} / {{ formatByteSize(vm.utilization.maxmem) }}
                  </Utilization>
                </div>
              </template>
              <template v-slot:secondary>
                <div class="d-flex w-100 gap-1 px-1">
                  <div class="label-c6 label-o6 label-w600 d-flex w-50"><span class="sf-icon label-c7 label-o3">􀆈</span>
                    {{ formatByteSize(vm.network.netin) }}
                  </div>
                  <div class="label-c6 label-o6 label-w600 d-flex w-50"><span class="sf-icon label-c7 label-o3">􀆇</span>
                    {{ formatByteSize(vm.network.netout) }}
                  </div>
                </div>
              </template>
            </ToggleView>
          </div>
        </div>
      </div>
    </div>

  </div>

</template>

<style lang="scss" scoped>

</style>