<template>
  <div class="object-scale-down h-24 w-48">
    <img :src="HellasAIOLogo" alt=""/>
    <p class="text-center font-semibold antialiased text-3xl pt-4">
      {{ title }}
    </p>
  </div>
</template>

<script>
import HellasAIOLogo from "@/assets/logo.png";

export default {
  name: "LogoImport",
  props: ['title'],
  setup() {
    return {
      HellasAIOLogo
    };
  },
}
</script>