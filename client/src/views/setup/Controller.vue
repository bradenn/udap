<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {config} from "@/config"
import type {Controller, Task} from "@/types"
import {PreferenceTypes, TaskType} from "@/types"
import axios from "axios"
import {Preference} from "@/preferences"
import {onMounted, reactive, ref} from "vue"
import TaskManager from "@/components/task/TaskManager.vue";


const tasks = [
    {
        title: "Controller",
        description: "Which core should the terminal authenticate with?",
        type: TaskType.Radio,
        options: [
            {title: "Primary", description: "10.0.1.2", value: "10.0.1.2:3020"},
            {title: "Backup", description: "10.0.1.18", value: "10.0.1.18"},
            {title: "Development", description: "10.0.1.11", value: "10.0.1.11:3020"},
        ],
        value: "10.0.1.2:3020",
        preview: "10.0.1.2"
    },
    {
        title: "Activation Code",
        description: "Enter the terminals designated activation code",
        type: TaskType.Passcode,
        value: "",
        preview: "---"
    },
]


// The local state of the component
let state = reactive({
    selected: getController(),
    auto: getController(),
    controllers: config.controllers,
    tasks: tasks
});

// When the view is mounted, test the default controllers to see which are live
onMounted(() => {
    // Access each controller
    for (let controller of state.controllers) {
        // Test its connectivity and update it's internal state to reflect it
        testController(controller)
    }
})

// Verify a controller is up and running, update the controller if it is
function testController(controller: Controller) {

    // Send a get request to the heartbeat endpoint of the provided controller app
    axios.get(`http://${controller.address}/status`, {
        httpsAgent: {
            rejectUnauthorized: false
        }
    }).then(res => {
        // Set the controller's status to reflect the successful request
        controller.status = true
        // Select the current controller, this will ensure a live controller will be selected by default (if possible)
        suggestController(controller.address)
    }).catch(err => {
        // Set the controller to reflect its down status
        controller.status = false
    })

}

// Recommend a working app controller to the user, typically the production node
function suggestController(address: string) {
    // Update the state to reflect the suggested address
    state.auto = address
}

// Set the preferred controller
function setController(address: string) {
    // Store the preferred controller in localStorage
    new Preference(PreferenceTypes.Controller).set(address);
    // Update the local state to reflect the selection
    state.selected = address
}

// Get the controller stored in localStorage
function getController() {
    return new Preference(PreferenceTypes.Controller).get()
}

// Errors incurred during verification
let errorMessage = ref("")

function authenticate(controller: string, key: string) {
    // Generated endpoint registration url using the security code and controller address
    let url = `http://${controller}:3020/endpoints/register/${key}`
    // Make the request to the controller app
    axios.get(url).then(res => {
        // Set the token in localStorage
        new Preference(PreferenceTypes.Token).set(res.data.token)
        // Redirect the user to the authenticated portal
        window.location.href = "/terminal"
    }).catch(err => {
        // Notify of failures
        errorMessage.value = "Invalid security code. Try again."
        // Reset the input dialog
    })
}

function finish(tasks: Task[]) {
    state.tasks = tasks

    const controller = tasks.find(t => t.title === "Controller");
    if (!controller) return;

    const code = tasks.find(t => t.title === "Activation Code");
    if (!code) return;

    setController(controller.value)

    authenticate(controller.value, code.value)


}

function next() {
    window.location.href = "/#/setup/authentication"
}

</script>

<template>

    <div class="d-flex justify-content-center w-100" style="margin-top: 6.25%;">
        <div style="width: 28rem">
            <TaskManager :on-complete="finish" :tasks="state.tasks" title="Wow">

            </TaskManager>
        </div>
    </div>

</template>


<style scoped>
.border-fog {
    border-color: rgba(255, 255, 255, 0.25) !important;
}

.border-transparent {
    border-color: transparent !important;
}
</style>