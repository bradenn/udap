export interface Weather {
    hourly: Hourly;
    generationtime_ms: number;
    daily_units: DailyUnits;
    longitude: number;
    current_weather: CurrentWeather;
    hourly_units: HourlyUnits;
    latitude: number;
    elevation: number;
    utc_offset_seconds: number;
    daily: Daily;
}

export interface Hourly {
    precipitation: number[];
    time: number[];
    relativehumidity_2m: number[];
    weathercode: number[];
    temperature_2m: number[];
}

export interface DailyUnits {
    sunset: string;
    time: string;
    temperature_2m_max: string;
    weathercode: string;
    temperature_2m_min: string;
    precipitation_sum: string;
    sunrise: string;
}

export interface CurrentWeather {
    temperature: number;
    time: number;
    winddirection: number;
    windspeed: number;
    weathercode: number;
}

export interface HourlyUnits {
    precipitation: string;
    time: string;
    relativehumidity_2m: string;
    weathercode: string;
    temperature_2m: string;
}

export interface Daily {
    sunset: number[];
    time: number[];
    temperature_2m_max: number[];
    weathercode: number[];
    temperature_2m_min: number[];
    precipitation_sum: number[];
    sunrise: number[];
}


interface WeatherCode {
    key: number,
    short: string,
    night: string,
    day: string
}

export let weatherCodes: WeatherCode[] = [
    {key: 0, night: '􀆮', day: '􀇁', short: 'Clear'},
    {key: 1, night: '􀇕', day: '􀇁', short: 'Mostly Clear'},
    {key: 2, night: '􀇕', day: '􀇛', short: 'Partly Cloudy'},
    {key: 3, night: '􀇃', day: '􀇃', short: 'Overcast'},
    {key: 45, night: '􀇋', day: '􀇋', short: 'Fog'},
    {key: 48, night: '􀇋', day: '􀇋', short: 'Depositing Rime Fog'},
    {key: 51, night: '􀇗', day: '􀇝', short: 'Slight Drizzle'},
    {key: 53, night: '􀇗', day: '􀇝', short: 'Moderate Drizzle'},
    {key: 55, night: '􀇗', day: '􀇝', short: 'Heavy Drizzle'},
    {key: 56, night: '􀇗', day: '􀇝', short: 'Light Frozen Drizzle'},
    {key: 57, night: '􀇗', day: '􀇝', short: 'Dense Frozen Drizzle'},
    {key: 61, night: '􀇅', day: '􀇅', short: 'Slight Rain'},
    {key: 63, night: '􀇇', day: '􀇇', short: 'Moderate Rain'},
    {key: 65, night: '􀇉', day: '􀇉', short: 'Heavy Rain'},
    {key: 71, night: '􀇏', day: '􀇏', short: 'Light Snow'},
    {key: 73, night: '􀇥', day: '􀇥', short: 'Moderate Snow'},
    {key: 75, night: '􀇥', day: '􀇥', short: 'Heavy Snow'},
    {key: 80, night: '􀇅', day: '􀇅', short: 'Light Rain'},
    {key: 81, night: '􀇇', day: '􀇇', short: 'Moderate Rain'},
    {key: 82, night: '􀇉', day: '􀇉', short: 'Heavy Rain'},
    {key: 85, night: '􀇏', day: '􀇏', short: 'Light Snow Showers'},
    {key: 86, night: '􀇏', day: '􀇏', short: 'Heavy Snow Showers'},
    {key: 95, night: '􀇟', day: '􀇟', short: 'Moderate Thunderstorms'}
]

function getWeatherState(code: number) {
    let meta = weatherCodes.find(w => w.key === code)
    if (!meta) return ""
    return meta.short
}

function getWeatherIcon(code: number, index: number) {
    let next = index % 24
    let meta = weatherCodes.find(w => w.key === code)
    if (!meta) return
    if (next >= 6 && next >= 20) {
        return meta.day
    }
    return meta.night
}

export {
    getWeatherIcon,
    getWeatherState
}
