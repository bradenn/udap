<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import {inject, onMounted, reactive, watchEffect} from "vue";
import Subplot from "@/components/plot/Subplot.vue";


import type {Entity} from "@/types";
import type {Remote} from "@/remote";


const voices = [
    {
        name: "Indian Male 1",
        score: 1,
        path: "aup",
    },
    {
        name: "Indian Female",
        score: 1,
        path: "axb",
    },
    {
        name: "Indian Male 2",
        score: 2,
        path: "gka",
    },
    {
        name: "Robot Male",
        score: 3,
        path: "jmk",
    },
    {
        name: "Indian Male 3",
        score: 3,
        path: "ksp",
    },
    {
        name: "Androgynous Robot",
        score: 1,
        path: "rxr",
    },
    {
        name: "Robotic Female 1",
        score: 1,
        path: "eey",
    },

]

const primaryVoices = [
    {
        name: "Default (British Male)",
        score: 4,
        path: "default"
    },
    {
        name: "US Female 1",
        score: 3,
        path: "slt",
    },
    {
        name: "Robot British Male",
        score: 2,
        path: "awb",
    },
    {
        name: "US Male 1",
        score: 2,
        path: "aew",
    },
    {
        name: "Ambiguous Male",
        score: 2,
        path: "ahw",
    },
    {
        name: "US Male 2",
        score: 3,
        path: "bdl",
    },
    {
        name: "US Female 2",
        score: 2,
        path: "clb",
    },
    {
        name: "Profound US Male",
        score: 3,
        path: "fem",
    },
    {
        name: "Robot Female",
        score: 3,
        path: "ljm",
    },
    {
        name: "Old Robot Male",
        score: 3,
        path: "rms",
    },

]


let state = reactive({
    voice: "default",
    entity: {} as Entity
})


let remote = inject("remote") as Remote

onMounted(() => {
    getAtlasEntity(remote)
})

watchEffect(() => {
    getAtlasEntity(remote)
    return remote.entities
})

function getAtlasEntity(remote: Remote) {
    let target = remote.entities.find(e => e.name === "atlas" && e.module === "atlas")
    if (!target) return
    state.entity = target
}

function setVoice(path: string) {
    state.voice = path
    remote.nexus.requestAttribute(state.entity.id, "voice", path)
}

</script>

<template>
    <div class="d-flex flex-column">
        <div class="label-lg mb-1">Voice Selection</div>
        <div class="d-flex gap">
            <div>
                <Plot :cols="1" :rows="8" alt="Clarity&nbsp;&nbsp;" style="width: 13rem" title="Primary Voices">
                    <Subplot v-for="voice in primaryVoices.sort((a, b) => b.score - a.score)"
                             :active="state.voice === voice.path"
                             :alt="`${Array(voice.score+1).join('•')}`"
                             :fn="() => setVoice(voice.path)"
                             :name="voice.name"></Subplot>
                </Plot>
            </div>
            <div>
                <Plot :cols="1" :rows="8" alt="Clarity&nbsp;&nbsp;" style="width: 13rem" title="Esoteric Voices">
                    <Subplot v-for="voice in voices.sort((a, b) => b.score - a.score)"
                             :active="state.voice === voice.path"
                             :alt="`${Array(voice.score+1).join('•')}`"
                             :fn="() => setVoice(voice.path)"
                             :name="voice.name"></Subplot>
                </Plot>
            </div>
        </div>


    </div>
</template>

<style lang="scss" scoped>
</style>