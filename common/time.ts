// Copyright (c) 2023 Braden Nicholson


export interface TimeUnit {
    value: number
    units: string
}

function pluralize(value: number, unit: string, suffix: boolean): string {
    let suffixString = ""
    if (suffix) {
        suffixString = " ago"
    }
    if (value == 1) {
        return `${value} ${unit}${suffixString}`
    }
    return `${value} ${unit}s${suffixString}`
}

export function since(date: string, suffix: boolean = true): string {
    let event = new Date(date).valueOf()
    let now = Date.now()
    let value = now - event

    let miliseconds = value % 1000
    let seconds = Math.round(value / 1000 % 60)
    let minutes = Math.round((value / 1000 / 60) % 60)
    let hours = Math.round((value / 1000 / 60 / 60) % 60)
    let days = Math.round((value / 1000 / 60 / 60 / 24) % 24)
    let years = Math.round((value / 1000 / 60 / 60 / 24 / 365))

    if (years >= 1) {
        return `${pluralize(years, "year", suffix)}`
    } else if (days >= 1) {
        return `${pluralize(days, "day", suffix)}`
    } else if (hours >= 1) {
        return `${pluralize(hours, "hour", suffix)}`
    } else if (minutes >= 1) {
        return `${pluralize(minutes, "minute", suffix)}`
    } else if (seconds >= 1) {
        return `${pluralize(seconds, "second", suffix)}`
    } else if (miliseconds >= 1) {
        return `${pluralize(miliseconds, "millisecond", suffix)}`
    }

    return `${value}`
}

export function fromDuration(nanoseconds: number, isNano: boolean = false, isPrecise: boolean = false): TimeUnit[] {
    if (nanoseconds < 0) {
        nanoseconds = Math.abs(nanoseconds)
        // throw new Error("Input must be a non-negative number.");
    }
    if (!isNano) {
        let value = nanoseconds;
        let miliseconds = value % 1000
        let seconds = value / 1000 % 60
        let minutes = (value / 1000 / 60) % 60
        let hours = (value / 1000 / 60 / 60) % 60
        let days = Math.round((value / 1000 / 60 / 60 / 24) % 24)
        let years = Math.round((value / 1000 / 60 / 60 / 24 / 365))
        let timeUnits = [] as TimeUnit[]
        if (years > 0) {
            timeUnits.push({
                value: Math.round(years),
                units: "y"
            } as TimeUnit)
        }
        if (years > 0) {
            timeUnits.push({
                value: Math.round(days),
                units: "d"
            } as TimeUnit)
        }
        if (years > 0) {
            timeUnits.push({
                value: Math.round(hours),
                units: "h"
            } as TimeUnit)
        }
        if (years > 0) {
            timeUnits.push({
                value: Math.round(minutes),
                units: "m"
            } as TimeUnit)
        }
        timeUnits.push({
            value: Math.round(seconds),
            units: "s"
        } as TimeUnit)
        if (isPrecise) {
            timeUnits.push({
                value: Math.round(miliseconds),
                units: "ms"
            } as TimeUnit)
        }

        return timeUnits;
    }

    if (nanoseconds === 0) {
        return [{value: 0, units: "ns"} as TimeUnit];
    }

    const units = ["ns", "µs", "ms", "s", "ks"];
    let value = nanoseconds;
    let unitIndex = 0;

    while (value >= 1000 && unitIndex < units.length - 1) {
        if (unitIndex >= 3) {
            let seconds = value % 60
            let minutes = (value / 60) % 60
            let hours = (value / 60 / 60) % 60
            let days = Math.round((value / 60 / 60 / 24) % 24)
            let years = Math.round((value / 60 / 60 / 24 / 365))
            let timeUnits = [] as TimeUnit[]
            if (years > 0) {
                timeUnits.push({
                    value: Math.round(years),
                    units: "y"
                } as TimeUnit)
            }
            if (years > 0) {
                timeUnits.push({
                    value: Math.round(days),
                    units: "d"
                } as TimeUnit)
            }
            if (years > 0) {
                timeUnits.push({
                    value: Math.round(hours),
                    units: "h"
                } as TimeUnit)
            }
            timeUnits.push({
                value: Math.round(minutes),
                units: "m"
            } as TimeUnit)
            timeUnits.push({
                value: Math.round(seconds),
                units: "s"
            } as TimeUnit)

            return timeUnits;

        }
        value /= 1000;
        unitIndex++;
    }
    return [{value: value, units: units[unitIndex]} as TimeUnit]

}