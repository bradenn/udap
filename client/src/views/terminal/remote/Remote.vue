<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity as EntityType} from "@/types";
import Slider from "@/components/Slider.vue";
import Entity from "@/components/Entity.vue";
import Switch from "@/components/Switch.vue";
import attributeService from "@/services/attributeService";
import moment from "moment";
import type {Remote} from "@/remote";
import FixedScroll from "@/components/scroll/FixedScroll.vue";

const state = reactive({
    entities: [] as EntityType[],
    selected: "",

    target: {} as EntityType,
    attributes: [] as Attribute[],
})


const remote = inject("remote") as Remote

function sortByModule(a: EntityType, b: EntityType): number {
    if (a.name >= b.name) {
        return 1
    } else if (a.name < b.name) {
        return -1
    }
    return 0

}

onMounted(() => {
    updateData()
})

watchEffect(() => {
    updateData()
    return remote.entities
})

function updateData() {
    state.entities = remote.entities.filter(e => e.module === 'govee' || e.module === 'macmeta').sort(sortByModule)
}

function selectEntity(id: string) {
    state.selected = id
    // if (state.selected.includes(id)) {
    //     state.selected = state.selected.filter(s => s != id)
    // } else {
    //     state.selected.push(id)
    // }
    fetchAttributes()

}

function fetchAttributes() {
    if (state.selected.length > 0) {
        let a = remote.attributes.filter(a => a.entity === state.selected)
        if (!a) return
        state.attributes = a
    }
}


function attributeRequestAll(key: string, value: string) {
    let attrs = remote.attributes.filter(a => state.selected.includes(a.entity) && a.key == key)
    if (!attrs) return;

    for (let i = 0; i < attrs.length; i++) {
        let attr = attrs[i];
        attr.request = value;
        attributeService.request(attr)

    }


}


function applyState(attribute: Attribute) {
    attributeService.request(attribute)
}

const colorTemps = [
    {
        name: "Warm",
        value: 2400,
        color: "255, 157, 63"
    },
    {
        name: "Tungsten",
        value: 3000,
        color: "255, 180, 107"
    },
    {
        name: "Moonlight",
        value: 4100,
        color: "255, 211, 168"
    }, {
        name: "Daylight",
        value: 5000,
        color: "255, 236, 224"
    },
    {
        name: "Cloudy",
        value: 5500,
        color: "255, 236, 224"
    },
    {
        name: "White",
        value: 6500,
        color: "255, 249, 253"
    }
]

function timeSince(time: string): string {
    return moment(time).fromNow()
}

</script>

<template>
    <div class="remote-grid h-100">
        <FixedScroll class="h-100 overflow-scroll" style="grid-column: 1 / span 3;">
            <div class="entity-grid h-100">
                <Entity v-for="entity in state.entities" :key="entity.id" :entity="entity"
                        :selected="state.selected === entity.id" @mousedown="() => selectEntity(entity.id)">

                </Entity>
            </div>
        </FixedScroll>
        <div class="d-flex flex-column" style="grid-column: 4 / span 6;">
            <div v-for="attribute in state.attributes" :key="`${attribute.id}`">
                <Switch v-if="attribute.key === 'on'"
                        :active="attribute.value === 'true'"
                        :change="(val) => {attribute.request = `${val?'true':'false'}`; applyState(attribute)}"
                        name="Power"></Switch>
                <Slider v-else-if="attribute.key === 'dim'"
                        :change="(val) => {attribute.request = `${val}`; applyState(attribute)}" :live="false"
                        :max="100"
                        :min="0"
                        :step="5"
                        :value="parseInt(attribute.request)" name="Brightness" unit="%"></Slider>
                <div v-else-if="attribute.key === 'cct'" class="">
                    <Slider :change="(val) => {attribute.request = `${colorTemps[val].value}`; applyState(attribute)}"
                            :live="false"
                            :max="colorTemps.length-1" :min="0" :step="1"
                            :tags="colorTemps.map(c => c.name)"
                            :value="0" name="Color Temperature" unit="K"></Slider>

                </div>
                <div v-else-if="attribute.key === 'api'">

                </div>
            </div>
        </div>
        <div class="d-flex flex-column gap-1" style="grid-column: 10 / span 3;">
            <div class="element">
                <div class="d-flex label-o4 label-c2 justify-content-between">
                    <div class="label-w600 label-r label-c2">
                        Last Updated
                    </div>
                    <div class="label-w500 label-o3 label-c2">
                        {{ timeSince(state.target.updated) }}
                    </div>
                </div>
                <div class="d-flex label-o4 label-c2 justify-content-between">
                    <div class="label-w600 label-r label-c2">
                        Created
                    </div>
                    <div class="label-w500 label-o3 label-c2">
                        {{ timeSince(state.target.created) }}
                    </div>
                </div>
                <div class="d-flex label-o4 label-c2 justify-content-between">
                    <div class="label-w600 label-r label-c2">
                        Neural
                    </div>
                    <div class="label-w500 label-o3 label-c2">
                        {{ state.target.neural }}
                    </div>
                </div>
                <div class="d-flex label-o4 label-c2 justify-content-between">
                    <div class="label-w600 label-r label-c2">
                        Alias
                    </div>
                    <div class="label-w500 label-o3 label-c2">
                        {{ state.target.alias }}
                    </div>
                </div>
            </div>
            <div class="element d-flex d-flex justify-content-between">
                <div class="label-w500 label-c1 label-o4 px-1">Neural</div>
                <div class="label-w500 label-c1 label-o3 px-1">{{ state.target.neural }}</div>
            </div>
        </div>

    </div>
</template>

<style scoped>

.remote-controls {
    grid-column: 4 / span 6;
}

.entity-grid {

    display: grid;
    grid-column-gap: 0.25rem;
    grid-row-gap: 0.25rem;
    grid-template-rows: repeat(10, 1fr);
    grid-template-columns: repeat(3, 1fr);
}

.remote-grid {
    display: grid;
    grid-column-gap: 0.25rem;
    grid-row-gap: 0.25rem;
    grid-template-rows: repeat(1, 1fr);
    grid-template-columns: repeat(12, 1fr);
}
</style>