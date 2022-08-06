// Copyright (c) 2022 Braden Nicholson

// Weather Routes
import WeatherApp from "./Weather.vue";
import Summary from "./Summary.vue";

const weatherRoutes = {
    path: '/terminal/weather',
    name: 'Weather',
    redirect: '/terminal/weather/summary',
    component: WeatherApp,
    icon: 'fa-cloud-sun',
    children: [
        {
            path: '/terminal/weather/summary',
            name: 'Summary',
            component: Summary,
        },
    ]
}

export default weatherRoutes