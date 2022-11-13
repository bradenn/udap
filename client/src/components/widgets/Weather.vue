<script lang="ts" setup>

import {inject, onMounted, reactive, watchEffect} from "vue";
import moment from "moment";
import type {Attribute, Entity, Remote} from "@/types";
import type {CurrentWeather, Weather} from "@/weather"
import {getWeatherIcon, getWeatherState} from "@/weather"

function degToCompass(num: number) {
    const val = Math.floor((num / 45) + 0.5);
    const arr = ["North", "North East", "East", "South East", "South", "South West", "West", "North West"];
    return arr[(val % 8)]
}


interface StateType {
    entity: Entity
    forecast: Attribute
    interval: number
    loading: boolean
    gradient: number
    weather: Weather
    ranges: any
    menu: boolean
    current: CurrentWeather
    timeSince: string
    lastPull: number
    properties: any
    now: any
}

const state = reactive<StateType>({
    entity: {} as Entity,
    forecast: {} as Attribute,
    interval: 0,
    loading: true,
    menu: false,
    gradient: 0,
    weather: {} as Weather,
    current: {} as CurrentWeather,
    lastPull: 0,
    timeSince: "",
    properties: {
        updateTime: "",
        generatedAt: "",
    },
    ranges: {
        temp: {
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
    now: {}
})


onMounted(() => {
    state.loading = true
    handleUpdates(remote)
})


let remote: Remote = inject("remote") || {} as Remote

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
    if (!remote) return
    let entity = remote.entities.find(e => e.name === 'weather')
    if (!entity) return;
    state.entity = entity
    let attribute = remote.attributes.find(e => e.entity === (entity as Entity).id && e.key === 'forecast')
    if (!attribute) return;
    state.forecast = attribute
    parseWeather(JSON.parse(attribute.value) as Weather)
    return remote.attributes
}

function parseWeather(we: Weather) {

    state.properties.updateTime = moment(we.generationtime_ms).fromNow()

    state.properties.generatedAt = moment(we.generationtime_ms).fromNow()
    state.current = we.current_weather

    if (we.hourly.temperature_2m.length <= 0) return
    for (let i = 0; i < Math.min(we.hourly.temperature_2m.length, 24); i++) {
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

    }


    state.weather = we
    state.lastPull = new Date(we.generationtime_ms).valueOf()
    state.timeSince = moment(state.lastPull).fromNow()
    state.loading = false
}

function toggleMenu() {
    state.menu = !state.menu
}

</script>

<template>
    <div v-if="state.loading" class="element gap w-100 h-100 p-2">

        <div class="d-flex justify-content-between align-items-center">
            <div class="label-xl label-o6 label-w500 lh-1 ">
                --°
            </div>
            <div class="label-c2 label-w400 label-o4 lh-1 d-flex flex-column align-items-end" style="gap: 0.125rem">
                <div>

                    <span>--</span>
                </div>
                <div class="label-c3 label-o2">H: --° L:
                                               --°
                </div>
            </div>
        </div>

        <h6 class="mb-1 mt-1">Today</h6>
        <div class=" d-flex flex-row justify-content-between px-1">
            <div class="d-flex flex-row justify-content-between align-items-end w-100">
                <div v-for="v in Array(7).keys()"
                     class="d-flex flex-column justify-content-center align-items-center">
                    <div class="label-c4 label-w400 label-o4">{{
                            moment(new Date().setHours(new Date().getHours() + v)).format("hA")
                        }}
                    </div>
                    <div class="label-c2 label-o4">
                        --
                    </div>
                    <div class="label-c4 label-o3 label-w500">&nbsp;--°</div>

                </div>
            </div>
        </div>

        <div class="label-c2 label-o3 py-1 pb-0">
        </div>
        <div class="h-sep my-0"></div>
        <div class="label-c2 label-w400 label-o2 lh-lg">
            <div>Winds at -- mph going --
            </div>
        </div>

    </div>
    <div v-else-if="state.menu" class="element gap w-100 h-100 p-2 weather-widget" @click="state.menu = false">
        <div class="d-flex justify-content-between align-items-center">
            <div class="label-xl label-o6 label-w500 lh-1 ">{{ Math.round(state.current.temperature) }}°</div>
            <div class="label-c2 label-w400 label-o4 lh-1 d-flex flex-column align-items-end" style="gap: 0.125rem">
                <div>
                    <div v-if="getWeatherState(state.current.weathercode).toLowerCase().includes('rain')">
                        Tut, Tut; Looks like {{ getWeatherState(state.current.weathercode) }}
                    </div>
                    <div v-else>
                        {{ getWeatherState(state.current.weathercode) }}
                    </div>
                </div>
                <div class="label-c3 label-o2">H: {{ Math.round(state.ranges.temp.max) }}° L:
                    {{ Math.round(state.ranges.temp.min) }}°
                </div>
            </div>
        </div>

        <h6 class="mb-0 mt-1">This Week</h6>
        <div class=" d-flex flex-row justify-content-between px-1">
            <div class="d-flex flex-row justify-content-between align-items-end w-100">
                <div v-for="v in Array(7).keys()"
                     class="d-flex flex-column justify-content-center align-items-center">
                    <div class="label-c4 label-w400 label-o4">{{
                            moment(new Date()).add(v, "day").format("ddd")
                        }}
                    </div>


                    <div class="label-c2 label-o4">{{ getWeatherIcon(state.weather.daily.weathercode[v], 24) }}</div>
                    <div v-if="state.weather.daily.precipitation_sum[v] === 0" class="label-c4 label-o3 label-w500">
                        &nbsp;{{ Math.round(state.weather.daily.temperature_2m_max[v]) }}°
                    </div>
                    <div v-else class="label-c4 label-o3 label-w500" style="color: rgba(92,177,246,0.6)">{{
                            state.weather.daily.precipitation_sum[v]
                        }}"
                    </div>

                </div>
            </div>
        </div>
        <div v-show="false">
            <div class="h-sep my-1"></div>
            <h6 class="mb-1 mt-2">Rain</h6>
            <div class=" d-flex flex-row justify-content-between px-1">
                <div class="d-flex flex-row justify-content-between align-items-end w-100">
                    <div v-for="p in state.weather.hourly.precipitation.slice(0, 24)">
                        <div :style="`height: ${(p - state.ranges.rain.min)/(state.ranges.rain.max - state.ranges.rain.min)}rem;`"
                             class="temp-bar"></div>
                    </div>
                </div>
                <div class="d-flex flex-column justify-content-start align-items-end px-1">
                    <div class="label-subtext">{{ state.ranges.rain.max }}in</div>
                    <div class="label-subtext mt-1">{{ state.ranges.rain.min }}in</div>
                </div>
            </div>
        </div>
        <div class="label-c2 label-o3 py-1 pb-0">
        </div>
        <div class="h-sep my-0"></div>
        <div class="label-c2 label-w400 label-o2 lh-lg">
            <div>Winds at {{ state.current.windspeed }} mph going {{
                    degToCompass(state.current.winddirection)
                }}
            </div>
        </div>

    </div>
    <div v-else class="element gap w-100 h-100 p-1 px-2 weather-widget" @click="toggleMenu">
        <div class="d-flex justify-content-start align-items-center">
            <div class="label-c2 label-w400 label-o4 lh-1 d-flex flex-column align-items-start "
                 style="gap: 0.125rem; width: 4rem">
                <div class="label-cn6 label-o5 label-w600 pt-1 pb-0 lh-1 label-r">{{
                        Math.round(state.current.temperature)
                    }}°
                </div>
                <div v-if="state.current">
                    <div v-if="getWeatherState(state.current.weathercode || 0).toLowerCase().includes('rain')">
                        Tut, Tut; Looks like {{ getWeatherState(state.current.weathercode) }}
                    </div>
                    <div v-else>
                        {{ getWeatherState(state.current.weathercode) }}
                    </div>
                </div>
                <div class="label-c3 label-o2 label-w500">H: {{ Math.round(state.ranges.temp.max) }}° L:
                    {{ Math.round(state.ranges.temp.min) }}°
                </div>
            </div>
            <div class="flex-grow-1 pt-1">
                <div class=" d-flex flex-row justify-content-between px-1">
                    <div class="d-flex flex-row justify-content-between align-items-end w-100">
                        <div v-for="v in Array(6).keys()"
                             class="d-flex flex-column justify-content-center align-items-center">
                            <div class="label-c4 label-w400 label-o4">{{
                                    moment(new Date()).add(v, "hour").format("hA")
                                }}
                            </div>


                            <div class="label-c2 label-o4">
                                {{ getWeatherIcon(state.weather.hourly.weathercode[new Date().getHours() + v], v) }}
                            </div>
                            <div v-if="state.weather.hourly.precipitation[new Date().getHours()+v] <= 0.0001"
                                 class="label-c4 label-o3 label-w500">
                                &nbsp;{{ Math.round(state.weather.hourly.temperature_2m[new Date().getHours() + v]) }}°
                            </div>
                            <div v-else class="label-c4 label-o3 label-w500" style="color: rgba(92,177,246,0.6)">{{
                                    state.weather.hourly.precipitation[new Date().getHours() + v]
                                }}"
                            </div>

                        </div>
                    </div>
                </div>

            </div>
        </div>


        <div class="label-c2 label-o3 py-1 pb-0">
        </div>
        <div class="label-c2 label-w400 label-o2 lh-lg">
            <!--      <div>Winds at {{ state.current.windspeed }} mph going {{-->
            <!--    <div class="h-sep my-0"></div>-->
            <!--          degToCompass(state.current.winddirection)-->
            <!--        }}-->
            <!--      </div>-->
        </div>

    </div>


</template>

<style lang="scss" scoped>

.contain {
  display: inline-block;
  width: 6rem;
  max-width: 6rem;
  height: 1rem;
  overflow: clip !important;
  text-overflow: ellipsis !important;
  text-wrap: none !important;
}

.night {
  background-color: rgba(255, 255, 255, 0.2) !important;
}

.weather-widget:active {
  animation: click 100ms ease forwards;
}


@keyframes click {
  0% {
    transform: scale(1);
  }
  15% {
    transform: scale(0.97);
  }
  25% {
    transform: scale(0.98);
  }
  100% {
    transform: scale(0.97);
  }
}

.label-subtext {
  font-size: 0.45rem;
  width: 8px;
  opacity: 0.5;
  line-height: 0.4rem;
}

.border-box {
  border: 1px solid white;
}

.temp-bar {
  width: 8px;
  border-radius: 3px;
  background-color: rgba(255, 255, 255, 0.5);
}

.widget {


}

.sky-gradient {
  height: 100%;
  border-radius: 0.4rem;
  transition: all 1s ease-in;
}

.g0 {
  background: #012459;
  background: -moz-linear-gradient(top, #012459 0%, #001322 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #012459), color-stop(100%, #001322));
  background: -webkit-linear-gradient(top, #012459 0%, #001322 100%);
  background: -o-linear-gradient(top, #012459 0%, #001322 100%);
  background: -ms-linear-gradient(top, #012459 0%, #001322 100%);
  background: linear-gradient(to bottom, #012459 0%, #001322 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#012459', endColorstr='#001322', GradientType=0);

}

.g1 {
  background: #003972;
  background: -moz-linear-gradient(top, #003972 0%, #001322 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #003972), color-stop(100%, #001322));
  background: -webkit-linear-gradient(top, #003972 0%, #001322 100%);
  background: -o-linear-gradient(top, #003972 0%, #001322 100%);
  background: -ms-linear-gradient(top, #003972 0%, #001322 100%);
  background: linear-gradient(to bottom, #003972 0%, #001322 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#003972', endColorstr='#001322', GradientType=0);
}

.g2 {
  background: #003972;
  background: -moz-linear-gradient(top, #003972 0%, #001322 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #003972), color-stop(100%, #001322));
  background: -webkit-linear-gradient(top, #003972 0%, #001322 100%);
  background: -o-linear-gradient(top, #003972 0%, #001322 100%);
  background: -ms-linear-gradient(top, #003972 0%, #001322 100%);
  background: linear-gradient(to bottom, #003972 0%, #001322 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#003972', endColorstr='#001322', GradientType=0);
}

.g3 {
  background: #004372;
  background: -moz-linear-gradient(top, #004372 0%, #00182b 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #004372), color-stop(100%, #00182b));
  background: -webkit-linear-gradient(top, #004372 0%, #00182b 100%);
  background: -o-linear-gradient(top, #004372 0%, #00182b 100%);
  background: -ms-linear-gradient(top, #004372 0%, #00182b 100%);
  background: linear-gradient(to bottom, #004372 0%, #00182b 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#004372', endColorstr='#00182b', GradientType=0);

}

.g4 {
  background: #004372;
  background: -moz-linear-gradient(top, #004372 0%, #011d34 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #004372), color-stop(100%, #011d34));
  background: -webkit-linear-gradient(top, #004372 0%, #011d34 100%);
  background: -o-linear-gradient(top, #004372 0%, #011d34 100%);
  background: -ms-linear-gradient(top, #004372 0%, #011d34 100%);
  background: linear-gradient(to bottom, #004372 0%, #011d34 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#004372', endColorstr='#011d34', GradientType=0);

}

.g5 {
  background: #016792;
  background: -moz-linear-gradient(top, #016792 0%, #00182b 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #016792), color-stop(100%, #00182b));
  background: -webkit-linear-gradient(top, #016792 0%, #00182b 100%);
  background: -o-linear-gradient(top, #016792 0%, #00182b 100%);
  background: -ms-linear-gradient(top, #016792 0%, #00182b 100%);
  background: linear-gradient(to bottom, #016792 0%, #00182b 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#016792', endColorstr='#00182b', GradientType=0);

}

.g6 {
  background: #07729f;
  background: -moz-linear-gradient(top, #07729f 0%, #042c47 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #07729f), color-stop(100%, #042c47));
  background: -webkit-linear-gradient(top, #07729f 0%, #042c47 100%);
  background: -o-linear-gradient(top, #07729f 0%, #042c47 100%);
  background: -ms-linear-gradient(top, #07729f 0%, #042c47 100%);
  background: linear-gradient(to bottom, #07729f 0%, #042c47 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#07729f', endColorstr='#042c47', GradientType=0);

}

.g7 {
  background: #12a1c0;
  background: -moz-linear-gradient(top, #12a1c0 0%, #07506e 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #12a1c0), color-stop(100%, #07506e));
  background: -webkit-linear-gradient(top, #12a1c0 0%, #07506e 100%);
  background: -o-linear-gradient(top, #12a1c0 0%, #07506e 100%);
  background: -ms-linear-gradient(top, #12a1c0 0%, #07506e 100%);
  background: linear-gradient(to bottom, #12a1c0 0%, #07506e 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#12a1c0', endColorstr='#07506e', GradientType=0);

}

.g8 {
  background: #74d4cc;
  background: -moz-linear-gradient(top, #74d4cc 0%, #1386a6 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #74d4cc), color-stop(100%, #1386a6));
  background: -webkit-linear-gradient(top, #74d4cc 0%, #1386a6 100%);
  background: -o-linear-gradient(top, #74d4cc 0%, #1386a6 100%);
  background: -ms-linear-gradient(top, #74d4cc 0%, #1386a6 100%);
  background: linear-gradient(to bottom, #74d4cc 0%, #1386a6 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#74d4cc', endColorstr='#1386a6', GradientType=0);

}

.g9 {
  background: #efeebc;
  background: -moz-linear-gradient(top, #efeebc 0%, #61d0cf 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #efeebc), color-stop(100%, #61d0cf));
  background: -webkit-linear-gradient(top, #efeebc 0%, #61d0cf 100%);
  background: -o-linear-gradient(top, #efeebc 0%, #61d0cf 100%);
  background: -ms-linear-gradient(top, #efeebc 0%, #61d0cf 100%);
  background: linear-gradient(to bottom, #efeebc 0%, #61d0cf 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#efeebc', endColorstr='#61d0cf', GradientType=0);

}

.g10 {
  background: #fee154;
  background: -moz-linear-gradient(top, #fee154 0%, #a3dec6 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #fee154), color-stop(100%, #a3dec6));
  background: -webkit-linear-gradient(top, #fee154 0%, #a3dec6 100%);
  background: -o-linear-gradient(top, #fee154 0%, #a3dec6 100%);
  background: -ms-linear-gradient(top, #fee154 0%, #a3dec6 100%);
  background: linear-gradient(to bottom, #fee154 0%, #a3dec6 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#fee154', endColorstr='#a3dec6', GradientType=0);

}

.g11 {
  background: #fdc352;
  background: -moz-linear-gradient(top, #fdc352 0%, #e8ed92 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #fdc352), color-stop(100%, #e8ed92));
  background: -webkit-linear-gradient(top, #fdc352 0%, #e8ed92 100%);
  background: -o-linear-gradient(top, #fdc352 0%, #e8ed92 100%);
  background: -ms-linear-gradient(top, #fdc352 0%, #e8ed92 100%);
  background: linear-gradient(to bottom, #fdc352 0%, #e8ed92 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#fdc352', endColorstr='#e8ed92', GradientType=0);

}

.g12 {
  background: #ffac6f;
  background: -moz-linear-gradient(top, #ffac6f 0%, #ffe467 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #ffac6f), color-stop(100%, #ffe467));
  background: -webkit-linear-gradient(top, #ffac6f 0%, #ffe467 100%);
  background: -o-linear-gradient(top, #ffac6f 0%, #ffe467 100%);
  background: -ms-linear-gradient(top, #ffac6f 0%, #ffe467 100%);
  background: linear-gradient(to bottom, #ffac6f 0%, #ffe467 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#ffac6f', endColorstr='#ffe467', GradientType=0);

}

.rain {
  box-shadow: inset 0 0 12px 0.5px opacify(rgba(44, 44, 44, 0.6), 0.00125) !important;
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, rgba(54, 104, 169, 0.5)), color-stop(40%, rgba(41, 92, 147, 0.5)), color-stop(100%, rgba(61, 77, 105, 0.5)));

}


.overcast {
  box-shadow: inset 0 0 12px 0.5px opacify(rgba(44, 44, 44, 0.6), 0.00125) !important;
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, rgba(144, 144, 144, 0.6)), color-stop(20%, rgba(144, 144, 144, 0.6)), color-stop(100%, rgba(77, 77, 77, 0.4)));

}

.g14 {
  background: #fd9e58;
  background: -moz-linear-gradient(top, #fd9e58 0%, #ffe467 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #fd9e58), color-stop(100%, #ffe467));
  background: -webkit-linear-gradient(top, #fd9e58 0%, #ffe467 100%);
  background: -o-linear-gradient(top, #fd9e58 0%, #ffe467 100%);
  background: -ms-linear-gradient(top, #fd9e58 0%, #ffe467 100%);
  background: linear-gradient(to bottom, #fd9e58 0%, #ffe467 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#fd9e58', endColorstr='#ffe467', GradientType=0);

}

.g15 {
  background: #f18448;
  background: -moz-linear-gradient(top, #f18448 0%, #ffd364 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #f18448), color-stop(100%, #ffd364));
  background: -webkit-linear-gradient(top, #f18448 0%, #ffd364 100%);
  background: -o-linear-gradient(top, #f18448 0%, #ffd364 100%);
  background: -ms-linear-gradient(top, #f18448 0%, #ffd364 100%);
  background: linear-gradient(to bottom, #f18448 0%, #ffd364 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#f18448', endColorstr='#ffd364', GradientType=0);

}

.g16 {
  background: #f06b7e;
  background: -moz-linear-gradient(top, #f06b7e 0%, #f9a856 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #f06b7e), color-stop(100%, #f9a856));
  background: -webkit-linear-gradient(top, #f06b7e 0%, #f9a856 100%);
  background: -o-linear-gradient(top, #f06b7e 0%, #f9a856 100%);
  background: -ms-linear-gradient(top, #f06b7e 0%, #f9a856 100%);
  background: linear-gradient(to bottom, #f06b7e 0%, #f9a856 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#f06b7e', endColorstr='#f9a856', GradientType=0);

}

.g17 {
  background: #ca5a92;
  background: -moz-linear-gradient(top, #ca5a92 0%, #f4896b 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #ca5a92), color-stop(100%, #f4896b));
  background: -webkit-linear-gradient(top, #ca5a92 0%, #f4896b 100%);
  background: -o-linear-gradient(top, #ca5a92 0%, #f4896b 100%);
  background: -ms-linear-gradient(top, #ca5a92 0%, #f4896b 100%);
  background: linear-gradient(to bottom, #ca5a92 0%, #f4896b 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#ca5a92', endColorstr='#f4896b', GradientType=0);

}

.g18 {
  background: #5b2c83;
  background: -moz-linear-gradient(top, #5b2c83 0%, #d1628b 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #5b2c83), color-stop(100%, #d1628b));
  background: -webkit-linear-gradient(top, #5b2c83 0%, #d1628b 100%);
  background: -o-linear-gradient(top, #5b2c83 0%, #d1628b 100%);
  background: -ms-linear-gradient(top, #5b2c83 0%, #d1628b 100%);
  background: linear-gradient(to bottom, #5b2c83 0%, #d1628b 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#5b2c83', endColorstr='#d1628b', GradientType=0);

}

.g19 {
  background: #371a79;
  background: -moz-linear-gradient(top, #371a79 0%, #713684 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #371a79), color-stop(100%, #713684));
  background: -webkit-linear-gradient(top, #371a79 0%, #713684 100%);
  background: -o-linear-gradient(top, #371a79 0%, #713684 100%);
  background: -ms-linear-gradient(top, #371a79 0%, #713684 100%);
  background: linear-gradient(to bottom, #371a79 0%, #713684 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#371a79', endColorstr='#713684', GradientType=0);

}

.g20 {
  background: #28166b;
  background: -moz-linear-gradient(top, #28166b 0%, #45217c 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #28166b), color-stop(100%, #45217c));
  background: -webkit-linear-gradient(top, #28166b 0%, #45217c 100%);
  background: -o-linear-gradient(top, #28166b 0%, #45217c 100%);
  background: -ms-linear-gradient(top, #28166b 0%, #45217c 100%);
  background: linear-gradient(to bottom, #28166b 0%, #45217c 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#28166b', endColorstr='#45217c', GradientType=0);

}

.g21 {
  background: #192861;
  background: -moz-linear-gradient(top, #192861 0%, #372074 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #192861), color-stop(100%, #372074));
  background: -webkit-linear-gradient(top, #192861 0%, #372074 100%);
  background: -o-linear-gradient(top, #192861 0%, #372074 100%);
  background: -ms-linear-gradient(top, #192861 0%, #372074 100%);
  background: linear-gradient(to bottom, #192861 0%, #372074 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#192861', endColorstr='#372074', GradientType=0);

}

.g22 {
  background: #040b3c;
  background: -moz-linear-gradient(top, #040b3c 0%, #233072 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #040b3c), color-stop(100%, #233072));
  background: -webkit-linear-gradient(top, #040b3c 0%, #233072 100%);
  background: -o-linear-gradient(top, #040b3c 0%, #233072 100%);
  background: -ms-linear-gradient(top, #040b3c 0%, #233072 100%);
  background: linear-gradient(to bottom, #040b3c 0%, #233072 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#040b3c', endColorstr='#233072', GradientType=0);

}

.g23 {
  background: #040b3c;
  background: -moz-linear-gradient(top, #040b3c 0%, #012459 100%);
  background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, #040b3c), color-stop(100%, #012459));
  background: -webkit-linear-gradient(top, #040b3c 0%, #012459 100%);
  background: -o-linear-gradient(top, #040b3c 0%, #012459 100%);
  background: -ms-linear-gradient(top, #040b3c 0%, #012459 100%);
  background: linear-gradient(to bottom, #040b3c 0%, #012459 100%);
  filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#040b3c', endColorstr='#012459', GradientType=0);

}

</style>
