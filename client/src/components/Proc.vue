<script>

function map_range(value, low1, high1, low2, high2) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

export default {
  name: "Proc",
  data() {
    return {
      canvas: null,
      config: {
        min: 250,
        max: 6000,
      },
      colors: [
        "rgba(255, 69, 58,0.6)",
        "rgba(255,159,10,0.6)",
        "rgba(48, 209, 88,0.6)",
        "rgba(100,210,255,0.6)",
        "rgba(10, 132, 255, 0.6)",
        "rgba(191, 90, 242, 0.6)",
      ],
      cycle: 0,
      width: 0,
      height: 0,
      frequencies: Array(),
      interval: null,
    }
  },
  props: {
    small: Boolean,
  },
  computed: {
    timings: function () {
      return this.$root.timings.sort((a, b) => b.delta - a.delta)
    }
  },
  mounted() {
    let canvas = document.getElementById("canvas")
    let dpr = 1;

    canvas.style.width = canvas.width + 'px'
    canvas.width *= dpr
    this.width = canvas.width

    canvas.style.height = canvas.height + 'px'
    canvas.height *= dpr
    this.height = canvas.height

    let ctx = canvas.getContext("2d");
    ctx.translate(0.5, 0.5);
    ctx.scale(1, 1);
    this.canvas = ctx
    this.cycle = 0;

    this.interval = setInterval(() => {
      this.cycle = (this.cycle - (Math.PI * 2) / 60 / 30)
      this.cycle %= (Math.PI * 2)
      this.drawAll(ctx)
    }, 1000 / 60)
  },
  beforeUnmount() {
    clearTimeout(this.interval)
  },
  methods: {
    groupBy(xs, key) {
      return xs.reduce(function (rv, x) {
        (rv[x[key]] = rv[x[key]] || []).push(x);
        return rv;
      }, {});
    },
    drawBounds() {
      let ctx = this.canvas
      ctx.beginPath();
      ctx.rect(0, 0, this.width, this.height);
      ctx.stroke();
    },
    drawAll(ctx) {
      ctx.clearRect(0, 0, this.width, this.height)
      for (let a of this.timings) {
        //
        this.draw(ctx, 0, a.delta / 100, a.frequency, this.colors[this.timings.indexOf(a) % this.colors.length])
      }
    },
    draw(ctx, x, y, frequency, color) {
      if (frequency > this.config.max || frequency < this.config.min) return
      let width = this.width;
      let height = this.height;

      ctx.beginPath();
      ctx.lineWidth = map_range(frequency, this.config.min, this.config.max, 6, 6);
      ctx.shadowBlur = 4;
      ctx.shadowColor = `rgba(255, 255, 255, 0.25)`;
      ctx.strokeStyle = `rgba(255,255,255,${map_range(frequency, this.config.min, this.config.max, 0.0625, 0.5)})`;

      let amplitude = this.height;
      let f = 0

      while (f <= width) {
        let k = height / 2 // Vertical Offset
        let b = (100 / frequency)
        amplitude = this.height / map_range(frequency, this.config.min, this.config.max, 0, this.height / 3);

        y = amplitude * Math.sin(((f) / width + this.cycle / 2 + x) / (b)) + k;
        ctx.lineTo(f, y);
        f = f + 1;
      }
      ctx.stroke();
    }
  }
}

</script>

<template>
  <div>
    <div class="label-ys label-o2 position-absolute mt-3">{{ timings.length }} processes</div>
    <div class="frequencies">
      <canvas id="canvas" class="frequencies"></canvas>
    </div>
  </div>
</template>

<style lang="scss" scoped>

.frequencies {
  width: 8rem !important;
  height: 2rem !important;
  padding: 0;
  margin: 0;
}
</style>
