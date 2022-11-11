<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
// Copyright (c) 2022 Braden Nicholson
import moment from "moment";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Remote} from "@/types";
import type {CurrentWeather, Weather} from "@/weather";
import {getWeatherIcon, getWeatherState} from "@/weather"
import PaneList from "@/components/pane/PaneList.vue";
import PaneListItemInline from "@/components/pane/PaneListItemInline.vue";
import Scroll from "@/components/scroll/Scroll.vue";
import HorizontalChart from "@/components/charts/HorizontalChart.vue";

interface WeatherProps {
    current: CurrentWeather
    latest: Weather
    entity: Entity
    forecast: Attribute
    loading: boolean
    ranges: any,
    sun: any,
    rain: any,
    lastUpdate: number,
    canvas: HTMLCanvasElement
}

let state = reactive<WeatherProps>({
    canvas: {} as HTMLCanvasElement,
    current: {} as CurrentWeather,
    latest: {} as Weather,
    entity: {} as Entity,
    forecast: {} as Attribute,
    loading: false,
    lastUpdate: 0 as number,
    sun: {
        setting: "" as string,
        rising: "" as string,
    },
    rain: {
        rainfall: 0 as number
    },
    ranges: {
        temp: {
            min: 100,
            max: -100,
        },
        humidity: {
            min: 100,
            max: -100,
        },
        wind: {
            min: 100,
            max: 0,
        }, rain: {
            min: 0,
            max: 0,
        },
    },
});


onMounted(() => {
    state.loading = true
    handleUpdates(remote)

})

let remote = inject("remote") as Remote

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
    if (!remote) return
    let entity = remote.entities.find(e => e.name === 'weather')
    if (!entity) return;
    state.entity = entity as Entity
    let attribute = remote.attributes.find(e => e.entity === state.entity.id && e.key === 'forecast')
    if (!attribute) return;
    state.forecast = attribute
    parseWeather(JSON.parse(attribute.value) as Weather)
    return remote.attributes
}

function parseWeather(we: Weather) {

    state.current = we.current_weather
    if (!we) return;
    if (we.hourly.temperature_2m.length <= 0) return
    for (let i = 0; i < we.hourly.temperature_2m.length; i++) {
        if (state.ranges.temp.max < we.hourly.temperature_2m[i]) {
            state.ranges.temp.max = we.hourly.temperature_2m[i]
        } else if (state.ranges.temp.min > we.hourly.temperature_2m[i]) {
            state.ranges.temp.min = we.hourly.temperature_2m[i]
        }

        if (state.ranges.rain.max < we.hourly.precipitation[i]) {
            state.ranges.rain.max = we.hourly.precipitation[i]
        } else if (state.ranges.rain.min > we.hourly.precipitation[i]) {
            state.ranges.rain.min = we.hourly.precipitation[i]
        }

        if (state.ranges.humidity.max < we.hourly.relativehumidity_2m[i]) {
            state.ranges.humidity.max = we.hourly.relativehumidity_2m[i]
        } else if (state.ranges.humidity.min > we.hourly.relativehumidity_2m[i]) {
            state.ranges.humidity.min = we.hourly.relativehumidity_2m[i]
        }

    }
    state.latest = we as Weather

    state.sun.rising = asDate(state.latest.daily.sunrise[0] * 1000)
    state.sun.setting = asDate(state.latest.daily.sunset[0] * 1000)

    state.rain.rainfall = state.latest.daily.precipitation_sum[0]


    state.lastUpdate = new Date().valueOf();
    state.loading = false
    loadCanvas()
}

function roundDecimal(input: number, places: number) {
    return Math.round(input * Math.pow(10, places)) / Math.pow(10, places)
}

function timeSince(ms: number): string {
    return moment(ms).fromNow()
}

function asDate(ms: number): string {
    return moment(ms).utc(true).format("h:mm A")
}

let renderer = {} as THREE.WebGLRenderer
let camera = {} as THREE.OrthographicCamera
let scene = {} as THREE.Scene


function setCamera(x: number, y: number, z: number) {
    camera.position.set(x, y, z);

}


function render() {
    renderer.setClearColor(0x000000, 0);

    renderer.render(scene, camera);
}

function animate() {

    render()
}


function loadCanvas() {
    state.canvas = document.getElementById("weather-chart") as HTMLCanvasElement
    const ctx = state.canvas.getContext('2d')
    if (!ctx) return

    const scale = 1

    ctx.scale(scale, scale)
    state.canvas.width = state.canvas.clientWidth / scale
    state.canvas.height = state.canvas.clientHeight / scale

    drawCanvas();
}

function lerp(v0: number, v1: number, t: number): number {
    return (1 - t) * v0 + t * v1;
}

function trendLine(ctx: CanvasRenderingContext2D, points: number[]) {

}

function drawWeek(days: number) {

}

function drawCanvas() {

    const ctx = state.canvas.getContext('2d')

    if (!ctx) return
    ctx.clearRect(0, 0, state.canvas.width, state.canvas.height)


    ctx.strokeStyle = `rgba(255,200,64,${0.5})`
    ctx.fillStyle = `rgba(255,255,255,${0.6})`
    ctx.lineWidth = 3
    ctx.lineCap = "round"
    ctx.lineJoin = "round"
    ctx.beginPath()
    // ctx.moveTo(-2, state.canvas.height / 2)

    let div = state.latest.hourly.temperature_2m.length;
    let chunk = state.canvas.width / div;
    let avg = (state.ranges.temp.max + state.ranges.temp.min) / 2
    for (let i = 0; i < state.canvas.width; i++) {
        let index = Math.floor(i / chunk)
        let floor = state.latest.hourly.temperature_2m[index]
        let next = state.latest.hourly.temperature_2m[index + 1]
        let dist = i % chunk;
        let interp = lerp(floor, next, dist / chunk);

        ctx.lineTo(i, (state.canvas.height / 2) - (interp - avg) * 4)
    }
    ctx.moveTo(state.canvas.width, (state.canvas.height / 2))
    ctx.closePath()
    ctx.stroke()


    ctx.beginPath()
    ctx.lineWidth = 2
    ctx.strokeStyle = `rgba(255,255,255,${0.1})`
    for (let i = 0; i < state.canvas.width; i++) {
        let index = Math.floor(i / chunk)
        let time = state.latest.hourly.time[index] * 1000
        if (new Date(time).getHours() === 0) {
            let fmt = moment(time).format("dddd")
            ctx.font = "20px SF Pro Display"

            ctx.fillText(fmt, i + 10, 20)
            if (index !== 0) {
                ctx.moveTo(i, 0)
                ctx.lineTo(i, (state.canvas.height))
            }
            ctx.closePath()
            i += chunk;
        }
        if (new Date(time).getDate() == new Date().getDate() && new Date(time).getHours() == new Date().getHours()) {
            ctx.font = "20px SF Pro Display"
            ctx.lineWidth = 5
            ctx.moveTo(i, 60)
            ctx.lineTo(i, (state.canvas.height) - 60)
            ctx.closePath()
            i += chunk;
        }
        ctx.lineWidth = 2
    }
    ctx.stroke()

    ctx.beginPath()
    ctx.lineWidth = 2
    ctx.strokeStyle = `rgba(255,255,255,${0.1})`
    let nowX = new Date().getTime();

    ctx.closePath()
    ctx.stroke()

    ctx.lineWidth = 2
    ctx.lineDashOffset = 10
    ctx.strokeStyle = `rgba(255,255,255,${0.1})`
    ctx.setLineDash([2, 8])
    ctx.beginPath()
    let maxY = (state.canvas.height / 2) - (state.ranges.temp.max - avg) * 4;
    ctx.fillText(state.ranges.temp.max + " F", 10, maxY - 10)
    ctx.moveTo(0, maxY)
    ctx.lineTo(state.canvas.width, maxY)
    ctx.closePath()

    ctx.stroke()

    ctx.beginPath()
    let minY = (state.canvas.height / 2) - (state.ranges.temp.min - avg) * 4;
    ctx.fillText(state.ranges.temp.min + " F", 10, minY + 25)
    ctx.moveTo(0, minY)
    ctx.lineTo(state.canvas.width, minY)
    ctx.closePath()

    ctx.stroke()
    ctx.setLineDash([0, 0])
    ctx.beginPath()
    // ctx.moveTo(-2, state.canvas.height / 2)
    div = state.latest.hourly.precipitation.length;
    chunk = state.canvas.width / div;
    avg = (state.ranges.rain.max + state.ranges.rain.min) / 2
    ctx.strokeStyle = `rgba(128,128,255,${0.5})`
    for (let i = 0; i < state.canvas.width; i++) {
        let index = Math.floor(i / chunk)
        let floor = state.latest.hourly.precipitation[index]
        let next = state.latest.hourly.precipitation[index + 1]
        let dist = i % chunk;
        let interp = lerp(floor, next, dist / chunk);

        ctx.lineTo(i, minY - interp * (state.canvas.height) * 2)
    }
    ctx.moveTo(state.canvas.width, (state.canvas.height / 2))
    ctx.closePath()
    ctx.stroke()

    ctx.beginPath()
    // ctx.moveTo(-2, state.canvas.height / 2)
    div = state.latest.hourly.relativehumidity_2m.length;
    chunk = state.canvas.width / div;
    ctx.strokeStyle = `rgba(128,255,255,${0.5})`
    for (let i = 0; i < state.canvas.width; i++) {
        let index = Math.floor(i / chunk)
        let floor = state.latest.hourly.relativehumidity_2m[index]
        let next = state.latest.hourly.relativehumidity_2m[index + 1]
        let dist = i % chunk;
        let interp = lerp(floor, next, dist / chunk);

        ctx.lineTo(i, minY - interp)
    }
    ctx.moveTo(state.canvas.width, (state.canvas.height / 2))
    ctx.closePath()
    ctx.stroke()


}

function getWeekday(ms: number): string {
    return moment(ms * 1000).format("dddd")
}

</script>
<template>
    <div class="d-flex flex-column gap-1">
        <div v-if="!state.loading" class="element p-2 pt-1">
            <div class=" d-flex flex-row align-items-center">
                <div class="flex-shrink-1 " style="min-width: 9rem; padding-left: 0.25rem">
                    <h2 class="lh-md">{{ roundDecimal(state.latest.current_weather?.temperature, 0) }}° F</h2>
                    <div class="label-c1 label-r label-w400 label-o4 lh-1">{{
                            getWeatherState(state.latest.current_weather?.weathercode)
                        }}
                    </div>
                    <div class="label-c1 label-r label-o3 label-w400 ">High {{ Math.round(state.ranges.temp.max) }}° •
                        Low
                        {{ Math.round(state.ranges.temp.min) }}°
                    </div>
                </div>
                <div v-if="state.latest.hourly" class="d-flex flex-row justify-content-between align-items-end"
                     style="width: 100%">
                    <div
                            v-for="(hour) in Array(12).keys()">
                        <div v-if="state.latest.hourly?.temperature_2m[hour]"
                             :class="new Date().getHours()===hour?'':''"
                             class=" d-flex flex-column align-items-center justify-content-center px-3">
                            <div class="label-c3 label-w400 label-o3 mt-1">
                                {{ moment(new Date().setHours(new Date().getHours() + hour)).format("hA") }}
                            </div>
                            <div class="label-sm label-o4">
                                {{
                                    getWeatherIcon(state.latest.hourly.weathercode[new Date().getHours() + hour], hour)
                                }}
                            </div>
                            <div v-if="state.latest.hourly.precipitation[new Date().getHours()+hour] <= 0.0001"
                                 class="d-flex align-items-center justify-content-center">
                                <div class="label-c3 label-w500 label-o4 mt-1">
                                    {{ state.latest.hourly.temperature_2m[new Date().getHours() + hour] }}
                                </div>

                                <div class="label-c3 label-w400 label-o3 mt-1">
                                    {{ state.latest.hourly_units.temperature_2m }}
                                </div>

                            </div>
                            <div v-else class="d-flex align-items-center justify-content-center">
                                <div class="label-c3 label-w500 label-o4 mt-1 rain">
                                    {{ state.latest.hourly.precipitation[new Date().getHours() + hour] }}"
                                </div>
                            </div>
                        </div>


                    </div>

                </div>
            </div>
        </div>
        <div v-if="!state.loading" class="d-flex gap-1 flex-row">
            <PaneList :alt='timeSince(state.lastUpdate)' style="width: 13rem !important; height: 14rem !important;"
                      title="Today">
                <PaneListItemInline :subtext="state.sun.rising" icon="􀆱"
                                    title="Sunrise"></PaneListItemInline>
                <PaneListItemInline :subtext="state.sun.setting" icon="􀆳"
                                    title="Sunset"></PaneListItemInline>
                <PaneListItemInline :subtext="`${state.rain.rainfall} in`" icon="􀇆"
                                    title="Rainfall"></PaneListItemInline>
            </PaneList>
            <div class=" d-flex flex-column gap-1 w-100">
                <div class="element d-flex flex-column p-2 pt-1">
                    <div class="d-flex justify-content-between lh-1">
                        <div class="label-o4 label-w500 label-c1 py-1">This Week</div>
                        <div v-if="state.latest.hourly.precipitation[new Date().getHours()] >= 0.0001"
                             class="label-o3 label-w500 label-r label-c2 py-1">Tut, Tut, Looks Like Rain</div>
                    </div>
                    <Scroll horizontal style="overflow-x: scroll; height: 100%">
                        <div class="d-flex flex-column" style="width: 200%">
                            <HorizontalChart :marker="new Date().getHours()" :scale="2"
                                             :sections-names="state.latest.hourly.time.filter(t => (t*1000)%(24*60*60*1000) === 0).map(t => getWeekday(t))"
                                             :values="state.latest.hourly.temperature_2m"
                                             :valuesPerSection="24"
                                             color="rgba(255,159,10,0.4)"
                                             name="Temperature"
                                             unit="° F">

                            </HorizontalChart>
                            <HorizontalChart :marker="new Date().getHours()" :scale="2"
                                             :sections-names="state.latest.hourly.time.filter(t => (t*1000)%(24*60*60*1000) === 0).map(t => '')"
                                             :values="state.latest.hourly.precipitation"
                                             :valuesPerSection="24"
                                             color="rgba(10,132,255,0.4)"
                                             name="Rain"
                                             unit='\"'>

                            </HorizontalChart>

                            <HorizontalChart :marker="new Date().getHours()" :scale="2"
                                             :sections-names="state.latest.hourly.time.filter(t => (t*1000)%(24*60*60*1000) === 0).map(t => '')"
                                             :values="state.latest.hourly.relativehumidity_2m"
                                             :valuesPerSection="24"
                                             color="rgba(100,210,255,0.4)"
                                             name="Humidity"
                                             unit="%">

                            </HorizontalChart>
                        </div>
                    </Scroll>
                </div>
                <div>

                    <!--          <Scroll :horizontal="true" style=" overflow-x: scroll; max-width: 100%;" class="">-->
                    <!--            <canvas id="weather-chart"-->
                    <!--                    style=" min-width:400px; width: 150%; height: 300px; "></canvas>-->
                    <!--          </Scroll>-->
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>

.condition {
    width: 4rem;
    height: 4rem;
    align-items: center;
    justify-content: center;
    display: flex;
    font-size: 1.25rem;
}

.temp-chart {
    display: flex;
    align-items: end;
    gap: 2px;
}

.temp-bar {
    content: ' ';
    width: 6px;
    border-radius: 2px;
    background-color: rgba(255, 255, 255, 0.5);
}
</style>