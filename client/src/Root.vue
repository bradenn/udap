<script>
import Element from "./components/Element.vue";
import Frame from "./components/Pane.vue";
import SimpleKeyboard from "./components/Keyboard.vue";
import Loading from "./components/Loading.vue";
import Context from "./components/Context.vue";

import Dock from "./components/Dock.vue";

function parseJwt(token) {
  let base64Url = token.split('.')[1];
  let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
  let jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
  }).join(''));

  return JSON.parse(jsonPayload);
}

export default {
  components: {Context, Loading, SimpleKeyboard, Element, Frame, Dock},
  data() {
    return {
      state: {
        keyboard: false,
        waiting: false,
        accepting: false,
        error: null,
        input: "",
        last: new Date(),
        history: []
      },
      preferences: {
        accent: "slate",
        background: "viridian",
        blur: 5,
        input: "mouse",
        padding: 4,
        scale: "1.25",
        theme: "dark"
      },
      config: {
        host: "localhost",
        port: 3020,
      },
      connection: {
        connected: false,
        connecting: false,
        websocket: undefined,
        error: {}
      },
      entities: [],
      attributes: [],
      devices: [],
      networks: [],
      endpoints: [],
      timings: [],
      session: {
        token: "unset",
        subscriptions: [],
        metadata: {
          endpoint: {},
          modules: [],
          entities: [],
          timings: [],
        }
      },
    }
  },
  watch: {
    preferences: {
      handler(a, b) {
        this.saveConfig()
      }, deep: true
    },
    config: {
      handler(a, b) {
        this.saveConfig()
      }, deep: true
    },
    session: {
      handler(a, b) {
        this.saveConfig()
      }, deep: true
    }
  },
  created() {
    this.isConfig()
    this.loadConfig()
    this.connect()
  },
  computed: {
    media: function () {
      let entity = this.entities.find(e => e.type === 'media')
      let attrs = this.attributes.filter(e => e.entity === entity.id)
      return {
        entity: entity,
        attributes: attrs,
      }
    },
    background() {
      return `/custom/${this.preferences.background || "viridian"}@4x.png`
    }
  },
  beforeUnmount() {
    this.disconnect()
  },
  methods: {

    saveConfig() {

      let obj = {};
      obj.preferences = this.preferences
      obj.config = this.config
      obj.session = this.session
      let str = JSON.stringify(obj)
      localStorage.setItem("context", str)
    },
    isConfig() {
      if (localStorage.getItem("context") === null) {
        this.saveConfig()
        this.$router.push("/")
      }
    },
    loadConfig() {
      let cnf = localStorage.getItem("context")

      let obj = JSON.parse(cnf);

      if (obj.preferences != null) this.preferences = obj.preferences
      if (obj.config != null) this.config = obj.config
      if (obj.session != null) this.session = obj.session
      let str = JSON.stringify(obj)
      localStorage.setItem("context", str)
    },
    request(target, operation, body) {
      this.connection.websocket.send(JSON.stringify({
            target: target,
            operation: operation,
            body: body
          }
      ));
    },
    requestId(target, operation, body, id) {
      this.connection.websocket.send(JSON.stringify({
            target: target,
            operation: operation,
            payload: JSON.stringify(body),
            id: id
          }
      ));
    },
    connect() {
      this.timings = []
      if (this.connection.connecting || this.connection.connected || this.session.token === "") return

      let host = `ws://${this.config.host}:${this.config.port}/socket/${this.session.token}`

      this.connection.websocket = new WebSocket(host)
      this.connection.connecting = true

      this.connection.websocket.onmessage = this.onMessage
      this.connection.websocket.onopen = this.onConnect
      this.connection.websocket.onclose = this.onClose
      this.connection.websocket.onerror = this.onError
    },
    disconnect() {
      this.connection.websocket.close()
      this.connection.connecting = false
      this.connection.connected = false
    },
    onError(event) {
      this.state.error = event
    },
    onConnect(event) {
      this.accepting = true
      this.connection.connecting = false
      this.connection.connected = true
      this.accepting = false
    },
    onMessage(event) {
      this.accepting = true
      let data = JSON.parse(event.data)
      if (!data) console.log("Invalid JSON received")
      this.state.last = new Date()
      switch (data.operation) {
        case "metadata":
          this.state.waiting = false
          this.session.metadata = data.body
          break
        case "entity":
          if (this.entities.find(e => e.id === data.body.id)) {
            this.entities = this.entities.map(a => a.id === data.body.id ? data.body : a)
          } else {
            this.entities.push(data.body)
          }
          break
        case "attribute":
          if (this.attributes.find(e => e.id === data.body.id)) {
            this.attributes = this.attributes.map(a => a.id === data.body.id ? data.body : a)
          } else {
            this.attributes.push(data.body)
          }
          break
        case "device":
          if (this.devices.find(e => e.id === data.body.id)) {
            this.devices = this.devices.map(a => a.id === data.body.id ? data.body : a)
          } else {
            this.devices.push(data.body)
          }
          break
        case "network":
          if (this.networks.find(e => e.id === data.body.id)) {
            this.networks = this.networks.map(a => a.id === data.body.id ? data.body : a)
          } else {
            this.networks.push(data.body)
          }
          break
        case "endpoint":
          if (this.endpoints.find(e => e.id === data.body.id)) {
            this.endpoints = this.endpoints.map(a => a.id === data.body.id ? data.body : a)
          } else {
            this.endpoints.push(data.body)
          }
          break
        case "timing":
          if (this.timings.find(e => e.pointer === data.body.pointer)) {
            this.timings = this.timings.map(a => a.pointer === data.body.pointer ? data.body : a)
          } else {
            this.timings.push(data.body)
          }
          break
        default:
          console.log(data);
      }
      this.accepting = false
    },
    onClose(event) {
      this.connection.connecting = false
      this.connection.connected = false
      setTimeout(this.connect, 1000)
    },
    rootClasses() {
      return `theme-${this.preferences.theme} ${this.preferences.input === 'touchscreen' ? 'input-touch' : ''} accent-${this.preferences.accent} blurs-${this.preferences.blur}`
    },
  }

}

</script>

<template>
  <div class="root" v-bind:class="rootClasses()">
    <img :src="background" alt="Background" class="backdrop"/>
    <router-view/>
  </div>
</template>

<style lang="scss">

</style>
