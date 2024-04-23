<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import List from "udap-ui/components/List.vue";
</script>

<template>
  <div>
    <Element>
      <List scroll-y style="height: 78vh">
        <List>
          <ElementHeader title="Prompt"></ElementHeader>
          <Element foreground>
            There is a ray consisting of two points in a 3d volume. One point of the ray is fixed at the north-most pole
            of a sphere. The second point can be at any point on a plane which bisects the sphere. The point on the
            plane
            bisecting the sphere is represented in a 2d space as its third dimension is constrained to the plane. The
            ray
            can pivot from the polar point around the y axis, and around its own x axis. This is a pan and tilt, or
            azimuth and elevation.
          </Element>
        </List>
        <List>
          <ElementHeader title="Spherical Coordinate System"></ElementHeader>
          <Element foreground>
            <img alt="example" class="invert" src="https://static.udap.app/sphere.svg"/>
          </Element>
          <Element foreground>
            Use the above figure, translated to the northernmost pole of the sphere for your analysis.
          </Element>
        </List>
        <List>
          <ElementHeader title="Problem Diagram"></ElementHeader>
          <Element foreground>
            <img alt="example" class="invert-full" src="https://static.udap.app/diagram-dark.png"/>
          </Element>
          <Element foreground>
            This figure demonstrates the implementation of the problem.
          </Element>
        </List>
        <List>
          <ElementHeader title="Problem Parameters"></ElementHeader>
          <Element foreground>
            1. Create a function <span>f(x, y)</span> which takes in a 2d point on the plane bisecting the sphere. Have
            the function return a tuple (&theta;, &phi;) representing the spherical coordinates.
          </Element>
          <Element foreground>
            2. Create an additional function <span>f(&theta;, &phi;, h)</span> which takes in the same spherical
            coordinates in addition to the height from the pole to the center of the sphere, or simply the radius. Have
            this function return a point (x, y) on the plane.
          </Element>
        </List>
        <List>
          <ElementHeader title="Bonus Question"></ElementHeader>
          <Element foreground>
            3. Moving the target on the bisecting plane in a straight line will require several non-linear parallel
            movements. What considerations should be made in order to smoothly move between two 2D points on the plane?
            Is there an equation that can interpolate the the path from point A to point B on the plane given N discrete
            movements?
            <br>
            <small>Use the following material to inform your thought process.</small>
          </Element>
        </List>
        <List>
          <ElementHeader title="Interpolation"></ElementHeader>
          <Element foreground>
            Consider the equation:
            <pre>(x-min)/(max-min)*(tMax-tMin)+tMin</pre>
            This equation is part of a well known function known as <i>lerp</i>, short for linear interpolation. The
            function takes in a value, and its expected range, then normalizes the value and converts it to another
            range between tMin and tMax.
          </Element>
          <Element foreground>
            Linear interpolation can be used to move between two points by running lerp on the points Ax and Bx and Ay
            and By. Using x as the current time on an interval 0 (min) and T (max). And tMin and tMax being Ax and Bx
            respectively. This results in a jumpy path between the two points.
          </Element>
        </List>
      </List>

    </Element>
  </div>
</template>

<style lang="scss" scoped>
.invert { /* svg on an img tag */
  -webkit-filter: invert(.75) brightness(120%); /* safari 6.0 - 9.0 */
  filter: invert(.75) brightness(100%);
}

.invert-full { /* svg on an img tag */


  background: transparent;
  //filter: contrast(1);/**/

}
</style>