<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup></script>

<template>
    <div>
        <div class="d-flex justify-content-between">
            <div>Weather</div>
            <div>21deg</div>
        </div>
        <div class="d-flex flex-row justify-content-between px-4">
            <div v-for="a in Array(8)" :key="a" class="d-flex flex-column justify-content-center align-items-center">
                <div>
                    􀇃
                </div>
                <div>23*</div>
            </div>
        </div>
        <div>Lorem ipsum dolor sit amet, consectetur adipiscing elit. In faucibus quam vitae massa fringilla
            condimentum. Fusce ipsum massa, dapibus in lorem vitae, mollis tristique enim. Pellentesque nibh mi, rhoncus
            a est id, tristique rhoncus nisl. Vestibulum molestie odio nec arcu mollis accumsan et nec turpis. Aenean
            porta euismod convallis. Nunc nisi urna, aliquet eu quam quis, commodo ultricies elit. Quisque et enim nunc.
            Nullam ut luctus massa. Aenean a sollicitudin dolor. Aenean sit amet libero sagittis, pellentesque augue a,
            pretium nisi. Maecenas tempor sem quis arcu faucibus, eu ornare turpis rhoncus. Sed nec condimentum dolor,
            sed tristique magna. Aliquam erat volutpat.
        </div>
    </div>

</template>

<style scoped>

</style>