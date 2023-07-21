// Copyright (c) 2023 Braden Nicholson

import {Attribute, Entity} from "../types";
import {CurrentWeather, getWeatherIcon, Weather} from "../weather";
import {inject, onMounted, reactive, watchEffect} from "vue";
//@ts-ignore
import moment from 'moment';

import {Remote} from "../remote";

export interface Forecast {
    entity: Entity
    forecast: Attribute
    interval: number
    loading: boolean
    gradient: number
    weather: Weather
    ranges: any
    hourly: HourlyData[]
    menu: boolean
    current: CurrentWeather
    timeSince: string
    lastPull: number
    properties: any
    now: any
}

export interface HourlyData {
    time: string
    icon: string
    temp: number
    rain: number
}

export function useForecast(): Forecast {


    const state = reactive<Forecast>({
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
        hourly: [] as HourlyData[],
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
        let now = moment(Date.now()).hour()

        let current = state.weather.hourly.time.find(e => now == moment(e * 1000).hour()) || 0
        let id = state.weather.hourly.time.indexOf(current)
        state.hourly = []
        let supl: HourlyData[] = [] as HourlyData[]
        for (let i = id; i < id + 6; i++) {
            supl.push({
                time: moment(state.weather.hourly.time[i] * 1000).format("h A") as string,
                temp: state.weather.hourly.temperature_2m[i],
                rain: state.weather.hourly.precipitation[i],
                icon: getWeatherIcon(state.weather.hourly.weathercode[i], i) || ""
            } as HourlyData)

        }
        state.hourly = supl
        state.lastPull = new Date(we.generationtime_ms).valueOf()
        state.timeSince = moment(state.lastPull).fromNow()
        state.loading = false
    }

    return state
}

