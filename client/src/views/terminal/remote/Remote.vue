<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity as EntityType, Preferences} from "@/types";
import Slider from "@/components/Slider.vue";
import Entity from "@/components/Entity.vue";
import Switch from "@/components/Switch.vue";
import attributeService from "@/services/attributeService";
import moment from "moment";
import type {Remote} from "@/remote";
import FixedScroll from "@/components/scroll/FixedScroll.vue";
import Color from "@/components/Color.vue";
import core from "@/core";

const state = reactive({
    entities: [] as EntityType[],
    selected: [] as string[],

    target: {} as EntityType,
    attributes: [] as Attribute[],
})


const remote = inject("remote") as Remote
const prefs = inject("preferences") as Preferences

function sortByModule(a: EntityType, b: EntityType): number {
    if (a.name >= b.name) {
        return 1
    } else if (a.name < b.name) {
        return -1
    }
    return 0

}

function compareOrder(a: Attribute, b: Attribute): number {
    if (a.order >= b.order) {
        return 1
    } else if (a.order < b.order) {
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
    // state.selected = id
    if (state.selected.includes(id)) {
        state.selected = state.selected.filter(s => s != id)
    } else {
        state.selected.push(id)
    }
    fetchAttributes()
    if (state.selected.length == 0) {
        state.attributes = []
    }

}

function groupBy<T>(xs: T[], key: string): T[] {
    return xs.reduce(function (rv: any, x: any): T {
        (rv[x[key]] = x);
        return rv;
    }, {});
}

function fetchAttributes() {
    if (state.selected.length > 0) {
        let a = remote.attributes.filter(a => state.selected.includes(a.entity))
        if (!a) return
        let grouped = groupBy(a, "key")

        state.attributes = Object.values(grouped) as Attribute[]
        state.attributes = state.attributes.sort(compareOrder)

    }
}


const notify = core.notify()

function attributeRequestAll(key: string, value: string) {
    let attrs = remote.attributes.filter(a => state.selected.includes(a.entity) && a.key === key) as Attribute[]
    if (!attrs) return;

    for (let attr of attrs) {
        attr.request = value
        attributeService.request(attr).then(res => {
        }).catch(err => {
            notify.fail("Core", `Could not change '${attr.key}'to '${attr.request}'`, err)
        })
    }
    notify.success("Core", `Changed ${key} to '${value}'`)


}


function applyState(attribute: Attribute) {
    attributeRequestAll(attribute.key, attribute.value)
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

function hexToRgb(hex: string) {
    let result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16)
    } : {r: 0, g: 0, b: 0};
}

function RGBToHSL(r: number, g: number, b: number): any {
    // Make r, g, and b fractions of 1
    r /= 255;
    g /= 255;
    b /= 255;

    // Find greatest and smallest channel values
    let cmin = Math.min(r, g, b),
        cmax = Math.max(r, g, b),
        delta = cmax - cmin,
        h = 0,
        s = 0,
        l = 0;
    // Calculate hue
    // No difference
    if (delta == 0)
        h = 0;
    // Red is max
    else if (cmax == r)
        h = ((g - b) / delta) % 6;
    // Green is max
    else if (cmax == g)
        h = (b - r) / delta + 2;
    // Blue is max
    else
        h = (r - g) / delta + 4;

    h = Math.round(h * 60);

    // Make negative hues positive behind 360°
    if (h < 0)
        h += 360;
    // Calculate lightness
    l = (cmax + cmin) / 2;

    // Calculate saturation
    s = delta == 0 ? 0 : delta / (1 - Math.abs(2 * l - 1));

    // Multiply l and s by 100
    s = +(s * 100).toFixed(1);
    l = +(l * 100).toFixed(1);


    return h
}

function hueFromHex(hex: string): number {
    let hexToRgb1 = hexToRgb(hex);
    return RGBToHSL(hexToRgb1.r, hexToRgb1.g, hexToRgb1.b)

}

function timeSince(time: string): string {
    return moment(time).fromNow()
}

</script>

<template>
    <div class="remote-grid h-100">
        <FixedScroll class="h-100 overflow-scroll" style="grid-column: 1 / span 3;">
            <div class="entity-grid h-100">
                <Entity v-for="entity in state.entities" :key="entity.id" :entity="entity"
                        :selected="state.selected.includes(entity.id)" @mousedown="() => selectEntity(entity.id)">

                </Entity>
            </div>
        </FixedScroll>
        <div class="d-flex flex-column gap-1" style="grid-column: 4 / span 6;">
            <div v-for="attribute in state.attributes" :key="`${attribute.id}`">
                <Switch v-if="attribute.key === 'on'"
                        :active="attribute.value === 'true'"
                        :change="(val) => attributeRequestAll(attribute.key, `${val?'true':'false'}`)"
                        name="Power"></Switch>
                <Slider v-else-if="attribute.key === 'dim'"
                        :change="(val) => {attributeRequestAll(attribute.key, `${val}`)}" :live="false"
                        :max="100"
                        :min="0"
                        :step="5"
                        :value="parseInt(attribute.request)" name="Brightness" unit="%"></Slider>
                <div v-else-if="attribute.key === 'cct'" class="">
                    <Slider :change="(val) => {attributeRequestAll(attribute.key, `${colorTemps[val].value}`)}"
                            :live="false"
                            :max="colorTemps.length-1" :min="0" :step="1"
                            :tags="colorTemps.map(c => c.name)"
                            :value="0" name="Color Temperature" unit="K"></Slider>

                </div>
                <div v-else-if="attribute.key === 'hue'">
                    <div class="element d-flex flex-column" style="display: flex; height: 4rem;">
                        <div class="label-w500 label-r label-c1 label-o3 px-1">
                            Color
                        </div>
                        <FixedScroll :horizontal="true" style="overflow-x: scroll !important;">
                            <div class="color-grid">
                                <div v-for="color in prefs.appdata.colors">
                                    <Color :color="color"
                                           @click="(e) => attributeRequestAll(attribute.key, `${hueFromHex(color)}`)"></Color>
                                </div>
                            </div>
                        </FixedScroll>
                    </div>

                </div>
                <div v-else-if="attribute.key === 'api'">

                </div>
                <div v-else style="display: none;"></div>
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

.color-grid {
    align-items: start;
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    padding: 0.25rem;
    gap: 0.25rem;
    width: 100%;

}

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