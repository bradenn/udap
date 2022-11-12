<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {Preference} from "@/preferences";
import {PreferenceTypes} from "@/types";
import core from "@/core";

const index = core.router()

let storedToken = new Preference(PreferenceTypes.Token).get()

if (storedToken === "unset") {
    index.push("/setup/controller")
} else {
    index.push("/terminal/home")
}

</script>

<template>
    <div>
        <router-view/>
    </div>
</template>