<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Device, Endpoint} from "@/types";
import Loader from "@/components/Loader.vue";
import endpointService from "@/services/endpointService";
import Toolbar from "@/components/toolbar/Toolbar.vue";
import ToolbarButton from "@/components/ToolbarButton.vue";
import CreateEndpoint from "@/views/terminal/settings/endpoint/CreateEndpoint.vue";
import type {Remote} from "@/remote";
import EndpointKey from "@/components/EndpointKey.vue";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
    endpoints: {} as Endpoint[],
    devices: {} as Device[],
    loading: true,
    toCreate: {
        name: ""
    } as Endpoint,
    mode: "list"
})


onMounted(() => {
    state.loading = true
    handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
    state.endpoints = remote.endpoints
    state.devices = remote.devices
    state.loading = false
    return remote
}

function setMode(mode: string) {
    state.mode = mode
}

function createEndpoint() {
    endpointService.createEndpoint(state.toCreate).then(res => {
        state.toCreate.name = `${res}`
    }).catch(err => {
        state.toCreate.name = `${err}`
    })
}

</script>

<template>
    <div v-if="!state.loading">

        <Toolbar class="mb-1" icon="􀏞" title="Endpoints">
            <div class="w-100"></div>
            <ToolbarButton :active="false" icon="􀅼" style="height: 1.5rem"
                           text="Endpoint"
                           @click="(e) => setMode(state.mode === 'create'?'list':'create')"></ToolbarButton>
        </Toolbar>

        <div v-if="state.mode === 'list'" class="d-flex flex-column gap-1 pt-1">
            <div class="label-c2 label-w700 label-o4 px-2">Active</div>
            <div class="endpoint-container w-100">
                <div v-for="endpoint in state.endpoints.filter(e => e.connected)"
                     :key="endpoint.id" class="">
                    <div class="element d-flex justify-content-between flex-column"
                         style="height: 3.25rem">
                        <div
                                class="d-flex gap-1 align-items-center justify-content-between">
                            <div class="d-flex gap-1 align-items-center p-1 py-0">
                                <div class="label-xs label-o3">􀏃</div>
                                <div class="label-xxs label-w500 label-o4">{{ endpoint.name }}</div>
                            </div>

                            <div v-if="endpoint.connected"
                                 class="label-xxs label-o5 text-accent subplot p-1 px-2 d-flex align-items-center gap-1 mx-1">
                                <div class="label-c2 lh-1">Connected</div>
                            </div>
                            <div v-else
                                 class="label-xxs label-o5 text-muted subplot p-1 px-2 d-flex align-items-center gap-1 mx-1">
                                <div class="label-c2 lh-1">Disconnected</div>
                            </div>
                        </div>
                        <EndpointKey :value="endpoint.key"></EndpointKey>
                        <!--                        -->

                    </div>

                </div>
            </div>

            <div class="label-c2 label-w700 label-o4 px-2">Inactive</div>
            <div class="endpoint-container w-100">
                <div v-for="endpoint in state.endpoints.filter(e => !e.connected)"
                     :key="endpoint.id" class="">
                    <div class="element d-flex justify-content-between flex-column align-content-start"
                         style="height: 3.25rem">
                        <div :key="endpoint.id"
                             class="d-flex gap-1 align-items-center justify-content-between">
                            <div class="d-flex gap-1 align-items-center p-1 py-0">
                                <div class="label-xs label-o3">􀏃</div>
                                <div class="label-xxs label-w500 label-o4">{{ endpoint.name }}</div>
                            </div>
                            <div v-if="endpoint.connected"
                                 class="label-xxs label-o5 text-accent subplot p-1 px-2 d-flex align-items-center gap-1 mx-1">
                                <div class="label-c2 lh-1">Connected</div>
                            </div>
                            <div v-else
                                 class="label-xxs label-o5 text-muted subplot p-1 px-2 d-flex align-items-center gap-1 mx-1">
                                <div class="label-c2 lh-1">Disconnected</div>
                            </div>


                        </div>
                        <EndpointKey :value="endpoint.key"></EndpointKey>
                    </div>

                </div>

            </div>
        </div>
        <div v-else-if="state.mode === 'create'">
            <CreateEndpoint :done="() => {state.mode = 'list'}"></CreateEndpoint>
        </div>
    </div>

    <div v-else>
        <div class="d-flex justify-content-start ali py-2 px-1">
            <div class="label-w500 label-o4 label-xxl"><i
                    :class="`fa-solid fa-expand fa-fw`"></i></div>
            <div class="label-w500 opacity-100 label-xxl px-2">Endpoints</div>
            <div class="flex-fill"></div>

        </div>
        <div class="element p-2">
            <div class="label-c1 label-o4 d-flex align-content-center gap-1">
                <div>
                    <Loader size="sm"></Loader>
                </div>
                <div class="">Loading...</div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.endpoint-container {
    width: 100%;
    display: grid;
    grid-column-gap: 0.25rem;
    grid-row-gap: 0.25rem;

    grid-template-columns: repeat(4, 1fr);
}

</style>